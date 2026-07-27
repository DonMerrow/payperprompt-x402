# PayPerPrompt Independent Rust Receipt Verifier

This service independently verifies receipt envelopes issued by the durable Go
core.

It checks:

- required settlement semantics
- replay-protection declaration
- positive amount
- canonical receipt fields
- constant-time HMAC-SHA256 verification
- tampering with amount, payer, route, transaction, or timestamp

## Run

Use the same local development secret as the Go core:

```bash
cd ~/payperprompt-x402/receipt-verifier-rust
export RECEIPT_HMAC_SECRET="local-test-secret-shared-with-rust"
cargo test
cargo run
```

The verifier listens on:

```text
http://127.0.0.1:8085
```

Start the Go core in another terminal with the same secret.

The secret is local development configuration. Never commit a production
receipt-signing secret.
