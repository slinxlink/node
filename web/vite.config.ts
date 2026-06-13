import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'

export default defineConfig(({ mode }) => ({
  plugins: [
    vue(),
    vueDevTools(),
    AutoImport({
      imports: ['vue', 'vue-router'],
      dts: true
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  build: {
    assetsDir: 'assets',
    rollupOptions: {
      input: {
        main: fileURLToPath(new URL('./index.html', import.meta.url)),
        sub: fileURLToPath(new URL('./sub.html', import.meta.url)),
      },
      output: {
        entryFileNames: 'assets/[name]-[hash].js',
      }
    }
  },
  server: mode === 'development' ? {
    proxy: {
      '/api': {
        target: 'http://localhost:11111',
        changeOrigin: true,
        ws: true
      }
    }
  } : {}
}))