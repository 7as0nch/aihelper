<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import type { SupportedLocale } from '@/i18n/types';
import { useLocaleStore } from '@/stores/locale';
import PublicSiteHeader from '@/components/landing/PublicSiteHeader.vue';
import LandingBackToTop from '@/components/landing/LandingBackToTop.vue';
import LandingFooter from '@/components/landing/LandingFooter.vue';
import LandingBanner2 from '@/components/landing/ant/LandingBanner2.vue';
import LandingBanner4 from '@/components/landing/ant/LandingBanner4.vue';
import LandingContent0 from '@/components/landing/ant/LandingContent0.vue';
import LandingAccessCta from '@/components/landing/ant/LandingAccessCta.vue';
import ScrollRevealSection from '@/components/landing/ScrollRevealSection.vue';

type SectionId = 'products' | 'access';

interface Content0Row {
  icon?: string;
  /** 公开路径列表，用于卡片内轮播；无则留白 */
  images?: string[];
  layout?: 'default' | 'split' | 'split-browser' | 'coming-soon';
  title: string;
  content: string;
}

interface TemplateCopy {
  banner2: {
    title: string;
    content: string;
    button: string;
    backgroundImage?: string;
  };
  banner4: {
    title: string;
    content: string;
    button: string;
    image?: string;
  };
  content0: {
    title: string;
    items: Content0Row[];
  };
  cta: {
    eyebrow: string;
    title: string;
    description: string;
    points: string[];
    primary: string;
    secondary: string;
  };
}

const currentYear = new Date().getFullYear();
const companySiteLink = 'https://aihelper.chat/';
const githubLink = 'https://github.com/7as0nch/aihelper.git';
const giteeLink = 'https://gitee.com/jas0nch/aichat.git';
const icp = '蜀ICP备2026011723号-1';
const icpLink = 'https://beian.miit.gov.cn';
const contactEmail = '7as0nch@gmail.com';

const route = useRoute();
const localeStore = useLocaleStore();
const { t, tm } = useI18n();
const activeSection = ref<SectionId>('products');
const tpl = computed(() => tm('landing.template') as TemplateCopy);
const homePath = computed(() => (route.path === '/info' ? '/info' : '/'));

const navItems = computed(() => [
  {
    label: t('landing.nav.products'),
    to: { path: homePath.value, hash: '#products' },
    active: activeSection.value === 'products',
  },
  {
    label: t('landing.nav.apply'),
    to: '/apply',
    active: activeSection.value === 'access',
  },
]);

const updateLocale = (nextLocale: SupportedLocale) => {
  localeStore.applyLocale(nextLocale);
};

let sectionObserver: IntersectionObserver | null = null;

onMounted(() => {
  sectionObserver = new IntersectionObserver(
    (entries) => {
      const visible = entries
        .filter((entry) => entry.isIntersecting)
        .sort((left, right) => right.intersectionRatio - left.intersectionRatio)[0];

      if (visible?.target instanceof HTMLElement) {
        activeSection.value = visible.target.dataset.section as SectionId;
      }
    },
    {
      threshold: [0.25, 0.45, 0.65],
      rootMargin: '-12% 0px -30% 0px',
    },
  );

  document.querySelectorAll<HTMLElement>('[data-section]').forEach((element) => {
    sectionObserver?.observe(element);
  });
});

onBeforeUnmount(() => {
  sectionObserver?.disconnect();
  sectionObserver = null;
});
</script>

<template>
  <div class="landing-ant site-page">
    <PublicSiteHeader
      variant="ant"
      :items="navItems"
      :model-value="localeStore.locale"
      :action-label="t('landing.nav.apply')"
      action-to="/apply"
      @update:model-value="updateLocale"
    />

    <main class="site-main">
      <LandingBanner2
        :title="tpl.banner2.title"
        :content="tpl.banner2.content"
        :button="tpl.banner2.button"
        :background-image="tpl.banner2.backgroundImage"
      />

      <ScrollRevealSection stagger>
        <LandingBanner4
          :title="tpl.banner4.title"
          :content="tpl.banner4.content"
          :button="tpl.banner4.button"
          :image="tpl.banner4.image"
        />
      </ScrollRevealSection>

      <ScrollRevealSection stagger>
        <LandingContent0 :title="tpl.content0.title" :items="tpl.content0.items" />
      </ScrollRevealSection>

      <ScrollRevealSection stagger>
        <LandingAccessCta
          :home-path="homePath"
          :eyebrow="tpl.cta.eyebrow"
          :title="tpl.cta.title"
          :description="tpl.cta.description"
          :points="tpl.cta.points"
          :primary="tpl.cta.primary"
          :secondary="tpl.cta.secondary"
        />
      </ScrollRevealSection>
    </main>

    <LandingBackToTop />

    <LandingFooter
      variant="ant"
      :home-path="homePath"
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

.site-page {
  min-height: 100vh;
  /* 与固定顶栏 min-height 一致，避免多出几像素露出 #ececec 的「缝」 */
  padding-top: var(--landing-header-height, 64px);
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
  background: var(--landing-page-bg, #ececec);
  font-family:
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    'Helvetica Neue',
    Arial,
    'Noto Sans',
    'PingFang SC',
    'Hiragino Sans GB',
    sans-serif;
}

.site-main {
  display: block;
}
</style>
