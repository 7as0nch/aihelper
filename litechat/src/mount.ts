import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import 'ant-design-vue/dist/reset.css';
import App from './App.vue'
import { createAppRouter } from './router'
import { createWebHistory, createMemoryHistory } from 'vue-router'
import { useThemeStore } from './stores/theme'
import { setRuntimeConfig, type RuntimeConfig } from './config';

export interface MountOptions {
    config?: Partial<RuntimeConfig>;
    container?: string | HTMLElement;
    routerMode?: 'web' | 'memory';
    initialPath?: string;
    styleContainer?: HTMLElement;
}

export function mountApp(options: MountOptions = {}) {
    if (options.config) {
        setRuntimeConfig(options.config);
    }

    const pinia = createPinia()
    const app = createApp(App, { styleContainer: options.styleContainer })

    app.use(pinia)

    const history = options.routerMode === 'memory' ? createMemoryHistory() : createWebHistory();
    const router = createAppRouter(history);
    app.use(router)

    const container = options.container || '#app';
    app.mount(container)

    if (options.initialPath) {
        router.push(options.initialPath);
    }

    // Initialize theme store to apply theme on load
    useThemeStore()

    return { app, router, pinia };
}
