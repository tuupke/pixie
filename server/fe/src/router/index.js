import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
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
    },
    {
      path: '/map',
      name: 'map',
      component: () => import('../views/MapDrawer.vue')
    },
    {
      path: '/playground',
      name: 'playground',
      component: () => import('../views/Playground.vue')
    }
  ]
})

export default router
