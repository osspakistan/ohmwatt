import { createRouter, createWebHistory } from 'vue-router'

import AppLayout from '@/layouts/app-layout.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/features/landing/landing-view.vue'),
    },
    {
      path: '/app',
      component: AppLayout,
      redirect: '/app/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/features/dashboard/dashboard-view.vue'),
        },
        {
          path: 'health',
          name: 'health',
          component: () => import('@/features/health/health-view.vue'),
        },
        {
          path: 'history',
          name: 'history',
          component: () => import('@/features/history/history-view.vue'),
        },
        {
          path: 'device',
          name: 'device',
          component: () => import('@/features/device/device-view.vue'),
        },
        {
          path: 'alerts',
          name: 'alerts',
          component: () => import('@/features/alerts/alerts-view.vue'),
        },
        {
          path: 'hibernate',
          name: 'hibernate',
          component: () => import('@/features/hibernate/hibernate-view.vue'),
        },
        {
          path: 'reports',
          name: 'reports',
          component: () => import('@/features/reports/reports-view.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/features/settings/settings-view.vue'),
        },
      ],
    },
  ],
})
