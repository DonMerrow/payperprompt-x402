# PayPerPrompt x402 v0.38.2

## Pre-settlement Solidity semantic proof

v0.38.2 closes the quality gap discovered in the real paid SimpleVault
explanation. Coverage alone is no longer enough: known contradictions are
rejected before x402 settlement.

### Deterministic source facts

The Go work engine derives and validates:

- explicit `msg.sender` authorization checks
- fixed `transfer` withdrawal destinations
- limited-gas `transfer` behavior versus call-based reentrancy
- `receive()` value flow into the contract balance
- Solidity `contract` syntax versus invented `class` declarations
- complete constructor, receive, fallback, and named-function coverage

### Corrective retry without charging

If the first Ollama draft contradicts a derived fact, the validator supplies
the failed check to one corrective generation attempt. If the corrected draft
still fails, the protected service returns HTTP 503 and official x402
settlement does not proceed.

### Visible semantic evidence

Successful Smart Contract Studio work now includes:

```json
{
  "semantic_validation": {
    "version": "solidity-semantic-v1",
    "valid": true,
    "checks": []
  }
}
```

The dashboard lists every passed deterministic semantic check beside the
function-coverage evidence.
