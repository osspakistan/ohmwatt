<script setup lang="ts">
import ThemeToggle from '@/components/common/theme-toggle.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useTheme } from '@/composables/use-theme'

const { mode } = useTheme()
</script>

<template>
  <div class="flex h-screen w-screen overflow-hidden bg-background text-foreground font-sans">
    <!-- Slim Navigation Icon Rail -->
    <aside
      class="flex w-16 flex-col items-center justify-between border-r border-border bg-card py-5"
    >
      <div class="flex flex-col items-center gap-6">
        <!-- Antlion Brand Mark -->
        <div
          class="group relative flex h-10 w-10 items-center justify-center rounded-lg bg-muted p-2 transition-colors hover:bg-accent"
          title="Antlion Battery Monitor"
        >
          <img src="/antlion-mark.svg" alt="Antlion Logo" class="h-6 w-6" />
        </div>
      </div>

      <!-- Bottom controls (Theme Toggle) -->
      <div class="flex flex-col items-center gap-3">
        <ThemeToggle />
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 overflow-y-auto p-8">
      <header class="mb-8 flex items-center justify-between">
        <div>
          <h1 class="font-serif text-3xl font-semibold tracking-tight">Antlion</h1>
          <p class="text-sm text-muted-foreground">Desktop battery health and power management</p>
        </div>
        <div class="flex items-center gap-2 text-xs font-mono text-muted-foreground">
          <span>Theme:</span>
          <span class="rounded bg-muted px-2 py-0.5 text-foreground capitalize">{{ mode }}</span>
        </div>
      </header>

      <!-- Welcome Overview Card demonstrating Caffeine Tokens in Dark/Light -->
      <div class="grid gap-6 md:grid-cols-2">
        <Card class="border-border bg-card">
          <CardHeader>
            <CardTitle class="font-serif text-xl font-medium">Battery Overview</CardTitle>
            <CardDescription>Live telemetry and health status</CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div
              class="flex items-center justify-between rounded-md border border-border bg-muted/40 p-4"
            >
              <span class="text-sm text-muted-foreground">Current Level</span>
              <span class="font-serif text-3xl font-bold tabular-nums text-foreground">78%</span>
            </div>

            <div class="flex items-center justify-between rounded-md border border-border p-4">
              <span class="text-sm text-muted-foreground">Status</span>
              <span
                class="inline-flex items-center gap-1.5 rounded-full bg-secondary/20 px-2.5 py-0.5 text-xs font-medium text-secondary-foreground border border-secondary/40"
              >
                ⚡ Plugged in, charging
              </span>
            </div>
          </CardContent>
        </Card>

        <Card class="border-border bg-card">
          <CardHeader>
            <CardTitle class="font-serif text-xl font-medium">Caffeine Theme System</CardTitle>
            <CardDescription>Theme verification & visual indicators</CardDescription>
          </CardHeader>
          <CardContent class="space-y-3 text-sm">
            <div class="flex items-center justify-between">
              <span class="text-muted-foreground">Primary Accent:</span>
              <span class="inline-block h-6 w-16 rounded bg-primary"></span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-muted-foreground">90% Heads-Up Marker:</span>
              <span class="inline-block h-6 w-16 rounded bg-secondary"></span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-muted-foreground">Alert / Critical:</span>
              <span class="inline-block h-6 w-16 rounded bg-destructive"></span>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  </div>
</template>
