<script setup lang="ts">
import { Paperclip, X } from 'lucide-vue-next';
import type { Attachment } from '../../../stores/chat';

defineProps<{
  attachments: Attachment[];
}>();

const emit = defineEmits<{
  (e: 'remove', id: string): void;
}>();
</script>

<template>
  <div v-if="attachments.length > 0" class="px-4 pt-4 flex gap-2 flex-wrap">
    <div 
      v-for="att in attachments" 
      :key="att.id"
      class="relative group"
    >
      <div v-if="att.type === 'image'" class="w-16 h-16 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
        <img :src="att.url" class="w-full h-full object-cover" />
      </div>
      <div v-else class="w-16 h-16 rounded-lg bg-gray-100 dark:bg-gray-800 flex items-center justify-center border border-gray-200 dark:border-gray-600">
        <Paperclip class="w-6 h-6 text-gray-400" />
      </div>
      
      <button 
        @click="emit('remove', att.id)"
        class="absolute -top-1 -right-1 bg-gray-900/50 text-white rounded-full p-0.5 opacity-0 group-hover:opacity-100 transition-opacity"
      >
        <X class="w-3 h-3" />
      </button>
    </div>
  </div>
</template>
