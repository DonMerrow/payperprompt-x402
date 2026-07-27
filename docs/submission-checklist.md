# Final Submission Checklist

Use this checklist after the code-freeze test succeeds.

## Repository

- [ ] Public repository URL opens in a private browser window.
- [ ] `LICENSE` is visible and identifies the MIT licence.
- [ ] `README.md` and `START-HERE-OFFICIAL.md` render correctly.
- [ ] No `.env`, `.env.local`, private key, seed phrase, token, control token, build output,
      durable ledger, or dependency cache is committed.
- [ ] The frozen release ZIP and its SHA-256 checksum are attached or tagged.

Build the public-safe archive with:

```bash
cd ~/payperprompt-x402
./scripts/build-submission-package.sh
```

## Clean restart

- [ ] Ollama reports `llama3.1:8b` ready.
- [ ] `./scripts/test-submission-proof-kit.sh` passes.
- [ ] `./scripts/start-judge-stack.sh` starts ports 8082, 8084, and 8085.
- [ ] The local dashboard and temporary HTTPS URL both load.
- [ ] Preparing a task does not open MetaMask.
- [ ] One representative payment settles and releases the committed work.
- [ ] `./scripts/test-live-official-evidence.sh` passes without signing or
      sending another payment.

## Submission fields

- [ ] Copy the canonical story from `submission/DEVPOST-PASTE.md`.
- [ ] Replace `[CURRENT_PUBLIC_DEMO_URL]` and `[PUBLIC_REPOSITORY_URL]`.
- [ ] Paying customer: AI builders, agent platforms, and small SaaS teams.
- [ ] Business model: 0.01, 0.02, or 0.04 USDC per completed AI request.
- [ ] Explain Challenge → Sign → Retry → Settle.
- [ ] Show the receipt in the successful response.
- [ ] Describe Base Sepolia as testnet proof, not a mainnet deployment.
- [ ] Link the representative BaseScan transaction.
- [ ] Link the public source repository.
- [ ] Link the live demo or explain its temporary tunnel requirement.
- [ ] Enable the static `docs/` proof page so offline evidence remains available.
- [ ] Do not make a temporary `trycloudflare.com` URL the only Try It Out link.

## Media

- [ ] Hero screenshot: prepared-work proof before payment.
- [ ] Wallet screenshot: exact USDC authorization on Base Sepolia.
- [ ] Result screenshot: committed work released with transaction link.
- [ ] Evidence screenshot: BaseScan transfer and live verification.
- [ ] Developer screenshot: audit trail, CLI, or Go SDK.
- [ ] Demo video is under three minutes and follows `docs/video-script.md`.

## Final honesty check

- [ ] No claim of formal audit, guaranteed security, production deployment, or
      mainnet completion.
- [ ] No claim that AI deployed a contract, moved assets, or ran tests unless
      separately proven.
- [ ] No post-quantum or mining capability is described as implemented.
- [ ] Historical audit is bounded and hidden until requested.
