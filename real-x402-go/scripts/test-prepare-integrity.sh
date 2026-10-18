#!/usr/bin/env bash
# Isolated rejection tests: no real dependency download, build or payment.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf -- "$TEST_ROOT"' EXIT
mkdir -p "$TEST_ROOT/scripts" "$TEST_ROOT/.deps" "$TEST_ROOT/mock-bin"
cp "$SCRIPT_DIR/prepare-official.sh" "$TEST_ROOT/scripts/"
cat > "$TEST_ROOT/mock-bin/go" <<'MOCK'
#!/usr/bin/env bash
printf 'unexpected Go execution\n' >> "$PREPARE_TEST_MARKER"
exit 99
MOCK
cat > "$TEST_ROOT/mock-bin/curl" <<'MOCK'
#!/usr/bin/env bash
while [[ $# -gt 0 ]]; do
  if [[ "$1" == -o ]]; then printf 'untrusted download' > "$2"; exit 0; fi
  shift
done
exit 98
MOCK
chmod +x "$TEST_ROOT/mock-bin/"*
export PATH="$TEST_ROOT/mock-bin:$PATH"
export PREPARE_TEST_MARKER="$TEST_ROOT/go-called"
unset X402_OFFICIAL_GO_DIR
expect_failure() {
  if bash "$TEST_ROOT/scripts/prepare-official.sh" > "$TEST_ROOT/output" 2>&1; then
    echo "FAIL: accepted $1" >&2; exit 1
  fi
  grep -q "$2" "$TEST_ROOT/output"
  test ! -e "$PREPARE_TEST_MARKER"
  test -z "$(find "$TEST_ROOT/.deps" -maxdepth 1 -name '.prepare-*' -print -quit)"
  echo "PASS: $1 rejected before Go execution"
}
expect_failure 'tampered download' 'archive checksum mismatch'
ARCHIVE_NAME="$(sed -n 's/^X402_COMMIT=//p' "$SCRIPT_DIR/prepare-official.sh")"
printf 'untrusted cache' > "$TEST_ROOT/.deps/x402-$ARCHIVE_NAME.zip"
expect_failure 'tampered cached archive' 'Cached x402 archive checksum mismatch'
export X402_OFFICIAL_GO_DIR="$TEST_ROOT/unverified"
expect_failure 'unverified local override' 'no longer accepted'
