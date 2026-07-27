# PayPerPrompt v0.35 — Trusted Browser Wallet x402

## Outcome

The judge page can now connect a standard browser wallet and make a new real
0.01 USDC x402 payment on Base Sepolia. This is separate from the optional
fake-fund simulator and from the read-only recorded-proof demonstration.

## Supported wallets

- MetaMask (primary)
- Coinbase Wallet
- Rabby

Wallets are discovered with EIP-6963 and used through EIP-1193. The interface
allowlists those three provider identities and does not display arbitrary
injected extensions.

## Safety rules

- Disposable Base Sepolia test account only.
- The public server accepts only the configured `PAYER_ADDRESS`.
- No seed phrase or private-key field exists.
- No WalletConnect relay or third-party wallet directory is used.
- Connecting and preparing do not sign or spend.
- The final button states the exact 0.01 USDC charge.
- The browser asks the wallet to sign EIP-712 typed data.
- Go compares the signed payload with a newly fetched official HTTP 402.
- Go enforces the durable payer policy before sending the payment.
- Ambiguous network outcomes stop with reconciliation guidance.

## Judge flow

1. Start the one-command stack.
2. Open the Cloudflare URL in the browser profile containing the trusted wallet.
3. Select MetaMask, Coinbase Wallet, or Rabby.
4. Confirm the disposable payer and Base Sepolia balances.
5. Click `Prepare 0.01 USDC Payment`.
6. Inspect the decoded charge.
7. Click `Pay 0.01 USDC & Run AI`.
8. Approve the typed authorization in the wallet.
9. Open the returned transaction on BaseScan.

The successful response contains the paid AI result, official settlement
response, transaction hash, exact amount, payer, and durable spend total.
