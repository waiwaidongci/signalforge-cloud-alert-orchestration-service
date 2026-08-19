#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

FILES="$(find . -type f -name '*.go' ! -name '*_test.go' | sort)"
FILE_COUNT="$(printf '%s\n' "$FILES" | sed '/^$/d' | wc -l | tr -d ' ')"
LINES="$(printf '%s\n' "$FILES" | xargs wc -l | tail -n 1 | awk '{print $1}')"
echo "non-test Go files: ${FILE_COUNT}"
echo "non-test Go lines: ${LINES}"

