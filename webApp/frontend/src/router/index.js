import { createRouter, createWebHistory } from 'vue-router'
import NotFoundView from '@/views/NotFoundView.vue'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'hello',
      component: HomeView,
    },
    {
        path: '/:cathAll(.*)',
        name: 'not-fount',
        component: NotFoundView,
    },
  ],
})

export default router
