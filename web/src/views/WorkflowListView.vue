<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTable, NEmpty, NSpin, NSpace, NTag, useMessage, useDialog } from 'naive-ui'
import { api, type Workflow } from '../api/client'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const workflows = ref<Workflow[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await api.listWorkflows()
    workflows.value = data.data
  } catch {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function edit(w: Workflow) {
  router.push({ name: 'workflow-edit', params: { id: w.id } })
}

function remove(w: Workflow) {
  dialog.warning({
    title: '删除 Workflow',
    content: `确定删除 "${w.key}" 吗?`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.removeWorkflow(w.id)
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
      <h2>Workflows</h2>
      <n-button type="primary" @click="router.push({ name: 'workflow-new' })">
        + 新建 Workflow
      </n-button>
    </div>

    <n-spin :show="loading">
      <n-empty v-if="!workflows.length" description="暂无 Workflow" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>名称</th>
            <th>步骤数</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in workflows" :key="w.id">
            <td><code>{{ w.key }}</code></td>
            <td>{{ w.name || '-' }}</td>
            <td><n-tag size="small">{{ w.steps.length }} 步</n-tag></td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="edit(w)">编辑 / 运行</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(w)">删除</n-button>
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
