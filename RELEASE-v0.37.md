# PayPerPrompt x402 v0.37

## Paid AI workspace

v0.37 turns the working x402 payment and verification stack into a useful AI application.

- Users can request general assistance, code review, bug summarization, meeting action items, document analysis, or prompt-security review.
- Ollama identifies the task and selects the cost, latency, or quality route before any wallet signature.
- If Ollama does not answer during planning, MetaMask payment remains disabled.
- The browser shows the exact provider and USDC charge before opening the wallet.
- After official x402 settlement, Ollama returns a structured completed work product with a title, summary, deliverable, action items, and caveats.
- The completed work is displayed separately from the raw payment and proof evidence.
- Go still enforces the durable payer policy and records settlement.
- Rust and live Base Sepolia JSON-RPC verification remain unchanged.

No seed phrase or private key is requested by the browser application.

## Demonstration

1. Connect the disposable Base Sepolia wallet.
2. Select a task type or allow automatic detection.
3. Enter a useful work request.
4. Ask AI to plan the payment.
5. Inspect the selected task, provider, route, and exact USDC amount.
6. Approve the visible MetaMask signature.
7. Read the completed AI deliverable.
8. Inspect the verified transaction and independent proof.
