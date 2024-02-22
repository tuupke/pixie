import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: import('../views/HomeView.vue')
    },
    {
      path: '/settings/map',
      name: 'settings',
      component: () => import('../views/SettingsView.vue')
    },
    {
      path: '/settings/area',
      name: 'area designer',
      component: () => import('../views/AreaDesigner.vue')
    }
  ]
})

export default router
