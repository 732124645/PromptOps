<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTable, NEmpty, NSpin, NSpace, useMessage, useDialog } from 'naive-ui'
import { api, type Agent } from '../api/client'
import { useWorkspaceStore } from '../stores/workspace'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const workspace = useWorkspaceStore()

const agents = ref<Agent[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await api.listAgents({ workspace: workspace.currentId })
    agents.value = data.data
  } catch {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

watch(() => workspace.currentId, load)

function edit(a: Agent) {
  router.push({ name: 'agent-edit', params: { id: a.id } })
}

function remove(a: Agent) {
  dialog.warning({
    title: '删除 Agent',
    content: `确定删除 "${a.key}" 吗?`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.removeAgent(a.id)
      message.success('已删除')
      load()
    },
  })
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>Agents</h2>
      <n-button type="primary" @click="router.push({ name: 'agent-new' })">+ 新建 Agent</n-button>
    </div>

    <n-spin :show="loading">
      <n-empty v-if="!agents.length" description="暂无 Agent" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>名称</th>
            <th>提供方</th>
            <th>模型</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in agents" :key="a.id">
            <td><code>{{ a.key }}</code></td>
            <td>{{ a.name || '-' }}</td>
            <td>{{ a.provider }}</td>
            <td>{{ a.model || '-' }}</td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="edit(a)">编辑 / 运行</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(a)">删除</n-button>
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
  margin-bottom: 16px;
}
code {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
</style>
