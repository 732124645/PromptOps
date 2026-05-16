import WebSocket from 'ws'

const DEFAULT_TOKEN = 'promptops-dev-token'
const VAR_RE = /{{\s*([\w.]+)\s*}}/g

/**
 * Replace {{variable}} placeholders. Unknown variables are left untouched.
 */
export function renderTemplate(content, vars = {}) {
  return String(content ?? '').replace(VAR_RE, (full, key) =>
    Object.prototype.hasOwnProperty.call(vars, key) ? String(vars[key]) : full,
  )
}

/**
 * Client for the PromptOps runtime API. Fetches prompts by key, caches them,
 * and (via watch()) keeps the cache fresh over a WebSocket connection.
 */
export class PromptOpsClient {
  constructor(options = {}) {
    if (!options.server) throw new Error('PromptOps: "server" is required')
    this.server = String(options.server).replace(/\/+$/, '')
    this.namespace = options.namespace || 'prod'
    this.token = options.token || DEFAULT_TOKEN
    this.appName = options.appName || ''
    this.cache = new Map()
    this._listeners = new Map()
    this._ws = null
  }

  async getPrompt(key, { refresh = false } = {}) {
    if (!refresh && this.cache.has(key)) return this.cache.get(key)
    const url =
      `${this.server}/api/sdk/prompts/${encodeURIComponent(key)}` +
      `?env=${encodeURIComponent(this.namespace)}`
    const res = await fetch(url, { headers: { Authorization: `Bearer ${this.token}` } })
    if (!res.ok) {
      throw new Error(`PromptOps: failed to fetch "${key}" (HTTP ${res.status})`)
    }
    const prompt = await res.json()
    this.cache.set(key, prompt)
    return prompt
  }

  async render(key, vars = {}) {
    const prompt = await this.getPrompt(key)
    return renderTemplate(prompt.content, vars)
  }

  on(event, cb) {
    if (!this._listeners.has(event)) this._listeners.set(event, new Set())
    this._listeners.get(event).add(cb)
    return () => this._listeners.get(event)?.delete(cb)
  }

  _emit(event, payload) {
    for (const cb of this._listeners.get(event) || []) {
      try {
        cb(payload)
      } catch {
        /* listener errors must not break the SDK */
      }
    }
  }

  /**
   * Open a WebSocket to the server and refresh any cached prompt whose key
   * appears in a hot-reload event. Emits 'update', 'connect', 'disconnect'.
   */
  watch() {
    if (this._ws) return this._ws
    // Identify this connection to the server's client registry.
    const params = new URLSearchParams({ client: 'node-sdk', namespace: this.namespace })
    if (this.appName) params.set('app', this.appName)
    const wsUrl = `${this.server.replace(/^http/, 'ws')}/ws?${params}`
    const ws = new WebSocket(wsUrl)
    this._ws = ws
    ws.addEventListener('open', () => this._emit('connect'))
    ws.addEventListener('close', () => {
      this._emit('disconnect')
      this._ws = null
    })
    ws.addEventListener('error', () => this._emit('error', new Error('PromptOps: websocket error')))
    ws.addEventListener('message', async (ev) => {
      // ev.data is a string on global WebSocket, but a Buffer on the `ws`
      // package — normalise both to text before parsing.
      const raw = ev.data
      let text = ''
      if (typeof raw === 'string') text = raw
      else if (Buffer.isBuffer(raw)) text = raw.toString('utf8')
      else if (raw instanceof ArrayBuffer) text = Buffer.from(raw).toString('utf8')
      else if (Array.isArray(raw)) text = Buffer.concat(raw).toString('utf8')
      let evt
      try {
        evt = JSON.parse(text)
      } catch {
        return
      }
      if (evt && evt.key && this.cache.has(evt.key)) {
        try {
          const prompt = await this.getPrompt(evt.key, { refresh: true })
          this._emit('update', { key: evt.key, prompt, event: evt.type })
        } catch (err) {
          this._emit('error', err)
        }
      }
    })
    return ws
  }

  close() {
    this._ws?.close()
    this._ws = null
  }
}

export default PromptOpsClient
