import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'


// https://vite.dev/config/
export default {
  plugins: [vue(), tailwindcss()],
  build: {
    // 提高chunk大小警告阈值到1000kB
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      external: [
        // @ant-design/icons依赖react，但项目中未安装react
        'react'
      ]
    }
  }
}
