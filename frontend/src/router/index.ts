import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import h5Routes from './h5Routes';
import adminRoutes from './adminRoutes';

// 404路由
const notFoundRoute: RouteRecordRaw = {
  path: '/:pathMatch(.*)*',
  name: 'not-found',
  component: () => import('../views/h5/HomeView.vue'),
  meta: {
    title: '页面不存在',
    auth: false
  }
};

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/home'
    },
    ...h5Routes,
    ...adminRoutes,
    notFoundRoute
  ]
});

// 路由导航守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string;
  }
  
  // 获取用户认证信息
  const token = localStorage.getItem('token');
  const isLoggedIn = !!token;
  const userRole = localStorage.getItem('userRole') || 'user';
  
  // 需要认证但未登录的情况
  if (to.meta.auth && !isLoggedIn) {
    // 判断是管理后台还是H5页面，重定向到对应的登录页面
    if (to.path.startsWith('/admin')) {
      next('/admin/login');
    } else {
      next('/login');
    }
    return;
  }
  
  // 需要管理员权限但不是管理员的情况
  if (to.meta.role === 'admin' && userRole !== 'admin') {
    if (isLoggedIn) {
      // 如果已登录但没有管理员权限，重定向到首页
      next('/');
    } else {
      // 如果未登录，重定向到管理后台登录页
      next('/admin/login');
    }
    return;
  }
  
  // 已登录用户访问登录或注册页面时，重定向到首页
  if (isLoggedIn && ['/login', '/register', '/admin/login'].includes(to.path)) {
    if (to.path.startsWith('/admin')) {
      next('/admin/dashboard');
    } else {
      next('/home');
    }
    return;
  }
  
  // 其他情况正常放行
  next();
});

export default router;