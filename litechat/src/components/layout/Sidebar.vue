<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { useChatStore } from '../../stores/chat';
import SidebarHeader from './sidebar/SidebarHeader.vue';
import SidebarNavigation from './sidebar/SidebarNavigation.vue';
import SidebarHistory from './sidebar/SidebarHistory.vue';
import SidebarFooter from './sidebar/SidebarFooter.vue';
import SettingsModal from '../settings/SettingsModal.vue';

defineProps<{
  isOpen: boolean;
  isCollapsed: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'toggleCollapse'): void;
}>();

const authStore = useAuthStore();
const store = useChatStore();
const footerRef = ref<InstanceType<typeof SidebarFooter>>();

// Watch auth state to fetch/clear history
watch(() => authStore.checkAuth(), (isAuthenticated) => {
  if (isAuthenticated) {
    store.fetchHistoryList();
  } else {
    store.clearHistoryList();
  }
});

  onMounted(() => {
  if (authStore.checkAuth()) {
    store.fetchHistoryList();
  }
});
</script>

<template>
  <div class="contents">
    <!-- Mobile Overlay -->
    <div 
      v-if="isOpen" 
      class="fixed inset-0 bg-black/50 z-40 md:hidden"
      @click="emit('close')"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed md:static inset-y-0 left-0 z-50 bg-[#f9f9f9] dark:bg-[#1a1a1a] border-r border-gray-200 dark:border-gray-800 transform transition-all duration-300 ease-in-out flex flex-col"
      :class="[
        isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
        isCollapsed ? 'md:w-20' : 'w-64'
      ]"
    >
      <!-- Header -->
      <SidebarHeader 
        :is-collapsed="isCollapsed"
        @close="emit('close')"
        @toggle-collapse="emit('toggleCollapse')"
      />
      <SidebarNavigation :is-collapsed="isCollapsed" />
      <!-- Nav Items -->
      <nav class="flex-1 flex flex-col min-h-0">
        <SidebarHistory 
          :is-collapsed="isCollapsed" 
          class="flex-1 min-h-0"
          @scroll-end="store.fetchHistoryList(true)"
        />
        
        <!-- Loading Indicator -->
        <div v-if="store.isHistoryLoading" class="flex justify-center p-2 shrink-0">
          <div class="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
        </div>
      </nav>

      <!-- Footer -->
      <SidebarFooter 
        ref="footerRef"
        :is-collapsed="isCollapsed"
      />
    </aside>

    <SettingsModal 
      :is-open="footerRef?.isSettingsOpen || false" 
      @close="footerRef && (footerRef.isSettingsOpen = false)" 
    />
  </div>
</template>
