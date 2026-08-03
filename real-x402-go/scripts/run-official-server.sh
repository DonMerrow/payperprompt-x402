#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

if [[ -z "${MERCHANT_ADDRESS:-}" ]]; then
  echo "MERCHANT_ADDRESS is missing."
  echo "Set it to the PUBLIC address of the receiving test wallet."
  echo "Do not use the payer wallet."
  exit 1
fi

export X402_NETWORK="${X402_NETWORK:-eip155:84532}"
export X402_FACILITATOR_URL="${X402_FACILITATOR_URL:-https://x402.org/facilitator}"
# Optional ordered redundancy:
# export X402_FACILITATOR_URLS="https://primary.example,https://backup.example"
export OLLAMA_URL="${OLLAMA_URL:-http://127.0.0.1:11434}"
export OLLAMA_MODEL="${OLLAMA_MODEL:-qwen3-coder:30b}"
unset EVM_PRIVATE_KEY

go run ./cmd/official-server
