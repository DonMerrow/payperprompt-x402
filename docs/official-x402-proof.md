# Official x402 Base Sepolia Proof

## Representative verified result — July 26, 2026

| Field | Verified value |
|---|---|
| Status | Successful and reconciled |
| Network | Base Sepolia (`eip155:84532`) |
| Token | Official Base Sepolia USDC |
| Amount | 0.02 USDC (`20000` atomic units) |
| Provider | Rapid Policy |
| Route | `guardrail-fast` |
| Payer | `0x826154a3d58aea3fbd2aa64aad424594ade927ef` |
| Merchant | `0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630` |
| Transaction | `0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5` |
| Work commitment | `fb1667e4c2b5c0de9fc8e2ae2d6f1c33f1496c0569d2475180147c069ab98338` |
| AI | Ollama `llama3.1:8b`, `ai_used: true` |

Explorer:

<https://sepolia.basescan.org/tx/0xfcf4b744479fed99b483d9bd665b67af014c6cb56e55e1492e1bee121dc4c3e5>

The successful response reported:

- `prepared_work_released: true`
- matching prepared and delivered SHA-256 commitments
- `payment_sent: true`
- `settled: true`
- successful live-chain verification
- independent Rust verification
- durable spend recording

At capture time, the dashboard reported 17 recorded official settlements, 17
independently verified proofs, and 0.34 USDC of Base Sepolia volume.

## What this proof establishes

The representative transaction establishes that:

1. a protected resource advertised an exact x402 requirement;
2. a recognized browser wallet authorized the required USDC amount;
3. official middleware accepted the paid retry;
4. the facilitator settled the payment;
5. the API returned the transaction and paid work;
6. the committed prepared work was released;
7. the ERC-20 transfer matched token, payer, merchant, and amount;
8. Rust independently accepted the proof structure.

It does not establish mainnet deployment, a formal security audit, guaranteed
AI correctness, or production readiness.

## Acceptance checklist

| Check | Required evidence |
|---|---|
| Challenge | Initial request returns HTTP 402 and `PAYMENT-REQUIRED` |
| Network | `eip155:84532` |
| Token | `0x036CbD53842c5426634e7929541eC2318f3dCF7e` |
| Wallet separation | Payer differs from `payTo` |
| Prepared work | Summary, validation, expiry, price, and commitment visible before signing |
| Settlement | `PAYMENT-RESPONSE.success` is true and includes a transaction |
| Work release | Released commitment equals prepared commitment |
| Explorer | BaseScan confirms successful exact USDC transfer |
| Live check | Go validates transaction and Transfer event through JSON-RPC |
| Independent check | Rust accepts cross-field proof consistency |

## Reproduce without another payment

```bash
cd ~/payperprompt-x402
./scripts/test-live-official-evidence.sh
```

This creates a fresh unpaid HTTP 402 challenge and re-verifies the recorded
transaction. It does not read a private key, request a signature, or send
another payment.

## Run one explicit browser-wallet payment

```bash
cd ~/payperprompt-x402
./scripts/start-judge-stack.sh
```

Then:

1. Open the local or printed HTTPS URL.
2. Connect a recognized disposable test wallet.
3. Confirm Base Sepolia, payer, merchant, and exact USDC price.
4. Prepare work and inspect the commitment.
5. Approve only the exact expected wallet authorization.
6. Open the returned BaseScan URL.
7. Run the no-payment live evidence action.

Never put a private key, seed phrase, production secret, or control token in
source, screenshots, prompts, documentation, or the merchant service.
