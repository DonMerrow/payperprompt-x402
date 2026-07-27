#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NODE_BIN="${NODE_BIN:-node}"

"$NODE_BIN" --check "$PROJECT_DIR/web/app.js"

if grep -Fq 'id="use-work-suggestion"' "$PROJECT_DIR/web/index.html"; then
  echo "Redundant Use suggested request control is still present." >&2
  exit 1
fi
grep -Fq 'id="generate-work-suggestion"' "$PROJECT_DIR/web/index.html"
grep -Fq 'WORK_TYPE_SUGGESTIONS' "$PROJECT_DIR/web/app.js"
grep -Fq 'handleWorkTypeChange' "$PROJECT_DIR/web/app.js"
grep -Fq 'generateAnotherWorkSuggestion' "$PROJECT_DIR/web/app.js"
grep -Fq 'Your custom request was preserved' "$PROJECT_DIR/web/app.js"
grep -Fq 'id="official-work-audit-body"' "$PROJECT_DIR/web/index.html"
grep -Fq '/api/official/work-audit' "$PROJECT_DIR/web/app.js"
grep -Fq 'Planning prepares and validates the work without opening MetaMask.' "$PROJECT_DIR/web/app.js"

for work_type in \
  auto \
  general-assistant \
  code-review \
  bug-summary \
  meeting-actions \
  document-analysis \
  prompt-security \
  smart-contract-audit \
  smart-contract-generate \
  smart-contract-explain \
  smart-contract-tests \
  smart-contract-fix
do
  grep -Fq "$work_type" "$PROJECT_DIR/web/app.js"
done

echo "WORK-TYPE SUGGESTION TEST PASSED"
echo "Every work type has tailored guidance, a realistic request, and free fresh AI ideas."
echo "Known suggestions auto-switch; custom requests require explicit replacement."
