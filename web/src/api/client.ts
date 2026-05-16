import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import router from '../router'

export interface Prompt {
  id: string
  workspace_id?: string
  key: string
  name: string
  content: string
  version: string
  env: string
  category: string
  tags: string
  model: string
  created_at?: string
  updated_at?: string
}

export interface PromptVersion {
  id: string
  prompt_id: string
  key: string
  version: string
  env: string
  content: string
  created_at: string
}

export interface PlaygroundRequest {
  id?: string
  content?: string
  variables?: Record<string, string>
  provider: string
  model?: string
  api_key?: string
  base_url?: string
}

export interface PlaygroundResponse {
  rendered: string
  result: {
    provider: string
    model: string
    output: string
  }
}

export interface Agent {
  id: string
  workspace_id?: string
  key: string
  name: string
  description: string
  prompt: string
  provider: string
  model: string
  created_at?: string
  updated_at?: string
}

export type WorkflowStepType = 'render' | 'model' | 'transform'

export interface WorkflowStep {
  name: string
  type: WorkflowStepType
  template?: string
  provider?: string
  model?: string
  op?: string
}

export interface Workflow {
  id: string
  workspace_id?: string
  key: string
  name: string
  description: string
  steps: WorkflowStep[]
  created_at?: string
  updated_at?: string
}

export interface WorkflowRunResult {
  output: string
  steps: { name: string; type: string; output: string }[]
}

export interface AuditEntry {
  id: string
  action: string
  resource: string
  resource_id: string
  key: string
  summary: string
  created_at: string
}

export interface RunEntry {
  id: string
  source: string
  ref_key: string
  provider: string
  model: string
  prompt_tokens: number
  output_tokens: number
  latency_ms: number
  status: string
  error: string
  created_at: string
}

export interface RunStats {
  total: number
  ok: number
  error: number
  prompt_tokens: number
  output_tokens: number
  avg_latency_ms: number
  by_provider: { provider: string; count: number }[]
}

export interface User {
  id: string
  username: string
  role: string
  created_at: string
}

export interface Rollout {
  id: string
  key: string
  env: string
  enabled: boolean
  variant_a: string
  variant_b: string
  weight_a: number
}

export interface Workspace {
  id: string
  name: string
  slug: string
  created_at: string
}

export interface LoginResponse {
  ok: boolean
  token: string
  role: string
  username: string
}

const http = axios.create({ baseURL: '/api' })

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

http.interceptors.response.use(
  (resp) => resp,
  (error) => {
    if (error.response?.status === 401 && router.currentRoute.value.name !== 'login') {
      useAuthStore().clear()
      router.push({ name: 'login' })
    }
    return Promise.reject(error)
  },
)

export const api = {
  login: (payload: { username?: string; password?: string; token?: string }) =>
    http.post<LoginResponse>('/login', payload),
  logout: () => http.post('/logout'),
  me: () => http.get<{ username: string; role: string }>('/me'),
  listUsers: () => http.get<{ data: User[] }>('/users'),
  createUser: (u: { username: string; password: string; role: string }) =>
    http.post<{ data: User }>('/users', u),
  updateUser: (id: string, u: { password?: string; role?: string }) =>
    http.put<{ data: User }>(`/users/${id}`, u),
  removeUser: (id: string) => http.delete(`/users/${id}`),
  list: (params: Record<string, string>) => http.get<{ data: Prompt[] }>('/prompts', { params }),
  get: (id: string) => http.get<{ data: Prompt }>(`/prompts/${id}`),
  create: (p: Partial<Prompt>) => http.post<{ data: Prompt }>('/prompts', p),
  update: (id: string, p: Partial<Prompt>) => http.put<{ data: Prompt }>(`/prompts/${id}`, p),
  remove: (id: string) => http.delete(`/prompts/${id}`),
  versions: (id: string) => http.get<{ data: PromptVersion[] }>(`/prompts/${id}/versions`),
  publish: (id: string) => http.post('/prompts/publish', { id }),
  rollback: (id: string, version: string) => http.post('/prompts/rollback', { id, version }),
  getRollout: (id: string) => http.get<{ data: Rollout | null }>(`/prompts/${id}/rollout`),
  setRollout: (
    id: string,
    payload: { enabled: boolean; variant_a: string; variant_b: string; weight_a: number },
  ) => http.put<{ data: Rollout }>(`/prompts/${id}/rollout`, payload),
  deleteRollout: (id: string) => http.delete(`/prompts/${id}/rollout`),
  providers: () => http.get<{ data: string[] }>('/playground/providers'),
  runPlayground: (payload: PlaygroundRequest) =>
    http.post<PlaygroundResponse>('/playground/run', payload),

  listWorkspaces: () => http.get<{ data: Workspace[] }>('/workspaces'),
  createWorkspace: (w: { name: string; slug?: string }) =>
    http.post<{ data: Workspace }>('/workspaces', w),
  removeWorkspace: (id: string) => http.delete(`/workspaces/${id}`),

  listAgents: (params?: Record<string, string>) =>
    http.get<{ data: Agent[] }>('/agents', { params }),
  getAgent: (id: string) => http.get<{ data: Agent }>(`/agents/${id}`),
  createAgent: (a: Partial<Agent>) => http.post<{ data: Agent }>('/agents', a),
  updateAgent: (id: string, a: Partial<Agent>) => http.put<{ data: Agent }>(`/agents/${id}`, a),
  removeAgent: (id: string) => http.delete(`/agents/${id}`),
  runAgent: (id: string, payload: { variables: Record<string, string>; api_key?: string }) =>
    http.post<PlaygroundResponse>(`/agents/${id}/run`, payload),

  listWorkflows: (params?: Record<string, string>) =>
    http.get<{ data: Workflow[] }>('/workflows', { params }),
  getWorkflow: (id: string) => http.get<{ data: Workflow }>(`/workflows/${id}`),
  createWorkflow: (w: Partial<Workflow>) => http.post<{ data: Workflow }>('/workflows', w),
  updateWorkflow: (id: string, w: Partial<Workflow>) =>
    http.put<{ data: Workflow }>(`/workflows/${id}`, w),
  removeWorkflow: (id: string) => http.delete(`/workflows/${id}`),
  runWorkflow: (id: string, payload: { variables: Record<string, string>; api_key?: string }) =>
    http.post<WorkflowRunResult>(`/workflows/${id}/run`, payload),

  listAudit: () => http.get<{ data: AuditEntry[] }>('/audit'),
  listRuns: () => http.get<{ data: RunEntry[] }>('/runs'),
  runStats: () => http.get<RunStats>('/runs/stats'),
}

export default http
