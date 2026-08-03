#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_FILE="${PAYPERPROMPT_CONFIG_FILE:-$PROJECT_DIR/.env.local}"

valid_address() {
  [[ "$1" =~ ^0x[0-9a-fA-F]{40}$ ]]
}

echo "PayPerPrompt public demo configuration"
echo
echo "Enter public Base Sepolia addresses only."
echo "Never enter a private key, seed phrase, wallet password, or production secret."
echo

read -rp "Disposable payer public address: " payer_address
if ! valid_address "$payer_address"; then
  echo "The payer must be a 0x-prefixed 20-byte EVM address."
  exit 1
fi

read -rp "Distinct merchant public address: " merchant_address
if ! valid_address "$merchant_address"; then
  echo "The merchant must be a 0x-prefixed 20-byte EVM address."
  exit 1
fi
if [[ "${payer_address,,}" == "${merchant_address,,}" ]]; then
  echo "Payer and merchant must be different addresses."
  exit 1
fi

read -rp "Maximum USDC per call [0.05]: " max_per_call
max_per_call="${max_per_call:-0.05}"
read -rp "Daily USDC limit [2.00]: " daily_limit
daily_limit="${daily_limit:-2.00}"

if ! [[ "$max_per_call" =~ ^[0-9]+([.][0-9]+)?$ ]] ||
   ! [[ "$daily_limit" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  echo "Policy limits must be non-negative decimal numbers."
  exit 1
fi

umask 077
{
  echo "# Public addresses and local demo settings only. Never add private keys."
  printf 'PAYER_ADDRESS=%s\n' "$payer_address"
  printf 'MERCHANT_ADDRESS=%s\n' "$merchant_address"
  printf 'MAX_PER_CALL_USD=%s\n' "$max_per_call"
  printf 'DAILY_LIMIT_USD=%s\n' "$daily_limit"
  echo "X402_NETWORK=eip155:84532"
  echo "X402_ASSET=0x036CbD53842c5426634e7929541eC2318f3dCF7e"
  echo "X402_EXPECTED_ASSET=0x036CbD53842c5426634e7929541eC2318f3dCF7e"
  echo "X402_FACILITATOR_URL=https://x402.org/facilitator"
  echo "BASE_SEPOLIA_RPC_URL=https://sepolia.base.org"
  echo "OLLAMA_URL=http://127.0.0.1:11434"
  echo "OLLAMA_MODEL=qwen3-coder:30b"
  echo "CLOUDFLARED_BIN=/tmp/cloudflared"
} >"$CONFIG_FILE"
chmod 600 "$CONFIG_FILE"

echo
echo "Saved $CONFIG_FILE"
echo "No private key was requested or stored."
echo "Run: ./scripts/start-judge-stack.sh"
