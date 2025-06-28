import { createRouter, createWebHistory } from 'vue-router'
import NotFoundView from '@/views/NotFoundView.vue'
import ServerView from '@/views/ServerView.vue'
import SettingsVue from '@/views/SettingsVue.vue'

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
      component: SettingsVue,
    },
    {
        path: '/:cathAll(.*)',
        name: 'not-fount',
        component: NotFoundView,
    },
  ],
})

export default router
