<script setup lang="ts">
import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
import { useChatStore, type Attachment } from '../../stores/chat';
import { Sparkles, Brain, Zap } from 'lucide-vue-next';
import { useAuthStore } from '../../stores/auth';
import { mentionTypes, type MentionType, type MentionOption } from '../../config/mentions';
import type { Dayjs } from 'dayjs';

// Sub-components
import InputAttachments from './input/InputAttachments.vue';
import InputQuote from './input/InputQuote.vue';
import InputMentionMenu from './input/InputMentionMenu.vue';
import InputTextarea from './input/InputTextarea.vue';
import InputToolbarMobile from './input/InputToolbarMobile.vue';
import InputToolbarDesktop from './input/InputToolbarDesktop.vue';

const props = defineProps<{
  quotedContent?: {id: string, content: string} | null;
}>();

const emit = defineEmits<{
  (e: 'clearQuote'): void;
  (e: 'toggleScreenshot'): void;
}>();

const store = useChatStore();
const authStore = useAuthStore();
const input = ref('');
const pendingAttachments = ref<Attachment[]>([]);
const isRecording = ref(false);

// Refs to sub-components
const textareaComponentRef = ref<InstanceType<typeof InputTextarea> | null>(null);
const mobileToolbarRef = ref<InstanceType<typeof InputToolbarMobile> | null>(null);
const desktopToolbarRef = ref<InstanceType<typeof InputToolbarDesktop> | null>(null);

// Mention state
const showMentionMenu = ref(false);
const mentionMenuPosition = ref({ top: 0, left: 0 });
const mentionQuery = ref('');
const activeMentionType = ref<MentionType | null>(null);
const mentionOptions = ref<MentionOption[]>([]);
const mentionCursorIndex = ref(0);
const selectedMentionIndex = ref(0);

// Date Range Picker State
const showDatePicker = ref(false);
const pickerMode = ref<'date' | 'time' | 'datetime'>('date');
const dateRangeValue = ref<[Dayjs, Dayjs] | undefined>(undefined);

const activeButton = ref<string | null>(null);

const triggerBounce = (key: string) => {
  activeButton.value = key;
  setTimeout(() => {
    activeButton.value = null;
  }, 300);
};

const getDateFormat = computed(() => {
  if (pickerMode.value === 'time') return 'HH:mm:ss';
  if (pickerMode.value === 'datetime') return 'YYYY-MM-DD HH:mm:ss';
  return 'YYYY-MM-DD';
});

const confirmDateRange = () => {
  if (dateRangeValue.value) {
    const start = dateRangeValue.value[0].format(getDateFormat.value);
    const end = dateRangeValue.value[1].format(getDateFormat.value);
    insertMention(` [ ${start} 至 ${end} ] `);
    showDatePicker.value = false;
    dateRangeValue.value = undefined;
  }
};

const modes = [
  { value: 'smart', label: '智能思考', desc: '智能决策动态搜索', icon: Sparkles, color: 'text-purple-500' },
  { value: 'deep', label: '深度思考', desc: '深入推理给出答案', icon: Brain, color: 'text-blue-500' },
  { value: 'quick', label: '快速回答', desc: '跳过推理直达结果', icon: Zap, color: 'text-orange-500' },
] as const;

const currentModeLabel = computed(() => {
  return modes.find(m => m.value === store.thinkingMode)?.label || '智能思考';
});

const currentModeIcon = computed(() => {
  return modes.find(m => m.value === store.thinkingMode)?.icon || Sparkles;
});

const selectMode = (mode: string) => {
  store.thinkingMode = mode as 'smart' | 'deep' | 'quick';
  mobileToolbarRef.value?.closeDropdown();
  desktopToolbarRef.value?.closeDropdown();
};

// Close dropdown when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  // Close mention menu if clicking outside
  if (showMentionMenu.value && 
      !(event.target as HTMLElement).closest('.mention-menu') && 
      !(event.target as HTMLElement).closest('.ant-picker-dropdown')
  ) {
    showMentionMenu.value = false;
    activeMentionType.value = null;
  }
  
  // Dropdowns in toolbars handle their own outside clicks usually via v-if logic or we can delegate
  // But here we rely on the toolbars to handle their own state or expose methods.
  // The original code had global click listeners.
  // For simplicity, let's assume toolbars handle their own "click outside" if they use a global listener or similar mechanism.
  // Actually, the original code used refs to check containment.
  // Since we moved logic to sub-components, we rely on them or a global click handler here if we pass refs up?
  // The sub-components use a simple v-if toggle. A global click handler is better for "click outside".
  // Let's keep a global click handler here and call close methods on refs.
  
  // const target = event.target as Node;
  // We can't easily check refs inside components from here without exposing the root element.
  // Instead, let's make the sub-components handle "click outside" if possible, or just close them on any click that isn't inside them.
  // For now, let's just close them on any click on document, and stopPropagation in the component itself.
  
  mobileToolbarRef.value?.closeDropdown();
  desktopToolbarRef.value?.closeDropdown();
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
  } else {
    insertMention(type.label + ' ');
  }
};

const handleOptionSelect = (option: MentionOption) => {
  if (option.id === 'custom_range') {
    showDatePicker.value = true;
    return;
  }
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
    const textarea = textareaComponentRef.value?.textarea;
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
  const textarea = textareaComponentRef.value?.textarea;
  if (textarea) textarea.style.height = 'auto';
  
  await store.sendMessage(content, attachments, quote);
};

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
      const textarea = textareaComponentRef.value?.textarea;
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

  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault();
    handleSend();
  }
};

const handleScreenshot = () => {
  emit('toggleScreenshot');
};
</script>

<template>
  <div class="max-w-3xl mx-auto w-full p-2">
    <!-- Mobile Floating Toolbar -->
    <InputToolbarMobile 
      ref="mobileToolbarRef"
      :current-mode="store.thinkingMode"
      :modes="modes"
      :active-button="activeButton"
      :current-mode-label="currentModeLabel"
      :current-mode-icon="currentModeIcon"
      @select-mode="selectMode"
      @trigger-bounce="triggerBounce"
      @screenshot="handleScreenshot"
      @upload="desktopToolbarRef?.triggerFileUpload" 
    />
    <!-- Note: upload on mobile might need a separate hidden input or reuse desktop's if visible/accessible, 
         but desktop toolbar is hidden on mobile. 
         Let's add hidden inputs here or in the mobile toolbar if needed. 
         Actually, InputToolbarMobile emits 'upload', we need to handle it.
         We can add a hidden file input in this component or reuse the one in DesktopToolbar if it was mounted (but it's v-if hidden?).
         Wait, DesktopToolbar is hidden with CSS (hidden md:flex), so it IS mounted. We can access it.
    -->

    <div class="relative bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 transition-all duration-300 ease-in-out hover:shadow-lg focus-within:shadow-lg">
      
      <!-- Pending Attachments -->
      <InputAttachments 
        :attachments="pendingAttachments"
        @remove="removeAttachment"
      />

      <!-- Quoted Content Display -->
      <InputQuote 
        v-if="quotedContent"
        :content="quotedContent.content"
        @clear="emit('clearQuote')"
      />

      <!-- Mention Menu -->
      <InputMentionMenu 
        :show="showMentionMenu"
        :position="mentionMenuPosition"
        :options="mentionOptions"
        :mention-types="mentionTypes"
        :active-index="selectedMentionIndex"
        :active-type="activeMentionType"
        :show-date-picker="showDatePicker"
        v-model:date-range-value="dateRangeValue"
        v-model:picker-mode="pickerMode"
        :get-date-format="getDateFormat"
        @select-type="handleMentionSelect"
        @select-option="handleOptionSelect"
        @close="showMentionMenu = false"
        @confirm-date="confirmDateRange"
        @update:active-type="activeMentionType = $event"
        @update:show-date-picker="showDatePicker = $event"
      />

      <!-- Textarea -->
      <InputTextarea 
        ref="textareaComponentRef"
        v-model="input"
        :is-loading="store.isLoading"
        :is-recording="isRecording"
        :has-attachments="pendingAttachments.length > 0"
        @send="handleSend"
        @stop="store.stopGeneration()"
        @keydown="handleKeydown"
        @paste="handlePaste"
        @focus="handleFocus"
        @resize="autoResize"
        @voice-input="handleVoiceInput"
      />
      
      <!-- Desktop Toolbar -->
      <InputToolbarDesktop 
        ref="desktopToolbarRef"
        :current-mode="store.thinkingMode"
        :modes="modes"
        :current-mode-label="currentModeLabel"
        :current-mode-icon="currentModeIcon"
        :input="input"
        :is-loading="store.isLoading"
        :has-attachments="pendingAttachments.length > 0"
        :is-recording="isRecording"
        @select-mode="selectMode"
        @screenshot="handleScreenshot"
        @voice-input="handleVoiceInput"
        @send="handleSend"
        @stop="store.stopGeneration()"
        @file-change="handleFileChange"
      />
    </div>
  </div>
</template>

<style scoped>
/* Animations moved to src/style.css */
</style>
