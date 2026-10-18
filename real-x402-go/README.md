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

Requirements: Go 1.24.0 or newer (allow Go toolchain selection for dependencies),
Bash, `curl`, `unzip`, and `sha256sum`.

This project imports the official `github.com/x402-foundation/x402/go/v2`
library for its x402 payment protocol implementation. The original preparation
workaround used a source ZIP because the requested `v2.0.0` module tag was not
available. The ZIP is dependency source, not a separate server.

Preparation now downloads exactly upstream commit
`90688e52e58ae9185f2860988bd2c46d2801ceda` (July 24, 2026), from:

```text
https://github.com/x402-foundation/x402/archive/90688e52e58ae9185f2860988bd2c46d2801ceda.zip
SHA-256: 2204efbdb16f674e8189d5bbd202c4eacd76da9e78b31c637028db86cef9599e
```

This revision was selected from the original build period; it is not claimed to
be a recovered hash of the original unpinned download.

The script verifies downloaded and cached ZIP bytes **before extraction**, then
recreates `.deps/x402-<commit>/go` from those verified bytes. The project-relative
`replace` entry in `go.mod` tells Go where to resolve imports of the official
library. Nothing needs copying into `src`, and no owner-specific home path is
required. Other Go modules retain their normal `go.sum` verification. Preparation
runs `go mod tidy`, `go mod verify`, all project Go tests, and builds
`bin/payperprompt`; it does not start a server or submit a payment.

Old `.deps/x402-main` and Downloads copies are no longer selected automatically.
`X402_OFFICIAL_GO_DIR` is rejected rather than silently accepting unverified
source. A checksum mismatch stops preparation; remove only the named cached ZIP
and rerun to fetch the pinned archive. Do not replace the expected checksum merely
to silence a failure. GitHub archive packaging changes may also require review.

To update x402 intentionally, review the upstream changes, choose an exact commit,
verify the archive, update the commit/checksum in `scripts/prepare-official.sh`
and the relative path in `go.mod`, then rerun preparation and commit resulting
`go.mod`/`go.sum` changes. Update this README with the new provenance. Pinning
prevents silent source changes; it is not a claim that all dependency code has
been audited.

To check the download rejection behavior without network access or a real Go
build, run:

```bash
bash scripts/test-prepare-integrity.sh
```

This exercises corrupt downloads, corrupt cached archives, and rejected local
source overrides. Each must stop before any Go command runs.

## Exact Launch Order

Ollama must be running with `qwen3-coder:30b`.

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
