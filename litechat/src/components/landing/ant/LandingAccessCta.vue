<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { ArrowRight } from 'lucide-vue-next';

const props = defineProps<{
  homePath: string;
  eyebrow: string;
  title: string;
  description: string;
  points: string[];
  primary: string;
  secondary: string;
}>();
</script>

<template>
  <section id="access" data-section="access" class="access-cta">
    <div class="access-cta-inner">
      <span class="access-cta-eyebrow scroll-reveal-stagger-item">{{ props.eyebrow }}</span>
      <h2 class="access-cta-title scroll-reveal-stagger-item">{{ props.title }}</h2>
      <p class="access-cta-desc scroll-reveal-stagger-item">{{ props.description }}</p>
      <ul class="access-cta-points scroll-reveal-stagger-item">
        <li v-for="(p, i) in props.points" :key="i">{{ p }}</li>
      </ul>
      <div class="access-cta-actions scroll-reveal-stagger-item">
        <RouterLink to="/apply" class="access-cta-primary">
          {{ props.primary }}
          <ArrowRight :size="18" />
        </RouterLink>
        <RouterLink :to="{ path: props.homePath, hash: '#products' }" class="access-cta-secondary">
          {{ props.secondary }}
        </RouterLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.access-cta {
  padding: 56px 24px 64px;
  background: #f5f5f5;
}

.access-cta-inner {
  max-width: 720px;
  margin: 0 auto;
  text-align: center;
}

.access-cta-eyebrow {
  display: inline-block;
  margin-bottom: 12px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--landing-primary, #1890ff);
}

.access-cta-title {
  margin: 0 0 12px;
  font-size: clamp(1.5rem, 3vw, 2rem);
  font-weight: 600;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
}

.access-cta-desc {
  margin: 0 0 20px;
  font-size: 14px;
  line-height: 1.75;
  color: var(--landing-text-secondary, rgba(0, 0, 0, 0.55));
}

.access-cta-points {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
  margin: 0 0 28px;
  padding: 0;
  list-style: none;
}

.access-cta-points li {
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 13px;
  background: #fff;
  border: 1px solid var(--landing-line, #e9e9e9);
  color: var(--landing-text, rgba(0, 0, 0, 0.75));
}

.access-cta-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
}

.access-cta-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 0 22px;
  border-radius: 4px;
  background: var(--landing-primary, #1890ff);
  color: #fff;
  font-weight: 600;
  text-decoration: none;
  transition: opacity 0.2s ease;
}

.access-cta-primary:hover {
  opacity: 0.9;
}

.access-cta-secondary {
  display: inline-flex;
  align-items: center;
  min-height: 44px;
  padding: 0 22px;
  border-radius: 4px;
  border: 1px solid var(--landing-line, #d9d9d9);
  background: #fff;
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
  font-weight: 500;
  text-decoration: none;
  transition:
    border-color 0.2s ease,
    color 0.2s ease;
}

.access-cta-secondary:hover {
  border-color: var(--landing-primary, #1890ff);
  color: var(--landing-primary, #1890ff);
}
</style>
