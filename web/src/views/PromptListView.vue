<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useWorkspaceStore } from '../stores/workspace'
import { useAuthStore } from '../stores/auth'
import {
  NButton,
  NInput,
  NSelect,
  NTag,
  NSpace,
  NTable,
  NEmpty,
  NSpin,
  useMessage,
  useDialog,
} from 'naive-ui'
import { api, type Prompt } from '../api/client'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const workspace = useWorkspaceStore()
const auth = useAuthStore()
const { t } = useI18n()

const prompts = ref<Prompt[]>([])
const q = ref('')
const env = ref<string | null>(null)
const loading = ref(false)
const live = ref(false)

const envOptions = computed(() => [
  { label: t('prompt.allEnvs'), value: '' },
  { label: 'dev', value: 'dev' },
  { label: 'test', value: 'test' },
  { label: 'prod', value: 'prod' },
])

let socket: WebSocket | null = null

async function load() {
  loading.value = true
  try {
    const params: Record<string, string> = { workspace: workspace.currentId }
    if (q.value.trim()) params.q = q.value.trim()
    if (env.value) params.env = env.value
    const { data } = await api.list(params)
    prompts.value = data.data
  } catch {
    message.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

function tagsOf(p: Prompt): string[] {
  return (p.tags || '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

function edit(p: Prompt) {
  router.push({ name: 'prompt-edit', params: { id: p.id } })
}

function remove(p: Prompt) {
  dialog.warning({
    title: t('prompt.deleteTitle'),
    content: t('prompt.deleteConfirm', { key: p.key, env: p.env }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await api.remove(p.id)
      message.success(t('common.deleted'))
      load()
    },
  })
}

function connectWS() {
  try {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    const params = new URLSearchParams({
      token: auth.token || '',
      client: 'browser',
      app: 'promptops-ui',
    })
    socket = new WebSocket(`${proto}://${location.host}/ws?${params}`)
    socket.onopen = () => (live.value = true)
    socket.onclose = () => (live.value = false)
    socket.onmessage = () => {
      message.info(t('prompt.changeDetected'))
      load()
    }
  } catch {
    live.value = false
  }
}

watch(() => workspace.currentId, load)

onMounted(() => {
  load()
  connectWS()
})

onBeforeUnmount(() => socket?.close())
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <div class="left">
        <h2>{{ t('prompt.title') }}</h2>
        <n-tag v-if="live" type="success" size="small" round>{{ t('prompt.liveConnected') }}</n-tag>
        <n-tag v-else size="small" round>{{ t('prompt.offline') }}</n-tag>
      </div>
      <n-button type="primary" @click="router.push({ name: 'prompt-new' })">
        {{ t('prompt.newPrompt') }}
      </n-button>
    </div>

    <n-space class="filters">
      <n-input
        v-model:value="q"
        :placeholder="t('prompt.searchPlaceholder')"
        clearable
        style="width: 280px"
        @keyup.enter="load"
      />
      <n-select
        v-model:value="env"
        :options="envOptions"
        :placeholder="t('common.env')"
        style="width: 140px"
        @update:value="load"
      />
      <n-button @click="load">{{ t('common.search') }}</n-button>
    </n-space>

    <n-spin :show="loading">
      <n-empty v-if="!prompts.length" :description="t('prompt.empty')" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>{{ t('common.name') }}</th>
            <th>{{ t('common.env') }}</th>
            <th>{{ t('common.version') }}</th>
            <th>{{ t('prompt.colCategory') }}</th>
            <th>{{ t('prompt.colTags') }}</th>
            <th>{{ t('common.model') }}</th>
            <th>{{ t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in prompts" :key="p.id">
            <td><code>{{ p.key }}</code></td>
            <td>{{ p.name || '-' }}</td>
            <td><n-tag size="small">{{ p.env }}</n-tag></td>
            <td>{{ p.version }}</td>
            <td>{{ p.category || '-' }}</td>
            <td>
              <n-space :size="4">
                <n-tag v-for="tg in tagsOf(p)" :key="tg" size="small" type="info">{{ tg }}</n-tag>
                <span v-if="!tagsOf(p).length">-</span>
              </n-space>
            </td>
            <td>{{ p.model || '-' }}</td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="edit(p)">{{ t('common.edit') }}</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(p)">
                  {{ t('common.delete') }}
                </n-button>
              </n-space>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-spin>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.toolbar .left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.filters {
  margin: 8px 0 16px;
}
code {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
</style>
