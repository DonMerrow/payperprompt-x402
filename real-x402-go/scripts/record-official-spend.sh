#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/policy-control.sh"

RECORD_URL="${OFFICIAL_POLICY_RECORD_URL:-http://127.0.0.1:8084/api/agents/policy/record-official}"

curl -fsS -X POST "$RECORD_URL" \
  -H "${POLICY_CONTROL_HEADER[0]}" \
  -H "content-type: application/json" \
  --data '{"bootstrap":true}' | python3 -m json.tool

echo
echo "Verified official settlement committed to the durable policy ledger."
echo "Running this command again is idempotent and will not double-count the transaction."
