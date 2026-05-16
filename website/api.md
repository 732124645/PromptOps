# API Reference

All endpoints are prefixed with `/api`. Except for `/api/login` and `/health`,
every request must carry an `Authorization: Bearer <token>` header. The `/ws`
WebSocket also requires a token — passed as a `?token=` query parameter on the
handshake, since browsers cannot set headers on a WebSocket connection.

Permissions: write operations require the `editor` role or above; user
management requires `admin`; `viewer` may only read and run.

## Authentication

| Method | Path | Description |
|---|---|---|
| POST | `/api/login` | Log in (username/password, or static token) |
| POST | `/api/logout` | Revoke the current session |
| GET | `/api/me` | Current user and role |

## Prompt

| Method | Path | Description |
|---|---|---|
| GET | `/api/prompts` | List / search (`q`, `env`, `category`, `tag`, `workspace`) |
| POST | `/api/prompts` | Create |
| GET / PUT / DELETE | `/api/prompts/:id` | Get / update / delete |
| GET | `/api/prompts/:id/versions` | Version history |
| POST | `/api/prompts/publish` | Publish (snapshot current content as a version) |
| POST | `/api/prompts/rollback` | Roll back to a version |
| GET / PUT / DELETE | `/api/prompts/:id/rollout` | Gray-release (A/B) config |
| GET | `/api/sdk/prompts/:key?env=` | SDK runtime fetch by key |

## Playground / Agent / Workflow

| Method | Path | Description |
|---|---|---|
| POST | `/api/playground/run` | Render a prompt and call a model provider |
| GET | `/api/playground/providers` | List available model providers |
| GET / POST / PUT / DELETE | `/api/agents` `/api/agents/:id` | Agent CRUD |
| POST | `/api/agents/:id/run` | Run an agent |
| GET / POST / PUT / DELETE | `/api/workflows` `/api/workflows/:id` | Workflow CRUD |
| POST | `/api/workflows/:id/run` | Run a workflow, returning a per-step trace |

## Governance & Observability

| Method | Path | Description |
|---|---|---|
| GET | `/api/audit` | Audit log (paged: `?limit=` `?offset=`) |
| GET | `/api/runs` | Run log (paged: `?limit=` `?offset=`) |
| GET | `/api/runs/stats` | Aggregated run metrics |
| GET | `/api/clients` | Live hot-reload connections (SDKs & browser) |
| GET / POST / PUT / DELETE | `/api/users` `/api/users/:id` | User management (admin only) |
| GET / POST / DELETE | `/api/workspaces` `/api/workspaces/:id` | Team workspaces |

## Realtime

| Method | Path | Description |
|---|---|---|
| GET | `/ws` | WebSocket hot-reload event stream (auth via `?token=`) |
