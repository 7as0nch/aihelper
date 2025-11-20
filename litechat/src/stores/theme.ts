import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export type ThemeMode = 'light' | 'dark' | 'system';

export const useThemeStore = defineStore('theme', () => {
    const mode = ref<ThemeMode>('system');

    // Get initial theme from localStorage or default to system
    const storedTheme = localStorage.getItem('theme-mode') as ThemeMode;
    if (storedTheme && ['light', 'dark', 'system'].includes(storedTheme)) {
        mode.value = storedTheme;
    }

    // Apply theme to document
    const applyTheme = () => {
        const root = document.documentElement;

        if (mode.value === 'system') {
            // Follow system preference
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            if (prefersDark) {
                root.classList.add('dark');
            } else {
                root.classList.remove('dark');
            }
        } else if (mode.value === 'dark') {
            root.classList.add('dark');
        } else {
            root.classList.remove('dark');
        }
    };

    // Watch for mode changes
    watch(mode, () => {
        localStorage.setItem('theme-mode', mode.value);
        applyTheme();
    }, { immediate: true });

    // Listen for system theme changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    mediaQuery.addEventListener('change', () => {
        if (mode.value === 'system') {
            applyTheme();
        }
    });

    const setMode = (newMode: ThemeMode) => {
        mode.value = newMode;
    };

    return {
        mode,
        setMode,
    };
});
