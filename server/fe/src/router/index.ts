import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/settings/map',
      name: 'settings',
      component: () => import('../views/MapEditor.vue')
    },
    {
      path: '/settings/map/:roomId',
      name: 'room-editor',
      component: () => import('../views/RoomEditor.vue'),
      props: route => ({selectedRoomIndex: Number.parseInt(route.params.roomId)})
    },
    {
      path: '/settings/area',
      name: 'area designer',
      component: () => import('../views/AreaDesigner.vue')
    }
  ]
})

export default router
