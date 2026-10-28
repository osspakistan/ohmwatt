<script setup lang="ts">
import { ref } from 'vue'

const headsUpEnabled = ref(true)
const threshold = ref(90) // 80 - 100%
const repeatInterval = ref<'off' | '10' | '30'>('off')

const heatWarningEnabled = ref(true)
const heatThreshold = ref(40) // 35 - 55 °C
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8">
      <h1 class="font-serif text-4xl tracking-tight text-foreground">Alerts & Notifications</h1>
      <p class="text-muted-foreground mt-1.5">
        Proactive reminders to prevent high-voltage overcharging and thermal stress.
      </p>
    </header>

    <div class="space-y-6">
      <!-- 90% Heads-Up Unplug Alarm (F17) -->
      <div class="rounded-xl border border-border bg-card p-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-serif text-lg font-medium text-foreground">
              Plugged-in Charge Reminder
            </h2>
            <p class="text-xs text-muted-foreground mt-0.5">
              Fires a desktop alert when battery reaches your target threshold while plugged into
              AC.
            </p>
          </div>
          <button
            type="button"
            :class="[
              'h-6 w-11 rounded-full transition-colors relative focus-visible:outline-2 focus-visible:outline-primary',
              headsUpEnabled ? 'bg-primary' : 'bg-muted border border-border',
            ]"
            role="switch"
            :aria-checked="headsUpEnabled"
            @click="headsUpEnabled = !headsUpEnabled"
          >
            <span
              :class="[
                'block h-5 w-5 rounded-full bg-background shadow transition-transform',
                headsUpEnabled ? 'translate-x-5' : 'translate-x-0.5',
              ]"
            />
          </button>
        </div>

        <div v-if="headsUpEnabled" class="mt-6 pt-6 border-t border-border space-y-5">
          <!-- Threshold Slider -->
          <div>
            <div class="flex items-center justify-between text-sm mb-2">
              <span class="text-muted-foreground">Alert Threshold</span>
              <span class="font-serif text-2xl font-medium text-foreground">{{ threshold }}%</span>
            </div>
            <input
              v-model.number="threshold"
              type="range"
              min="80"
              max="100"
              step="1"
              class="w-full accent-primary h-2 bg-muted rounded-lg cursor-pointer"
            />
            <div class="flex justify-between text-[11px] font-mono text-muted-foreground mt-1">
              <span>80%</span>
              <span>85%</span>
              <span>90% (Recommended)</span>
              <span>95%</span>
              <span>100%</span>
            </div>
          </div>

          <!-- Repeat Frequency -->
          <div class="flex items-center justify-between pt-2">
            <div>
              <div class="text-sm font-medium">Repeat Interval</div>
              <div class="text-xs text-muted-foreground">
                Repeat notification if cable remains connected
              </div>
            </div>
            <div class="inline-flex rounded-lg border border-border bg-muted p-0.5 text-xs">
              <button
                v-for="opt in [
                  { value: 'off', label: 'Once' },
                  { value: '10', label: 'Every 10m' },
                  { value: '30', label: 'Every 30m' },
                ]"
                :key="opt.value"
                type="button"
                :class="[
                  'px-3 py-1 rounded transition font-medium',
                  repeatInterval === opt.value
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground',
                ]"
                @click="repeatInterval = opt.value as any"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Thermal Protection Warning (F11-lite) -->
      <div class="rounded-xl border border-border bg-card p-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-serif text-lg font-medium text-foreground">Battery Heat Warning</h2>
            <p class="text-xs text-muted-foreground mt-0.5">
              Dispatches an advisory notification when battery temperature exceeds safety limits
              while charging.
            </p>
          </div>
          <button
            type="button"
            :class="[
              'h-6 w-11 rounded-full transition-colors relative focus-visible:outline-2 focus-visible:outline-primary',
              heatWarningEnabled ? 'bg-primary' : 'bg-muted border border-border',
            ]"
            role="switch"
            :aria-checked="heatWarningEnabled"
            @click="heatWarningEnabled = !heatWarningEnabled"
          >
            <span
              :class="[
                'block h-5 w-5 rounded-full bg-background shadow transition-transform',
                heatWarningEnabled ? 'translate-x-5' : 'translate-x-0.5',
              ]"
            />
          </button>
        </div>

        <div v-if="heatWarningEnabled" class="mt-6 pt-6 border-t border-border space-y-4">
          <div class="flex items-center justify-between text-sm">
            <span class="text-muted-foreground">Temperature Trigger</span>
            <span class="font-mono text-lg font-medium text-destructive"
              >{{ heatThreshold }} °C</span
            >
          </div>
          <input
            v-model.number="heatThreshold"
            type="range"
            min="35"
            max="50"
            step="1"
            class="w-full accent-destructive h-2 bg-muted rounded-lg cursor-pointer"
          />
        </div>
      </div>
    </div>
  </div>
</template>
