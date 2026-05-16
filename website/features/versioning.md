# Versioning & Publishing

## What it is

A prompt you edit and save is **mutable** — it changes in place. **Publishing**
takes the current content and freezes it into an **immutable version
snapshot**. Versions are what you roll back to when a change goes wrong.

Think of it like Git: editing is your working copy, publishing is a commit.

## Why it matters

- **Safety net** — if a new prompt performs badly, roll back in one click.
- **Auditability** — every published version is kept, with its timestamp.
- **Gray release** — [Gray Release](./rollout) splits traffic between two
  *published* versions, so you must publish before you can A/B test.

## Publishing a version

1. Open a prompt in the editor.
2. Set the `version` label (e.g. bump `v1` → `v2`).
3. Click **Publish version**. The current content is snapshotted under that
   label.

A snapshot records the `version`, `env`, content, and the publish time. It is
never modified afterwards.

## Version history & rollback

In the prompt editor, click **Version history** to see every published
snapshot. For each one you can:

- **Compare with current** — a line-by-line diff between that version and the
  content you are editing now, so you see exactly what changed.
- **Roll back** — replace the working content with that version's content.

## How it maps to the API

| Action | Endpoint |
|---|---|
| Publish current content | `POST /api/prompts/publish` |
| List version history | `GET /api/prompts/:id/versions` |
| Roll back to a version | `POST /api/prompts/rollback` |

## Next steps

- To serve a specific published version to a percentage of traffic, see
  [Gray Release](./rollout).
- Every publish and rollback is recorded — see [Observability](./observability).
