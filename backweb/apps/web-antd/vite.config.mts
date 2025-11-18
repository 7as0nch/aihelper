import { resolve } from 'node:path';
import { cwd } from 'node:process';

import { defineConfig } from '@vben/vite-config';

// 自行取消注释来启用按需导入功能
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers';
import Components from 'unplugin-vue-components/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      define: {
        // 注入项目根路径到运行时（ruoyi-admin目录层级）
        __PROJECT_ROOT__: JSON.stringify(
          (() => {
            // 获取当前工作目录字符串
            const currentWorkingDir = cwd();

            console.warn('Vite当前工作目录:', currentWorkingDir);
            // 规范化路径分隔符，统一处理 Windows 和 macOS/Linux
            const normalizedPath = currentWorkingDir.replaceAll('\\', '/');

            // 如果当前目录包含 /apps/web-antd，则向上两级到 ruoyi-admin
            if (normalizedPath.includes('/apps/web-antd')) {
              return resolve(currentWorkingDir, '../..');
            }
            // 如果当前目录包含 /apps，则向上一级到 ruoyi-admin
            else if (normalizedPath.includes('/apps')) {
              return resolve(currentWorkingDir, '..');
            }
            // 否则返回当前目录
            else {
              return currentWorkingDir;
            }
          })(),
        ),
      },
      plugins: [
        Components({
          dirs: [], // 默认会导入src/components目录下所有组件 不需要
          dts: './types/components.d.ts', // 输出类型文件
          resolvers: [
            AntDesignVueResolver({
              // 需要排除Button组件 全局已经默认导入了
              exclude: ['Button'],
              importStyle: false, // css in js
            }),
          ],
        }),
      ],
      server: {
        proxy: {
          '/api': {
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api/, ''),
            // mock代理目标地址
            target: 'http://127.0.0.1:6039',
            ws: true,
          },
        },
      },
    },
  };
});
