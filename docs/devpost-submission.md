# PayPerPrompt x402 — Devpost Submission

The exact field-by-field paste copy, Built With tags, links, and media order
are maintained in `submission/DEVPOST-PASTE.md`.

## Project name

PayPerPrompt x402

## Elevator pitch

PayPerPrompt lets a user inspect completed, validated AI work before deciding
whether to pay. Ollama prepares the work, Go commits and controls it, MetaMask
authorizes the exact x402 charge, official middleware settles USDC on Base
Sepolia, and Rust independently verifies the proof.

## Paying customer and business model

The paying customers are AI builders, agent platforms, and small SaaS teams
that need useful AI work without subscriptions, account provisioning, or
unsafe wallet practices.

Each completed request costs an exact amount:

- Local Guard: 0.01 USDC
- Rapid Policy: 0.02 USDC
- Deep Shield: 0.04 USDC

Preparation is free. A failed quality gate, expired commitment, policy denial,
or unavailable AI cannot open the wallet or charge the user.

## Inspiration

Machine-payable APIs solve checkout friction, but payment alone does not prove
that useful work exists. An autonomous buyer also needs to know what is being
purchased, what it costs, whether policy permits it, and whether the delivered
result is the work that was inspected before payment.

PayPerPrompt turns those requirements into one auditable lifecycle.

## What it does

1. The user submits a real work request.
2. Ollama identifies the task and selects cost, speed, or enhanced quality.
3. Go completes deterministic grounding and semantic validation.
4. Go keeps the full deliverable hidden and shows its summary, coverage,
   expiry, exact route, exact price, and SHA-256 commitment.
5. Official x402 middleware returns HTTP 402 with the payment requirement.
6. The trusted browser wallet signs the exact Base Sepolia USDC authorization.
7. The paid request retries and the facilitator settles it.
8. Only the committed work is released.
9. The response returns the receipt, transaction, route, amount, payer,
   merchant, and commitment.
10. Go checks the chain through JSON-RPC and Rust independently checks proof
    consistency.

The workspace supports writing, analysis, code review, bug triage, meeting
actions, document work, prompt-security review, and a bounded Smart Contract
Studio for Solidity generation, explanation, tests, repair, and defensive
review.

## How we built it

- Go control plane with durable policy, atomic reservations, idempotency,
  prepared-work commitments, release control, audit metadata, and recovery
- official `github.com/x402-foundation/x402/go/v2` HTTP middleware
- Base Sepolia and official Base Sepolia USDC
- trusted injected-wallet support for MetaMask, Coinbase Wallet, and Rabby
- local Ollama `llama3.1:8b`
- deterministic source-grounding and semantic gates
- independent Rust proof verifier
- reusable dependency-free Go SDK
- Go CLI, challenge debugger, simulator, preflight, facilitator diagnostics,
  proof verification, and reconciliation tools
- dependency-light HTML, CSS, and browser JavaScript dashboard

## Complete payment flow

```text
Prepare → Commit → Challenge → Sign → Retry → Settle → Release → Verify
```

The required x402 sequence is directly visible inside that larger workflow:

```text
Challenge → Sign → Retry → Settle
```

## Challenges

### Paying only for credible work

Earlier versions generated work after payment. The final design completes and
validates it first, commits the hidden result, and releases only that exact
result after settlement.

### Treating AI output as untrusted

Model output can omit facts or invent claims. Go derives source facts,
rejects known contradictions, requires task-specific coverage, and blocks
wallet signing when validation fails.

### Avoiding duplicate payments

Read-only facilitator operations may fail over. Settlement is attempted once.
An ambiguous result enters reconciliation instead of being sent blindly to a
second endpoint.

### Making testnet evidence honest

Local fake funds remain an optional developer simulator. Official claims come
from real Base Sepolia USDC transfers verified through JSON-RPC and BaseScan.

## Accomplishments

- Real browser-wallet x402 payments on Base Sepolia
- Three exact pay-per-call service tiers
- Prepared-work SHA-256 commitments and one-time release
- Durable per-wallet spending policy and atomic budget reservations
- Crash-safe settlement recovery and exactly-once reconciliation
- Health-aware facilitator selection without duplicate-settlement failover
- Independent Rust proof verification and tamper rejection
- Current-only audit by default with bounded history on request
- Go SDK and CLI covering discovery, planning, debugging, proof, and audit
- Smart Contract Studio with deterministic Solidity grounding
- One coherent dashboard mapping all five hackathon tracks to evidence

## Representative official proof

| Field | Value |
|---|---|
| Network | Base Sepolia (`eip155:84532`) |
| Asset | Official Base Sepolia USDC |
| Provider | Rapid Policy |
| Route | `guardrail-fast` |
| Amount | 0.02 USDC (`20000` atomic units) |
| Payer | `0x826154a3d58aea3fbd2aa64aad424594ade927ef` |
| Merchant | `0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630` |
| Transaction | `0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5` |
| Prepared-work commitment | `fb1667e4c2b5c0de9fc8e2ae2d6f1c33f1496c0569d2475180147c069ab98338` |

Explorer:
<https://sepolia.basescan.org/tx/0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5>

At the time this proof was captured, the dashboard reported 17 recorded
settlements, 17 independently verified proofs, and 0.34 USDC of official Base
Sepolia volume.

## Five-track coverage

### Track 1 — x402-powered AI applications

Useful prepared AI work is released only after an exact x402 settlement.

### Track 2 — agentic commerce and payment infrastructure

AI routing, service discovery, spend policy, atomic reservations, durable
audit, receipt verification, facilitator safety, and reconciliation.

### Track 3 — developer tools and SDKs

Reusable Go SDK, CLI, challenge debugger, simulator, preflight, payment
commands, proof verifier, and inspection example.

### Track 4 — DeFi, Web3, and tokenized finance

Trusted browser-wallet integration and official USDC settlement on Base
Sepolia, independently checked against the ERC-20 Transfer event.

### Track 5 — open innovation

A practical trust layer for AI agents purchasing cybersecurity, code,
document, and business services under explicit authorization limits.

## What we learned

x402 is most powerful when the payment proof is connected to the purchased
result. The transaction should prove not only that money moved, but which
service, route, price, policy decision, and committed work the buyer received.

AI should recommend and prepare. Deterministic policy should authorize. The
wallet should consent. x402 should settle. Independent verification should
prove.

## What is next

- stable production hosting and monitoring
- multi-tenant merchant and policy isolation
- managed production secret storage and incident response
- independent providers with signed service manifests
- external security review before mainnet operation
- research into standardized post-quantum signatures for long-lived service
  manifests when compatible wallet and x402 tooling exists

The current implementation is testnet-proven and mainnet-configurable. It is
not represented as a production or formally audited system.

## Technologies

Go, Rust, x402, Base Sepolia, USDC, Ollama, llama3.1:8b, EIP-712, JSON-RPC,
SHA-256, HMAC-SHA256, HTML, CSS, JavaScript, and Cloudflare Tunnel.
