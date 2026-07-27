# PayPerPrompt x402 v0.39.3

## Draft-Aware Corrective Generation

The local 8B model could repeat incomplete Solidity test output because the
correction prompt included the validation error but not the rejected draft.

v0.39.3:

- includes the rejected JSON work product in each corrective request
- names every source element requiring executable test evidence
- allows two corrective attempts after the initial draft
- requires a complete replacement instead of a prose explanation
- still blocks MetaMask unless deterministic validation succeeds

Failed preparation remains free and displays no active charge.
