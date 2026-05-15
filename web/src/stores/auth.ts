import { defineStore } from 'pinia'
import { ref } from 'vue'

const STORAGE_KEY = 'promptops_token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem(STORAGE_KEY) || '')

  function setToken(value: string) {
    token.value = value
    localStorage.setItem(STORAGE_KEY, value)
  }

  function clear() {
    token.value = ''
    localStorage.removeItem(STORAGE_KEY)
  }

  return { token, setToken, clear }
})
