<script setup lang="ts">
import { computed } from 'vue';
import { ArrowRight } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';
const props = defineProps<{
  repoLink: string;
}>();
const { t, tm } = useI18n();
const featured = computed(
  () => tm('landing.projectLab.featured') as {
    tag: string;
    title: string;
    description: string;
    bullets: string[];
    promptLabel: string;
    repoLabel: string;
  },
);
const backlog = computed(
  () => tm('landing.projectLab.backlog') as Array<{ title: string; description: string }>,
);
</script>
<template>
  <div class="section-shell project-shell">
    <div class="section-header">
      <span>{{ t('landing.projectLab.eyebrow') }}</span>
      <h2>{{ t('landing.projectLab.title') }}</h2>
      <p>{{ t('landing.projectLab.description') }}</p>
    </div>
    <div class="project-grid">
      <article class="featured-project">
        <div class="featured-tag">{{ featured.tag }}</div>
        <h3>{{ featured.title }}</h3>
        <p>{{ featured.description }}</p>
        <ul>
          <li v-for="item in featured.bullets" :key="item">{{ item }}</li>
        </ul>
        <div class="project-prompt">
          <span>{{ featured.promptLabel }}</span>
          <ArrowRight :size="18" />
        </div>
        <a class="project-repo" :href="props.repoLink" target="_blank" rel="noreferrer">
          {{ featured.repoLabel }}
          <ArrowRight :size="18" />
        </a>
      </article>
      <div class="project-side">
        <article v-for="item in backlog" :key="item.title" class="future-card">
          <h4>{{ item.title }}</h4>
          <p>{{ item.description }}</p>
        </article>
      </div>
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
.project-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(280px, 0.75fr);
  gap: 20px;
  margin-top: 40px;
}
.featured-project,
.future-card {
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.56);
  border: 1px solid rgba(148, 163, 184, 0.22);
}
.featured-project {
  padding: 30px;
  min-height: 540px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.84), rgba(239, 246, 255, 0.74)),
    radial-gradient(circle at top right, rgba(96, 165, 250, 0.2), transparent 34%);
}
.featured-tag {
  display: inline-flex;
  align-items: center;
  min-height: 38px;
  padding: 0 14px;
  border-radius: 999px;
  background: rgba(37, 99, 235, 0.12);
  color: var(--accent-strong);
  font-size: 0.82rem;
  font-weight: 700;
}
.featured-project h3 {
  margin: 22px 0 0;
  color: var(--text-strong);
  font-size: clamp(2rem, 3vw, 3rem);
  line-height: 1.04;
  letter-spacing: -0.05em;
}
.featured-project p {
  margin: 18px 0 0;
  max-width: 640px;
  color: var(--text-muted);
  line-height: 1.9;
}
.featured-project ul {
  display: grid;
  gap: 12px;
  margin: 28px 0 0;
  padding: 0;
  list-style: none;
}
.featured-project li {
  position: relative;
  padding-left: 18px;
  color: var(--text-soft);
}
.featured-project li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 10px;
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--accent-strong);
}
.project-prompt {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 36px;
  padding: 18px 20px;
  border-radius: 18px;
  background: rgba(248, 250, 252, 0.92);
  border: 1px solid rgba(148, 163, 184, 0.28);
  color: var(--text-strong);
  font-weight: 700;
}
.project-repo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 48px;
  margin-top: 18px;
  padding: 0 18px;
  border-radius: 14px;
  color: #eff6ff;
  text-decoration: none;
  font-weight: 700;
  background: linear-gradient(135deg, #1d4ed8 0%, #38bdf8 100%);
  box-shadow: 0 18px 36px rgba(37, 99, 235, 0.2);
}
.project-side {
  display: grid;
  gap: 14px;
}
.future-card {
  padding: 22px;
}
.future-card h4 {
  margin: 0;
  color: var(--text-strong);
  font-size: 1.08rem;
}
.future-card p {
  margin: 12px 0 0;
  color: var(--text-soft);
  line-height: 1.75;
}
:global(.theme-dark) .featured-project,
:global(.theme-dark) .future-card,
:global(.theme-dark) .project-prompt {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(148, 163, 184, 0.18);
}
@media (max-width: 960px) {
  .project-grid {
    grid-template-columns: 1fr;
  }
}
</style>