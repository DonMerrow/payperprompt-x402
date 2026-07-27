# Browser Wallet Configuration

PayPerPrompt uses a browser wallet for explicit Base Sepolia test payments.
The application never needs, requests, or stores a private key, seed phrase,
wallet password, or recovery phrase.

## Safe configuration

From the project root:

```bash
./scripts/configure-public-demo.sh
```

The script asks for:

1. a disposable Base Sepolia payer **public address**;
2. a different merchant **public address**;
3. a maximum USDC amount per call;
4. a daily USDC spending limit.

It writes `.env.local` with owner-only permissions. That file is ignored by Git
and excluded from release archives.

Manual configuration is also possible:

```bash
cp .env.example .env.local
```

Then replace the two address placeholders. Keep the network and official Base
Sepolia USDC contract unchanged for the demonstrated flow.

## Wallet requirements

- MetaMask, Coinbase Wallet, or Rabby;
- the configured disposable payer account selected;
- Base Sepolia selected;
- enough Base Sepolia ETH for gas;
- enough Base Sepolia USDC for the exact displayed charge.

The dashboard loads the allowed public payer from `/api/config/public`. A
different account is rejected before payment preparation. This allowlist is a
testnet safety interlock, not production customer onboarding.

## What never belongs in this repository

- private keys;
- seed or recovery phrases;
- wallet passwords;
- production API tokens;
- policy control tokens;
- browser extension exports;
- real customer secrets.

If a private key or seed phrase has ever been pasted into a file, terminal
history, issue, screenshot, or prompt, treat that wallet as compromised and
move assets to a new wallet. Removing the text from Git does not make an
exposed key safe again.

## Production path

A production service should replace the single configured payer with
authenticated wallet registration, tenant isolation, database-backed policy,
rate limiting, monitoring, managed application secrets, and explicit mainnet
controls. Browser signing should still remain in the user's wallet.
