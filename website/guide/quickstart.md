# Quickstart

## Requirements

- Go 1.24+
- Node.js 18+

## Start the backend

```bash
cd server
go mod tidy
go run .
```

The backend listens on `:8080` by default; the SQLite database lives at
`server/data/promptops.db`.

Environment variables:

| Variable | Default | Description |
|---|---|---|
| `PROMPTOPS_ADDR` | `:8080` | Listen address |
| `PROMPTOPS_DB` | `data/promptops.db` | SQLite file path |
| `PROMPTOPS_TOKEN` | `promptops-dev-token` | Static admin token |

## Start the frontend

```bash
cd web
npm install
npm run dev
```

Open `http://localhost:5173`. The default admin account is `admin` / `admin`.

## Create your first prompt

1. After logging in, open the **Prompts** page and click "New Prompt".
2. Fill in a `key` (e.g. `sql.generator`), an environment and the content;
   the content can use `{{variable}}` placeholders.
3. Save, then click "Publish" to create an immutable version snapshot.

## Try it in the Playground

Open the **Playground**, pick the prompt you created, fill in variable values,
choose a model provider (`mock` needs no API key and runs offline), and click
"Run" to see the rendered prompt and the model output.

## Use it from an SDK

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
})

const text = await client.render('sql.generator', { question: 'list all users' })

client.on('update', (e) => console.log('prompt hot-reloaded:', e.key))
client.watch()
```

See the [SDK docs](/sdk/node) and the [API reference](/api) for more.
