#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/policy-control.sh"

if [[ -z "${PAYER_ADDRESS:-}" ]]; then
  echo "PAYER_ADDRESS is required. Use the public testnet wallet address only."
  exit 1
fi

RESERVE_URL="${OFFICIAL_POLICY_RESERVE_URL:-http://127.0.0.1:8084/api/agents/policy/reserve}"
RELEASE_URL="${OFFICIAL_POLICY_RELEASE_URL:-http://127.0.0.1:8084/api/agents/policy/release}"
unauthorized_file="$(mktemp)"
authorized_file="$(mktemp)"
trap 'rm -f "$unauthorized_file" "$authorized_file"' EXIT

payload="$(
  python3 - "$PAYER_ADDRESS" <<'PY'
import json
import sys
print(json.dumps({
    "wallet": sys.argv[1],
    "resource": "/api/check-prompt",
    "route_id": "guardrail-economy",
    "provider": "Local Guard",
    "amount_usd": 0.01
}))
PY
)"

unauthorized_status="$(curl -sS -o "$unauthorized_file" -w '%{http_code}' \
  -X POST "$RESERVE_URL" -H "content-type: application/json" --data "$payload")"
authorized_status="$(curl -sS -o "$authorized_file" -w '%{http_code}' \
  -X POST "$RESERVE_URL" -H "content-type: application/json" \
  -H "${POLICY_CONTROL_HEADER[0]}" --data "$payload")"

echo "WITHOUT CONTROL TOKEN"
python3 -m json.tool "$unauthorized_file"
echo
echo "WITH LOCAL CONTROL TOKEN"
python3 -m json.tool "$authorized_file"

if [[ "$unauthorized_status" != "401" || "$authorized_status" != "201" ]]; then
  echo "Expected unauthorized HTTP 401 and authorized HTTP 201."
  exit 1
fi

authorization_id="$(
  python3 - "$authorized_file" <<'PY'
import json
import sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["authorization_id"])
PY
)"
release_payload="$(python3 - "$authorization_id" <<'PY'
import json
import sys
print(json.dumps({"authorization_id": sys.argv[1], "reason": "policy control test completed"}))
PY
)"
curl -fsS -X POST "$RELEASE_URL" \
  -H "content-type: application/json" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  --data "$release_payload" >/dev/null

echo
echo "POLICY CONTROL TEST PASSED"
echo "Public mutation was rejected; the same request with the local token was authorized."
echo "The test reservation was released."
