# PayPerPrompt v0.35.3 — Public-Chain Ledger Recovery

This repair handles a real settlement whose local browser-wallet ledger row was
lost before proof reconciliation completed.

## Recovery command

```bash
./scripts/recover-browser-wallet-proof.sh 0xTRANSACTION_HASH
```

The recovery path:

1. Accepts only a public transaction hash.
2. Restricts recovery to the configured disposable payer and exact $0.01 route.
3. Verifies transaction success and the exact payer-to-merchant Base Sepolia
   USDC transfer through JSON-RPC.
4. Requires independent Rust proof verification.
5. Recreates the durable Go ledger record exactly once.
6. Writes the latest official proof atomically.
7. Appends official proof history exactly once.

It does not connect a wallet, read a key, request a signature, or send a
payment.

For transaction
`0xaf05d3640cc9369c26eadcad1a030e8f3c2cdf10f78cae282cd61270e35f395d`,
the corrected totals are three official settlements and 0.04 USDC.
