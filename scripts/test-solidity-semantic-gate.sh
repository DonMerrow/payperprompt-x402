#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$PROJECT_DIR/real-x402-go"
go test ./internal/payperprompt -run 'TestSolidity|TestWorkerRetriesRejectedSolidityDraft'

echo
echo "SOLIDITY SEMANTIC GATE TEST PASSED"
echo "Contradictory access-control, transfer/reentrancy, destination, receive-flow, and Solidity-syntax claims were rejected before settlement."
