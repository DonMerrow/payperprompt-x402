#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/policy-control.sh"

if [[ -z "${PAYER_ADDRESS:-}" ]]; then
  echo "PAYER_ADDRESS is required. This is the public testnet wallet address, not its private key."
  exit 1
fi

POLICY_URL="${OFFICIAL_POLICY_CONFIG_URL:-http://127.0.0.1:8084/api/agents/policy}"
MAX_PER_CALL_USD="${MAX_PER_CALL_USD:-0.05}"
DAILY_LIMIT_USD="${DAILY_LIMIT_USD:-0.25}"

python3 - "$PAYER_ADDRESS" "$MAX_PER_CALL_USD" "$DAILY_LIMIT_USD" <<'PY' |
import json
import sys

print(json.dumps({
    "wallet": sys.argv[1],
    "enabled": True,
    "max_per_call_usd": float(sys.argv[2]),
    "daily_limit_usd": float(sys.argv[3]),
    "allowed_resources": [
        "/api/check-prompt",
        "/api/services/rapid-policy/check-prompt",
        "/api/services/deep-shield/check-prompt"
    ]
}))
PY
curl -fsS -X POST "$POLICY_URL" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" \
  --data-binary @- | python3 -m json.tool

echo
echo "Official payer policy configured. No private key was read and no payment was signed."
