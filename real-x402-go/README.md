# Official x402 Go + Ollama Integration

This folder is PayPerPrompt's real Base Sepolia lane. It uses the official x402 Go SDK, the x402.org test facilitator, and a local Ollama model.

It is intentionally separate from the durable Go/Rust sandbox. A chain or facilitator problem cannot break the offline developer demo.

## Why Go

The official x402 repository documents a Go module:

```bash
go get github.com/x402-foundation/x402/go/v2
```

Go gives this project a serious infrastructure path without tying the main demo to npm packages.

## Verified Target Flow

```text
POST /api/check-prompt
  no payment payload
  -> 402 Payment Required + PAYMENT-REQUIRED header

POST /api/check-prompt
  PAYMENT-SIGNATURE: <base64 payment payload>
  -> facilitator /verify
  -> run grounded Ollama prompt-safety analysis
  -> facilitator /settle
  -> 200 OK + PAYMENT-RESPONSE header + receipt JSON
```

The acceptance test is strict:

1. Payer and merchant are two different test-only wallets.
2. The unpaid request returns official HTTP 402 plus `PAYMENT-REQUIRED`.
3. The client validates Base Sepolia, the USDC contract, recipient, and amount.
4. Only the explicit `pay` command signs and retries.
5. The facilitator returns a successful transaction hash.
6. BaseScan shows the transfer from payer to merchant.
7. The client independently queries Base Sepolia and matches the successful USDC
   `Transfer` event by token, payer, merchant, and exact atomic amount.
8. `proof/official-settlement.json` records the latest public proof and paid AI response.
9. `proof/official-settlements.jsonl` preserves append-only settlement history.
10. The proof preserves the durable Go authorization ID before accounting commit.
11. If accounting is interrupted after settlement, reconciliation verifies the
    proof and commits the reservation exactly once without signing or paying again.
12. Ordered facilitator failover is allowed for support discovery and
    verification, but settlement is never automatically retried after an
    ambiguous response.

## Official Agentic Service Tiers

The official server exposes three separately priced x402 resources:

| AI strategy | Service | Resource | Price |
|---|---|---|---:|
| `lowest-cost` | Local Guard | `/api/check-prompt` | $0.01 |
| `lowest-latency` | Rapid Policy | `/api/services/rapid-policy/check-prompt` | $0.02 |
| `highest-quality` | Deep Shield | `/api/services/deep-shield/check-prompt` | $0.04 |

Ollama selects a strategy from risk and urgency. `agent-preflight` proves the
selected endpoint's HTTP 402 challenge matches the selected price without
signing. Only `agent-pay` authorizes a real testnet settlement.

## Security Controls To Show Judges

- idempotent `X-Request-Id`
- replay protection for payment signatures
- max prompt size
- clear retention policy
- no account required
- no API key required for buyers
- merchant wallet controlled by project owner

## First-Time Preparation

Run once:

```bash
cd ~/payperprompt-x402/real-x402-go
./scripts/prepare-official.sh
```

The script uses `~/Downloads/x402-official-readonly/go` when present. Otherwise it downloads the version pinned in `go.mod`. It runs all Go tests before reporting success.

## Exact Launch Order

Ollama must be running with `llama3.1:8b`.

Terminal 1 — receiving server. This wallet is public and must differ from the payer:

```bash
cd ~/payperprompt-x402/real-x402-go
export MERCHANT_ADDRESS="0xPUBLIC_ADDRESS_OF_SECOND_TEST_WALLET"
./scripts/run-official-server.sh
```

Terminal 2 — payer. Enter a test-only private key locally; do not paste it into chat, files, screenshots, or source:

```bash
cd ~/payperprompt-x402/real-x402-go
read -rsp "Test-only payer private key: " EVM_PRIVATE_KEY
echo
export EVM_PRIVATE_KEY
./scripts/official-preflight.sh
```

Preflight makes no payment. It must print:

```text
Checks:   distinct wallets ✓  network ✓  token ✓  nonzero price ✓
PREFLIGHT PASSED — no payment was signed or sent.
```

Only after checking both addresses and balances, make one $0.01 test-USDC payment:

```bash
./scripts/official-pay.sh
```

Success prints the transaction, BaseScan URL, proof path, and paid Ollama response.

## CLI and Transaction Debugger

Build the standalone Go CLI:

```bash
./scripts/build-cli.sh
```

Read-only commands do not require a private key:

```bash
./bin/payperprompt catalog
./bin/payperprompt analyze
./bin/payperprompt debug-challenge
./bin/payperprompt verify-proof
./bin/payperprompt facilitators
```

Preview the AI-selected official route without paying:

```bash
./scripts/official-agent-preflight.sh
```

After reviewing its provider, endpoint, recipient, token, network, and exact
price, explicitly authorize that route:

```bash
./scripts/official-agent-pay.sh
```

If a confirmed settlement reports that durable policy accounting failed, do
not run the payment again. Reconcile the preserved proof:

```bash
./scripts/reconcile-official-settlement.sh
```

This recovery command does not read a private key, sign a payment, or send a
payment. Its protected Go endpoint verifies the latest proof live on Base
Sepolia and commits the matching authorization idempotently.

Test the recovery state machine without real funds:

```bash
./scripts/test-settlement-recovery.sh
```

Test facilitator failover and duplicate-payment prevention without real funds:

```bash
./scripts/test-facilitator-failover.sh
```

If multiple compatible facilitators are available, configure them in priority
order:

```bash
export X402_FACILITATOR_URLS="https://primary.example,https://backup.example"
```

`GetSupported` and `Verify` can use the next endpoint. `Settle` never does a
blind second attempt. Inspect the live read-only status with:

```bash
./scripts/facilitator-status.sh
```

## Files

- `cmd/official-server`: official protected Ollama API
- `cmd/official-client`: safe preflight and explicit paying client
- `internal/payperprompt/analysis.go`: grounded local-AI contract and fallback
- `scripts/prepare-official.sh`: dependency preparation and tests
- `scripts/run-official-server.sh`: receiving server
- `scripts/official-preflight.sh`: no-spend validation
- `scripts/official-pay.sh`: explicit one-call payment
- `scripts/official-agent-preflight.sh`: AI route selection plus no-spend validation
- `scripts/official-agent-pay.sh`: explicit AI-selected official settlement
- `scripts/reconcile-official-settlement.sh`: no-payment recovery after settlement/accounting interruption
- `scripts/test-settlement-recovery.sh`: deterministic recovery and idempotency tests
- `scripts/test-facilitator-failover.sh`: no-payment failover and ambiguous-settlement tests
- `scripts/facilitator-status.sh`: read-only endpoint health and settlement-safety diagnostics
- `scripts/debug-official-challenge.sh`: decode and validate HTTP 402 without paying
- `scripts/build-cli.sh`: build the standalone Go developer CLI
- `proof/official-settlement.json`: generated public proof after success
- `proof/official-settlements.jsonl`: append-only official settlement evidence

## Security Boundary

- The server receives only the merchant's public address.
- The payer private key stays in one terminal environment.
- Preflight refuses payer=merchant.
- Preflight refuses the wrong network, token contract, or a zero amount.
- The proof contains public chain facts only.
- An ambiguous settlement is never sent to a second facilitator; reconcile the
  chain/proof state before any retry.
- `https://x402.org/facilitator` is testnet-only. Do not present this configuration as mainnet.
