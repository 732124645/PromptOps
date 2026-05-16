import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const TOKEN_KEY = 'promptops_token'
const ROLE_KEY = 'promptops_role'
const USER_KEY = 'promptops_user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const role = ref<string>(localStorage.getItem(ROLE_KEY) || '')
  const username = ref<string>(localStorage.getItem(USER_KEY) || '')

  const isAdmin = computed(() => role.value === 'admin')
  const canWrite = computed(() => role.value === 'admin' || role.value === 'editor')

  function setSession(t: string, r: string, u: string) {
    token.value = t
    role.value = r
    username.value = u
    localStorage.setItem(TOKEN_KEY, t)
    localStorage.setItem(ROLE_KEY, r)
    localStorage.setItem(USER_KEY, u)
  }

  function clear() {
    token.value = ''
    role.value = ''
    username.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(ROLE_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, role, username, isAdmin, canWrite, setSession, clear }
})
