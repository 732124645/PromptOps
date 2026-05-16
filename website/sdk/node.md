# Node SDK

Node 18+. Uses the built-in `fetch`; WebSocket hot-reload is powered by
[`ws`](https://www.npmjs.com/package/ws), its only dependency. Source lives
in `sdk/node`.

## Usage

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
  token: 'promptops-dev-token',
})

// Fetch a prompt by key (locally cached)
const prompt = await client.getPrompt('sql.generator')

// Render template variables
const text = await client.render('sql.generator', { question: 'list all users' })

// Subscribe to hot-reload: the cache refreshes when the prompt changes
client.on('update', (e) => console.log('hot-reloaded:', e.key))
client.watch()
```

## API

| Method | Description |
|---|---|
| `new PromptOpsClient({ server, namespace?, token?, appName? })` | Create a client (`appName` shows in the server's connected-clients panel) |
| `getPrompt(key, { refresh? })` | Fetch a prompt; cached by default |
| `render(key, vars)` | Fetch and render `{{variables}}` |
| `on(event, cb)` | Listen for `update` / `connect` / `disconnect` / `error` |
| `watch()` | Open the WebSocket and auto-refresh the cache |
| `close()` | Close the connection |

`renderTemplate(content, vars)` is also exported for standalone use.
