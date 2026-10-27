<script setup lang="ts">
import { ref } from 'vue'

import ThemeToggle from '@/components/common/theme-toggle.vue'

const tempUnit = ref<'C' | 'F'>('C')
const launchAtLogin = ref(true)
const mockDualBattery = ref(false)
const mockUnavailableSensors = ref(false)
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8">
      <h1 class="font-serif text-4xl tracking-tight text-foreground">Settings</h1>
      <p class="text-muted-foreground mt-1.5">
        Application preferences, appearance, and sensor simulation.
      </p>
    </header>

    <div class="space-y-6">
      <!-- General Preferences -->
      <div class="rounded-xl border border-border bg-card p-6 divide-y divide-border">
        <div class="flex items-center justify-between py-3">
          <div>
            <div class="text-sm font-medium">Temperature Unit</div>
            <div class="text-xs text-muted-foreground">Select Celsius (°C) or Fahrenheit (°F)</div>
          </div>
          <div
            class="inline-flex rounded-lg border border-border bg-muted p-0.5 text-xs font-medium"
          >
            <button
              v-for="u in ['C', 'F']"
              :key="u"
              type="button"
              :class="[
                'px-3 py-1 rounded transition',
                tempUnit === u
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground',
              ]"
              @click="tempUnit = u as any"
            >
              °{{ u }}
            </button>
          </div>
        </div>

        <div class="flex items-center justify-between py-3">
          <div>
            <div class="text-sm font-medium">Theme Appearance</div>
            <div class="text-xs text-muted-foreground">
              Light, dark, or follow operating system palette
            </div>
          </div>
          <div class="flex items-center gap-2">
            <ThemeToggle />
          </div>
        </div>

        <div class="flex items-center justify-between py-3">
          <div>
            <div class="text-sm font-medium">Launch at Login</div>
            <div class="text-xs text-muted-foreground">
              Start silently in the system tray when desktop boots
            </div>
          </div>
          <button
            type="button"
            :class="[
              'h-6 w-11 rounded-full transition-colors relative',
              launchAtLogin ? 'bg-primary' : 'bg-muted border border-border',
            ]"
            role="switch"
            :aria-checked="launchAtLogin"
            @click="launchAtLogin = !launchAtLogin"
          >
            <span
              :class="[
                'block h-5 w-5 rounded-full bg-background shadow transition-transform',
                launchAtLogin ? 'translate-x-5' : 'translate-x-0.5',
              ]"
            />
          </button>
        </div>
      </div>

      <!-- Developer & Design Preview Switches -->
      <div class="rounded-xl border border-border bg-card p-6">
        <h2 class="text-xs uppercase font-semibold text-muted-foreground tracking-wider mb-4">
          Hardware Simulation Preview
        </h2>

        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm font-medium">Simulate Dual Battery Setup (F22)</div>
              <div class="text-xs text-muted-foreground">
                Displays split dual vessels on dashboard
              </div>
            </div>
            <input
              v-model="mockDualBattery"
              type="checkbox"
              class="accent-primary h-4 w-4 rounded"
            />
          </div>

          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm font-medium">
                Simulate Limited Sensors (Windows without Drivers)
              </div>
              <div class="text-xs text-muted-foreground">
                Tests graceful "Not available" and "Unknown" states
              </div>
            </div>
            <input
              v-model="mockUnavailableSensors"
              type="checkbox"
              class="accent-primary h-4 w-4 rounded"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
