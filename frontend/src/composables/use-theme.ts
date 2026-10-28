import { useColorMode } from '@vueuse/core'
import { computed } from 'vue'

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
