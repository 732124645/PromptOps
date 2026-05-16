<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
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
const { t } = useI18n()

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

const providerOptions = computed(() => [
  { label: t('providers.mock'), value: 'mock' },
  { label: 'OpenAI', value: 'openai' },
  { label: 'Claude', value: 'claude' },
  { label: t('providers.ollama'), value: 'ollama' },
  { label: 'Gemini', value: 'gemini' },
])

const variables = ref<Record<string, string>>({})
const apiKey = ref('')
const running = ref(false)
const result = ref<PlaygroundResponse | null>(null)
const error = ref('')

const detectedVars = computed(() => {
  const matches = (form.value.prompt || '').match(/{{\s*([\w.]+)\s*}}/g) || []
  return [...new Set(matches.map((m) => m.replace(/[{}]/g, '').trim()))]
})
const needsKey = computed(() =>
  ['openai', 'claude', 'gemini'].includes(form.value.provider || 'mock'),
)

async function load() {
  if (isNew.value) return
  try {
    const { data } = await api.getAgent(id.value)
    form.value = data.data
  } catch {
    message.error(t('common.loadFailed'))
  }
}

async function save() {
  if (!form.value.key) {
    message.warning(t('common.keyRequired'))
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
      message.success(t('common.created'))
      router.replace({ name: 'agent-edit', params: { id: id.value } })
    } else {
      await api.updateAgent(id.value, form.value)
      message.success(t('common.saved'))
    }
  } catch {
    message.error(t('common.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function run() {
  if (isNew.value) {
    message.warning(t('common.saveFirst', { name: t('entity.agent') }))
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
    error.value = resp?.data?.error || t('common.runFailed')
  } finally {
    running.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ isNew ? t('agent.newTitle') : t('agent.editTitle', { key: form.key }) }}</h2>
      <n-space>
        <n-button @click="router.push('/agents')">{{ t('common.back') }}</n-button>
        <n-button type="primary" :loading="saving" @click="save">{{ t('common.save') }}</n-button>
      </n-space>
    </div>

    <n-card :title="t('agent.configCard')" :bordered="true">
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="Key">
              <n-input v-model:value="form.key" :disabled="!isNew" placeholder="support.agent" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.name')">
              <n-input v-model:value="form.name" :placeholder="t('agent.namePlaceholder')" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('agent.providerLabel')">
              <n-select v-model:value="form.provider" :options="providerOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.model')">
              <n-input v-model:value="form.model" :placeholder="t('agent.modelPlaceholder')" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item :span="2">
            <n-form-item :label="t('agent.descriptionLabel')">
              <n-input
                v-model:value="form.description"
                :placeholder="t('agent.descriptionPlaceholder')"
              />
            </n-form-item>
          </n-grid-item>
        </n-grid>
        <n-form-item label="Prompt">
          <n-input
            v-model:value="form.prompt"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 8, maxRows: 20 }"
            :placeholder="t('agent.promptPlaceholder')"
          />
        </n-form-item>
      </n-form>
    </n-card>

    <n-card :title="t('agent.runCard')" :bordered="true" style="margin-top: 16px">
      <p v-if="isNew" class="muted">{{ t('agent.runSaveFirst') }}</p>
      <template v-else>
        <n-form>
          <n-form-item v-if="detectedVars.length" :label="t('agent.variables')">
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
          <n-form-item v-if="needsKey" label="API Key">
            <n-input
              v-model:value="apiKey"
              type="password"
              show-password-on="click"
              :placeholder="t('agent.apiKeyPlaceholder')"
            />
          </n-form-item>
          <n-button type="primary" :loading="running" @click="run">
            {{ t('agent.runAgent') }}
          </n-button>
        </n-form>

        <n-alert v-if="error" type="error" :title="t('common.runFailed')" style="margin-top: 14px">
          {{ error }}
        </n-alert>
        <template v-if="result">
          <div class="result-label">{{ t('result.rendered') }}</div>
          <pre class="block">{{ result.rendered }}</pre>
          <div class="result-label">
            {{ t('result.modelOutput') }} · {{ result.result.provider }} / {{ result.result.model }}
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
