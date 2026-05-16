# Gray Release

## What it is

**Gray Release** (also called a rollout, or A/B release) serves *two* published
versions of a prompt side by side, splitting traffic between them by weight.
Instead of switching everyone to a new prompt at once, you send, say, 20% of
requests to the new version and watch.

## Why it matters

A new prompt can behave worse in ways you won't see until real traffic hits it.
Gray Release de-risks the change: roll it out to a slice, compare, then widen
or roll back.

## How it works

A rollout is configured **per prompt** (by `key` + `env`):

| Setting | Meaning |
|---|---|
| `enabled` | Whether the rollout is active |
| `variant_a` | A **published** version label, e.g. `v1` |
| `variant_b` | Another **published** version label, e.g. `v2` |
| `weight_a` | Percentage of traffic sent to variant A (the rest goes to B) |

When an SDK fetches the prompt by `key` + `env`, the server rolls the dice by
`weight_a` and returns variant A or variant B's published content.

::: tip
Both variant versions must be **published first** — see
[Versioning & Publishing](./versioning). The rollout chooses between published
snapshots, not the mutable working content.
:::

## Using it

1. Publish at least two versions of the prompt (e.g. `v1` and `v2`).
2. Open the prompt editor and find the **Gray Release (A/B)** card.
3. Turn on **Enable rollout**.
4. Set variant A and variant B to two published version labels.
5. Set **Variant A traffic (%)** — e.g. `80` means 80% A, 20% B.
6. Click **Save rollout config**.

## How it maps to the API

| Action | Endpoint |
|---|---|
| Get rollout config | `GET /api/prompts/:id/rollout` |
| Set rollout config | `PUT /api/prompts/:id/rollout` |
| Remove rollout | `DELETE /api/prompts/:id/rollout` |

## Next steps

- [Versioning & Publishing](./versioning) — publish the versions a rollout
  chooses between.
- [Observability](./observability) — watch how each variant performs.
