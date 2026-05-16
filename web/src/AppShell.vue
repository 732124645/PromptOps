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
        <a
          class="ext-link"
          href="https://github.com/732124645/PromptOps"
          target="_blank"
          rel="noopener noreferrer"
          title="GitHub"
        >
          <svg viewBox="0 0 24 24" width="17" height="17" aria-hidden="true">
            <path
              fill="currentColor"
              d="M12 .5C5.37.5 0 5.78 0 12.29c0 5.2 3.44 9.6 8.21 11.16.6.11.82-.25.82-.56 0-.28-.01-1.02-.02-2-3.34.71-4.04-1.58-4.04-1.58-.55-1.37-1.34-1.73-1.34-1.73-1.09-.73.08-.72.08-.72 1.2.08 1.84 1.21 1.84 1.21 1.07 1.8 2.81 1.28 3.5.98.11-.76.42-1.28.76-1.58-2.67-.3-5.47-1.31-5.47-5.83 0-1.29.47-2.34 1.24-3.17-.13-.3-.54-1.52.12-3.16 0 0 1.01-.32 3.3 1.21a11.6 11.6 0 0 1 6 0c2.29-1.53 3.3-1.21 3.3-1.21.66 1.64.25 2.86.12 3.16.77.83 1.23 1.88 1.23 3.17 0 4.53-2.8 5.52-5.48 5.82.43.36.81 1.08.81 2.18 0 1.58-.01 2.85-.01 3.24 0 .31.21.68.83.56A12.04 12.04 0 0 0 24 12.29C24 5.78 18.63.5 12 .5Z"
            />
          </svg>
        </a>
        <a
          class="ext-link"
          href="https://732124645.github.io/PromptOps/"
          target="_blank"
          rel="noopener noreferrer"
          title="Docs"
        >
          <svg viewBox="0 0 24 24" width="17" height="17" aria-hidden="true">
            <path
              fill="currentColor"
              d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6Zm0 2 4 4h-4V4ZM8 13h8v1.6H8V13Zm0 3.4h8V18H8v-1.6Z"
            />
          </svg>
        </a>
        <span class="divider" />
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
  height: 60px;
  padding: 0 26px;
  background: rgba(7, 11, 22, 0.8);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(34, 211, 238, 0.16);
  box-shadow: 0 1px 24px -6px rgba(34, 211, 238, 0.3);
}
.left {
  display: flex;
  align-items: center;
  gap: 26px;
}
.brand {
  display: flex;
  align-items: center;
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 700;
  font-size: 19px;
  letter-spacing: 0.02em;
  color: #ffffff;
  cursor: pointer;
  text-shadow: 0 0 20px rgba(34, 211, 238, 0.45);
}
.nav {
  display: flex;
  gap: 4px;
}
.nav a {
  position: relative;
  cursor: pointer;
  font-size: 13.5px;
  font-weight: 500;
  color: #8493a8;
  padding: 7px 12px;
  border-radius: 7px;
  transition:
    color 0.16s ease,
    background 0.16s ease;
}
.nav a:hover {
  color: #d8edf2;
  background: rgba(34, 211, 238, 0.07);
}
.nav a.active {
  color: #22d3ee;
}
.nav a.active::after {
  content: '';
  position: absolute;
  left: 12px;
  right: 12px;
  bottom: -1px;
  height: 2px;
  border-radius: 2px;
  background: #22d3ee;
  box-shadow: 0 0 10px rgba(34, 211, 238, 0.95);
}
.right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.ext-link {
  display: flex;
  align-items: center;
  color: #8493a8;
  cursor: pointer;
  transition:
    color 0.16s ease,
    filter 0.16s ease;
}
.ext-link:hover {
  color: #22d3ee;
  filter: drop-shadow(0 0 6px rgba(34, 211, 238, 0.7));
}
.divider {
  width: 1px;
  height: 18px;
  margin: 0 2px;
  background: rgba(120, 160, 200, 0.22);
}
.who {
  font-size: 13px;
  color: #93a3b8;
}
.tag {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.11em;
  text-transform: uppercase;
  margin-left: 9px;
  padding: 3px 7px;
  border-radius: 5px;
  background: rgba(34, 211, 238, 0.12);
  border: 1px solid rgba(34, 211, 238, 0.32);
  color: #22d3ee;
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
