import { toPng } from 'html-to-image'
import { ref } from 'vue'

import { ScreenshotService } from '../../bindings/changeme/service'

const isCapturing = ref(false)
const lastSavedPath = ref<string | null>(null)

export function useScreenshot() {
  async function captureWindow(customSelector = '#app'): Promise<string | null> {
    if (isCapturing.value) return null
    isCapturing.value = true

    try {
      const node = document.querySelector(customSelector) as HTMLElement | null
      const target = node || document.body

      // High-resolution Retina capture (pixelRatio: 2)
      const dataUrl = await toPng(target, {
        quality: 1.0,
        pixelRatio: 2,
        backgroundColor: '#111111',
        filter: (child) => {
          // Exclude any temporary toast overlays if desired
          if (child instanceof HTMLElement && child.dataset.noScreenshot) {
            return false
          }
          return true
        },
      })

      const savedPath = await ScreenshotService.SaveScreenshot(dataUrl)
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
    __antlion_capture_screenshot?: (selector?: string) => Promise<string | null>
  }
}

if (typeof window !== 'undefined') {
  const { captureWindow } = useScreenshot()
  window.__antlion_capture_screenshot = captureWindow
}
