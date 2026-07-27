#!/usr/bin/env bash
set -euo pipefail

echo "Rust verifier"
curl -fsS http://127.0.0.1:8085/api/health
echo
echo

echo "Go core"
curl -fsS http://127.0.0.1:8084/api/health
echo
echo

echo "Official x402 service"
curl -fsS http://127.0.0.1:8082/api/health
echo
echo

echo "Fresh official HTTP 402 evidence (no signature, no payment)"
curl -fsS \
  -H "content-type: application/json" \
  -d '{"prompt":"Please improve this public product description."}' \
  http://127.0.0.1:8084/api/official/challenge
echo
echo

echo "Latest official x402 proof"
curl -fsS http://127.0.0.1:8084/api/proof/official
echo
