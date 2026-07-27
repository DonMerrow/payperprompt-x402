#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_FILE="${PAYPERPROMPT_CONFIG_FILE:-$PROJECT_DIR/.env.local}"
if [[ -f "$CONFIG_FILE" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$CONFIG_FILE"
  set +a
fi
CLOUDFLARED_BIN="${CLOUDFLARED_BIN:-/tmp/cloudflared}"
RECEIPT_SECRET="${RECEIPT_HMAC_SECRET:-local-test-secret-shared-with-rust}"
RUST_LOG="${TMPDIR:-/tmp}/payperprompt-rust-verifier.log"
GO_LOG="${TMPDIR:-/tmp}/payperprompt-go-core.log"
OFFICIAL_LOG="${TMPDIR:-/tmp}/payperprompt-official-x402.log"
PREPARE_LOG="${TMPDIR:-/tmp}/payperprompt-official-prepare.log"
MERCHANT_ADDRESS="${MERCHANT_ADDRESS:-}"
PAYER_ADDRESS="${PAYER_ADDRESS:-}"
MAX_PER_CALL_USD="${MAX_PER_CALL_USD:-0.05}"
DAILY_LIMIT_USD="${DAILY_LIMIT_USD:-2.00}"

if [[ ! "$PAYER_ADDRESS" =~ ^0x[0-9a-fA-F]{40}$ ]] ||
   [[ ! "$MERCHANT_ADDRESS" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
  echo "Public payer and merchant addresses are not configured."
  echo "Run: ./scripts/configure-public-demo.sh"
  echo "Never enter a private key or seed phrase."
  exit 1
fi
if [[ "${PAYER_ADDRESS,,}" == "${MERCHANT_ADDRESS,,}" ]]; then
  echo "PAYER_ADDRESS and MERCHANT_ADDRESS must be different."
  exit 1
fi

if [[ ! -x "$CLOUDFLARED_BIN" ]]; then
  echo "cloudflared was not found at $CLOUDFLARED_BIN"
  exit 1
fi

if ! curl -fsS http://127.0.0.1:11434/api/tags >/dev/null; then
  echo "Ollama is not available at http://127.0.0.1:11434"
  exit 1
fi

for url in \
  http://127.0.0.1:8082/api/health \
  http://127.0.0.1:8084/api/health \
  http://127.0.0.1:8085/api/health; do
  if curl -fsS "$url" >/dev/null 2>&1; then
    echo "A PayPerPrompt service is already running at $url"
    echo "Stop the old process before using the one-command launcher."
    exit 1
  fi
done

if ! (
  cd "$PROJECT_DIR/real-x402-go"
  ./scripts/prepare-official.sh
) >"$PREPARE_LOG" 2>&1; then
  echo "Official x402 preparation failed."
  tail -80 "$PREPARE_LOG"
  exit 1
fi
echo "Official x402 dependencies ready."

cleanup() {
  trap - EXIT INT TERM
  [[ -n "${go_pid:-}" ]] && kill "$go_pid" 2>/dev/null || true
  [[ -n "${rust_pid:-}" ]] && kill "$rust_pid" 2>/dev/null || true
  [[ -n "${official_pid:-}" ]] && kill "$official_pid" 2>/dev/null || true
  wait "${go_pid:-}" "${rust_pid:-}" "${official_pid:-}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

(
  cd "$PROJECT_DIR/receipt-verifier-rust"
  export RECEIPT_HMAC_SECRET="$RECEIPT_SECRET"
  cargo run
) >"$RUST_LOG" 2>&1 &
rust_pid=$!

(
  cd "$PROJECT_DIR/real-x402-go"
  unset EVM_PRIVATE_KEY
  export MERCHANT_ADDRESS
  export X402_NETWORK="${X402_NETWORK:-eip155:84532}"
  export X402_FACILITATOR_URL="${X402_FACILITATOR_URL:-https://x402.org/facilitator}"
  go run ./cmd/official-server
) >"$OFFICIAL_LOG" 2>&1 &
official_pid=$!

(
  cd "$PROJECT_DIR/go-core"
  export RECEIPT_HMAC_SECRET="$RECEIPT_SECRET"
  export MERCHANT_ADDRESS
  export PAYER_ADDRESS
  export X402_NETWORK="${X402_NETWORK:-eip155:84532}"
  export X402_ASSET="${X402_ASSET:-0x036CbD53842c5426634e7929541eC2318f3dCF7e}"
  go run ./cmd/server
) >"$GO_LOG" 2>&1 &
go_pid=$!

wait_for_url() {
  local name="$1"
  local url="$2"
  local pid="$3"
  for _ in $(seq 1 60); do
    if ! kill -0 "$pid" 2>/dev/null; then
      echo "$name stopped during startup."
      return 1
    fi
    if curl -fsS "$url" >/dev/null 2>&1; then
      echo "$name ready: $url"
      return 0
    fi
    sleep 1
  done
  echo "$name did not become ready."
  return 1
}

if ! wait_for_url "Rust verifier" "http://127.0.0.1:8085/api/health" "$rust_pid"; then
  tail -40 "$RUST_LOG"
  exit 1
fi

if ! wait_for_url "Official x402 service" "http://127.0.0.1:8082/api/health" "$official_pid"; then
  tail -60 "$OFFICIAL_LOG"
  echo "See $PREPARE_LOG for the dependency preparation result."
  exit 1
fi

if ! wait_for_url "Go core" "http://127.0.0.1:8084/api/health" "$go_pid"; then
  tail -40 "$GO_LOG"
  exit 1
fi

if ! (
  cd "$PROJECT_DIR/real-x402-go"
  export PAYER_ADDRESS MAX_PER_CALL_USD DAILY_LIMIT_USD
  ./scripts/configure-official-policy.sh
) >/dev/null; then
  echo "The configured payer policy could not be installed."
  exit 1
fi
echo "Configured payer policy ready."

echo
echo "Judge stack is ready."
echo "Local dashboard: http://127.0.0.1:8084"
echo "Official x402 service: http://127.0.0.1:8082"
echo "Press Ctrl+C once to stop the tunnel, Go core, official x402 service, and Rust verifier."
echo

"$CLOUDFLARED_BIN" tunnel --no-autoupdate --url http://127.0.0.1:8084
