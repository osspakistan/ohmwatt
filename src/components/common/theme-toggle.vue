<script setup lang="ts">
import { Moon02Icon, Sun01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'

import { Button } from '@/components/ui/button'
import { useTheme } from '@/composables/use-theme'

const { isDark, toggleTheme } = useTheme()
</script>

<template>
  <Button
    variant="ghost"
    size="icon"
    class="relative h-9 w-9 rounded-md border border-border text-foreground transition-colors hover:bg-muted"
    :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
    aria-label="Toggle Theme"
    @click="toggleTheme"
  >
    <HugeiconsIcon
      v-if="isDark"
      :icon="Sun01Icon"
      class="h-4 w-4 text-primary transition-all duration-300"
    />
    <HugeiconsIcon
      v-else
      :icon="Moon02Icon"
      class="h-4 w-4 text-foreground transition-all duration-300"
    />
    <span class="sr-only">Toggle theme</span>
  </Button>
</template>
