# Workflows

## What it is

A **Workflow** is a multi-step pipeline. Where an [Agent](./agents) is a single
prompt-and-model call, a Workflow chains several steps in order — each step's
output becomes the next step's input.

## Step types

| Type | What it does |
|---|---|
| `render` | Render a template, filling `{{variables}}` (and `{{input}}`) |
| `model` | Call a model provider with the current text |
| `transform` | Apply a text transform: `upper`, `lower` or `trim` |

The special variable `{{input}}` in a `render` step holds the **output of the
previous step**, which is how steps are chained together.

## Why it matters

Real tasks are rarely one call. For example: *render a summarisation prompt →
call a model → trim the result*. A Workflow makes that sequence a single,
saved, runnable unit instead of glue code scattered across your app.

## Using it

1. Open the **Workflows** page and click **New Workflow**.
2. Give it a key and name.
3. Click **Add step** for each step. For each one, pick a type and configure it:
   - `render` — write the template
   - `model` — pick a provider and model
   - `transform` — pick the operation
4. Reorder steps with the ↑ / ↓ buttons.
5. Click **Save**, then fill in any variables and click **Run Workflow**.

## Reading the result

After a run, the result panel shows:

- **Step trace** — the output of every step, in order, so you can see exactly
  where the text changed.
- **Final output** — the result of the last step.

## Behind the scenes

Running a workflow calls `POST /api/workflows/:id/run`, which returns the full
step trace. Each run is recorded for [Observability](./observability).

## Next steps

- [Agents](./agents) — the single-step counterpart.
- [Playground](./playground) — try individual prompts before wiring them in.
