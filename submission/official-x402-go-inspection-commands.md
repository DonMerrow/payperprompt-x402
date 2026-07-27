# Official x402 Go Inspection Commands

Run these inside the official clone:

```bash
cd ~/Downloads/x402-official-readonly/go
```

Find public constructors and types:

```bash
grep -R "func New\\|func With\\|type .*Client\\|type .*Server\\|type Payment" -n . | head -160
```

Find verification and settlement functions:

```bash
grep -R "func .*Verify\\|func .*Settle\\|\\.Verify\\|\\.Settle" -n . | head -160
```

Find HTTP header constants:

```bash
grep -R "PAYMENT-REQUIRED\\|PAYMENT-SIGNATURE\\|PAYMENT-RESPONSE\\|Payment-Required\\|Payment-Signature\\|Payment-Response" -n . | head -160
```

Find examples:

```bash
find ../examples . -type f | grep -E '\\.(go|md)$' | sort | head -160
```

Find Go tests that show usage:

```bash
grep -R "Newx402ResourceServer\\|WithFacilitatorClient\\|WithSchemeServer\\|PaymentRequirements\\|PaymentPayload" -n . | head -200
```

Paste the output back into the chat. The goal is to wire the smallest safe official SDK path into `real-x402-go/internal/x402adapter/adapter.go`.
