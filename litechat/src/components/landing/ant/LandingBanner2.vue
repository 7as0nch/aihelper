<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ChevronDown } from 'lucide-vue-next';

const props = defineProps<{
  title: string;
  content: string;
  button: string;
  /** Optional full URL for background; empty uses gradient placeholder */
  backgroundImage?: string;
}>();

const hasBgImage = computed(() => Boolean(props.backgroundImage?.trim()));

/** -1 … 1，相对横幅中心：左负右正；用于景深视差 */
const parallaxTilt = ref(0);
const bannerRef = ref<HTMLElement | null>(null);
const reduceMotion = ref(false);

/** 背景横向偏移 px（与文字反向） */
const BG_SHIFT_MAX = 22;
/** 文字横向偏移 px */
const TEXT_SHIFT_MAX = 14;

onMounted(() => {
  reduceMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
});

const bgStyle = computed(() => {
  const url = props.backgroundImage?.trim();
  if (!url) {
    return undefined;
  }
  const t = reduceMotion.value ? 0 : parallaxTilt.value;
  const bgPx = t * -BG_SHIFT_MAX;
  return {
    backgroundImage: `url(${url})`,
    backgroundSize: 'cover',
    backgroundPosition: `calc(50% + ${bgPx}px) center`,
    backgroundRepeat: 'no-repeat',
  };
});

const textParallaxStyle = computed(() => {
  if (!hasBgImage.value || reduceMotion.value) {
    return undefined;
  }
  const t = parallaxTilt.value;
  const tx = t * TEXT_SHIFT_MAX;
  return {
    transform: `translate3d(${tx}px, 0, 0)`,
  };
});

const onBannerMouseMove = (e: MouseEvent) => {
  if (reduceMotion.value || !hasBgImage.value) {
    return;
  }
  const el = bannerRef.value;
  if (!el) {
    return;
  }
  const r = el.getBoundingClientRect();
  const w = r.width || 1;
  const t = ((e.clientX - r.left) / w - 0.5) * 2;
  parallaxTilt.value = Math.max(-1, Math.min(1, t));
};

const onBannerMouseLeave = () => {
  parallaxTilt.value = 0;
};

const scrollToIntro = () => {
  document.getElementById('landing-intro')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
};

const scrollWithButton = () => {
  scrollToIntro();
};
</script>

<template>
  <section
    ref="bannerRef"
    class="banner2"
    :class="{ 'banner2--fixed-bg': hasBgImage, 'banner2--parallax': hasBgImage }"
    :style="bgStyle"
    @mousemove="onBannerMouseMove"
    @mouseleave="onBannerMouseLeave"
  >
    <div class="banner2-overlay" aria-hidden="true" />
    <div class="banner2-page">
      <div class="banner2-text-wrapper" :style="textParallaxStyle">
        <h1 class="banner2-title banner2-queue">{{ title }}</h1>
        <p class="banner2-content banner2-queue">{{ content }}</p>
        <button type="button" class="banner2-button banner2-queue" @click="scrollWithButton">
          {{ button }}
        </button>
      </div>
    </div>
    <button type="button" class="banner2-scroll" aria-label="Scroll down" @click="scrollToIntro">
      <ChevronDown class="banner2-scroll-icon" :size="28" stroke-width="2" />
    </button>
  </section>
</template>

<style scoped>
.banner2 {
  position: relative;
  width: 100%;
  height: calc(100vh - var(--landing-header-height, 64px));
  min-height: 420px;
  overflow: hidden;
  background:
    linear-gradient(135deg, #1d3754 0%, #243b55 40%, #141e30 100%);
  background-size: cover;
  background-position: center;
}

.banner2[style*='background-image'] {
  background-color: #1a1a1a;
}

/* 与 Ant 首屏类似：背景相对视口固定，滚动时画面不跟着走 */
.banner2--fixed-bg {
  background-attachment: fixed;
}

@media (max-width: 767px) {
  .banner2--fixed-bg {
    background-attachment: scroll;
  }
}

/* 鼠标横向视差：背景与文字反向微移（不在 mousemove 上加 transition，避免滞后） */
.banner2--parallax .banner2-text-wrapper {
  will-change: transform;
}

.banner2-overlay {
  position: absolute;
  inset: 0;
  background: var(--landing-banner-overlay, rgba(0, 0, 0, 0.35));
  pointer-events: none;
}

.banner2-page {
  position: relative;
  z-index: 1;
  max-width: var(--landing-content-max-width, 1200px);
  margin: 0 auto;
  padding: 0 clamp(16px, 4vw, 48px);
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  box-sizing: border-box;
}

.banner2-text-wrapper {
  max-width: min(440px, 100%);
  text-align: left;
  color: #fff;
  font-size: 14px;
}

/* rc-queue-anim: first screen enter */
.banner2-queue {
  animation: banner2-queue-in 0.65s cubic-bezier(0.25, 0.46, 0.45, 0.94) backwards;
}

.banner2-title.banner2-queue {
  animation-delay: 0.08s;
}

.banner2-content.banner2-queue {
  animation-delay: 0.2s;
}

.banner2-button.banner2-queue {
  animation-delay: 0.34s;
}

@keyframes banner2-queue-in {
  from {
    opacity: 0;
    transform: translate3d(0, 36px, 0);
  }

  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}

.banner2-title {
  margin: 0 0 16px;
  font-size: clamp(2.2rem, 6vw, 3.5rem);
  line-height: 1.15;
  font-weight: 600;
  letter-spacing: -0.03em;
}

.banner2-content {
  margin: 0 0 20px;
  font-size: 1rem;
  line-height: 1.75;
  font-weight: 300;
  opacity: 0.95;
}

.banner2-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
  padding: 0 18px;
  border: 1px solid #fff;
  border-radius: 2px;
  background: transparent;
  color: #fff;
  font-size: 16px;
  line-height: 36px;
  cursor: pointer;
  transition:
    background 0.35s ease,
    border-color 0.35s ease;
}

.banner2-button:hover {
  border-color: var(--landing-primary, #1890ff);
  background: var(--landing-primary, #1890ff);
}

.banner2-scroll {
  position: absolute;
  bottom: 36px;
  left: 50%;
  z-index: 2;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  min-height: 44px;
  padding: 8px;
  border: none;
  border-radius: 0;
  background: transparent;
  color: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  animation: banner2-bounce 1s ease-in-out infinite;
  transition: color 0.2s ease, opacity 0.2s ease;
}

.banner2-scroll-icon {
  display: block;
}

.banner2-scroll:hover {
  color: #fff;
  opacity: 1;
}

@keyframes banner2-bounce {
  0%,
  100% {
    transform: translate(-50%, 0);
  }

  50% {
    transform: translate(-50%, -10px);
  }
}

@media (max-width: 767px) {
  .banner2-page {
    justify-content: center;
    padding: 0 20px;
  }

  .banner2-text-wrapper {
    text-align: center;
  }

  .banner2-title {
    font-size: 2rem;
  }
}
</style>
