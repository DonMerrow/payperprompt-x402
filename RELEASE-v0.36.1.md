# PayPerPrompt v0.36.1 — Unified Official Settlement Accounting

This maintenance release fixes two judge-dashboard accounting gaps discovered
after a successful AI-routed browser-wallet settlement.

## Durable payer visibility

`GET /api/agents` now includes wallets found in:

- balances
- policies
- transactions
- active or historical reservations

A real payer represented only by verified settlement records therefore remains
visible with its default policy, daily spending, allowed count, denied count,
and pending reservation total.

## Deduplicated official analytics

`GET /api/proof/official/analytics` now merges:

- append-only official proof history
- durable Go transactions carrying `official-x402-onchain` evidence

The transaction hash is the deduplication key. A settlement present in both
sources contributes exactly once to count and USDC volume.

This release performs no payment, signature, wallet connection, or chain
mutation.
