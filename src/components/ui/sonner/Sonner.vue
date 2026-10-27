<script lang="ts" setup>
import type { ToasterProps } from 'vue-sonner'
import {
  CheckmarkCircle01Icon,
  InformationCircleIcon,
  Loading03Icon,
  CancelCircleIcon,
  Alert02Icon,
  Cancel01Icon,
} from '@hugeicons/core-free-icons'
import { cn } from 'cn'
import { Toaster as Sonner } from 'vue-sonner'

const props = defineProps<ToasterProps>()
</script>

<template>
  <Sonner
    :class="cn('toaster group', props.class)"
    :style="{
      '--normal-bg': 'var(--popover)',
      '--normal-text': 'var(--popover-foreground)',
      '--normal-border': 'var(--border)',
      '--border-radius': 'var(--radius)',
    }"
    v-bind="props"
  >
    <template #success-icon>
      <CheckmarkCircle01Icon class="size-4" />
    </template>
    <template #info-icon>
      <InformationCircleIcon class="size-4" />
    </template>
    <template #warning-icon>
      <Alert02Icon class="size-4" />
    </template>
    <template #error-icon>
      <CancelCircleIcon class="size-4" />
    </template>
    <template #loading-icon>
      <div>
        <Loading03Icon class="size-4 animate-spin" />
      </div>
    </template>
    <template #close-icon>
      <Cancel01Icon class="size-4" />
    </template>
  </Sonner>
</template>
