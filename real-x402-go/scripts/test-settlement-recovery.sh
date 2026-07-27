#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$PROJECT_DIR/../go-core"
go test ./cmd/server -run 'TestOfficialSettlementRecoveryCommitsExpiredReservationExactlyOnce' -count=1

cd "$PROJECT_DIR"
go test ./cmd/official-client \
  -run 'Test(PaidSettlementFailureDoesNotAuthorizeReservationRelease|ProofPreservesPolicyAuthorizationForRecovery)' \
  -count=1

echo
echo "SETTLEMENT RECOVERY TEST PASSED"
echo "A verified payment can be committed after interruption, exactly once."
echo "The recovery path does not require a private key and does not send another payment."
