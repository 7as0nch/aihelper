import { createApp } from 'vue';
import { createPinia } from 'pinia';
import 'ant-design-vue/dist/reset.css';
import './style.css';
import './styles/landing-ant.css';

import { StyleProvider } from 'ant-design-vue';
import App from './App.vue';
import { i18n } from './i18n';
import { createAppRouter } from './router';
import { createMemoryHistory, createWebHistory } from 'vue-router';
import { setRuntimeConfig, type RuntimeConfig } from './config';
import { useLocaleStore } from './stores/locale';
import { useThemeStore } from './stores/theme';
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

  const pinia = createPinia();
  const app = createApp(App, { styleContainer: options.styleContainer });

  app.component('StyleProvider', StyleProvider);
  app.use(pinia);
  app.use(i18n);

  useThemeStore();
  useLocaleStore().initialize();

  const history = options.routerMode === 'memory' ? createMemoryHistory() : createWebHistory();
  const router = createAppRouter(history);

  if ((import.meta as any).env.VITE_TRACKER_ENABLE === 'true') {
    tracker
      .setUserId(0)
      .before((logItem) => {
        if (logItem.type === 'api' && logItem.data.method === 'OPTIONS') {
          return false;
        }

        logItem.data.lang = navigator.language;
        return true;
      })
      .installRouter(router)
      .installApi(request)
      .installDirective(app);
  } else {
    app.directive('tracker', {});
  }

  app.use(router);

  const container = options.container || '#app';
  app.mount(container);

  if (options.initialPath) {
    router.push(options.initialPath);
  }

  return { app, router, pinia };
}