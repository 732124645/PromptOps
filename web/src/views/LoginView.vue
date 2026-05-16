<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NCard, NInput, NButton, NForm, NFormItem, useMessage } from 'naive-ui'
import { api } from '../api/client'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const message = useMessage()
const { t } = useI18n()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function login() {
  if (!username.value || !password.value) {
    message.warning(t('login.enterCredentials'))
    return
  }
  loading.value = true
  try {
    const { data } = await api.login({
      username: username.value,
      password: password.value,
    })
    auth.setSession(data.token, data.role, data.username)
    router.push({ name: 'prompts' })
  } catch {
    message.error(t('login.invalidCredentials'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrap">
    <n-card class="login-card" :bordered="true">
      <div class="title">PromptOps</div>
      <div class="subtitle">{{ t('login.subtitle') }}</div>
      <n-form @submit.prevent="login">
        <n-form-item :label="t('login.username')">
          <n-input
            v-model:value="username"
            :placeholder="t('login.username')"
            @keyup.enter="login"
          />
        </n-form-item>
        <n-form-item :label="t('login.password')">
          <n-input
            v-model:value="password"
            type="password"
            show-password-on="click"
            :placeholder="t('login.password')"
            @keyup.enter="login"
          />
        </n-form-item>
        <n-button type="primary" block :loading="loading" @click="login">
          {{ t('login.signIn') }}
        </n-button>
      </n-form>
      <div class="hint">{{ t('login.hint') }}<code>admin</code> / <code>admin</code></div>
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
