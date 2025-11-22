<script setup lang="ts">
import { ref } from 'vue';
import { ChevronDown, Globe, Camera, Paperclip, Mic, ArrowUp, Square, Check } from 'lucide-vue-next';

defineProps<{
  currentMode: string;
  modes: readonly any[];
  currentModeLabel: string;
  currentModeIcon: any;
  input: string;
  isLoading: boolean;
  hasAttachments: boolean;
  isRecording: boolean;
}>();

const emit = defineEmits<{
  (e: 'select-mode', mode: string): void;
  (e: 'screenshot'): void;
  (e: 'upload'): void;
  (e: 'voice-input'): void;
  (e: 'send'): void;
  (e: 'stop'): void;
  (e: 'file-change', event: Event, type: 'file' | 'image'): void;
}>();

const isDropdownOpen = ref(false);

const fileInput = ref<HTMLInputElement | null>(null);
const imageInput = ref<HTMLInputElement | null>(null);

const triggerFileUpload = () => {
  fileInput.value?.click();
};

const closeDropdown = () => {
  isDropdownOpen.value = false;
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
          @click.stop="isDropdownOpen = !isDropdownOpen"
          class="flex items-center gap-1 px-3 py-1.5 text-sm font-medium text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg transition-colors"
        >
          <component :is="currentModeIcon" class="w-4 h-4 text-primary" />
          <span class="hidden md:inline">{{ currentModeLabel }}</span>
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

      <button class="flex items-center gap-1 px-2 py-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors text-sm">
        <Globe class="w-4 h-4" />
        <!-- 是否启用联网搜索的MCP Tool -->
        <span class="hidden md:inline">联网搜索</span>
      </button>

      <button 
        @click="emit('screenshot')"
        class="flex items-center gap-1 px-2 py-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors text-sm"
        title="截屏当前会话"
      >
        <Camera class="w-4 h-4" />
        <span class="hidden md:inline">截屏</span>
      </button>
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
      
      <button 
        @click="emit('voice-input')" 
        class="p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors"
        :class="{ 'text-red-500 bg-red-50 dark:bg-red-900/20': isRecording }"
        title="语音输入"
      >
        <Mic class="w-5 h-5" />
      </button>

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
