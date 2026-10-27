import { useColorMode } from '@vueuse/core'
import { computed, watch } from 'vue'

export type ThemeMode = 'auto' | 'light' | 'dark'

export function useTheme() {
  const mode = useColorMode({
    attribute: 'class',
    modes: {
      light: '',
      dark: 'dark',
    },
    storageKey: 'antlion-theme',
  })

  const isDark = computed(() => mode.value === 'dark')

  // Dynamically update runtime favicon when user toggles theme in the app
  watch(
    isDark,
    (dark) => {
      if (typeof document === 'undefined') return
      const favicon = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
      if (favicon) {
        favicon.href = dark ? '/antlion-mark-dark.svg' : '/antlion-mark-light.svg'
      }
    },
    { immediate: true },
  )

  function toggleTheme() {
    mode.value = mode.value === 'dark' ? 'light' : 'dark'
  }

  function setTheme(newMode: ThemeMode) {
    mode.value = newMode
  }

  return {
    mode,
    isDark,
    toggleTheme,
    setTheme,
  }
}
