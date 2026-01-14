<script setup lang="ts">
import { useChatStore } from '../../stores/chat';
import MessageItem from './MessageItem.vue';
import { storeToRefs } from 'pinia';
import { onUpdated, ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { Bot, ArrowDown, Check } from 'lucide-vue-next';

const emit = defineEmits<{
  (e: 'quote', messageId: string, content: string): void;
  (e: 'regenerate', id: string): void;
  (e: 'previewImage', url: string): void;
  (e: 'toggleSelect', id: string): void;
}>();

defineProps<{
  isScreenshotMode?: boolean;
  selectedIds?: Set<string>;
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
  showScrollButton.value = distanceToBottom > 50;
};

onMounted(() => {
  scrollToBottom(false);
  containerRef.value?.addEventListener('scroll', handleScroll);
});

// Watch for messages change to auto-scroll
watch(() => store.messages, async () => {
  await nextTick();
  if (!showScrollButton.value) {
    scrollToBottom(true);
  }
}, { deep: true });

// Watch for specific message content updates (streaming)
watch(() => store.messages[store.messages.length - 1]?.content, async () => {
  if (!showScrollButton.value) {
     bottomRef.value?.scrollIntoView({ behavior: 'auto' });
  }
});

// Also watch reasoning content
watch(() => store.messages[store.messages.length - 1]?.reasoningContent, async () => {
  if (!showScrollButton.value) {
     bottomRef.value?.scrollIntoView({ behavior: 'auto' });
  }
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

defineExpose({
  containerRef
});
</script>

<template>
  <div 
    ref="containerRef"
    id="message-list-container"
    class="flex-1 min-w-0 overflow-y-auto space-y-6 relative scroll-smooth message-list-container custom-scrollbar message-list-mask"
  >
    <!-- Message List -->
    <div class="space-y-6 pb-48 pt-4">
      <MessageItem 
        v-for="(msg, index) in messages" 
        :key="msg.id" 
        :message="msg"
        :is-last-message="index === messages.length - 1"
        @quote="handleQuote"
        @regenerate="handleRegenerate"
        @preview-image="handlePreviewImage"
        :class="{ 
          'opacity-50': isScreenshotMode && selectedIds && !selectedIds.has(msg.id), 
          'cursor-pointer': isScreenshotMode,
          'ring-2 ring-primary ring-offset-2 ring-offset-white dark:ring-offset-[#1a1a1a] rounded-lg': isScreenshotMode && selectedIds?.has(msg.id)
        }"
        @click="isScreenshotMode && $emit('toggleSelect', msg.id)"
        :data-message-id="msg.id"
        class="message-item relative transition-all duration-200"
      >
        <template #prefix v-if="isScreenshotMode">
           <div class="absolute left-0 top-1/2 -translate-y-1/2 -ml-8 selection-checkbox">
             <div 
               class="w-5 h-5 rounded-full border-2 flex items-center justify-center transition-colors"
               :class="selectedIds?.has(msg.id) ? 'bg-primary border-primary' : 'border-gray-300 dark:border-gray-600'"
             >
               <Check v-if="selectedIds?.has(msg.id)" class="w-3 h-3 text-white" />
             </div>
           </div>
        </template>
      </MessageItem>
      
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
      class="fixed bottom-[200px] left-1/2 -translate-x-1/2 p-3 bg-white dark:bg-gray-800 rounded-full shadow-lg border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-300 hover:text-primary hover:border-primary transition-all duration-300 z-[99999] animate-fade-in"
      title="回到底部"
    >
      <ArrowDown class="w-5 h-5" />
    </button>
  </div>
</template>

<style scoped>
/* 消息列表特有的遮罩效果，滚动条样式已移至全局 style.css */
</style>
