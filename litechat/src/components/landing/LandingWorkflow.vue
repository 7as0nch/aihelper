<script setup lang="ts">
import { computed, type Component } from 'vue';
import { BrainCircuit, LayoutDashboard, Sparkles } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

const { t, tm } = useI18n();

const icons: Component[] = [BrainCircuit, LayoutDashboard, Sparkles];
const panels = computed(
  () => tm('landing.workflow.panels') as Array<{ title: string; description: string }>,
);
const callouts = computed(() => tm('landing.workflow.callouts') as string[]);
</script>

<template>
  <div class="section-shell workflow-shell">
    <div class="workflow-copy">
      <span>{{ t('landing.workflow.eyebrow') }}</span>
      <h2>{{ t('landing.workflow.title') }}</h2>
      <p>{{ t('landing.workflow.description') }}</p>

      <ul class="workflow-callouts">
        <li v-for="item in callouts" :key="item">{{ item }}</li>
      </ul>
    </div>

    <div class="workflow-visual">
      <div class="workflow-rail" />
      <article v-for="(panel, index) in panels" :key="panel.title" class="workflow-panel">
        <div class="workflow-icon">
          <component :is="icons[index]" :size="22" />
        </div>
        <div>
          <h3>{{ panel.title }}</h3>
          <p>{{ panel.description }}</p>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.workflow-shell {
  display: grid;
  grid-template-columns: minmax(0, 0.95fr) minmax(0, 1.05fr);
  gap: 34px;
  align-items: start;
}

.workflow-copy {
  position: sticky;
  top: 112px;
}

.workflow-copy span {
  display: inline-block;
  color: var(--accent-strong);
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.workflow-copy h2 {
  margin: 18px 0 0;
  color: var(--text-strong);
  font-size: clamp(2.4rem, 4vw, 4rem);
  line-height: 1;
  letter-spacing: -0.05em;
}

.workflow-copy p {
  margin: 18px 0 0;
  color: var(--text-muted);
  line-height: 1.9;
}

.workflow-callouts {
  display: grid;
  gap: 10px;
  margin: 28px 0 0;
  padding: 0;
  list-style: none;
}

.workflow-callouts li {
  color: var(--text-soft);
}

.workflow-visual {
  position: relative;
  display: grid;
  gap: 18px;
  padding-left: 20px;
}

.workflow-rail {
  position: absolute;
  top: 8px;
  bottom: 8px;
  left: 11px;
  width: 2px;
  background: linear-gradient(180deg, rgba(37, 99, 235, 0.76), rgba(125, 211, 252, 0.12));
}

.workflow-panel {
  position: relative;
  display: grid;
  grid-template-columns: 62px minmax(0, 1fr);
  gap: 16px;
  padding: 22px;
  border-radius: 28px;
  background: rgba(255, 255, 255, 0.54);
  border: 1px solid rgba(148, 163, 184, 0.22);
}

.workflow-panel::before {
  content: '';
  position: absolute;
  left: -14px;
  top: 34px;
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: var(--accent-strong);
  box-shadow: 0 0 0 8px rgba(37, 99, 235, 0.1);
}

.workflow-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 18px;
  color: var(--accent-strong);
  background: rgba(37, 99, 235, 0.1);
}

.workflow-panel h3 {
  margin: 0;
  color: var(--text-strong);
  font-size: 1.16rem;
}

.workflow-panel p {
  margin: 12px 0 0;
  color: var(--text-soft);
  line-height: 1.75;
}

:global(.theme-dark) .workflow-panel {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(148, 163, 184, 0.18);
}

@media (max-width: 960px) {
  .workflow-shell {
    grid-template-columns: 1fr;
  }

  .workflow-copy {
    position: static;
  }
}
</style>