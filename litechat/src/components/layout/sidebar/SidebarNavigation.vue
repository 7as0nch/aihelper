<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Search, BookOpen, Bookmark, FileText } from 'lucide-vue-next';

defineProps<{
  isCollapsed: boolean;
}>();

const route = useRoute();
const router = useRouter();

// Icon mapping
const iconMap: Record<string, any> = {
  Search,
  BookOpen,
  Bookmark,
  FileText
};

// Dynamic nav items from router
const navItems = computed(() => {
  return router.options.routes
    .filter((r: any) => !r.hidden && r.meta?.title)
    .map((r: any) => ({
      label: r.meta.title,
      path: r.path,
      icon: iconMap[r.meta.icon] || Search,
      active: false
    }));
});
</script>

<template>
  <router-link 
    v-for="item in navItems" 
    :key="item.label"
    :to="item.path" 
    class="flex items-center gap-3 px-3 py-3 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 ease-in-out hover:translate-x-1"
    :class="[
      { 'bg-blue-50 text-primary dark:bg-blue-900/20 dark:text-blue-400': route.path === item.path },
      isCollapsed ? 'justify-center' : ''
    ]"
    :title="isCollapsed ? item.label : ''"
  >
    <component :is="item.icon" class="w-5 h-5 shrink-0" />
    <span v-show="!isCollapsed" class="whitespace-nowrap overflow-hidden">{{ item.label }}</span>
  </router-link>
</template>
