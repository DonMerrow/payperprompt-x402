package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"payperprompt-x402-real/internal/payperprompt"
)

func TestSelectOfficialService(t *testing.T) {
	cases := []struct {
		strategy string
		routeID  string
		amount   string
	}{
		{"lowest-cost", "guardrail-economy", "10000"},
		{"lowest-latency", "guardrail-fast", "20000"},
		{"highest-quality", "guardrail-deep", "40000"},
	}
	for _, tc := range cases {
		selected := selectOfficialService(payperprompt.Analysis{Strategy: tc.strategy})
		if selected.RouteID != tc.routeID || selected.AmountAtomic != tc.amount {
			t.Fatalf("%s selected unexpected service: %+v", tc.strategy, selected)
		}
	}
}

func TestAtomicUSDC(t *testing.T) {
	if got := atomicUSDC("40000"); got != "0.04" {
		t.Fatalf("expected 0.04, got %s", got)
	}
}

func TestPaidSettlementFailureDoesNotAuthorizeReservationRelease(t *testing.T) {
	if shouldReleaseReservation(true) {
		t.Fatal("a confirmed settlement would release its reserved budget")
	}
	if !shouldReleaseReservation(false) {
		t.Fatal("a pre-settlement failure would leave reserved budget active")
	}
}

func TestProofPreservesPolicyAuthorizationForRecovery(t *testing.T) {
	tempDir := t.TempDir()
	proofPath := filepath.Join(tempDir, "official-settlement.json")
	t.Setenv("OFFICIAL_PROOF_HISTORY_PATH", filepath.Join(tempDir, "official-settlements.jsonl"))
	if err := saveProof(proofPath, officialProof{
		ProofVersion:          "payperprompt-official-v1",
		PolicyAuthorizationID: "auth-recovery-proof",
	}); err != nil {
		t.Fatal(err)
	}
	payload, err := os.ReadFile(proofPath)
	if err != nil {
		t.Fatal(err)
	}
	var saved map[string]any
	if err := json.Unmarshal(payload, &saved); err != nil {
		t.Fatal(err)
	}
	if saved["policy_authorization_id"] != "auth-recovery-proof" {
		t.Fatalf("proof lost its policy authorization: %v", saved)
	}
}

func TestVerifySettlementOnChainMatchesSelectedAmount(t *testing.T) {
	rpc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"jsonrpc":"2.0",
			"id":1,
			"result":{
				"status":"0x1",
				"logs":[{
					"address":"0x036CbD53842c5426634e7929541eC2318f3dCF7e",
					"topics":[
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000826154a3d58aea3fbd2aa64aad424594ade927ef",
						"0x00000000000000000000000007fb6cdd24cf265f8ea01a323708db34d6dbb630"
					],
					"data":"0x9c40"
				}]
			}
		}`))
	}))
	defer rpc.Close()
	t.Setenv("BASE_SEPOLIA_RPC_URL", rpc.URL)

	proof := officialProof{
		Payer: "0x826154a3d58aeA3FBD2aa64aAD424594ade927eF",
		Merchant: "0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630",
		Asset: baseSepoliaUSDC,
		AmountAtomic: "40000",
		Settlement: settlement{Transaction: "0xtest"},
	}
	verification := verifySettlementOnChain(context.Background(), proof)
	if !verification.Valid || !verification.USDCTransferMatched {
		t.Fatalf("expected valid enhanced-route transfer: %+v", verification)
	}
}

func TestVerifySettlementOnChainRejectsWrongAmount(t *testing.T) {
	rpc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"jsonrpc":"2.0",
			"id":1,
			"result":{
				"status":"0x1",
				"logs":[{
					"address":"0x036CbD53842c5426634e7929541eC2318f3dCF7e",
					"topics":[
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						"0x000000000000000000000000826154a3d58aea3fbd2aa64aad424594ade927ef",
						"0x00000000000000000000000007fb6cdd24cf265f8ea01a323708db34d6dbb630"
					],
					"data":"0x2710"
				}]
			}
		}`))
	}))
	defer rpc.Close()
	t.Setenv("BASE_SEPOLIA_RPC_URL", rpc.URL)

	proof := officialProof{
		Payer: "0x826154a3d58aeA3FBD2aa64aAD424594ade927eF",
		Merchant: "0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630",
		Asset: baseSepoliaUSDC,
		AmountAtomic: "40000",
		Settlement: settlement{Transaction: "0xtest"},
	}
	verification := verifySettlementOnChain(context.Background(), proof)
	if verification.Valid || verification.USDCTransferMatched {
		t.Fatalf("wrong transfer amount was accepted: %+v", verification)
	}
}
