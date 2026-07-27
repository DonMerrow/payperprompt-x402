# PayPerPrompt v0.34 — Live Official Judge Flow

This release makes real x402 evidence the primary public demonstration.

## What judges see

- the official x402 service is started by the one-command judge stack
- a button sends a fresh unpaid request to official x402 middleware
- the resulting HTTP 402 `PAYMENT-REQUIRED` header is decoded and validated
- the recorded real Base Sepolia USDC settlement is checked live by JSON-RPC
- the independent Rust service verifies the official proof
- the fake-fund developer simulator is collapsed and clearly labeled optional

The public action never reads a wallet key, signs a payment, or moves funds.
Creating a new on-chain settlement remains an explicit local-terminal action.

## Run

```bash
cd ~/payperprompt-x402
./scripts/start-judge-stack.sh
```

The launcher now starts:

- official x402 AI service on `127.0.0.1:8082`
- durable Go control plane and dashboard on `127.0.0.1:8084`
- independent Rust verifier on `127.0.0.1:8085`
- a temporary Cloudflare tunnel to the dashboard

## No-payment proof test

With the judge stack running:

```bash
cd ~/payperprompt-x402
./scripts/test-live-official-evidence.sh
```

Expected:

```text
LIVE OFFICIAL JUDGE EVIDENCE PASSED
Fresh official HTTP 402 decoded and validated.
Recorded Base Sepolia settlement verified live.
Independent Rust proof verification passed.
No payment was signed or sent.
```
