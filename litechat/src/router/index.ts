import { createRouter } from 'vue-router';
import ChatView from '@/views/ChatView.vue';
import KnowledgeBaseView from '@/views/KnowledgeBaseView.vue';
import HistoryView from '@/views/HistoryView.vue';
import LandingView from '@/views/LandingView.vue';
import ApplyBetaView from '@/views/ApplyBetaView.vue';

interface Route {
  path: string;
  name: string;
  hidden?: boolean;
  component: any;
  meta?: {
    title?: string;
    icon?: string;
  };
}

const routes: Route[] = [
  {
    path: '/',
    name: 'page',
    hidden: true,
    component: LandingView,
    meta: {
      title: '首页',
    },
  },
  {
    path: '/info',
    name: 'Landing',
    hidden: true,
    component: LandingView,
    meta: {
      title: '官网',
    },
  },
  {
    path: '/apply',
    name: 'ApplyBeta',
    hidden: true,
    component: ApplyBetaView,
    meta: {
      title: '申请内测',
    },
  },
  {
    path: '/app',
    name: 'Chat',
    hidden: false,
    component: ChatView,
    meta: {
      title: '搜索',
      icon: 'Search',
    },
  },
  {
    path: '/app/:id',
    name: 'ChatHistory',
    hidden: true,
    component: ChatView,
  },
  {
    path: '/knowledge',
    name: 'KnowledgeBase',
    hidden: true,
    component: KnowledgeBaseView,
    meta: {
      title: '知识库',
      icon: 'BookOpen',
    },
  },
  {
    path: '/meeting_summary',
    name: 'MeetingSummary',
    hidden: false,
    component: () => import('@/views/MeetingSummaryView.vue'),
    meta: {
      title: '会议总结',
      icon: 'FileText',
    },
  },
  {
    path: '/collections',
    name: 'Collections',
    hidden: false,
    component: () => import('@/views/FavoritesView.vue'),
    meta: {
      title: '收藏',
      icon: 'Bookmark',
    },
  },
  {
    path: '/history',
    name: 'History',
    hidden: true,
    component: HistoryView,
  },
  {
    path: '/mock-qq-login',
    name: 'MockQQLogin',
    hidden: true,
    component: () => import('@/views/MockQQLogin.vue'),
    meta: {
      title: 'QQ登录模拟',
    },
  },
];

export function createAppRouter(history: any) {
  const router = createRouter({
    history: history,
    routes: routes,
    scrollBehavior(to, _from, savedPosition) {
      if (savedPosition) {
        return savedPosition;
      }

      if (to.hash) {
        return {
          el: to.hash,
          top: 112,
          behavior: 'smooth',
        };
      }

      return { top: 0 };
    },
  });

  router.beforeEach((to, _from, next) => {
    const publicRoutes = ['page', 'Landing', 'ApplyBeta', 'Chat', 'MockQQLogin'];

    if (!publicRoutes.includes(to.name as string)) {
      const token = localStorage.getItem('token');
      if (!token) {
        import('@/stores/auth').then(({ useAuthStore }) => {
          const authStore = useAuthStore();
          if (!authStore.checkAuth()) {
            authStore.openModal();
            next('/app');
          } else {
            next();
          }
        });
        return;
      }
    }
    next();
  });

  return router;
}