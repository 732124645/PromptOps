<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
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
const { t } = useI18n()

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

const providerOptions = computed(() => [
  { label: t('providers.mock'), value: 'mock' },
  { label: 'OpenAI', value: 'openai' },
  { label: 'Claude', value: 'claude' },
  { label: t('providers.ollama'), value: 'ollama' },
  { label: 'Gemini', value: 'gemini' },
])

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
    message.error(t('playground.loadListFailed'))
  }
}

async function loadPromptContent(id: string) {
  try {
    const { data } = await api.get(id)
    content.value = data.data.content
    if (data.data.model) model.value = data.data.model
  } catch {
    message.error(t('common.loadFailed'))
  }
}

async function run() {
  if (!content.value.trim()) {
    message.warning(t('playground.contentRequired'))
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
    error.value = resp?.data?.error || t('common.runFailed')
  } finally {
    running.value = false
  }
}

onMounted(loadPrompts)
</script>

<template>
  <div class="page">
    <h2>{{ t('playground.title') }}</h2>

    <n-card :bordered="true">
      <n-form>
        <n-form-item :label="t('playground.loadFromPrompt')">
          <n-select
            v-model:value="selectedId"
            :options="promptOptions"
            :placeholder="t('playground.selectPrompt')"
            clearable
            @update:value="(v) => v && loadPromptContent(v)"
          />
        </n-form-item>

        <n-form-item :label="t('playground.contentLabel')">
          <n-input
            v-model:value="content"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 6, maxRows: 18 }"
            :placeholder="t('playground.contentPlaceholder')"
          />
        </n-form-item>

        <n-form-item v-if="detectedVars.length" :label="t('playground.variables')">
          <n-space vertical style="width: 100%">
            <div v-for="v in detectedVars" :key="v" class="var-row">
              <n-tag type="warning" size="small">{{ v }}</n-tag>
              <n-input
                v-model:value="variables[v]"
                :placeholder="t('agent.varValuePlaceholder', { name: v })"
              />
            </div>
          </n-space>
        </n-form-item>

        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item :label="t('playground.providerLabel')">
              <n-select v-model:value="provider" :options="providerOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.model')">
              <n-input v-model:value="model" :placeholder="t('playground.modelPlaceholder')" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item v-if="needsKey">
            <n-form-item label="API Key">
              <n-input
                v-model:value="apiKey"
                type="password"
                show-password-on="click"
                :placeholder="t('playground.apiKeyPlaceholder')"
              />
            </n-form-item>
          </n-grid-item>
          <n-grid-item v-if="needsBaseUrl">
            <n-form-item label="Base URL">
              <n-input v-model:value="baseUrl" placeholder="http://localhost:11434" />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-button type="primary" :loading="running" @click="run">{{ t('common.run') }}</n-button>
      </n-form>
    </n-card>

    <n-alert v-if="error" type="error" :title="t('common.runFailed')" style="margin-top: 16px">
      {{ error }}
    </n-alert>

    <n-card
      v-if="result"
      :title="t('playground.resultCard')"
      :bordered="true"
      style="margin-top: 16px"
    >
      <div class="result-label">{{ t('result.rendered') }}</div>
      <pre class="block">{{ result.rendered }}</pre>
      <div class="result-label">
        {{ t('result.modelOutput') }} · {{ result.result.provider }} / {{ result.result.model }}
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
