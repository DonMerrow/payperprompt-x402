# Architecture

PayPerPrompt has four independently visible parts:

- Ollama worker: identifies the task, selects cost, latency, or quality, and
  prepares the requested work
- Go control plane: grounds model output, commits hidden work, protects
  resources, enforces durable spend policy, and reserves budget atomically
- official x402 lane: returns `402`, verifies payment, settles Base Sepolia
  USDC, and records public proof
- Rust verifier: independently rejects altered local and official evidence

The payment layer is isolated so the demo verifier can be replaced with the official x402 SDK or facilitator flow.

## Request Lifecycle

```text
request
prepare and validate work
commit hidden deliverable
return 402 challenge
wallet signs exact amount
retry with payment payload
verify and settle
release matching committed work
return work and receipt
independently verify proof
```

Preparation failure, policy denial, expiry, wallet mismatch, or semantic
rejection stops before signing. A successful release must match the exact
SHA-256 commitment shown before payment.

## Facilitator Failure Boundary

Support discovery and verification are safe to fail over because they do not
move funds. Settlement is different:

1. The successful verifier is pinned to the exact payload and requirements.
2. The server calls that facilitator once for settlement.
3. A successful response records the transaction.
4. A definitive rejection is returned as a rejection.
5. A transport/protocol error after the settlement call begins is marked
   `unknown_requires_reconciliation`.
6. No second facilitator receives the settlement payload.
7. The durable authorization remains held until chain/proof reconciliation
   resolves the outcome.

## Production Controls

- Require HTTPS.
- Persist receipt and request IDs.
- Reject duplicate settlement attempts.
- Rate limit unpaid and failed requests.
- Store only receipt metadata by default.
- Keep prompt retention configurable.
- Separate payment verification from prompt analysis.
- Never automatically retry an ambiguous settlement through another
  facilitator.
