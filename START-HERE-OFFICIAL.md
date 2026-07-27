# Start Here — PayPerPrompt x402 v0.41.0.1

The primary judge flow uses MetaMask, Coinbase Wallet, or Rabby on Base Sepolia.
The page never asks for a seed phrase or private key.

## Requirements

- Ollama running with `llama3.1:8b`
- the configured disposable payer account in a recognized browser wallet
- Base Sepolia ETH for gas
- Base Sepolia USDC
- a distinct public merchant address

## 1 — Configure public testnet addresses

Run the guided setup:

```bash
cd ~/payperprompt-x402
./scripts/configure-public-demo.sh
```

Enter only a disposable payer's **public address** and a distinct merchant's
**public address**. The script never requests a private key, seed phrase,
wallet password, or production secret. It creates an ignored `.env.local`
file. Browser-wallet signing stays inside MetaMask, Coinbase Wallet, or Rabby.

See `docs/WALLET-CONFIGURATION.md` for manual configuration and safety details.

## 2 — Run the acceptance gate

```bash
cd ~/payperprompt-x402
./scripts/test-submission-proof-kit.sh
```

Stop if any Go, Rust, SDK, wallet-safety, or grounding test fails.

## 3 — Start the complete judge stack

```bash
cd ~/payperprompt-x402
./scripts/start-judge-stack.sh
```

Keep that terminal open. Use the new Cloudflare URL printed by the launcher.

## 4 — Prepare work without payment

1. Connect the disposable browser wallet.
2. Confirm Base Sepolia.
3. Choose a work type and review the request.
4. Select **Ask AI to Plan Payment**.
5. Inspect the route, exact price, work summary, coverage, validation checks,
   expiry, and SHA-256 commitment.

If preparation fails, no wallet signature or payment should occur.

## 5 — Make one explicit wallet payment

Select the payment button only when the prepared proof is credible. Confirm the
exact USDC authorization in the wallet.

Success must display:

- completed committed work;
- official x402 settlement;
- transaction hash and BaseScan link;
- payer, merchant, network, token, and exact amount;
- durable audit event;
- independent live-chain and Rust verification.

## 6 — Verify the final evidence

Open BaseScan and run **Live x402 Evidence**. Confirm:

- successful Base Sepolia transaction;
- official Base Sepolia USDC contract;
- expected payer and distinct merchant;
- exact atomic amount;
- fresh official HTTP 402 challenge;
- independent Rust proof acceptance.

Historical audit metadata is hidden by default. Select **Show History** only
when a judge asks to inspect earlier preparations or settlements.

See `docs/final-submission-runbook.md` for the large-contract proof, video path,
and final completion schedule. Use `docs/submission-checklist.md` before
publishing the repository or submitting the project.

The quick Cloudflare URL works only while the host computer and launcher remain
online. For review after the host is offline, publish `docs/` with GitHub Pages
and submit the repository, video, static proof page, and BaseScan transaction.
See `docs/HOSTING-AND-JUDGE-AVAILABILITY.md`.

Build a clean public archive only after the acceptance gate passes:

```bash
./scripts/build-submission-package.sh
```
