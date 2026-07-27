# PayPerPrompt x402 v0.40.6.1

## Go assignment correction

- Replaces an invalid short declaration involving `product.Summary` with a
  declared local count and normal multi-value assignment.
- Scans the Go tree for other selector-based short declarations.
- Runtime grounding, prepared-work reuse, payment, and settlement behavior are
  unchanged.

## Origin-stable prepared work

- Removes bounded, source-provably false AI sentences before commitment.
- Discloses the number of deterministic corrections in the work caveats.
- Keeps all source-derived access, pause, accounting, transfer, and SafeMath
  findings visible.
- Reuses the same unused prepared-work commitment for an identical prompt,
  task, and wallet across localhost and Cloudflare.
- Refuses reuse after settlement begins, after use, after expiry, or for a
  different wallet.
- Prevents origin-dependent Ollama sampling from producing conflicting judge
  outcomes.

## Sentence-bounded semantic claims

- Fixes the full-review false positive observed on localhost.
- Prevents a correct SafeMath negation from matching text in a neighboring
  sentence.
- Continues rejecting claims that SafeMath provides or should be used as
  reentrancy protection.
- Avoids unnecessary corrective model calls that could exhaust the 90-second
  Ollama client timeout.
- Adds the exact combined-text regression exposed by the public and local runs.

## Corrective retry regression

- Updates the stale retry fixture exposed by the full Legion acceptance run.
- Omission-only drafts now correctly pass on the first attempt after
  deterministic grounding.
- Explicitly contradictory drafts still fail and trigger a corrected second
  Ollama attempt.
- Runtime payment, settlement, durable data, and proof behavior are unchanged.

## Deterministically grounded asynchronous review

- Extends the background deadline to cover the full bounded preparation
  pipeline.
- Adds authoritative Go-derived Solidity findings when Ollama omits facts that
  are directly provable from the submitted source.
- Grounds coverage in actual submitted functions and state transitions.
- Still rejects contradictory access-control, pause, accounting, SafeMath, and
  transfer claims.
- Keeps the prepared-work commitment stable when the paid response revalidates
  it after settlement.
- Separates transient polling failures from terminal preparation rejection and
  shows the exact connection error while safely retrying.

## Async prepared work and semantic-intent correction

- Runs long Ollama work preparation in a bounded background job.
- Returns `202 Accepted` immediately and exposes compact polling responses.
- Reuses duplicate active requests and limits concurrent preparation jobs.
- Keeps wallet signing and payment disabled until validated work is ready.
- Reports intermediary HTML and other non-JSON responses without misleading
  JavaScript parse errors.
- Accepts correct SafeMath negations and equivalent descriptions of ineffective
  pause enforcement.
- Retains fail-closed source coverage, access-control, ETH-accounting, and
  transfer-semantic checks.
- Adds regression tests for ready, duplicate, and failed asynchronous jobs.

## Submission proof kit

- Shows one current work-audit event by default.
- Loads bounded history only after **Show History** is requested.
- Adds a dependency-free reusable Go SDK and inspection example.
- Adds an on-page matrix for every formal submission expectation.
- Rejects observed non-compiling Foundry patterns before wallet approval.
- Requires imports, real caller switching, bounded fuzz inputs, and executable
  receive-path evidence in generated Solidity tests.
- Normalizes defensive smart-contract work so technical vulnerabilities are
  not misrepresented as malicious user intent.
- Adds three large-contract prompt fixtures and a twelve-element coverage
  regression test.
- Expands generated Smart Contract Studio ideas beyond single-vault examples.

Signing remains in the official x402 client or trusted browser wallet. No seed
phrase or private-key input was added.

This corrective release also supplies the required prepared-work identifier in
the unconfigured-payer regression test and removes an obsolete unused variable
that blocked the official Go lane from compiling.

Version 0.40.2 closes the general code-review routing gap for pasted Solidity.
It rejects the observed SharedWallet hallucinations and requires grounded
coverage of access control, pause enforcement, ETH accounting, transfer
semantics, and all submitted functions before payment can be signed.

The all-in-one acceptance script now runs official dependency preparation
before testing that lane.

Version 0.40.3 gives general Solidity code reviews a source-derived checklist
on the first generation attempt. Corrective retries receive that checklist,
the exact semantic failures, and the complete rejected draft. Regression tests
prove both a valid first attempt and a repaired second attempt.
