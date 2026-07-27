# PayPerPrompt v0.35.1 — Browser Settlement Reconciliation

This release reconciles a successful browser-wallet payment into every official
evidence layer without making another payment.

## Recovery path

`POST /api/official/reconcile-browser-wallet`:

1. Finds the latest browser-wallet settlement in the durable Go ledger.
2. Reconstructs its exact Base Sepolia USDC requirement.
3. Queries the transaction receipt through Base Sepolia JSON-RPC.
4. Requires transaction success and an exact payer-to-merchant USDC Transfer.
5. Sends the reconstructed proof to the independent Rust verifier.
6. Atomically replaces the latest official proof.
7. Appends the proof to official history only when the transaction is new.

The browser calls this recovery path when the dashboard loads. It is
idempotent, does not request a signature, and does not send a payment.

For the existing `0xaf05…395d` browser settlement, the expected corrected
dashboard totals are:

- 3 official settlements
- 3 independently verified settlements
- 0.04 USDC official Base Sepolia volume
- $0.04 durable payer spending
