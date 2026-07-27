# PayPerPrompt Durable Go Core

This is the production-oriented local gateway for PayPerPrompt x402 Lab.

It uses only the Go standard library and provides:

- durable atomic JSON state
- per-agent balances and spend policies
- daily spend enforcement across restarts
- policy-aware service routing
- local Ollama analysis
- HMAC-SHA256 signed receipt envelopes
- optional independent Rust receipt verification
- audit export
- the existing browser UI

## Run

Start the Rust verifier first if available, then:

```bash
cd ~/payperprompt-x402/go-core
export RECEIPT_HMAC_SECRET="local-test-secret-shared-with-rust"
go run ./cmd/server
```

Open:

```text
http://127.0.0.1:8084
```

Durable state is written atomically to:

```text
go-core/data/runtime-state.json
```

Do not commit runtime state or production secrets.

## Restart Test

1. Fund the sandbox agent.
2. Run an AI mission.
3. Stop the Go server.
4. Start it again with the same `RECEIPT_HMAC_SECRET`.
5. Refresh the page.

The balance, policy, receipts, transaction history, and daily spending total
must remain unchanged.
