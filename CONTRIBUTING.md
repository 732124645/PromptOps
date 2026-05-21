# Contributing to PromptOps

Thanks for taking the time to contribute. PromptOps is an early-stage project,
so the bar for "useful contribution" is low — bug reports, doc tweaks, example
projects and small features are all welcome.

## Quick links

- Bugs and feature requests: open a GitHub issue (use the templates)
- Security issues: see [`SECURITY.md`](SECURITY.md), do **not** file public issues
- Larger changes: open an issue first to discuss the design before coding

## Local development

The project has three independent components plus three SDKs. Each can be
developed in isolation. See the per-component instructions in [`README.md`](README.md).

### Backend (`server/`)

```bash
cd server
go mod tidy
go run .            # listens on :8080; SQLite file at server/data/promptops.db
```

Run tests:

```bash
cd server
go test ./...
```

### Frontend (`web/`)

```bash
cd web
npm install
npm run dev         # http://localhost:5173, proxied to :8080
```

### SDKs (`sdk/{node,python,java}`)

Each SDK has its own tests:

```bash
cd sdk/node && node --test
cd sdk/python && python -m unittest discover
cd sdk/java && ./gradlew test       # or mvn test, depending on your setup
```

## End-to-end verification

Before opening a PR that touches the backend, an SDK, or the WebSocket
contract, run the integration script. It boots a real server and exercises
each SDK's `fetch → render → WebSocket hot-reload` path:

```bash
bash scripts/sdk-integration.sh
```

For backend-only API changes, the smoke test is faster:

```bash
bash scripts/smoke.sh
```

Both scripts are also run in CI on every PR (`.github/workflows/ci.yml`).

## Commit messages

Match the existing style — short imperative subject, optionally with a
component scope. Examples from recent history:

```
Authenticate /ws, paginate logs, aggregate stats in SQL
feat(i18n): add internationalization support with English and Chinese translations
Add live client tracking, cyberpunk UI redesign, and feature docs
```

Keep the subject under ~72 characters. Use the body for the *why*, not the
*what* — the diff already shows what changed.

## Pull request checklist

- [ ] Touched code has tests (unit, integration, or both as appropriate)
- [ ] `go test ./...` passes for backend changes
- [ ] `npm run build` passes for frontend changes
- [ ] `bash scripts/sdk-integration.sh` passes for SDK or `/ws` changes
- [ ] User-facing changes are reflected in `README.md` / `README.zh-CN.md`
- [ ] Notable changes are added to `CHANGELOG.md` under `[Unreleased]`
- [ ] No new defaults that weaken security (see [`SECURITY.md`](SECURITY.md))

## Code style

- **Go**: `gofmt`-ed (CI enforces this). Prefer small focused packages under
  `server/internal/`. Handler functions go on the `*Handler` receiver.
- **Vue / TypeScript**: project ESLint config. Components in `web/src/`.
  Keep stores small and centred on Pinia.
- **SDKs**: keep dependencies minimal. Node depends only on `ws`; Python and
  Java are zero-dependency. New SDK methods should have parity across all
  three languages where it makes sense.

## What to work on if you don't have your own itch

Look for issues labelled `good first issue` or `help wanted`. Things that are
always welcome:

- New examples under `examples/` (Python, Java, Spring Boot, customer-service
  agent, workflow demos, etc.)
- Improvements to the [VitePress docs site](website/)
- Bug reports — even ones you can't fix — with clear reproduction steps
- Translations into languages other than English and Chinese

## Code of conduct

Be respectful. We don't have a formal CoC yet; if a situation arises that
needs one, we'll adopt the
[Contributor Covenant](https://www.contributor-covenant.org/).
