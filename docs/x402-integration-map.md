# x402 Integration Map

PayPerPrompt supports an optional local simulator and a primary official x402
browser-wallet flow. Developers can debug policy without funds, then buy the
same prepared-work service on Base Sepolia through a recognized injected
wallet.

## Current Demo Flow

| Step | Optional simulator | Official browser-wallet lane |
|---|---|---|
| Prepare | Local demo work | Ollama work plus deterministic Go validation |
| Commit | Local receipt metadata | Hidden deliverable bound to SHA-256 commitment |
| Challenge | `paymentRequired` returns HTTP 402 | Official middleware returns `PAYMENT-REQUIRED` |
| Sign | UI creates a deterministic local signature | MetaMask, Coinbase Wallet, or Rabby signs exact EIP-712 authorization |
| Retry | UI sends `PAYMENT-SIGNATURE` | Browser sends the wallet payment payload |
| Verify | local ledger validates signature and request ID | Official middleware and facilitator verify the payment |
| Work | local demonstration result | Previously prepared committed work is selected |
| Settle | local fake balances are debited | Official middleware and facilitator settle on Base Sepolia |
| Release | local result returned | Exact committed deliverable is released once |
| Receipt | JSON simulator receipt | Work, commitment, payer, route, amount, transaction, and `PAYMENT-RESPONSE` |

## Ordered Facilitator Pool

The official lane wraps the SDK's `FacilitatorClient` interface:

| Operation | Failover | Reason |
|---|---:|---|
| `GetSupported` | Yes | Capability discovery cannot transfer funds |
| `Verify` | Yes | Verification cannot transfer funds |
| `Settle` | No | A timeout may hide a successful submission; retrying elsewhere could duplicate payment |

The selected verifier is keyed to the exact payload plus requirements and
pinned for settlement. Ambiguous settlement errors are exposed in diagnostics
and resolved through the crash-safe proof reconciliation path.

## Facilitator Responsibilities

The facilitator:

- verifies that the payment payload satisfies the declared payment requirements
- submits validated payments to the blockchain
- monitors confirmation
- returns verification and settlement responses to the resource server

The resource server still decides when to perform the paid work and what to return.

## Testnet Verification Rule

For judging, lead with prepared work and the browser-wallet flow. Show the
commitment before signing, exact authorization, released work, decoded
`PAYMENT-RESPONSE`, and matching BaseScan transfer. The simulator is optional.
Do not claim protocol execution until the explorer record is confirmed.
