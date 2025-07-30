import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
 

// import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
// import path from 'path'


// https://vite.dev/config/
export default defineConfig({
  base: './', // 关键：使用相对路径，避免打包后路径包含 `src`
  build: {
    assetsDir: 'assets', // 打包后静态资源存储在 `dist/assets` 目录
    rollupOptions: {
      output: {
        assetFileNames:  (assetInfo) => {
          if (assetInfo.name && assetInfo.name.endsWith('.css')) {
            return 'assets/[name]-[hash][extname]'; // 保留 .css 后缀
          }
          return 'assets/[name]-[hash][extname]';
        }
      },
    },
  },
  plugins: [
    vue(),
    tailwindcss(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 你的后端 API 地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api/'),
      },
    },
  },
})
