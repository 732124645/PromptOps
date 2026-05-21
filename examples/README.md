# PromptOps Examples

> 中文索引见 [`zh.md`](zh.md).

End-to-end runnable demos that show how PromptOps is used in real scenarios.
Each example assumes you already have PromptOps running locally — see the
Quick Start or Docker section in the repo root [`README.md`](../README.md).

## Available examples

| Directory | Language | What it shows |
|---|---|---|
| [`node-hot-reload`](node-hot-reload) | Node.js | Minimal consumer that fetches one prompt and subscribes to hot-reload via WebSocket. Edit and publish the prompt in the Web UI; the terminal reflects the new content immediately without a restart. |
| [`sql-agent`](sql-agent) | Node.js | Natural-language → SQL agent. Pulls the `sql.generator` template from PromptOps, calls the Playground's `mock` provider to produce SQL, and updates its generation rules at runtime by editing the prompt. |

## Examples we'd love to receive

Pick any of these (or propose your own) — see [`../CONTRIBUTING.md`](../CONTRIBUTING.md):

- `python-hot-reload` — Python equivalent of the minimal hot-reload demo
- `java-springboot-hot-reload` — Spring Boot integration example
- `customer-service-agent` — Versioned reply templates with gray-release rollout
- `workflow-demo` — Multi-step Workflow (`render → model → transform`)
- `code-review-agent` — Prompt-driven code review agent

## Conventions

- Every example has its own `README.md` with a 30-second walkthrough
- Examples consume the local SDK via `"file:../../sdk/node"` (or equivalent
  for other languages) so they work without publishing to a package registry
- API keys are never required — examples default to the built-in `mock`
  provider; instructions for swapping in a real provider are in each README
- All examples respect the same env-var conventions:
  `PROMPTOPS_SERVER`, `PROMPTOPS_TOKEN`, `PROMPTOPS_NAMESPACE`
