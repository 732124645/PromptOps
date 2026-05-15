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
}

export default http
