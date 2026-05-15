<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NInput, NButton, NForm, NFormItem, useMessage } from 'naive-ui'
import { api } from '../api/client'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const message = useMessage()

const token = ref('')
const loading = ref(false)

async function login() {
  if (!token.value) {
    message.warning('请输入访问 Token')
    return
  }
  loading.value = true
  try {
    await api.login(token.value)
    auth.setToken(token.value)
    router.push({ name: 'prompts' })
  } catch {
    message.error('Token 无效')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrap">
    <n-card class="login-card" :bordered="true">
      <div class="title">PromptOps</div>
      <div class="subtitle">AI Prompt Runtime 平台</div>
      <n-form @submit.prevent="login">
        <n-form-item label="访问 Token">
          <n-input
            v-model:value="token"
            type="password"
            show-password-on="click"
            placeholder="默认: promptops-dev-token"
            @keyup.enter="login"
          />
        </n-form-item>
        <n-button type="primary" block :loading="loading" @click="login">登录</n-button>
      </n-form>
      <div class="hint">默认 Token: <code>promptops-dev-token</code>(可用 PROMPTOPS_TOKEN 配置)</div>
    </n-card>
  </div>
</template>

<style scoped>
.login-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
}
.login-card {
  width: 360px;
}
.title {
  font-size: 24px;
  font-weight: 700;
}
.subtitle {
  color: #888;
  margin: 4px 0 20px;
  font-size: 13px;
}
.hint {
  margin-top: 14px;
  font-size: 12px;
  color: #888;
}
</style>
