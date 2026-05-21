// Minimal PromptOps consumer.
//
// Renders one prompt every few seconds and prints the result, while a
// WebSocket subscription keeps the local cache in sync with the server. Edit
// the prompt in the PromptOps Web UI and click Publish — the next render
// here picks up the new content without restarting.

import { PromptOpsClient } from '@promptops/client'

const SERVER = process.env.PROMPTOPS_SERVER || 'http://localhost:8080'
const TOKEN = process.env.PROMPTOPS_TOKEN || 'promptops-dev-token'
const NAMESPACE = process.env.PROMPTOPS_NAMESPACE || 'prod'
const PROMPT_KEY = process.env.PROMPT_KEY || 'demo.greeting'
const INTERVAL_MS = Number(process.env.RENDER_INTERVAL_MS || 3000)

const client = new PromptOpsClient({
  server: SERVER,
  namespace: NAMESPACE,
  token: TOKEN,
  appName: 'example-node-hot-reload',
})

client.on('connect', () => console.log(`[ws] connected to ${SERVER}/ws`))
client.on('disconnect', () => console.log('[ws] disconnected'))
client.on('error', (err) => console.error('[ws] error:', err.message))
client.on('update', (e) =>
  console.log(`[update] "${e.key}" -> version ${e.prompt.version} (${e.event})`),
)

client.watch()

async function renderOnce() {
  try {
    const text = await client.render(PROMPT_KEY, {
      name: 'World',
      now: new Date().toISOString(),
    })
    console.log(`\n[${new Date().toLocaleTimeString()}] ${PROMPT_KEY}:`)
    console.log(text)
  } catch (err) {
    console.error(
      `[render] failed (${err.message}). ` +
        `Make sure a prompt with key "${PROMPT_KEY}" exists in env "${NAMESPACE}".`,
    )
  }
}

await renderOnce()
const timer = setInterval(renderOnce, INTERVAL_MS)

const shutdown = () => {
  clearInterval(timer)
  client.close()
  process.exit(0)
}
process.on('SIGINT', shutdown)
process.on('SIGTERM', shutdown)
