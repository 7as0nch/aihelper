<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import MessageList from './MessageList.vue';
import InputArea from './InputArea.vue';
import { Menu, X, BookOpen, MessageCircle, Star, Share2, Plus, ChevronLeft } from 'lucide-vue-next';
import { useChatStore } from '../../stores/chat';
import { useAuthStore } from '../../stores/auth';
import { useRecommendationStore } from '../../stores/recommendation';

const store = useChatStore();
const authStore = useAuthStore();
const recommendationStore = useRecommendationStore();
const quotedContent = ref<{id: string, content: string} | null>(null);

// ... (keep existing emits and handlers)

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

// Recommendations
onMounted(() => {
  recommendationStore.fetchRecommendations();
});

const randomQuestions = computed(() => {
  return recommendationStore.getRandomItems(recommendationStore.questionList, 3);
});

const randomKnowledgeBase = computed(() => {
  const items = recommendationStore.getRandomItems(recommendationStore.knowledgeBaseList, 1);
  return items.length > 0 ? items[0] : null;
});

const handleQuestionClick = (content: string) => {
    if (!authStore.checkAuth()) return;
    store.sendMessage(content);
};
</script>

<template>
  <div class="flex-1 min-w-0 flex flex-col h-full relative bg-white/70 dark:bg-[#242424]/70 backdrop-blur-md transition-colors">
    <!-- Mobile Header -->
    <div class="md:hidden h-14 flex items-center justify-between px-4 border-b border-gray-100 dark:border-gray-800">
      <div class="flex items-center gap-2">
        <button @click="$emit('toggleSidebar')" class="p-2 -ml-2 text-gray-600 dark:text-gray-300">
          <Menu class="w-6 h-6" />
        </button>
        <span class="font-medium text-gray-900 dark:text-white">{{ store.historyItems.find(h => h.id === store.currentChatId)?.title || '新对话' }}</span>
      </div>
      <div class="flex items-center gap-1">
        <button 
          @click="$router.push('/')"
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-500 dark:text-gray-400"
          title="新对话"
        >
          <Plus class="w-5 h-5" />
        </button>
        <button 
          v-if="store.currentChatId"
          @click="store.toggleFavorite(store.currentChatId)"
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
          :class="store.isFavorite(store.currentChatId) ? 'text-yellow-500' : 'text-gray-400'"
          title="收藏"
        >
          <Star class="w-5 h-5" :fill="store.isFavorite(store.currentChatId) ? 'currentColor' : 'none'" />
        </button>
        <button 
          class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-500 dark:text-gray-400"
          title="分享"
        >
          <Share2 class="w-5 h-5" />
        </button>
      </div>
    </div>

    <!-- Welcome State (Centered) -->
    <div v-if="store.messages.length === 0" class="flex-1 flex flex-col items-center justify-center p-4 overflow-y-auto">
      <div class="w-full max-w-3xl space-y-12 -mt-20">
        <h1 class="text-4xl font-bold text-center text-gray-800 dark:text-gray-100 tracking-wide">
          用 <span class="text-primary">提问</span> 发现世界
        </h1>
        
        <InputArea ref="inputAreaRef" :quoted-content="quotedContent" @clear-quote="quotedContent = null" />
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
          <!-- Random Question -->
          <!-- Random Questions -->
          <div v-if="recommendationStore.showQuestions && randomQuestions.length > 0" class="flex flex-col gap-4">
            <h3 class="font-medium text-lg text-gray-900 dark:text-gray-100 mb-2">想了解点什么？</h3>
            <div class="flex flex-col gap-3 items-start">
              <button
                v-for="(question, index) in randomQuestions"
                :key="question.id"
                @click="handleQuestionClick(question.content)"
                class="px-5 py-3 bg-white dark:bg-[#2a2a2a] rounded-full border border-gray-200 dark:border-gray-700 hover:border-primary dark:hover:border-primary hover:text-primary dark:hover:text-primary transition-all shadow-sm hover:shadow-md text-left text-sm text-gray-600 dark:text-gray-300 animate-float"
                :style="{ animationDelay: `${index * 0.1}s` }"
              >
                {{ question.content }}
              </button>
            </div>
          </div>

          <!-- Random Knowledge Base -->
          <div 
            v-if="recommendationStore.showKnowledgeBase && randomKnowledgeBase"
            class="p-4 bg-white dark:bg-[#2a2a2a] rounded-xl border border-gray-100 dark:border-gray-800 hover:shadow-md transition-shadow cursor-pointer text-left group"
          >
             <div class="flex items-center gap-2 mb-2">
                <BookOpen class="w-4 h-4 text-primary" />
                <h3 class="font-medium text-gray-900 dark:text-gray-100 group-hover:text-primary transition-colors">推荐知识库</h3>
            </div>
            <p class="text-sm text-gray-500 line-clamp-2">{{ randomKnowledgeBase.title }}</p>
          </div>
        </div>
      </div>
      
      <div class="fixed bottom-4 text-xs text-gray-400">
        AI 生成的内容可能不准确，请谨慎参考
      </div>
    </div>

    <!-- Standard Chat Layout -->
    <template v-else>
      <!-- Desktop Header -->
      <div class="hidden md:flex items-center justify-between px-6 py-3 border-b border-gray-100 dark:border-gray-800">
        <div class="flex items-center gap-3">
          <button 
            @click="$router.push('/')"
            class="p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            title="返回"
          >
            <ChevronLeft class="w-5 h-5" />
          </button>
          <div>
            <h2 class="text-base font-medium text-gray-900 dark:text-white">
              {{ store.historyItems.find(h => h.id === store.currentChatId)?.title || '新对话' }}
            </h2>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ store.messages.length }} 条对话
            </div>
          </div>
        </div>
        <div class="flex items-center gap-1">
          <button 
            @click="$router.push('/')"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-500 dark:text-gray-400"
            title="新对话"
          >
            <Plus class="w-5 h-5" />
          </button>
          <button 
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-500 dark:text-gray-400"
            title="评论"
          >
            <MessageCircle class="w-5 h-5" />
          </button>
          <button 
            v-if="store.currentChatId"
            @click="store.toggleFavorite(store.currentChatId)"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            :class="store.isFavorite(store.currentChatId) ? 'text-yellow-500' : 'text-gray-400'"
            title="收藏"
          >
            <Star class="w-5 h-5" :fill="store.isFavorite(store.currentChatId) ? 'currentColor' : 'none'" />
          </button>
          <button 
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-500 dark:text-gray-400"
            title="分享"
          >
            <Share2 class="w-5 h-5" />
          </button>
        </div>
      </div>

      <MessageList 
        @quote="handleQuote" 
        @regenerate="handleRegenerate"
        @preview-image="handlePreviewImage"
      />
      <div class="pb-6">
        <!-- Persistent Suggestions -->
        <div v-if="recommendationStore.showQuestions" class="px-4 mb-2 overflow-x-auto no-scrollbar relative z-10">
          <div class="flex gap-2 w-max">
            <button
              v-for="question in recommendationStore.questionList"
              :key="question.id"
              @click="handleQuestionClick(question.content)"
              class="px-4 py-2 bg-white/50 dark:bg-[#2a2a2a]/50 backdrop-blur-sm border border-gray-200 dark:border-gray-700 rounded-full text-sm text-gray-600 dark:text-gray-300 hover:bg-white dark:hover:bg-[#2a2a2a] hover:border-primary dark:hover:border-primary hover:text-primary dark:hover:text-primary transition-all whitespace-nowrap shadow-sm"
            >
              {{ question.content }}
            </button>
          </div>
        </div>

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

<style scoped>
@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-5px); }
  100% { transform: translateY(0px); }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
