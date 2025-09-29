import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/admin',
    name: 'admin',
    redirect: '/admin/dashboard',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: {
      title: '管理后台',
      auth: true,
      role: 'admin'
    },
    children: [
      {
        path: 'dashboard',
        name: 'adminDashboard',
        component: () => import('../views/admin/DashboardView.vue'),
        meta: {
          title: '仪表盘',
          auth: true,
          role: 'admin',
          menuName: '仪表盘',
          menuIcon: 'el-icon-data-board'
        }
      },
      {
        path: 'users',
        name: 'adminUsers',
        component: () => import('../views/admin/UserManagementView.vue'),
        meta: {
          title: '用户管理',
          auth: true,
          role: 'admin',
          menuName: '用户管理',
          menuIcon: 'el-icon-user-solid'
        }
      },
      {
        path: 'functions',
        name: 'adminFunctionTools',
        component: () => import('../views/admin/FunctionToolManagementView.vue'),
        meta: {
          title: '函数工具管理',
          auth: true,
          role: 'admin',
          menuName: '函数工具管理',
          menuIcon: 'el-icon-s-tools'
        }
      },
      {
        path: 'workflows',
        name: 'adminWorkflows',
        component: () => import('../views/admin/WorkflowManagementView.vue'),
        meta: {
          title: '工作流管理',
          auth: true,
          role: 'admin',
          menuName: '工作流管理',
          menuIcon: 'el-icon-s-order'
        }
      },
      {
        path: 'settings',
        name: 'adminSettings',
        component: () => import('../views/admin/SystemSettingsView.vue'),
        meta: {
          title: '系统设置',
          auth: true,
          role: 'admin',
          menuName: '系统设置',
          menuIcon: 'el-icon-setting'
        }
      }
    ]
  },
  {
    path: '/admin/login',
    name: 'adminLogin',
    component: () => import('../views/admin/AdminLoginView.vue'),
    meta: {
      title: '管理员登录',
      auth: false
    }
  }
];

export default routes;