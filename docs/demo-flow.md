# Demo Flow

1. Open the site.
2. Submit a risky prompt.
3. The API returns `402 Payment Required`.
4. Click `Simulate x402 Payment`.
5. The request retries with `PAYMENT-SIGNATURE`.
6. The API returns a safety report and receipt.

For the real x402 integration, the simulated payment button will be replaced by wallet/facilitator-backed signing and settlement.

## Judge Checklist

- Paying customer: AI developers and agent builders.
- Pay-per-call model: one guardrail check per payment.
- x402 flow: challenge, sign, retry, settle, receipt.
- Receipt: included in response body and `PAYMENT-RESPONSE` header.
- Mainnet path: documented facilitator verification and settlement replacement.
