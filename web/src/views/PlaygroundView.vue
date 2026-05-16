<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import {
  NCard,
  NSelect,
  NInput,
  NButton,
  NSpace,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NTag,
  NAlert,
  useMessage,
} from 'naive-ui'
import { api, type Prompt, type PlaygroundResponse } from '../api/client'

const message = useMessage()

const prompts = ref<Prompt[]>([])
const selectedId = ref<string | null>(null)
const content = ref('')
const variables = ref<Record<string, string>>({})

const provider = ref('mock')
const model = ref('')
const apiKey = ref('')
const baseUrl = ref('')

const running = ref(false)
const result = ref<PlaygroundResponse | null>(null)
const error = ref('')

const providerOptions = [
  { label: 'mock — 离线,无需 API Key', value: 'mock' },
  { label: 'OpenAI', value: 'openai' },
  { label: 'Claude', value: 'claude' },
  { label: 'Ollama — 本地', value: 'ollama' },
  { label: 'Gemini', value: 'gemini' },
]

const promptOptions = computed(() =>
  prompts.value.map((p) => ({ label: `${p.key} (${p.env})`, value: p.id })),
)

const needsKey = computed(() => ['openai', 'claude', 'gemini'].includes(provider.value))
const needsBaseUrl = computed(() => provider.value === 'ollama')

const detectedVars = computed(() => {
  const matches = content.value.match(/{{\s*([\w.]+)\s*}}/g) || []
  return [...new Set(matches.map((m) => m.replace(/[{}]/g, '').trim()))]
})

watch(
  detectedVars,
  (vars) => {
    const next: Record<string, string> = {}
    for (const v of vars) next[v] = variables.value[v] ?? ''
    variables.value = next
  },
  { immediate: true },
)

async function loadPrompts() {
  try {
    const { data } = await api.list({})
    prompts.value = data.data
  } catch {
    message.error('加载 Prompt 列表失败')
  }
}

async function loadPromptContent(id: string) {
  try {
    const { data } = await api.get(id)
    content.value = data.data.content
    if (data.data.model) model.value = data.data.model
  } catch {
    message.error('加载 Prompt 失败')
  }
}

async function run() {
  if (!content.value.trim()) {
    message.warning('Prompt 内容不能为空')
    return
  }
  running.value = true
  error.value = ''
  result.value = null
  try {
    const { data } = await api.runPlayground({
      content: content.value,
      variables: variables.value,
      provider: provider.value,
      model: model.value || undefined,
      api_key: apiKey.value || undefined,
      base_url: baseUrl.value || undefined,
    })
    result.value = data
  } catch (e) {
    const resp = (e as { response?: { data?: { error?: string } } }).response
    error.value = resp?.data?.error || '运行失败'
  } finally {
    running.value = false
  }
}

onMounted(loadPrompts)
</script>

<template>
  <div class="page">
    <h2>Playground</h2>

    <n-card :bordered="true">
      <n-form>
        <n-form-item label="从已有 Prompt 载入(可选)">
          <n-select
            v-model:value="selectedId"
            :options="promptOptions"
            placeholder="选择一个 Prompt"
            clearable
            @update:value="(v) => v && loadPromptContent(v)"
          />
        </n-form-item>

        <n-form-item label="Prompt 内容">
          <n-input
            v-model:value="content"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 6, maxRows: 18 }"
            placeholder="输入 Prompt,使用 双花括号变量 作为占位符"
          />
        </n-form-item>

        <n-form-item v-if="detectedVars.length" label="变量">
          <n-space vertical style="width: 100%">
            <div v-for="v in detectedVars" :key="v" class="var-row">
              <n-tag type="warning" size="small">{{ v }}</n-tag>
              <n-input v-model:value="variables[v]" :placeholder="`${v} 的值`" />
            </div>
          </n-space>
        </n-form-item>

        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="模型提供方">
              <n-select v-model:value="provider" :options="providerOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="模型">
              <n-input v-model:value="model" placeholder="留空使用默认模型" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item v-if="needsKey">
            <n-form-item label="API Key">
              <n-input
                v-model:value="apiKey"
                type="password"
                show-password-on="click"
                placeholder="该提供方需要 API Key"
              />
            </n-form-item>
          </n-grid-item>
          <n-grid-item v-if="needsBaseUrl">
            <n-form-item label="Base URL">
              <n-input v-model:value="baseUrl" placeholder="http://localhost:11434" />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-button type="primary" :loading="running" @click="run">运行</n-button>
      </n-form>
    </n-card>

    <n-alert v-if="error" type="error" title="运行失败" style="margin-top: 16px">
      {{ error }}
    </n-alert>

    <n-card v-if="result" title="结果" :bordered="true" style="margin-top: 16px">
      <div class="result-label">渲染后的 Prompt</div>
      <pre class="block">{{ result.rendered }}</pre>
      <div class="result-label">
        模型输出 · {{ result.result.provider }} / {{ result.result.model }}
      </div>
      <pre class="block">{{ result.result.output }}</pre>
    </n-card>
  </div>
</template>

<style scoped>
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
