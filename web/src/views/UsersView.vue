<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard,
  NTable,
  NButton,
  NInput,
  NSelect,
  NSpace,
  NTag,
  NEmpty,
  NSpin,
  NModal,
  NForm,
  NFormItem,
  useMessage,
  useDialog,
} from 'naive-ui'
import { api, type User } from '../api/client'
import { useAuthStore } from '../stores/auth'

const message = useMessage()
const dialog = useDialog()
const auth = useAuthStore()
const { t } = useI18n()

const users = ref<User[]>([])
const loading = ref(false)
const showCreate = ref(false)
const form = ref({ username: '', password: '', role: 'viewer' })

const roleOptions = computed(() => [
  { label: t('users.roleAdmin'), value: 'admin' },
  { label: t('users.roleEditor'), value: 'editor' },
  { label: t('users.roleViewer'), value: 'viewer' },
])

async function load() {
  loading.value = true
  try {
    const { data } = await api.listUsers()
    users.value = data.data
  } catch {
    message.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function create() {
  if (!form.value.username || !form.value.password) {
    message.warning(t('users.credentialsRequired'))
    return
  }
  try {
    await api.createUser(form.value)
    message.success(t('common.created'))
    showCreate.value = false
    form.value = { username: '', password: '', role: 'viewer' }
    load()
  } catch (e) {
    const resp = (e as { response?: { data?: { error?: string } } }).response
    message.error(resp?.data?.error || t('common.createFailed'))
  }
}

async function changeRole(u: User, role: string) {
  await api.updateUser(u.id, { role })
  message.success(t('users.roleChanged', { name: u.username, role }))
  load()
}

function remove(u: User) {
  dialog.warning({
    title: t('users.deleteTitle'),
    content: t('users.deleteConfirm', { name: u.username }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      await api.removeUser(u.id)
      message.success(t('common.deleted'))
      load()
    },
  })
}

function fmtTime(s: string) {
  return new Date(s).toLocaleString()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ t('users.title') }}</h2>
      <n-button type="primary" @click="showCreate = true">{{ t('users.newUser') }}</n-button>
    </div>

    <n-spin :show="loading">
      <n-empty v-if="!users.length" :description="t('users.empty')" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>{{ t('users.colUsername') }}</th>
            <th>{{ t('users.colRole') }}</th>
            <th>{{ t('users.colCreatedAt') }}</th>
            <th>{{ t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>
              <code>{{ u.username }}</code>
              <n-tag v-if="u.username === auth.username" size="tiny" type="success">
                {{ t('common.current') }}
              </n-tag>
            </td>
            <td>
              <n-select
                :value="u.role"
                :options="roleOptions"
                size="small"
                style="width: 200px"
                @update:value="(v) => changeRole(u, v)"
              />
            </td>
            <td>{{ fmtTime(u.created_at) }}</td>
            <td>
              <n-button
                size="tiny"
                type="error"
                ghost
                :disabled="u.username === auth.username"
                @click="remove(u)"
              >
                {{ t('common.delete') }}
              </n-button>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-spin>

    <n-modal
      v-model:show="showCreate"
      preset="card"
      :title="t('users.createTitle')"
      style="width: 420px"
    >
      <n-form>
        <n-form-item :label="t('users.colUsername')">
          <n-input v-model:value="form.username" :placeholder="t('users.usernamePlaceholder')" />
        </n-form-item>
        <n-form-item :label="t('login.password')">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            :placeholder="t('users.passwordPlaceholder')"
          />
        </n-form-item>
        <n-form-item :label="t('users.roleLabel')">
          <n-select v-model:value="form.role" :options="roleOptions" />
        </n-form-item>
        <n-space justify="end">
          <n-button @click="showCreate = false">{{ t('common.cancel') }}</n-button>
          <n-button type="primary" @click="create">{{ t('common.create') }}</n-button>
        </n-space>
      </n-form>
    </n-modal>
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
  margin-right: 6px;
}
</style>
