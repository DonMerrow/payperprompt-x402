#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT_NAME="payperprompt-x402"
OUTPUT_DIR="${1:-$PROJECT_DIR/dist}"
STAGE_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$STAGE_DIR"
}
trap cleanup EXIT

mkdir -p "$OUTPUT_DIR"

rsync -a \
  --exclude '.git/' \
  --exclude '.env' \
  --exclude '.env.local' \
  --exclude '*.log' \
  --exclude 'dist/' \
  --exclude 'node_modules/' \
  --exclude '.npm-cache/' \
  --exclude 'go-core/data/' \
  --exclude 'real-x402-go/proof/' \
  --exclude 'real-x402-go/.deps/' \
  --exclude 'real-x402-go/bin/' \
  --exclude 'real-x402-rust/target/' \
  --exclude 'receipt-verifier-rust/target/' \
  "$PROJECT_DIR/" \
  "$STAGE_DIR/$PROJECT_NAME/"

if find "$STAGE_DIR/$PROJECT_NAME" -type f \
  \( -name '.env' -o -name '.env.local' -o -name '*.pem' -o -name '*.key' -o -name '*private*key*' \) \
  -print -quit | grep -q .; then
  echo "Unsafe secret-like file found in submission staging."
  exit 1
fi

ARCHIVE="$OUTPUT_DIR/payperprompt-x402-v0.41.0.1-public-repository.zip"
(
  cd "$STAGE_DIR"
  zip -qr "$ARCHIVE" "$PROJECT_NAME"
)
(
  cd "$OUTPUT_DIR"
  sha256sum "$(basename "$ARCHIVE")" > "$(basename "$ARCHIVE").sha256"
)

echo "Submission package: $ARCHIVE"
echo "Checksum: $ARCHIVE.sha256"
