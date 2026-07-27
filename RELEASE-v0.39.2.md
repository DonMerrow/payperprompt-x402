# PayPerPrompt x402 v0.39.2

## Test-Aware Solidity Coverage

The pre-settlement semantic gate now treats generated tests as executable code
rather than prose:

- `new ContractName(...)` proves constructor coverage
- a value-bearing low-level call or `sendTransaction` proves receive coverage
- `.functionName(...)` proves named-function coverage
- missing executable evidence still blocks wallet signing

Go adds only evidence-backed elements to the returned coverage list. It does
not blindly trust AI-supplied coverage metadata.

The wallet panel now clears stale prices when a work type changes and displays
`No charge · preparation failed` after a preflight rejection.

No wallet permissions, settlement rules, private-key handling, or directory
layout changed.
