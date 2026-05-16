# Playground

## What it is

The **Playground** is an interactive console for trying a prompt against a real
(or mock) model — without writing any code. Fill in the variables, pick a
provider, run, and see both the rendered prompt and the model's output.

## Why it matters

It closes the loop between *editing* a prompt and *knowing if it works*. You
can iterate on wording and immediately see the effect, before publishing.

## Model providers

| Provider | Needs |
|---|---|
| `mock` | Nothing — runs offline, returns a deterministic echo. Great for testing. |
| `openai` | An API key |
| `claude` | An API key |
| `gemini` | An API key |
| `ollama` | A Base URL (e.g. `http://localhost:11434`) for your local Ollama |

API keys are entered per-run and used only for that call — they are not stored.

## Using it

1. Open the **Playground**.
2. Optionally pick an existing prompt to load its content, or type content
   directly.
3. The Playground detects every `{{variable}}` and shows an input for each —
   fill them in.
4. Choose a **model provider** and model. Start with `mock` if you just want to
   see the rendered prompt with no API key.
5. Click **Run**.

The result panel shows:

- **Rendered prompt** — your content with all `{{variables}}` substituted.
- **Model output** — what the provider returned, with the provider and model
  used.

## Behind the scenes

The Playground calls `POST /api/playground/run`. Every run is also recorded for
[Observability](./observability) — you can see it later in the run log.

## Next steps

- To save a prompt + provider + model as a reusable, runnable unit, see
  [Agents](./agents).
- To chain multiple steps together, see [Workflows](./workflows).
