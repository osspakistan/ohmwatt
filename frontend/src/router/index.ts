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
        },
        {
          path: 'health',
          name: 'health',
          component: HealthView,
        },
        {
          path: 'history',
          name: 'history',
          component: HistoryView,
        },
        {
          path: 'device',
          name: 'device',
          component: DeviceView,
        },
        {
          path: 'alerts',
          name: 'alerts',
          component: AlertsView,
        },
        {
          path: 'hibernate',
          name: 'hibernate',
          component: HibernateView,
        },
        {
          path: 'reports',
          name: 'reports',
          component: ReportsView,
        },
        {
          path: 'settings',
          name: 'settings',
          component: SettingsView,
        },
      ],
    },
  ],
})
