<script setup lang="ts">
import { computed } from 'vue';
import { Carousel } from 'ant-design-vue';

const props = defineProps<{
  title: string;
  content: string;
  button: string;
  /** 若填写则仅展示单图；否则轮播 public/image 下文件名含 main 的图片 */
  image?: string;
}>();

/**
 * public/image 下文件名包含 main 的资源（顺序轮播）。
 * 新增同类资源时请把路径加进此列表。
 */
const MAIN_SLIDE_SRCS = [
  '/image/aichat/main.jpg',
  '/image/aichat/main_page.png',
  '/image/aicook/cook_main.png',
] as const;

const slideSrcs = computed(() => {
  const single = props.image?.trim();
  if (single) {
    return [single];
  }
  return [...MAIN_SLIDE_SRCS];
});

const scrollToProducts = () => {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
};
</script>

<template>
  <section id="landing-intro" class="banner4">
    <div class="banner4-page">
      <div class="banner4-grid">
        <div class="banner4-copy scroll-reveal-stagger-item">
          <h2 class="banner4-title">{{ props.title }}</h2>
          <p class="banner4-content">{{ props.content }}</p>
          <button type="button" class="banner4-button" @click="scrollToProducts">
            {{ props.button }}
          </button>
        </div>
        <div class="banner4-visual scroll-reveal-stagger-item">
          <Carousel
            v-if="slideSrcs.length > 1"
            class="banner4-carousel"
            autoplay
            :autoplay-speed="5000"
            effect="fade"
            dots
            :pause-on-hover="true"
          >
            <div v-for="src in slideSrcs" :key="src" class="banner4-slide">
              <img class="banner4-carousel-img" :src="src" alt="" loading="lazy" />
            </div>
          </Carousel>
          <img
            v-else-if="slideSrcs.length === 1"
            :src="slideSrcs[0]"
            alt=""
            class="banner4-image"
            loading="lazy"
          />
          <div v-else class="banner4-placeholder" role="img" aria-hidden="true" />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.banner4 {
  width: 100%;
  background: #4b5564;
  color: #fff;
  padding: 48px 0 64px;
}

.banner4-page {
  max-width: var(--landing-content-max-width, 1200px);
  margin: 0 auto;
  padding: 0 24px;
}

.banner4-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 40px;
  align-items: center;
}

.banner4-title {
  margin: 0 0 12px;
  font-size: clamp(1.75rem, 4vw, 2.5rem);
  line-height: 1.35;
  font-weight: 600;
}

.banner4-content {
  margin: 0 0 20px;
  font-size: 14px;
  line-height: 1.7;
  opacity: 0.95;
}

.banner4-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
  padding: 0 20px;
  border: none;
  border-radius: 2px;
  background: var(--landing-primary, #1890ff);
  color: #fff;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.banner4-button:hover {
  opacity: 0.9;
}

.banner4-visual {
  min-height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner4-carousel {
  width: 100%;
  max-width: 100%;
}

.banner4-carousel :deep(.slick-dots) {
  bottom: -28px;
}

.banner4-carousel :deep(.slick-dots li button) {
  background: rgba(255, 255, 255, 0.35);
}

.banner4-carousel :deep(.slick-dots li.slick-active button) {
  background: rgba(255, 255, 255, 0.95);
}

.banner4-slide {
  outline: none;
}

.banner4-carousel-img {
  width: 100%;
  aspect-ratio: 16 / 10;
  max-height: 320px;
  object-fit: cover;
  border-radius: 8px;
  display: block;
}

.banner4-placeholder {
  width: 100%;
  aspect-ratio: 16 / 10;
  max-height: 320px;
  border-radius: 8px;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.12), rgba(0, 0, 0, 0.15));
  border: 1px dashed rgba(255, 255, 255, 0.25);
}

.banner4-image {
  width: 100%;
  max-height: 360px;
  object-fit: contain;
  border-radius: 8px;
}

@media (max-width: 767px) {
  .banner4-grid {
    grid-template-columns: 1fr;
  }

  .banner4-visual {
    order: -1;
  }
}
</style>
