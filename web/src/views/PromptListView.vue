<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
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

const prompts = ref<Prompt[]>([])
const q = ref('')
const env = ref<string | null>(null)
const loading = ref(false)
const live = ref(false)

const envOptions = [
  { label: '全部环境', value: '' },
  { label: 'dev', value: 'dev' },
  { label: 'test', value: 'test' },
  { label: 'prod', value: 'prod' },
]

let socket: WebSocket | null = null

async function load() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (q.value.trim()) params.q = q.value.trim()
    if (env.value) params.env = env.value
    const { data } = await api.list(params)
    prompts.value = data.data
  } catch {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function tagsOf(p: Prompt): string[] {
  return (p.tags || '')
    .split(',')
    .map((t) => t.trim())
    .filter(Boolean)
}

function edit(p: Prompt) {
  router.push({ name: 'prompt-edit', params: { id: p.id } })
}

function remove(p: Prompt) {
  dialog.warning({
    title: '删除 Prompt',
    content: `确定删除 "${p.key}" (${p.env}) 吗?此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.remove(p.id)
      message.success('已删除')
      load()
    },
  })
}

function connectWS() {
  try {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    socket = new WebSocket(`${proto}://${location.host}/ws`)
    socket.onopen = () => (live.value = true)
    socket.onclose = () => (live.value = false)
    socket.onmessage = () => {
      message.info('检测到 Prompt 变更,已刷新')
      load()
    }
  } catch {
    live.value = false
  }
}

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
        <h2>Prompts</h2>
        <n-tag v-if="live" type="success" size="small" round>● 热更新已连接</n-tag>
        <n-tag v-else size="small" round>○ 离线</n-tag>
      </div>
      <n-button type="primary" @click="router.push({ name: 'prompt-new' })">
        + 新建 Prompt
      </n-button>
    </div>

    <n-space class="filters">
      <n-input
        v-model:value="q"
        placeholder="搜索 key / 名称 / 内容"
        clearable
        style="width: 280px"
        @keyup.enter="load"
      />
      <n-select
        v-model:value="env"
        :options="envOptions"
        placeholder="环境"
        style="width: 140px"
        @update:value="load"
      />
      <n-button @click="load">搜索</n-button>
    </n-space>

    <n-spin :show="loading">
      <n-empty v-if="!prompts.length" description="暂无 Prompt" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>名称</th>
            <th>环境</th>
            <th>版本</th>
            <th>分类</th>
            <th>标签</th>
            <th>模型</th>
            <th>操作</th>
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
                <n-tag v-for="t in tagsOf(p)" :key="t" size="small" type="info">{{ t }}</n-tag>
                <span v-if="!tagsOf(p).length">-</span>
              </n-space>
            </td>
            <td>{{ p.model || '-' }}</td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="edit(p)">编辑</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(p)">删除</n-button>
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
