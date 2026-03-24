<script setup lang="ts">
import { Globe, Mail } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';
const props = defineProps<{
  currentYear: number;
  companySiteLink: string;
  githubLink: string;
  giteeLink: string;
  icp: string;
  icpLink: string;
  contactEmail: string;
}>();
const { t } = useI18n();
</script>
<template>
  <footer class="footer-shell">
    <div class="footer-grid">
      <div class="footer-brand">
        <span class="footer-title">LiteChat</span>
        <p>{{ t('landing.footer.description') }}</p>
        <div class="footer-badge">{{ t('landing.footer.builtWith') }}</div>
      </div>
      <div class="footer-column">
        <h4>{{ t('landing.footer.ecosystem') }}</h4>
        <a href="#capabilities">{{ t('landing.nav.features') }}</a>
        <a href="#workflow">{{ t('landing.nav.workflow') }}</a>
        <a href="#project-lab">{{ t('landing.nav.projects') }}</a>
        <a href="#open-source">{{ t('landing.nav.openSource') }}</a>
      </div>
      <div class="footer-column">
        <h4>{{ t('landing.footer.resources') }}</h4>
        <a :href="props.githubLink" target="_blank" rel="noreferrer">{{ t('landing.footer.github') }}</a>
        <a :href="props.giteeLink" target="_blank" rel="noreferrer">{{ t('landing.footer.gitee') }}</a>
        <a :href="props.companySiteLink" target="_blank" rel="noreferrer">{{ t('landing.footer.officialSite') }}</a>
      </div>
      <div class="footer-column">
        <h4>{{ t('landing.footer.status') }}</h4>
        <div class="status-card">
          <span>{{ t('landing.footer.statusLabel') }}</span>
          <strong>{{ t('landing.footer.statusValue') }}</strong>
        </div>
      </div>
    </div>
    <div class="footer-bottom">
      <div class="footer-bottom-meta">
        <span>© {{ props.currentYear }} LiteChat. {{ t('landing.footer.copyright') }}</span>
        <a class="icp-link" :href="props.icpLink" target="_blank" rel="noreferrer">{{ props.icp }}</a>
      </div>
      <div class="footer-icons">
        <a :href="props.companySiteLink" target="_blank" rel="noreferrer" aria-label="official site">
          <Globe :size="18" />
        </a>
        <a :href="'mailto:' + props.contactEmail" aria-label="email">
          <Mail :size="18" />
        </a>
      </div>
    </div>
  </footer>
</template>
<style scoped>
.footer-shell {
  --footer-bg: linear-gradient(180deg, rgba(8, 15, 28, 0.94), rgba(4, 10, 20, 0.98));
  --footer-border: rgba(148, 163, 184, 0.18);
  --footer-panel: rgba(255, 255, 255, 0.04);
  --footer-text: #f8fbff;
  --footer-muted: #a9b8cf;
  --footer-accent: #7dd3fc;
  position: relative;
  z-index: 1;
  margin-top: 40px;
  padding: 58px 24px 32px;
  background: var(--footer-bg);
  border-top: 1px solid var(--footer-border);
}
.footer-grid,
.footer-bottom {
  max-width: 1320px;
  margin: 0 auto;
}
.footer-grid {
  display: grid;
  grid-template-columns: 1.2fr repeat(3, minmax(0, 1fr));
  gap: 28px;
}
.footer-title {
  color: var(--footer-text);
  font-size: 1.8rem;
  font-weight: 700;
  letter-spacing: -0.05em;
}
.footer-brand p {
  max-width: 320px;
  margin: 16px 0 0;
  color: var(--footer-muted);
  line-height: 1.8;
}
.footer-badge {
  display: inline-flex;
  align-items: center;
  min-height: 38px;
  margin-top: 22px;
  padding: 0 14px;
  border-radius: 999px;
  background: var(--footer-panel);
  border: 1px solid var(--footer-border);
  color: var(--footer-accent);
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.16em;
}
.footer-column {
  display: grid;
  gap: 12px;
  align-content: start;
}
.footer-column h4 {
  margin: 0 0 6px;
  color: var(--footer-text);
  font-size: 0.92rem;
}
.footer-column a,
.footer-icons a,
.icp-link {
  color: var(--footer-muted);
  text-decoration: none;
  transition: color 0.2s ease;
}
.footer-column a:hover,
.footer-icons a:hover,
.icp-link:hover {
  color: var(--footer-text);
}
.status-card {
  display: grid;
  gap: 12px;
  padding: 18px;
  border-radius: 16px;
  background: var(--footer-panel);
  border: 1px solid var(--footer-border);
}
.status-card span {
  color: var(--footer-muted);
  font-size: 0.82rem;
}
.status-card strong {
  color: var(--footer-text);
}
.footer-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 46px;
  padding-top: 24px;
  border-top: 1px solid var(--footer-border);
  color: var(--footer-muted);
  font-size: 0.92rem;
}
.footer-bottom-meta {
  display: grid;
  gap: 6px;
}
.footer-icons {
  display: flex;
  align-items: center;
  gap: 16px;
}
@media (max-width: 960px) {
  .footer-grid {
    grid-template-columns: 1fr;
  }
  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>