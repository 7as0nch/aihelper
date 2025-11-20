<script setup lang="ts">
import { ref } from 'vue';
import Sidebar from '@/components/layout/Sidebar.vue';
import AuthModal from '@/components/auth/AuthModal.vue';

const isSidebarOpen = ref(false);
const isSidebarCollapsed = ref(false);
</script>

<template>
  <div class="flex h-screen w-full overflow-hidden bg-white dark:bg-[#242424]">
    <Sidebar 
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
