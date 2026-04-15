<script setup lang="ts">
import { Carousel, Image } from 'ant-design-vue';

export interface Content0Item {
  icon?: string;
  images?: string[];
  /** split：左文右手机图；split-browser：左 Mac 窗口轮播 + 右文；coming-soon：无图提示 */
  layout?: 'default' | 'split' | 'split-browser' | 'coming-soon';
  title: string;
  content: string;
}

const props = defineProps<{
  title: string;
  items: Content0Item[];
}>();

const slidesFor = (item: Content0Item) => {
  const raw = item.images?.map((s) => s.trim()).filter(Boolean) ?? [];
  if (raw.length > 0) {
    return raw;
  }
  const legacy = item.icon?.trim();
  return legacy ? [legacy] : [];
};
</script>

<template>
  <section id="products" data-section="products" class="content0-wrapper">
    <div class="content0">
      <div class="title-wrapper scroll-reveal-stagger-item">
        <h2 class="title-h1">{{ title }}</h2>
      </div>

      <div class="content0-block-wrapper">
        <article
          v-for="(item, index) in items"
          :key="`row-${index}-${item.title}`"
          class="content0-row scroll-reveal-stagger-item"
        >
          <!-- AI Chat：左 Mac 浏览器框 + 轮播，右文案 -->
          <div v-if="item.layout === 'split-browser'" class="content0-split content0-split--browser">
            <div class="content0-mac">
              <div class="content0-mac-toolbar" aria-hidden="true">
                <div class="content0-mac-traffic">
                  <span class="content0-mac-dot content0-mac-dot--red" />
                  <span class="content0-mac-dot content0-mac-dot--yellow" />
                  <span class="content0-mac-dot content0-mac-dot--green" />
                </div>
                <div class="content0-mac-address">
                  <span class="content0-mac-url">aichat.app</span>
                </div>
              </div>
              <div class="content0-mac-viewport">
                <Carousel
                  v-if="slidesFor(item).length > 1"
                  class="content0-carousel content0-carousel--mac"
                  autoplay
                  :autoplay-speed="4500"
                  effect="fade"
                  dots
                  :pause-on-hover="true"
                >
                  <div v-for="src in slidesFor(item)" :key="src" class="content0-slide">
                    <div class="content0-img-frame content0-img-frame--mac">
                      <Image
                        class="content0-preview-image content0-preview-image--mac"
                        :src="src"
                        alt=""
                        :preview="true"
                        loading="lazy"
                      />
                    </div>
                  </div>
                </Carousel>
                <div v-else-if="slidesFor(item).length === 1" class="content0-img-frame content0-img-frame--mac">
                  <Image
                    class="content0-preview-image content0-preview-image--mac"
                    :src="slidesFor(item)[0]"
                    alt=""
                    :preview="true"
                    loading="lazy"
                  />
                </div>
              </div>
            </div>
            <div class="content0-split-text">
              <h3 class="content0-block-title content0-split-title">{{ item.title }}</h3>
              <p class="content0-block-text content0-split-desc">{{ item.content }}</p>
            </div>
          </div>

          <!-- AICook：左文右手机 -->
          <div v-else-if="item.layout === 'split'" class="content0-split content0-split--phone">
            <div class="content0-split-text">
              <h3 class="content0-block-title content0-split-title">{{ item.title }}</h3>
              <p class="content0-block-text content0-split-desc">{{ item.content }}</p>
            </div>
            <div class="content0-split-visual">
              <div class="content0-device-frame">
                <div class="content0-device-screen">
                  <Carousel
                    v-if="slidesFor(item).length > 1"
                    class="content0-carousel content0-carousel--phone"
                    autoplay
                    :autoplay-speed="4500"
                    effect="fade"
                    dots
                    :pause-on-hover="true"
                  >
                    <div v-for="src in slidesFor(item)" :key="src" class="content0-slide">
                      <div class="content0-img-frame content0-img-frame--phone">
                        <Image
                          class="content0-preview-image content0-preview-image--phone"
                          :src="src"
                          alt=""
                          :preview="true"
                          loading="lazy"
                        />
                      </div>
                    </div>
                  </Carousel>
                  <div v-else-if="slidesFor(item).length === 1" class="content0-img-frame content0-img-frame--phone">
                    <Image
                      class="content0-preview-image content0-preview-image--phone"
                      :src="slidesFor(item)[0]"
                      alt=""
                      :preview="true"
                      loading="lazy"
                    />
                  </div>
                  <div v-else class="content0-visual--empty content0-visual--empty-phone" aria-hidden="true" />
                </div>
              </div>
            </div>
          </div>

          <!-- Tech Sandbox：无图，敬请期待 -->
          <div v-else-if="item.layout === 'coming-soon'" class="content0-coming-soon">
            <h3 class="content0-coming-title">{{ item.title }}</h3>
            <p class="content0-coming-text">{{ item.content }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.content0-wrapper {
  min-height: 200px;
  overflow: hidden;
  background: var(--landing-page-bg, #ececec);
}

.content0 {
  max-width: var(--landing-content-max-width, 1200px);
  margin: 0 auto;
  padding: 64px 24px;
}

.title-wrapper {
  text-align: center;
  margin-bottom: 48px;
}

.title-h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 500;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.content0-block-wrapper {
  display: flex;
  flex-direction: column;
  gap: 56px;
}

.content0-row {
  width: 100%;
}

.content0-block-title {
  margin: 10px auto;
  font-size: 22px;
  line-height: 32px;
  font-weight: 500;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.content0-block-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.75;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.45));
}

/* split-browser：左图右文 */
.content0-split--browser {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(0, 0.95fr);
  gap: 40px;
  align-items: center;
}

.content0-mac {
  width: 100%;
  max-width: 520px;
  margin: 0 auto;
  border-radius: 10px;
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(0, 0, 0, 0.12),
    0 18px 48px rgba(0, 0, 0, 0.14),
    0 4px 12px rgba(0, 0, 0, 0.08);
  background: #ececec;
}

.content0-mac-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 40px;
  padding: 0 12px;
  background: linear-gradient(180deg, #e4e4e4 0%, #d8d8d8 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.content0-mac-traffic {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.content0-mac-dot {
  width: 11px;
  height: 11px;
  border-radius: 50%;
  box-shadow: 0 0 0 0.5px rgba(0, 0, 0, 0.15) inset;
}

.content0-mac-dot--red {
  background: #ff5f57;
}

.content0-mac-dot--yellow {
  background: #febc2e;
}

.content0-mac-dot--green {
  background: #28c840;
}

.content0-mac-address {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.content0-mac-url {
  max-width: 100%;
  padding: 4px 14px;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.55);
  background: rgba(255, 255, 255, 0.65);
  border-radius: 5px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.content0-mac-viewport {
  background: #f5f5f5;
  padding: 10px;
}

.content0-carousel--mac :deep(.slick-dots) {
  bottom: -22px;
  position: relative;
  margin-top: 8px;
}

.content0-carousel--mac :deep(.slick-dots li button) {
  background: rgba(0, 0, 0, 0.25);
}

.content0-carousel--mac :deep(.slick-dots li.slick-active button) {
  background: rgba(0, 0, 0, 0.5);
}

.content0-img-frame--mac {
  min-height: 160px;
  max-height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border-radius: 6px;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.content0-preview-image--mac :deep(.ant-image-img) {
  max-height: 300px;
  max-width: 100%;
  width: auto !important;
  height: auto !important;
  object-fit: contain;
  cursor: zoom-in;
}

.content0-split--browser .content0-split-text {
  text-align: left;
  padding-left: 8px;
}

.content0-split--browser .content0-split-title {
  margin: 0 0 12px;
  text-align: left;
}

.content0-split--browser .content0-split-desc {
  text-align: left;
}

/* AICook：左文右手机 */
.content0-split--phone {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.1fr);
  gap: 40px;
  align-items: center;
}

.content0-split--phone .content0-split-text {
  text-align: left;
  padding-right: 8px;
}

.content0-split-title {
  margin: 0 0 12px;
}

.content0-split-desc {
  text-align: left;
}

.content0-split-visual {
  display: flex;
  justify-content: center;
  align-items: center;
}

.content0-carousel {
  width: 100%;
  max-width: 100%;
}

.content0-carousel--phone :deep(.slick-dots) {
  bottom: -28px;
}

.content0-carousel--phone :deep(.slick-dots li button) {
  background: rgba(255, 255, 255, 0.35);
}

.content0-carousel--phone :deep(.slick-dots li.slick-active button) {
  background: rgba(255, 255, 255, 0.85);
}

.content0-slide {
  outline: none;
}

.content0-img-frame--phone {
  min-height: 0;
  max-height: none;
  padding: 0;
  background: transparent;
}

.content0-preview-image :deep(.ant-image) {
  display: flex !important;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.content0-preview-image--phone :deep(.ant-image-img) {
  max-height: min(420px, 52vh);
  border-radius: 0;
  cursor: zoom-in;
  object-fit: contain;
  width: auto !important;
  height: auto !important;
}

.content0-visual--empty-phone {
  aspect-ratio: 9 / 19;
  max-height: 420px;
  margin: 0 auto;
  max-width: 220px;
}

/* iPhone 外框 */
.content0-device-frame {
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
  padding: 10px;
  border-radius: 42px;
  background: linear-gradient(160deg, #2c2c2e 0%, #1c1c1e 40%, #0d0d0f 100%);
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.06) inset,
    0 16px 48px rgba(0, 0, 0, 0.28),
    0 4px 12px rgba(0, 0, 0, 0.12);
}

.content0-device-screen {
  border-radius: 32px;
  overflow: visible;
  background: #000;
  position: relative;
}

.content0-device-screen .content0-carousel--phone :deep(.slick-list) {
  border-radius: 28px;
  overflow: hidden;
  background: #000;
}

/* Tech Sandbox */
.content0-coming-soon {
  text-align: center;
  padding: 40px 24px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.03);
  border: 1px dashed var(--landing-line, #d9d9d9);
  max-width: 560px;
  margin: 0 auto;
}

.content0-coming-title {
  margin: 0 0 12px;
  font-size: 22px;
  font-weight: 500;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.content0-coming-text {
  margin: 0;
  font-size: 15px;
  line-height: 1.7;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.5));
}

@media (max-width: 767px) {
  .content0-split--browser,
  .content0-split--phone {
    grid-template-columns: 1fr;
    gap: 28px;
  }

  .content0-split--browser .content0-split-text,
  .content0-split--phone .content0-split-text {
    text-align: center;
    padding-left: 0;
    padding-right: 0;
  }

  .content0-split--browser .content0-split-title,
  .content0-split--browser .content0-split-desc,
  .content0-split--phone .content0-split-title,
  .content0-split--phone .content0-split-desc {
    text-align: center;
  }

  .content0-mac {
    max-width: 100%;
  }

  .content0-device-frame {
    max-width: min(300px, 88vw);
  }

  .content0-wrapper {
    min-height: 0;
  }
}
</style>
