#!/usr/bin/env bash
# Starts the built PromptOps server and runs each SDK's end-to-end
# fetch + render + WebSocket hot-reload check against it.
set -euo pipefail

if [ ! -x server/promptops ]; then
  echo "server/promptops binary not found — run 'cd server && go build -o promptops .' first"
  exit 1
fi

echo "==> starting server"
./server/promptops &
SRV=$!
trap 'kill "$SRV" 2>/dev/null || true' EXIT

for _ in $(seq 1 40); do
  if curl -sf http://localhost:8080/health >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -sf http://localhost:8080/health >/dev/null || { echo "FAIL: server did not start"; exit 1; }

echo "==> Node SDK integration"
node sdk/node/examples/hot-reload.mjs

echo "==> Python SDK integration"
python3 sdk/python/examples/hot_reload.py

echo "==> all SDK integration tests passed"
