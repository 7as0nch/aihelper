import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'



// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  build: {
    // 提高chunk大小警告阈值到1000kB
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        // 代码分割配置
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('element-plus')) {
              return 'element-plus'
            }
            if (id.includes('vant')) {
              return 'vant'
            }
            if (id.includes('axios')) {
              return 'axios'
            }
            if (id.includes('pinia')) {
              return 'pinia'
            }
            if (id.includes('vue-router')) {
              return 'vue-router'
            }
            // 其他第三方库放入vendor chunk
            return 'vendor'
          }
        }
      }
    }
  }
})
