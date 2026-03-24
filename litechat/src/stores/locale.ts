import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

import { i18n, setI18nLanguage } from '@/i18n';
import {
  DEFAULT_LOCALE,
  LOCALE_STORAGE_KEY,
  isSupportedLocale,
  type SupportedLocale,
} from '@/i18n/types';

function resolveStoredLocale(): SupportedLocale {
  if (typeof window === 'undefined') {
    return DEFAULT_LOCALE;
  }

  const stored = window.localStorage.getItem(LOCALE_STORAGE_KEY);
  if (stored && isSupportedLocale(stored)) {
    return stored;
  }

  return DEFAULT_LOCALE;
}

export const useLocaleStore = defineStore('locale', () => {
  const locale = ref<SupportedLocale>(resolveStoredLocale());

  const applyLocale = (nextLocale: SupportedLocale) => {
    locale.value = nextLocale;
  };

  watch(
    locale,
    (nextLocale) => {
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(LOCALE_STORAGE_KEY, nextLocale);
      }

      setI18nLanguage(nextLocale);
    },
    { immediate: true },
  );

  const initialize = () => {
    const initial = resolveStoredLocale();
    locale.value = initial;
    i18n.global.locale.value = initial;
    setI18nLanguage(initial);
  };

  return {
    locale,
    applyLocale,
    initialize,
  };
});