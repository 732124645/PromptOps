# Prompts

## What it is

A **Prompt** is the core resource in PromptOps — a named, versioned piece of
text that your AI application fetches at runtime instead of hard-coding it.

Instead of writing prompt text inside your code and redeploying every time it
changes, you store it here and fetch it by `key`. Editors can change the
wording; your app picks it up without a release.

## Anatomy of a prompt

| Field | Description |
|---|---|
| `key` | Stable identifier the SDK fetches by, e.g. `sql.generator` |
| `name` | Human-readable label |
| `content` | The prompt text, with `{{variable}}` placeholders |
| `env` | Environment — `dev`, `test` or `prod` |
| `version` | The current working version label, e.g. `v1` |
| `category` | Free-text grouping, e.g. `code` |
| `tags` | Comma-separated tags, e.g. `review, sql` |
| `model` | Suggested model for this prompt, e.g. `gpt-4o` |

The same `key` can exist once per environment — `sql.generator` in `dev` and
`sql.generator` in `prod` are two separate prompts.

## Template variables

Prompt content can contain `{{variable}}` placeholders. They are filled in at
render time — by the SDK's `render()`, by the Playground, or by an Agent:

```
You are a {{language}} expert.
Please review the code below:
{{code}}
```

The editor automatically detects every `{{variable}}` in the content and shows
them, so you always know what inputs a prompt expects.

## Using it

1. Open the **Prompts** page and click **New Prompt**.
2. Fill in `key`, environment, content and any optional metadata.
3. Click **Save**.
4. Use the search box (matches key / name / content) and the environment
   filter to find prompts later.

## Fetching a prompt from your app

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({ server: 'http://localhost:8080', namespace: 'prod' })

// renders {{variable}} placeholders with the values you pass
const text = await client.render('sql.generator', { question: 'list all users' })
```

## Next steps

- Saving a prompt is editable and mutable. To freeze it, see
  [Versioning & Publishing](./versioning).
- To push changes to running apps instantly, see [Hot Reload](./hot-reload).
