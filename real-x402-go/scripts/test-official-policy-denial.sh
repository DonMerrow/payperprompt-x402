#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/policy-control.sh"

if [[ -z "${PAYER_ADDRESS:-}" ]]; then
  echo "PAYER_ADDRESS is required. Use the public testnet wallet address only."
  exit 1
fi

POLICY_CONFIG_URL="${OFFICIAL_POLICY_CONFIG_URL:-http://127.0.0.1:8084/api/agents/policy}"
POLICY_EVALUATE_URL="${OFFICIAL_POLICY_URL:-http://127.0.0.1:8084/api/agents/policy/evaluate}"

python3 - "$PAYER_ADDRESS" <<'PY' |
import json
import sys

print(json.dumps({
    "wallet": sys.argv[1],
    "enabled": True,
    "max_per_call_usd": 0.01,
    "daily_limit_usd": 0.25,
    "allowed_resources": [
        "/api/check-prompt",
        "/api/services/rapid-policy/check-prompt",
        "/api/services/deep-shield/check-prompt"
    ]
}))
PY
curl -fsS -X POST "$POLICY_CONFIG_URL" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" \
  --data-binary @- >/dev/null

response_file="$(mktemp)"
status="$(
  python3 - "$PAYER_ADDRESS" <<'PY' |
import json
import sys
print(json.dumps({
    "wallet": sys.argv[1],
    "resource": "/api/services/rapid-policy/check-prompt",
    "amount_usd": 0.02
}))
PY
  curl -sS -o "$response_file" -w '%{http_code}' -X POST "$POLICY_EVALUATE_URL" \
    -H "content-type: application/json" \
    --data-binary @-
)"

python3 -m json.tool "$response_file"
if [[ "$status" != "403" ]]; then
  echo "Expected HTTP 403 policy denial; received $status."
  exit 1
fi
if ! python3 - "$response_file" <<'PY'
import json
import sys
data = json.load(open(sys.argv[1], encoding="utf-8"))
raise SystemExit(0 if data.get("allowed") is False and data.get("signed") is False and data.get("settled") is False else 1)
PY
then
  echo "Policy response did not prove fail-closed behavior."
  exit 1
fi

echo
echo "OFFICIAL POLICY DENIAL PASSED"
echo "The selected \$0.02 route was denied before any signature or settlement."
