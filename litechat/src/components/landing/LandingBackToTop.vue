<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { Button } from 'ant-design-vue';
import { VerticalAlignTopOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const visible = ref(false);
const SHOW_AFTER = 240;

let raf = 0;
const onScroll = () => {
  if (raf) {
    return;
  }
  raf = requestAnimationFrame(() => {
    raf = 0;
    const next = window.scrollY > SHOW_AFTER;
    if (visible.value !== next) {
      visible.value = next;
    }
  });
};

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

onMounted(() => {
  onScroll();
  window.addEventListener('scroll', onScroll, { passive: true });
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll);
  if (raf) {
    cancelAnimationFrame(raf);
    raf = 0;
  }
});
</script>

<template>
  <Transition name="landing-btt-fade">
    <Button
      v-if="visible"
      type="primary"
      shape="circle"
      class="landing-back-to-top"
      :aria-label="t('landing.backToTop')"
      :title="t('landing.backToTop')"
      @click="scrollToTop"
    >
      <template #icon>
        <VerticalAlignTopOutlined />
      </template>
    </Button>
  </Transition>
</template>

<style scoped>
.landing-back-to-top {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 99;
  width: 44px !important;
  height: 44px !important;
  min-width: 44px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.landing-btt-fade-enter-active,
.landing-btt-fade-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.landing-btt-fade-enter-from,
.landing-btt-fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
