<script setup lang="ts">
import { computed } from 'vue';
import { ExternalLink } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  githubLink: string;
}>();

const { t, tm } = useI18n();
const bullets = computed(() => tm('landing.openSource.bullets') as string[]);
</script>

<template>
  <div class="section-shell open-source-shell">
    <div class="open-source-panel">
      <span>{{ t('landing.openSource.eyebrow') }}</span>
      <h2>{{ t('landing.openSource.title') }}</h2>
      <p>{{ t('landing.openSource.description') }}</p>
    </div>

    <div class="open-source-side">
      <ul>
        <li v-for="item in bullets" :key="item">{{ item }}</li>
      </ul>

      <a :href="props.githubLink" target="_blank" rel="noreferrer" class="roadmap-link">
        {{ t('landing.openSource.cta') }}
        <ExternalLink :size="18" />
      </a>
    </div>
  </div>
</template>

<style scoped>
.open-source-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(280px, 0.9fr);
  gap: 22px;
  align-items: center;
}

.open-source-panel,
.open-source-side {
  padding: 28px;
  border-radius: 32px;
  background: rgba(255, 255, 255, 0.56);
  border: 1px solid rgba(226, 192, 204, 0.56);
}

.open-source-panel span {
  display: inline-block;
  color: var(--accent-strong);
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.open-source-panel h2 {
  margin: 18px 0 0;
  color: var(--text-strong);
  font-size: clamp(2.4rem, 4vw, 4rem);
  line-height: 1;
  letter-spacing: -0.05em;
}

.open-source-panel p {
  margin: 18px 0 0;
  color: var(--text-muted);
  line-height: 1.9;
}

.open-source-side ul {
  display: grid;
  gap: 12px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.open-source-side li {
  position: relative;
  padding-left: 18px;
  color: var(--text-soft);
}

.open-source-side li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 10px;
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--accent-strong);
}

.roadmap-link {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 52px;
  margin-top: 28px;
  padding: 0 20px;
  border-radius: 999px;
  color: var(--text-strong);
  text-decoration: none;
  font-weight: 700;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(228, 193, 206, 0.7);
}

:global(.theme-dark) .open-source-panel,
:global(.theme-dark) .open-source-side,
:global(.theme-dark) .roadmap-link {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.12);
}

@media (max-width: 960px) {
  .open-source-shell {
    grid-template-columns: 1fr;
  }
}
</style>