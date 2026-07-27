# PayPerPrompt Final Submission Runbook

## Code-freeze acceptance gate

Run:

```bash
cd ~/payperprompt-x402
./scripts/test-submission-proof-kit.sh
./scripts/test-live-official-evidence.sh
```

Do not record a final demo until both commands pass.

## Large-contract proof

The three prompts in `examples/large-contract-prompts/` cover generation,
multi-contract auditing, and Foundry test creation.

For each prompt:

1. Choose the matching Smart Contract Studio work type.
2. Paste the fixture.
3. Select **Ask AI to Plan Payment**.
4. Confirm the prepared-work summary, SHA-256 commitment, complete coverage,
   and deterministic checks.
5. Do not pay every preparation. Pay one representative large-contract job
   only after the prepared proof is credible.

Generated Foundry work must not contain `vm.deposit`, a named `receive()` call,
an invented `balance()` method, undefined `random()`, missing imports, unbounded
fuzz input, or an authorization test that never switches callers.

## Ninety-second judge path

1. Open UI v0.41.0.1 and identify the paying customer and exact per-call prices.
2. Connect MetaMask on Base Sepolia.
3. Choose a realistic AI task and prepare it without opening the wallet.
4. Show the hidden-work commitment, coverage, validation, and exact price.
5. Approve one wallet signature.
6. Show the released deliverable and transaction link.
7. Show the current audit event, then explicitly open bounded history.
8. Run **Live x402 Evidence**.
9. Show the submission matrix and Track 3 Go SDK.

## Four-day finish

### Day 1 — Freeze behavior

- Run all Go, Rust, SDK, browser-wallet, and proof tests.
- Prepare all three large-contract fixtures.
- Correct any rejected output without weakening the gate.
- Make one final representative wallet payment only if needed.

### Day 2 — Public presentation

- Replace the temporary quick tunnel with a stable named tunnel or deployment.
- Confirm the public page, BaseScan links, mobile layout, and no-secret checks.
- Update screenshots and README commands from the frozen build.

### Day 3 — Submission media

- Record one uninterrupted judge flow.
- Keep the video under three minutes.
- Show the transaction, paid deliverable, audit, live verification, and SDK.
- Finish the Devpost text and repository links.

### Day 4 — Contingency

- Re-run the acceptance gate on a clean restart.
- Verify the stable URL from another device.
- Submit early enough to correct upload or link failures.

## Mainnet-readiness statement

The demonstrated network is Base Sepolia. The architecture separates network
and token configuration, policy control, signing, settlement, reconciliation,
and independent verification. A production launch still requires a production
facilitator/network configuration, managed secret storage, stable hosting,
monitoring, incident response, and an external security review.

## Frozen representative proof

- Provider: Rapid Policy
- Route: `guardrail-fast`
- Amount: 0.02 USDC
- Transaction:
  `0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5`
- Prepared/released commitment:
  `fb1667e4c2b5c0de9fc8e2ae2d6f1c33f1496c0569d2475180147c069ab98338`

Do not spend again merely to obtain a newer transaction. Use the no-payment
live evidence action unless a new representative payment is genuinely needed.
