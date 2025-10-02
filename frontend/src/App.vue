<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

// 路由和状态管理
const route = useRoute();
const router = useRouter();

// 计算页面标题
const pageTitle = computed(() => {
  if (route.meta.title) {
    return `${route.meta.title} | ${import.meta.env.VITE_APP_TITLE || 'AI聊天助手'}`;
  }
  return import.meta.env.VITE_APP_TITLE || 'AI聊天助手';
});

// 监听路由变化更新页面标题
onMounted(() => {
  // 初始化页面标题
  document.title = pageTitle.value;
  
  // 监听路由变化
  const unwatch = route.meta.title && router.afterEach((to) => {
    if (to.meta.title) {
      document.title = `${to.meta.title} | ${import.meta.env.VITE_APP_TITLE || 'AI聊天助手'}`;
    } else {
      document.title = import.meta.env.VITE_APP_TITLE || 'AI聊天助手';
    }
  });
  
  // 在组件卸载时取消监听
  return () => {
    if (typeof unwatch === 'function') {
      unwatch();
    }
  };
});
</script>

<template>
  <div id="app">
    <!-- 路由视图 -->
    <router-view />
  </div>
</template>

<style scoped>
/* 全局应用样式已在 style.css 中定义 */
</style>
