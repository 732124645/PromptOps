# Introduction

PromptOps is a **Prompt Runtime platform** for AI applications. It decouples
prompts, agents and workflows from code, so AI applications can hot-reload,
version, publish, gray-release and observe their prompts at runtime.

Core idea:

> Make prompts a properly engineered resource — Build → Version → Deploy → Runtime → Observe.

## Why it exists

| Problem | Description |
|---|---|
| Prompts hard-coded | Changing a prompt requires a redeploy |
| No version control | Hard to roll back when something breaks |
| No hot reload | Updating a prompt means restarting the AI service |
| Messy team collaboration | No ownership or permissions |
| No gray release | New prompts go fully live at once — risky |
| No observability | No idea how often prompts run, or how well |

## Capabilities

| Module | Description |
|---|---|
| Prompt Runtime | CRUD, search, categories, multi-environment, version publish & rollback |
| Hot Reload | Real-time prompt-change push over WebSocket |
| SDK | Node / Python / Java — Python & Java zero-dependency, Node depends only on `ws` |
| Playground | Fill variables, call models, view results; version diff |
| Agent | A config-driven prompt + provider + model |
| Workflow | A render → model → transform step-orchestration engine |
| Observability | Audit logs, run logs, token & latency stats, live connected clients |
| RBAC | admin / editor / viewer roles |
| Gray Release | Weighted A/B version traffic split |
| Workspace | Resources grouped by team workspace |

## Architecture

- **Backend**: Go + Gin + GORM + SQLite, with WebSocket hot-reload.
- **Frontend**: Vue 3 + Vite + Naive UI + Pinia.
- **SDKs**: Node (built-in fetch, `ws` for WebSocket), Python (standard library), Java (java.net.http).
- **Deployment**: single binary / Docker.

Next: [Quickstart](./quickstart) · [Deployment](./deployment).
