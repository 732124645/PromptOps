<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutHeader, NLayoutContent, NButton, NTag } from 'naive-ui'
import { useAuthStore } from './stores/auth'
import { api } from './api/client'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const showHeader = computed(() => route.name !== 'login')

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

async function logout() {
  await api.logout().catch(() => {})
  auth.clear()
  router.push({ name: 'login' })
}
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
        <span class="who">{{ auth.username || '—' }}</span>
        <n-tag size="tiny" :type="auth.isAdmin ? 'success' : 'default'">{{ auth.role }}</n-tag>
        <n-button quaternary size="small" @click="logout">退出登录</n-button>
      </div>
    </n-layout-header>
    <n-layout-content :content-style="'min-height: 100%'">
      <router-view />
    </n-layout-content>
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
</style>
