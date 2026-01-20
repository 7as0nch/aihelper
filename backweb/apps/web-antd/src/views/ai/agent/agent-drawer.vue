<script setup lang="ts">
import type {
  CreateAgentRequest,
  UpdateAgentRequest,
} from '#/api/ai/agent/model';

import { computed, ref } from 'vue';

import { Collapse, CollapsePanel } from 'ant-design-vue';
import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { cloneDeep } from '@vben/utils';

import { useVbenForm } from '#/adapter/form';
import { agentAdd, agentInfo, agentUpdate } from '#/api/ai/agent';
import PromptEditorWrapper from '#/components/PromptEditorWrapper/index.vue';
import ChatTestPanel from './components/ChatTestPanel.vue';

import {
  modelConfigSchema,
  basicInfoSchema,
  agentConfigSchema,
  featuresSchema,
} from './data';

const emit = defineEmits<{ reload: [] }>();

const isUpdate = ref(false);
const title = computed(() => {
  return isUpdate.value ? $t('pages.common.edit') : $t('pages.common.add');
});

// 表单数据
const formData = ref<Partial<CreateAgentRequest & UpdateAgentRequest>>({});
const systemPrompt = ref('');
const userInputPrompt = ref('');
const activeKeys = ref(['model', 'basic', 'agent', 'features']);

// 分组表单
const [ModelForm, modelFormApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-1',
  },
  layout: 'vertical',
  schema: modelConfigSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-1',
});

const [BasicForm, basicFormApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-1',
  },
  layout: 'vertical',
  schema: basicInfoSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-1',
});

const [AgentForm, agentFormApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-1',
  },
  layout: 'vertical',
  schema: agentConfigSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-1',
});

const [FeaturesForm, featuresFormApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-1',
  },
  layout: 'vertical',
  schema: featuresSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-1',
});

const [BasicDrawer, drawerApi] = useVbenDrawer({
  onCancel: handleCancel,
  onConfirm: handleConfirm,
  async onOpenChange(isOpen) {
    if (!isOpen) {
      return null;
    }
    drawerApi.drawerLoading(true);
    const { id } = drawerApi.getData() as { id?: number | string };
    isUpdate.value = !!id;

    if (isUpdate.value && id) {
      const record = await agentInfo(id);
      formData.value = cloneDeep(record);
      systemPrompt.value = record.systemPrompt || '';
      userInputPrompt.value = record.userInputPrompt || '';

      await modelFormApi.setValues({ id: record.id, originalModelId: record.originalModelId });
      await basicFormApi.setValues({
        name: record.name,
        code: record.code,
        description: record.description,
      });
      await agentFormApi.setValues({
        adapterType: record.adapterType,
        maxIteration: record.maxIteration,
        type: record.type,
        status: record.status,
        order: record.order,
      });
      await featuresFormApi.setValues({
        withWriteTodos: record.withWriteTodos,
        withWebSearchAgent: record.withWebSearchAgent,
      });
    } else {
      formData.value = {};
      systemPrompt.value = '';
      userInputPrompt.value = '';
      await modelFormApi.resetForm();
      await basicFormApi.resetForm();
      await agentFormApi.resetForm();
      await featuresFormApi.resetForm();
    }

    drawerApi.drawerLoading(false);
  },
});

async function handleConfirm() {
  try {
    drawerApi.drawerLoading(true);

    // 验证所有表单
    const [modelValid, basicValid, agentValid, featuresValid] = await Promise.all([
      modelFormApi.validate(),
      basicFormApi.validate(),
      agentFormApi.validate(),
      featuresFormApi.validate(),
    ]);

    if (!modelValid.valid || !basicValid.valid || !agentValid.valid || !featuresValid.valid) {
      return;
    }

    // 合并所有表单数据
    const [modelValues, basicValues, agentValues, featuresValues] = await Promise.all([
      modelFormApi.getValues(),
      basicFormApi.getValues(),
      agentFormApi.getValues(),
      featuresFormApi.getValues(),
    ]);

    const data = {
      ...modelValues,
      ...basicValues,
      ...agentValues,
      ...featuresValues,
      systemPrompt: systemPrompt.value,
      userInputPrompt: userInputPrompt.value,
    } as CreateAgentRequest & UpdateAgentRequest;

    if (isUpdate.value && formData.value.id) {
      data.id = formData.value.id;
    }

    await (isUpdate.value ? agentUpdate(data) : agentAdd(data));
    emit('reload');
    await handleCancel();
  } catch (error) {
    console.error(error);
  } finally {
    drawerApi.drawerLoading(false);
  }
}

async function handleCancel() {
  drawerApi.close();
  await Promise.all([
    modelFormApi.resetForm(),
    basicFormApi.resetForm(),
    agentFormApi.resetForm(),
    featuresFormApi.resetForm(),
  ]);
  formData.value = {};
  systemPrompt.value = '';
  userInputPrompt.value = '';
}
</script>

<template>
  <BasicDrawer :close-on-click-modal="true" :title="title" class="w-[1400px]">
    <div class="h-[calc(100vh-120px)] flex gap-4">
      <!-- 左侧：提示词编辑器 -->
      <div class="w-[400px] flex-shrink-0">
        <PromptEditorWrapper
          v-model:system-prompt="systemPrompt"
          v-model:user-input-prompt="userInputPrompt"
          class="h-full"
        />
      </div>

      <!-- 中间：Agent 配置 -->
      <div class="flex-1 min-w-0 overflow-y-auto pr-2">
        <Collapse v-model:activeKey="activeKeys" class="bg-white">
          <CollapsePanel key="model" header="AI 模型配置">
            <div class="py-2">
              <ModelForm />
            </div>
          </CollapsePanel>
          <CollapsePanel key="basic" header="基本信息">
            <div class="py-2">
              <BasicForm />
            </div>
          </CollapsePanel>
          <CollapsePanel key="agent" header="Agent 配置">
            <div class="py-2">
              <AgentForm />
            </div>
          </CollapsePanel>
          <CollapsePanel key="features" header="功能开关">
            <div class="py-2">
              <FeaturesForm />
            </div>
          </CollapsePanel>
        </Collapse>
      </div>

      <!-- 右侧：聊天测试 -->
      <div class="w-[400px] flex-shrink-0">
        <ChatTestPanel class="h-full" />
      </div>
    </div>
  </BasicDrawer>
</template>
