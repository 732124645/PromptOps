# Hot Reload

## What it is

**Hot Reload** pushes prompt changes to running applications in real time.
When an editor saves a prompt, every connected SDK refreshes its cache
automatically — no restart, no redeploy.

## Why it matters

Without it, updating a prompt means restarting your AI service. With it, the
prompt is a runtime resource: change the wording, hit save, and production
picks it up within milliseconds.

## How it works

1. The SDK calls `watch()`, which opens a WebSocket to the server at `/ws`.
2. When a prompt is created, updated or deleted, the server **broadcasts** a
   small event — `{ type, key, env, version }` — to every connected client.
3. Each SDK checks the event: if it has that prompt `key` cached, it re-fetches
   the latest content and fires an `update` event for your code to react to.

The web UI connects the same way — the Prompts page shows a
**"Live updates connected"** badge when its WebSocket is open.

## Using it from the SDK

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
  appName: 'checkout-service', // shown in the Connected clients panel
})

// prime the cache
await client.render('sql.generator', { question: 'list all users' })

// react to live changes
client.on('update', (e) => console.log('prompt hot-reloaded:', e.key))
client.on('connect', () => console.log('watching'))
client.watch()
```

After `watch()`, any edit to `sql.generator` in the UI triggers the `update`
callback and refreshes the cached content automatically.

## Seeing who is connected

Every `watch()` connection is tracked. The Observability page lists them under
**Connected clients** — see [Observability](./observability).

## Next steps

- [Prompts](./prompts) — the resource that gets hot-reloaded.
- [SDK](/sdk/node) — full client reference.
