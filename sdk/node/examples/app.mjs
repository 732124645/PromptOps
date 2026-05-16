/**
 * Realistic usage example for the PromptOps Node SDK.
 *
 * It runs a long-lived "support assistant" that:
 *   1. fetches a prompt from PromptOps and renders it with ticket variables,
 *   2. opens a WebSocket (watch) so the prompt hot-reloads when an editor
 *      publishes a change in the platform — no redeploy needed,
 *   3. stays alive, so the process appears in the platform UI under
 *      Observability -> Connected clients (app = "support-assistant").
 *
 * Run it:
 *   node examples/app.mjs
 *   # or, against a remote server:
 *   PROMPTOPS_SERVER=http://localhost:8080 node examples/app.mjs
 *
 * Then, in the PromptOps web UI, edit the prompt "support.reply" and save —
 * watch this process pick up the new version live.
 */
import { PromptOpsClient } from '../src/index.js'

const SERVER = process.env.PROMPTOPS_SERVER || 'http://localhost:8080'
const TOKEN = process.env.PROMPTOPS_TOKEN || 'promptops-dev-token'
const PROMPT_KEY = 'support.reply'
const NAMESPACE = 'prod'

// A sample support ticket. In a real app this comes from your queue / inbox.
const ticket = {
  product: 'PromptOps',
  tone: 'friendly and concise',
  message: 'I changed my prompt but my app still serves the old text — why?',
}

/** Minimal API call helper for the one-time setup step below. */
async function api(path, method = 'GET', body) {
  const res = await fetch(`${SERVER}${path}`, {
    method,
    headers: { Authorization: `Bearer ${TOKEN}`, 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok && res.status !== 404) {
    throw new Error(`${method} ${path} -> HTTP ${res.status}`)
  }
  return res
}

/**
 * Make sure the demo prompt exists so the example is runnable with no manual
 * setup. A real consumer app would NOT do this — prompts are authored in the
 * PromptOps UI; the SDK only reads them.
 */
async function ensureDemoPrompt() {
  const check = await api(`/api/sdk/prompts/${PROMPT_KEY}?env=${NAMESPACE}`)
  if (check.ok) return
  console.log(`· demo prompt "${PROMPT_KEY}" not found — creating it`)
  await api('/api/prompts', 'POST', {
    key: PROMPT_KEY,
    name: 'Support reply',
    content:
      'You are a support agent for {{product}}.\n' +
      'Reply to the customer in a {{tone}} tone.\n\n' +
      'Customer message:\n{{message}}',
    env: NAMESPACE,
    model: 'gpt-4o',
    workspace_id: 'default', // so it shows under the default workspace in the UI
  })
}

/** Render the prompt with the ticket and print it as the assistant would. */
async function buildReplyPrompt(client) {
  const rendered = await client.render(PROMPT_KEY, ticket)
  console.log('\n─── rendered prompt (feed this to your LLM) ───')
  console.log(rendered)
  console.log('───────────────────────────────────────────────\n')
  return rendered
}

async function main() {
  await ensureDemoPrompt()

  // appName makes this process identifiable in Connected clients.
  const client = new PromptOpsClient({
    server: SERVER,
    namespace: NAMESPACE,
    token: TOKEN,
    appName: 'support-assistant',
  })

  await buildReplyPrompt(client) // initial render

  // Hot-reload lifecycle.
  client.on('connect', () => console.log('· watch connected — live updates on'))
  client.on('disconnect', () => console.log('· watch disconnected'))
  client.on('error', (err) => console.error('· watch error:', err.message))
  client.on('update', async (e) => {
    console.log(`· prompt "${e.key}" hot-reloaded (event: ${e.event})`)
    await buildReplyPrompt(client) // re-render with the new content
  })

  client.watch()

  console.log('Assistant is running. Try this:')
  console.log(`  1. open the PromptOps UI and edit the prompt "${PROMPT_KEY}"`)
  console.log('  2. watch the re-render appear here automatically')
  console.log('  3. check Observability -> Connected clients (app: support-assistant)')
  console.log('Press Ctrl+C to stop.\n')

  // Keep the process alive; shut the WebSocket down cleanly on exit.
  process.on('SIGINT', () => {
    console.log('\n· shutting down')
    client.close()
    process.exit(0)
  })
}

main().catch((err) => {
  console.error('example failed:', err.message)
  console.error(`is the PromptOps server running at ${SERVER}?`)
  process.exit(1)
})
