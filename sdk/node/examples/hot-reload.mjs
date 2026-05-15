// End-to-end check: fetch + render + WebSocket hot-reload against a running
// PromptOps server. Used by CI's integration job.
import { PromptOpsClient } from '../src/index.js'

const SERVER = process.env.PROMPTOPS_SERVER || 'http://localhost:8080'
const TOKEN = process.env.PROMPTOPS_TOKEN || 'promptops-dev-token'

async function api(path, method = 'GET', body) {
  const res = await fetch(`${SERVER}${path}`, {
    method,
    headers: { Authorization: `Bearer ${TOKEN}`, 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) throw new Error(`${method} ${path} -> HTTP ${res.status}`)
  return res.json()
}

function fail(msg) {
  console.error('FAIL:', msg)
  process.exit(1)
}

const created = await api('/api/prompts', 'POST', {
  key: 'sdk.demo',
  name: 'SDK Demo',
  content: 'Hello {{name}}, v1',
  env: 'prod',
  model: 'gpt-4o',
})
const id = created.data.id

const client = new PromptOpsClient({ server: SERVER, namespace: 'prod', token: TOKEN })

const rendered = await client.render('sdk.demo', { name: 'Ada' })
if (rendered !== 'Hello Ada, v1') fail(`render mismatch: ${rendered}`)
console.log('ok - getPrompt + render ->', rendered)

const updateEvent = new Promise((resolve, reject) => {
  const timer = setTimeout(() => reject(new Error('no hot-reload event within 5s')), 5000)
  client.on('update', (e) => {
    clearTimeout(timer)
    resolve(e)
  })
})

client.watch()
await new Promise((r) => setTimeout(r, 500)) // allow the WebSocket to connect

await api(`/api/prompts/${id}`, 'PUT', {
  name: 'SDK Demo',
  content: 'Hello {{name}}, v2',
  version: 'v2',
  env: 'prod',
})

let evt
try {
  evt = await updateEvent
} catch (err) {
  fail(err.message)
}

const after = await client.render('sdk.demo', { name: 'Ada' })
if (after !== 'Hello Ada, v2') fail(`cache not refreshed after hot-reload: ${after}`)
console.log(`ok - hot-reload via WebSocket -> ${after} (event: ${evt.event})`)

client.close()
await api(`/api/prompts/${id}`, 'DELETE')
console.log('Node SDK integration test passed')
process.exit(0)
