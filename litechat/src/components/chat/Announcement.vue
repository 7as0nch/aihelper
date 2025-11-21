<script setup lang="ts">
import { X, Megaphone } from 'lucide-vue-next';

defineProps<{
  content: string;
  type?: 'info' | 'warning' | 'success';
}>();

const emit = defineEmits<{
  (e: 'dismiss'): void;
}>();
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-[#2a2a2a] rounded-xl shadow-2xl max-w-md w-full overflow-hidden border border-gray-100 dark:border-gray-700 transform transition-all scale-100">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-gray-100 dark:border-gray-700 flex items-center justify-between bg-gray-50/50 dark:bg-gray-800/50">
        <div class="flex items-center gap-2">
          <div class="p-1.5 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
            <Megaphone class="w-4 h-4 text-blue-600 dark:text-blue-400" />
          </div>
          <h3 class="font-medium text-gray-900 dark:text-gray-100">系统公告</h3>
        </div>
        <button 
          @click="emit('dismiss')"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>
      
      <!-- Content -->
      <div class="p-6">
        <p class="text-gray-600 dark:text-gray-300 leading-relaxed text-sm">
          {{ content }}
        </p>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 bg-gray-50 dark:bg-gray-800/30 flex justify-end">
        <button 
          @click="emit('dismiss')"
          class="px-4 py-2 bg-primary text-white text-sm font-medium rounded-lg hover:bg-blue-600 transition-colors shadow-sm"
        >
          我知道了
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
