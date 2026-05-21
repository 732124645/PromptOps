# Changelog

All notable changes to PromptOps are recorded here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project aims to follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
once it reaches `v1.0.0`. Until then, **breaking API or schema changes may
happen on any minor version bump** during the `v0.x` preview series.

## [Unreleased]

_Nothing yet — see [0.1.0] for the current release._

## [0.1.0] - 2026-05-21

First public preview release. The MVP runtime, SDKs, hot-reload, Playground,
Agent / Workflow runtime, observability, audit and RBAC are all in place,
and the project ships with the open-source hygiene needed for external
adoption (LICENSE, security policy, contributor guide, runnable examples).

### Added

- **Backend** — Go + Gin + GORM + SQLite. Prompt CRUD, search, versioning,
  multi-environment support, SDK API (`/api/sdk/prompts/:key`)
- **Hot reload** — WebSocket endpoint `/ws` pushes prompt-change events to
  connected SDK clients in real time; `/ws` is authenticated via the same
  bearer token as the REST API
- **Frontend** — Vue 3 + Vite + Naive UI + Pinia. Pages for Prompts, Agents,
  Workflows, Playground and Observability with English / Chinese i18n
- **SDKs** — Node (depends only on `ws`), Python (zero-dependency stdlib) and
  Java (zero-dependency JDK). All three implement fetch + `{{variable}}`
  template render + WebSocket hot-reload watch
- **Playground** — Model gateway with `mock` / `openai` / `claude` / `ollama`
  / `gemini` providers. Fill variables, call a model, view results, diff
  versions
- **Agent runtime** — Save reusable agent configs (prompt + provider + model)
  and run them with one click
- **Workflow runtime** — Step engine (render → model → transform) that
  orchestrates multi-step flows, chains outputs, and returns a per-step trace
- **Observability** — Audit log for every important change, run log for every
  model call, aggregated run metrics computed in SQL, paginated log endpoints,
  and a live connected-clients registry (`/api/clients`)
- **RBAC** — Users, roles and sessions with PBKDF2-hashed passwords. Three
  roles: `admin`, `editor`, `viewer`. Route access enforced by role
- **Gray release** — Rollout config (A/B traffic split by key + environment).
  The SDK returns one of two versions by weight
- **Workspaces** — Group prompts, agents and workflows by workspace; switch
  between workspaces from the UI
- **Deployment** — `Dockerfile` and `docker-compose.yml` build the frontend
  and backend together and serve them on `:8080`
- **CI** — GitHub Actions builds the server, runs Go tests and the smoke
  test, builds the frontend, unit-tests every SDK, and runs an integration
  job that exercises each SDK's fetch → render → WebSocket hot-reload path
- **Docs site** — VitePress source under `website/`, published to GitHub
  Pages by `.github/workflows/docs.yml`
- **Open-source hygiene** — `LICENSE` (Apache-2.0), `SECURITY.md`,
  `CONTRIBUTING.md`, this changelog, and `.github/` issue / PR templates
- **Examples** — `examples/node-hot-reload` (minimal consumer with live
  WebSocket updates) and `examples/sql-agent` (natural-language → SQL agent
  whose generation rules can be changed at runtime by editing the prompt);
  both linked from the README and the VitePress docs site
- **READMEs** — Tagline, Security Notice and Examples sections in both
  English and Chinese

### Changed

- `sdk/node/package.json` license updated from `MIT` to `Apache-2.0` to
  match the repository LICENSE

[Unreleased]: https://github.com/732124645/PromptOps/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/732124645/PromptOps/releases/tag/v0.1.0
