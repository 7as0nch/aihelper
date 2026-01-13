import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import * as path from 'path';

export default defineConfig(({ mode }) => {
    // Load env file based on `mode` in the current working directory.
    // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
    const env = loadEnv(mode, process.cwd(), '');

    return {
        plugins: [vue()],
        resolve: {
            alias: {
                '@': path.resolve(__dirname, './src'),
            },
        },
        server: {
            port: 5174,
            proxy: {
                '/api': {
                    changeOrigin: true,
                    rewrite: (path) => path.replace(/^\/api/, ''),
                    // Use VITE_API_TARGET if defined, otherwise fallback to local backend
                    // You can add VITE_API_TARGET=http://127.0.0.1:6039 to your .env files
                    target: env.VITE_OPENAI_BASE_URL || 'http://127.0.0.1:6039',
                    ws: true,
                },
                '^/log(/|$)': {
                    changeOrigin: true,
                    rewrite: (path) => path.replace(/^\/log/, ''),
                    // Use VITE_API_TARGET if defined, otherwise fallback to local backend
                    // You can add VITE_API_TARGET=http://127.0.0.1:6039 to your .env files
                    target: env.VITE_OPENAI_BASE_URL || 'http://127.0.0.1:6039',
                    ws: true,
                },
            },
        },
    };
});