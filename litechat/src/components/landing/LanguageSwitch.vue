<script setup lang="ts">
import type { SupportedLocale } from '@/i18n/types';

import { useI18n } from 'vue-i18n';

const props = defineProps<{
  modelValue: SupportedLocale;
}>();

const emit = defineEmits<{
  (event: 'update:modelValue', value: SupportedLocale): void;
}>();

const { t } = useI18n();

const setLocale = (locale: SupportedLocale) => {
  if (locale !== props.modelValue) {
    emit('update:modelValue', locale);
  }
};
</script>

<template>
  <div class="language-switch" :aria-label="t('landing.language.label')">
    <button
      type="button"
      class="language-option"
      :class="{ 'is-active': modelValue === 'zh-CN' }"
      @click="setLocale('zh-CN')"
    >
      {{ t('landing.language.zh') }}
    </button>
    <button
      type="button"
      class="language-option"
      :class="{ 'is-active': modelValue === 'en-US' }"
      @click="setLocale('en-US')"
    >
      {{ t('landing.language.en') }}
    </button>
  </div>
</template>

<style scoped>
.language-switch {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px;
  border-radius: 14px;
  background: var(--surface-strong);
  border: 1px solid var(--border-soft);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.18);
}

.language-option {
  border: 0;
  cursor: pointer;
  min-width: 52px;
  padding: 9px 12px;
  border-radius: 14px;
  background: transparent;
  color: var(--text-muted);
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  transition: background 0.25s ease, color 0.25s ease, transform 0.25s ease;
}

.language-option.is-active {
  background: rgba(255, 255, 255, 0.9);
  color: var(--text-strong);
  transform: translateY(-1px);
}

:global(.theme-dark) .language-option.is-active {
  background: rgba(125, 211, 252, 0.14);
}
</style>