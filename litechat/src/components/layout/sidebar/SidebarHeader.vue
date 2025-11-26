<script setup lang="ts">
import { X, PanelLeftClose, PanelLeftOpen } from 'lucide-vue-next';
import { getConfig } from '../../../config';

defineProps<{
  isCollapsed: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'toggleCollapse'): void;
}>();

const appTitle = getConfig('VITE_APP_TITLE') || 'AI Chat';
const appLogo = getConfig('VITE_APP_LOGO') || '/logo.png';
</script>

<template>
  <div 
    class="flex flex-col transition-all duration-300"
    :class="[
      isCollapsed 
        ? 'py-4 gap-4 items-center' 
        : 'py-4 px-4 gap-4'
    ]"
  >
    <!-- Desktop Collapse Button -->
    <div :class="[isCollapsed ? 'w-full flex justify-center' : 'w-full flex justify-start']">
      <button 
        class="hidden md:flex text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-800 p-1 rounded-lg transition-colors shrink-0"
        @click="emit('toggleCollapse')"
        :title="isCollapsed ? '展开侧边栏' : '收起侧边栏'"
      >
        <PanelLeftOpen v-if="isCollapsed" class="w-5 h-5" />
        <PanelLeftClose v-else class="w-5 h-5" />
      </button>
    </div>

    <!-- Logo & Title -->
    <div class="flex items-center gap-2 overflow-hidden whitespace-nowrap" :class="{ 'justify-center': isCollapsed }">
      <img :src="appLogo" class="w-8 h-8 rounded shrink-0" alt="Logo" />
      <span v-show="!isCollapsed" class="font-bold text-xl text-gray-900 dark:text-white truncate">{{ appTitle }}</span>
    </div>

    <!-- Mobile Close Button -->
    <button 
      class="md:hidden absolute top-4 right-4 text-gray-500"
      @click="emit('close')"
    >
      <X class="w-6 h-6" />
    </button>
  </div>
</template>
