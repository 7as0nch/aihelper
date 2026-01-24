<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';
import type { Message, Attachment, QuoteSearchLink } from '../../stores/chat';
import { useChatStore } from '../../stores/chat';
import { Quote, FileText, ChevronDown, ChevronRight, Copy, Search, Sparkles } from 'lucide-vue-next';
import { formatMessageTime } from '../../utils/time';
import MarkdownRenderer from '../common/MarkdownRenderer.vue';
import MessageActions from './MessageActions.vue';
import ExternalLinkModal from '../common/ExternalLinkModal.vue';

const props = defineProps<{
  message: Message;
  isLastMessage: boolean;
}>();

// 缓存 searchLinks，防止流式输出时被覆盖
const cachedQuoteSearchLinks = ref<QuoteSearchLink[]>([]);

watch(() => props.message.quoteSearchLinks, (newLinks) => {
  if (newLinks && newLinks.length > 0) {
    cachedQuoteSearchLinks.value = [...newLinks];
  }
}, { immediate: true });

const displayQuoteSearchLinks = computed(() => {
  return cachedQuoteSearchLinks.value.length > 0 ? cachedQuoteSearchLinks.value : (props.message.quoteSearchLinks || []);
});

const linkModalVisible = ref(false);
const pendingUrl = ref('');

const openExternalLink = (url: string) => {
  pendingUrl.value = url;
  linkModalVisible.value = true;
};

const confirmLinkJump = () => {
  if (pendingUrl.value) {
    window.open(pendingUrl.value, '_blank', 'noopener noreferrer');
  }
  linkModalVisible.value = false;
  pendingUrl.value = '';
};

const emit = defineEmits<{
  (e: 'quote', id: string, content: string): void;
  (e: 'regenerate', id: string): void;
  (e: 'image-click', url: string): void;
  (e: 'file-click', file: Attachment): void;
}>();

const chatStore = useChatStore();
const isReasoningCollapsed = ref(false);
const isSourcesCollapsed = ref(true);
const reasoningScrollRef = ref<HTMLElement | null>(null);

// Auto-scroll reasoning content while generating
watch(() => props.message.reasoningContent, () => {
  if (props.isLastMessage && chatStore.isLoading && !props.message.content && reasoningScrollRef.value) {
    // 使用 requestAnimationFrame 或更频繁的 nextTick 确保滚动
    const scrollContainer = reasoningScrollRef.value;
    setTimeout(() => {
      scrollContainer.scrollTo({
        top: scrollContainer.scrollHeight,
        behavior: 'smooth'
      });
    }, 10);
  }
}, { immediate: true });

// Auto-collapse reasoning when content starts generating
watch(() => props.message.content, (newContent, oldContent) => {
  if (newContent && !oldContent && props.message.reasoningContent) {
    isReasoningCollapsed.value = true;
  }
});

if (!chatStore.isLoading) {
  isReasoningCollapsed.value = true;
}

const handleQuote = (id: string, content: string) => {
  emit('quote', id, content);
};

const handleRegenerate = (id: string) => {
  emit('regenerate', id);
};

const handleImageClick = (url: string) => {
  emit('image-click', url);
};

const handleFileClick = (file: Attachment) => {
  emit('file-click', file);
};

const copyToClipboard = async (text: string) => {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = text;
      textArea.style.position = 'fixed';
      textArea.style.left = '-9999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    }
    // Optional: Show a toast or feedback
  } catch (err) {
    console.error('Failed to copy:', err);
  }
};

// Partial Quote Logic
const showQuoteBtn = ref(false);
const quoteBtnPosition = ref({ top: 0, left: 0 });
const selectedText = ref('');

const handleSelection = () => {
  const selection = window.getSelection();
  if (!selection || selection.isCollapsed) {
    showQuoteBtn.value = false;
    return;
  }

  const text = selection.toString().trim();
  if (!text) {
    showQuoteBtn.value = false;
    return;
  }
  
  const range = selection.getRangeAt(0);
  const rect = range.getBoundingClientRect();
  
  const btnWidth = 67; // Approximated width of the button
  const padding = 16;
  let left = rect.left + (rect.width / 2) - (btnWidth / 2);
  
  // Boundary detection
  if (left < padding) {
    left = padding;
  } else if (left + btnWidth > window.innerWidth - padding) {
    left = window.innerWidth - btnWidth - padding;
  }

  selectedText.value = text;
  quoteBtnPosition.value = {
    top: rect.top - 40, // Position above selection
    left: left
  };
  showQuoteBtn.value = true;
};

const handleQuoteSelection = () => {
  if (selectedText.value) {
    emit('quote', props.message.id, selectedText.value);
    showQuoteBtn.value = false;
    window.getSelection()?.removeAllRanges();
  }
};

// Close quote button on click outside or scroll
const handleClickOutside = (e: MouseEvent) => {
  if (showQuoteBtn.value && !(e.target as HTMLElement).closest('button')) {
    showQuoteBtn.value = false;
  }
};

const handleScroll = () => {
  if (showQuoteBtn.value) {
    showQuoteBtn.value = false;
  }
};

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside);
  window.addEventListener('scroll', handleScroll, true);
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
  window.removeEventListener('scroll', handleScroll, true);
});
</script>

<template>
  <div 
    class="flex gap-4 max-w-3xl mx-auto w-full px-4 py-2 animate-fade-in"
    :class="[message.role === 'user' ? 'flex-row-reverse' : 'flex-row']"
  >
    <!-- Content -->
    <div 
      class="flex-1 min-w-0 overflow-visible"
      :class="[message.role === 'user' ? 'text-right' : 'text-left']"
    >
      <div 
        class="text-left max-w-full"
        :class="[
          message.role === 'user' 
            ? 'inline-block bg-primary text-white rounded-2xl rounded-tr-sm px-4 py-2 user-message-bubble' 
            : 'block w-full prose dark:prose-invert max-w-none'
        ]"
        :style="message.role === 'user' ? { backgroundColor: '#3b82f6', color: '#ffffff' } : {}"
        @mouseup="message.role === 'assistant' ? handleSelection() : null"
      >
        <!-- Reasoning & Tools (Search Results) -->
        <div v-if="message.reasoningContent || displayQuoteSearchLinks.length > 0 || (isLastMessage && chatStore.isLoading && message.aiModel?.searchByWeb)" class="mb-4 space-y-3">
          <!-- Search Results Header (Clickable Pill) -->
          <div v-if="displayQuoteSearchLinks.length > 0 || (isLastMessage && chatStore.isLoading && message.aiModel?.searchByWeb)" class="flex flex-col gap-3">
            <div class="flex items-center">
              <button 
                @click="isSourcesCollapsed = !isSourcesCollapsed"
                class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 text-xs font-medium text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 transition-all cursor-pointer group"
                :class="{'animate-searching': isLastMessage && chatStore.isLoading && message.aiModel?.searchByWeb}"
              >
                <Search class="w-3.5 h-3.5 text-blue-500" :class="{'animate-pulse-subtle': isLastMessage && chatStore.isLoading && message.aiModel?.searchByWeb}" />
                <span>{{ displayQuoteSearchLinks.length > 0 ? `已检索到 ${displayQuoteSearchLinks.length} 个网页` : '正在检索网页...' }}</span>
                <component v-if="displayQuoteSearchLinks.length > 0" :is="isSourcesCollapsed ? ChevronRight : ChevronDown" class="w-3.5 h-3.5 opacity-50 group-hover:opacity-100 transition-all" />
              </button>
            </div>

            <!-- Expanded Sources List -->
            <div v-show="!isSourcesCollapsed" class="grid grid-cols-1 md:grid-cols-2 gap-3 animate-fade-in pl-1">
              <div 
                v-for="(link, index) in displayQuoteSearchLinks" 
                :key="index"
                @click="openExternalLink(link.url)"
                class="flex flex-col gap-1 p-2.5 rounded-xl bg-white dark:bg-gray-800/30 border border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-800 transition-all group shadow-sm cursor-pointer"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <span class="shrink-0 text-[10px] font-bold text-gray-400 dark:text-gray-500 w-4 h-4 flex items-center justify-center rounded bg-gray-100 dark:bg-gray-700 group-hover:bg-primary group-hover:text-white transition-colors">
                    {{ index + 1 }}
                  </span>
                  <div class="text-xs font-medium text-gray-700 dark:text-gray-200 group-hover:text-primary truncate transition-colors">
                    {{ link.title }}
                  </div>
                </div>
                <div class="text-[11px] text-gray-500 dark:text-gray-400 line-clamp-1 mt-0.5">
                  {{ link.content }}
                </div>
              </div>
            </div>
          </div>

          <!-- Deep Thinking Block -->
          <div v-if="message.reasoningContent" class="group">
            <button 
              @click="isReasoningCollapsed = !isReasoningCollapsed"
              class="flex items-center gap-2 mb-2 text-sm font-semibold text-gray-800 dark:text-gray-200 hover:opacity-80 transition-opacity"
            >
              <div class="flex items-center gap-2">
                <Sparkles class="w-4 h-4 text-purple-500 animate-pulse" v-if="isLastMessage && chatStore.isLoading && !message.content" />
                <span>{{ isLastMessage && chatStore.isLoading && !message.content ? '正在深度思考中...' : '已深度思考' }}</span>
                <span v-if="message.extra?.generateTime && !(isLastMessage && chatStore.isLoading)" class="text-xs font-normal text-gray-400">（用时 {{ message.extra.generateTime }}）</span>
              </div>
              <component :is="isReasoningCollapsed ? ChevronRight : ChevronDown" class="w-4 h-4 text-gray-400" />
            </button>
            
            <div v-show="!isReasoningCollapsed" class="relative group/reasoning">
              <div 
                ref="reasoningScrollRef"
                class="relative pl-4 border-l-2 border-gray-200 dark:border-gray-700 py-1 transition-all duration-500 ease-in-out scroll-smooth"
                :style="isLastMessage && chatStore.isLoading && !message.content ? { maxHeight: '160px', overflowY: 'auto' } : { maxHeight: 'none' }"
                :class="[
                  isLastMessage && chatStore.isLoading && !message.content 
                    ? 'custom-scrollbar-thin' 
                    : ''
                ]"
              >
                <div class="prose dark:prose-invert max-w-none text-[13px] text-gray-500 dark:text-gray-400 leading-relaxed whitespace-pre-wrap italic opacity-85">
                  {{ message.reasoningContent }}
                </div>
              </div>
              
              <!-- Gradient Masks for scrolling -->
              <template v-if="props.isLastMessage && chatStore.isLoading && !message.content">
                <div class="absolute top-0 left-4 right-0 h-8 bg-gradient-to-b from-[#ffffff]/40 dark:from-[#242424]/40 to-transparent pointer-events-none z-10"></div>
                <div class="absolute bottom-0 left-4 right-0 h-8 bg-gradient-to-t from-[#ffffff]/40 dark:from-[#242424]/40 to-transparent pointer-events-none z-10"></div>
              </template>
            </div>
          </div>
        </div>

        <div v-if="message.role === 'user'" class="relative group/user-msg">
          <div class="pr-6">
            {{ message.content }}
          </div>
          <button 
            @click="copyToClipboard(message.content)"
            class="absolute top-0 right-0 p-1 text-white/50 hover:text-white opacity-0 group-hover/user-msg:opacity-100 transition-all"
            title="复制"
          >
            <Copy class="w-3.5 h-3.5" />
          </button>
        </div>
        
        <!-- Markdown Renderer for Assistant -->
        <MarkdownRenderer 
          v-else 
          :content="message.content" 
          :loading="isLastMessage && chatStore.isLoading"
          :quote-search-links="displayQuoteSearchLinks"
        />

        <!-- Quote Display (for user messages with quotes) -->
        <div v-if="message.role === 'user' && message.quoteContent" class="mt-2 pt-2 border-t border-white/20">
          <div class="flex items-start gap-1.5 opacity-80 bg-white/10 p-2 rounded text-sm">
            <Quote class="w-3 h-3 shrink-0 mt-0.5" />
            <p class="line-clamp-2">{{ message.quoteContent }}</p>
          </div>
        </div>

        <!-- Timestamp -->
        <div class="mt-1 text-[10px] opacity-60 text-right">
          {{ formatMessageTime(message.timestamp) }}
        </div>

        <!-- Image Attachments -->
        <div v-if="message.attachments?.some(a => a.type === 'image')" class="flex gap-2 mb-2 overflow-x-auto pb-2 custom-scrollbar">
          <div 
            v-for="img in message.attachments.filter(a => a.type === 'image')" 
            :key="img.id"
            class="relative group cursor-zoom-in shrink-0"
            @click="handleImageClick(img.url || '')"
          >
            <img 
              :src="img.url" 
              class="h-48 w-auto rounded-lg border border-gray-200 dark:border-gray-700 hover:opacity-90 transition-opacity object-cover"
              :alt="img.name"
            />
          </div>
        </div>

        <!-- File Attachments -->
        <div v-if="message.attachments?.some(a => a.type === 'file')" class="flex flex-col gap-2 mb-2">
          <div 
            v-for="file in message.attachments.filter(a => a.type === 'file')" 
            :key="file.id"
            class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors cursor-pointer group"
            @click="handleFileClick(file)"
          >
            <div class="p-2 bg-white dark:bg-gray-700 rounded-lg text-gray-500 group-hover:text-primary transition-colors">
              <FileText class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ file.name }}</div>
              <div class="text-xs text-gray-500">点击下载</div>
            </div>
          </div>
        </div>
        
        <!-- Streaming Cursor -->
        <span 
          v-if="isLastMessage && chatStore.isLoading" 
          class="inline-block w-2 h-4 ml-1 bg-primary align-middle animate-pulse"
        ></span>
      </div>
      
      <!-- Message Actions (for assistant messages only) -->
      <MessageActions 
        v-if="message.role === 'assistant' && !(isLastMessage && chatStore.isLoading)"
        :message-content="message.content"
        :message-id="message.id"
        @quote="handleQuote"
        @regenerate="handleRegenerate"
      />

      <!-- Floating Quote Button -->
      <div 
        v-if="showQuoteBtn"
        class="fixed z-50 animate-fade-in"
        :style="{ top: quoteBtnPosition.top + 'px', left: quoteBtnPosition.left + 'px' }"
      >
        <button 
          @click="handleQuoteSelection"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-gray-900 text-white text-xs font-medium rounded shadow-lg hover:bg-gray-800 transition-colors"
        >
          <Quote class="w-3 h-3" />
          引用
        </button>
      </div>

      <!-- External Link Confirmation Modal -->
      <ExternalLinkModal 
        :visible="linkModalVisible" 
        :url="pendingUrl"
        appName="LiteChat"
        @close="linkModalVisible = false"
        @confirm="confirmLinkJump"
      />
    </div>
  </div>
</template>

<style scoped>
@keyframes gradient-flow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes pulse-subtle {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.85); }
}

.animate-searching {
  background: linear-gradient(-45deg, rgba(243, 244, 246, 0.6), rgba(224, 231, 255, 0.6), rgba(253, 242, 248, 0.6), rgba(224, 231, 255, 0.6));
  background-size: 400% 400%;
  animation: gradient-flow 3s ease infinite;
  border-color: rgba(199, 210, 254, 0.5) !important;
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.05);
}

.dark .animate-searching {
  background: linear-gradient(-45deg, rgba(31, 41, 55, 0.6), rgba(49, 46, 129, 0.6), rgba(49, 46, 129, 0.6), rgba(31, 41, 55, 0.6));
  background-size: 400% 400%;
  border-color: rgba(55, 48, 163, 0.5) !important;
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.1);
}

.animate-pulse-subtle {
  animation: pulse-subtle 2s ease-in-out infinite;
}

.custom-scrollbar-thin::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar-thin::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}

.dark .custom-scrollbar-thin::-webkit-scrollbar-thumb {
  background: #374151;
}
</style>
