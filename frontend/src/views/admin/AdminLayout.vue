<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { ElMenu, ElMenuItem, ElAvatar, ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus';
import 'element-plus/es/components/menu/style/css';
import 'element-plus/es/components/avatar/style/css';
import 'element-plus/es/components/dropdown/style/css';
import 'element-plus/es/components/message/style/css';

const router = useRouter();
const userStore = useUserStore();
const collapsed = ref(false);
const userInfo = ref<any>({});

// 侧边栏菜单
const menuItems = [
  {
    index: 'dashboard',
    label: '仪表盘',
    icon: 'Document',
    path: '/admin/dashboard'
  },
  {
    index: 'user',
    label: '用户管理',
    icon: 'User',
    path: '/admin/users'
  },
  {
    index: 'function',
    label: '函数工具管理',
    icon: 'Code',
    path: '/admin/functions'
  },
  {
    index: 'workflow',
    label: '工作流管理',
    icon: 'Operation',
    path: '/admin/workflows'
  },
  {
    index: 'settings',
    label: '系统设置',
    icon: 'Setting',
    path: '/admin/settings'
  }
];



// 初始化
onMounted(() => {
  // 检查用户是否已登录且为管理员
  if (!userStore.isLoggedIn || userStore.userInfo?.role !== 'admin') {
    router.push('/admin/login');
    return;
  }
  
  // 获取用户信息
  userInfo.value = userStore.userInfo || {};
});

// 菜单点击处理
const handleMenuClick = (index: string) => {
  const menu = menuItems.find(item => item.index === index);
  if (menu) {
    router.push(menu.path);
  }
};

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  collapsed.value = !collapsed.value;
};


</script>

<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed }">
      <div class="sidebar-header">
        <div class="logo-container">
          <h1 class="logo-text" v-if="!collapsed">管理后台</h1>
          <h1 class="logo-text collapsed" v-else>后台</h1>
        </div>
        <button class="collapse-btn" @click="toggleCollapse">
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            width="16" 
            height="16" 
            fill="currentColor" 
            viewBox="0 0 16 16"
          >
            <path 
              d="M7.5 12a.5.5 0 0 1-.5-.5V3.707L2.354 9.854a.5.5 0 1 1-.708-.708l5-5a.5.5 0 0 1 .708 0l5 5a.5.5 0 0 1-.708.708L8 3.707V11.5a.5.5 0 0 1-.5.5z" 
              v-if="!collapsed"
            />
            <path 
              d="M7.646 15.854a.5.5 0 0 0 .708 0l6-6a.5.5 0 0 0-.708-.708L8 14.293l-5.646-5.647a.5.5 0 0 0-.708.708l6 6z" 
              v-if="collapsed"
            />
          </svg>
        </button>
      </div>
      
      <ElMenu 
        active-text-color="#409EFF"
        background-color="#1f2937"
        class="menu" 
        default-active="dashboard"
        text-color="#9ca3af"
        unique-opened
        :collapse="collapsed"
        :collapse-transition="false"
        @select="handleMenuClick"
      >
        <ElMenuItem 
          v-for="item in menuItems" 
          :key="item.index"
          :index="item.index"
        >
          <template #title>
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              fill="currentColor" 
              viewBox="0 0 16 16"
            >
              <path v-if="item.icon === 'Document'" d="M14 14V4.5L9.5 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2zM9.5 3A1.5 1.5 0 0 0 11 4.5h2V14a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1h5.5v2z"/>
              <path v-if="item.icon === 'User'" d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z"/>
              <path v-if="item.icon === 'Code'" d="M5.854 4.646a.5.5 0 0 1 0 .708L2.707 8l3.147 3.146a.5.5 0 0 1-.708.708l-3.5-3.5a.5.5 0 0 1 0-.708l3.5-3.5a.5.5 0 0 1 .708 0zm4.292 0a.5.5 0 0 0 0 .708L13.293 8l-3.147 3.146a.5.5 0 0 0 .708.708l3.5-3.5a.5.5 0 0 0 0-.708l-3.5-3.5a.5.5 0 0 0-.708 0z"/>
              <path v-if="item.icon === 'Operation'" d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
              <path v-if="item.icon === 'Operation'" d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
              <path v-if="item.icon === 'Setting'" d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
              <path v-if="item.icon === 'Setting'" d="M5.255 5.786a.237.237 0 0 0 .241.247h.754a.238.238 0 0 0 .238-.247l-.005-.861a.248.248 0 0 0-.237-.255H5.29a.25.25 0 0 0-.25.25l.006.861zm1.508-.934a.242.242 0 0 0 0 .254l.861.006a.25.25 0 0 0 .25-.25l-.006-.861a.242.242 0 0 0-.254 0zm2.496 0a.242.242 0 0 1 0 .254l.861.006a.25.25 0 0 1 .25-.25l-.006-.861a.242.242 0 0 1-.254 0zm1.23-.61a.25.25 0 0 1 .25.25l-.006.861a.242.242 0 0 1-.254 0l-.861-.006a.25.25 0 0 1-.25-.25l.006-.861a.25.25 0 0 1 .25-.25h.754zM5.255 12.21a.25.25 0 0 1 .25-.25h.754a.25.25 0 0 1 .25.25l-.006.861a.242.242 0 0 1-.254 0l-.861-.006a.25.25 0 0 1-.25-.25zm1.508 0a.25.25 0 0 0 .25-.25l-.006-.861a.242.242 0 0 0-.254 0l-.861.006a.25.25 0 0 0-.25.25l.006.861a.25.25 0 0 0 .25.25h.754zm2.496-.01a.25.25 0 0 0 .25-.25l-.006-.861a.242.242 0 0 0-.254 0l-.861.006a.25.25 0 0 0-.25.25l.006.861a.25.25 0 0 0 .25.25h.754zm1.23-1.734a.262.262 0 0 1-.064.125.252.252 0 0 1-.098.07 2.67 2.67 0 0 1-.566.148c-.37.028-.68.21-.865.501-.184.29-.207.653-.175 1.005l.021.253a.25.25 0 0 1-.25.245h-.754a.25.25 0 0 1-.25-.25l.006-.861a.25.25 0 0 1 .25-.25h.504c.176 0 .329.08.405.218.076.138.09.312.07 1.027a7.152 7.152 0 0 0-.12 1.277 2.67 2.67 0 0 0-.429.564c-.105.17-.296.293-.63.478-.334.185-.815.306-1.277.306a3.12 3.12 0 0 1-1.027-.175 2.67 2.67 0 0 1-.564-.429 2.67 2.67 0 0 1-.148-.566 2.658 2.658 0 0 1 .148-.566c.2-.37.312-.814.306-1.277a3.12 3.12 0 0 0-.175-1.027 2.67 2.67 0 0 1 .429-.564c.171-.105.477-.295.63-.478.185-.335.306-.815.306-1.277a3.12 3.12 0 0 0-.175-1.027 2.67 2.67 0 0 1-.564-.429 2.67 2.67 0 0 1-.566-.148 2.658 2.658 0 0 1-.566.148 2.67 2.67 0 0 1-.429.564A3.12 3.12 0 0 0 9.79 5a2.658 2.658 0 0 1-.12-1.277.25.25 0 0 1 .25-.25h.754a.25.25 0 0 1 .25.25l-.021.253a2.658 2.658 0 0 1 .175 1.005c.028.37.21.68.501.865a2.67 2.67 0 0 1 .566.148 2.67 2.67 0 0 1 .566-.148c.29-.185.501-.495.501-.865a2.658 2.658 0 0 0-.175-1.005l.021-.253a.25.25 0 0 1 .25-.245h.754a.25.25 0 0 1 .25.25l-.006.861a.25.25 0 0 1-.25.25h-.504a2.658 2.658 0 0 0-.405.218 2.658 2.658 0 0 0-.07 1.027 7.152 7.152 0 0 1 .12 1.277 2.67 2.67 0 0 1 .429.564c.105.17.295.293.63.478.335.185.815.306 1.277.306a3.12 3.12 0 0 1 1.027-.175 2.67 2.67 0 0 1 .564.429 2.67 2.67 0 0 1 .148.566 2.658 2.658 0 0 1-.148.566 2.67 2.67 0 0 1-.429.564 3.12 3.12 0 0 1-1.027.175z"/>
            </svg>
            <span v-if="!collapsed">{{ item.label }}</span>
          </template>
        </ElMenuItem>
      </ElMenu>
    </aside>

    <!-- 主内容区域 -->
    <div class="main-content" :class="{ 'sidebar-collapsed': collapsed }">
      <!-- 顶部导航栏 -->
      <header class="top-header">
        <div class="header-left">
          <button class="toggle-sidebar-btn" @click="toggleCollapse">
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              fill="currentColor" 
              viewBox="0 0 16 16"
            >
              <path 
                d="M7.5 12a.5.5 0 0 1-.5-.5V3.707L2.354 9.854a.5.5 0 1 1-.708-.708l5-5a.5.5 0 0 1 .708 0l5 5a.5.5 0 0 1-.708.708L8 3.707V11.5a.5.5 0 0 1-.5.5z" 
                v-if="!collapsed"
              />
              <path 
                d="M7.646 15.854a.5.5 0 0 0 .708 0l6-6a.5.5 0 0 0-.708-.708L8 14.293l-5.646-5.647a.5.5 0 0 0-.708.708l6 6z" 
                v-if="collapsed"
              />
            </svg>
          </button>
        </div>
        
        <div class="header-right">
          <ElDropdown class="user-dropdown">
            <span class="user-info">
              <ElAvatar class="user-avatar" :icon="'User'" />
              <span class="username" v-if="!collapsed">{{ userInfo.username || '管理员' }}</span>
            </span>
            <ElDropdownMenu slot="dropdown">
              <ElDropdownItem divided>
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  width="14" 
                  height="14" 
                  fill="currentColor" 
                  viewBox="0 0 16 16" 
                  class="mr-2"
                >
                  <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
                  <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
                </svg>
                退出登录
              </ElDropdownItem>
            </ElDropdownMenu>
          </ElDropdown>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  background-color: var(--bg-primary);
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  width: 200px;
  background-color: #1f2937;
  color: #9ca3af;
  transition: width 0.3s ease;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.sidebar.collapsed {
  width: 60px;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #374151;
}

.logo-container {
  display: flex;
  align-items: center;
}

.logo-text {
  font-size: 1.2rem;
  font-weight: 600;
  color: white;
  margin: 0;
}

.logo-text.collapsed {
  display: none;
}

.collapse-btn {
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 5px;
  border-radius: var(--radius-sm);
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background-color: #374151;
  color: white;
}

.menu {
  flex: 1;
  border-right: none;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
}

.main-content.sidebar-collapsed {
  margin-left: 0;
}

/* 顶部导航栏 */
.top-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  background-color: white;
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
}

.toggle-sidebar-btn {
  background: none;
  border: none;
  color: var(--text-primary);
  cursor: pointer;
  padding: 10px;
  border-radius: var(--radius-sm);
  transition: all 0.3s ease;
}

.toggle-sidebar-btn:hover {
  background-color: var(--bg-secondary);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  transition: background-color 0.3s ease;
}

.user-info:hover {
  background-color: var(--bg-secondary);
}

.user-avatar {
  margin-right: 10px;
}

.username {
  font-size: 0.9rem;
  color: var(--text-primary);
}

/* 内容区域 */
.content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: var(--bg-primary);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    position: fixed;
    left: -200px;
    top: 0;
    bottom: 0;
    z-index: 200;
  }
  
  .sidebar.collapsed {
    left: 0;
    width: 60px;
  }
  
  .main-content {
    margin-left: 0;
  }
  
  .toggle-sidebar-btn {
    display: block;
  }
}

@media (max-width: 768px) {
  .content {
    padding: 15px;
  }
  
  .top-header {
    padding: 0 15px;
  }
}
</style>