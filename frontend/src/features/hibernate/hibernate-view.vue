<script setup lang="ts">
import { ref, onMounted } from 'vue'

import type { Diagnostic } from '../../../bindings/changeme/core/hibernate'
import { HibernateService } from '../../../bindings/changeme/service'

import CommandBlock from '@/components/common/command-block.vue'
import SegmentControl from '@/components/common/segment-control.vue'

const activeOS = ref<'linux' | 'windows' | 'macos'>('linux')

const osOptions = [
  { value: 'linux', label: 'Ubuntu / Debian' },
  { value: 'windows', label: 'Windows' },
  { value: 'macos', label: 'macOS' },
]

const diag = ref<Diagnostic | null>(null)

onMounted(async () => {
  try {
    const res = await HibernateService.GetDiagnostic()
    if (res) {
      diag.value = res
    }
  } catch (err) {
    console.warn('Failed to load hibernate diagnostic:', err)
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">Hibernate Guide</h1>
        <p class="text-muted-foreground mt-1.5">
          Zero-drain power management setup and system diagnostics.
        </p>
      </div>
      <SegmentControl v-model="activeOS" :options="osOptions" />
    </header>

    <!-- Prominent Risk Warning Banner -->
    <div class="rounded-xl border border-destructive/30 bg-destructive/10 p-5 mb-8">
      <div class="flex items-start gap-3">
        <span class="text-destructive text-lg font-bold">⚠</span>
        <div>
          <div class="text-sm font-semibold text-destructive uppercase tracking-wide">
            System Configuration Notice
          </div>
          <p class="text-sm text-foreground/80 mt-1">
            Hibernation writes active RAM contents directly to disk swap. Follow the verified
            commands below to prevent boot-looping or lost sessions.
          </p>
        </div>
      </div>
    </div>

    <!-- Live System Readiness Checklist (F26) -->
    <div v-if="activeOS === 'linux'" class="mb-10">
      <h2 class="text-xs uppercase font-mono tracking-wider text-muted-foreground mb-4">
        Live System Readiness (Kernel Diagnostic)
      </h2>
      <div class="rounded-xl border border-border bg-card divide-y divide-border overflow-hidden">
        <div
          v-for="item in diag?.checks || []"
          :key="item.title"
          class="p-4 flex items-center justify-between flex-wrap gap-3"
        >
          <div class="space-y-1">
            <div class="text-sm font-medium text-foreground">{{ item.title }}</div>
            <div class="text-xs text-muted-foreground">{{ item.detail }}</div>
          </div>
          <span
            :class="[
              'text-xs font-mono uppercase font-semibold px-2 py-0.5 rounded',
              item.passed
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
            ]"
          >
            {{ item.passed ? 'Ready' : 'Action Needed' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Instructions Tabs -->
    <div v-if="activeOS === 'linux'" class="space-y-6">
      <h2 class="text-xs uppercase font-mono tracking-wider text-muted-foreground">
        Step-by-Step Ubuntu / Debian Commands
      </h2>

      <div class="space-y-3">
        <div class="text-sm font-medium">1. Ensure Sufficient Swap Allocation</div>
        <p class="text-xs text-muted-foreground">
          Allocate a swapfile that comfortably accommodates your {{ diag?.ramSizeGiB ?? 22 }} GiB
          RAM:
        </p>
        <CommandBlock
          command="sudo fallocate -l 24G /swap.img && sudo chmod 600 /swap.img && sudo mkswap /swap.img && sudo swapon /swap.img"
        />
      </div>

      <div class="space-y-3">
        <div class="text-sm font-medium">2. Find Swap Physical Offset</div>
        <p class="text-xs text-muted-foreground">
          Query the first block physical offset for the kernel resume loader:
        </p>
        <CommandBlock command="sudo filefrag -v /swap.img | head -n 4" />
      </div>

      <div class="space-y-3">
        <div class="text-sm font-medium">3. Update GRUB Resume Target</div>
        <p class="text-xs text-muted-foreground">
          Append resume UUID and resume_offset to /etc/default/grub:
        </p>
        <CommandBlock
          command="sudo nano /etc/default/grub && sudo update-grub && sudo update-initramfs -u -k all"
        />
      </div>
    </div>

    <!-- Windows Tab -->
    <div v-else-if="activeOS === 'windows'" class="space-y-6">
      <h2 class="text-xs uppercase font-mono tracking-wider text-muted-foreground">
        Windows Powercfg Setup
      </h2>
      <div class="space-y-3">
        <div class="text-sm font-medium">Enable Full Hibernation File</div>
        <CommandBlock command="powercfg.exe /hibernate on" />
        <CommandBlock command="powercfg.exe /h /type full" />
      </div>
    </div>

    <!-- macOS Tab -->
    <div v-else class="space-y-6">
      <h2 class="text-xs uppercase font-mono tracking-wider text-muted-foreground">
        macOS Standby / SafeSleep
      </h2>
      <div class="space-y-3">
        <div class="text-sm font-medium">Set Hibernatemode 25 (Deep Zero-Drain)</div>
        <CommandBlock command="sudo pmset -a hibernatemode 25" />
      </div>
    </div>
  </div>
</template>
