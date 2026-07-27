#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

if [[ -z "${EVM_PRIVATE_KEY:-}" ]]; then
  echo "EVM_PRIVATE_KEY is missing."
  echo "Load a test-only payer key into this terminal; never paste it into chat or source code."
  exit 1
fi

export SERVER_URL="${SERVER_URL:-http://127.0.0.1:8082/api/check-prompt}"
go run ./cmd/official-client pay
