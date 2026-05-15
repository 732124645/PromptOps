<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutHeader, NLayoutContent, NButton } from 'naive-ui'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const showHeader = computed(() => route.name !== 'login')

function logout() {
  auth.clear()
  router.push({ name: 'login' })
}
</script>

<template>
  <n-layout style="height: 100vh">
    <n-layout-header v-if="showHeader" bordered class="header">
      <div class="brand" @click="router.push('/')">
        PromptOps
        <span class="tag">Runtime</span>
      </div>
      <n-button quaternary size="small" @click="logout">退出登录</n-button>
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
.brand {
  font-weight: 700;
  font-size: 18px;
  cursor: pointer;
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
