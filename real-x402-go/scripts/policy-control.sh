#!/usr/bin/env bash

policy_project_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
policy_token_file="${POLICY_CONTROL_TOKEN_FILE:-$policy_project_dir/../go-core/data/policy-control-token}"

if [[ -z "${POLICY_CONTROL_TOKEN:-}" ]]; then
  if [[ ! -f "$policy_token_file" ]]; then
    echo "Policy control token file is missing: $policy_token_file" >&2
    echo "Start the Go core once so it can create the protected local token." >&2
    exit 1
  fi
  POLICY_CONTROL_TOKEN="$(<"$policy_token_file")"
  export POLICY_CONTROL_TOKEN
fi

POLICY_CONTROL_HEADER=("X-Policy-Control-Token: $POLICY_CONTROL_TOKEN")
