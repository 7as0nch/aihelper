<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';
import type { Message, Attachment } from '../../stores/chat';
import type { CallingTool } from '../../api/chat';
import { useChatStore } from '../../stores/chat';
import { User, Bot, Quote, FileText, Wrench, Link as LinkIcon, ChevronDown, ChevronRight, Copy } from 'lucide-vue-next';
import { formatMessageTime } from '../../utils/time';
import MarkdownRenderer from '../common/MarkdownRenderer.vue';
import MessageActions from './MessageActions.vue';

const props = defineProps<{
  message: Message;
  isLastMessage: boolean;
}>();

// 添加一个 ref 来缓存 callingTools
const cachedCallingTools = ref<CallingTool[]>([]);

// 监听 message.callingTools 的变化，更新缓存
watch(() => props.message.callingTools, (newCallingTools) => {
  if (newCallingTools && newCallingTools.length > 0) {
    cachedCallingTools.value = [...newCallingTools];
  }
}, { immediate: true });

// 创建一个计算属性来返回缓存的 callingTools
const displayCallingTools = computed(() => {
  return cachedCallingTools.value.length > 0 ? cachedCallingTools.value : (props.message.callingTools || []);
});

const emit = defineEmits<{
  (e: 'quote', id: string, content: string): void;
  (e: 'regenerate', id: string): void;
  (e: 'image-click', url: string): void;
  (e: 'file-click', file: Attachment): void;
}>();

const chatStore = useChatStore();
const isReasoningCollapsed = ref(false);

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
  
  selectedText.value = text;
  quoteBtnPosition.value = {
    top: rect.top - 40, // Position above selection
    left: rect.left + (rect.width / 2) - 200 // Center horizontally
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

// Close quote button on click outside
const handleClickOutside = (e: MouseEvent) => {
  if (showQuoteBtn.value && !(e.target as HTMLElement).closest('button')) {
    showQuoteBtn.value = false;
  }
};

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
});
</script>

<template>
  <div 
    class="flex gap-4 max-w-3xl mx-auto w-full p-4 animate-fade-in"
    :class="[message.role === 'user' ? 'flex-row-reverse' : 'flex-row']"
  >
    <!-- Avatar -->
    <div 
      class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
      :class="[
        message.role === 'user' 
          ? 'bg-blue-100 text-primary dark:bg-blue-900/30' 
          : 'bg-orange-100 text-orange-600 dark:bg-orange-900/30'
      ]"
    >
      <User v-if="message.role === 'user'" class="w-5 h-5" />
      <Bot v-else class="w-5 h-5" />
    </div>

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
        <!-- Reasoning Content -->
        <div v-if="message.reasoningContent" class="mb-4">
          <div 
            class="bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
          >
            <button 
              @click="isReasoningCollapsed = !isReasoningCollapsed"
              class="w-full flex items-center justify-between px-3 py-2 text-xs font-medium text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded-full bg-purple-100 dark:bg-purple-900/30 text-purple-600 flex items-center justify-center">
                  <span class="text-[10px]">R</span>
                </div>
                <span>深度思考过程</span>
              </div>
              <component :is="isReasoningCollapsed ? ChevronRight : ChevronDown" class="w-4 h-4" />
            </button>
            
            <div v-show="!isReasoningCollapsed" class="px-3 py-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-800/30">
              <div class="prose dark:prose-invert max-w-none text-xs text-gray-600 dark:text-gray-400 leading-relaxed whitespace-pre-wrap font-mono">
                {{ message.reasoningContent }}
              </div>
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
        
        <!-- Calling Tools -->
        <div v-if="displayCallingTools && displayCallingTools.length > 0" class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700/50">
          <div class="text-xs font-medium text-gray-500 mb-2 flex items-center gap-1.5">
            <Wrench class="w-3.5 h-3.5" />
            <span>使用的工具</span>
          </div>
          <div class="flex flex-wrap gap-2">
            <div 
              v-for="(tool, index) in displayCallingTools" 
              :key="index"
              class="flex items-center gap-1.5 px-2.5 py-1.5 bg-gray-50 dark:bg-gray-800/50 border border-gray-100 dark:border-gray-700 rounded-md text-xs text-gray-600 dark:text-gray-300"
            >
              <span class="font-medium">{{ tool.name }}</span>
            </div>
          </div>
        </div>

        <!-- Quote Search Links -->
        <div v-if="message.quoteSearchLinks && message.quoteSearchLinks.length > 0" class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700/50">
          <div class="text-xs font-medium text-gray-500 mb-2 flex items-center gap-1.5">
            <LinkIcon class="w-3.5 h-3.5" />
            <span>引用来源</span>
          </div>
          <div class="flex flex-col gap-1.5">
            <a 
              v-for="(link, index) in message.quoteSearchLinks" 
              :key="index"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-start gap-2 p-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors group border border-transparent hover:border-gray-100 dark:hover:border-gray-700"
            >
              <div class="shrink-0 mt-0.5 text-gray-400 group-hover:text-primary transition-colors">
                <span class="text-[10px] font-mono border border-current rounded px-1 flex items-center justify-center min-w-[18px] h-[18px]">{{ index + 1 }}</span>
              </div>
              <div class="min-w-0">
                <div class="text-xs font-medium text-gray-700 dark:text-gray-200 group-hover:text-primary truncate transition-colors">
                  {{ link.title }}
                </div>
                <div class="text-[10px] text-gray-400 truncate mt-0.5">
                  {{ link.content }}
                </div>
              </div>
            </a>
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
    </div>
  </div>
</template>
