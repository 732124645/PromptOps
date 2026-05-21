# PromptOps

English | [简体中文](README.zh-CN.md)

Open-source runtime platform for AI prompts, agents, and workflows.

> **Stop hardcoding prompts in your application.**

PromptOps is a **Prompt Runtime platform** for AI applications. It decouples
prompts, agents and workflows from code, so AI applications can hot-reload,
version, publish, gray-release and observe their prompts at runtime.

Core idea: **Build → Version → Deploy → Runtime → Observe** prompts.

## Security Notice

The defaults below are intentionally easy so `docker compose up` works out of
the box for local development. **Do not run a default-configured PromptOps
instance on the public internet.**

Before any production deployment, change at minimum:

- The static token `PROMPTOPS_TOKEN` (default `promptops-dev-token`) — replace
  with a high-entropy secret
- The seeded `admin` / `admin` account password (log in once, then change it
  in the Web UI)
- Terminate TLS at a reverse proxy (Nginx / Caddy) in front of `:8080`
- Restrict CORS and admin-UI access to trusted hosts / networks
- Mount the SQLite file on a persistent, backed-up volume

The full hardening checklist and vulnerability-reporting process live in
[`SECURITY.md`](SECURITY.md).

## Status

The MVP, SDKs / hot-reload, Playground, Agent / Workflow runtime, observability
and audit are all complete.

| Module | Stack | Description |
|---|---|---|
| Backend | Go + Gin + GORM + SQLite | Prompt CRUD / search / versioning / multi-env / SDK API |
| Hot Reload | WebSocket (`/ws`) | Real-time prompt-change push to connected clients |
| Frontend | Vue 3 + Vite + Naive UI + Pinia | Prompt / Agent / Workflow / Playground / Observability pages |
| SDK | Node (depends on `ws`) / Python / Java (zero-dependency) | Fetch prompts by key, render templates, WebSocket hot-reload |
| Playground | Model gateway (mock / OpenAI / Claude / Ollama / Gemini) | Fill variables, call models, view results; version diff |
| Agent | Config-driven agent (prompt + provider + model) | Save reusable configs and run with one click |
| Workflow | Step engine (render → model → transform) | Orchestrate multi-step flows, chain output, view step trace |
| Observability | Audit log + run log + token stats + connected clients | Records every change and model call, aggregates run metrics |
| RBAC | Users / roles / sessions (PBKDF2 passwords) | admin / editor / viewer roles, route access by role |
| Gray Release | Rollout (A/B traffic split by key + env) | The SDK returns one of two versions by weight |
| Workspace | Workspace | Prompts / Agents / Workflows grouped by workspace, switchable in the UI |

## Project structure

```txt
PromptOps/
├── server/      # Go backend (Gin + GORM + SQLite)
├── web/         # Vue 3 frontend (Vite + Naive UI)
├── sdk/         # Runtime SDKs (node / python / java)
├── docs/        # Design docs
├── scripts/     # Smoke tests and other scripts
├── Dockerfile
└── docker-compose.yml
```

## Local development

### Backend

```bash
cd server
go mod tidy
go run .            # listens on :8080 by default; SQLite file at server/data/promptops.db
```

Environment variables: `PROMPTOPS_ADDR` (default `:8080`), `PROMPTOPS_DB`,
`PROMPTOPS_TOKEN` (default `promptops-dev-token`).

### Frontend

```bash
cd web
npm install
npm run dev         # http://localhost:5173, proxied to the backend on :8080
```

Default login: `admin` / `admin` (created on first start). The static token
`promptops-dev-token` also works as an admin credential for SDKs and scripts.

## Docker deployment

```bash
docker compose up --build   # builds frontend + backend, open http://localhost:8080
```

## API overview

| Method | Path | Description |
|---|---|---|
| POST | `/api/login` | Log in (username / password, or static token) |
| GET | `/api/me` | Current user and role |
| GET/POST/PUT/DELETE | `/api/users` `/api/users/:id` | User management (admin only) |
| GET/POST/DELETE | `/api/workspaces` `/api/workspaces/:id` | Workspace management |
| GET | `/api/prompts` | List / search (`q`, `env`, `category`, `tag`) |
| POST | `/api/prompts` | Create |
| GET/PUT/DELETE | `/api/prompts/:id` | Get / update / delete |
| GET | `/api/prompts/:id/versions` | Version history |
| POST | `/api/prompts/publish` | Publish (snapshot current content as a version) |
| POST | `/api/prompts/rollback` | Roll back to a version |
| GET/PUT/DELETE | `/api/prompts/:id/rollout` | Gray release (A/B) config |
| GET | `/api/sdk/prompts/:key?env=` | SDK runtime fetch by key |
| POST | `/api/playground/run` | Render a prompt and call a model provider |
| GET | `/api/playground/providers` | List available model providers |
| GET/POST/PUT/DELETE | `/api/agents` `/api/agents/:id` | Agent CRUD |
| POST | `/api/agents/:id/run` | Run an agent |
| GET/POST/PUT/DELETE | `/api/workflows` `/api/workflows/:id` | Workflow CRUD |
| POST | `/api/workflows/:id/run` | Run a workflow, returns the step trace |
| GET | `/api/audit` | Audit log (change records) |
| GET | `/api/runs` | Run log (model-call records) |
| GET | `/api/runs/stats` | Aggregated run metrics (count / tokens / latency) |
| GET | `/api/clients` | Live hot-reload connections (SDKs & browser) |
| GET | `/ws` | WebSocket hot-reload event stream |

All `/api/*` routes except `/api/login` and `/health` require
`Authorization: Bearer <token>`; the `/ws` WebSocket authenticates via a
`?token=` query parameter. Write operations require the `editor` role or above,
user management requires `admin`; `viewer` can only read and run.

## SDK

PromptOps provides lightweight runtime SDKs for fetching prompts by key,
rendering `{{variable}}` templates, and hot-reload over WebSocket (`watch()` —
the local cache refreshes automatically when a prompt changes on the server).

- `sdk/node` — Node.js (Node 18+, depends only on `ws`; uses the built-in `fetch`)
- `sdk/python` — Python (zero-dependency, standard library only, with a minimal WebSocket client)
- `sdk/java` — Java (zero-dependency, pure JDK, uses `java.net.http`)

Node example:

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({ server: 'http://localhost:8080', namespace: 'prod' })

const text = await client.render('sql.generator', { question: 'list all users' })

client.on('update', (e) => console.log('prompt hot-reloaded:', e.key))
client.watch()
```

## Examples

End-to-end runnable demos live under [`examples/`](examples/):

- [`examples/node-hot-reload`](examples/node-hot-reload) — minimal Node
  consumer that fetches a prompt every few seconds and reflects edits made
  in the Web UI without restarting
- [`examples/sql-agent`](examples/sql-agent) — natural-language → SQL agent
  showing how to change generation rules at runtime by editing the prompt
  (uses the built-in `mock` provider, so no API key is required)

See [`examples/README.md`](examples/README.md) for the full list and
contribution ideas.

## CI

`.github/workflows/ci.yml` verifies on GitHub runners: the backend
(`go build`, `go test`, smoke test), the frontend (`npm run build`), unit tests
for all three SDKs, and an integration job — starting a real server and running
each SDK's full "fetch → render → WebSocket hot-reload" path.

## Documentation

Online documentation site (VitePress, hosted on GitHub Pages):
<https://732124645.github.io/PromptOps/> — the source lives in `website/` and
is built and published by `.github/workflows/docs.yml`.

The full design document is at [docs/PromptOps.md](docs/PromptOps.md), covering
positioning, core concepts, MVP scope, architecture, database and API design,
SDKs, hot-reload, deployment and the roadmap.
