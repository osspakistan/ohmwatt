<script setup lang="ts">
import { ref } from 'vue'

import CommandBlock from '@/components/common/command-block.vue'
import SegmentControl from '@/components/common/segment-control.vue'

const activeOS = ref<'linux' | 'windows' | 'macos'>('linux')

const osOptions = [
  { value: 'linux', label: 'Ubuntu / Debian' },
  { value: 'windows', label: 'Windows' },
  { value: 'macos', label: 'macOS' },
]

// Diagnostic checklist for user's EliteBook
const checklist = [
  {
    label: 'Swap space size (19.8 GiB) vs RAM (22 GiB)',
    status: 'warning',
    text: 'Recommend extending to 24 GiB',
  },
  { label: 'Physical swap file offset', status: 'ready', text: 'Offset: 34816 (/dev/nvme0n1p3)' },
  {
    label: 'GRUB kernel resume parameter',
    status: 'missing',
    text: 'Missing resume=UUID=... in /etc/default/grub',
  },
  { label: 'Secure Boot Lockdown', status: 'ready', text: 'Disabled (Will not block hibernation)' },
]
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
          <p class="text-xs text-muted-foreground mt-1">
            Editing bootloader parameters and swap images alters core system initialization. Always
            ensure critical files are backed up before applying system modifications.
          </p>
        </div>
      </div>
    </div>

    <!-- Diagnostic State -->
    <div v-if="activeOS === 'linux'" class="space-y-6">
      <div class="rounded-xl border border-border bg-card p-5">
        <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-4">
          Diagnostic Analysis: HP EliteBook 850 G7
        </h2>

        <div class="space-y-3">
          <div
            v-for="(item, idx) in checklist"
            :key="idx"
            class="flex items-center justify-between border-b border-border py-2 last:border-0"
          >
            <div>
              <div class="text-sm font-medium text-foreground">{{ item.label }}</div>
              <div class="text-xs text-muted-foreground">{{ item.text }}</div>
            </div>
            <span
              :class="[
                'text-xs font-mono uppercase font-semibold px-2 py-0.5 rounded',
                item.status === 'ready'
                  ? 'bg-secondary/20 text-secondary-foreground'
                  : item.status === 'warning'
                    ? 'bg-destructive/10 text-destructive'
                    : 'bg-muted text-muted-foreground',
              ]"
            >
              {{ item.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Step-by-Step Tailored Recipe Blocks -->
      <div class="space-y-4">
        <h3 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          Configuration Steps (Copy & Paste to Terminal)
        </h3>

        <CommandBlock
          label="1. Add Resume Parameters to GRUB (/etc/default/grub)"
          command='GRUB_CMDLINE_LINUX_DEFAULT="quiet splash resume=UUID=88bbefa2-dd83-4049-8ede-c17fbc61898c resume_offset=34816"'
        />

        <CommandBlock
          label="2. Rebuild GRUB & Initramfs"
          command="sudo update-grub && sudo update-initramfs -u -k all"
        />

        <CommandBlock label="3. Test System Hibernation" command="sudo systemctl hibernate" />
      </div>
    </div>

    <!-- Windows Tab -->
    <div v-else-if="activeOS === 'windows'" class="space-y-4">
      <div class="rounded-xl border border-border bg-card p-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-2">
          Windows Power Configuration
        </h2>
        <p class="text-xs text-muted-foreground mb-4">
          Windows handles hibernation through the standard Windows power management utility. Open an
          elevated Command Prompt or PowerShell and run:
        </p>
        <CommandBlock
          command="powercfg /hibernate on"
          label="Enable Hibernation (Admin PowerShell)"
        />
      </div>
    </div>

    <!-- macOS Tab -->
    <div v-else class="space-y-4">
      <div class="rounded-xl border border-border bg-card p-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-2">
          macOS Standby Modes
        </h2>
        <p class="text-xs text-muted-foreground mb-4">
          macOS regulates sleep and hibernation via
          <code class="font-mono text-xs">pmset hibernatemode</code>. Inspect your active power
          management profile:
        </p>
        <CommandBlock command="pmset -g | grep hibernatemode" label="Inspect macOS Hibernatemode" />
      </div>
    </div>
  </div>
</template>
