# Node Hot-Reload Example

A 60-line Node.js script that subscribes to a single PromptOps prompt and
prints its rendered content every few seconds. When you edit and publish
the prompt in the Web UI, the next render here reflects the new content
**without restarting the process**.

This is the smallest possible demonstration of PromptOps' core idea: change
prompts at runtime, observe the change in a live application.

## Prerequisites

- PromptOps running locally (`docker compose up` from the repo root, or the
  manual backend / frontend setup in the root [`README.md`](../../README.md))
- Node.js 18+

## Run

```bash
cd examples/node-hot-reload
npm install
cp .env.example .env   # optional — defaults work for a local PromptOps
npm start
```

You should see the script connect to the WebSocket and then complain that
the prompt does not exist yet. That's expected — create it in the next step.

## Create the prompt

1. Open the PromptOps Web UI (default <http://localhost:8080>, log in as
   `admin` / `admin`)
2. Go to **Prompts → New**
3. Fill in:
   - **Key**: `demo.greeting`
   - **Env**: `prod`
   - **Content**:
     ```
     Hello {{name}}! The time is {{now}}.
     ```
4. Save, then click **Publish**

The terminal should immediately start printing the rendered greeting every
few seconds.

## See hot-reload in action

1. Back in the Web UI, edit the same prompt to:
   ```
   👋 Hello {{name}} from PromptOps! Now: {{now}}.
   ```
2. Click **Publish**
3. Watch the terminal — within ~1 second you'll see a `[update]` line, and
   the next render will show the new content

No restart. No redeploy. The SDK's WebSocket subscription invalidated the
local cache and pulled the new version.

## What this example uses

- `@promptops/client` (resolved via `"file:../../sdk/node"` in
  [`package.json`](package.json))
- The runtime fetch endpoint `GET /api/sdk/prompts/:key?env=`
- The WebSocket endpoint `GET /ws?token=...&client=node-sdk&namespace=...`

The full SDK reference is in
[`sdk/node/index.d.ts`](../../sdk/node/index.d.ts).

## Configuration

All settings are environment variables — see [`.env.example`](.env.example).
The defaults assume a vanilla local PromptOps install.
