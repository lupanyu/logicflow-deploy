import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
 

// import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
// import path from 'path'


// https://vite.dev/config/
export default defineConfig({

  plugins: [
    vue(),
    tailwindcss(),
    vueDevTools(),
    // createSvgIconsPlugin({
    //   //将svg图标放入src文件下面中的assets下面中的icons文件夹中
    //     iconDirs: [path.resolve(process.cwd(), "src/assets/icons")],
    //     symbolId: "icon-[dir]-[name]",
    //   }),
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
