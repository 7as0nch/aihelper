<script setup lang="ts">
const props = defineProps<{
  title: string;
  content: string;
  image?: string;
}>();
</script>

<template>
  <section id="studio" data-section="studio" class="content1-wrapper">
    <div class="content1">
      <div class="content1-visual scroll-reveal-stagger-item">
        <div v-if="!props.image?.trim()" class="content1-placeholder" role="img" aria-label="" />
        <img v-else :src="props.image" alt="" class="content1-img" />
      </div>
      <div class="content1-text scroll-reveal-stagger-item">
        <h2 class="content1-title">{{ title }}</h2>
        <p class="content1-content">{{ content }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.content1-wrapper {
  overflow: hidden;
  background: #f5f5f5;
}

.content1 {
  max-width: var(--landing-content-max-width, 1200px);
  margin: 0 auto;
  padding: 64px 24px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 40px;
  align-items: center;
}

.content1-visual {
  min-height: 200px;
}

.content1-placeholder {
  width: 100%;
  aspect-ratio: 5 / 4;
  max-height: 360px;
  border-radius: 8px;
  background: linear-gradient(160deg, #e0e0e0, #bdbdbd);
  border: 1px solid var(--landing-line, #e9e9e9);
}

.content1-img {
  width: 100%;
  height: auto;
  max-height: 400px;
  object-fit: contain;
  border-radius: 8px;
}

.content1-title {
  margin: 0 0 16px;
  font-size: 26px;
  font-weight: 500;
  line-height: 1.4;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.content1-content {
  margin: 0;
  font-size: 14px;
  line-height: 1.85;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.65));
}

@media (max-width: 767px) {
  .content1 {
    grid-template-columns: 1fr;
  }

  .content1-visual {
    order: -1;
  }
}
</style>
