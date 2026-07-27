# PayPerPrompt Go SDK

This dependency-free client integrates applications with the durable
PayPerPrompt control plane.

It can:

- discover priced services;
- prepare and validate committed AI work without signing;
- read the current audit event or request bounded history;
- run the read-only official x402 evidence check.

It deliberately does not accept private keys or implement payment signatures.
Use the official x402 Go client or the trusted browser-wallet flow for signing
and settlement.

```go
client := payperprompt.New("http://127.0.0.1:8084")
services, err := client.Services(context.Background())
```

Inspect a running gateway:

```bash
cd sdk/go
go test ./...
go run ./cmd/inspect -url http://127.0.0.1:8084
```

Add `-history` only when historical audit metadata is needed.
