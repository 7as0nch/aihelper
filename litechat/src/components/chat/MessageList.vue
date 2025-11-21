<script setup lang="ts">
import { useChatStore } from '../../stores/chat';
import MessageItem from './MessageItem.vue';
import { storeToRefs } from 'pinia';
import { onUpdated, ref, onMounted, onUnmounted } from 'vue';
import { Bot, ArrowDown } from 'lucide-vue-next';

const emit = defineEmits<{
  (e: 'quote', messageId: string, content: string): void;
  (e: 'regenerate', id: string): void;
  (e: 'previewImage', url: string): void;
}>();

const store = useChatStore();
const { messages } = storeToRefs(store);
const bottomRef = ref<HTMLElement | null>(null);
const containerRef = ref<HTMLElement | null>(null);
const showScrollButton = ref(false);

const scrollToBottom = (smooth = true) => {
  bottomRef.value?.scrollIntoView({ behavior: smooth ? 'smooth' : 'auto' });
};

// Auto scroll on new messages if near bottom
onUpdated(() => {
  if (!showScrollButton.value) {
    scrollToBottom();
  }
});

const handleScroll = () => {
  if (!containerRef.value) return;
  
  const { scrollTop, scrollHeight, clientHeight } = containerRef.value;
  // Show button if we are more than 200px away from bottom
  const distanceToBottom = scrollHeight - scrollTop - clientHeight;
  showScrollButton.value = distanceToBottom > 200;
};

onMounted(() => {
  scrollToBottom(false);
  containerRef.value?.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  containerRef.value?.removeEventListener('scroll', handleScroll);
});

const handleQuote = (messageId: string, content: string) => {
  emit('quote', messageId, content);
};

const handleRegenerate = (messageId: string) => {
  emit('regenerate', messageId);
};

const handlePreviewImage = (url: string) => {
  emit('previewImage', url);
};
</script>

<template>
  <div 
    ref="containerRef"
    class="flex-1 min-w-0 overflow-y-auto space-y-6 relative scroll-smooth"
  >
    <!-- Message List -->
    <div class="space-y-6 pb-4 pt-4">
      <MessageItem 
        v-for="msg in messages" 
        :key="msg.id" 
        :message="msg"
        @quote="handleQuote"
        @regenerate="handleRegenerate"
        @preview-image="handlePreviewImage"
      />
      
      <!-- Thinking Indicator -->
      <div v-if="store.isThinking" class="flex gap-4 max-w-3xl mx-auto w-full p-4 animate-fade-in">
        <div class="w-8 h-8 rounded-full bg-orange-100 text-orange-600 dark:bg-orange-900/30 flex items-center justify-center shrink-0">
          <Bot class="w-5 h-5" />
        </div>
        <div class="flex items-center">
          <div class="flex space-x-1">
            <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0s"></div>
            <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
            <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.4s"></div>
          </div>
          <span class="ml-2 text-sm text-gray-500">思考中...</span>
        </div>
      </div>

      <div ref="bottomRef"></div>
    </div>

    <!-- Scroll to Bottom Button -->
    <button
      v-show="showScrollButton"
      @click="scrollToBottom(true)"
      class="fixed bottom-[250px] left-1/2 -translate-x-1/2 p-3 bg-white dark:bg-gray-800 rounded-full shadow-lg border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:text-primary hover:border-primary transition-all duration-300 z-30 animate-fade-in"
      title="回到底部"
    >
      <ArrowDown class="w-5 h-5" />
    </button>
  </div>
</template>
