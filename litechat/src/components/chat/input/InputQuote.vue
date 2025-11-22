<script setup lang="ts">
import { Quote, X } from 'lucide-vue-next';

defineProps<{
  content: string;
}>();

const emit = defineEmits<{
  (e: 'clear'): void;
}>();
</script>

<template>
  <div class="px-4 pt-3 pb-2 border-b border-gray-200 dark:border-gray-700">
    <div class="flex items-start gap-2 p-2 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
      <Quote class="w-4 h-4 text-gray-400 shrink-0 mt-0.5" />
      <div class="flex-1 min-w-0">
        <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ content }}</p>
      </div>
      <button 
        @click="emit('clear')"
        class="shrink-0 p-0.5 hover:bg-gray-200 dark:hover:bg-gray-700 rounded transition-colors"
      >
        <X class="w-4 h-4 text-gray-400" />
      </button>
    </div>
  </div>
</template>
