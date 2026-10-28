import { createRouter, createWebHashHistory } from 'vue-router'

import AppLayout from '@/layouts/app-layout.vue'
import DashboardView from '@/features/dashboard/dashboard-view.vue'
import HealthView from '@/features/health/health-view.vue'
import HistoryView from '@/features/history/history-view.vue'
import DeviceView from '@/features/device/device-view.vue'
import AlertsView from '@/features/alerts/alerts-view.vue'
import HibernateView from '@/features/hibernate/hibernate-view.vue'
import ReportsView from '@/features/reports/reports-view.vue'
import SettingsView from '@/features/settings/settings-view.vue'

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: AppLayout,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: DashboardView,
          meta: { index: 0 },
        },
        {
          path: 'health',
          name: 'health',
          component: HealthView,
          meta: { index: 1 },
        },
        {
          path: 'history',
          name: 'history',
          component: HistoryView,
          meta: { index: 2 },
        },
        {
          path: 'device',
          name: 'device',
          component: DeviceView,
          meta: { index: 3 },
        },
        {
          path: 'alerts',
          name: 'alerts',
          component: AlertsView,
          meta: { index: 4 },
        },
        {
          path: 'hibernate',
          name: 'hibernate',
          component: HibernateView,
          meta: { index: 5 },
        },
        {
          path: 'reports',
          name: 'reports',
          component: ReportsView,
          meta: { index: 6 },
        },
        {
          path: 'settings',
          name: 'settings',
          component: SettingsView,
          meta: { index: 7 },
        },
      ],
    },
    // Fallback: redirect any unknown paths to /dashboard
    {
      path: '/:pathMatch(.*)*',
      redirect: '/dashboard',
    },
  ],
})
