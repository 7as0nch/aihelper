<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../../stores/auth';
import ThemeToggle from '../../ThemeToggle.vue';
import { LogOut, User, Settings } from 'lucide-vue-next';
import { getConfig } from '@/config';

defineProps<{
  isCollapsed: boolean;
}>();

const router = useRouter();
const authStore = useAuthStore();
const isSettingsOpen = ref(false);

const handleLogout = () => {
  authStore.logout();
  router.push('/');
};

// Expose settings state to parent
defineExpose({
  isSettingsOpen
});
</script>

<template>
  <div class="p-3 border-t border-gray-200 dark:border-gray-800">
    <!-- Expanded State -->
    <div v-if="!isCollapsed" class="flex items-center justify-between gap-2">
      <ThemeToggle />
      
      <div v-if="authStore.isAuthenticated" class="flex items-center gap-2 flex-1 min-w-0">
        <img 
          :src="authStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'" 
          class="w-8 h-8 rounded-full bg-gray-100"
        />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ authStore.user?.userName }}</div>
        </div>
        <button 
          @click="isSettingsOpen = true"
          class="p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
          title="设置"
        >
          <Settings class="w-4 h-4" />
        </button>
        <button 
          @click="handleLogout"
          class="p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
          title="退出登录"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>

      <button 
        v-else-if="getConfig('VITE_AI_TYPE') !== 'frontend'"
        @click="authStore.openModal()"
        class="flex items-center justify-center gap-2 px-3 py-2 rounded-lg bg-primary text-white hover:bg-blue-600 transition-colors flex-1 text-sm font-medium"
      >
        登录
      </button>
    </div>

    <!-- Collapsed State -->
    <div v-else class="flex flex-col gap-2 items-center">
      <ThemeToggle />
      
      <div v-if="authStore.isAuthenticated" class="relative group">
        <img 
          :src="authStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'" 
          class="w-8 h-8 rounded-full bg-gray-100 cursor-pointer"
        />
        <div class="absolute left-full bottom-0 ml-2 mb-[-4px] hidden group-hover:block z-50">
          <div class="flex flex-col bg-white dark:bg-[#2a2a2a] border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg overflow-hidden">
            <button 
              @click="isSettingsOpen = true"
              class="flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 whitespace-nowrap"
            >
              <Settings class="w-4 h-4" />
              设置
            </button>
            <button 
              @click="handleLogout"
              class="flex items-center gap-2 px-3 py-2 text-sm text-red-500 hover:bg-gray-50 dark:hover:bg-gray-800 whitespace-nowrap"
            >
              <LogOut class="w-4 h-4" />
              退出
            </button>
          </div>
        </div>
      </div>

      <button 
        v-else-if="getConfig('VITE_AI_TYPE') !== 'frontend'"
        @click="authStore.openModal()"
        class="w-8 h-8 flex items-center justify-center rounded-full bg-primary text-white hover:bg-blue-600 transition-colors"
        title="登录"
      >
        <User class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
