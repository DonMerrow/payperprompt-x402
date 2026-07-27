#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OFFICIAL_LOCAL="${X402_OFFICIAL_GO_DIR:-}"

cd "$PROJECT_DIR"

if [[ -z "$OFFICIAL_LOCAL" ]]; then
  for candidate in \
    "$PROJECT_DIR/.deps/x402-main/go" \
    "$HOME/Downloads/x402-main/go" \
    "$HOME/Downloads/x402-official-readonly/go"
  do
    if [[ -f "$candidate/go.mod" ]]; then
      OFFICIAL_LOCAL="$candidate"
      break
    fi
  done
fi

if [[ -z "$OFFICIAL_LOCAL" || ! -f "$OFFICIAL_LOCAL/go.mod" ]]; then
  DEPS_ROOT="$PROJECT_DIR/.deps"
  ARCHIVE="$DEPS_ROOT/x402-official-main.zip"
  SOURCE_ROOT="$DEPS_ROOT/x402-main"
  EXTRACT_ROOT="$(mktemp -d)"
  trap 'rm -rf "$EXTRACT_ROOT"' EXIT

  echo "Official x402 Go source is missing."
  echo "Downloading the x402 Foundation source ZIP — no Git and no npm."
  mkdir -p "$DEPS_ROOT"
  curl -fL \
    https://github.com/x402-foundation/x402/archive/refs/heads/main.zip \
    -o "$ARCHIVE"
  unzip -q "$ARCHIVE" -d "$EXTRACT_ROOT"
  if [[ ! -f "$EXTRACT_ROOT/x402-main/go/go.mod" ]]; then
    echo "Downloaded x402 source did not contain go/go.mod."
    exit 1
  fi
  if [[ -e "$SOURCE_ROOT" && ! -f "$SOURCE_ROOT/go/go.mod" ]]; then
    echo "Replacing incomplete internal dependency cache:"
    echo "  $SOURCE_ROOT"
    rm -rf -- "$SOURCE_ROOT"
  fi
  if [[ ! -e "$SOURCE_ROOT" ]]; then
    mv "$EXTRACT_ROOT/x402-main" "$SOURCE_ROOT"
  fi
  OFFICIAL_LOCAL="$SOURCE_ROOT/go"
fi

if [[ ! -f "$OFFICIAL_LOCAL/go.mod" ]]; then
  echo "Official x402 Go source was not found after download:"
  echo "  $OFFICIAL_LOCAL"
  exit 1
fi

echo "Using your read-only official x402 Go source:"
echo "  $OFFICIAL_LOCAL"
go mod edit -replace "github.com/x402-foundation/x402/go/v2=$OFFICIAL_LOCAL"

go mod tidy
go test ./...
mkdir -p bin
go build -buildvcs=false -trimpath -o bin/payperprompt ./cmd/official-client

echo
echo "Official Go lane is prepared, tests passed, and the CLI was rebuilt."
