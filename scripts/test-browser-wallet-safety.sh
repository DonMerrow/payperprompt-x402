#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NODE_BIN="${NODE_BIN:-node}"

"$NODE_BIN" --check "$PROJECT_DIR/web/app.js"

grep -Fq '["io.metamask"' "$PROJECT_DIR/web/app.js"
grep -Fq '["com.coinbase.wallet"' "$PROJECT_DIR/web/app.js"
grep -Fq '["io.rabby"' "$PROJECT_DIR/web/app.js"
grep -Fq 'eth_signTypedData_v4' "$PROJECT_DIR/web/app.js"
grep -Fq 'wallet_switchEthereumChain' "$PROJECT_DIR/web/app.js"
grep -Fq 'PAYMENT-SIGNATURE' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'Only the configured disposable Base Sepolia payer wallet' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq '/api/official/plan' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq '/api/official/plan-jobs' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq '/api/official/plan-jobs' "$PROJECT_DIR/web/app.js"
grep -Fq 'returned non-JSON HTTP' "$PROJECT_DIR/web/app.js"
grep -Fq '/api/ai/work-suggestion' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'officialServiceForStrategy' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'selectedAmountAtomic' "$PROJECT_DIR/web/app.js"
grep -Fq '/api/official/reconcile-browser-wallet' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'appendOfficialHistoryOnce' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'verifyOfficialProofPayloadWithRust' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'Ollama did not answer. Wallet payment is disabled' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'task_type: preparedOfficialPlan.work_order.task_type' "$PROJECT_DIR/web/app.js"
grep -Fq 'work_completed' "$PROJECT_DIR/real-x402-go/cmd/official-server/main.go"
grep -Fq 'Complete the user'\''s authorized task' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'smart-contract-audit' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'Never deploy or execute transactions' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'It cannot deploy contracts, request private keys, sign transactions, move assets' "$PROJECT_DIR/web/index.html"
grep -Fq 'validateWorkProduct' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'Solidity elements covered before settlement' "$PROJECT_DIR/web/index.html"
grep -Fq 'grounded-work-v7' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'source-identifier consistency' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'source-accounting consistency' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'Deterministic Go source findings' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'reused_ready_job' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'num_predict' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq '/api/official/work-audit' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'pre-settlement work quality gate rejected' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'Deterministic semantic checks passed before settlement' "$PROJECT_DIR/web/index.html"
grep -Fq 'strategyForPaidService' "$PROJECT_DIR/real-x402-go/cmd/official-server/main.go"
grep -Fq 'prepared_work_id' "$PROJECT_DIR/web/app.js"
grep -Fq 'deliverable_commitment_sha256' "$PROJECT_DIR/web/app.js"
grep -Fq 'Prepared work already entered settlement' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'prepared_work_released' "$PROJECT_DIR/real-x402-go/cmd/official-server/main.go"

if grep -Eqi 'seed[ _-]?phrase[\"[:space:]]*:' "$PROJECT_DIR/web/app.js"; then
  echo "Unsafe seed-phrase input detected."
  exit 1
fi

if grep -Eqi 'private[ _-]?key[\"[:space:]]*:' "$PROJECT_DIR/web/app.js"; then
  echo "Unsafe private-key input detected."
  exit 1
fi

echo "BROWSER WALLET SAFETY TEST PASSED"
echo "Trusted provider allowlist, free AI ideas, prepared-work commitment, one-time release, exact pricing, Base Sepolia signing, Go validation, live-proof reconciliation, and Rust verification are present."
echo "No seed-phrase or private-key input contract was found."
