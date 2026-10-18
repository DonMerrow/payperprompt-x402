#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

# The original v2.0.0 module tag was unavailable. Use a verified source
# snapshot, not changing main. Update these values, go.mod and README together.
X402_COMMIT=90688e52e58ae9185f2860988bd2c46d2801ceda
X402_SHA256=2204efbdb16f674e8189d5bbd202c4eacd76da9e78b31c637028db86cef9599e
DEPS_ROOT="$PROJECT_DIR/.deps"
ARCHIVE="$DEPS_ROOT/x402-$X402_COMMIT.zip"
SOURCE_ROOT="$DEPS_ROOT/x402-$X402_COMMIT"

if [[ -n "${X402_OFFICIAL_GO_DIR:-}" ]]; then
  echo "X402_OFFICIAL_GO_DIR is no longer accepted: preparation verifies the pinned archive only." >&2
  exit 1
fi
for required in curl unzip sha256sum go; do
  command -v "$required" >/dev/null || { echo "Required command missing: $required" >&2; exit 1; }
done
mkdir -p "$DEPS_ROOT"
WORK_DIR="$(mktemp -d "$DEPS_ROOT/.prepare-XXXXXX")"
trap 'rm -rf -- "$WORK_DIR"' EXIT
verify_archive() {
  printf '%s  %s\n' "$X402_SHA256" "$1" | sha256sum --check --status
}
if [[ ! -f "$ARCHIVE" ]]; then
  echo "Downloading official x402 source at commit $X402_COMMIT"
  curl --fail --location --proto '=https' --proto-redir '=https' \
    --connect-timeout 15 --max-time 300 \
    "https://github.com/x402-foundation/x402/archive/$X402_COMMIT.zip" \
    -o "$WORK_DIR/source.zip"
  if ! verify_archive "$WORK_DIR/source.zip"; then
    echo "x402 archive checksum mismatch; nothing will be extracted or built." >&2
    exit 1
  fi
  mv "$WORK_DIR/source.zip" "$ARCHIVE"
fi
# Cached downloads must pass the same verification before extraction/build.
if ! verify_archive "$ARCHIVE"; then
  echo "Cached x402 archive checksum mismatch: $ARCHIVE" >&2
  echo "Remove that archive and rerun to download the pinned copy." >&2
  exit 1
fi
unzip -q "$ARCHIVE" -d "$WORK_DIR"
if [[ ! -f "$WORK_DIR/x402-$X402_COMMIT/go/go.mod" ]]; then
  echo "Pinned archive is missing the expected Go module." >&2
  exit 1
fi
# Re-extract verified bytes rather than trusting a previously edited cache.
rm -rf -- "$SOURCE_ROOT"
mv "$WORK_DIR/x402-$X402_COMMIT" "$SOURCE_ROOT"
echo "Using verified official source: $SOURCE_ROOT/go"
go mod edit -replace "github.com/x402-foundation/x402/go/v2=./.deps/x402-$X402_COMMIT/go"
go mod tidy
go mod verify
go test ./...
mkdir -p bin
go build -buildvcs=false -trimpath -o bin/payperprompt ./cmd/official-client
echo "Official Go lane prepared: checksums verified, tests passed, CLI rebuilt."
