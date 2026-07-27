#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

if [[ ! -x ./bin/payperprompt ]]; then
  ./scripts/build-cli.sh >/dev/null
fi

./bin/payperprompt facilitators
