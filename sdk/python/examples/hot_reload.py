"""End-to-end check: fetch + render + WebSocket hot-reload against a running
PromptOps server. Used by CI's integration job."""

import json
import os
import sys
import threading
import time
import urllib.request

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from promptops import PromptOpsClient

SERVER = os.environ.get("PROMPTOPS_SERVER", "http://localhost:8080")
TOKEN = os.environ.get("PROMPTOPS_TOKEN", "promptops-dev-token")


def api(path, method="GET", body=None):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(
        SERVER + path,
        data=data,
        method=method,
        headers={"Authorization": "Bearer " + TOKEN, "Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        return json.loads(resp.read().decode("utf-8"))


def fail(message):
    print("FAIL:", message)
    sys.exit(1)


created = api(
    "/api/prompts",
    "POST",
    {
        "key": "sdk.py.demo",
        "name": "Py Demo",
        "content": "Hello {{name}}, v1",
        "env": "prod",
        "model": "gpt-4o",
    },
)
prompt_id = created["data"]["id"]

client = PromptOpsClient(SERVER, namespace="prod", token=TOKEN)

rendered = client.render("sdk.py.demo", {"name": "Ada"})
if rendered != "Hello Ada, v1":
    fail("render mismatch: " + rendered)
print("ok - get_prompt + render ->", rendered)

updated_event = threading.Event()
client.watch(on_update=lambda _e: updated_event.set())
time.sleep(0.7)  # allow the WebSocket handshake to complete

api(
    "/api/prompts/" + prompt_id,
    "PUT",
    {"name": "Py Demo", "content": "Hello {{name}}, v2", "version": "v2", "env": "prod"},
)

if not updated_event.wait(timeout=5):
    fail("no hot-reload event within 5s")

after = client.render("sdk.py.demo", {"name": "Ada"})
if after != "Hello Ada, v2":
    fail("cache not refreshed after hot-reload: " + after)
print("ok - hot-reload via WebSocket ->", after)

client.close()
api("/api/prompts/" + prompt_id, "DELETE")
print("Python SDK integration test passed")
