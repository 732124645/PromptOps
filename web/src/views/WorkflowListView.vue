<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NButton, NTable, NEmpty, NSpin, NSpace, NTag, useMessage, useDialog } from 'naive-ui'
import { api, type Workflow } from '../api/client'
import { useWorkspaceStore } from '../stores/workspace'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const workspace = useWorkspaceStore()
const { t } = useI18n()

const workflows = ref<Workflow[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await api.listWorkflows({ workspace: workspace.currentId })
    workflows.value = data.data
  } catch {
    message.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

watch(() => workspace.currentId, load)

function edit(w: Workflow) {
  router.push({ name: 'workflow-edit', params: { id: w.id } })
}

function remove(w: Workflow) {
  dialog.warning({
    title: t('workflow.deleteTitle'),
    content: t('workflow.deleteConfirm', { key: w.key }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await api.removeWorkflow(w.id)
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
      <h2>{{ t('workflow.title') }}</h2>
      <n-button type="primary" @click="router.push({ name: 'workflow-new' })">
        {{ t('workflow.newWorkflow') }}
      </n-button>
    </div>

    <n-spin :show="loading">
      <n-empty
        v-if="!workflows.length"
        :description="t('workflow.empty')"
        style="margin: 48px 0"
      />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>Key</th>
            <th>{{ t('common.name') }}</th>
            <th>{{ t('workflow.colSteps') }}</th>
            <th>{{ t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in workflows" :key="w.id">
            <td><code>{{ w.key }}</code></td>
            <td>{{ w.name || '-' }}</td>
            <td>
              <n-tag size="small">{{ t('workflow.stepsCount', { count: w.steps.length }) }}</n-tag>
            </td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="edit(w)">{{ t('common.editRun') }}</n-button>
                <n-button size="tiny" type="error" ghost @click="remove(w)">
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
