<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

const users = ref<User[]>([])
const loading = ref(false)
const showCreate = ref(false)
const form = ref({ username: '', password: '', role: 'viewer' })

const roleOptions = [
  { label: 'admin — 全部权限', value: 'admin' },
  { label: 'editor — 可读写', value: 'editor' },
  { label: 'viewer — 只读', value: 'viewer' },
]

async function load() {
  loading.value = true
  try {
    const { data } = await api.listUsers()
    users.value = data.data
  } catch {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

async function create() {
  if (!form.value.username || !form.value.password) {
    message.warning('用户名和密码必填')
    return
  }
  try {
    await api.createUser(form.value)
    message.success('已创建')
    showCreate.value = false
    form.value = { username: '', password: '', role: 'viewer' }
    load()
  } catch (e) {
    const resp = (e as { response?: { data?: { error?: string } } }).response
    message.error(resp?.data?.error || '创建失败')
  }
}

async function changeRole(u: User, role: string) {
  await api.updateUser(u.id, { role })
  message.success(`已将 ${u.username} 设为 ${role}`)
  load()
}

function remove(u: User) {
  dialog.warning({
    title: '删除用户',
    content: `确定删除用户 "${u.username}" 吗?`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.removeUser(u.id)
      message.success('已删除')
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
      <h2>用户管理</h2>
      <n-button type="primary" @click="showCreate = true">+ 新建用户</n-button>
    </div>

    <n-spin :show="loading">
      <n-empty v-if="!users.length" description="暂无用户" style="margin: 48px 0" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>用户名</th>
            <th>角色</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>
              <code>{{ u.username }}</code>
              <n-tag v-if="u.username === auth.username" size="tiny" type="success">当前</n-tag>
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
                删除
              </n-button>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-spin>

    <n-modal v-model:show="showCreate" preset="card" title="新建用户" style="width: 420px">
      <n-form>
        <n-form-item label="用户名">
          <n-input v-model:value="form.username" placeholder="用户名" />
        </n-form-item>
        <n-form-item label="密码">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            placeholder="密码"
          />
        </n-form-item>
        <n-form-item label="角色">
          <n-select v-model:value="form.role" :options="roleOptions" />
        </n-form-item>
        <n-space justify="end">
          <n-button @click="showCreate = false">取消</n-button>
          <n-button type="primary" @click="create">创建</n-button>
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
