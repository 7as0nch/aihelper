import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export type ThemeMode = 'light' | 'dark' | 'system';

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>('light');

  if (typeof window !== 'undefined') {
    const storedTheme = window.localStorage.getItem('theme-mode') as ThemeMode | null;
    if (storedTheme && ['light', 'dark', 'system'].includes(storedTheme)) {
      mode.value = storedTheme;
    }
  }

  const applyTheme = () => {
    if (typeof document === 'undefined' || typeof window === 'undefined') {
      return;
    }

    const root = document.documentElement;

    if (mode.value === 'system') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      root.classList.toggle('dark', prefersDark);
      return;
    }

    root.classList.toggle('dark', mode.value === 'dark');
  };

  watch(
    mode,
    (nextMode) => {
      if (typeof window !== 'undefined') {
        window.localStorage.setItem('theme-mode', nextMode);
      }
      applyTheme();
    },
    { immediate: true },
  );

  if (typeof window !== 'undefined') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    mediaQuery.addEventListener('change', () => {
      if (mode.value === 'system') {
        applyTheme();
      }
    });
  }

  const setMode = (newMode: ThemeMode) => {
    mode.value = newMode;
  };

  return {
    mode,
    setMode,
  };
});