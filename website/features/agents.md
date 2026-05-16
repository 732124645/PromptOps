# Agents

## What it is

An **Agent** is a saved, reusable bundle of a prompt plus the model setup to
run it: `prompt + provider + model`. Where the [Playground](./playground) is a
scratchpad, an Agent is a named configuration you keep and run again.

## Why it matters

In practice you don't just have prompt text — you have "this prompt, run on
this provider, with this model". An Agent captures that whole unit once, so
anyone on the team can run it consistently without re-entering the setup.

## Anatomy of an agent

| Field | Description |
|---|---|
| `key` | Stable identifier, e.g. `support.agent` |
| `name` | Human-readable label |
| `description` | What this agent is for |
| `prompt` | The prompt text, with `{{variable}}` placeholders |
| `provider` | `mock` / `openai` / `claude` / `ollama` / `gemini` |
| `model` | Model name; leave empty for the provider default |

## Using it

1. Open the **Agents** page and click **New Agent**.
2. Fill in the key, choose a provider and model, and write the prompt.
3. Click **Save**.
4. On the agent's page, fill in any detected `{{variables}}`, supply an API key
   if the provider needs one, and click **Run Agent**.

The result shows the rendered prompt and the model output, just like the
Playground — but the configuration is saved for next time.

## Behind the scenes

Running an agent calls `POST /api/agents/:id/run`. Each run is recorded for
[Observability](./observability).

## Agent vs. Workflow

- An **Agent** is a *single* prompt-and-model step.
- A [**Workflow**](./workflows) chains *multiple* steps, feeding each step's
  output into the next.

## Next steps

- [Workflows](./workflows) — orchestrate multiple steps.
- [Observability](./observability) — review agent runs.
