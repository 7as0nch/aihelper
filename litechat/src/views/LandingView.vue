<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';
import { Mail, Moon, Sun } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';
import AmbientCanvas from '@/components/landing/AmbientCanvas.vue';
import LandingCapabilities from '@/components/landing/LandingCapabilities.vue';
import LandingFooter from '@/components/landing/LandingFooter.vue';
import LandingHero from '@/components/landing/LandingHero.vue';
import LandingOpenSource from '@/components/landing/LandingOpenSource.vue';
import LandingProjectLab from '@/components/landing/LandingProjectLab.vue';
import LandingWorkflow from '@/components/landing/LandingWorkflow.vue';
import LanguageSwitch from '@/components/landing/LanguageSwitch.vue';
import type { SupportedLocale } from '@/i18n/types';
import { useLocaleStore } from '@/stores/locale';
import { useThemeStore } from '@/stores/theme';
const currentYear = new Date().getFullYear();
const companySiteLink = 'https://aihelper.chat/';
const githubLink = 'https://github.com/7as0nch/aihelper.git';
const giteeLink = 'https://gitee.com/jas0nch/aichat.git';
const icp = '蜀ICP备2026011723号-1';
const icpLink = 'https://beian.miit.gov.cn';
const contactEmail = '7as0nch@gmail.com';
const aicookLink = 'https://github.com/7as0nch/aicook';
type SectionId = 'hero' | 'capabilities' | 'workflow' | 'project-lab' | 'open-source';
const themeStore = useThemeStore();
const localeStore = useLocaleStore();
const { t } = useI18n();
const activeSection = ref<SectionId>('hero');
const scrollProgress = ref(0);
const sections = computed(() => [
  { id: 'capabilities' as SectionId, label: t('landing.nav.features') },
  { id: 'workflow' as SectionId, label: t('landing.nav.workflow') },
  { id: 'project-lab' as SectionId, label: t('landing.nav.projects') },
  { id: 'open-source' as SectionId, label: t('landing.nav.openSource') },
]);
const resolvedTheme = computed<'light' | 'dark'>(() => {
  return themeStore.mode === 'dark' ? 'dark' : 'light';
});
const themeClass = computed(() => 'theme-' + resolvedTheme.value);
const handleThemeToggle = () => {
  themeStore.setMode(resolvedTheme.value === 'dark' ? 'light' : 'dark');
};
const updateLocale = (nextLocale: SupportedLocale) => {
  localeStore.applyLocale(nextLocale);
};
let observer: IntersectionObserver | null = null;
let frameId = 0;
const updateProgress = () => {
  const totalHeight = document.documentElement.scrollHeight - window.innerHeight;
  scrollProgress.value = totalHeight > 0 ? window.scrollY / totalHeight : 0;
};
const handleScroll = () => {
  if (frameId) {
    return;
  }
  frameId = window.requestAnimationFrame(() => {
    updateProgress();
    frameId = 0;
  });
};
onMounted(() => {
  updateProgress();
  observer = new IntersectionObserver(
    (entries) => {
      const visible = entries
        .filter((entry) => entry.isIntersecting)
        .sort((left, right) => right.intersectionRatio - left.intersectionRatio)[0];
      if (visible?.target instanceof HTMLElement) {
        activeSection.value = visible.target.dataset.section as SectionId;
      }
    },
    {
      threshold: [0.25, 0.45, 0.7],
      rootMargin: '-10% 0px -24% 0px',
    },
  );
  document.querySelectorAll<HTMLElement>('[data-section]').forEach((element) => {
    observer?.observe(element);
  });
  window.addEventListener('scroll', handleScroll, { passive: true });
});
onBeforeUnmount(() => {
  observer?.disconnect();
  observer = null;
  if (frameId) {
    window.cancelAnimationFrame(frameId);
  }
  window.removeEventListener('scroll', handleScroll);
});
</script>
<template>
  <div class="landing-page" :class="themeClass">
    <AmbientCanvas :theme="resolvedTheme" :progress="scrollProgress" />
    <div class="landing-mesh" aria-hidden="true" />
    <header class="landing-header">
      <div class="header-inner">
        <div class="brand-block">
          <img class="brand-logo" src="/logo.png" alt="LiteChat" />
          <div>
            <span class="brand-name">LiteChat</span>
            <span class="brand-meta">AI delivery layer</span>
          </div>
        </div>
        <nav class="header-nav">
          <a
            v-for="item in sections"
            :key="item.id"
            class="nav-link"
            :class="{ 'is-active': activeSection === item.id }"
            :href="'#' + item.id"
          >
            {{ item.label }}
          </a>
          <a class="nav-link" :href="githubLink" target="_blank" rel="noreferrer">
            {{ t('landing.nav.docs') }}
          </a>
        </nav>
        <div class="header-actions">
          <a class="contact-link" :href="'mailto:' + contactEmail">
            <Mail :size="16" />
            {{ t('landing.footer.contact') }}
          </a>
          <LanguageSwitch :model-value="localeStore.locale" @update:model-value="updateLocale" />
          <button
            class="theme-toggle"
            type="button"
            :aria-label="resolvedTheme === 'dark' ? '切换到浅色主题' : '切换到暗色主题'"
            @click="handleThemeToggle"
          >
            <Sun v-if="resolvedTheme === 'dark'" :size="18" />
            <Moon v-else :size="18" />
          </button>
          <RouterLink to="/app" class="launch-link">
            {{ t('landing.nav.launch') }}
          </RouterLink>
        </div>
      </div>
    </header>
    <aside class="section-indicator" aria-label="page progress">
      <a
        v-for="item in sections"
        :key="item.id"
        class="indicator-dot"
        :class="{ 'is-active': activeSection === item.id }"
        :href="'#' + item.id"
      >
        <span>{{ item.label }}</span>
      </a>
    </aside>
    <main class="landing-main">
      <section id="hero" data-section="hero" class="landing-section hero-section is-active">
        <div class="section-inner hero-inner">
          <LandingHero :github-link="githubLink" />
        </div>
      </section>
      <section
        id="capabilities"
        data-section="capabilities"
        class="landing-section"
        :class="{ 'is-active': activeSection === 'capabilities' }"
      >
        <div class="section-inner">
          <LandingCapabilities />
        </div>
      </section>
      <section
        id="workflow"
        data-section="workflow"
        class="landing-section accent-section"
        :class="{ 'is-active': activeSection === 'workflow' }"
      >
        <div class="section-inner">
          <LandingWorkflow />
        </div>
      </section>
      <section
        id="project-lab"
        data-section="project-lab"
        class="landing-section"
        :class="{ 'is-active': activeSection === 'project-lab' }"
      >
        <div class="section-inner">
          <LandingProjectLab :repo-link="aicookLink" />
        </div>
      </section>
      <section
        id="open-source"
        data-section="open-source"
        class="landing-section accent-section"
        :class="{ 'is-active': activeSection === 'open-source' }"
      >
        <div class="section-inner">
          <LandingOpenSource :github-link="githubLink" />
        </div>
      </section>
    </main>
    <LandingFooter
      :current-year="currentYear"
      :company-site-link="companySiteLink"
      :github-link="githubLink"
      :gitee-link="giteeLink"
      :icp="icp"
      :icp-link="icpLink"
      :contact-email="contactEmail"
    />
  </div>
</template>
<style scoped>
:global(html) {
  scroll-behavior: smooth;
}
.landing-page {
  --page-bg: #f5f7fb;
  --page-bg-secondary: #fbfcfe;
  --surface-strong: rgba(255, 255, 255, 0.72);
  --border-soft: rgba(148, 163, 184, 0.26);
  --text-strong: #0f172a;
  --text-soft: #475569;
  --text-muted: #64748b;
  --accent-strong: #2563eb;
  position: relative;
  min-height: 100vh;
  color: var(--text-strong);
  background:
    radial-gradient(circle at top left, rgba(96, 165, 250, 0.18), transparent 32%),
    linear-gradient(180deg, var(--page-bg) 0%, var(--page-bg-secondary) 62%, #f2f6fc 100%);
  font-family: 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  overflow-x: hidden;
}
.theme-dark {
  --page-bg: #06111f;
  --page-bg-secondary: #0a1425;
  --surface-strong: rgba(15, 23, 42, 0.62);
  --border-soft: rgba(148, 163, 184, 0.22);
  --text-strong: #f8fbff;
  --text-soft: #dbe7f6;
  --text-muted: #a9b8cf;
  --accent-strong: #7dd3fc;
  background:
    radial-gradient(circle at top left, rgba(56, 189, 248, 0.18), transparent 34%),
    linear-gradient(180deg, var(--page-bg) 0%, var(--page-bg-secondary) 65%, #081221 100%);
}
.landing-mesh {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.4) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.4) 1px, transparent 1px);
  background-size: 52px 52px;
  mask-image: linear-gradient(to bottom, rgba(0, 0, 0, 0.28), transparent 76%);
  opacity: 0.1;
}
.theme-dark .landing-mesh {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.08) 1px, transparent 1px);
}
.landing-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 30;
  padding: 18px 24px;
}
.header-inner {
  max-width: 1320px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  min-height: 72px;
  padding: 0 16px;
  border-radius: 14px;
  background: rgba(248, 250, 252, 0.82);
  border: 1px solid var(--border-soft);
  backdrop-filter: blur(22px);
  box-shadow: 0 20px 44px rgba(15, 23, 42, 0.08);
}
.theme-dark .header-inner {
  background: rgba(8, 16, 30, 0.9);
  box-shadow: 0 18px 38px rgba(0, 0, 0, 0.28);
}
.brand-block {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.brand-logo {
  width: 42px;
  height: 42px;
  border-radius: 8px;
  object-fit: cover;
}
.brand-name,
.brand-meta {
  display: block;
}
.brand-name {
  color: var(--text-strong);
  font-size: 1.05rem;
  font-weight: 800;
  letter-spacing: -0.03em;
}
.brand-meta {
  margin-top: 3px;
  color: var(--text-muted);
  font-size: 0.76rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
.header-nav,
.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
.nav-link {
  padding: 10px 12px;
  border-radius: 12px;
  color: var(--text-muted);
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.25s ease, background 0.25s ease;
}
.nav-link.is-active,
.nav-link:hover {
  color: var(--text-strong);
  background: rgba(37, 99, 235, 0.08);
}
.theme-dark .nav-link.is-active,
.theme-dark .nav-link:hover {
  background: rgba(125, 211, 252, 0.12);
}
.theme-toggle,
.launch-link,
.contact-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 46px;
  border-radius: 12px;
  border: 0;
  text-decoration: none;
}
.theme-toggle {
  width: 46px;
  color: var(--text-strong);
  cursor: pointer;
  background: var(--surface-strong);
  border: 1px solid var(--border-soft);
}
.contact-link {
  gap: 8px;
  padding: 0 16px;
  color: var(--text-strong);
  background: var(--surface-strong);
  border: 1px solid var(--border-soft);
}
.launch-link {
  padding: 0 18px;
  color: #eff6ff;
  background: linear-gradient(135deg, #2563eb 0%, #38bdf8 100%);
  font-weight: 700;
  box-shadow: 0 18px 36px rgba(37, 99, 235, 0.24);
}
.section-indicator {
  position: fixed;
  top: 50%;
  right: 28px;
  z-index: 20;
  display: grid;
  gap: 10px;
  transform: translateY(-50%);
}
.indicator-dot {
  display: flex;
  justify-content: flex-end;
  text-decoration: none;
}
.indicator-dot span {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 14px;
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.68);
  border: 1px solid transparent;
  color: transparent;
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  transition: color 0.28s ease, background 0.28s ease, border-color 0.28s ease;
}
.indicator-dot.is-active span {
  color: var(--text-strong);
  background: rgba(37, 99, 235, 0.08);
  border-color: var(--border-soft);
}
.theme-dark .indicator-dot span {
  background: rgba(255, 255, 255, 0.06);
}
.landing-main {
  position: relative;
  z-index: 1;
  padding-top: 96px;
}
.landing-section {
  padding: 44px 24px 88px;
}
.hero-section {
  min-height: calc(100svh - 24px);
  display: flex;
  align-items: center;
}
.section-inner {
  width: min(1320px, 100%);
  margin: 0 auto;
}
.accent-section .section-inner {
  position: relative;
}
.accent-section .section-inner::before {
  content: '';
  position: absolute;
  inset: -18px;
  z-index: -1;
  border-radius: 26px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.38), rgba(226, 232, 240, 0.16));
  border: 1px solid rgba(148, 163, 184, 0.16);
}
.theme-dark .accent-section .section-inner::before {
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.24), rgba(15, 23, 42, 0.08));
}
.landing-section :deep(.section-shell) {
  opacity: 0.76;
  transform: translateY(24px) scale(0.99);
  filter: saturate(0.94);
  transition: opacity 0.7s ease, transform 0.7s ease, filter 0.7s ease;
}
.hero-section :deep(.section-shell),
.landing-section.is-active :deep(.section-shell) {
  opacity: 1;
  transform: translateY(0) scale(1);
  filter: none;
}
@media (max-width: 1180px) {
  .section-indicator {
    display: none;
  }
}
@media (max-width: 920px) {
  .landing-header {
    padding: 16px;
  }
  .header-inner {
    min-height: auto;
    border-radius: 12px;
    padding: 14px;
    flex-wrap: wrap;
  }
  .header-nav {
    order: 3;
    width: 100%;
    overflow-x: auto;
  }
  .landing-section {
    padding: 30px 16px 72px;
  }
}
@media (max-width: 720px) {
  .header-nav {
    display: none;
  }
  .header-actions {
    width: 100%;
    justify-content: space-between;
    flex-wrap: wrap;
  }
  .launch-link {
    padding: 0 16px;
  }
  .hero-section {
    min-height: auto;
    padding-top: 18px;
  }
}
</style>