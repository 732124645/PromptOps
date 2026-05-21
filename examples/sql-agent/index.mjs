// Natural-language → SQL agent driven by a PromptOps-managed prompt.
//
// Flow:
//   1. Fetch the `sql.generator` prompt from PromptOps via the SDK
//   2. POST it (with the user's question + dialect) to /api/playground/run,
//      which renders + calls the chosen model provider and returns the result
//   3. Print the generated SQL
//
// The whole point: to change how SQL is generated, edit the prompt in the
// PromptOps Web UI and click Publish. No code change, no redeploy.

import { PromptOpsClient } from '@promptops/client'

const SERVER = process.env.PROMPTOPS_SERVER || 'http://localhost:8080'
const TOKEN = process.env.PROMPTOPS_TOKEN || 'promptops-dev-token'
const NAMESPACE = process.env.PROMPTOPS_NAMESPACE || 'prod'
const PROMPT_KEY = process.env.PROMPT_KEY || 'sql.generator'
const DIALECT = process.env.DIALECT || 'mysql'
const PROVIDER = process.env.PROVIDER || 'mock'
const MODEL = process.env.MODEL || 'mock-1'
const PROVIDER_API_KEY = process.env.PROVIDER_API_KEY || ''
const PROVIDER_BASE_URL = process.env.PROVIDER_BASE_URL || ''

const question = process.argv.slice(2).join(' ').trim()
if (!question) {
  console.error(
    'Usage: node index.mjs "<natural-language question>"\n' +
      '       e.g. node index.mjs "Top 10 customers by revenue this month"',
  )
  process.exit(1)
}

const client = new PromptOpsClient({
  server: SERVER,
  namespace: NAMESPACE,
  token: TOKEN,
  appName: 'example-sql-agent',
})

// Subscribe so the local prompt cache stays fresh if you publish a new
// version while the agent is running.
client.on('update', (e) =>
  console.error(`[update] "${e.key}" -> version ${e.prompt.version}`),
)
client.watch()
// Give the WebSocket a brief moment to connect before we make the first
// request. Not strictly required, but keeps the [update] notice in order.
await new Promise((r) => setTimeout(r, 200))

const prompt = await client.getPrompt(PROMPT_KEY).catch((err) => {
  console.error(
    `[error] could not load prompt "${PROMPT_KEY}" from ${SERVER}: ${err.message}\n` +
      `        Create it first — see prompts/sql-generator.md.`,
  )
  process.exit(2)
})

const playgroundRes = await fetch(`${SERVER}/api/playground/run`, {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    content: prompt.content,
    variables: { dialect: DIALECT, question },
    provider: PROVIDER,
    model: MODEL,
    api_key: PROVIDER_API_KEY,
    base_url: PROVIDER_BASE_URL,
  }),
})

if (!playgroundRes.ok) {
  const detail = await playgroundRes.text().catch(() => '')
  console.error(
    `[error] /api/playground/run -> HTTP ${playgroundRes.status} ${detail}`,
  )
  process.exit(3)
}

const { rendered, result } = await playgroundRes.json()

console.log('--- Rendered prompt (what the model saw) ---')
console.log(rendered)
console.log()
console.log(`--- Model output (${result.provider} / ${result.model}) ---`)
console.log(result.output)

client.close()
