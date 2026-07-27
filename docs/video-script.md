# Demo Video Script — Final Judge Cut

Target length: 120 to 150 seconds. Make no payment before recording unless the
prepared proof is credible.

## 0:00–0:18 — Problem and customer

“PayPerPrompt is a prepared-work x402 marketplace for AI builders, agent
platforms, and small SaaS teams. AI completes and validates the work first.
The buyer sees proof and the exact price before deciding whether to pay.”

Show the paying customer and 0.01, 0.02, and 0.04 USDC pricing.

## 0:18–0:45 — Prepare without payment

Choose a realistic task and select **Ask AI to Plan Payment**.

“Ollama identifies the task and route. Go grounds the result and keeps the full
deliverable hidden. The page shows its summary, coverage, expiry, semantic
checks, exact price, and SHA-256 commitment. MetaMask has not opened and
nothing has been charged.”

Show the prepared-work proof and commitment.

## 0:45–1:10 — Exact wallet authorization

Select the payment button.

“The protected endpoint returns its official HTTP 402 requirement. MetaMask
shows the exact Base Sepolia USDC authorization. The buyer explicitly
approves that amount.”

Show the wallet network, token, amount, and merchant. Do not show secrets.

## 1:10–1:35 — Settle and release

“The paid request retries, official x402 middleware settles, and only the
previously committed work is released.”

Show:

- completed deliverable;
- matching prepared and delivered commitments;
- transaction hash and BaseScan link;
- route and exact amount;
- current durable audit event.

## 1:35–1:58 — Independent proof

Run **Live x402 Evidence**.

“This creates a fresh unpaid HTTP 402 challenge, checks the recorded transfer
through Base Sepolia JSON-RPC, and asks the independent Rust service to verify
proof consistency. It does not sign or send another payment.”

Show the five green proof stages and BaseScan transaction.

## 1:58–2:20 — Developer tooling and close

Show the submission matrix and Go SDK/CLI commands.

“PayPerPrompt covers the complete x402 lifecycle and the tooling around it:
AI prepares, Go controls, the wallet consents, x402 settles, and Rust verifies.
Failed work cannot charge, ambiguous settlement cannot be blindly retried, and
the successful response contains the receipt.”

End on the repository URL and project name.
