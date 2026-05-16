<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const ws = useWorkspaceStore()
const message = useMessage()
const dialog = useDialog()

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
    message.success('工作区已创建')
  } catch {
    message.error('创建工作区失败')
  }
}

function removeWs(id: string, name: string) {
  dialog.warning({
    title: '删除工作区',
    content: `确定删除工作区 "${name}" 吗?其中的资源不会被删除,但将不再归属此工作区。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await api.removeWorkspace(id)
        await ws.loadList()
        message.success('已删除')
      } catch {
        message.error('删除失败')
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
          <a :class="{ active: isActive('prompts') }" @click="router.push('/')">Prompts</a>
          <a :class="{ active: isActive('agents') }" @click="router.push('/agents')">Agents</a>
          <a :class="{ active: isActive('workflows') }" @click="router.push('/workflows')">
            Workflows
          </a>
          <a :class="{ active: isActive('playground') }" @click="router.push('/playground')">
            Playground
          </a>
          <a
            :class="{ active: isActive('observability') }"
            @click="router.push('/observability')"
          >
            观测
          </a>
          <a
            v-if="auth.isAdmin"
            :class="{ active: isActive('users') }"
            @click="router.push('/users')"
          >
            用户
          </a>
        </nav>
      </div>
      <div class="right">
        <n-select
          :value="ws.currentId"
          :options="wsOptions"
          size="small"
          style="width: 150px"
          placeholder="工作区"
          @update:value="ws.setCurrent"
        />
        <n-button quaternary size="tiny" @click="showWs = true">管理</n-button>
        <span class="who">{{ auth.username || '—' }}</span>
        <n-tag size="tiny" :type="auth.isAdmin ? 'success' : 'default'">{{ auth.role }}</n-tag>
        <n-button quaternary size="small" @click="logout">退出登录</n-button>
      </div>
    </n-layout-header>
    <n-layout-content :content-style="'min-height: 100%'">
      <router-view />
    </n-layout-content>

    <n-modal v-model:show="showWs" preset="card" title="工作区管理" style="width: 440px">
      <div v-for="w in ws.list" :key="w.id" class="ws-row">
        <span>{{ w.name }}</span>
        <n-tag v-if="w.id === 'default'" size="tiny">默认</n-tag>
        <n-button v-else size="tiny" type="error" ghost @click="removeWs(w.id, w.name)">
          删除
        </n-button>
      </div>
      <div class="ws-create">
        <n-input v-model:value="newWsName" placeholder="新工作区名称" size="small" />
        <n-button size="small" type="primary" @click="createWs">创建</n-button>
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
