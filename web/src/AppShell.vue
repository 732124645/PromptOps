<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NButton,
  NTag,
  NSelect,
  NModal,
  NInput,
  useMessage,
  useDialog,
} from 'naive-ui'
import { useAuthStore } from './stores/auth'
import { useWorkspaceStore } from './stores/workspace'
import { api } from './api/client'
import { setLocale, type Locale } from './i18n'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const ws = useWorkspaceStore()
const message = useMessage()
const dialog = useDialog()
const { t, locale } = useI18n()

const localeOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en' },
]
function changeLocale(value: Locale) {
  setLocale(value)
}

const showHeader = computed(() => route.name !== 'login')
const wsOptions = computed(() => ws.list.map((w) => ({ label: w.name, value: w.id })))

const showWs = ref(false)
const newWsName = ref('')

function isActive(
  section: 'prompts' | 'playground' | 'agents' | 'workflows' | 'observability' | 'users',
) {
  const name = String(route.name || '')
  const groups: Record<string, string[]> = {
    prompts: ['prompts', 'prompt-new', 'prompt-edit'],
    playground: ['playground'],
    agents: ['agents', 'agent-new', 'agent-edit'],
    workflows: ['workflows', 'workflow-new', 'workflow-edit'],
    observability: ['observability'],
    users: ['users'],
  }
  return groups[section].includes(name)
}

async function createWs() {
  if (!newWsName.value.trim()) {
    return
  }
  try {
    const { data } = await api.createWorkspace({ name: newWsName.value.trim() })
    newWsName.value = ''
    await ws.loadList()
    ws.setCurrent(data.data.id)
    message.success(t('app.workspaceCreated'))
  } catch {
    message.error(t('app.workspaceCreateFailed'))
  }
}

function removeWs(id: string, name: string) {
  dialog.warning({
    title: t('app.deleteWorkspaceTitle'),
    content: t('app.deleteWorkspaceConfirm', { name }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await api.removeWorkspace(id)
        await ws.loadList()
        message.success(t('common.deleted'))
      } catch {
        message.error(t('common.deleteFailed'))
      }
    },
  })
}

async function logout() {
  await api.logout().catch(() => {})
  auth.clear()
  router.push({ name: 'login' })
}

// Load workspaces once authenticated (and reload on login).
watch(
  () => auth.token,
  (token) => {
    if (token) ws.loadList()
  },
  { immediate: true },
)
</script>

<template>
  <n-layout style="height: 100vh">
    <n-layout-header v-if="showHeader" bordered class="header">
      <div class="left">
        <div class="brand" @click="router.push('/')">
          PromptOps
          <span class="tag">Runtime</span>
        </div>
        <nav class="nav">
          <a :class="{ active: isActive('prompts') }" @click="router.push('/')">
            {{ t('nav.prompts') }}
          </a>
          <a :class="{ active: isActive('agents') }" @click="router.push('/agents')">
            {{ t('nav.agents') }}
          </a>
          <a :class="{ active: isActive('workflows') }" @click="router.push('/workflows')">
            {{ t('nav.workflows') }}
          </a>
          <a :class="{ active: isActive('playground') }" @click="router.push('/playground')">
            {{ t('nav.playground') }}
          </a>
          <a
            :class="{ active: isActive('observability') }"
            @click="router.push('/observability')"
          >
            {{ t('nav.observability') }}
          </a>
          <a
            v-if="auth.isAdmin"
            :class="{ active: isActive('users') }"
            @click="router.push('/users')"
          >
            {{ t('nav.users') }}
          </a>
        </nav>
      </div>
      <div class="right">
        <n-select
          :value="locale"
          :options="localeOptions"
          size="small"
          style="width: 96px"
          @update:value="changeLocale"
        />
        <n-select
          :value="ws.currentId"
          :options="wsOptions"
          size="small"
          style="width: 150px"
          :placeholder="t('app.workspacePlaceholder')"
          @update:value="ws.setCurrent"
        />
        <n-button quaternary size="tiny" @click="showWs = true">{{ t('app.manage') }}</n-button>
        <span class="who">{{ auth.username || '—' }}</span>
        <n-tag size="tiny" :type="auth.isAdmin ? 'success' : 'default'">{{ auth.role }}</n-tag>
        <n-button quaternary size="small" @click="logout">{{ t('app.logout') }}</n-button>
      </div>
    </n-layout-header>
    <n-layout-content :content-style="'min-height: 100%'">
      <router-view />
    </n-layout-content>

    <n-modal
      v-model:show="showWs"
      preset="card"
      :title="t('app.workspaceManage')"
      style="width: 440px"
    >
      <div v-for="w in ws.list" :key="w.id" class="ws-row">
        <span>{{ w.name }}</span>
        <n-tag v-if="w.id === 'default'" size="tiny">{{ t('app.default') }}</n-tag>
        <n-button v-else size="tiny" type="error" ghost @click="removeWs(w.id, w.name)">
          {{ t('common.delete') }}
        </n-button>
      </div>
      <div class="ws-create">
        <n-input
          v-model:value="newWsName"
          :placeholder="t('app.newWorkspacePlaceholder')"
          size="small"
        />
        <n-button size="small" type="primary" @click="createWs">{{ t('common.create') }}</n-button>
      </div>
    </n-modal>
  </n-layout>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 24px;
}
.left {
  display: flex;
  align-items: center;
  gap: 28px;
}
.brand {
  font-weight: 700;
  font-size: 18px;
  cursor: pointer;
}
.nav {
  display: flex;
  gap: 18px;
}
.nav a {
  cursor: pointer;
  font-size: 14px;
  color: #999;
  transition: color 0.15s;
}
.nav a:hover {
  color: #fff;
}
.nav a.active {
  color: #63e2b7;
}
.right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.who {
  font-size: 13px;
  color: #aaa;
}
.tag {
  font-size: 11px;
  font-weight: 500;
  margin-left: 6px;
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(99, 226, 183, 0.15);
  color: #63e2b7;
}
.ws-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
}
.ws-create {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}
</style>
