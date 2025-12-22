```vue
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import MessageList from './MessageList.vue';
import InputArea from './InputArea.vue';
import ScreenshotManager from './ScreenshotManager.vue';
import WelcomeMobile from './welcome/WelcomeMobile.vue';
import WelcomeDesktop from './welcome/WelcomeDesktop.vue';
import { Menu, X, MessageCircle, Star, Share2, Plus, ChevronLeft } from 'lucide-vue-next';
import { useChatStore } from '../../stores/chat';
import { useAuthStore } from '../../stores/auth';
import { useRecommendationStore } from '../../stores/recommendation';
import { getConfig } from '../../config';

const store = useChatStore();
const authStore = useAuthStore();
const recommendationStore = useRecommendationStore();
const quotedContent = ref<{id: string, content: string} | null>(null);

const props = defineProps<{
  isScreenshotMode: boolean;
  selectedMessageIds: Set<string>;
}>();

const emit = defineEmits<{
  (e: 'toggleSidebar'): void;
  (e: 'update:isScreenshotMode', value: boolean): void;
  (e: 'toggleSelect', id: string): void;
}>();

const messageListRef = ref<InstanceType<typeof MessageList> | null>(null);

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

const toggleScreenshotMode = () => {
  emit('update:isScreenshotMode', !props.isScreenshotMode);
  // Clearing selection is handled by the parent or ScreenshotManager when mode changes
};

// const handleQuestionClick = (content: string) => {
//   if (!authStore.checkAuth()) return;
//   store.sendMessage(content);
// };

// Announcement Logic
import Announcement from './Announcement.vue';
import { fetchAnnouncement, type Announcement as AnnouncementType } from '../../api/announcement';
import { useStorage } from '@vueuse/core';

const announcement = ref<AnnouncementType | null>(null);
const dismissedAnnouncements = useStorage<string[]>('dismissed-announcements', []);

const isAnnouncementDismissed = computed(() => {
  return announcement.value ? dismissedAnnouncements.value.includes(announcement.value.id) : true;
});

const dismissAnnouncement = () => {
  if (announcement.value) {
    dismissedAnnouncements.value.push(announcement.value.id);
  }
};

onMounted(async () => {
  recommendationStore.fetchRecommendations();
  
  // Fetch announcement
  const data = await fetchAnnouncement();
  if (data && !dismissedAnnouncements.value.includes(data.id)) {
    announcement.value = data;
  }
});
</script>

<template>
  <div class="flex-1 min-w-0 flex flex-col h-full relative bg-white/70 dark:bg-[#242424]/70 backdrop-blur-md transition-colors">
    <!-- Mobile Header -->
    <div class="md:hidden h-14 flex items-center justify-between px-4 border-b border-gray-100 dark:border-gray-800">
      <div class="flex items-center gap-2">
        <button @click="emit('toggleSidebar')" class="p-2 -ml-2 text-gray-600 dark:text-gray-300">
          <Menu class="w-6 h-6" />
        </button>
        <span class="font-medium text-gray-900 dark:text-white">{{ store.historyItems.find(h => h.id === store.currentChatId)?.title || '新对话' }}</span>
      </div>
      <div class="flex items-center gap-1">
        <button 
          @click="$router.push('/chat')"
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
    <div v-if="store.messages.length === 0" class="flex-1 h-full relative">
      <WelcomeMobile 
        class="md:hidden"
        :quoted-content="quotedContent"
        @clear-quote="quotedContent = null"
        @toggle-screenshot="toggleScreenshotMode"
      />
      <WelcomeDesktop
        class="hidden md:flex"
        :quoted-content="quotedContent"
        @clear-quote="quotedContent = null"
        @toggle-screenshot="toggleScreenshotMode"
      />
    </div>

    <!-- Standard Chat Layout -->
    <template v-else>
      <!-- Desktop Header -->
      <div 
        v-if="getConfig('VITE_SHOW_HEADER') !== 'false'"
        class="hidden md:flex items-center justify-between px-6 py-3 border-b border-gray-100 dark:border-gray-800"
      >
        <div class="flex items-center gap-3">
          <button 
            @click="$router.push('/chat')"
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
            @click="$router.push('/chat')"
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
        ref="messageListRef"
        :messages="store.messages" 
        :is-thinking="store.isThinking"
        :is-screenshot-mode="props.isScreenshotMode"
        :selected-ids="props.selectedMessageIds"

        @quote="handleQuote" 
        @regenerate="handleRegenerate"
        @preview-image="handlePreviewImage"
        @toggle-select="emit('toggleSelect', $event)"
        class="flex-1 overflow-y-auto"
      />
      <ScreenshotManager 
        :is-screenshot-mode="isScreenshotMode"
        @update:is-screenshot-mode="$emit('update:isScreenshotMode', $event)"
        :selected-ids="selectedMessageIds"
        :target-element="messageListRef?.containerRef || null"
        @clear-selection="selectedMessageIds.clear()"
      />
      <div class="relative z-20">
        <div class="bg-white/70 dark:bg-[#242424]/70 backdrop-blur-md pb-4">
          <InputArea 
            ref="inputAreaRef" 
            :quoted-content="quotedContent" 
            @clear-quote="quotedContent = null" 
            @toggle-screenshot="toggleScreenshotMode"
          />
          <div 
            v-if="getConfig('VITE_SHOW_FOOTER') !== 'false'"
            class="text-center mt-2 text-xs text-gray-400"
          >
            {{ getConfig('VITE_FOOTER_TEXT') }}
          </div>
        </div>
      </div>
    </template>

    <!-- Announcement Modal -->
    <Announcement 
      v-if="announcement && !isAnnouncementDismissed"
      :content="announcement.content"
      :type="announcement.type"
      @dismiss="dismissAnnouncement"
    />

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
.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
