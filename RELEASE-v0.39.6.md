# PayPerPrompt x402 v0.39.6

## Grounded work audit

- Canonicalizes checksum and lowercase EVM addresses into one durable policy
  identity, so the configured $2 daily limit is found in either representation.
  Conflicting duplicate policies fail toward the stricter limit.
- Retries malformed or truncated Ollama work JSON up to four bounded attempts
  with a larger generation budget.
- Applies deterministic grounding checks to general AI work as well as
  Solidity work.
- Rejects unsupported code-coverage labels and claims of deployments, test
  runs, formal audits, or guaranteed security.
- Requires requested architecture diagrams and interface mockups to be
  represented in the prepared Markdown before wallet approval.
- Records preparation, policy, and settlement metadata in a durable official
  work audit.
- Shows that audit on the main page with BaseScan links for settled payments.
- Separates official transactions from the fake local simulator in exported
  evidence.
- Removes the redundant **Use suggested request** button while preserving
  custom requests and free AI-generated alternatives.

No seed phrase or private key input is present. Failed preparation and policy
checks do not open MetaMask or charge USDC.
