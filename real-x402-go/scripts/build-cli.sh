#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

mkdir -p bin
go build -buildvcs=false -trimpath -o bin/payperprompt ./cmd/official-client

echo "Built: $PROJECT_DIR/bin/payperprompt"
echo
bin/payperprompt help
