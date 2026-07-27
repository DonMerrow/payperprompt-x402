#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

(cd "$PROJECT_DIR/go-core" && go test ./...)
(cd "$PROJECT_DIR/real-x402-go" && ./scripts/prepare-official.sh)
(cd "$PROJECT_DIR/real-x402-go" && go test ./...)
(cd "$PROJECT_DIR/sdk/go" && go test ./...)
(cd "$PROJECT_DIR/receipt-verifier-rust" && cargo test)

node --check "$PROJECT_DIR/web/app.js"
bash "$PROJECT_DIR/scripts/test-work-type-suggestions.sh"
bash "$PROJECT_DIR/scripts/test-browser-wallet-safety.sh"

grep -Fq 'UI v0.41.0.1' "$PROJECT_DIR/web/index.html"
grep -Fq 'proof-before-payment AI commerce' "$PROJECT_DIR/submission/DEVPOST-PASTE.md"
grep -Fq '/api/official/plan-jobs' "$PROJECT_DIR/web/app.js"
grep -Fq 'SafeMath is not a reentrancy guard' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'applyDeterministicSolidityReviewFindings' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'removeKnownSolidityReviewContradictions' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'submission expectations proof matrix' "$PROJECT_DIR/web/index.html"
grep -Fq 'toggle-work-audit-history' "$PROJECT_DIR/web/index.html"
grep -Fq 'history=${workAuditHistoryVisible}' "$PROJECT_DIR/web/app.js"
grep -Fq 'vm.deposit is not a Foundry cheatcode' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'SafeMath and reentrancy consistency' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'balance accounting coverage' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'SOURCE-DERIVED SOLIDITY REVIEW OBLIGATIONS' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'sdk/go' "$PROJECT_DIR/docs/hackathon-coverage.md"
grep -Fq '0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5' "$PROJECT_DIR/docs/official-x402-proof.md"
grep -Fq 'grounded-work-v7' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'MIT License' "$PROJECT_DIR/LICENSE"
test -s "$PROJECT_DIR/docs/submission-checklist.md"
test -s "$PROJECT_DIR/docs/WALLET-CONFIGURATION.md"
test -s "$PROJECT_DIR/docs/HOSTING-AND-JUDGE-AVAILABILITY.md"
test -s "$PROJECT_DIR/docs/PUBLIC-REPOSITORY.md"
test -s "$PROJECT_DIR/docs/index.html"
test -x "$PROJECT_DIR/scripts/build-submission-package.sh"
test -x "$PROJECT_DIR/scripts/configure-public-demo.sh"
grep -Fq '/api/config/public' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq '/api/config/public' "$PROJECT_DIR/web/app.js"
grep -Fq -- "--exclude '.env.local'" "$PROJECT_DIR/scripts/build-submission-package.sh"
grep -Fq '.env.local' "$PROJECT_DIR/.gitignore"
if grep -Fq '0x826154a3d58aeA3FBD2aa64aAD424594ade927eF' "$PROJECT_DIR/scripts/start-judge-stack.sh" ||
   grep -Fqi '0x826154a3d58aea3fbd2aa64aad424594ade927ef' "$PROJECT_DIR/web/app.js"; then
  echo "A historical payer address is still hardcoded into the configurable runtime."
  exit 1
fi

for fixture in "$PROJECT_DIR"/examples/large-contract-prompts/*.txt; do
  test -s "$fixture"
done

echo "SUBMISSION PROOF KIT TEST PASSED"
echo "Go, Rust, SDK, wallet safety, bounded audit history, submission mapping, and large-contract fixtures passed."
