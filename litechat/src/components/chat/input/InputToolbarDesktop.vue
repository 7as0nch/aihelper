<script setup lang="ts">
import { ref } from 'vue';
import { ChevronDown, Globe, Camera, Paperclip, ArrowUp, Square, Check, MoreHorizontal, Sparkles, Image as ImageIcon } from 'lucide-vue-next';
import VoiceInput from '../../common/VoiceInput.vue';

defineProps<{
  currentMode: string;
  modes: readonly any[];
  currentModeLabel: string;
  currentModeIcon: any;
  input: string;
  isLoading: boolean;
  hasAttachments: boolean;
  extButtons?: Array<{ id: string; name: string; api: string; desc: string; icon?: any }>;
  searchByWeb?: boolean;
}>();

const emit = defineEmits<{
  (e: 'select-mode', mode: string): void;
  (e: 'screenshot'): void;
  (e: 'upload'): void;
  (e: 'voice-start'): void;
  (e: 'voice-stop'): void;
  (e: 'voice-result', text: string): void;
  (e: 'voice-interim', text: string): void;
  (e: 'voice-error', message: string): void;
  (e: 'send'): void;
  (e: 'stop'): void;
  (e: 'file-change', event: Event, type: 'file' | 'image'): void;
  (e: 'ext-action', btn: any): void;
  (e: 'toggle-web-search'): void;
}>();

const isDropdownOpen = ref(false);
const isMoreDropdownOpen = ref(false);

const fileInput = ref<HTMLInputElement | null>(null);
const imageInput = ref<HTMLInputElement | null>(null);

const triggerFileUpload = () => {
  fileInput.value?.click();
};

const closeDropdown = () => {
  isDropdownOpen.value = false;
  isMoreDropdownOpen.value = false;
};

// Icon mapping for external buttons
const iconMap: Record<string, any> = {
  'Sparkles': Sparkles,
  'Image': ImageIcon
};

const getIcon = (iconName: string) => {
  return iconMap[iconName] || Sparkles;
};

defineExpose({
  closeDropdown,
  triggerFileUpload
});
</script>

<template>
  <div class="hidden md:flex items-center justify-between px-2 pb-2">
    <div class="flex items-center gap-1">
      <!-- Smart Thinking Dropdown -->
      <div class="relative">
        <button 
          @click.stop="isDropdownOpen = !isDropdownOpen; isMoreDropdownOpen = false"
          class="flex items-center gap-1 px-3 py-1.5 text-sm font-medium text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors"
        >
          <component :is="currentModeIcon" class="w-4 h-4 text-primary" />
          <span class="hidden md:inline">{{ currentModeLabel }}</span>
          <ChevronDown class="w-3 h-3 text-gray-500" />
        </button>

        <!-- Dropdown Menu -->
        <div 
          v-if="isDropdownOpen"
          class="absolute bottom-full left-0 mb-2 w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-50 p-1"
        >
          <button 
            v-for="mode in modes" 
            :key="mode.value"
            @click="emit('select-mode', mode.value); isDropdownOpen = false"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
            :class="{ 'bg-blue-50 dark:bg-blue-900/20': currentMode === mode.value }"
          >
            <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
              <component :is="mode.icon" class="w-4 h-4" :class="mode.color" />
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ mode.label }}</span>
              <span class="text-xs text-gray-500">{{ mode.desc }}</span>
            </div>
            <Check v-if="currentMode === mode.value" class="w-4 h-4 text-primary ml-auto" />
          </button>
        </div>
      </div>
      
      <div class="h-4 w-px bg-gray-200 dark:bg-gray-700 mx-2"></div>

      <button 
        @click="emit('screenshot')"
        class="flex items-center gap-1 px-2 py-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors text-sm"
        title="截屏当前会话"
      >
        <Camera class="w-4 h-4" />
        <span class="hidden md:inline">截屏</span>
      </button>

      <!-- More Dropdown -->
      <div class="relative">
        <button 
          @click.stop="isMoreDropdownOpen = !isMoreDropdownOpen; isDropdownOpen = false"
          class="flex items-center gap-1 px-2 py-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors text-sm"
          title="更多选项"
        >
          <MoreHorizontal class="w-4 h-4" />
        </button>

        <!-- More Menu -->
        <div 
          v-if="isMoreDropdownOpen"
          class="absolute bottom-full left-0 mb-2 w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden z-50 p-1"
        >
          <!-- Web Search Toggle -->
          <button 
            @click="emit('toggle-web-search')"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
            :class="{ 'bg-blue-50 dark:bg-blue-900/20': searchByWeb }"
          >
            <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
              <Globe class="w-4 h-4" :class="searchByWeb ? 'text-blue-500' : 'text-gray-500'" />
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-medium text-gray-900 dark:text-gray-100">联网搜索</span>
              <span class="text-xs text-gray-500">启用实时网络搜索</span>
            </div>
            <Check v-if="searchByWeb" class="w-4 h-4 text-primary ml-auto" />
          </button>
          
          <div v-if="extButtons && extButtons.length > 0" class="h-px bg-gray-100 dark:bg-gray-700 my-1"></div>

          <!-- External Buttons -->
          <button 
            v-for="btn in extButtons"
            :key="btn.id"
            @click="emit('ext-action', btn); isMoreDropdownOpen = false"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors text-left"
          >
            <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">
              <component :is="getIcon(btn.icon)" class="w-4 h-4 text-gray-500" />
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ btn.name }}</span>
              <span class="text-xs text-gray-500">{{ btn.desc }}</span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-1">
       <input 
        type="file" 
        ref="fileInput" 
        class="hidden" 
        @change="(e) => emit('file-change', e, 'file')" 
      />
      <button @click="triggerFileUpload" class="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors" title="上传文件">
        <Paperclip class="w-5 h-5" />
      </button>
      
      <input 
        type="file" 
        ref="imageInput" 
        class="hidden" 
        accept="image/*"
        @change="(e) => emit('file-change', e, 'image')" 
      />
      
      <VoiceInput
        size="md"
        @start="emit('voice-start')"
        @stop="emit('voice-stop')"
        @result="emit('voice-result', $event)"
        @interim="emit('voice-interim', $event)"
        @error="emit('voice-error', $event)"
      />

      <button 
        class="p-2 rounded-full transition-all duration-200 ease-in-out ml-1"
        :class="[
          (input.trim() || hasAttachments) || isLoading
            ? 'bg-primary text-white hover:bg-blue-600 shadow-md hover:shadow-lg transform hover:-translate-y-0.5' 
            : 'bg-gray-100 text-gray-400 dark:bg-gray-700 dark:text-gray-500 cursor-not-allowed',
          isLoading ? 'bg-red-500 hover:bg-red-600' : ''
        ]"
        :disabled="(!input.trim() && !hasAttachments) && !isLoading"
        @click="isLoading ? emit('stop') : emit('send')"
      >
        <div v-if="isLoading" class="flex items-center justify-center w-5 h-5">
          <Square class="w-3 h-3 fill-current" />
        </div>
        <ArrowUp v-else class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>
