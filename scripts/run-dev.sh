#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

mkdir -p data
export SIGNALFORGE_CONFIG="${SIGNALFORGE_CONFIG:-$ROOT/configs/config.yaml}"
export SIGNALFORGE_SERVER_ADDR="${SIGNALFORGE_SERVER_ADDR:-:8080}"

echo "starting signalforge server on ${SIGNALFORGE_SERVER_ADDR}"
go run ./cmd/server &
SERVER_PID=$!

cleanup() {
  if kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT INT TERM

for _ in $(seq 1 60); do
  if curl -fsS "http://127.0.0.1:${SIGNALFORGE_SERVER_ADDR#:}/healthz" >/dev/null 2>&1; then
    echo "server is ready"
    wait "$SERVER_PID"
    exit 0
  fi
  sleep 0.25
done

echo "server did not become ready" >&2
exit 1

