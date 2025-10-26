import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue'),
    },
    {
      path: '/metrics',
      name: 'metrics',
      component: () => import('../views/Metrics.vue'),
    },
    {
      path: '/infrastructure',
      name: 'infrastructure',
      component: () => import('../views/Infrastructure.vue'),
    },
    {
      path: '/security',
      name: 'security',
      component: () => import('../views/Security.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/Settings.vue'),
    },
  ],
})

export default router
