# PayPerPrompt x402 v0.38.1

## Smart Contract Studio quality gate

v0.38.1 corrects issues found during the first real paid Solidity explanation.

### Deterministic pricing

- Explain Solidity contract: Local Guard, $0.01 USDC
- Generate Solidity tests: Local Guard, $0.01 USDC
- Solidity security audit: Deep Shield, $0.04 USDC
- Generate Solidity contract: Deep Shield, $0.04 USDC
- Repair Solidity contract: Deep Shield, $0.04 USDC

The explicit task selection is authoritative. The paid report receives the same
task type and route strategy used during planning, so a second model analysis
cannot contradict the settled service.

### Pre-settlement completeness

For Solidity explanation, audit, repair, and test work, the worker extracts
every constructor, receive function, fallback function, and named function from
the submitted source. The AI response must list and discuss every extracted
element. Undersized or incomplete work returns an error before official x402
settlement, so no USDC moves for an incomplete deliverable.

### Security accuracy

The worker now distinguishes the limited gas stipend and gas-brittleness of
Solidity `transfer` from call-based reentrancy. It is instructed not to invent
requirements, controls, vulnerabilities, or unsupported assurances.
