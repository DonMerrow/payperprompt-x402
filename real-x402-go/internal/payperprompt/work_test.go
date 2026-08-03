package payperprompt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeTaskType(t *testing.T) {
	cases := []struct {
		input string
		text  string
		want  string
	}{
		{"code-review", "anything", "code-review"},
		{"auto", "Please turn these meeting notes into action items.", "meeting-actions"},
		{"auto", "Review this stack trace and bug report.", "bug-summary"},
		{"auto", "Audit this smart contract for reentrancy.", "smart-contract-audit"},
		{"auto", "Generate a smart contract for milestone escrow.", "smart-contract-generate"},
		{"auto", "Rewrite this x402 micropayment description.", "general-assistant"},
		{"auto", "Explain why the sky is blue.", "general-assistant"},
	}
	for _, testCase := range cases {
		if got := NormalizeTaskType(testCase.input, testCase.text); got != testCase.want {
			t.Fatalf("NormalizeTaskType(%q, %q) = %q, want %q", testCase.input, testCase.text, got, testCase.want)
		}
	}
}

func TestWorkerReturnsCompletedDeliverable(t *testing.T) {
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			http.NotFound(w, r)
			return
		}
		content, _ := json.Marshal(WorkProduct{
			TaskType: "meeting-actions", Title: "Meeting Action Plan",
			Summary: "Three decisions were recorded.",
			Deliverable: strings.Repeat(
				"## Decisions\n\nShip the release after validation and document the result. ",
				4,
			),
			ActionItems: []string{"Run the release tests."},
			Caveats:     []string{"No owner was named."},
		})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	worker := NewWorker(ollama.URL, "test-model")
	product, err := worker.Complete(
		context.Background(),
		"meeting-actions",
		"Turn these meeting notes into action items.",
		"standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if product.TaskType != "meeting-actions" ||
		product.Deliverable == "" ||
		len(product.ActionItems) != 1 {
		t.Fatalf("unexpected work product: %+v", product)
	}
}

func TestWorkerRetriesTruncatedJSON(t *testing.T) {
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		content := `{"task_type":"general-assistant","title":"Truncated"`
		if attempts > 1 {
			product, _ := json.Marshal(WorkProduct{
				TaskType: "general-assistant", Title: "Complete response",
				Summary: "A grounded completed response.",
				Deliverable: strings.Repeat(
					"This finished deliverable directly addresses the supplied request without claiming external actions. ",
					3,
				),
				ActionItems: []string{}, Caveats: []string{}, Coverage: []string{"requested response"},
			})
			content = string(product)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": content},
		})
	}))
	defer ollama.Close()

	worker := NewWorker(ollama.URL, "test-model")
	product, err := worker.Complete(context.Background(), "general-assistant", "Write a useful response.", "standard")
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 || product.Title != "Complete response" {
		t.Fatalf("expected one JSON retry, attempts=%d product=%+v", attempts, product)
	}
}

func TestWorkerNormalizesObjectCoverageFromOllama(t *testing.T) {
	source := `Explain this Solidity contract.
pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		content, _ := json.Marshal(map[string]any{
			"task_type": "smart-contract-explain",
			"title":     "SimpleVault explanation",
			"summary":   "The deployer owns withdrawals and the vault accepts ETH.",
			"deliverable": strings.Repeat(
				"Constructor sets owner. Receive accepts Ether into the contract balance. Withdraw checks msg.sender and sends Ether to owner. Transfer is gas-brittle but is not by itself proof of reentrancy. ", 4,
			),
			"action_items": []string{"Consider call with explicit failure handling."},
			"caveats":      []string{"The owner is fully trusted to withdraw vault funds."},
			"coverage": map[string]any{
				"constructor": "deployer ownership",
				"receive":     "Ether receipt",
				"withdraw":    "owner-only value flow",
			},
		})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	product, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(), "smart-contract-explain", source, "standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("object coverage should normalize without retry, attempts=%d", attempts)
	}
	joined := strings.Join(product.Coverage, ",")
	for _, expected := range []string{"constructor", "receive", "withdraw"} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("normalized coverage omitted %q: %v", expected, product.Coverage)
		}
	}
	if product.SemanticValidation == nil || !product.SemanticValidation.Valid {
		t.Fatalf("normalized work lacks semantic proof: %+v", product)
	}
}

func TestGroundingGateRejectsLanguageTokensAsCoverage(t *testing.T) {
	product := WorkProduct{
		TaskType: "code-review", Title: "Review", Summary: "Review summary.",
		Deliverable: strings.Repeat("The Transfer function needs synchronized access and rejects negative amounts. ", 4),
		Coverage:    []string{"Transfer", "make", "defer"},
	}
	if _, err := validateWorkProduct(product, "Review this Go function: func Transfer() {}"); err == nil ||
		!strings.Contains(err.Error(), "language token") {
		t.Fatalf("unsupported coverage was not rejected: %v", err)
	}
}

func TestGroundingGateRequiresRequestedArtifacts(t *testing.T) {
	product := WorkProduct{
		TaskType: "general-assistant", Title: "Architecture", Summary: "System design.",
		Deliverable: strings.Repeat("A cloud service coordinates teachers and students with secure tenancy. ", 4),
		Coverage:    []string{"architecture"},
	}
	request := "Provide a detailed architecture diagram and user interface mockups."
	if _, err := validateWorkProduct(product, request); err == nil ||
		!strings.Contains(err.Error(), "Mermaid") {
		t.Fatalf("missing requested diagram was not rejected: %v", err)
	}
}

func TestSolidityCodeReviewRejectsObservedUngroundedCoverage(t *testing.T) {
	source := `Review this Solidity contract.
pragma solidity ^0.8.0;
contract SharedWallet {
    address[] public admins;
    mapping(address => bool) public isAdmin;
    uint public balance;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint amount) public {
        require(amount <= balance, "Insufficient funds");
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function addAdmin(address admin) public { admins.push(admin); isAdmin[admin] = true; }
    function removeAdmin(address admin) public { delete isAdmin[admin]; }
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
    bool public paused;
}`
	product := WorkProduct{
		TaskType: "code-review", Title: "SharedWallet review",
		Summary: "A defensive review of the supplied contract.",
		Deliverable: strings.Repeat(
			"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause were reviewed for security and accounting behavior. ",
			3,
		),
		Coverage: []string{
			"withdraw",
			"removeAdmin",
			"reentrancy modifier from OpenZeppelin's SafeMath library",
			"Check for ongoing transactions",
		},
	}
	_, err := validateWorkProduct(product, source)
	if err == nil || !strings.Contains(err.Error(), "SafeMath") {
		t.Fatalf("observed invented coverage was accepted: %v", err)
	}
}

func TestSolidityCodeReviewRejectsObservedSemanticClaims(t *testing.T) {
	source := `Review this Solidity contract.
pragma solidity ^0.8.0;
contract SharedWallet {
    address[] public admins;
    mapping(address => bool) public isAdmin;
    uint public balance;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint amount) public {
        require(amount <= balance, "Insufficient funds");
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function addAdmin(address admin) public { admins.push(admin); isAdmin[admin] = true; }
    function removeAdmin(address admin) public { delete isAdmin[admin]; }
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
    bool public paused;
}`
	product := WorkProduct{
		TaskType: "code-review", Title: "SharedWallet review",
		Summary: "The review discusses every supplied function.",
		Deliverable: strings.Repeat(
			"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause were reviewed. "+
				"The contract is vulnerable to reentrancy due to the use of transfer. "+
				"Use a reentrancy modifier from OpenZeppelin's SafeMath library. ",
			3,
		),
		Coverage: []string{"deposit", "withdraw", "addAdmin", "removeAdmin", "pause", "unpause"},
	}
	validation, err := validateWorkProduct(product, source)
	if err == nil {
		t.Fatal("observed incorrect SharedWallet review was accepted")
	}
	for _, expected := range []string{
		"SafeMath and reentrancy consistency",
		"transfer and reentrancy consistency",
		"unrestricted access coverage",
		"pause enforcement coverage",
		"balance accounting coverage",
	} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("missing %q in semantic rejection: validation=%+v err=%v", expected, validation, err)
		}
	}
}

func TestSolidityCodeReviewAcceptsGroundedSharedWalletFindings(t *testing.T) {
	source := `Review this Solidity contract.
pragma solidity ^0.8.0;
contract SharedWallet {
    address[] public admins;
    mapping(address => bool) public isAdmin;
    uint public balance;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint amount) public {
        require(amount <= balance, "Insufficient funds");
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function addAdmin(address admin) public { admins.push(admin); isAdmin[admin] = true; }
    function removeAdmin(address admin) public { delete isAdmin[admin]; }
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
    bool public paused;
}`
	product := WorkProduct{
		TaskType: "code-review", Title: "Grounded SharedWallet review",
		Summary: "The contract has critical authorization and accounting defects.",
		Deliverable: strings.Repeat(
			"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause are covered. "+
				"Withdraw and every administration function are unrestricted because anyone can call them. "+
				"The paused flag is not enforced, so operations remain available while paused. "+
				"The tracked balance may diverge from address(this).balance through forced ETH. "+
				"Transfer is gas-brittle because of its stipend, but is not by itself proof of reentrancy. ",
			2,
		),
		Coverage: []string{"deposit", "withdraw", "addAdmin", "removeAdmin", "pause", "unpause"},
	}
	validation, err := validateWorkProduct(product, source)
	if err != nil {
		t.Fatalf("grounded SharedWallet review was rejected: validation=%+v err=%v", validation, err)
	}
	if !validation.Valid {
		t.Fatalf("grounded SharedWallet review was marked invalid: %+v", validation)
	}
}

func TestSemanticIntentMatchingAcceptsCorrectNegationsAndEquivalentPauseWording(t *testing.T) {
	safeMathCases := []struct {
		text string
		want bool
	}{
		{"Use a reentrancy modifier from OpenZeppelin's SafeMath library.", true},
		{"SafeMath provides reentrancy protection.", true},
		{"Use SafeMath as a reentrancy guard.", true},
		{"Use SafeMath to prevent reentrancy.", true},
		{"SafeMath is not a reentrancy guard.", false},
		{"Do not use SafeMath as reentrancy protection; use a dedicated guard.", false},
		{"Transfer is gas-brittle but is not proof of reentrancy. SafeMath is arithmetic protection and is not a reentrancy guard.", false},
		{"A dedicated reentrancy guard is appropriate. SafeMath only checks arithmetic.", false},
	}
	for _, testCase := range safeMathCases {
		if got := claimsSafeMathReentrancyProtection(strings.ToLower(testCase.text)); got != testCase.want {
			t.Fatalf("claimsSafeMathReentrancyProtection(%q)=%v want %v", testCase.text, got, testCase.want)
		}
	}
	for _, text := range []string{
		"The paused flag is never checked by deposit or withdraw.",
		"Pausing does not prevent deposits or withdrawals.",
		"Operations are still available while paused.",
		"The pause control is ineffective because protected operations ignore it.",
	} {
		if !reportsMissingPauseEnforcement(strings.ToLower(text)) {
			t.Fatalf("equivalent pause finding was not recognized: %q", text)
		}
	}
}

func TestDeterministicSolidityReviewFindingsRepairOmissions(t *testing.T) {
	source := `Perform a defensive review.
pragma solidity ^0.8.0;
contract SharedWallet {
    uint256 public balance;
    bool public paused;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint256 amount) public {
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
}`
	product := WorkProduct{
		TaskType: "code-review",
		Title:    "SharedWallet review",
		Summary:  "A defensive review of the submitted contract.",
		Deliverable: strings.Repeat(
			"The submitted contract requires a careful security review of its value flow and administrative behavior. ",
			3,
		),
		Coverage: []string{"security"},
	}
	grounded, err := ValidatePreparedWork(product, "code-review", source)
	if err != nil {
		t.Fatalf("deterministic source findings did not repair omissions: %v", err)
	}
	for _, expected := range []string{
		"### Deterministic Go source findings",
		"deposit, withdraw, pause, unpause",
		"unrestricted and any address can call them: withdraw, pause, unpause",
		"pausing does not block deposits",
		"address(this).balance",
		"SafeMath is not a reentrancy guard",
	} {
		if !strings.Contains(grounded.Deliverable, expected) {
			t.Fatalf("grounded review omitted %q: %s", expected, grounded.Deliverable)
		}
	}
	if strings.Join(grounded.Coverage, ",") != "deposit,withdraw,pause,unpause" {
		t.Fatalf("coverage was not source-derived: %v", grounded.Coverage)
	}
	firstCommitment, err := WorkProductCommitment(grounded)
	if err != nil {
		t.Fatal(err)
	}
	revalidated, err := ValidatePreparedWork(grounded, "code-review", source)
	if err != nil {
		t.Fatalf("grounded work was not release-safe: %v", err)
	}
	secondCommitment, err := WorkProductCommitment(revalidated)
	if err != nil {
		t.Fatal(err)
	}
	if firstCommitment != secondCommitment {
		t.Fatalf("release validation changed the prepared-work commitment: %s != %s", firstCommitment, secondCommitment)
	}
}

func TestDeterministicSolidityReviewFindingsRemoveKnownContradictions(t *testing.T) {
	source := `pragma solidity ^0.8.0;
contract SharedWallet {
    bool public paused;
    function withdraw(uint256 amount) public {}
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
}`
	product := WorkProduct{
		TaskType: "code-review",
		Title:    "Incorrect review",
		Summary:  "An incorrect defensive review.",
		Deliverable: strings.Repeat(
			"Access control is properly implemented, only admins can call withdraw, and the pause mechanism is enforced. ",
			3,
		),
	}
	corrected, err := ValidatePreparedWork(product, "code-review", source)
	if err != nil {
		t.Fatalf("known source contradiction was not deterministically corrected: %v", err)
	}
	for _, rejected := range []string{"access control is properly implemented", "only admins can call", "pause mechanism is enforced"} {
		if strings.Contains(strings.ToLower(corrected.Deliverable), rejected) {
			t.Fatalf("unsupported claim survived deterministic correction: %q in %s", rejected, corrected.Deliverable)
		}
	}
	if !strings.Contains(strings.Join(corrected.Caveats, "\n"), "Deterministic Go grounding removed") {
		t.Fatalf("correction disclosure was omitted: %+v", corrected)
	}
}

func TestWorkerIncludesSourceDerivedSolidityReviewObligationsOnFirstAttempt(t *testing.T) {
	source := `Review this Solidity contract.
pragma solidity ^0.8.0;
contract SharedWallet {
    uint256 public balance;
    bool public paused;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint256 amount) public {
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function addAdmin(address admin) public {}
    function removeAdmin(address admin) public {}
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
}`
	product := WorkProduct{
		TaskType: "code-review", Title: "Grounded SharedWallet review",
		Summary: "The contract has critical authorization and accounting defects.",
		Deliverable: strings.Repeat(
			"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause are covered. "+
				"Sensitive operations are unrestricted because anyone can call them. "+
				"The paused flag is not enforced. The tracked balance may diverge from address(this).balance through forced ETH. "+
				"Transfer is gas-brittle but is not by itself proof of reentrancy. ",
			2,
		),
		Coverage: []string{"deposit", "withdraw", "addAdmin", "removeAdmin", "pause", "unpause"},
	}
	var firstSystem string
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var requestPayload struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			t.Fatal(err)
		}
		firstSystem = requestPayload.Messages[0].Content
		content, _ := json.Marshal(product)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	result, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(),
		"code-review",
		source,
		"standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 || result.SemanticValidation == nil || !result.SemanticValidation.Valid {
		t.Fatalf("grounded first attempt did not pass: attempts=%d result=%+v", attempts, result)
	}
	for _, expected := range []string{
		"SOURCE-DERIVED SOLIDITY REVIEW OBLIGATIONS",
		"withdraw, addadmin, removeadmin, pause, unpause",
		"paused flag is changed but never enforced",
		"address(this).balance",
		"SafeMath is arithmetic protection",
	} {
		if !strings.Contains(firstSystem, expected) {
			t.Fatalf("first-attempt prompt omitted %q: %s", expected, firstSystem)
		}
	}
}

func TestWorkerGroundsKnownGeneralSolidityContradictionsWithoutRetry(t *testing.T) {
	source := `Review this Solidity contract.
pragma solidity ^0.8.0;
contract SharedWallet {
    uint256 public balance;
    bool public paused;
    function deposit() public payable { balance += msg.value; }
    function withdraw(uint256 amount) public {
        balance -= amount;
        payable(msg.sender).transfer(amount);
    }
    function addAdmin(address admin) public {}
    function removeAdmin(address admin) public {}
    function pause() public { paused = true; }
    function unpause() public { paused = false; }
}`
	validProduct := WorkProduct{
		TaskType: "code-review", Title: "Corrected SharedWallet review",
		Summary: "The corrected review covers authorization, pausing, accounting, and value transfer.",
		Deliverable: strings.Repeat(
			"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause are covered. "+
				"Sensitive operations are unrestricted because anyone can call them. "+
				"The paused flag is not enforced. The tracked balance may diverge from address(this).balance through forced ETH. "+
				"Transfer is gas-brittle but is not by itself proof of reentrancy. ",
			2,
		),
		Coverage: []string{"deposit", "withdraw", "addAdmin", "removeAdmin", "pause", "unpause"},
	}
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var requestPayload struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			t.Fatal(err)
		}
		product := validProduct
		if attempts == 1 {
			product.Deliverable = strings.Repeat(
				"Deposit, withdraw, addAdmin, removeAdmin, pause, and unpause were reviewed. "+
					"Access control is properly implemented, only admins can call withdraw, and the pause mechanism is enforced. ",
				3,
			)
			product.Deliverable +=
				"The contract uses address(this).balance to track its funds. " +
					"SafeMath is used for arithmetic protection. " +
					"The transfer can lead to reentrancy attacks."
			product.Caveats = []string{
				"The contract's use of SafeMath does not serve as a reentrancy guard.",
			}
		}
		content, _ := json.Marshal(product)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	result, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(),
		"code-review",
		source,
		"standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 || result.SemanticValidation == nil || !result.SemanticValidation.Valid {
		t.Fatalf("deterministic Solidity grounding did not pass on one attempt: attempts=%d result=%+v", attempts, result)
	}
	for _, rejected := range []string{
		"access control is properly implemented",
		"only admins can call",
		"pause mechanism is enforced",
		"uses address(this).balance",
		"safemath is used",
		"contract's use of safemath",
		"can lead to reentrancy attacks",
	} {
		if strings.Contains(strings.ToLower(result.Deliverable), rejected) {
			t.Fatalf("unsupported claim survived worker grounding: %q in %s", rejected, result.Deliverable)
		}
	}
	if strings.Contains(strings.ToLower(strings.Join(result.Caveats, "\n")), "contract's use of safemath") {
		t.Fatalf("invented SafeMath caveat survived worker grounding: %+v", result.Caveats)
	}
	if !strings.Contains(strings.Join(result.Caveats, "\n"), "Deterministic Go grounding removed") {
		t.Fatalf("worker omitted correction disclosure: %+v", result)
	}
}

func TestSmartContractGuidanceIsDefined(t *testing.T) {
	for _, taskType := range []string{
		"smart-contract-audit",
		"smart-contract-generate",
		"smart-contract-explain",
		"smart-contract-tests",
		"smart-contract-fix",
	} {
		if TaskLabel(taskType) == "AI Work Product" {
			t.Fatalf("missing Smart Contract Studio label for %s", taskType)
		}
		if taskGuidance(taskType) == "" {
			t.Fatalf("missing Smart Contract Studio guidance for %s", taskType)
		}
	}
}

func TestWorkerCorrectsClassTerminologyBeforeValidation(t *testing.T) {
	source := `Explain this Solidity contract.
pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		content, _ := json.Marshal(WorkProduct{
			TaskType: "smart-contract-explain",
			Title:    "SimpleVault explanation",
			Summary:  "The deployer controls class SimpleVault.",
			Deliverable: strings.Repeat(
				"Class SimpleVault has a constructor that sets owner. Receive accepts Ether into the contract balance. Withdraw checks msg.sender and transfers Ether to owner. Transfer is gas-brittle but is not by itself proof of reentrancy. ", 4,
			),
			ActionItems: []string{},
			Caveats:     []string{"The owner is trusted to withdraw funds."},
			Coverage:    []string{"constructor", "receive", "withdraw"},
		})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	product, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(), "smart-contract-explain", source, "standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("terminology correction should pass without retry, attempts=%d", attempts)
	}
	if strings.Contains(strings.ToLower(product.Deliverable), "class simplevault") {
		t.Fatalf("incorrect Solidity terminology survived: %s", product.Deliverable)
	}
	if !strings.Contains(strings.ToLower(product.Deliverable), "contract simplevault") {
		t.Fatalf("corrected contract terminology is missing: %s", product.Deliverable)
	}
	if !strings.Contains(strings.Join(product.Caveats, "\n"), "Deterministic Go normalization corrected") {
		t.Fatalf("correction disclosure is missing: %v", product.Caveats)
	}
	if product.SemanticValidation == nil || !product.SemanticValidation.Valid {
		t.Fatalf("corrected work lacks semantic proof: %+v", product)
	}
}

func TestSolidityExplanationRequiresCompleteCoverage(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    constructor() {}
    receive() external payable {}
    function withdraw(uint256 amount) external {}
}`
	incomplete := WorkProduct{
		TaskType:    "smart-contract-explain",
		Title:       "SimpleVault",
		Summary:     "Incomplete explanation.",
		Deliverable: strings.Repeat("Constructor explanation. ", 20),
		Coverage:    []string{"constructor"},
	}
	if _, err := validateWorkProduct(incomplete, source); err == nil ||
		!strings.Contains(err.Error(), "receive") ||
		!strings.Contains(err.Error(), "withdraw") {
		t.Fatalf("incomplete Solidity coverage was not rejected correctly: %v", err)
	}
	complete := WorkProduct{
		TaskType: "smart-contract-explain",
		Title:    "SimpleVault",
		Summary:  "Complete explanation.",
		Deliverable: strings.Repeat(
			"Constructor establishes ownership. Receive accepts ETH. Withdraw transfers the requested amount. ",
			6,
		),
		Coverage: []string{"constructor()", "receive()", "withdraw(uint256)"},
	}
	validation, err := validateWorkProduct(complete, source)
	if err != nil {
		t.Fatalf("complete Solidity coverage was rejected: %v", err)
	}
	if !validation.Valid || len(validation.Checks) < 2 {
		t.Fatalf("semantic validation evidence was not returned: %+v", validation)
	}
}

func TestFoundryCoverageUsesExecutableTestEvidence(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	deliverable := `pragma solidity ^0.8.20;
import {Test} from "forge-std/Test.sol";
import {SimpleVault} from "../src/SimpleVault.sol";
contract SimpleVaultTest is Test {
    SimpleVault vault;
    function setUp() public {
        vault = new SimpleVault();
        vm.deal(address(this), 10 ether);
    }
    function testAcceptsEther() public {
        (bool ok,) = address(vault).call{value: 1 ether}("");
        assertTrue(ok);
        assertEq(address(vault).balance, 1 ether);
    }
    function testOwnerWithdraws() public {
        (bool ok,) = address(vault).call{value: 1 ether}("");
        assertTrue(ok);
        vault.withdraw(1 ether);
        assertEq(address(vault).balance, 0);
    }
    function testNonOwnerCannotWithdraw() public {
        address stranger = address(0xBEEF);
        vm.prank(stranger);
        vm.expectRevert(bytes("not owner"));
        vault.withdraw(1);
    }
}`
	product := WorkProduct{
		TaskType: "smart-contract-tests", Title: "SimpleVault Foundry tests",
		Summary:     "Tests deployment, deposits, and owner withdrawals.",
		Deliverable: deliverable,
		Coverage:    []string{},
	}
	normalizeSolidityTestCoverage(&product, source)
	validation, err := validateWorkProduct(product, source)
	if err != nil {
		t.Fatalf("executable Foundry coverage was rejected: %v", err)
	}
	if !validation.Valid || strings.Join(product.Coverage, ",") != "constructor,receive,withdraw" {
		t.Fatalf("test-aware coverage was not normalized: product=%+v validation=%+v", product, validation)
	}
}

func TestFoundryCoverageRejectsMissingConstructionEvidence(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    constructor() {}
    function withdraw(uint256 amount) external {}
}`
	product := WorkProduct{
		TaskType: "smart-contract-tests", Title: "Incomplete tests",
		Summary: "Only withdrawal is exercised.",
		Deliverable: `pragma solidity ^0.8.20;
import {Test} from "forge-std/Test.sol";
import {SimpleVault} from "../src/SimpleVault.sol";
contract IncompleteTest is Test {` + strings.Repeat(
			`function testWithdraw() public { vault.withdraw(1); assertTrue(true); }`,
			6,
		) + `}`,
	}
	if _, err := validateWorkProduct(product, source); err == nil ||
		!strings.Contains(err.Error(), "constructor") {
		t.Fatalf("missing construction evidence was not rejected: %v", err)
	}
}

func TestFoundryValidityRejectsObservedNonCompilingOutput(t *testing.T) {
	source := `Generate a bounded fuzz test and confirm a non-owner reverts.
pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	deliverable := `pragma solidity ^0.8.20;
import "forge-std/Test.sol";
contract SimpleVaultTest is Test {
    SimpleVault simpleVault;
    function setUp() public { simpleVault = new SimpleVault(); }
    function testReceive() public { vm.deposit(address(this), 1 ether); simpleVault.receive(); }
    function testWithdraw() public { simpleVault.withdraw(1); assertEq(simpleVault.balance(), 0); }
    function testNonOwner() public { vm.expectRevert("not owner"); simpleVault.withdraw(1); }
    function testFuzz() public { uint256 amount = random(1, 10); simpleVault.withdraw(amount); }
}` + strings.Repeat("\n// invalid generated test evidence", 8)
	product := WorkProduct{
		TaskType: "smart-contract-tests", Title: "Invalid tests",
		Summary: "Observed non-compiling output.", Deliverable: deliverable,
		Coverage: []string{"constructor", "receive", "withdraw"},
	}
	_, err := validateWorkProduct(product, source)
	if err == nil {
		t.Fatal("known non-compiling Foundry output was accepted")
	}
	for _, expected := range []string{"vm.deposit", "receive cannot be called", "no balance() function", "random() is undefined", "does not import"} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("missing %q in validity failure: %v", expected, err)
		}
	}
}

func TestLargeContractCoverageScalesBeyondSimpleVault(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract TeamTreasury {
    constructor(address[3] memory approvers) {}
    receive() external payable {}
    fallback() external payable {}
    function propose(address to, uint256 value, bytes calldata data) external returns (uint256) {}
    function approve(uint256 id) external {}
    function revoke(uint256 id) external {}
    function cancel(uint256 id) external {}
    function execute(uint256 id) external {}
    function proposal(uint256 id) external view returns (address, uint256) {}
    function approvedBy(uint256 id, address signer) external view returns (bool) {}
    function quorum(uint256 id) external view returns (uint256) {}
    function expired(uint256 id) external view returns (bool) {}
}`
	elements := solidityCoverageElements(source)
	if len(elements) != 12 {
		t.Fatalf("large contract coverage found %d elements, want 12: %v", len(elements), elements)
	}
	product := WorkProduct{
		TaskType: "smart-contract-audit", Title: "Incomplete large audit",
		Summary: "One function was omitted.",
		Deliverable: strings.Repeat(
			"constructor receive fallback propose approve revoke cancel execute proposal approvedBy quorum are reviewed. ",
			4,
		),
		Coverage: []string{
			"constructor", "receive", "fallback", "propose", "approve", "revoke",
			"cancel", "execute", "proposal", "approvedBy", "quorum",
		},
	}
	if _, err := validateWorkProduct(product, source); err == nil ||
		!strings.Contains(err.Error(), "expired") {
		t.Fatalf("large-contract omitted function was not rejected: %v", err)
	}
}

func TestNormalizeSolidityTestBalanceAccessRepairsInventedGetter(t *testing.T) {
	source := `Generate Foundry tests.
pragma solidity ^0.8.20;
contract SimpleVault { receive() external payable {} }`
	product := WorkProduct{
		TaskType: "smart-contract-tests",
		Deliverable: `function testBalance() public {
    assertEq(vault.balance(), 1 ether);
}`,
	}
	normalizeSolidityTestBalanceAccess(&product, source)
	if strings.Contains(product.Deliverable, "vault.balance()") ||
		!strings.Contains(product.Deliverable, "address(vault).balance") {
		t.Fatalf("invented balance getter was not normalized: %s", product.Deliverable)
	}
	if len(product.Caveats) != 1 || !strings.Contains(product.Caveats[0], "address(contract).balance") {
		t.Fatalf("normalization was not disclosed: %v", product.Caveats)
	}
}

func TestNormalizeSolidityTestBalanceAccessPreservesDeclaredBalanceFunction(t *testing.T) {
	source := `Generate Foundry tests.
pragma solidity ^0.8.20;
contract Ledger { function balance() external view returns (uint256) { return 1; } }`
	product := WorkProduct{TaskType: "smart-contract-tests", Deliverable: "ledger.balance()"}
	normalizeSolidityTestBalanceAccess(&product, source)
	if product.Deliverable != "ledger.balance()" {
		t.Fatalf("declared balance function was incorrectly rewritten: %s", product.Deliverable)
	}
}

func TestSolidityTestGenerationInstructionRequiresExecutableFoundryPatterns(t *testing.T) {
	request := `Generate a complete Foundry test suite with a bounded fuzz test.
pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
    }
}`
	instruction := solidityTestGenerationInstruction(request, solidityCoverageElements(request))
	for _, required := range []string{
		`Import {Test} from "forge-std/Test.sol"`,
		`import {SimpleVault} from "../src/SimpleVault.sol";`,
		"vm.prank(nonOwner)",
		"vm.expectRevert",
		"bound(value, minimum, maximum)",
		"constructor, withdraw",
	} {
		if !strings.Contains(instruction, required) {
			t.Fatalf("generation instruction omitted %q: %s", required, instruction)
		}
	}
}

func TestSolidityTestGenerationInstructionDoesNotAffectHardhatSelection(t *testing.T) {
	request := `Generate Hardhat tests.
pragma solidity ^0.8.20;
contract Counter { function increment() external {} }`
	instruction := solidityTestGenerationInstruction(request, solidityCoverageElements(request))
	if !strings.Contains(instruction, "Hardhat") || !strings.Contains(instruction, "ethers deployment fixtures") {
		t.Fatalf("Hardhat instructions were not preserved: %s", instruction)
	}
	if strings.Contains(instruction, "forge-std") {
		t.Fatalf("Hardhat request received Foundry imports: %s", instruction)
	}
}

func TestSolidityContractPatternRequiresDeclarationBody(t *testing.T) {
	prose := "Create a Solidity contract that implements an ERC-721 token."
	if matches := solidityContractPattern.FindAllStringSubmatch(prose, -1); len(matches) != 0 {
		t.Fatalf("natural-language prose became a contract declaration: %v", matches)
	}
	source := "pragma solidity ^0.8.20; contract Token { }"
	matches := solidityContractPattern.FindAllStringSubmatch(source, -1)
	if len(matches) != 1 || matches[0][1] != "Token" {
		t.Fatalf("actual contract declaration was not recognized: %v", matches)
	}
}

func TestHardhatCoverageRecognizesDeployAndValueTransfer(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    constructor() {}
    receive() external payable {}
    function withdraw(uint256 amount) external {}
}`
	tests := `
const Vault = await ethers.getContractFactory("SimpleVault");
const vault = await Vault.deploy();
await owner.sendTransaction({to: await vault.getAddress(), value: ethers.parseEther("1")});
await vault.withdraw(1);
`
	for _, name := range []string{"constructor", "receive", "withdraw"} {
		if passed, _ := solidityTestCoverageEvidence(name, tests, source); !passed {
			t.Fatalf("Hardhat evidence did not prove %s coverage", name)
		}
	}
}

func TestSoliditySemanticValidatorRejectsContradictoryClaims(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	base := WorkProduct{
		TaskType: "smart-contract-explain",
		Title:    "SimpleVault explanation",
		Summary:  "The owner controls withdrawals.",
		Deliverable: strings.Repeat(
			"Constructor sets owner. Receive accepts Ether into the contract balance. Withdraw checks the owner and transfers Ether to owner. ",
			4,
		),
		Coverage: []string{"constructor", "receive", "withdraw"},
	}
	cases := []struct {
		name       string
		mutate     func(*WorkProduct)
		wantReason string
	}{
		{
			name: "denies existing access control",
			mutate: func(product *WorkProduct) {
				product.Caveats = []string{"The contract does not have any access control."}
			},
			wantReason: "access-control consistency",
		},
		{
			name: "mislabels transfer as reentrancy",
			mutate: func(product *WorkProduct) {
				product.Caveats = []string{"The contract has no protection against reentrancy."}
			},
			wantReason: "transfer and reentrancy consistency",
		},
		{
			name: "invents arbitrary recipient",
			mutate: func(product *WorkProduct) {
				product.ActionItems = []string{"The owner can choose the recipient."}
			},
			wantReason: "fixed-destination value flow",
		},
		{
			name: "claims receive forwards Ether",
			mutate: func(product *WorkProduct) {
				product.Deliverable += strings.Repeat(" The receive function forwards Ether elsewhere.", 12)
			},
			wantReason: "receive-function value flow",
		},
		{
			name: "rewrites Solidity contract as class",
			mutate: func(product *WorkProduct) {
				product.Deliverable += strings.Repeat(" class SimpleVault is declared here.", 15)
			},
			wantReason: "Solidity syntax consistency",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			product := base
			testCase.mutate(&product)
			validation, err := validateWorkProduct(product, source)
			if err == nil || !strings.Contains(err.Error(), testCase.wantReason) {
				t.Fatalf("contradiction was not rejected correctly: validation=%+v err=%v", validation, err)
			}
			failedCheckRecorded := false
			for _, check := range validation.Checks {
				if check.Name == testCase.wantReason && !check.Passed {
					failedCheckRecorded = true
				}
			}
			if validation.Valid || !failedCheckRecorded {
				t.Fatalf("failed semantic check was not recorded: %+v", validation)
			}
		})
	}
}

func TestSoliditySemanticValidatorRejectsObservedSimpleVaultFailure(t *testing.T) {
	source := `pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	product := WorkProduct{
		TaskType: "smart-contract-explain",
		Title:    "SimpleVault Contract Explanation",
		Summary:  "The owner controls withdrawals.",
		Deliverable: strings.Repeat(
			"The constructor sets the owner. The withdraw function transfers value. "+
				"The receive function forwards any received Ether to the contract's balance. "+
				"The contract does not have any protection against reentrancy. "+
				"The owner must not transfer Ether to unauthorized addresses. "+
				"class SimpleVault represents the contract. ",
			3,
		),
		Coverage: []string{"constructor", "receive", "withdraw"},
	}
	validation, err := validateWorkProduct(product, source)
	if err == nil {
		t.Fatal("the observed contradictory paid explanation was accepted")
	}
	for _, checkName := range []string{
		"access-control consistency",
		"transfer and reentrancy consistency",
		"fixed-destination value flow",
		"receive-function value flow",
		"Solidity syntax consistency",
	} {
		if !strings.Contains(err.Error(), checkName) {
			t.Fatalf("missing %s failure in %q", checkName, err)
		}
	}
	if validation.Valid {
		t.Fatalf("observed failure was marked valid: %+v", validation)
	}
}

func TestWorkerRetriesRejectedSolidityDraftBeforeReturningWork(t *testing.T) {
	source := `Explain this Solidity contract.
pragma solidity ^0.8.20;
contract SimpleVault {
    address public owner;
    constructor() { owner = msg.sender; }
    receive() external payable {}
    function withdraw(uint256 amount) external {
        require(msg.sender == owner, "not owner");
        payable(owner).transfer(amount);
    }
}`
	attempts := 0
	sawFreshCorrection := false
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		var requestPayload struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			t.Fatal(err)
		}
		if attempts > 1 && len(requestPayload.Messages) > 1 &&
			strings.Contains(requestPayload.Messages[0].Content, "Generate a fresh answer") &&
			strings.Contains(requestPayload.Messages[0].Content, "constructor, receive, withdraw") &&
			!strings.Contains(requestPayload.Messages[1].Content, "REJECTED DRAFT FROM THE PREVIOUS ATTEMPT") {
			sawFreshCorrection = true
		}
		deliverable := strings.Repeat(
			"Constructor sets owner. Receive accepts Ether into the contract balance. Withdraw checks owner and transfers to owner. ",
			4,
		)
		caveats := []string{"transfer has a limited gas stipend and can be gas-brittle."}
		if attempts == 1 {
			caveats = []string{"The contract does not have any access control."}
		}
		content, _ := json.Marshal(WorkProduct{
			TaskType:    "smart-contract-explain",
			Title:       "SimpleVault explanation",
			Summary:     "The deployer becomes owner and controls withdrawals.",
			Deliverable: deliverable,
			Caveats:     caveats,
			Coverage:    []string{"constructor", "receive", "withdraw"},
		})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	product, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(),
		"smart-contract-explain",
		source,
		"standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("expected one corrective retry, got %d attempts", attempts)
	}
	if !sawFreshCorrection {
		t.Fatal("corrective retry did not receive fresh source-grounded obligations")
	}
	if product.SemanticValidation == nil || !product.SemanticValidation.Valid {
		t.Fatalf("corrected work lacks semantic proof: %+v", product)
	}
}

func TestWorkerAllowsSecondCorrectiveAttemptForGeneratedTests(t *testing.T) {
	source := `Generate Foundry tests.
pragma solidity ^0.8.20;
contract SimpleVault {
    constructor() {}
    receive() external payable {}
    function withdraw(uint256 amount) external {}
}`
	attempts := 0
	ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		deliverable := strings.Repeat(
			"An incomplete draft that describes testing without executable contract calls. ",
			6,
		)
		if attempts == 3 {
			deliverable = `pragma solidity ^0.8.20;
import {Test} from "forge-std/Test.sol";
import {SimpleVault} from "../src/SimpleVault.sol";
contract SimpleVaultTest is Test {
    SimpleVault vault;
    function setUp() public { vault = new SimpleVault(); vm.deal(address(this), 2 ether); }
    function testReceive() public {
        (bool ok,) = address(vault).call{value: 1 ether}("");
        assertTrue(ok);
    }
    function testWithdraw() public {
        vault.withdraw(0);
        assertEq(address(vault).balance, 0);
    }
}`
		}
		content, _ := json.Marshal(WorkProduct{
			TaskType: "smart-contract-tests", Title: "SimpleVault tests",
			Summary: "Generated Foundry coverage.", Deliverable: deliverable,
		})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": string(content)},
		})
	}))
	defer ollama.Close()

	product, err := NewWorker(ollama.URL, "test-model").Complete(
		context.Background(),
		"smart-contract-tests",
		source,
		"standard",
	)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Fatalf("expected two corrective attempts, got %d total attempts", attempts)
	}
	if product.SemanticValidation == nil || !product.SemanticValidation.Valid {
		t.Fatalf("third-attempt test suite lacks semantic proof: %+v", product)
	}
}

func TestSolidityWorkCoverageInstructionRequiresExactElements(t *testing.T) {
	request := `Explain this contract.
pragma solidity ^0.8.20;
contract SimpleVault {
    constructor() {}
    receive() external payable {}
    function withdraw(uint256 amount) external {}
}`

	required := solidityCoverageElements(request)
	instruction := solidityWorkCoverageInstruction(
		"smart-contract-explain",
		request,
		required,
	)

	for _, element := range []string{"constructor", "receive", "withdraw"} {
		if !strings.Contains(instruction, element) {
			t.Fatalf("instruction omitted %q: %s", element, instruction)
		}
	}
	if !strings.Contains(instruction, "coverage JSON array") {
		t.Fatalf("instruction omitted exact coverage-array requirement: %s", instruction)
	}
}

func TestPreparedWorkCommitmentChangesWhenDeliverableChanges(t *testing.T) {
	product := WorkProduct{
		TaskType: "general-assistant",
		Title:    "Prepared work", Summary: "Ready before payment.",
		Deliverable: "The exact completed result.",
		ActionItems: []string{}, Caveats: []string{}, Coverage: []string{},
	}
	first, err := WorkProductCommitment(product)
	if err != nil {
		t.Fatal(err)
	}
	product.Deliverable = "An altered completed result."
	second, err := WorkProductCommitment(product)
	if err != nil {
		t.Fatal(err)
	}
	if first == second || len(first) != 64 || len(second) != 64 {
		t.Fatalf("prepared-work commitments are not content-bound: %q %q", first, second)
	}
}
