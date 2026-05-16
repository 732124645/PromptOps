import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api, type Workspace } from '../api/client'

const KEY = 'promptops_workspace'

export const useWorkspaceStore = defineStore('workspace', () => {
  const list = ref<Workspace[]>([])
  const currentId = ref<string>(localStorage.getItem(KEY) || 'default')

  async function loadList() {
    try {
      const { data } = await api.listWorkspaces()
      list.value = data.data
      if (list.value.length && !list.value.some((w) => w.id === currentId.value)) {
        setCurrent(list.value[0].id)
      }
    } catch {
      /* ignore — header will show the default workspace */
    }
  }

  function setCurrent(id: string) {
    currentId.value = id
    localStorage.setItem(KEY, id)
  }

  return { list, currentId, loadList, setCurrent }
})
