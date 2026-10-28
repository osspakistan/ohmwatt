<script setup lang="ts">
import { ref, onMounted } from 'vue'

import type { DeviceSpecs } from '../../../bindings/changeme/core/hardware'
import type { AdapterInfo } from '../../../bindings/changeme/core/power'
import { DeviceService } from '../../../bindings/changeme/service'

import MaskedField from '@/components/common/masked-field.vue'
import SegmentControl from '@/components/common/segment-control.vue'

const activeTab = ref<'battery' | 'adapter' | 'laptop'>('battery')

const tabOptions = [
  { value: 'battery', label: 'Battery' },
  { value: 'adapter', label: 'Adapter' },
  { value: 'laptop', label: 'Laptop' },
]

const specs = ref<DeviceSpecs>({
  batteryManufacturer: 'Loading...',
  batteryModel: 'Loading...',
  batterySerial: 'Loading...',
  batteryTechnology: 'Loading...',
  designCapacityWh: '—',
  fullChargeWh: '—',
  designVoltage: '—',
  laptopManufacturer: 'Loading...',
  laptopModel: 'Loading...',
  laptopSerial: 'Loading...',
  biosVersion: 'Loading...',
  totalMemoryRAM: 'Loading...',
  kernelVersion: 'Loading...',
})

const adapter = ref<AdapterInfo>({
  isConnected: false,
  type: 'Detecting...',
  negotiatedWatts: 0,
  currentDrawWatts: 0,
  voltageVolts: 0,
  amperes: 0,
  isUnderpowered: false,
})

onMounted(async () => {
  try {
    const s = await DeviceService.GetDeviceSpecs()
    if (s) {
      specs.value = s
    }
    const a = await DeviceService.GetAdapterInfo()
    if (a) {
      adapter.value = a
    }
  } catch (err) {
    console.warn('Failed to load device specs from Go:', err)
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">Device</h1>
        <p class="text-muted-foreground mt-1.5">
          Real hardware specifications and connected power supplies from Linux sysfs & DMI.
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
        <span class="font-mono text-sm font-medium">{{ specs.batteryManufacturer }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Model Name</span>
        <span class="font-mono text-sm font-medium">{{ specs.batteryModel }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Serial Number</span>
        <MaskedField :value="specs.batterySerial" />
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Chemistry</span>
        <span class="font-mono text-sm font-medium">{{ specs.batteryTechnology }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Design Capacity</span>
        <span class="font-mono text-sm font-medium">{{ specs.designCapacityWh }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Full Charge Capacity</span>
        <span class="font-mono text-sm font-medium">{{ specs.fullChargeWh }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Design Voltage</span>
        <span class="font-mono text-sm font-medium">{{ specs.designVoltage }}</span>
      </div>
    </div>

    <!-- Adapter Sub-tab (F29) -->
    <div
      v-else-if="activeTab === 'adapter'"
      class="rounded-xl border border-border bg-card p-6 divide-y divide-border"
    >
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Connection Status</span>
        <span
          class="font-mono text-sm font-medium"
          :class="adapter.isConnected ? 'text-primary' : 'text-muted-foreground'"
        >
          {{ adapter.isConnected ? 'Plugged in & Active' : 'Disconnected / On Battery' }}
        </span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Adapter Type</span>
        <span class="font-mono text-sm font-medium">{{ adapter.type }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Negotiated Wattage</span>
        <span class="font-mono text-sm font-medium"
          >{{ adapter.negotiatedWatts }} W ({{ adapter.voltageVolts }}V @
          {{ adapter.amperes }}A)</span
        >
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Active Power Draw</span>
        <span class="font-mono text-sm font-medium">{{ adapter.currentDrawWatts }} W</span>
      </div>
    </div>

    <!-- Laptop Host Sub-tab (F30) -->
    <div
      v-else-if="activeTab === 'laptop'"
      class="rounded-xl border border-border bg-card p-6 divide-y divide-border"
    >
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Laptop Model</span>
        <span class="font-mono text-sm font-medium"
          >{{ specs.laptopManufacturer }} {{ specs.laptopModel }}</span
        >
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">System Serial Number</span>
        <MaskedField :value="specs.laptopSerial" />
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">BIOS Firmware</span>
        <span class="font-mono text-sm font-medium">{{ specs.biosVersion }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Linux Kernel</span>
        <span class="font-mono text-sm font-medium">{{ specs.kernelVersion }}</span>
      </div>
      <div class="flex items-center justify-between py-3">
        <span class="text-sm text-muted-foreground">Installed RAM</span>
        <span class="font-mono text-sm font-medium">{{ specs.totalMemoryRAM }}</span>
      </div>
    </div>
  </div>
</template>
