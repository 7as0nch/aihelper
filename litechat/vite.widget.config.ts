import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import dts from 'vite-plugin-dts'
import * as path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        dts({
            include: ['src/widget/**/*.ts', 'src/config.ts', 'src/types/**/*.ts'],
            outDir: 'dist-widget',
            rollupTypes: true,
            copyDtsFiles: true,
            staticImport: true,
            insertTypesEntry: true,
        })
    ],
    build: {
        lib: {
            entry: path.resolve(__dirname, 'src/widget/index.ts'),
            name: 'AiChatWidget',
            fileName: (format) => `litechat-widget.${format}.js`,
            formats: ['es']
        },
        rollupOptions: {
            // Make sure to externalize deps that shouldn't be bundled
            // into your library
            external: [], // We want to bundle everything for a standalone widget
            output: {
                // Provide global variables to use in the UMD build
                // for externalized deps
                globals: {
                    vue: 'Vue'
                }
            }
        },
        outDir: 'dist-widget'
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    define: {
        'process.env': {}
    }
})
