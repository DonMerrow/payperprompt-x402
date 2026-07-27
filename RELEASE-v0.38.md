# PayPerPrompt x402 v0.38

## Smart Contract Studio

v0.38 adds dedicated paid Solidity work to the browser-wallet AI workspace.

### Operations

- Smart contract security audit
- Smart contract generation
- Smart contract explanation
- Foundry or Hardhat test generation
- Smart contract repair with change explanations

Security audits and new contract generation automatically use the enhanced
Deep Shield route at $0.04 USDC. Other contract operations retain the
AI-planned cost, speed, or quality route.

### Safety boundary

The studio works only with requirements and pasted source text. It does not:

- deploy contracts
- request or process private keys or seed phrases
- sign wallet transactions
- move assets
- claim that generated code is formally audited

Generated and repaired contracts include assumptions, remaining test
requirements, and security caveats. Human review remains required before any
deployment.

### Classification correction

Normal writing requests that mention x402, micropayments, or AI purchasing are
no longer treated as prompt-security work merely because of those product
terms. Explicit task selections always override automatic task detection.
