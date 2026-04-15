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

const isPublicPage = computed(() => ['Landing', 'page', 'ApplyBeta'].includes(route.name as string));

watch(
  isPublicPage,
  (publicPage) => {
    if (publicPage) {
      authStore.closeModal();
    }
  },
  { immediate: true },
);
</script>

<template>
  <StyleProvider :container="styleContainer" :hash-priority="'high'">
    <div
      :class="[
        isPublicPage
          ? 'min-h-screen bg-transparent'
          : 'flex h-screen w-full overflow-hidden bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-gray-900 dark:to-gray-800',
      ]"
    >
      <Sidebar
        v-if="!isPublicPage"
        :is-open="isSidebarOpen"
        :is-collapsed="isSidebarCollapsed"
        @close="isSidebarOpen = false"
        @toggle-collapse="isSidebarCollapsed = !isSidebarCollapsed"
      />

      <router-view v-slot="{ Component }">
        <component :is="Component" @toggle-sidebar="isSidebarOpen = !isSidebarOpen" />
      </router-view>

      <AuthModal v-if="!isPublicPage" />
    </div>
  </StyleProvider>
</template>