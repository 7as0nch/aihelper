<script setup lang="ts">
import { computed, type Component } from 'vue';
import {
  Camera,
  Globe,
  MessageSquare,
  Mic,
  Search,
  Workflow,
} from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const { t, tm } = useI18n();

const icons: Component[] = [MessageSquare, Search, Camera, Globe, Mic, Workflow];
const cards = computed(
  () => tm('landing.capabilities.cards') as Array<{ title: string; description: string }>,
);
</script>

<template>
  <div class="section-shell capabilities-shell">
    <div class="section-header">
      <span>{{ t('landing.capabilities.eyebrow') }}</span>
      <h2>{{ t('landing.capabilities.title') }}</h2>
      <p>{{ t('landing.capabilities.description') }}</p>
    </div>

    <div class="capabilities-grid">
      <article v-for="(card, index) in cards" :key="card.title" class="capability-card">
        <div class="capability-icon">
          <component :is="icons[index]" :size="24" />
        </div>
        <h3>{{ card.title }}</h3>
        <p>{{ card.description }}</p>
      </article>
    </div>
  </div>
</template>

<style scoped>
.section-header {
  max-width: 780px;
}

.section-header span {
  display: inline-block;
  color: var(--accent-strong);
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.section-header h2 {
  margin: 18px 0 0;
  color: var(--text-strong);
  font-size: clamp(2.4rem, 4vw, 4rem);
  line-height: 1;
  letter-spacing: -0.05em;
}

.section-header p {
  margin: 18px 0 0;
  color: var(--text-muted);
  line-height: 1.9;
}

.capabilities-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
  margin-top: 44px;
}

.capability-card {
  min-height: 240px;
  padding: 26px;
  border-radius: 28px;
  background: rgba(255, 255, 255, 0.58);
  border: 1px solid rgba(148, 163, 184, 0.22);
  box-shadow: 0 18px 38px rgba(15, 23, 42, 0.08);
  transition: transform 0.28s ease, border-color 0.28s ease;
}

.capability-card:hover {
  transform: translateY(-4px);
  border-color: rgba(37, 99, 235, 0.32);
}

.capability-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 18px;
  margin-bottom: 24px;
  color: var(--accent-strong);
  background: rgba(37, 99, 235, 0.1);
}

.capability-card h3 {
  margin: 0;
  color: var(--text-strong);
  font-size: 1.18rem;
}

.capability-card p {
  margin: 14px 0 0;
  color: var(--text-soft);
  line-height: 1.75;
}

:global(.theme-dark) .capability-card {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(148, 163, 184, 0.18);
  box-shadow: none;
}

@media (max-width: 1080px) {
  .capabilities-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .capabilities-grid {
    grid-template-columns: 1fr;
  }
}
</style>