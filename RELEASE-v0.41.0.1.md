# PayPerPrompt x402 v0.41.0.1

## Portable public-config test correction

The public configuration endpoint already returned the correct safe JSON. Its
new regression test incorrectly compared the pretty-printed response against
compact JSON fragments. The test now decodes the response and checks typed
field values, making it independent of harmless whitespace formatting.

Runtime wallet configuration, signing, settlement, and proof behavior are
unchanged. No payment was signed or sent for this correction.

## Qwen and prepared-work reliability update

The default local work model is now `qwen3-coder:30b`. Active configuration,
startup scripts, examples, and documentation were updated consistently.
Historical settlement proofs retain their original model metadata.

Solidity preparation now gives exact source-element obligations to explanation
and test tasks. Corrective retries regenerate from the original source instead
of repeating rejected drafts. The default Foundry request includes executable
fuzz bounds, and generated work suggestions cannot demand a conflicting
top-level JSON, XML, or YAML response schema.

The deterministic quality gate remains strict and continues to reject
incomplete or contradictory work before wallet approval. No payment was signed
or sent for these changes. The project favicon is now served by the web UI.
