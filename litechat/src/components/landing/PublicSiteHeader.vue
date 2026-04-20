<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { RouterLink, type RouteLocationRaw } from 'vue-router';
import { Layout, Menu, Button, Dropdown, Space } from 'ant-design-vue';
import { GlobalOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';
import type { SupportedLocale } from '@/i18n/types';
import LanguageSwitch from '@/components/landing/LanguageSwitch.vue';

interface NavItem {
  label: string;
  to: RouteLocationRaw;
  active?: boolean;
}

const props = withDefaults(
  defineProps<{
    items: NavItem[];
    modelValue: SupportedLocale;
    actionLabel: string;
    actionTo: RouteLocationRaw;
    variant?: 'default' | 'ant';
  }>(),
  { variant: 'default' },
);

const emit = defineEmits<{
  (event: 'update:modelValue', value: SupportedLocale): void;
}>();

const { t } = useI18n();
/** Nav0：mobile 下 .home-page.open 展开 inline 菜单 */
const inlineNavOpen = ref(false);
/** 滚动后顶栏变为两侧缩进的胶囊浮动条 */
const headerFloating = ref(false);
const FLOAT_SCROLL_THRESHOLD = 12;

let scrollRaf = 0;
const onWindowScroll = () => {
  if (scrollRaf) {
    return;
  }
  scrollRaf = requestAnimationFrame(() => {
    scrollRaf = 0;
    const next = window.scrollY > FLOAT_SCROLL_THRESHOLD;
    if (headerFloating.value !== next) {
      headerFloating.value = next;
    }
  });
};

onMounted(() => {
  onWindowScroll();
  window.addEventListener('scroll', onWindowScroll, { passive: true });
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onWindowScroll);
  if (scrollRaf) {
    cancelAnimationFrame(scrollRaf);
    scrollRaf = 0;
  }
});

const updateLocale = (locale: SupportedLocale) => {
  emit('update:modelValue', locale);
};

const selectedNavKeys = computed(() => {
  const i = props.items.findIndex((x) => x.active);
  return i >= 0 ? [`nav-${i}`] : [];
});

const localeLabel = computed(() => (props.modelValue === 'zh-CN' ? t('landing.language.zh') : t('landing.language.en')));

const onLocaleMenuClick = (info: { key: string | number }) => {
  const key = String(info.key) as SupportedLocale;
  if (key === 'zh-CN' || key === 'en-US') {
    updateLocale(key);
  }
};

const toggleInlineNav = () => {
  inlineNavOpen.value = !inlineNavOpen.value;
};

const closeInlineNav = () => {
  inlineNavOpen.value = false;
};

watch(
  () => props.items,
  () => {
    inlineNavOpen.value = false;
  },
);
</script>

<template>
  <!-- Ant landing: Nav0 风格 — 移动端 home-page.open + inline Menu，无 Drawer -->
  <Layout.Header
    v-if="props.variant === 'ant'"
    class="header-ant"
    :class="{ 'header-ant--floating': headerFloating }"
  >
    <div
      class="header-ant-bar home-page"
      :class="{ open: inlineNavOpen }"
    >
      <RouterLink to="/" class="header-ant-brand" @click="closeInlineNav">
        <span class="header-ant-logo">
          <img src="/logo.png" :alt="t('landing.brand.name')" />
        </span>
        <strong class="header-ant-title">{{ t('landing.brand.name') }}</strong>
      </RouterLink>

      <Menu
        class="header-ant-menu-desktop"
        theme="dark"
        mode="horizontal"
        :selected-keys="selectedNavKeys"
      >
        <Menu.Item v-for="(item, i) in props.items" :key="`nav-${i}`">
          <RouterLink :to="item.to">{{ item.label }}</RouterLink>
        </Menu.Item>
      </Menu>

      <Space class="header-ant-actions" :size="8" align="center">
        <Dropdown :trigger="['click']">
          <Button type="text" ghost class="header-ant-lang-btn">
            <GlobalOutlined />
            <span>{{ localeLabel }}</span>
          </Button>
          <template #overlay>
            <Menu theme="dark" :selected-keys="[props.modelValue]" @click="onLocaleMenuClick">
              <Menu.Item key="zh-CN">{{ t('landing.language.zh') }}</Menu.Item>
              <Menu.Item key="en-US">{{ t('landing.language.en') }}</Menu.Item>
            </Menu>
          </template>
        </Dropdown>

        <RouterLink v-slot="{ navigate }" :to="props.actionTo" custom>
          <Button
            type="primary"
            @click="
              () => {
                navigate();
                closeInlineNav();
              }
            "
          >
            {{ props.actionLabel }}
          </Button>
        </RouterLink>

        <button
          type="button"
          class="header-ant-mobile-menu"
          :aria-expanded="inlineNavOpen"
          aria-controls="header-ant-inline-nav"
          :aria-label="t('landing.footer.navigation')"
          @click="toggleInlineNav"
        >
          <span />
          <span />
          <span />
        </button>
      </Space>

      <div
        id="header-ant-inline-nav"
        class="header-ant-inline-nav header-ant-menu-mobile-wrap"
      >
        <Menu
          theme="dark"
          mode="inline"
          class="header-ant-menu-mobile"
          :selected-keys="selectedNavKeys"
        >
          <Menu.Item v-for="(item, i) in props.items" :key="`nav-${i}`">
            <RouterLink :to="item.to" @click="closeInlineNav">{{ item.label }}</RouterLink>
          </Menu.Item>
        </Menu>
      </div>
    </div>
  </Layout.Header>

  <header v-else class="site-header">
    <div class="site-header-inner">
      <RouterLink to="/" class="brand-link">
        <span class="brand-mark">
          <img src="/logo.png" :alt="t('landing.brand.name')" />
        </span>
        <span class="brand-copy">
          <strong>{{ t('landing.brand.name') }}</strong>
          <span>{{ t('landing.brand.meta') }}</span>
        </span>
      </RouterLink>

      <nav class="site-nav" :aria-label="t('landing.footer.navigation')">
        <RouterLink
          v-for="item in props.items"
          :key="item.label"
          :to="item.to"
          class="site-nav-link"
          :class="{ 'is-active': item.active }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="site-header-actions">
        <LanguageSwitch :model-value="props.modelValue" tone="light" @update:model-value="updateLocale" />
        <RouterLink :to="props.actionTo" class="site-header-cta">
          {{ props.actionLabel }}
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header-ant {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: auto;
  min-height: var(--landing-header-height, 64px);
  padding: 0;
  line-height: 1.5;
  background: var(--landing-nav-bg, #001529) !important;
  box-shadow: var(--landing-shadow, 0 2px 8px rgba(0, 0, 0, 0.15));
  /* Added transitions for a smooth morphing effect */
  transition:
    top 0.4s cubic-bezier(0.4, 0, 0.2, 1),
    padding 0.4s cubic-bezier(0.4, 0, 0.2, 1),
    background 0.3s ease,
    box-shadow 0.3s ease;
}

.header-ant--floating {
  top: 12px;
  background: transparent !important;
  box-shadow: none;
  padding: 0 clamp(12px, 3vw, 24px);
}

.header-ant--floating .header-ant-bar.home-page {
  background: var(--landing-nav-bg, #001529);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.header-ant-bar.home-page {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  position: relative;
  max-width: 1200px;
  width: 100%;
  min-height: var(--landing-header-height, 64px);
  margin: 0 auto;
  padding: 0 24px;
  gap: 8px;
  row-gap: 0;
  background: transparent;
  border-radius: 0;
  /* Control the border-radius and shadow during the transition */
  transition:
    border-radius 0.4s cubic-bezier(0.4, 0, 0.2, 1),
    box-shadow 0.4s cubic-bezier(0.4, 0, 0.2, 1),
    max-width 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.header-ant-brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  margin-right: 16px;
  color: #fff;
  text-decoration: none;
  z-index: 101;
}

.header-ant-logo {
  display: flex;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.header-ant-logo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-ant-title {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.header-ant-menu-desktop {
  flex: 1 1 auto;
  min-width: min(280px, 100%);
  border-bottom: none !important;
  background: transparent !important;
  line-height: 62px;
}

/* 桌面导航：选中态仅用底边指示，切换时过渡 */
.header-ant-menu-desktop :deep(.ant-menu-item) {
  padding-inline: 14px;
  margin-inline: 2px;
  background: transparent !important;
  color: rgba(255, 255, 255, 0.65) !important;
  border-radius: 0 !important;
  border-bottom: 2px solid transparent !important;
  transition:
    color 0.28s ease,
    border-color 0.28s ease,
    background-color 0.28s ease;
}

.header-ant-menu-desktop :deep(.ant-menu-item:hover) {
  color: rgba(255, 255, 255, 0.95) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}

.header-ant-menu-desktop :deep(.ant-menu-item-selected) {
  color: #fff !important;
  background: transparent !important;
  border-bottom-color: #1890ff !important;
}

.header-ant-menu-desktop :deep(.ant-menu-item::after) {
  display: none !important;
}

.header-ant-menu-desktop :deep(.ant-menu-item a) {
  color: inherit;
  transition: color 0.28s ease;
}

.header-ant-actions {
  flex-shrink: 0;
  margin-left: auto;
}

.header-ant-lang-btn {
  color: rgba(255, 255, 255, 0.85) !important;
}

/* Nav0 汉堡：三线，仅移动端显示 */
.header-ant-mobile-menu {
  display: none;
  width: 22px;
  height: 16px;
  padding: 0;
  margin: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  position: relative;
  z-index: 100;
  flex-shrink: 0;
}

.header-ant-mobile-menu span {
  display: block;
  width: 100%;
  height: 2px;
  background: #fff;
  margin-top: 5px;
  border-radius: 1px;
  transition:
    transform 0.3s ease-in-out,
    opacity 0.3s ease-in-out;
}

.header-ant-mobile-menu span:first-child {
  margin-top: 0;
}

.home-page.open .header-ant-mobile-menu span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}

.home-page.open .header-ant-mobile-menu span:nth-child(2) {
  opacity: 0;
}

.home-page.open .header-ant-mobile-menu span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

/* 移动端 inline 菜单容器：全宽换行，默认收起（对齐 nav0 > .header0-menu opacity） */
.header-ant-inline-nav {
  flex: 1 1 100%;
  max-height: 0;
  opacity: 0;
  overflow: hidden;
  pointer-events: none;
  margin-left: -24px;
  width: calc(100% + 48px);
  transition:
    opacity 0.3s ease-in-out,
    max-height 0.35s ease-in-out;
}

.header-ant-menu-mobile {
  border-inline: none !important;
  background: transparent !important;
}

.header-ant-menu-mobile :deep(.ant-menu-item) {
  height: auto !important;
  line-height: 1.5 !important;
  padding-inline: 24px !important;
  margin-inline: 0 !important;
  width: 100%;
}

.header-ant-menu-mobile :deep(.ant-menu-item a) {
  display: block;
  padding: 8px 0;
  color: rgba(255, 255, 255, 0.85);
}

.header-ant-menu-mobile :deep(.ant-menu-item-selected) {
  background: rgba(24, 144, 255, 0.12) !important;
}

.home-page.open .header-ant-inline-nav {
  max-height: 480px;
  opacity: 1;
  pointer-events: auto;
  padding-bottom: 8px;
}

@media (min-width: 992px) {
  .header-ant-inline-nav {
    display: none;
  }

  .header-ant-mobile-menu {
    display: none !important;
  }
}

@media (max-width: 991px) {
  .header-ant-menu-desktop {
    display: none !important;
  }

  .header-ant-mobile-menu {
    display: block;
  }

  .header-ant-bar.home-page {
    padding-top: 10px;
    padding-bottom: 10px;
  }
}

.site-header {
  position: relative;
  z-index: 100;
}

.site-header-inner {
  width: min(1240px, 100%);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  min-height: 78px;
  padding: 0 20px;
  border: 1px solid var(--site-border, rgba(45, 35, 27, 0.12));
  border-radius: 28px;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.72), rgba(255, 247, 234, 0.9)),
    var(--site-surface, rgba(255, 249, 241, 0.78));
  box-shadow: var(--site-shadow, 0 24px 60px rgba(58, 39, 20, 0.08));
  backdrop-filter: blur(18px);
}

.brand-link {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: var(--site-text, #1f1a16);
  text-decoration: none;
  min-width: 0;
}

.brand-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 46px;
  border-radius: 16px;
  background: rgba(25, 18, 14, 0.06);
  border: 1px solid rgba(25, 18, 14, 0.08);
  overflow: hidden;
}

.brand-mark img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.brand-copy {
  display: grid;
  gap: 2px;
}

.brand-copy strong {
  font-size: 1.12rem;
  letter-spacing: -0.04em;
}

.brand-copy span {
  color: var(--site-muted, #6f665d);
  font-size: 0.72rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.site-nav,
.site-header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.site-nav {
  flex-wrap: wrap;
  justify-content: center;
}

.site-nav-link {
  padding: 10px 14px;
  border-radius: 999px;
  color: var(--site-muted, #6f665d);
  text-decoration: none;
  transition:
    color 0.22s ease,
    background 0.22s ease,
    transform 0.22s ease;
}

.site-nav-link:hover,
.site-nav-link.is-active {
  color: var(--site-text, #1f1a16);
  background: rgba(199, 91, 44, 0.12);
  transform: translateY(-1px);
}

.site-header-cta {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 44px;
  padding: 0 18px;
  border-radius: 999px;
  background: var(--site-accent, #bf5a2a);
  color: #fff7ef;
  font-weight: 700;
  text-decoration: none;
  transition:
    transform 0.22s ease,
    box-shadow 0.22s ease,
    opacity 0.22s ease;
  box-shadow: 0 14px 28px rgba(160, 66, 24, 0.22);
}

.site-header-cta:hover {
  transform: translateY(-1px);
  box-shadow: 0 18px 32px rgba(160, 66, 24, 0.26);
}

@media (max-width: 980px) {
  .site-header-inner {
    flex-wrap: wrap;
    justify-content: center;
    padding: 16px;
  }

  .brand-link {
    width: 100%;
    justify-content: center;
  }

  .site-header-actions {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 680px) {
  .site-header {
    padding-inline: 16px;
  }

  .site-header-inner {
    border-radius: 22px;
  }

  .site-nav {
    width: 100%;
  }

  .site-nav-link {
    padding-inline: 12px;
  }
}
</style>
