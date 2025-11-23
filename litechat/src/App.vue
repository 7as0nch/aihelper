<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import Sidebar from '@/components/layout/Sidebar.vue';
import AuthModal from '@/components/auth/AuthModal.vue';

const route = useRoute();
const isSidebarOpen = ref(false);
const isSidebarCollapsed = ref(false);

const isLanding = computed(() => route.name === 'Landing');
</script>

<template>
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

    <AuthModal />
  </div>
</template>
