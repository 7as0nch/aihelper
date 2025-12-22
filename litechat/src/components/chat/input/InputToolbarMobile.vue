<script setup lang="ts">
import { ref } from 'vue';
import { ChevronDown, Globe, Scissors, Paperclip, Check, Sparkles, Image as ImageIcon } from 'lucide-vue-next';

defineProps<{
  currentMode: string;
  modes: readonly any[];
  activeButton: string | null;
  currentModeLabel: string;
  currentModeIcon: any;
  extButtons?: Array<{ id: string; name: string; api: string; desc: string; icon?: any }>;
  searchByWeb?: boolean;
  showUpload?: boolean;
}>();

const emit = defineEmits<{
  (e: 'select-mode', mode: string): void;
  (e: 'trigger-bounce', key: string): void;
  (e: 'toggle-dropdown'): void;
  (e: 'screenshot'): void;
  (e: 'upload'): void;
  (e: 'ext-action', btn: any): void;
  (e: 'toggle-web-search'): void;
}>();

const isDropdownOpen = ref(false);
const dropdownRef = ref<HTMLElement | null>(null);

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value;
  emit('toggle-dropdown');
};

const selectMode = (mode: string) => {
  emit('select-mode', mode);
  isDropdownOpen.value = false;
};

const closeDropdown = () => {
  isDropdownOpen.value = false;
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
  closeDropdown
});
</script>

<template>
  <div class="md:hidden mb-2 overflow-x-auto no-scrollbar touch-pan-x">
    <div class="flex items-center gap-2 w-max px-1">
      <!-- Model Selection -->
      <div class="relative shrink-0" ref="dropdownRef">
        <button 
          @click.stop="toggleDropdown"
          v-tracker="{ type: 'click', name: '模型选择' }"
          class="flex items-center gap-1.5 px-3 py-2 bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-full shadow-sm border border-gray-200 dark:border-gray-700 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-white dark:hover:bg-[#333] transition-colors"
        >
          <component :is="currentModeIcon" class="w-4 h-4 text-primary" />
          <span>{{ currentModeLabel }}</span>
          <ChevronDown class="w-3 h-3 text-gray-500" />
        </button>

        <!-- Dropdown Menu -->
        <div 
          v-if="isDropdownOpen"
          class="fixed bottom-[150px] left-10 mb-2 w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 p-1 animate-fade-in-up z-[99999]"
          @click.stop
        >
          <button 
            v-for="mode in modes" 
            :key="mode.value"
            @click="selectMode(mode.value)"
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

      <!-- Web Search Toggle -->
      <button 
        @click="emit('toggle-web-search'); emit('trigger-bounce', 'search')"
        v-tracker="{ type: 'click', name: '搜索' }"
        class="flex items-center justify-center h-9 px-3 py-2 backdrop-blur-xl rounded-full shadow-sm border transition-all shrink-0 gap-1.5 text-sm font-medium"
        :class="[
          searchByWeb 
            ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 border-blue-200 dark:border-blue-800' 
            : 'bg-white/80 dark:bg-[#2a2a2a]/80 text-gray-700 dark:text-gray-200 border-gray-200 dark:border-gray-700 hover:bg-white dark:hover:bg-[#333] hover:text-blue-500 dark:hover:text-blue-400',
          { 'animate-jump': activeButton === 'search' }
        ]"
      >
        <Globe class="w-4 h-4" :class="searchByWeb ? 'text-blue-600 dark:text-blue-400' : 'text-gray-500 group-hover:text-blue-500'" />
        搜索
      </button>

      <!-- Screenshot -->
      <button 
        @click="emit('screenshot'); emit('trigger-bounce', 'screenshot')"
        v-tracker="{ type: 'click', name: '截屏'}"
        class="flex items-center justify-center h-9 px-3 py-2 bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-full shadow-sm border border-gray-200 dark:border-gray-700 transition-all hover:bg-white dark:hover:bg-[#333] hover:text-blue-500 dark:hover:text-blue-400 shrink-0 gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-200"
        :class="{ 'animate-jump': activeButton === 'screenshot' }"
      >
        <Scissors class="w-4 h-4 text-gray-500 group-hover:text-blue-500" />
        截屏
      </button>

      <!-- Upload -->
      <button 
        v-if="showUpload !== false"
        v-tracker="{ type: 'click', name: '上传' }"
        @click="emit('upload'); emit('trigger-bounce', 'upload')"
        class="flex items-center justify-center h-9 px-3 py-2 bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-full shadow-sm border border-gray-200 dark:border-gray-700 transition-all hover:bg-white dark:hover:bg-[#333] hover:text-blue-500 dark:hover:text-blue-400 shrink-0 gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-200"
        :class="{ 'animate-jump': activeButton === 'upload' }"
      >
        <Paperclip class="w-4 h-4 text-gray-500 group-hover:text-blue-500" />
        上传
      </button>

      <!-- External Buttons -->
      <button 
        v-for="btn in extButtons"
        :key="btn.id"
        @click="emit('ext-action', btn)"
        class="flex items-center justify-center h-9 px-3 py-2 bg-white/80 dark:bg-[#2a2a2a]/80 backdrop-blur-xl rounded-full shadow-sm border border-gray-200 dark:border-gray-700 transition-all hover:bg-white dark:hover:bg-[#333] hover:text-blue-500 dark:hover:text-blue-400 shrink-0 gap-1.5 text-sm font-medium text-gray-700 dark:text-gray-200"
      >
        <component :is="getIcon(btn.icon)" class="w-4 h-4" />
        {{ btn.name }}
      </button>
    </div>
  </div>
</template>
