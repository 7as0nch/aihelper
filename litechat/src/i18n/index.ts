import { createI18n } from 'vue-i18n';

import zhLanding from '@/locales/zh-CN/landing';
import enLanding from '@/locales/en-US/landing';
import {
  DEFAULT_LOCALE,
  LOCALE_STORAGE_KEY,
  isSupportedLocale,
  type SupportedLocale,
} from './types';

const messages = {
  'zh-CN': {
    landing: zhLanding,
  },
  'en-US': {
    landing: enLanding,
  },
};

function resolveInitialLocale(): SupportedLocale {
  if (typeof window === 'undefined') {
    return DEFAULT_LOCALE;
  }

  const stored = window.localStorage.getItem(LOCALE_STORAGE_KEY);
  if (stored && isSupportedLocale(stored)) {
    return stored;
  }

  return DEFAULT_LOCALE;
}

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: resolveInitialLocale(),
  fallbackLocale: DEFAULT_LOCALE,
  messages,
});

export function setI18nLanguage(locale: SupportedLocale) {
  i18n.global.locale.value = locale;

  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('lang', locale);
  }
}