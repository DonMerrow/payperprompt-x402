#!/usr/bin/env bash
set -euo pipefail

GO_CORE_URL="${GO_CORE_URL:-http://127.0.0.1:8084}"
RESPONSE_FILE="$(mktemp)"
trap 'rm -f "$RESPONSE_FILE"' EXIT

HTTP_STATUS="$(curl -sS \
  -o "$RESPONSE_FILE" \
  -w "%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{}' \
  "$GO_CORE_URL/api/official/reconcile-browser-wallet")"

python3 -m json.tool <"$RESPONSE_FILE"

if [[ "$HTTP_STATUS" -lt 200 || "$HTTP_STATUS" -ge 300 ]]; then
  echo
  echo "Reconciliation did not complete (HTTP $HTTP_STATUS)."
  echo "If the local ledger row is missing, use recover-browser-wallet-proof.sh with the known transaction hash."
  exit 1
fi

echo
echo "BROWSER WALLET PROOF RECONCILIATION COMPLETE"
echo "No wallet key was read, no signature was requested, and no payment was sent."
