# Observability

## What it is

**Observability** is the record of what happened: every model call, every
change, and every connected client. It answers "how often does this run, how
well, and who is using it?"

Everything below lives on the **Observability** page in the UI.

## Run log

Every model invocation — from the [Playground](./playground),
[Agents](./agents) or [Workflows](./workflows) — is logged with:

- timestamp and source
- provider and model
- token counts (in / out, estimated)
- latency in milliseconds
- status (`ok` / `error`)

This is your record of what the AI actually did.

## Stats

The run log is aggregated into summary cards:

- **Total runs**
- **Success / failure** counts
- **Estimated total tokens**
- **Average latency**
- a breakdown **by provider**

## Audit log

Every *mutation* is recorded separately: who changed what, and when. Creating,
updating, deleting, publishing and rolling back a prompt all leave an audit
entry with the action, resource, key and a short summary.

The run log tells you what the AI did; the audit log tells you what people did.

## Connected clients

The **Connected clients** tab lists every live hot-reload connection — SDK
instances and browser sessions currently watching for updates. For each one:

| Column | Description |
|---|---|
| Type | `node-sdk`, `python-sdk`, `java-sdk` or `browser` |
| App | The `appName` the SDK was created with |
| Namespace | The environment it is watching |
| IP address | Where it connected from |
| Connected at | When the connection opened |

This is how you see, at a glance, which applications are currently bound to
PromptOps. See [Hot Reload](./hot-reload) for how clients connect.

## How it maps to the API

| Data | Endpoint |
|---|---|
| Run log | `GET /api/runs` |
| Run stats | `GET /api/runs/stats` |
| Audit log | `GET /api/audit` |
| Connected clients | `GET /api/clients` |

## Next steps

- [Access Control & Workspaces](./access-control) — who is allowed to do what.
