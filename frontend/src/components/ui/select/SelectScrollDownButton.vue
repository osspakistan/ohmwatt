<script setup lang="ts">
import type { SelectScrollDownButtonProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { ArrowDown01Icon } from '@hugeicons/core-free-icons'
import { reactiveOmit } from '@vueuse/core'
import { cn } from 'cn'
import { SelectScrollDownButton, useForwardProps } from 'reka-ui'

const props = defineProps<SelectScrollDownButtonProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = reactiveOmit(props, 'class')

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <SelectScrollDownButton
    data-slot="select-scroll-down-button"
    v-bind="forwardedProps"
    :class="cn('flex cursor-default items-center justify-center py-1', props.class)"
  >
    <slot>
      <ArrowDown01Icon class="size-4" />
    </slot>
  </SelectScrollDownButton>
</template>
