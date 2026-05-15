import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('../views/LoginView.vue') },
    { path: '/', name: 'prompts', component: () => import('../views/PromptListView.vue') },
    { path: '/prompts/new', name: 'prompt-new', component: () => import('../views/PromptEditView.vue') },
    { path: '/prompts/:id', name: 'prompt-edit', component: () => import('../views/PromptEditView.vue') },
    { path: '/playground', name: 'playground', component: () => import('../views/PlaygroundView.vue') },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.name !== 'login' && !auth.token) return { name: 'login' }
  if (to.name === 'login' && auth.token) return { name: 'prompts' }
  return true
})

export default router
