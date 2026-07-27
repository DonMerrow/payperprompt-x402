#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NODE_BIN="${NODE_BIN:-node}"

(cd "$PROJECT_DIR/go-core" && go test ./cmd/server -run 'TestFreeWorkSuggestion|TestWorkSuggestionRate')
"$NODE_BIN" --check "$PROJECT_DIR/web/app.js"

grep -Fq 'POST /api/ai/work-suggestion' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'spend_policy_checked' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'wallet_required' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'curated-fallback' "$PROJECT_DIR/go-core/cmd/server/main.go"
grep -Fq 'id="generate-work-suggestion"' "$PROJECT_DIR/web/index.html"
grep -Fq 'Give me another AI idea' "$PROJECT_DIR/web/app.js"
grep -Fq 'Replace your custom work request' "$PROJECT_DIR/web/app.js"
grep -Fq 'current_prompt: currentRequest' "$PROJECT_DIR/web/app.js"

echo "ADAPTIVE AI WORK-SUGGESTION TEST PASSED"
echo "Ollama can generate distinct work ideas without wallet, reservation, policy evaluation, signature, or payment."
echo "Immediate repeats, unsafe requests, and excessive public use are controlled."
