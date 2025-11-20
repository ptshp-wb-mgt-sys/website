import { fileURLToPath, URL } from 'node:url'

import { defineConfig, type PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'

const enableDevtools = process.env.VITE_ENABLE_DEVTOOLS === 'true'

// https://vite.dev/config/
export default defineConfig(async () => {
  const plugins: PluginOption[] = [vue()]

  if (enableDevtools) {
    const { default: vueDevTools } = await import('vite-plugin-vue-devtools')
    plugins.push(vueDevTools())
  }

  return {
    plugins,
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
