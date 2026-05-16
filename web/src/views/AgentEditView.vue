<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton,
  NInput,
  NSelect,
  NSpace,
  NCard,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NTag,
  NAlert,
  useMessage,
} from 'naive-ui'
import { api, type Agent, type PlaygroundResponse } from '../api/client'
import { useWorkspaceStore } from '../stores/workspace'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const workspace = useWorkspaceStore()

const isNew = computed(() => !route.params.id)
const id = ref<string>((route.params.id as string) || '')
const saving = ref(false)

const form = ref<Partial<Agent>>({
  key: '',
  name: '',
  description: '',
  prompt: '',
  provider: 'mock',
  model: '',
})

const providerOptions = [
  { label: 'mock — 离线,无需 API Key', value: 'mock' },
  { label: 'OpenAI', value: 'openai' },
  { label: 'Claude', value: 'claude' },
  { label: 'Ollama — 本地', value: 'ollama' },
  { label: 'Gemini', value: 'gemini' },
]

const variables = ref<Record<string, string>>({})
const apiKey = ref('')
const running = ref(false)
const result = ref<PlaygroundResponse | null>(null)
const error = ref('')

const detectedVars = computed(() => {
  const matches = (form.value.prompt || '').match(/{{\s*([\w.]+)\s*}}/g) || []
  return [...new Set(matches.map((m) => m.replace(/[{}]/g, '').trim()))]
})
const needsKey = computed(() => ['openai', 'claude', 'gemini'].includes(form.value.provider || 'mock'))

async function load() {
  if (isNew.value) return
  try {
    const { data } = await api.getAgent(id.value)
    form.value = data.data
  } catch {
    message.error('加载失败')
  }
}

async function save() {
  if (!form.value.key) {
    message.warning('Key 不能为空')
    return
  }
  saving.value = true
  try {
    if (isNew.value) {
      const { data } = await api.createAgent({
        ...form.value,
        workspace_id: workspace.currentId,
      })
      id.value = data.data.id
      message.success('已创建')
      router.replace({ name: 'agent-edit', params: { id: id.value } })
    } else {
      await api.updateAgent(id.value, form.value)
      message.success('已保存')
    }
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function run() {
  if (isNew.value) {
    message.warning('请先保存 Agent')
    return
  }
  running.value = true
  error.value = ''
  result.value = null
  try {
    const { data } = await api.runAgent(id.value, {
      variables: variables.value,
      api_key: apiKey.value || undefined,
    })
    result.value = data
  } catch (e) {
    const resp = (e as { response?: { data?: { error?: string } } }).response
    error.value = resp?.data?.error || '运行失败'
  } finally {
    running.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ isNew ? '新建 Agent' : `Agent: ${form.key}` }}</h2>
      <n-space>
        <n-button @click="router.push('/agents')">返回</n-button>
        <n-button type="primary" :loading="saving" @click="save">保存</n-button>
      </n-space>
    </div>

    <n-card title="配置" :bordered="true">
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="Key">
              <n-input v-model:value="form.key" :disabled="!isNew" placeholder="support.agent" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="名称">
              <n-input v-model:value="form.name" placeholder="Agent 名称" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="模型提供方">
              <n-select v-model:value="form.provider" :options="providerOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="模型">
              <n-input v-model:value="form.model" placeholder="留空使用默认模型" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item :span="2">
            <n-form-item label="描述">
              <n-input v-model:value="form.description" placeholder="这个 Agent 做什么" />
            </n-form-item>
          </n-grid-item>
        </n-grid>
        <n-form-item label="Prompt">
          <n-input
            v-model:value="form.prompt"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 8, maxRows: 20 }"
            placeholder="Agent 的 Prompt,使用 双花括号变量 作为占位符"
          />
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="运行" :bordered="true" style="margin-top: 16px">
      <p v-if="isNew" class="muted">请先保存 Agent 再运行。</p>
      <template v-else>
        <n-form>
          <n-form-item v-if="detectedVars.length" label="变量">
            <n-space vertical style="width: 100%">
              <div v-for="v in detectedVars" :key="v" class="var-row">
                <n-tag type="warning" size="small">{{ v }}</n-tag>
                <n-input v-model:value="variables[v]" :placeholder="`${v} 的值`" />
              </div>
            </n-space>
          </n-form-item>
          <n-form-item v-if="needsKey" label="API Key">
            <n-input
              v-model:value="apiKey"
              type="password"
              show-password-on="click"
              placeholder="该提供方需要 API Key"
            />
          </n-form-item>
          <n-button type="primary" :loading="running" @click="run">运行 Agent</n-button>
        </n-form>

        <n-alert v-if="error" type="error" title="运行失败" style="margin-top: 14px">
          {{ error }}
        </n-alert>
        <template v-if="result">
          <div class="result-label">渲染后的 Prompt</div>
          <pre class="block">{{ result.rendered }}</pre>
          <div class="result-label">
            模型输出 · {{ result.result.provider }} / {{ result.result.model }}
          </div>
          <pre class="block">{{ result.result.output }}</pre>
        </template>
      </template>
    </n-card>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.muted {
  color: #888;
  font-size: 13px;
}
.var-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.result-label {
  font-size: 13px;
  color: #888;
  margin: 12px 0 6px;
}
.block {
  margin: 0;
  padding: 12px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.05);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
