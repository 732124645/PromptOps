# Security Policy

PromptOps is an early-stage open-source project. We take security seriously and
welcome reports of any vulnerability you find.

## Supported Versions

Only the latest `v0.1.x` preview is actively supported. Older preview builds do
not receive security fixes — upgrade before reporting issues against them.

| Version | Status |
|---|---|
| `v0.1.x` (current preview) | Supported |
| `< v0.1.0` | Not supported |

## Reporting a Vulnerability

**Do not file a public GitHub issue for security problems.** Public reports give
attackers a head start on unpatched users.

Please email **leizhuo666@gmail.com** with:

- A clear description of the issue and its impact
- Step-by-step reproduction (request payloads, configs, screenshots if useful)
- The PromptOps version (backend commit, frontend build, SDK version) and your
  environment (OS, Docker / Go / Node version)
- Any suggested fix or mitigation, if you have one

We aim to acknowledge reports within **3 business days** and to ship a fix or
mitigation within **30 days** for confirmed vulnerabilities. Once a fix is
released we are happy to credit reporters (or keep the report anonymous on
request).

## Known Default-Credential Risks

PromptOps ships with developer-friendly defaults so that `docker compose up`
works out of the box. **These defaults must be changed before any production
deployment.**

| Item | Default value | Where it lives | What to change it to |
|---|---|---|---|
| Static SDK / admin token | `promptops-dev-token` | env var `PROMPTOPS_TOKEN`, fallback in [`server/internal/handlers/middleware.go`](server/internal/handlers/middleware.go) | A long random secret, set via `PROMPTOPS_TOKEN` |
| Admin account | `admin` / `admin` | Seeded on first start when no users exist | Log in once, then change the password in the Web UI |
| Database location | `server/data/promptops.db` (SQLite) | env var `PROMPTOPS_DB` | A path on a persistent, backed-up volume |
| Listen address | `:8080` | env var `PROMPTOPS_ADDR` | Bind behind a reverse proxy (Nginx / Caddy) with HTTPS |

## Production Hardening Checklist

Before exposing PromptOps to the public internet, at minimum:

- [ ] Replace `PROMPTOPS_TOKEN` with a high-entropy secret (32+ random bytes)
- [ ] Change the seeded `admin` password
- [ ] Put the server behind HTTPS via a reverse proxy
- [ ] Restrict CORS / origins to your own application hosts
- [ ] Mount the SQLite file (or future Postgres database) on a persistent,
      regularly backed-up volume
- [ ] Restrict access to the admin Web UI — IP allowlist, VPN, or SSO upstream
- [ ] Rotate API tokens periodically and revoke unused ones

## Scope

In scope:

- The PromptOps server (`server/`), Web UI (`web/`), SDKs (`sdk/node`,
  `sdk/python`, `sdk/java`), and the published Docker images
- The default deployment as described in `README.md` and `docker-compose.yml`

Out of scope:

- Vulnerabilities in third-party LLM providers, model gateways, or any service
  PromptOps merely forwards requests to
- Self-inflicted misconfigurations (e.g. running with default credentials on
  the public internet, disabling auth in custom forks)
- Issues that require a pre-authenticated `admin` role and only affect that
  role's own data
