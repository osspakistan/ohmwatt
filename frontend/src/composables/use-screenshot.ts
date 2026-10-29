import { toPng } from 'html-to-image'
import { ref } from 'vue'

import { ScreenshotService } from '../../bindings/changeme/service'

const isCapturing = ref(false)
const lastSavedPath = ref<string | null>(null)

export function useScreenshot() {
  async function captureWindow(
    customSelector = '#app',
    screenName?: string,
  ): Promise<string | null> {
    if (isCapturing.value) return null
    isCapturing.value = true

    try {
      const node = document.querySelector(customSelector) as HTMLElement | null
      const target = node || document.body

      let resolvedScreen = screenName
      if (!resolvedScreen && typeof window !== 'undefined') {
        const hash = window.location.hash.replace(/^#\/?/, '').split('/')[0]
        const path = window.location.pathname.replace(/^\//, '').split('/')[0]
        resolvedScreen = hash || path || 'dashboard'
      }

      // High-resolution Retina capture (pixelRatio: 2)
      const dataUrl = await toPng(target, {
        quality: 1.0,
        pixelRatio: 2,
        backgroundColor: '#111111',
        filter: (child) => {
          if (child instanceof HTMLElement && child.dataset.noScreenshot) {
            return false
          }
          return true
        },
      })

      const savedPath = await ScreenshotService.SaveScreenshot(
        resolvedScreen || 'dashboard',
        dataUrl,
      )
      lastSavedPath.value = savedPath
      return savedPath
    } catch (err) {
      console.error('Failed to capture window screenshot:', err)
      return null
    } finally {
      isCapturing.value = false
    }
  }

  return {
    captureWindow,
    isCapturing,
    lastSavedPath,
  }
}

// Expose helper on window for external scripts/automation
declare global {
  interface Window {
    __ohmwatt_capture_screenshot?: (
      selector?: string,
      screenName?: string,
    ) => Promise<string | null>
  }
}

if (typeof window !== 'undefined') {
  const { captureWindow } = useScreenshot()
  window.__ohmwatt_capture_screenshot = captureWindow
}
