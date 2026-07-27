# PayPerPrompt x402 v0.39.1

## Canonical Prepared-Work Commitment Fix

v0.39 correctly blocked payment when the prepared-work commitment did not
match, but the transport compared two harmlessly different JSON
serializations. This caused valid prepared work to fail before MetaMask opened.

v0.39.1 transports the canonical work bytes as Base64 and hashes those exact
bytes in Go. The official paid endpoint independently reconstructs and checks
the same work product before settlement.

No payment behavior, wallet permissions, pricing, private-key handling, or
directory layout changed.
