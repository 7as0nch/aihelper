import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import dts from 'vite-plugin-dts'
import viteCompression from 'vite-plugin-compression'
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
        }),
        viteCompression({
            verbose: true,
            disable: false,
            threshold: 10240,
            algorithm: 'gzip',
            ext: '.gz',
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
                },
                // Exclude image assets from bundle
                assetFileNames: (assetInfo) => {
                    // Keep CSS in root with fixed name
                    if (assetInfo.name && assetInfo.name.endsWith('.css')) {
                        return 'litechat.css';
                    }
                    // Skip copying image files to dist-widget
                    if (assetInfo.name && /\.(png|jpe?g|gif|svg|ico|webp)$/i.test(assetInfo.name)) {
                        return 'assets/[name]-[hash][extname]';
                    }
                    return 'assets/[name]-[hash][extname]';
                }
            }
        },
        // Don't inline any assets - keep them as separate files that won't be copied
        assetsInlineLimit: 0,
        outDir: 'dist-widget',
        // Exclude image files from being copied
        copyPublicDir: false
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
