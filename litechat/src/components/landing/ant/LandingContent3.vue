<script setup lang="ts">
export interface Content3Item {
  icon?: string;
  title: string;
  description: string;
}

defineProps<{
  title: string;
  subtitle: string;
  items: Content3Item[];
}>();
</script>

<template>
  <section id="updates" data-section="updates" class="content3-wrapper">
    <div class="content3">
      <div class="title-wrapper scroll-reveal-stagger-item">
        <h2 class="title-h1">{{ title }}</h2>
        <p class="title-content">{{ subtitle }}</p>
      </div>
      <div class="content3-block-wrapper">
        <article
          v-for="(item, index) in items"
          :key="index"
          class="content3-block scroll-reveal-stagger-item"
        >
          <div v-if="!item.icon?.trim()" class="content3-icon content3-icon--placeholder" />
          <img v-else :src="item.icon" alt="" class="content3-icon" />
          <div class="content3-text">
            <h3 class="content3-title">{{ item.title }}</h3>
            <p class="content3-content">{{ item.description }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.content3-wrapper {
  overflow: hidden;
  background: var(--landing-page-bg, #ececec);
}

.content3 {
  max-width: var(--landing-content-max-width, 1200px);
  margin: 0 auto;
  padding: 64px 24px;
}

.title-wrapper {
  text-align: center;
  margin-bottom: 48px;
}

.title-h1 {
  margin: 0 0 8px;
  font-size: 28px;
  font-weight: 500;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.title-content {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.45));
}

.content3-block-wrapper {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 32px 24px;
}

.content3-block {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.content3-icon {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  object-fit: contain;
}

.content3-icon--placeholder {
  border-radius: 8px;
  background: #d9d9d9;
}

.content3-text {
  min-width: 0;
}

.content3-title {
  margin: 0 0 8px;
  font-size: 16px;
  font-weight: 500;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.content3-content {
  margin: 0;
  font-size: 13px;
  line-height: 1.7;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.45));
}

@media (max-width: 991px) {
  .content3-block-wrapper {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 575px) {
  .content3-block-wrapper {
    grid-template-columns: 1fr;
  }
}
</style>
