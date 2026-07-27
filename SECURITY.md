# Security Policy

PayPerPrompt is a hackathon MVP demonstrated on Base Sepolia. It is not
formally audited and should not be used with mainnet funds.

## Never disclose

Do not submit private keys, seed phrases, wallet passwords, policy control
tokens, production credentials, or customer secrets in issues, pull requests,
prompts, screenshots, logs, or configuration committed to Git.

The supported browser-wallet flow asks only for public addresses during local
configuration. Signing remains inside MetaMask, Coinbase Wallet, or Rabby.

## Reporting a vulnerability

Use a private security advisory in the public GitHub repository. Include:

- affected version and component;
- reproducible steps using testnet or fake funds;
- expected and observed behavior;
- payment, replay, authorization, or data exposure impact;
- a suggested mitigation if available.

Do not test with another person's wallet or assets. Do not publish an active
secret or exploit before a fix is available.

## Security boundaries

- Base Sepolia only in the demonstrated flow.
- One configured disposable payer is accepted by the public demo.
- Preparation and validation must pass before wallet signing.
- A prepared result is released once and must match its commitment.
- An ambiguous settlement enters reconciliation instead of being paid again.
- Solidity assistance does not deploy contracts or handle secrets.

Production use requires external review, authenticated multi-tenant wallet
onboarding, managed secrets, stable storage, rate limiting, monitoring, and
incident response.
