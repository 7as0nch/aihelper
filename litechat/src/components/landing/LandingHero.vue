<script setup lang="ts">
import { computed } from 'vue';
import { RouterLink } from 'vue-router';
import { ArrowRight, Github } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  githubLink: string;
}>();

const { t, tm } = useI18n();

const metrics = computed(
  () => tm('landing.hero.metrics') as Array<{ label: string; value: string }>,
);
const highlights = computed(() => tm('landing.hero.highlights') as string[]);
</script>

<template>
  <div class="section-shell hero-shell">
    <div class="hero-copy">
      <span class="hero-eyebrow">{{ t('landing.hero.eyebrow') }}</span>
      <div class="hero-badge">{{ t('landing.hero.badge') }}</div>

      <h1>
        {{ t('landing.hero.titleLead') }}
        <span>{{ t('landing.hero.titleHighlight') }}</span>
      </h1>

      <p class="hero-description">
        {{ t('landing.hero.description') }}
      </p>

      <div class="hero-actions">
        <RouterLink to="/app" class="primary-action">
          {{ t('landing.hero.primary') }}
          <ArrowRight :size="18" />
        </RouterLink>
        <a class="secondary-action" :href="props.githubLink" target="_blank" rel="noreferrer">
          <Github :size="18" />
          {{ t('landing.hero.secondary') }}
        </a>
      </div>

      <ul class="hero-highlights">
        <li v-for="item in highlights" :key="item">{{ item }}</li>
      </ul>
    </div>

    <div class="hero-visual">
      <div class="hero-stage">
        <div class="hero-stage-backdrop" />
        <div class="hero-main-panel">
          <div class="panel-topline">LiteChat</div>
          <div class="panel-title">Delivery-ready AI surface</div>
          <div class="panel-copy">Go orchestration, reusable frontend delivery, and scenario-driven product building.</div>
          <div class="panel-bars">
            <span />
            <span />
            <span />
          </div>
        </div>
        <div class="hero-float hero-float-top">Deep Search</div>
        <div class="hero-float hero-float-bottom">Widget Embed</div>
      </div>

      <div class="hero-metrics">
        <article v-for="item in metrics" :key="item.label" class="metric-card">
          <span>{{ item.label }}</span>
          <strong>{{ item.value }}</strong>
        </article>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hero-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr);
  gap: 48px;
  align-items: center;
}

.hero-copy {
  position: relative;
  z-index: 1;
}

.hero-eyebrow {
  display: inline-block;
  margin-bottom: 14px;
  color: var(--accent-strong);
  font-size: 0.82rem;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  min-height: 38px;
  padding: 0 16px;
  border-radius: 999px;
  margin-bottom: 22px;
  background: var(--surface-strong);
  border: 1px solid var(--border-soft);
  color: var(--text-muted);
  font-size: 0.85rem;
  font-weight: 700;
}

.hero-copy h1 {
  margin: 0;
  max-width: 720px;
  color: var(--text-strong);
  font-size: clamp(3.25rem, 8vw, 6.2rem);
  line-height: 0.94;
  letter-spacing: -0.06em;
}

.hero-copy h1 span {
  display: block;
  margin-top: 10px;
  color: var(--accent-strong);
}

.hero-description {
  max-width: 640px;
  margin: 28px 0 0;
  color: var(--text-muted);
  font-size: 1.08rem;
  line-height: 1.9;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 30px;
}

.primary-action,
.secondary-action {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 54px;
  padding: 0 22px;
  border-radius: 999px;
  text-decoration: none;
  font-weight: 700;
  transition: transform 0.25s ease, box-shadow 0.25s ease, background 0.25s ease;
}

.primary-action {
  color: #401426;
  background: linear-gradient(135deg, #2563eb 0%, #38bdf8 100%);
  box-shadow: 0 18px 40px rgba(37, 99, 235, 0.26);
}

.secondary-action {
  color: var(--text-strong);
  background: rgba(248, 250, 252, 0.76);
  border: 1px solid rgba(148, 163, 184, 0.34);
}

.primary-action:hover,
.secondary-action:hover {
  transform: translateY(-2px);
}

.hero-highlights {
  display: grid;
  gap: 10px;
  margin: 28px 0 0;
  padding: 0;
  list-style: none;
}

.hero-highlights li {
  position: relative;
  padding-left: 18px;
  color: var(--text-soft);
  line-height: 1.7;
}

.hero-highlights li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 11px;
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--accent-strong);
}

.hero-visual {
  display: grid;
  gap: 18px;
}

.hero-stage {
  position: relative;
  min-height: 520px;
  border-radius: 36px;
  padding: 28px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.86), rgba(239, 246, 255, 0.78));
  border: 1px solid rgba(148, 163, 184, 0.26);
  overflow: hidden;
  box-shadow: 0 40px 80px rgba(37, 99, 235, 0.12);
}

.hero-stage-backdrop {
  position: absolute;
  inset: 12% 8% auto auto;
  width: 260px;
  height: 260px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(96, 165, 250, 0.32), transparent 68%);
  filter: blur(8px);
}

.hero-main-panel {
  position: absolute;
  left: 28px;
  right: 84px;
  bottom: 28px;
  padding: 28px;
  border-radius: 28px;
  background: rgba(248, 250, 252, 0.78);
  border: 1px solid rgba(148, 163, 184, 0.3);
  box-shadow: 0 18px 45px rgba(37, 99, 235, 0.12);
  backdrop-filter: blur(20px);
}

.panel-topline {
  color: var(--text-muted);
  font-size: 0.74rem;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.panel-title {
  margin-top: 18px;
  color: var(--text-strong);
  font-size: 1.9rem;
  font-weight: 700;
  letter-spacing: -0.04em;
}

.panel-copy {
  margin-top: 14px;
  color: var(--text-muted);
  line-height: 1.8;
}

.panel-bars {
  display: grid;
  gap: 12px;
  margin-top: 22px;
}

.panel-bars span {
  display: block;
  height: 12px;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.82), rgba(125, 211, 252, 0.24));
}

.panel-bars span:nth-child(2) {
  width: 78%;
}

.panel-bars span:nth-child(3) {
  width: 52%;
}

.hero-float {
  position: absolute;
  display: inline-flex;
  align-items: center;
  min-height: 42px;
  padding: 0 16px;
  border-radius: 999px;
  background: rgba(248, 250, 252, 0.92);
  border: 1px solid rgba(148, 163, 184, 0.3);
  box-shadow: 0 16px 32px rgba(37, 99, 235, 0.12);
  color: var(--text-strong);
  font-size: 0.85rem;
  font-weight: 700;
}

.hero-float-top {
  top: 32px;
  right: 28px;
}

.hero-float-bottom {
  top: 124px;
  left: 28px;
}

.hero-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.metric-card {
  padding: 18px;
  border-radius: 22px;
  background: rgba(248, 250, 252, 0.74);
  border: 1px solid rgba(223, 194, 206, 0.56);
}

.metric-card span {
  display: block;
  color: var(--text-muted);
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.14em;
}

.metric-card strong {
  display: block;
  margin-top: 12px;
  color: var(--text-strong);
  font-size: 1rem;
}

:global(.theme-dark) .secondary-action {
  background: rgba(15, 23, 42, 0.74);
  border-color: rgba(148, 163, 184, 0.24);
}

:global(.theme-dark) .hero-stage {
  background: linear-gradient(180deg, rgba(10, 20, 37, 0.96), rgba(7, 16, 31, 0.96));
  border-color: rgba(96, 165, 250, 0.18);
  box-shadow: 0 40px 80px rgba(0, 0, 0, 0.24);
}

:global(.theme-dark) .hero-main-panel,
:global(.theme-dark) .metric-card,
:global(.theme-dark) .hero-float {
  background: rgba(15, 23, 42, 0.74);
  border-color: rgba(148, 163, 184, 0.22);
}

@media (max-width: 1100px) {
  .hero-shell {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .hero-stage {
    min-height: 420px;
    padding: 18px;
  }

  .hero-main-panel {
    left: 18px;
    right: 18px;
    bottom: 18px;
    padding: 20px;
  }

  .hero-metrics {
    grid-template-columns: 1fr;
  }
}
</style>