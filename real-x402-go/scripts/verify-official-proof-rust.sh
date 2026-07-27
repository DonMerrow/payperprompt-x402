#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROOF_FILE="${OFFICIAL_PROOF_FILE:-$PROJECT_DIR/proof/official-settlement.json}"
RUST_URL="${RUST_OFFICIAL_PROOF_URL:-http://127.0.0.1:8085/api/verify-official-proof}"

if [[ ! -f "$PROOF_FILE" ]]; then
  echo "Official proof file not found: $PROOF_FILE"
  exit 1
fi

curl -fsS -X POST "$RUST_URL" \
  -H "content-type: application/json" \
  --data-binary "@$PROOF_FILE" | python3 -m json.tool
