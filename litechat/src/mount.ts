import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import 'ant-design-vue/dist/reset.css';
import App from './App.vue'
import { createAppRouter } from './router'
import { createWebHistory, createMemoryHistory } from 'vue-router'
import { useThemeStore } from './stores/theme'
import { setRuntimeConfig, type RuntimeConfig } from './config';
import { tracker } from './utils/tracer';
import { request } from './utils/request';

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

    // 
    if ((import.meta as any).env.VITE_TRACKER_ENABLE === 'true') {
        tracker
            .setUserId(0)
            // 2. Before 拦截器：针对【每一条】数据处理
            // 作用：数据清洗、黑名单过滤、添加额外参数
            .before((logItem) => {
                if (logItem.type === 'api' && logItem.data.method === 'OPTIONS') {
                    return false; // 不入队，不记录
                }

                // 示例：给所有埋点增加当前语言环境
                logItem.data.lang = navigator.language;

                return true; // 继续执行
            })
            // 3. After 拦截器：针对【每一次批量发送】处理
            // 作用：监控上报状态，Debug 打印
            // .after((batchData) => {
            //     console.log(`[Tracker] 🚀 成功批量上报 ${batchData.length} 条数据:`, batchData);
            // })
            // 4. 安装插件
            .installRouter(router)
            .installApi(request)
            .installDirective(app);
    } else {
        // Prevent "Failed to resolve directive: tracker" warning when tracking is disabled
        app.directive('tracker', {});
    }

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
