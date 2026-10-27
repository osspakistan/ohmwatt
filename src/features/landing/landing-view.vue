<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div
    class="flex min-h-screen flex-col items-center justify-center p-6 text-center font-sans bg-background text-foreground"
  >
    <div class="max-w-md space-y-5">
      <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-muted p-3">
        <img src="/antlion-mark.svg" alt="Antlion Logo" class="h-10 w-10" />
      </div>

      <h1 class="font-serif text-4xl font-semibold tracking-tight">Antlion</h1>
      <p class="text-sm text-muted-foreground leading-relaxed">
        Desktop battery health, lifespan protection, and zero-drain power management for Linux,
        Windows, and macOS.
      </p>

      <div class="pt-4">
        <button
          type="button"
          class="rounded-lg bg-primary px-5 py-2.5 text-sm font-medium text-primary-foreground shadow transition hover:opacity-90"
          @click="router.push('/app/dashboard')"
        >
          Open Battery Monitor →
        </button>
      </div>
    </div>
  </div>
</template>
