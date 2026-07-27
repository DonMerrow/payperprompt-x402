# Devpost Fields — Final Paste Copy

## Project name

PayPerPrompt x402

## Tagline

See proof the AI work is ready. Then pay for exactly that work with x402.

## About the project

## Inspiration

AI agents can already call APIs and make decisions, but buying digital work
still creates a basic trust problem: **how does the buyer know useful work
exists before paying for it?**

A normal paywall proves that money moved. It does not prove what was purchased,
whether the result was ready, whether the price matched policy, or whether the
delivered result is the same work the buyer expected.

PayPerPrompt started as a paid prompt-safety API and grew into a broader idea:
**proof-before-payment AI commerce**.

The AI completes the requested work first. Go validates and commits the hidden
deliverable. The buyer sees a summary, coverage, deterministic checks, expiry,
route, exact USDC price, and SHA-256 commitment before a wallet opens. Official
x402 settlement then releases only that exact committed result.

The paying customers are AI builders, agent platforms, and small SaaS teams
that need useful work without subscriptions, account provisioning, or
uncontrolled agent spending.

## What it does

PayPerPrompt is a prepared-work AI workspace with three exact service tiers:

- **Local Guard:** 0.01 USDC
- **Rapid Policy:** 0.02 USDC
- **Deep Shield:** 0.04 USDC

The user submits a real request for writing, analysis, code review, bug triage,
meeting actions, document work, prompt-security review, or bounded Solidity
assistance.

The system then follows this lifecycle:

```text
Prepare → Ground → Commit → Challenge → Sign → Retry → Settle → Release → Verify
```

1. Ollama identifies the task and selects cost, latency, or enhanced quality.
2. The work is completed before payment.
3. Go applies task-specific quality and source-grounding checks.
4. The full deliverable stays hidden while Go exposes its proof and SHA-256
   commitment.
5. The protected API returns an official HTTP 402 payment requirement.
6. A trusted browser wallet displays the exact Base Sepolia USDC
   authorization.
7. The user explicitly approves the payment.
8. The request retries and official x402 middleware settles it.
9. Only the previously committed work is released.
10. The successful response includes the work, receipt, route, amount, payer,
    merchant, transaction, commitment, and BaseScan link.
11. Go checks the transfer through Base Sepolia JSON-RPC, and Rust
    independently checks proof consistency.

Preparation is free. If Ollama is unavailable, validation fails, the
commitment expires, or spending policy denies the route, MetaMask cannot open
and the user cannot be charged.

## Why it is different

PayPerPrompt does not ask the buyer to pay for an AI promise. It uses a
**commit-and-release pattern**:

- useful work is prepared first;
- the hidden result is bound to a cryptographic commitment;
- price and policy are visible before consent;
- settlement releases that exact result once;
- independent evidence connects payment to the purchased work.

It is not a custodial escrow and does not deploy contracts or hold customer
funds. The innovation is at the HTTP application layer: binding prepared work,
wallet consent, x402 settlement, and a verifiable receipt into one lifecycle.

## How we built it

The system has four independently visible proof layers:

1. **Ollama AI** plans the route and produces the work.
2. **Go control plane** grounds output, commits work, enforces durable spending
   policy, reserves budget atomically, prevents replay, and handles recovery.
3. **Official x402 middleware** performs the HTTP 402 challenge, paid retry,
   verification, and Base Sepolia USDC settlement.
4. **Rust verifier** independently checks canonical proof fields and rejects
   altered evidence.

The project also includes:

- a reusable dependency-free Go SDK;
- a Go CLI for catalog discovery, analysis, challenge debugging, preflight,
  payment, facilitator diagnostics, proof verification, and reconciliation;
- current-only audit evidence with bounded history on request;
- health-aware failover for read-only facilitator operations;
- a single-attempt settlement boundary;
- crash-safe proof reconciliation instead of blind duplicate payment;
- an optional fake-fund simulator for repeatable developer tests;
- a Smart Contract Studio that can generate, explain, test, repair, and
  defensively review Solidity without deploying contracts or handling secrets.

## Challenges we faced

### Treating AI output as untrusted

The model sometimes omitted source elements or produced confident but
incorrect contract claims. We added deterministic source findings, semantic
coverage checks, bounded correction, and fail-closed payment preparation.
Unsupported claims cannot simply pass because they sound convincing.

### Preventing payment for failed work

Earlier designs produced the result after settlement. We reversed the order:
prepare and validate first, commit the hidden result, then allow payment.

### Avoiding duplicate settlement

Failing over a read-only verification call is safe. Repeating an uncertain
settlement may pay twice. PayPerPrompt pins the selected facilitator for
settlement and enters reconciliation when the result is ambiguous.

### Keeping testnet evidence honest

The optional simulator uses fake local funds and is clearly labelled. Official
claims come from explorer-verifiable Base Sepolia USDC transfers checked
through JSON-RPC and independently verified by Rust.

### Making a complex system judgeable

The final dashboard presents the paying customer, prices, prepared-work proof,
wallet boundary, audit event, live chain evidence, four proof layers, all five
tracks, and developer tooling on one page.

## Accomplishments

- Completed real browser-wallet x402 payments on Base Sepolia.
- Returned transaction receipts in successful paid responses.
- Bound prepared and released work to matching SHA-256 commitments.
- Implemented three exact pay-per-call service tiers.
- Enforced durable per-wallet budgets and atomic reservations.
- Added idempotency, replay controls, exactly-once recovery, and safe
  facilitator behavior.
- Independently verified official proofs in Rust.
- Rejected altered receipts and contradictory Solidity analysis.
- Built a reusable Go SDK, CLI, challenge debugger, simulator, and proof tools.
- Demonstrated one coherent product across all five hackathon tracks.

## Verified proof

The representative prepared-work transaction is:

- **Network:** Base Sepolia (`eip155:84532`)
- **Asset:** official Base Sepolia USDC
- **Provider:** Rapid Policy
- **Route:** `guardrail-fast`
- **Amount:** 0.02 USDC (`20000` atomic units)
- **Transaction:**
  `0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5`
- **Prepared/released commitment:**
  `fb1667e4c2b5c0de9fc8e2ae2d6f1c33f1496c0569d2475180147c069ab98338`

[Open the transaction on BaseScan](https://sepolia.basescan.org/tx/0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5)

At the time of the final proof capture, the dashboard reported 17 recorded
settlements, 17 independently verified proofs, and 0.34 USDC of official Base
Sepolia volume.

## What we learned

x402 is more than a paywall. Its strongest use is machine-to-machine commerce
where the buyer can inspect the requirement, apply policy, consent to an exact
price, receive the result, and verify the transaction without creating an
account.

We also learned that AI payment systems need separation of duties:

- AI prepares and recommends.
- Deterministic policy controls.
- The wallet consents.
- x402 settles.
- Independent verification proves.

## MVP scope and production path

This is a working hackathon MVP, not a claim of finished production
infrastructure.

The public demonstration deliberately accepts a configured disposable Base
Sepolia payer. That allowlist is a safety interlock for repeatable judging, not
a claim of open multi-wallet customer onboarding. The architecture separates
wallet identity, policy, network, token, merchant, signing, settlement, and
verification so broader wallet registration and multi-tenant isolation can be
added without placing private keys in the application.

A production launch would still require stable hosting, authenticated wallet
registration, merchant isolation, managed secrets, database-backed
multi-tenancy, rate limiting, monitoring, incident response, an external
security review, and explicit mainnet controls.

## What is next

- authenticated multi-wallet and multi-merchant onboarding;
- independent service providers with signed service manifests;
- stable production hosting and observability;
- PostgreSQL-backed tenant and audit isolation;
- provider reputation and service discovery;
- external security review before mainnet operation;
- future research into standardized post-quantum signatures for long-lived
  service manifests when compatible wallet and x402 tooling exists.

PayPerPrompt is testnet-proven and mainnet-configurable. It is not represented
as a formally audited or production-deployed system.

## Built with

Use these tags:

1. Go
2. Rust
3. JavaScript
4. HTML5
5. CSS3
6. x402
7. Base
8. Base Sepolia
9. USDC
10. Ethereum
11. EIP-712
12. MetaMask
13. Coinbase Wallet
14. Rabby Wallet
15. Ollama
16. llama3.1
17. JSON-RPC
18. SHA-256
19. HMAC-SHA256
20. Cloudflare Tunnel
21. Solidity
22. Foundry
23. REST API
24. Web3
25. AI Agents

## Try it out links

Add these separately in Devpost:

1. **Live demo:** `[CURRENT_PUBLIC_DEMO_URL]`
2. **Source code:** `[PUBLIC_REPOSITORY_URL]`
3. **Verified BaseScan proof:**
   `https://sepolia.basescan.org/tx/0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5`

Do not submit the placeholders. Replace them after the public repository and
final demo URL are ready.

## Media plan

Recommended gallery order:

1. Hero: prepared-work proof before wallet signing.
2. Wallet: exact Base Sepolia USDC authorization.
3. Result: committed work released with transaction link.
4. Proof: BaseScan transfer plus live chain verification.
5. Architecture: Ollama, Go, Rust, and official x402 proof layers.
6. Developer tooling: SDK, CLI, audit, and submission matrix.

Use 3:2 images under 5 MB. Remove wallet pop-up account details beyond the
already-public disposable address, and never capture browser extensions,
private keys, seed phrases, control tokens, or unrelated tabs.
