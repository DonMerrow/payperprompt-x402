# PayPerPrompt v0.33 — Safe Facilitator Failover

This release adds ordered facilitator resilience without introducing a
duplicate-payment path.

## Safety model

- `GetSupported` may try the next configured facilitator.
- `Verify` may try the next configured facilitator because verification does
  not transfer funds.
- A successful verifier is pinned to the exact payment payload and requirements.
- `Settle` is attempted once against that pinned facilitator.
- Any settlement transport error is classified as
  `unknown_requires_reconciliation`.
- The server never sends the same settlement to a second facilitator.
- The official client already preserves the durable authorization after a
  payment attempt, so the existing proof-reconciliation path remains the only
  safe recovery action.

## Operator diagnostics

Configure an ordered list when more than one compatible facilitator is
available:

```bash
export X402_FACILITATOR_URLS="https://primary.example,https://backup.example"
```

With the official server running, probe support and health without loading a
private key:

```bash
cd ~/payperprompt-x402/real-x402-go
./scripts/facilitator-status.sh
```

The response reports endpoint order, health, operation counts, active
facilitator, last settlement state, and the no-automatic-retry policy.

## No-payment validation

```bash
cd ~/payperprompt-x402/real-x402-go
./scripts/test-facilitator-failover.sh
```

Expected result:

```text
FACILITATOR FAILOVER SAFETY TEST PASSED
Supported/verify operations can use the next healthy endpoint.
Settlement is attempted once; ambiguous outcomes require reconciliation.
No private key was read. No payment was signed or sent.
```

