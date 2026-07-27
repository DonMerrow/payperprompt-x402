# PayPerPrompt x402 v0.41.0.1

## Portable public-config test correction

The public configuration endpoint already returned the correct safe JSON. Its
new regression test incorrectly compared the pretty-printed response against
compact JSON fragments. The test now decodes the response and checks typed
field values, making it independent of harmless whitespace formatting.

Runtime wallet configuration, signing, settlement, and proof behavior are
unchanged. No payment was signed or sent for this correction.
