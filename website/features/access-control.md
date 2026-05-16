# Access Control & Workspaces

Two separate things keep a team's use of PromptOps organised: **roles** decide
*what a person may do*, **workspaces** decide *how resources are grouped*.

## Roles (RBAC)

Every user has one of three roles:

| Role | Can do |
|---|---|
| `viewer` | Read everything; run prompts, agents and workflows |
| `editor` | Everything a viewer can, plus create / edit / delete / publish |
| `admin` | Everything an editor can, plus user management and workspace deletion |

The rule of thumb: **viewer** reads and runs, **editor** also writes,
**admin** also administers.

### Signing in

You authenticate in one of two ways:

- **Username + password** — a normal user account. The server issues a session
  token; the SDK and UI send it as `Authorization: Bearer <token>`.
- **Static token** — the value of `PROMPTOPS_TOKEN` (default
  `promptops-dev-token`) always works as an `admin` token. It is convenient for
  SDKs and scripts; change it for any real deployment.

The default admin account is `admin` / `admin`, created on first start.

### Managing users

An `admin` opens the **Users** page to create accounts, change roles, or remove
users. Roles take effect immediately on the user's next request.

## Workspaces

A **workspace** groups prompts, agents and workflows for a team or project.
Each of those resources belongs to exactly one workspace.

- The UI has a workspace switcher in the header — pick one and the lists scope
  to it.
- There is always a built-in `default` workspace.
- Any signed-in user can create a workspace; deleting one is `admin`-only.
- Deleting a workspace does **not** delete its resources — they simply stop
  belonging to it.

Workspaces are about *organisation* (team / project separation). The `env`
field on a prompt is about *deployment stage* (`dev` / `test` / `prod`) — the
two are independent.

## How it maps to the API

| Action | Endpoint | Required role |
|---|---|---|
| Log in | `POST /api/login` | — |
| Current user | `GET /api/me` | viewer |
| User management | `/api/users` `…/:id` | admin |
| List / create workspaces | `GET` / `POST /api/workspaces` | viewer / editor |
| Delete a workspace | `DELETE /api/workspaces/:id` | admin |

All `/api/*` routes except `/api/login` require the `Authorization` header;
write routes require `editor`, user management requires `admin`.

## Next steps

- [Observability](./observability) — the audit log records who changed what.
