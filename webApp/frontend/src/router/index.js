import { createRouter, createWebHistory } from 'vue-router'
import HelloView from '@/views/HelloView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'hello',
      component: HelloView,
    },
    {
        path: '/:cathAll(.*)',
        name: 'not-fount',
        component: NotFoundView,
    },
  ],
})

export default router
