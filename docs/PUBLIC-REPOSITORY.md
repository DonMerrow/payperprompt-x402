# Public Repository Guide

This repository is intended to be understandable, reproducible, and safe to
inspect without access to the original developer's machine.

## First run

```bash
git clone YOUR_PUBLIC_REPOSITORY_URL
cd payperprompt-x402

ollama pull llama3.1:8b
./scripts/configure-public-demo.sh
./scripts/test-submission-proof-kit.sh
./scripts/start-judge-stack.sh
```

The guided configuration accepts public testnet addresses only. Read
`START-HERE-OFFICIAL.md` for the judge flow and
`docs/WALLET-CONFIGURATION.md` before connecting a wallet.

## Repository map

| Path | Purpose |
|---|---|
| `go-core/` | durable policy, prepared-work escrow, audit, reconciliation, dashboard API |
| `real-x402-go/` | official x402 challenge, payment verification, settlement, and AI work lane |
| `receipt-verifier-rust/` | independent semantic and cross-field proof verification |
| `sdk/go/` | dependency-free integration SDK and inspector |
| `web/` | browser-wallet workspace and evidence dashboard |
| `scripts/` | setup, acceptance, launch, recovery, and packaging commands |
| `docs/` | architecture, proof, hosting, submission, and operator guidance |
| `examples/` | bounded large-contract test prompts |

## Safe contribution workflow

1. Create a branch.
2. Keep runtime configuration in `.env.local`.
3. Run `./scripts/test-submission-proof-kit.sh`.
4. Confirm no private keys, seed phrases, policy tokens, durable data, proofs,
   dependencies, or build outputs are staged.
5. Explain payment-boundary changes explicitly in the pull request.

Changes that weaken fail-closed preparation, replay protection, payer policy,
one-time release, settlement pinning, or proof verification require focused
regression tests.

## Build a public-safe archive

```bash
./scripts/build-submission-package.sh
sha256sum -c dist/payperprompt-x402-v0.41.0.1-public-repository.zip.sha256
```

The packaging script excludes local configuration, runtime state, proof files,
downloaded dependencies, compiled binaries, and build caches.

## Evidence versus simulation

The developer simulator is clearly fake and intended for repeatable local
tests. Official claims use Base Sepolia USDC transactions, JSON-RPC transfer
verification, official x402 middleware, and an independent Rust verifier.

The representative frozen proof remains documented even if no local runtime
state exists. New operators should not mistake that historical evidence for
their own configured payer or merchant.
