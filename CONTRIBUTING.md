# Contributing

Contributions are welcome when they preserve PayPerPrompt's central guarantee:
work is prepared, grounded, committed, and priced before wallet consent, and
only the matching committed result is released after verified settlement.

## Development checks

```bash
./scripts/test-submission-proof-kit.sh
```

For payment-boundary changes, also run the relevant prepared-work, browser
wallet, reconciliation, and live-evidence tests described in `docs/`.

## Pull requests

- describe user-visible behavior;
- identify any change to signing, settlement, policy, or proof verification;
- add regression coverage;
- keep simulator evidence separate from on-chain evidence;
- do not commit `.env.local`, runtime data, proof artifacts, dependencies, or
  build output;
- never include private keys, seed phrases, or production secrets.

See `docs/PUBLIC-REPOSITORY.md` for the repository map and public release flow.
