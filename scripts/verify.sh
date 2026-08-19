#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

gofmt -w .
go mod tidy
go test ./...
go vet ./...
go build ./...
./scripts/stats.sh

