<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { useI18n } from 'vue-i18n';

const props = withDefaults(
  defineProps<{
    currentYear: number;
    companySiteLink: string;
    githubLink: string;
    giteeLink: string;
    icp: string;
    icpLink: string;
    contactEmail: string;
    homePath?: string;
    variant?: 'default' | 'ant';
  }>(),
  { homePath: '/', variant: 'ant' },
);

const { t } = useI18n();

const mailtoHref = `mailto:${props.contactEmail}`;
</script>

<template>
  <footer class="footer-root" :class="{ 'footer-root--ant': props.variant === 'ant' }">
    <div v-if="props.variant === 'ant'" class="footer-landing">
      <div class="footer-landing-inner">
        <div class="footer-grid">
          <!-- Brand -->
          <div class="footer-col footer-col-brand">
            <RouterLink :to="props.homePath" class="footer-brand-row" @click.stop>
              <span class="footer-brand-mark">
                <img src="/logo.png" :alt="t('landing.brand.name')" />
              </span>
              <span class="footer-brand-name">{{ t('landing.brand.name') }}</span>
            </RouterLink>
            <p class="footer-brand-desc">{{ t('landing.footer.brandDescription') }}</p>
          </div>

          <!-- Products -->
          <div class="footer-col">
            <h4 class="footer-col-title">{{ t('landing.footer.columnProducts') }}</h4>
            <ul class="footer-link-list">
              <li>
                <RouterLink :to="{ path: props.homePath, hash: '#products' }">
                  {{ t('landing.footer.linkProductIntro') }}
                </RouterLink>
              </li>
              <li>
                <RouterLink to="/apply">{{ t('landing.footer.linkBeta') }}</RouterLink>
              </li>
            </ul>
          </div>

          <!-- About -->
          <div class="footer-col">
            <h4 class="footer-col-title">{{ t('landing.footer.columnAbout') }}</h4>
            <ul class="footer-link-list">
              <li>
                <RouterLink :to="{ path: props.homePath, hash: '#products' }">
                  {{ t('landing.footer.linkFaq') }}
                </RouterLink>
              </li>
              <li>
                <a :href="mailtoHref">{{ t('landing.footer.linkContact') }}</a>
              </li>
            </ul>
          </div>

          <!-- Resources -->
          <div class="footer-col">
            <h4 class="footer-col-title">{{ t('landing.footer.columnResources') }}</h4>
            <ul class="footer-link-list">
              <li>
                <a :href="props.githubLink" target="_blank" rel="noreferrer">{{ t('landing.footer.github') }}</a>
              </li>
              <li>
                <a :href="props.giteeLink" target="_blank" rel="noreferrer">{{ t('landing.footer.gitee') }}</a>
              </li>
              <li>
                <a :href="props.companySiteLink" target="_blank" rel="noreferrer">
                  {{ t('landing.footer.officialSite') }}
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div class="footer-divider" />

        <div class="footer-bottom">
          <p class="footer-copyright">
            © {{ props.currentYear }} {{ t('landing.brand.name') }} · {{ t('landing.footer.rightsReserved') }}
            <a class="footer-icp" :href="props.icpLink" target="_blank" rel="noreferrer">{{ props.icp }}</a>
          </p>
        </div>
      </div>
    </div>

    <div v-else class="footer-fallback">
      <p class="footer-fallback-meta">
        © {{ props.currentYear }} {{ t('landing.brand.name') }}
        <a :href="props.icpLink" target="_blank" rel="noreferrer">{{ props.icp }}</a>
      </p>
    </div>
  </footer>
</template>

<style scoped>
.footer-root--ant {
  margin-top: 0;
}

.footer-landing {
  background: var(--landing-nav-bg, #001529);
  color: rgba(255, 255, 255, 0.65);
  padding: 64px 24px 28px;
}

.footer-landing-inner {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-grid {
  display: grid;
  grid-template-columns: 1.4fr repeat(3, minmax(0, 1fr));
  gap: 40px 32px;
  align-items: flex-start;
}

.footer-col-title {
  margin: 0 0 20px;
  font-size: 15px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.92);
  letter-spacing: 0.02em;
}

.footer-brand-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  text-decoration: none;
  color: #fff;
  transition: opacity 0.2s ease;
}

.footer-brand-row:hover {
  opacity: 0.88;
}

.footer-brand-mark {
  display: flex;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.footer-brand-mark img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.footer-brand-name {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.footer-brand-desc {
  margin: 0;
  max-width: 280px;
  font-size: 13px;
  line-height: 1.75;
  color: rgba(255, 255, 255, 0.45);
}

.footer-link-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 12px;
}

.footer-link-list a {
  font-size: 13px;
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.45);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-link-list a:hover {
  color: #fff;
}

.footer-divider {
  height: 1px;
  margin: 48px 0 20px;
  background: rgba(255, 255, 255, 0.08);
}

.footer-bottom {
  text-align: center;
  padding-bottom: 8px;
}

.footer-copyright {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.35);
}

.footer-icp {
  margin-left: 10px;
  color: rgba(255, 255, 255, 0.35);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-icp:hover {
  color: rgba(255, 255, 255, 0.75);
}

.footer-fallback {
  padding: 24px;
  text-align: center;
  background: var(--landing-page-bg, #ececec);
  font-size: 13px;
  color: rgba(0, 0, 0, 0.45);
}

.footer-fallback-meta a {
  margin-left: 8px;
  color: inherit;
}

@media (max-width: 991px) {
  .footer-grid {
    grid-template-columns: 1fr 1fr;
  }

  .footer-col-brand {
    grid-column: 1 / -1;
  }
}

@media (max-width: 575px) {
  .footer-landing {
    padding: 48px 16px 24px;
  }

  .footer-grid {
    grid-template-columns: 1fr;
    gap: 32px;
  }
}
</style>
