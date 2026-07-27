# PayPerPrompt x402 v0.40.7.1

## Dotted-identifier sanitizer correction

The v0.40.7 regression exposed a sentence-boundary edge case:
`address(this).balance` contains a period, so the generic sentence splitter
could divide the false claim before deterministic grounding recognized it.

This release:

- removes bounded false source-usage sentences before generic splitting;
- handles invented `address(this).balance` and `SafeMath` usage;
- handles claims that limited-gas `transfer` can itself cause reentrancy;
- preserves correct deterministic findings that compare manual accounting with
  `address(this).balance`;
- keeps the regression's one-attempt requirement, preventing an unnecessary
  second Ollama call.

No payment, wallet, settlement, policy, or secret-handling behavior changed.
