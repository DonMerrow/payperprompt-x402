#!/usr/bin/env bash
set -euo pipefail

response_file="$(mktemp)"
trap 'rm -f "$response_file"' EXIT

curl -fsS \
  -H "content-type: application/json" \
  -d '{"prompt":"Please improve the clarity of this public product description."}' \
  http://127.0.0.1:8084/api/official/evidence \
  >"$response_file"

python3 - "$response_file" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as handle:
    evidence = json.load(handle)

challenge = evidence["fresh_challenge"]
checks = challenge["checks"]
live = evidence["recorded_real_settlement"]["live_chain_verification"]
rust = evidence["rust_verification"]

assert challenge["observed_http_status"] == 402
assert challenge["payment_signed"] is False
assert challenge["payment_sent"] is False
assert all(checks.values()), checks
assert live["valid"] is True, live
assert rust["valid"] is True, rust
assert evidence["passed"] is True, evidence

print("LIVE OFFICIAL JUDGE EVIDENCE PASSED")
print("Fresh official HTTP 402 decoded and validated.")
print("Recorded Base Sepolia settlement verified live.")
print("Independent Rust proof verification passed.")
print("No payment was signed or sent.")
PY
