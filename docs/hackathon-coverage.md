# Hackathon Coverage and Inspectable Evidence

PayPerPrompt is one product across all five tracks: an AI prepares useful work,
Go grounds and commits it, policy authorizes the exact route, a trusted wallet
signs, official x402 settles, the committed work is released, and independent
verification proves the result.

## Track 1 — x402-powered AI applications

Implemented:

- paid AI workspace for writing, analysis, code, documents, meetings, prompt
  security, and Solidity work;
- three exact per-call prices;
- work preparation and validation before wallet signing;
- HTTP 402 challenge, paid retry, settlement, receipt, and committed release;
- no charge when preparation or policy fails.

Evidence:

- successful response contains `ai_used: true`, `work_completed: true`,
  `prepared_work_released: true`, amount, route, and transaction;
- representative 0.02 USDC Rapid Policy payment is permanent on BaseScan.

## Track 2 — agentic commerce and payment infrastructure

Implemented:

- AI choice of cost, latency, or quality;
- service catalog and exact route pricing;
- per-wallet per-call and daily policy;
- atomic pending-budget reservations;
- durable state and current/bounded audit views;
- idempotency and replay controls;
- health-aware read-only facilitator failover;
- one settlement attempt and proof reconciliation after uncertainty;
- independent receipt verification.

Evidence:

- denied policy cannot open the wallet;
- competing reservations cannot overspend;
- crash recovery commits an existing verified payment exactly once;
- ambiguous settlement cannot trigger blind duplicate payment;
- altered proof is rejected.

## Track 3 — developer tools and SDKs

Implemented:

- reusable dependency-free Go SDK;
- standalone `payperprompt` Go CLI;
- catalog and route analysis;
- HTTP 402 challenge debugger;
- no-payment preflight;
- explicit payment commands;
- facilitator diagnostics;
- proof verification and reconciliation;
- local simulator and tamper tests;
- large-contract fixtures.

Evidence:

```text
payperprompt catalog
payperprompt analyze
payperprompt debug-challenge
payperprompt agent-preflight
payperprompt agent-pay
payperprompt verify-proof
payperprompt facilitators
payperprompt reconcile
cd sdk/go && go test ./...
go run ./cmd/inspect
```

The SDK never accepts private keys. Signing remains in a trusted wallet or the
official explicit-payment client.

## Track 4 — DeFi, Web3, and tokenized finance

Implemented:

- official x402 Go middleware;
- Base Sepolia;
- official Base Sepolia USDC;
- trusted browser-wallet EIP-712 authorization;
- distinct payer and merchant;
- facilitator verification and settlement;
- JSON-RPC verification of transaction status and exact ERC-20 Transfer event.

Representative evidence:

- transaction:
  `0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5`;
- provider: Rapid Policy;
- amount: 0.02 USDC (`20000` atomic units);
- prepared work and released work share commitment
  `fb1667e4c2b5c0de9fc8e2ae2d6f1c33f1496c0569d2475180147c069ab98338`.

## Track 5 — open innovation

PayPerPrompt addresses a real AI-agent trust problem: machine buyers can spend
faster than humans can inspect each purchase. The system binds useful work,
price, policy, wallet consent, settlement, and independent evidence into one
auditable transaction.

Its defensive Smart Contract Studio also demonstrates how paid AI work can be
constrained: no deployment, no secret handling, no asset movement, and no
claim of formal audit.

## Submission expectations

| Expectation | PayPerPrompt evidence | Status |
|---|---|---|
| Identified paying customer | AI builders, agent platforms, and small SaaS teams | Complete |
| Pay-per-call model | 0.01, 0.02, or 0.04 USDC per completed request | Complete |
| Challenge → Sign → Retry → Settle | Official HTTP 402, wallet approval, paid retry, facilitator settlement | Complete |
| Receipt in successful response | Transaction, payer, merchant, amount, route, commitment, and explorer URL | Complete |
| Clean setup and docs | README, start guide, runbook, SDK, proof guide, and licence | Complete |
| Working MVP | Repeated browser-wallet and terminal evidence | Complete |
| Mainnet-ready architecture | Configurable and testnet-proven; production controls are documented | Honest qualification |

## Final no-spin acceptance gate

1. `./scripts/test-submission-proof-kit.sh` passes.
2. The stack starts cleanly after a reboot.
3. Preparing work cannot open the wallet.
4. Policy rejection and semantic rejection cannot charge.
5. One explicit prepared-work payment settles.
6. Released commitment equals prepared commitment.
7. The receipt is included in the successful response.
8. Live JSON-RPC matches the exact Base Sepolia USDC transfer.
9. Rust accepts the original proof and rejects altered evidence.
10. SDK tests and inspection example pass.
11. Public dashboard shows current audit first and bounded history only on
    request.
12. Public repository contains no secrets and includes the MIT licence.
