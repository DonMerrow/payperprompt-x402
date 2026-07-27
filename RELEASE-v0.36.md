# PayPerPrompt v0.36 — AI-Routed Browser Wallet x402

This release connects the local AI router directly to the real browser-wallet
x402 lane.

## Judge flow

1. Connect the configured disposable MetaMask, Coinbase Wallet, or Rabby
   account on Base Sepolia.
2. Enter a prompt.
3. Click **Ask AI to Plan Payment**.
4. Ollama selects:
   - Local Guard — lowest cost — $0.01 USDC
   - Rapid Policy — lowest latency — $0.02 USDC
   - Deep Shield — highest quality — $0.04 USDC
5. Go validates the selected service's fresh official HTTP 402 requirement and
   durable spend policy.
6. The page shows the route, strategy, merchant, and exact charge without
   requesting a signature.
7. The user explicitly approves the typed USDC authorization in the wallet.
8. Go re-fetches and validates the same selected route before official x402
   settlement.
9. The paid AI result, settlement receipt, BaseScan transaction, durable spend,
   live-chain proof, and independent Rust verification are recorded.

## Safety boundaries

- No automatic payment follows AI planning.
- No seed phrase or private key enters the page.
- Only the configured disposable payer is accepted.
- Only three fixed official routes and prices are accepted.
- The signed amount must match the fresh challenge and selected route exactly.
- Policy is checked again immediately before settlement.
- Ambiguous settlement outcomes remain fail-closed.
