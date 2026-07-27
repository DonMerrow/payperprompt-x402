# PayPerPrompt v0.32 — Crash-Safe Settlement Recovery

This release closes the failure window between a successful Base Sepolia x402
payment and durable Go policy accounting.

## What changed

- Official proofs now preserve the Go policy authorization ID.
- Once the official client attempts payment, it never releases reserved budget
  merely because the response, live check, or ledger commit fails.
- A protected reconciliation endpoint verifies the latest official proof live
  on Base Sepolia and commits the matching reservation exactly once.
- Expired or previously released reservations can be recovered only when the
  proof's payer, resource, amount, authorization, transaction, token contract,
  Transfer event, and live receipt are consistent.
- Repeating reconciliation is idempotent and never duplicates spend.
- The recovery command does not load a private key, sign, or send payment.
- The dashboard exposes whether durable settlement recovery is required.

## No-payment validation

```bash
cd ~/payperprompt-x402/real-x402-go
./scripts/test-settlement-recovery.sh
```

Expected result:

```text
SETTLEMENT RECOVERY TEST PASSED
A verified payment can be committed after interruption, exactly once.
The recovery path does not require a private key and does not send another payment.
```

## Operator recovery

Run this only when an official payment settled but durable accounting reported
an error:

```bash
cd ~/payperprompt-x402/real-x402-go
./scripts/reconcile-official-settlement.sh
```

The command reads the preserved proof and local policy-control token. It never
reads `EVM_PRIVATE_KEY`.
