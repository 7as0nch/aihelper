<script setup lang="ts">
import { computed } from 'vue';
import { BookOpen } from 'lucide-vue-next';
import InputArea from '../InputArea.vue';
import { useChatStore } from '../../../stores/chat';
import { useAuthStore } from '../../../stores/auth';
import { useRecommendationStore } from '../../../stores/recommendation';

const store = useChatStore();
const authStore = useAuthStore();
const recommendationStore = useRecommendationStore();

const props = defineProps<{
  quotedContent: { id: string; content: string } | null;
}>();

const emit = defineEmits<{
  (e: 'clearQuote'): void;
  (e: 'toggleScreenshot'): void;
}>();

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
  <div class="flex-1 flex flex-col h-full p-4 overflow-y-auto">
    <!-- PC Top Spacer -->
    <div class="flex-1">
      <!-- Title -->
      <h1 class="text-4xl font-bold text-center text-gray-800 dark:text-gray-100 tracking-wide mt-40">
      你的 <span class="text-primary">智慧</span> 帮手
      </h1>
    </div>

    <!-- Main Content Area -->
    <div class="w-full max-w-3xl mx-auto block h-auto justify-center">
      
      <!-- Input Area -->
      <div class="w-full mb-0">
        <InputArea 
          ref="inputAreaRef" 
          :quoted-content="props.quotedContent" 
          @clear-quote="emit('clearQuote')" 
          @toggle-screenshot="emit('toggleScreenshot')"
        />
      </div>
    </div>

    <!-- PC Bottom Spacer (contains Recommendations) -->
    <div class="flex-1 w-full max-w-3xl mx-auto relative">
        <div class="absolute top-8 left-0 w-full flex flex-col gap-4 animate-fade-in-up">
           <!-- Random Questions -->
          <div v-if="recommendationStore.showQuestions && randomQuestions.length > 0" class="flex flex-col gap-4 w-full">
            <h3 class="font-medium text-lg text-gray-900 dark:text-gray-100 opacity-80 text-center" :class="{ 'md:text-left': recommendationStore.showKnowledgeBase }">想了解点什么？</h3>
            <div class="flex flex-wrap gap-3 w-full" :class="recommendationStore.showKnowledgeBase ? 'justify-start' : 'justify-center'">
              <button
                v-for="(question, index) in randomQuestions"
                v-tracker="{ type: 'click', name: 'randomQuestion', data: { question: question.content } }"
                :key="question.id"
                @click="handleQuestionClick(question.content)"
                class="px-5 py-3 bg-white dark:bg-[#2a2a2a] rounded-2xl border border-gray-200 dark:border-gray-700 hover:border-primary dark:hover:border-primary hover:text-primary dark:hover:text-primary transition-all shadow-sm hover:shadow-md text-sm text-gray-600 dark:text-gray-300 animate-float"
                :style="{ animationDelay: `${index * 0.1}s` }"
              >
                {{ question.content }}
              </button>
            </div>
          </div>

          <!-- Random Knowledge Base -->
          <div 
            v-if="recommendationStore.showKnowledgeBase && randomKnowledgeBase"
            class="mt-4 p-4 bg-white dark:bg-[#2a2a2a] rounded-xl border border-gray-100 dark:border-gray-800 hover:shadow-md transition-shadow cursor-pointer text-left group max-w-md w-full"
          >
             <div class="flex items-center gap-2 mb-2">
                <BookOpen class="w-4 h-4 text-primary" />
                <h3 class="font-medium text-gray-900 dark:text-gray-100 group-hover:text-primary transition-colors">推荐知识库</h3>
            </div>
            <p class="text-sm text-gray-500 line-clamp-2">{{ randomKnowledgeBase.title }}</p>
          </div>
        </div>
    </div>
    
    <div class="absolute bottom-4 left-0 right-0 text-center text-xs text-gray-400">
      AI 生成的内容可能不准确，请谨慎参考
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
</style>
