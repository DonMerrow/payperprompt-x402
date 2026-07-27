# PayPerPrompt x402 v0.39

## Prepared-Work Escrow

v0.39 changes the real browser-wallet flow from post-payment generation to
pre-payment completion:

1. Ollama analyzes the request and selects the priced route.
2. Ollama completes the work before any wallet prompt.
3. Existing deterministic quality and Solidity semantic checks run.
4. Go stores the full work in a ten-minute, one-time memory escrow.
5. The browser receives only a safe preview and SHA-256 commitment.
6. MetaMask authorizes the exact route price.
7. Official x402 middleware settles USDC on Base Sepolia.
8. The paid endpoint validates and releases the exact committed work.

The escrow is bound to the prompt, task type, route, configured payer wallet,
expiry, and settlement state. Reuse, mutation, expiration, or a second attempt
after an ambiguous settlement is blocked.

No seed phrase or private key is entered into the application. AI cannot
deploy contracts, sign transactions, or move funds.

## Verification

```bash
cd ~/payperprompt-x402/go-core
go test ./...

cd ~/payperprompt-x402/receipt-verifier-rust
cargo test

cd ~/payperprompt-x402/real-x402-go
./scripts/prepare-official.sh

cd ~/payperprompt-x402
./scripts/test-prepared-work-escrow.sh
./scripts/test-solidity-semantic-gate.sh
./scripts/test-browser-wallet-safety.sh
```

Then launch:

```bash
cd ~/payperprompt-x402
./scripts/start-judge-stack.sh
```
