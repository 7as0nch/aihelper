import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/home',
    name: 'home',
    component: () => import('../views/h5/HomeView.vue'),
    meta: {
      title: '首页',
      auth: false
    }
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import('../views/h5/ChatView.vue'),
    meta: {
      title: '聊天',
      auth: true
    }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/h5/ProfileView.vue'),
    meta: {
      title: '个人中心',
      auth: true
    }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/h5/LoginView.vue'),
    meta: {
      title: '登录',
      auth: false
    }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/h5/RegisterView.vue'),
    meta: {
      title: '注册',
      auth: false
    }
  }
];

export default routes;