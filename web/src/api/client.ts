import axios from 'axios'
import { useAuthStore } from '../stores/auth'
import router from '../router'

export interface Prompt {
  id: string
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
  login: (token: string) => http.post('/login', { token }),
  list: (params: Record<string, string>) => http.get<{ data: Prompt[] }>('/prompts', { params }),
  get: (id: string) => http.get<{ data: Prompt }>(`/prompts/${id}`),
  create: (p: Partial<Prompt>) => http.post<{ data: Prompt }>('/prompts', p),
  update: (id: string, p: Partial<Prompt>) => http.put<{ data: Prompt }>(`/prompts/${id}`, p),
  remove: (id: string) => http.delete(`/prompts/${id}`),
  versions: (id: string) => http.get<{ data: PromptVersion[] }>(`/prompts/${id}/versions`),
  publish: (id: string) => http.post('/prompts/publish', { id }),
  rollback: (id: string, version: string) => http.post('/prompts/rollback', { id, version }),
  providers: () => http.get<{ data: string[] }>('/playground/providers'),
  runPlayground: (payload: PlaygroundRequest) =>
    http.post<PlaygroundResponse>('/playground/run', payload),

  listAgents: () => http.get<{ data: Agent[] }>('/agents'),
  getAgent: (id: string) => http.get<{ data: Agent }>(`/agents/${id}`),
  createAgent: (a: Partial<Agent>) => http.post<{ data: Agent }>('/agents', a),
  updateAgent: (id: string, a: Partial<Agent>) => http.put<{ data: Agent }>(`/agents/${id}`, a),
  removeAgent: (id: string) => http.delete(`/agents/${id}`),
  runAgent: (id: string, payload: { variables: Record<string, string>; api_key?: string }) =>
    http.post<PlaygroundResponse>(`/agents/${id}/run`, payload),

  listWorkflows: () => http.get<{ data: Workflow[] }>('/workflows'),
  getWorkflow: (id: string) => http.get<{ data: Workflow }>(`/workflows/${id}`),
  createWorkflow: (w: Partial<Workflow>) => http.post<{ data: Workflow }>('/workflows', w),
  updateWorkflow: (id: string, w: Partial<Workflow>) =>
    http.put<{ data: Workflow }>(`/workflows/${id}`, w),
  removeWorkflow: (id: string) => http.delete(`/workflows/${id}`),
  runWorkflow: (id: string, payload: { variables: Record<string, string>; api_key?: string }) =>
    http.post<WorkflowRunResult>(`/workflows/${id}/run`, payload),
}

export default http
