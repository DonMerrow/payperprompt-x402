#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/policy-control.sh"

if [[ -z "${PAYER_ADDRESS:-}" ]]; then
  echo "PAYER_ADDRESS is required. Use the public testnet wallet address only."
  exit 1
fi

BASE_URL="${OFFICIAL_POLICY_BASE_URL:-http://127.0.0.1:8084/api/agents/policy}"
first_file="$(mktemp)"
second_file="$(mktemp)"
trap 'rm -f "$first_file" "$second_file"' EXIT

policy_payload="$(
  python3 - "$PAYER_ADDRESS" <<'PY'
import json
import sys
print(json.dumps({
    "wallet": sys.argv[1],
    "enabled": True,
    "max_per_call_usd": 0.02,
    "daily_limit_usd": 0.05,
    "allowed_resources": ["/api/services/rapid-policy/check-prompt"]
}))
PY
)"
curl -fsS -X POST "$BASE_URL" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" \
  --data "$policy_payload" >/dev/null

reservation_payload="$(
  python3 - "$PAYER_ADDRESS" <<'PY'
import json
import sys
print(json.dumps({
    "wallet": sys.argv[1],
    "resource": "/api/services/rapid-policy/check-prompt",
    "route_id": "guardrail-fast",
    "provider": "Rapid Policy",
    "amount_usd": 0.02
}))
PY
)"

first_status="$(curl -sS -o "$first_file" -w '%{http_code}' -X POST "$BASE_URL/reserve" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" --data "$reservation_payload")"
second_status="$(curl -sS -o "$second_file" -w '%{http_code}' -X POST "$BASE_URL/reserve" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" --data "$reservation_payload")"

echo "FIRST RESERVATION"
python3 -m json.tool "$first_file"
echo
echo "SECOND RESERVATION"
python3 -m json.tool "$second_file"

if [[ "$first_status" != "201" || "$second_status" != "403" ]]; then
  echo "Expected first reservation HTTP 201 and second HTTP 403."
  exit 1
fi

authorization_id="$(
  python3 - "$first_file" <<'PY'
import json
import sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["authorization_id"])
PY
)"
curl -fsS -X POST "$BASE_URL/release" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" \
  --data "$(python3 - "$authorization_id" <<'PY'
import json
import sys
print(json.dumps({"authorization_id": sys.argv[1], "reason": "atomic reservation test completed"}))
PY
)" >/dev/null

MAX_PER_CALL_USD=0.05 DAILY_LIMIT_USD=0.25 \
  "$(dirname "$0")/configure-official-policy.sh" >/dev/null

echo
echo "ATOMIC RESERVATION TEST PASSED"
echo "The first \$0.02 authorization reserved budget; the competing authorization was denied."
echo "The test reservation was released and the normal policy was restored."
