# SQL Agent Example

A small Node.js agent that takes a natural-language question and returns a
SQL query. The SQL generation rules live entirely in a PromptOps-managed
prompt — change the rules in the Web UI, click **Publish**, and the next
run of the agent follows the new rules. No redeploy.

Typical use cases this pattern fits:

- BI / analytics assistants
- ERP report generators
- CRM "ask in English" data exploration
- Internal data-access tools where the SQL style needs to evolve

## Prerequisites

- PromptOps running locally (`docker compose up` from the repo root, or the
  manual setup in the root [`README.md`](../../README.md))
- Node.js 18+

## Step 1 — Create the prompt

The agent reads a prompt with key `sql.generator`. Create it in the Web UI
using the starter template in
[`prompts/sql-generator.md`](prompts/sql-generator.md):

1. Open the PromptOps Web UI (default <http://localhost:8080>, log in as
   `admin` / `admin`)
2. Go to **Prompts → New**
3. Set **Key** = `sql.generator`, **Env** = `prod`
4. Paste the **Content** block from `prompts/sql-generator.md`
5. Save and click **Publish**

## Step 2 — Run the agent

```bash
cd examples/sql-agent
npm install
cp .env.example .env   # optional — defaults work for a local PromptOps
npm start -- "Top 10 customers by revenue this month"
```

Expected output (with the default `mock` provider):

```
--- Rendered prompt (what the model saw) ---
You are a senior data engineer. Generate a single mysql SQL query that
answers the user's question. Follow these rules:
... (the full rendered template with your question filled in) ...

--- Model output (mock / mock-1) ---
[mock completion]
model: mock-1
prompt: ... chars, ... line(s)
first line: ...

This is a deterministic mock response. Select a real provider and supply
an API key to call an actual LLM.
```

The mock provider deliberately echoes the prompt back. That's enough to
prove the round-trip — the value the example shows is **prompt management
at runtime**, not the LLM itself.

## Step 3 — Change the SQL style at runtime

1. While the example is idle (or running), edit `sql.generator` in the Web UI.
   For example, replace `Use snake_case for all aliases.` with
   `Use CamelCase for all aliases.`
2. Click **Publish**
3. Run the example again — the rendered prompt section will reflect the new
   rule immediately

The SDK's WebSocket subscription invalidated the local cache, so the next
`client.getPrompt()` call fetches the new version. No process restart.

## Calling a real LLM

To swap the mock provider for a real one, set these in your `.env`:

```bash
PROVIDER=openai            # or claude / ollama / gemini
MODEL=gpt-4o-mini          # any model supported by your provider
PROVIDER_API_KEY=sk-...
# PROVIDER_BASE_URL=       # only if you need a custom endpoint
```

The Playground forwards these directly to the model gateway in
[`server/internal/providers/`](../../server/internal/providers/). The same
prompt now produces real SQL instead of the mock echo.

## What this example uses

- `@promptops/client` (resolved via `"file:../../sdk/node"` in
  [`package.json`](package.json))
- `GET /api/sdk/prompts/:key?env=` to fetch the prompt
- `POST /api/playground/run` to render + invoke the chosen model provider
- `GET /ws` for hot-reload notifications

See the API table in the root [`README.md`](../../README.md) for the full
endpoint reference.
