export interface Prompt {
  key: string
  version: string
  env: string
  model: string
  content: string
}

export interface PromptOpsClientOptions {
  /** Base server URL, e.g. "http://localhost:8080". */
  server: string
  /** Environment to fetch prompts from. Defaults to "prod". */
  namespace?: string
  /** Bearer token. Defaults to "promptops-dev-token". */
  token?: string
}

export interface UpdateEvent {
  key: string
  prompt: Prompt
  event: string
}

export function renderTemplate(content: string, vars?: Record<string, unknown>): string

export class PromptOpsClient {
  readonly server: string
  readonly namespace: string
  readonly cache: Map<string, Prompt>

  constructor(options: PromptOpsClientOptions)

  getPrompt(key: string, opts?: { refresh?: boolean }): Promise<Prompt>
  render(key: string, vars?: Record<string, unknown>): Promise<string>

  on(event: 'update', cb: (e: UpdateEvent) => void): () => void
  on(event: 'connect' | 'disconnect', cb: () => void): () => void
  on(event: 'error', cb: (e: Error) => void): () => void

  watch(): WebSocket
  close(): void
}

export default PromptOpsClient
