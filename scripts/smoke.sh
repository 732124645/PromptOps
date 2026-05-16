#!/usr/bin/env bash
# Starts the built PromptOps server and exercises the core API end-to-end.
set -euo pipefail

TOKEN="${PROMPTOPS_TOKEN:-promptops-dev-token}"
BASE="http://localhost:8080"

if [ ! -x server/promptops ]; then
  echo "server/promptops binary not found — run 'cd server && go build -o promptops .' first"
  exit 1
fi

echo "==> starting server"
./server/promptops &
SRV=$!
trap 'kill "$SRV" 2>/dev/null || true' EXIT

for _ in $(seq 1 40); do
  if curl -sf "$BASE/health" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -sf "$BASE/health" >/dev/null || { echo "FAIL: server did not start"; exit 1; }
echo "    health OK"

echo "==> login"
curl -sf -X POST "$BASE/api/login" -H 'Content-Type: application/json' \
  -d "{\"token\":\"$TOKEN\"}" >/dev/null
echo "    login OK"

echo "==> create prompt"
curl -sf -X POST "$BASE/api/prompts" \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"key":"smoke.test","name":"Smoke","content":"hello {{name}}","env":"prod","model":"gpt-4o"}' >/dev/null
echo "    create OK"

echo "==> SDK fetch"
OUT="$(curl -sf "$BASE/api/sdk/prompts/smoke.test?env=prod" -H "Authorization: Bearer $TOKEN")"
echo "    response: $OUT"
echo "$OUT" | grep -q 'hello {{name}}' || { echo "FAIL: SDK content mismatch"; exit 1; }

echo "==> unauthorized request is rejected"
CODE="$(curl -s -o /dev/null -w '%{http_code}' "$BASE/api/prompts" -H 'Authorization: Bearer wrong')"
[ "$CODE" = "401" ] || { echo "FAIL: expected 401, got $CODE"; exit 1; }
echo "    auth OK"

echo "==> playground run (mock provider)"
PG="$(curl -sf -X POST "$BASE/api/playground/run" \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"content":"Hi {{who}}","variables":{"who":"world"},"provider":"mock"}')"
echo "    response: $PG"
echo "$PG" | grep -q 'Hi world' || { echo "FAIL: playground did not render variables"; exit 1; }
echo "$PG" | grep -q 'mock completion' || { echo "FAIL: playground output mismatch"; exit 1; }

echo "==> workflow run (render -> model -> transform)"
WF="$(curl -sf -X POST "$BASE/api/workflows" \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"key":"smoke.flow","name":"Smoke Flow","steps":[{"name":"build","type":"render","template":"hi {{who}}"},{"name":"call","type":"model","provider":"mock"},{"name":"shout","type":"transform","op":"upper"}]}')"
WF_ID="$(echo "$WF" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)"
[ -n "$WF_ID" ] || { echo "FAIL: workflow not created"; exit 1; }
RUN="$(curl -sf -X POST "$BASE/api/workflows/$WF_ID/run" \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"variables":{"who":"world"}}')"
echo "    response: $RUN"
echo "$RUN" | grep -q 'HI WORLD' || { echo "FAIL: workflow chain output mismatch"; exit 1; }

echo "==> observability (audit + run stats)"
curl -sf "$BASE/api/audit" -H "Authorization: Bearer $TOKEN" | grep -q '"action"' \
  || { echo "FAIL: no audit entries recorded"; exit 1; }
STATS="$(curl -sf "$BASE/api/runs/stats" -H "Authorization: Bearer $TOKEN")"
echo "    stats: $STATS"
echo "$STATS" | grep -q '"total"' || { echo "FAIL: run stats missing"; exit 1; }

echo "==> smoke test passed"
