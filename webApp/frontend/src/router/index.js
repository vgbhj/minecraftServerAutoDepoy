import { createRouter, createWebHistory } from 'vue-router'
import NotFoundView from '@/views/NotFoundView.vue'
import ServerView from '@/views/ServerView.vue'
import SettingsView from '@/views/SettingsView.vue'
import ConsoleView from '@/views/ConsoleView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'server',
      component: ServerView,
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
    },
    {
      path: '/console',
      name: 'console',
      component: ConsoleView,
    },
    {
        path: '/:cathAll(.*)',
        name: 'not-fount',
        component: NotFoundView,
    },
  ],
})

export default router
