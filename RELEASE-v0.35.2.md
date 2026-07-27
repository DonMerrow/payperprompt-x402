# PayPerPrompt v0.35.2 — Self-Preparing Judge Stack

This maintenance release fixes two independent startup-test problems.

## Correct idempotency assertion

The v0.35.1 test searched raw JSON text for the transaction hash. A valid
single proof contains the hash in both `settlement.transaction` and
`explorer_url`, so the test incorrectly reported two history entries.

The corrected test:

- requires exactly one JSONL history record
- parses that record
- checks `settlement.transaction` directly

## Automatic official dependency preparation

`start-judge-stack.sh` now runs `real-x402-go/scripts/prepare-official.sh`
before launching services. This restores the official local replacement,
regenerates `go.sum`, runs official Go tests, and rebuilds the CLI.

No private key is read, no signature is requested, and no payment is sent
during preparation or reconciliation.
