<script setup lang="ts">
import type { SupportedLocale } from '@/i18n/types';
import { useI18n } from 'vue-i18n';

const props = withDefaults(
  defineProps<{
    modelValue: SupportedLocale;
    /** Dark nav bar (Ant landing) */
    tone?: 'light' | 'dark';
  }>(),
  { tone: 'light' },
);

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
  <div
    class="language-switch"
    :class="{ 'language-switch--dark': props.tone === 'dark' }"
    :aria-label="t('landing.language.label')"
  >
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
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.66);
  border: 1px solid var(--site-border);
}

.language-option {
  min-width: 52px;
  padding: 9px 12px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--site-muted);
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition:
    background 0.22s ease,
    color 0.22s ease,
    transform 0.22s ease;
}

.language-option.is-active {
  background: rgba(199, 91, 44, 0.14);
  color: var(--site-text);
  transform: translateY(-1px);
}

.language-switch--dark {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.language-switch--dark .language-option {
  color: rgba(255, 255, 255, 0.65);
}

.language-switch--dark .language-option.is-active {
  background: var(--landing-primary, #1890ff);
  color: #fff;
  transform: none;
}
</style>