<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { Popover, Image } from 'ant-design-vue';
import { VerticalAlignTopOutlined, WechatOutlined } from '@ant-design/icons-vue';

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
  <div class="floating-actions">
    <!-- Contact WeChat (Always visible or same threshold as BTT) -->
    <div class="action-item contact-action">
      <Popover placement="left" trigger="hover">
        <template #content>
          <div class="wechat-popover-content">
            <Image
              src="/wechat.jpg"
              :width="160"
              alt="WeChat QR Code"
              class="wechat-qr"
            />
            <p class="wechat-tip">扫码联系我</p>
          </div>
        </template>
        <button class="action-btn contact-btn" aria-label="Contact WeChat">
          <WechatOutlined />
        </button>
      </Popover>
    </div>

    <!-- Back to Top (Conditional visibility) -->
    <Transition name="action-slide">
      <div v-if="visible" class="action-item btt-action">
        <button
          class="action-btn btt-btn"
          :aria-label="t('landing.backToTop')"
          :title="t('landing.backToTop')"
          @click="scrollToTop"
        >
          <VerticalAlignTopOutlined />
        </button>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.floating-actions {
  position: fixed;
  right: 24px;
  bottom: 32px;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  color: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.action-btn:hover {
  background: #fff;
  color: #1890ff;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  border-color: rgba(24, 144, 255, 0.2);
}

.contact-btn {
  color: #52c41a; /* WeChat Green */
}

.contact-btn:hover {
  color: #52c41a;
  border-color: rgba(82, 196, 26, 0.2);
}

/* Popover Content */
.wechat-popover-content {
  text-align: center;
  padding: 4px;
}

.wechat-tip {
  margin-top: 8px;
  margin-bottom: 0;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.45);
}

/* Animations */
.action-slide-enter-active,
.action-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.action-slide-enter-from,
.action-slide-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.8);
}
</style>

<style>
/* Adjust Ant Design Image internal style if needed */
.wechat-qr img {
  display: block;
  border-radius: 4px;
}
</style>
