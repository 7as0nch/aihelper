<script setup lang="ts">
import { ref } from 'vue';
import MessageList from './MessageList.vue';
import InputArea from './InputArea.vue';
import { Menu, X } from 'lucide-vue-next';
import { useChatStore } from '../../stores/chat';

import { useAuthStore } from '../../stores/auth';

const store = useChatStore();
const authStore = useAuthStore();
const quotedContent = ref<{id: string, content: string} | null>(null);

defineEmits<{
  (e: 'toggleSidebar'): void
}>();

const handleQuote = (messageId: string, content: string) => {
  quotedContent.value = {
    id: messageId,
    content: content
  };
};

const handleRegenerate = (id: string) => {
  if (!authStore.checkAuth()) return;
  
  // Find the message and regenerate from that point
  const messageIndex = store.messages.findIndex(m => m.id === id);
  if (messageIndex !== -1) {
    // Get the last user message before this assistant message
    let userMessageIndex = messageIndex - 1;
    while (userMessageIndex >= 0 && store.messages[userMessageIndex].role !== 'user') {
      userMessageIndex--;
    }
    
    if (userMessageIndex >= 0) {
      // Remove messages from this assistant message onwards
      store.messages.splice(messageIndex);
      // Resend the user message
      const userMessage = store.messages[userMessageIndex];
      store.sendMessage(userMessage.content, userMessage.attachments || []);
    }
  }
};

const previewImageUrl = ref<string | null>(null);

const handlePreviewImage = (url: string) => {
  previewImageUrl.value = url;
};
</script>

<template>
  <div class="flex-1 min-w-0 flex flex-col h-full relative bg-white dark:bg-[#242424]">
    <!-- Mobile Header -->
    <div class="md:hidden h-14 flex items-center px-4 border-b border-gray-100 dark:border-gray-800">
      <button @click="$emit('toggleSidebar')" class="p-2 -ml-2 text-gray-600 dark:text-gray-300">
        <Menu class="w-6 h-6" />
      </button>
      <span class="ml-2 font-medium">新对话</span>
    </div>

    <!-- Welcome State (Centered) -->
    <div v-if="store.messages.length === 0" class="flex-1 flex flex-col items-center justify-center p-4 overflow-y-auto">
      <div class="w-full max-w-3xl space-y-12 -mt-20">
        <h1 class="text-4xl font-bold text-center text-gray-800 dark:text-gray-100 tracking-wide">
          用 <span class="text-primary">提问</span> 发现世界
        </h1>
        
        <InputArea ref="inputAreaRef" :quoted-content="quotedContent" @clear-quote="quotedContent = null" />
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
          <div class="p-4 bg-white dark:bg-[#2a2a2a] rounded-xl border border-gray-100 dark:border-gray-800 hover:shadow-md transition-shadow cursor-pointer text-left group">
            <h3 class="font-medium text-gray-900 dark:text-gray-100 mb-1 group-hover:text-primary transition-colors">想了解点什么？</h3>
            <p class="text-sm text-gray-500">房价下跌对普通人意味着什么？</p>
          </div>
          <div class="p-4 bg-white dark:bg-[#2a2a2a] rounded-xl border border-gray-100 dark:border-gray-800 hover:shadow-md transition-shadow cursor-pointer text-left group">
            <h3 class="font-medium text-gray-900 dark:text-gray-100 mb-1 group-hover:text-primary transition-colors">推荐知识库</h3>
            <p class="text-sm text-gray-500">正经的大学生活与不正经的PC</p>
          </div>
        </div>
      </div>
      
      <div class="fixed bottom-4 text-xs text-gray-400">
        AI 生成的内容可能不准确，请谨慎参考
      </div>
    </div>

    <!-- Standard Chat Layout -->
    <template v-else>
      <MessageList 
        @quote="handleQuote" 
        @regenerate="handleRegenerate"
        @preview-image="handlePreviewImage"
      />
      <div class="pb-6">
        <InputArea ref="inputAreaRef" :quoted-content="quotedContent" @clear-quote="quotedContent = null" />
        <div class="text-center mt-2 text-xs text-gray-400">
          AI 生成的内容可能不准确，请谨慎参考
        </div>
      </div>
    </template>

    <!-- Image Preview Overlay -->
    <div 
      v-if="previewImageUrl" 
      class="fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4 cursor-zoom-out"
      @click="previewImageUrl = null"
    >
      <img 
        :src="previewImageUrl" 
        class="max-w-full max-h-full object-contain rounded-lg shadow-2xl"
        @click.stop
      />
      <button 
        class="absolute top-4 right-4 text-white/70 hover:text-white p-2 rounded-full hover:bg-white/10 transition-colors"
        @click="previewImageUrl = null"
      >
        <X class="w-8 h-8" />
      </button>
    </div>
  </div>
</template>
