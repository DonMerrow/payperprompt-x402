# PayPerPrompt x402 Lab v0.25.2

## Submission Candidate

v0.25 turns the working system into a clear judge-facing product demonstration.

## Evidence Presented Together

- Ollama `llama3.1:8b` performs grounded prompt-risk analysis
- Go enforces per-call and daily budgets, routing, replay protection, durable
  state, settlement, and audit logging
- Rust independently verifies canonical HMAC-SHA256 receipts and rejects
  tampered copies
- official x402 Go middleware settled 0.01 Base Sepolia USDC
- BaseScan transaction:
  `0x03c3b1be51cedd392add099d95571e6ac4ec220e012b9670ee8dbd8b496387cb`

## Judge Demo

`Run Judge Demo · No Real Funds` performs:

1. restore a safe local agent policy
2. fund the fake local wallet
3. run an Ollama-planned paid mission
4. enforce and settle through the durable Go core
5. verify the receipt with the independent Rust service
6. alter the amount and prove both verifiers reject the tampered copy

The official on-chain proof is displayed separately and is never rerun by this
button.

## Evidence Clarity Correction

The local decision table is labeled `Local Policy Simulation Log`; it proves
routing, budgets, denials, and receipt-integrity behavior but is not described
as blockchain evidence.

The final table on the page is `Official x402 On-Chain Settlement Record`. It
shows the verified Base Sepolia network, 0.01 USDC amount, distinct payer and
merchant, and direct BaseScan transaction. The audit export preserves the same
separation under `official_x402_settlement` and `local_simulation`.

The Go core also performs a live Base Sepolia JSON-RPC check. It retrieves the
transaction receipt and accepts the proof only when all of these match:

- successful transaction status
- Base Sepolia USDC contract
- ERC-20 `Transfer` event
- expected payer
- expected merchant
- exactly `10000` atomic units (`0.01 USDC`)

If the RPC is unreachable, the interface says so and falls back to the direct
BaseScan evidence link instead of claiming a live verification.

## Honest Track Coverage

| Track | Evidence |
| --- | --- |
| 1 — x402 AI Applications | Paid AI prompt-guardrail API with real Ollama output |
| 2 — Agentic Payment Infrastructure | Policy engine, router, budgets, receipts, replay controls, and audit |
| 3 — Developer Tools & SDKs | Local x402 lab, calibration, one-click demo, tamper testing, and proof export |
| 4 — Web3 & Tokenized Finance | Official 0.01 USDC Base Sepolia x402 settlement; not presented as a full DeFi protocol |
| 5 — Open Innovation | Safe development and autonomous purchasing for paid AI services |
