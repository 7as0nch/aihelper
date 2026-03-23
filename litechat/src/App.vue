<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import { StyleProvider } from 'ant-design-vue';
import Sidebar from '@/components/layout/Sidebar.vue';
import AuthModal from '@/components/auth/AuthModal.vue';
import { useAuthStore } from '@/stores/auth';

defineProps<{
  styleContainer?: HTMLElement;
}>();

const route = useRoute();
const authStore = useAuthStore();
const isSidebarOpen = ref(false);
const isSidebarCollapsed = ref(false);

const isLanding = computed(() => ['Landing', 'page'].includes(route.name as string));

watch(isLanding, (landing) => {
  if (landing) {
    authStore.closeModal();
  }
}, { immediate: true });
</script>

<template>
  <StyleProvider :container="styleContainer" :hash-priority="'high'">
    <div :class="[
      'w-full bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-gray-900 dark:to-gray-800',
      isLanding ? 'min-h-screen overflow-auto' : 'flex h-screen overflow-hidden'
    ]">
      <Sidebar 
        v-if="!isLanding"
        :is-open="isSidebarOpen" 
        :is-collapsed="isSidebarCollapsed"
        @close="isSidebarOpen = false"
        @toggle-collapse="isSidebarCollapsed = !isSidebarCollapsed"
      />
      
      <router-view v-slot="{ Component }">
        <component :is="Component" @toggle-sidebar="isSidebarOpen = !isSidebarOpen" />
      </router-view>

      <AuthModal v-if="!isLanding" />
    </div>
  </StyleProvider>
</template>
