# PayPerPrompt x402 Lab v0.24.1

## Outcome

v0.24.1 adds a safe, official Base Sepolia path without disturbing the proven
Go/Rust sandbox.

## Added

- official x402 Go middleware protecting the paid AI endpoint
- grounded Ollama analysis after successful payment verification
- no-spend payer preflight
- automatic payer-address derivation without printing the key
- refusal when payer and merchant are the same wallet
- strict Base Sepolia network and USDC contract validation
- explicit `pay` command as the only spending path
- public settlement proof JSON with BaseScan link
- official-mode dashboard on port 8082
- one-command preparation and launch scripts
- Go tests for AI normalization and unsupported-claim grounding

## Security

- the merchant server unsets `EVM_PRIVATE_KEY`
- the payer key remains in one terminal environment
- no private key is written into the proof
- documentation no longer treats terminal output as chain verification

## Verification Status

- local Go/Rust sandbox: tested
- official lane compiled and tests passed on Legion
- official preflight passed with distinct payer and merchant wallets
- official payment moved 0.01 Base Sepolia USDC
- Ollama `llama3.1:8b` produced the paid result
- BaseScan transaction confirmed:
  `0x03c3b1be51cedd392add099d95571e6ac4ec220e012b9670ee8dbd8b496387cb`
