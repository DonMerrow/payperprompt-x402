#!/usr/bin/env bash
set -euo pipefail

GO_CORE_URL="${GO_CORE_URL:-http://127.0.0.1:8084}"
TRANSACTION="${1:-${BROWSER_SETTLEMENT_TRANSACTION:-}}"
PAYER_ADDRESS="${PAYER_ADDRESS:-}"
if [[ ! "$PAYER_ADDRESS" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
  echo "PAYER_ADDRESS is required. Use the public disposable testnet address only."
  exit 1
fi

if [[ -z "$TRANSACTION" ]]; then
  echo "Usage: $0 0xTRANSACTION_HASH"
  echo "This recovery reads public Base Sepolia evidence only. It never requests a wallet key or signature."
  exit 1
fi

python3 - "$TRANSACTION" "$PAYER_ADDRESS" <<'PY' |
import json
import sys

print(json.dumps({
    "transaction": sys.argv[1],
    "payer": sys.argv[2],
    "amount_usd": "0.01",
}))
PY
  curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    --data-binary @- \
    "$GO_CORE_URL/api/official/reconcile-browser-wallet" |
  python3 -m json.tool

echo
echo "BROWSER WALLET PROOF RECOVERY COMPLETE"
echo "The existing transaction was verified publicly; no key was read, no signature was requested, and no payment was sent."
