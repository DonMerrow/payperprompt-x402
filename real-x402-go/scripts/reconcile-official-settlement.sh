#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

echo "OFFICIAL x402 RECOVERY"
echo "This command verifies and records the latest settled proof."
echo "It does not read a private key, sign a payment, or send a payment."
echo

go run ./cmd/official-client reconcile
