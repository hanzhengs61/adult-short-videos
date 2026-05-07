import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      // 开发时代理头像等静态上传文件；生产用 Nginx/CDN 直接服务，此行无效
      '/uploads': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
  build: {
    outDir: '../web/dist',
  },
})
