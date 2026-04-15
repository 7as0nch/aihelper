<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { Card, Form, Input, Segmented, Button, Space, Collapse, Typography, Result } from 'ant-design-vue';
import { CopyOutlined, MailOutlined } from '@ant-design/icons-vue';
import { useI18n } from 'vue-i18n';
import { betaApi, type BetaApplicationPayload } from '@/api/beta';
import type { SupportedLocale } from '@/i18n/types';
import { useLocaleStore } from '@/stores/locale';
import LandingFooter from '@/components/landing/LandingFooter.vue';
import LandingBackToTop from '@/components/landing/LandingBackToTop.vue';
import PublicSiteHeader from '@/components/landing/PublicSiteHeader.vue';

interface ApplyCopy {
  badge: string;
  title: string;
  description: string;
  productLabel: string;
  contactTypeLabel: string;
  contactTypeEmail: string;
  contactTypeQQ: string;
  contactValueLabel: string;
  contactValuePlaceholderEmail: string;
  contactValuePlaceholderQQ: string;
  useCaseLabel: string;
  useCasePlaceholder: string;
  noteSectionTitle: string;
  noteLabel: string;
  notePlaceholder: string;
  submit: string;
  submitting: string;
  successTitle: string;
  successDescription: string;
  successAgain: string;
  backHome: string;
  altPrefix: string;
  copyQq: string;
  sendEmail: string;
  validation: Record<string, string>;
  products: {
    litechat: string;
    aicook: string;
    techSandbox: string;
  };
}

const currentYear = new Date().getFullYear();
const companySiteLink = 'https://aihelper.chat/';
const githubLink = 'https://github.com/7as0nch/aihelper.git';
const giteeLink = 'https://gitee.com/jas0nch/aichat.git';
const icp = '蜀ICP备2026011723号-1';
const icpLink = 'https://beian.miit.gov.cn';
const contactEmail = '7as0nch@gmail.com';
const qqNumber = '2538684421';

const localeStore = useLocaleStore();
const { t, tm } = useI18n();
const copied = ref(false);
const isSubmitting = ref(false);
const isSubmitted = ref(false);
const errors = reactive<Record<string, string>>({});

const form = reactive<BetaApplicationPayload>({
  productInterest: 'litechat',
  contactType: 'email',
  contactValue: '',
  useCase: '',
  note: '',
  sourcePage: '/apply',
});

const applyCopy = computed(() => tm('landing.apply') as ApplyCopy);
const navItems = computed(() => [
  { label: t('landing.nav.products'), to: { path: '/', hash: '#products' } },
  { label: t('landing.nav.apply'), to: '/apply', active: true },
]);
const emailLink = computed(
  () => `mailto:${contactEmail}?subject=${encodeURIComponent('Aichat Beta Access')}`,
);

const productSegmentOptions = computed(() => [
  { label: applyCopy.value.products.litechat, value: 'litechat' },
  { label: applyCopy.value.products.aicook, value: 'aicook' },
  { label: applyCopy.value.products.techSandbox, value: 'tech-sandbox' },
]);

const contactSegmentOptions = computed(() => [
  { label: applyCopy.value.contactTypeEmail, value: 'email' },
  { label: applyCopy.value.contactTypeQQ, value: 'qq' },
]);

const contactPlaceholder = computed(() =>
  form.contactType === 'email'
    ? applyCopy.value.contactValuePlaceholderEmail
    : applyCopy.value.contactValuePlaceholderQQ,
);

let copyTimer = 0;

const updateLocale = (nextLocale: SupportedLocale) => {
  localeStore.applyLocale(nextLocale);
};

const handleCopyQQ = async () => {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(qqNumber);
    } else {
      const textarea = document.createElement('textarea');
      textarea.value = qqNumber;
      textarea.setAttribute('readonly', 'true');
      textarea.style.position = 'absolute';
      textarea.style.left = '-9999px';
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }
    copied.value = true;
    window.clearTimeout(copyTimer);
    copyTimer = window.setTimeout(() => {
      copied.value = false;
    }, 1800);
  } catch (_error) {
    copied.value = false;
  }
};

const resetErrors = () => {
  Object.keys(errors).forEach((key) => {
    delete errors[key];
  });
};

const validateForm = () => {
  resetErrors();
  if (!form.productInterest) {
    errors.productInterest = applyCopy.value.validation.product;
  }
  if (!form.contactType) {
    errors.contactType = applyCopy.value.validation.contactType;
  }
  if (!form.contactValue.trim()) {
    errors.contactValue = applyCopy.value.validation.contactValue;
  }
  if (!form.useCase.trim()) {
    errors.useCase = applyCopy.value.validation.useCase;
  }
  return Object.keys(errors).length === 0;
};

const handleSubmit = async (event?: Event) => {
  event?.preventDefault();
  event?.stopPropagation();
  if (isSubmitting.value) {
    return;
  }
  if (!validateForm()) {
    return;
  }
  isSubmitting.value = true;
  try {
    await betaApi.submitApplication({
      ...form,
      contactValue: form.contactValue.trim(),
      useCase: form.useCase.trim(),
      note: form.note?.trim(),
    });
    isSubmitted.value = true;
    form.contactValue = '';
    form.useCase = '';
    form.note = '';
    form.contactType = 'email';
    form.productInterest = 'litechat';
    resetErrors();
  } finally {
    isSubmitting.value = false;
  }
};

const resetSuccess = () => {
  isSubmitted.value = false;
};
</script>

<template>
  <div class="landing-ant apply-page">
    <PublicSiteHeader
      variant="ant"
      :items="navItems"
      :model-value="localeStore.locale"
      :action-label="applyCopy.backHome"
      action-to="/"
      @update:model-value="updateLocale"
    />

    <main class="apply-main">
      <div class="apply-wrap">
        <div class="apply-intro">
          <Typography.Text type="secondary" class="apply-badge">{{ applyCopy.badge }}</Typography.Text>
          <Typography.Title :level="2" class="apply-title">{{ applyCopy.title }}</Typography.Title>
          <Typography.Paragraph type="secondary" class="apply-desc">{{ applyCopy.description }}</Typography.Paragraph>
        </div>

        <Card class="apply-card">
          <Result
            v-if="isSubmitted"
            status="success"
            :title="applyCopy.successTitle"
            :sub-title="applyCopy.successDescription"
          >
            <template #extra>
              <Button type="primary" @click="resetSuccess">{{ applyCopy.successAgain }}</Button>
            </template>
          </Result>

          <Form v-else layout="vertical" :model="form" class="apply-form" :colon="false" @submit="handleSubmit">
            <Form.Item :label="applyCopy.productLabel" v-bind="errors.productInterest ? { validateStatus: 'error', help: errors.productInterest } : {}">
              <Segmented
                v-model:value="form.productInterest"
                block
                :options="productSegmentOptions"
                class="apply-segmented"
              />
            </Form.Item>

            <Form.Item :label="applyCopy.contactTypeLabel" v-bind="errors.contactType ? { validateStatus: 'error', help: errors.contactType } : {}">
              <Segmented v-model:value="form.contactType" block :options="contactSegmentOptions" />
            </Form.Item>

            <Form.Item :label="applyCopy.contactValueLabel" v-bind="errors.contactValue ? { validateStatus: 'error', help: errors.contactValue } : {}">
              <Input v-model:value="form.contactValue" size="large" :placeholder="contactPlaceholder" allow-clear />
            </Form.Item>

            <Form.Item :label="applyCopy.useCaseLabel" v-bind="errors.useCase ? { validateStatus: 'error', help: errors.useCase } : {}">
              <Input.TextArea
                v-model:value="form.useCase"
                :rows="4"
                :placeholder="applyCopy.useCasePlaceholder"
                :maxlength="2000"
                show-count
              />
            </Form.Item>

            <Collapse ghost class="apply-note-collapse">
              <Collapse.Panel :header="applyCopy.noteSectionTitle" key="note">
                <Form.Item :label="applyCopy.noteLabel" class="apply-note-field">
                  <Input.TextArea v-model:value="form.note" :rows="3" :placeholder="applyCopy.notePlaceholder" />
                </Form.Item>
              </Collapse.Panel>
            </Collapse>

            <Form.Item>
              <Button type="primary" block size="large" html-type="submit" :loading="isSubmitting">
                {{ isSubmitting ? applyCopy.submitting : applyCopy.submit }}
              </Button>
            </Form.Item>
          </Form>

          <div class="apply-alt">
            <Typography.Text type="secondary">{{ applyCopy.altPrefix }}</Typography.Text>
            <Space :size="12" wrap>
              <Button type="link" size="small" @click="handleCopyQQ">
                <template #icon><CopyOutlined /></template>
                {{ copied ? applyCopy.copyQq + ' ✓' : applyCopy.copyQq }}
              </Button>
              <Typography.Text type="secondary">·</Typography.Text>
              <Button type="link" size="small" tag="a" :href="emailLink">
                <template #icon><MailOutlined /></template>
                {{ applyCopy.sendEmail }}
              </Button>
            </Space>
          </div>
        </Card>
      </div>
    </main>

    <LandingBackToTop />

    <LandingFooter
      variant="ant"
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
.apply-page {
  min-height: 100vh;
  padding-top: var(--landing-header-height, 64px);
  color: var(--landing-text, rgba(0, 0, 0, 0.85));
  background: var(--landing-page-bg, #ececec);
}

.apply-main {
  padding: 24px 16px 40px;
}

.apply-wrap {
  max-width: 560px;
  margin: 0 auto;
}

.apply-intro {
  margin-bottom: 20px;
  text-align: center;
}

.apply-badge {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.apply-title {
  margin: 8px 0 12px !important;
}

.apply-desc {
  margin-bottom: 0 !important;
}

.apply-card :deep(.ant-card-body) {
  padding: 28px 24px 20px;
}

.apply-form :deep(.ant-form-item-label > label) {
  font-weight: 600;
}

.apply-segmented :deep(.ant-segmented-item) {
  flex: 1;
  justify-content: center;
}

.apply-note-collapse :deep(.ant-collapse-header) {
  padding: 4px 0 !important;
  font-weight: 500;
}

.apply-note-field {
  margin-bottom: 0 !important;
}

.apply-alt {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--landing-line, #f0f0f0);
  text-align: center;
}

.apply-alt .ant-typography {
  display: block;
  margin-bottom: 8px;
}
</style>
