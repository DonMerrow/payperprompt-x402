#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

(cd "$PROJECT_DIR/go-core" && go test ./cmd/server)
(cd "$PROJECT_DIR/real-x402-go" && go test ./internal/payperprompt ./cmd/official-server)

grep -Fq 'POST /api/work-preflight' "$PROJECT_DIR/real-x402-go/cmd/official-server/main.go"
grep -Fq 'WorkProductCommitment' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'work_canonical_base64' "$PROJECT_DIR/real-x402-go/cmd/official-server/main.go"
grep -Fq 'decode canonical prepared work' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'solidityTestCoverageEvidence' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'No charge · preparation failed' "$PROJECT_DIR/web/app.js"
grep -Fq 'REJECTED DRAFT FROM THE PREVIOUS ATTEMPT' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'attempt < 3' "$PROJECT_DIR/real-x402-go/internal/payperprompt/work.go"
grep -Fq 'prepared_work_id' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'Prepared work already entered settlement' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'deliverable_commitment_sha256' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'prepared_work_released' "$PROJECT_DIR/web/app.js"
grep -Fq 'prepared-work-preview' "$PROJECT_DIR/web/index.html"
grep -Fq 'This one-time commitment is bound to the exact prompt, task, route, wallet, and payment.' "$PROJECT_DIR/web/index.html"

echo "PREPARED-WORK ESCROW TEST PASSED"
echo "AI work is completed and committed before wallet approval."
echo "Prompt, task, route, wallet, expiry, settlement state, and one-time release are enforced."
