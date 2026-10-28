<script setup lang="ts">
import { ref } from 'vue'

import MaskedField from '@/components/common/masked-field.vue'
import SegmentControl from '@/components/common/segment-control.vue'

const activeTab = ref<'battery' | 'adapter' | 'laptop'>('battery')

const tabOptions = [
  { value: 'battery', label: 'Battery' },
  { value: 'adapter', label: 'Adapter' },
  { value: 'laptop', label: 'Laptop' },
]
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">Device</h1>
        <p class="text-muted-foreground mt-1.5">
          Hardware specifications and connected power supplies.
        </p>
      </div>
      <SegmentControl v-model="activeTab" :options="tabOptions" />
    </header>

    <!-- Battery Specs Sub-tab (F27) -->
    <div
      v-if="activeTab === 'battery'"
      class="rounded-xl border border-border bg-card p-6 divide-y divide-border"
    >
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Manufacturer</span>
        <span class="font-mono text-sm font-medium">SMP (Simplo Technology)</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Model Name</span>
        <span class="font-mono text-sm font-medium">L38125-005</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Serial Number</span>
        <MaskedField value="SMP-492184-HP" />
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Chemistry</span>
        <span class="font-mono text-sm font-medium">Lithium-Polymer (Li-poly)</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Design Capacity</span>
        <span class="font-mono text-sm font-medium">53,010 mWh</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Full Charge Capacity</span>
        <span class="font-mono text-sm font-medium">48,240 mWh</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Cycle Count</span>
        <span class="font-mono text-sm font-medium">212</span>
      </div>
    </div>

    <!-- Adapter Sub-tab (F29) -->
    <div
      v-else-if="activeTab === 'adapter'"
      class="rounded-xl border border-border bg-card p-6 divide-y divide-border"
    >
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Connection State</span>
        <span
          class="inline-flex items-center gap-1.5 rounded-full bg-secondary/30 text-secondary-foreground border border-secondary/50 px-2.5 py-0.5 text-xs font-medium"
        >
          Connected (Online)
        </span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Protocol</span>
        <span class="font-mono text-sm font-medium">USB Power Delivery 3.0</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Rated Wattage</span>
        <span class="font-mono text-sm font-medium">65.0 W (20V @ 3.25A)</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Current Power Output</span>
        <span class="font-mono text-sm font-medium">32.4 W</span>
      </div>
      <div class="pt-4 text-xs text-destructive flex items-center gap-2">
        <span>⚠</span>
        <span>Laptop requires 90 W for peak performance under heavy GPU/CPU workloads.</span>
      </div>
    </div>

    <!-- Laptop Sub-tab (F30) -->
    <div v-else class="rounded-xl border border-border bg-card p-6 divide-y divide-border">
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">OEM & Model</span>
        <span class="font-mono text-sm font-medium">HP EliteBook 850 G7 Notebook PC</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Service Tag / Serial</span>
        <MaskedField value="5CG1392XYZ" />
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">BIOS Version</span>
        <span class="font-mono text-sm font-medium">HP S70 Ver. 01.12.01</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">OS & Kernel</span>
        <span class="font-mono text-sm font-medium"
          >Ubuntu 26.04 LTS (Kernel 7.0.0-38-generic)</span
        >
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">System RAM</span>
        <span class="font-mono text-sm font-medium">22 GiB</span>
      </div>
    </div>
  </div>
</template>
