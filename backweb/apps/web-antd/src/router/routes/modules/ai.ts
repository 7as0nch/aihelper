import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:bot',
      order: 0,
      title: 'AI 管理',
    },
    name: 'AI',
    path: '/ai',
    children: [
      {
        name: 'AIAgent',
        path: '/ai/agent',
        component: () => import('#/views/ai/agent/index.vue'),
        meta: {
          icon: 'lucide:robot',
          title: 'Agent 管理',
        },
      },
      {
        name: 'AIApplication',
        path: '/ai/app',
        component: () => import('#/views/ai/app/index.vue'),
        meta: {
          icon: 'lucide:workflow',
          title: 'AI 应用',
        },
      },
      {
        name: 'AIApplicationEditor',
        path: '/ai/app/editor/:id?',
        component: () => import('#/views/ai/app/app-editor.vue'),
        meta: {
          hideInMenu: true,
          title: '流程编辑器',
        },
      },
      {
        name: 'AIPrompt',
        path: '/ai/prompt',
        component: () => import('#/views/ai/prompt/index.vue'),
        meta: {
          icon: 'lucide:file-text',
          title: '提示词管理',
        },
      },
      {
        name: 'AIModel',
        path: '/ai/model',
        component: () => import('#/views/ai/model/index.vue'),
        meta: {
          icon: 'lucide:cpu',
          title: '模型管理',
        },
      },
      {
        name: 'AITool',
        path: '/ai/tool',
        component: () => import('#/views/ai/tool/index.vue'),
        meta: {
          icon: 'lucide:wrench',
          title: '工具管理',
        },
      },
    ],
  },
];

export default routes;
