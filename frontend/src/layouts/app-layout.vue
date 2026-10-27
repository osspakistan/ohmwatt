<script setup lang="ts">
import {
  DashboardSquare01Icon,
  FavouriteIcon,
  File01Icon,
  LaptopIcon,
  Moon02Icon,
  Notification01Icon,
  Settings01Icon,
  Time02Icon,
} from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AppToolbar from '@/components/common/app-toolbar.vue'
import ThemeToggle from '@/components/common/theme-toggle.vue'

const route = useRoute()
const router = useRouter()

const transitionName = ref('slide-down')

router.beforeEach((to, from) => {
  const toIndex = (to.meta.index as number) ?? 0
  const fromIndex = (from.meta.index as number) ?? 0
  transitionName.value = toIndex > fromIndex ? 'slide-down' : 'slide-up'
})

const navItems = [
  { path: '/dashboard', label: 'Now', icon: DashboardSquare01Icon },
  { path: '/health', label: 'Health', icon: FavouriteIcon },
  { path: '/history', label: 'History', icon: Time02Icon },
  { path: '/device', label: 'Device', icon: LaptopIcon },
  { path: '/alerts', label: 'Alerts', icon: Notification01Icon },
  { path: '/hibernate', label: 'Hibernate', icon: Moon02Icon },
  { path: '/reports', label: 'Reports', icon: File01Icon },
  { path: '/settings', label: 'Settings', icon: Settings01Icon },
]

function isActive(path: string) {
  return route.path.startsWith(path)
}
</script>

<template>
  <div
    class="flex h-screen w-screen flex-col overflow-hidden bg-background text-foreground font-sans antialiased"
  >
    <!-- Top Window Chrome / Draggable Toolbar -->
    <AppToolbar title="Antlion" />

    <!-- Application Body: Slim Left Rail + Main Screen Content -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Slim Navigation Rail -->
      <aside
        class="flex w-16 flex-none flex-col items-center justify-between border-r border-border bg-card/60 py-4 select-none"
      >
        <!-- Top Navigation Icons -->
        <div class="flex flex-col items-center gap-2">
          <button
            v-for="item in navItems"
            :key="item.path"
            type="button"
            :class="[
              'group flex h-11 w-12 flex-col items-center justify-center rounded-lg text-[10px] transition-all duration-150 cursor-pointer active:scale-95',
              isActive(item.path)
                ? 'bg-accent font-medium text-foreground shadow-sm ring-1 ring-border/50'
                : 'text-muted-foreground hover:bg-muted/70 hover:text-foreground',
            ]"
            :title="item.label"
            :aria-label="item.label"
            @click="router.push(item.path)"
          >
            <HugeiconsIcon :icon="item.icon" class="h-4 w-4" />
            <span class="mt-1 leading-none tracking-tight">{{ item.label }}</span>
          </button>
        </div>

        <!-- Bottom Actions (Theme Toggle) -->
        <div class="flex flex-col items-center gap-3">
          <ThemeToggle />
        </div>
      </aside>

      <!-- Active Routed Screen View Host with Directional Desktop Transition -->
      <main class="relative flex-1 overflow-y-auto overflow-x-hidden">
        <router-view v-slot="{ Component }">
          <transition :name="transitionName" mode="out-in">
            <keep-alive>
              <component :is="Component" :key="route.path" />
            </keep-alive>
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Professional Desktop Micro-Transitions (Apple/macOS Fluid Motion Style) */

/* Navigating downward in the tab rail (toIndex > fromIndex) */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: opacity 0.16s cubic-bezier(0.2, 0, 0, 1), transform 0.16s cubic-bezier(0.2, 0, 0, 1);
  will-change: opacity, transform;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* Navigating upward in the tab rail (toIndex < fromIndex) */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: opacity 0.16s cubic-bezier(0.2, 0, 0, 1), transform 0.16s cubic-bezier(0.2, 0, 0, 1);
  will-change: opacity, transform;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(-6px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(6px);
}
</style>
