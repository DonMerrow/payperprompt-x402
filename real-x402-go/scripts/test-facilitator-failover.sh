#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

go test ./internal/facilitatorpool ./cmd/official-server ./cmd/official-client

echo
echo "FACILITATOR FAILOVER SAFETY TEST PASSED"
echo "Supported/verify operations can use the next healthy endpoint."
echo "Settlement is attempted once; ambiguous outcomes require reconciliation."
echo "No private key was read. No payment was signed or sent."
