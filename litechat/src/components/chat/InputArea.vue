<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { useChatStore, type Attachment } from '../../stores/chat';
import { 
  Paperclip, 

  Mic, 
  ArrowUp,
  X,
  ChevronDown,
  Zap,
  Brain,
  Sparkles,
  Check,
  Globe
} from 'lucide-vue-next';

import { useAuthStore } from '../../stores/auth';
import { mentionTypes, type MentionType, type MentionOption } from '../../config/mentions';

const props = defineProps<{
  quotedContent?: {id: string, content: string} | null;
}>();

const emit = defineEmits<{
  (e: 'clearQuote'): void;
}>();

const store = useChatStore();
const authStore = useAuthStore();
const input = ref('');
const fileInput = ref<HTMLInputElement | null>(null);
const imageInput = ref<HTMLInputElement | null>(null);
const pendingAttachments = ref<Attachment[]>([]);
const isRecording = ref(false);
const isDropdownOpen = ref(false);
const dropdownRef = ref<HTMLElement | null>(null);

// Mention state
const showMentionMenu = ref(false);
const mentionMenuPosition = ref({ top: 0, left: 0 });
const mentionQuery = ref('');
const activeMentionType = ref<MentionType | null>(null);
const mentionOptions = ref<MentionOption[]>([]);
const mentionCursorIndex = ref(0);
const selectedMentionIndex = ref(0);

const modes = [
  { value: 'smart', label: '智能思考', desc: '智能决策动态搜索', icon: Sparkles, color: 'text-purple-500' },
  { value: 'deep', label: '深度思考', desc: '深入推理给出答案', icon: Brain, color: 'text-blue-500' },
  { value: 'quick', label: '快速回答', desc: '跳过推理直达结果', icon: Zap, color: 'text-orange-500' },
] as const;

const currentModeLabel = computed(() => {
  return modes.find(m => m.value === store.thinkingMode)?.label;
});

const currentModeIcon = computed(() => {
  return modes.find(m => m.value === store.thinkingMode)?.icon;
});

const selectMode = (mode: 'smart' | 'deep' | 'quick') => {
  store.thinkingMode = mode;
  isDropdownOpen.value = false;
};

// Close dropdown when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isDropdownOpen.value = false;
  }
  // Close mention menu if clicking outside
  if (showMentionMenu.value && !(event.target as HTMLElement).closest('.mention-menu')) {
    showMentionMenu.value = false;
    activeMentionType.value = null;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});

const handlePaste = (e: ClipboardEvent) => {
  const items = e.clipboardData?.items;
  if (items) {
    for (const item of items) {
      if (item.type.indexOf('image') !== -1) {
        e.preventDefault();
        const file = item.getAsFile();
        if (file) {
          const attachment: Attachment = {
            id: Date.now().toString(),
            type: 'image',
            name: file.name || 'pasted-image.png',
            url: URL.createObjectURL(file)
          };
          pendingAttachments.value.push(attachment);
        }
        return;
      }
    }
  }
};

const autoResize = (e: Event) => {
  const target = e.target as HTMLTextAreaElement;
  target.style.height = 'auto';
  target.style.height = target.scrollHeight + 'px';
  
  // max-h-48 is 12rem = 192px
  if (target.scrollHeight > 192) {
    target.style.overflowY = 'auto';
  } else {
    target.style.overflowY = 'hidden';
  }
  
  checkMentionTrigger(target);
};

const checkMentionTrigger = (textarea: HTMLTextAreaElement) => {
  const cursorPosition = textarea.selectionStart;
  const textBeforeCursor = textarea.value.substring(0, cursorPosition);
  const lastAtSymbol = textBeforeCursor.lastIndexOf('@');
  
  if (lastAtSymbol !== -1) {
    const textAfterAt = textBeforeCursor.substring(lastAtSymbol + 1);
    // Check if there are spaces, which might indicate we are past the mention
    if (!textAfterAt.includes(' ')) {
      mentionCursorIndex.value = lastAtSymbol;
      mentionQuery.value = textAfterAt;
      showMentionMenu.value = true;
      
      // Calculate position for menu
      // This is a simplified calculation. For a real textarea, you might need a library like get-caret-coordinates
      // For now, we'll position it relative to the textarea bottom-left + some offset based on text length
      // A better approach would be to use a hidden div to mirror text and get coordinates
      mentionMenuPosition.value = {
        top: -10, // Will be adjusted by CSS to be above
        left: 10
      };
      
      // If we are already in a type, filter options
      if (activeMentionType.value && activeMentionType.value.fetchOptions) {
         activeMentionType.value.fetchOptions(mentionQuery.value).then(options => {
           mentionOptions.value = options;
           selectedMentionIndex.value = 0;
         });
      }
      return;
    }
  }
  
  showMentionMenu.value = false;
  activeMentionType.value = null;
};

const handleMentionSelect = async (type: MentionType) => {
  if (type.hasSubMenu && type.fetchOptions) {
    activeMentionType.value = type;
    mentionOptions.value = await type.fetchOptions();
    selectedMentionIndex.value = 0;
    // Clear query so we can type to filter
    // Actually, we might want to keep the trigger text? 
    // For now, let's just show options
  } else {
    insertMention(type.label + ' ');
  }
};

const handleOptionSelect = (option: MentionOption) => {
  insertMention(option.label + ' ');
};

const insertMention = (text: string) => {
  const beforeAt = input.value.substring(0, mentionCursorIndex.value);
  const afterCursor = input.value.substring(mentionCursorIndex.value + mentionQuery.value.length + 1);
  
  input.value = beforeAt + text + afterCursor;
  showMentionMenu.value = false;
  activeMentionType.value = null;
  
  // Reset cursor position and focus
  nextTick(() => {
    const textarea = document.querySelector('textarea');
    if (textarea) {
      textarea.focus();
      const newCursorPos = beforeAt.length + text.length;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
    }
  });
};

const handleSend = async () => {
  if ((!input.value.trim() && pendingAttachments.value.length === 0) || store.isLoading) return;
  
  // Check authentication
  if (!authStore.checkAuth()) return;
  
  const content = input.value;
  const attachments = [...pendingAttachments.value];
  const quote = props.quotedContent ? {
    quoteId: props.quotedContent.id,
    quoteContent: props.quotedContent.content
  } : undefined;
  
  input.value = '';
  pendingAttachments.value = [];
  
  // Clear quote
  if (props.quotedContent) {
    emit('clearQuote');
  }
  
  // Reset height
  const textarea = document.querySelector('textarea');
  if (textarea) textarea.style.height = 'auto';
  
  await store.sendMessage(content, attachments, quote);
};

const triggerFileUpload = () => fileInput.value?.click();


const handleFileChange = (e: Event, type: 'file' | 'image') => {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    const file = target.files[0];
    const attachment: Attachment = {
      id: Date.now().toString(),
      type,
      name: file.name,
      url: URL.createObjectURL(file)
    };
    pendingAttachments.value.push(attachment);
    // Reset input so same file can be selected again if needed
    target.value = '';
  }
};

const removeAttachment = (id: string) => {
  pendingAttachments.value = pendingAttachments.value.filter(a => a.id !== id);
};

const handleVoiceInput = () => {
  if (!('webkitSpeechRecognition' in window)) {
    alert('您的浏览器不支持语音输入');
    return;
  }
  
  if (isRecording.value) {
    isRecording.value = false;
    return;
  }

  isRecording.value = true;
  // Mock voice input for now as we can't easily test real speech api in this env
  setTimeout(() => {
    input.value += ' (语音转文字测试内容) ';
    isRecording.value = false;
    // Trigger resize
    setTimeout(() => {
      const textarea = document.querySelector('textarea');
      if (textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = textarea.scrollHeight + 'px';
      }
    }, 0);
  }, 2000);
};

const handleFocus = (e: FocusEvent) => {
  if (!authStore.checkAuth()) {
    (e.target as HTMLTextAreaElement).blur();
  }
};

const handleKeydown = (e: KeyboardEvent) => {
  if (showMentionMenu.value) {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      const max = activeMentionType.value ? mentionOptions.value.length : mentionTypes.length;
      selectedMentionIndex.value = (selectedMentionIndex.value + 1) % max;
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      const max = activeMentionType.value ? mentionOptions.value.length : mentionTypes.length;
      selectedMentionIndex.value = (selectedMentionIndex.value - 1 + max) % max;
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (activeMentionType.value) {
        handleOptionSelect(mentionOptions.value[selectedMentionIndex.value]);
      } else {
        handleMentionSelect(mentionTypes[selectedMentionIndex.value]);
      }
    } else if (e.key === 'Escape') {
      showMentionMenu.value = false;
      activeMentionType.value = null;
    }
    return;
  }

  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    handleSend();
  }
};
</script>

<template>
  <div class="max-w-3xl mx-auto w-full p-4">
    <div class="relative bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 transition-all duration-300 ease-in-out hover:shadow-lg focus-within:shadow-lg">
      
      <!-- Pending Attachments -->
      <div v-if="pendingAttachments.length > 0" class="px-4 pt-4 flex gap-2 flex-wrap">
        <div 
          v-for="att in pendingAttachments" 
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
            @click="removeAttachment(att.id)"
            class="absolute -top-1 -right-1 bg-gray-900/50 text-white rounded-full p-0.5 opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <X class="w-3 h-3" />
          </button>
        </div>
      </div>

      <!-- Quoted Content Display -->
      <div v-if="quotedContent" class="px-4 pt-3 pb-2 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-start gap-2 p-2 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
          <Quote class="w-4 h-4 text-gray-400 shrink-0 mt-0.5" />
          <div class="flex-1 min-w-0">
            <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ quotedContent.content }}</p>
          </div>
          <button 
            @click="$emit('clearQuote')"
            class="shrink-0 p-0.5 hover:bg-gray-200 dark:hover:bg-gray-700 rounded transition-colors"
          >
            <X class="w-4 h-4 text-gray-400" />
          </button>
        </div>
      </div>

      <!-- Mention Menu -->
      <div 
        v-if="showMentionMenu"
        class="absolute bottom-full left-4 mb-2 w-64 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-20 mention-menu"
      >
        <div class="p-1">
          <template v-if="!activeMentionType">
            <div class="px-2 py-1 text-xs text-gray-400 font-medium">选择类型</div>
            <button 
              v-for="(type, index) in mentionTypes" 
              :key="type.key"
              @click="handleMentionSelect(type)"
              class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
              :class="{ 'bg-blue-50 dark:bg-blue-900/20': index === selectedMentionIndex }"
            >
              <div class="w-6 h-6 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
                <component :is="type.icon" class="w-3 h-3 text-gray-600 dark:text-gray-300" />
              </div>
              <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ type.label }}</span>
            </button>
          </template>
          
          <template v-else>
             <div class="px-2 py-1 text-xs text-gray-400 font-medium flex items-center gap-1">
               <button @click="activeMentionType = null" class="hover:text-primary">返回</button>
               <span>/</span>
               <span>{{ activeMentionType.label }}</span>
             </div>
             <button 
              v-for="(option, index) in mentionOptions" 
              :key="option.id"
              @click="handleOptionSelect(option)"
              class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
              :class="{ 'bg-blue-50 dark:bg-blue-900/20': index === selectedMentionIndex }"
            >
              <div class="flex flex-col">
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ option.label }}</span>
                <span v-if="option.description" class="text-xs text-gray-500">{{ option.description }}</span>
              </div>
            </button>
            <div v-if="mentionOptions.length === 0" class="px-3 py-4 text-center text-gray-500 text-sm">
              无匹配结果
            </div>
          </template>
        </div>
      </div>

      <textarea
        v-model="input"
        rows="1"
        class="w-full bg-transparent border-0 focus:ring-0 focus:outline-none resize-none py-4 pl-4 pr-12 max-h-48 text-gray-900 dark:text-gray-100 placeholder-gray-400 overflow-y-hidden"
        placeholder="输入你的问题，或使用「@快捷引用」对知乎答主、知识库进行提问"
        @input="autoResize"
        @keydown="handleKeydown"
        @paste="handlePaste"
        @focus="handleFocus"
      ></textarea>
      
      <div class="flex items-center justify-between px-2 pb-2">
        <div class="flex items-center gap-1">
          <!-- Smart Thinking Dropdown -->
          <div class="relative" ref="dropdownRef">
            <button 
              @click="isDropdownOpen = !isDropdownOpen"
              class="flex items-center gap-1 px-3 py-1.5 text-sm font-medium text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors"
            >
              <component :is="currentModeIcon" class="w-4 h-4 text-primary" />
              <span>{{ currentModeLabel }}</span>
              <ChevronDown class="w-3 h-3 text-gray-500" />
            </button>

            <!-- Dropdown Menu -->
            <div 
              v-if="isDropdownOpen"
              class="absolute bottom-full left-0 mb-2 w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-10 p-1"
            >
              <button 
                v-for="mode in modes" 
                :key="mode.value"
                @click="selectMode(mode.value)"
                class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
                :class="{ 'bg-blue-50 dark:bg-blue-900/20': store.thinkingMode === mode.value }"
              >
                <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
                  <component :is="mode.icon" class="w-4 h-4" :class="mode.color" />
                </div>
                <div class="flex flex-col">
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ mode.label }}</span>
                  <span class="text-xs text-gray-500">{{ mode.desc }}</span>
                </div>
                <Check v-if="store.thinkingMode === mode.value" class="w-4 h-4 text-primary ml-auto" />
              </button>
            </div>
          </div>
          
          <div class="h-4 w-px bg-gray-200 dark:bg-gray-700 mx-2"></div>

          <button class="flex items-center gap-1 px-2 py-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors text-sm">
            <Globe class="w-4 h-4" />
            <span>知 乎</span>
          </button>
        </div>

        <div class="flex items-center gap-1">
           <input 
            type="file" 
            ref="fileInput" 
            class="hidden" 
            @change="(e) => handleFileChange(e, 'file')" 
          />
          <button @click="triggerFileUpload" class="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors" title="上传文件">
            <Paperclip class="w-5 h-5" />
          </button>
          
          <input 
            type="file" 
            ref="imageInput" 
            class="hidden" 
            accept="image/*"
            @change="(e) => handleFileChange(e, 'image')" 
          />
          
          <button 
            @click="handleVoiceInput" 
            class="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors"
            :class="{ 'text-red-500 bg-red-50 dark:bg-red-900/20': isRecording }"
            title="语音输入"
          >
            <Mic class="w-5 h-5" />
          </button>

          <button 
            class="p-2 rounded-full transition-all duration-200 ease-in-out ml-1"
            :class="[
              (input.trim() || pendingAttachments.length > 0) && !store.isLoading
                ? 'bg-primary text-white hover:bg-blue-600 shadow-md hover:shadow-lg transform hover:-translate-y-0.5' 
                : 'bg-gray-100 text-gray-400 dark:bg-gray-700 dark:text-gray-500 cursor-not-allowed'
            ]"
            :disabled="(!input.trim() && pendingAttachments.length === 0) || store.isLoading"
            @click="handleSend"
          >
            <div v-if="store.isLoading" class="w-5 h-5 border-2 border-gray-400 border-t-transparent rounded-full animate-spin"></div>
            <ArrowUp v-else class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
