<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NButton, NTable, NEmpty, NSpin, NSpace, useMessage, useDialog } from 'naive-ui'
import { api, type Agent } from '../api/client'
import { useWorkspaceStore } from '../stores/workspace'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const workspace = useWorkspaceStore()
const { t } = useI18n()

const agents = ref<Agent[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await api.listAgents({ workspace: workspace.currentId })
    agents.value = data.data
  } catch {
    message.error(t('common.loadFailed'))
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
    title: t('agent.deleteTitle'),
    content: t('agent.deleteConfirm', { key: a.key }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await api.removeAgent(a.id)
      message.success(t('common.deleted'))
      load()
    },
  })
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ t('agent.title') }}</h2>
      <n-button type="primary" @click="router.push({ name: 'agent-new' })">
        {{ t('agent.newAgent') }}
      </n-button>
    </div>

    <n-spin :show="loading">
      <n-empty v-if="!agents.length" :description="t('agent.empty')" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>{{ t('common.name') }}</th>
            <th>{{ t('agent.colProvider') }}</th>
            <th>{{ t('common.model') }}</th>
            <th>{{ t('common.actions') }}</th>
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
                <n-button size="tiny" @click="edit(a)">{{ t('common.editRun') }}</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(a)">
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
  margin-bottom: 16px;
}
code {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
</style>
