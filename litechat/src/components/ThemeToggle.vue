<script setup lang="ts">
import { ref, computed } from 'vue';
import { useThemeStore } from '../stores/theme';
import { Sun, Moon, Monitor, Check } from 'lucide-vue-next';

const themeStore = useThemeStore();
const isOpen = ref(false);

const themes = [
  { value: 'light' as const, label: '亮色', icon: Sun },
  { value: 'dark' as const, label: '暗色', icon: Moon },
  { value: 'system' as const, label: '跟随系统', icon: Monitor },
];

const currentTheme = computed(() => themes.find(t => t.value === themeStore.mode));
</script>

<template>
  <div class="relative">
    <button
      @click="isOpen = !isOpen"
      class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
      :title="currentTheme?.label"
    >
      <component :is="currentTheme?.icon" class="w-5 h-5 text-gray-600 dark:text-gray-300" />
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      @click="isOpen = false"
      class="absolute left-0 bottom-full mb-2 w-40 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-50"
    >
      <button
        v-for="theme in themes"
        :key="theme.value"
        @click="themeStore.setMode(theme.value)"
        class="w-full px-3 py-2 flex items-center gap-2 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        :class="{ 'bg-gray-50 dark:bg-gray-700/50': themeStore.mode === theme.value }"
      >
        <component :is="theme.icon" class="w-4 h-4 text-gray-500 dark:text-gray-400" />
        <span class="text-sm text-gray-700 dark:text-gray-200 flex-1 text-left">{{ theme.label }}</span>
        <Check v-if="themeStore.mode === theme.value" class="w-4 h-4 text-primary" />
      </button>
    </div>
  </div>
</template>
