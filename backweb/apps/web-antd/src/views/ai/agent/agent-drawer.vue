<script setup lang="ts">
import type {
  CreateAgentRequest,
  UpdateAgentRequest,
} from '#/api/ai/agent/model';

import { computed, ref, onMounted } from 'vue';

import { Collapse, CollapsePanel, Select, message } from 'ant-design-vue';
import {
  CodeOutlined,
  MessageOutlined,
  SettingOutlined,
  ThunderboltOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
} from '@ant-design/icons-vue';
import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { cloneDeep } from '@vben/utils';

import { useVbenForm } from '#/adapter/form';
import type { ModelInfo } from '#/api/ai/model/model';
import { agentAdd, agentInfo, agentUpdate } from '#/api/ai/agent';
import { modelInfo as getModelInfo } from '#/api/ai/model';
import { promptList } from '#/api/ai/prompt';
import type { PromptInfo } from '#/api/ai/prompt/model';
import PromptEditorWrapper from '#/components/PromptEditorWrapper/index.vue';
import ChatTestPanel from './components/ChatTestPanel.vue';

import {
  modelConfigSchema,
  basicInfoSchema,
  agentConfigSchema,
  featuresSchema,
} from './data';

const panelClass =
  'rounded-xl border border-slate-200/60 bg-white shadow-lg shadow-slate-200/50 flex flex-col min-h-0 overflow-hidden';
const panelHeaderClass =
  'px-5 py-4 border-b border-slate-100 bg-gradient-to-r from-slate-50 to-transparent flex items-center justify-between';
const panelBodyClass = 'flex-1 min-h-0 p-4';
const dividerClass =
  'group flex items-stretch justify-center cursor-col-resize select-none';
const rightPanelTransitionClass =
  'transition-[opacity,transform] duration-300 ease-in-out';

const emit = defineEmits<{ reload: [] }>();

const isUpdate = ref(false);
const title = computed(() => {
  return isUpdate.value ? $t('pages.common.edit') : $t('pages.common.add');
});

// 提示词模板相关
const activePromptTab = ref('system');
const templates = ref<PromptInfo[]>([]);
const templateLoading = ref(false);

const fetchTemplates = async () => {
  try {
    templateLoading.value = true;
    const res = await promptList({ page: 1, pageSize: 1000 });
    templates.value = res.list || [];
  } catch (error) {
    console.error('获取模板列表失败:', error);
  } finally {
    templateLoading.value = false;
  }
};

const currentTemplates = computed(() => {
  const type = activePromptTab.value === 'system' ? 1 : 2;
  return templates.value.filter((item) => item.type === type);
});

const handleTemplateSelect = (value: any) => {
  const id = value as number;
  const template = templates.value.find((item) => item.id === id);
  if (template) {
    if (activePromptTab.value === 'system') {
      systemPrompt.value = template.text;
    } else {
      userInputPrompt.value = template.text;
    }
    message.success(`已应用模板: ${template.name}`);
  }
};

onMounted(() => {
  fetchTemplates();
});

// 表单数据
const formData = ref<Partial<CreateAgentRequest & UpdateAgentRequest>>({});
const systemPrompt = ref('');
const userInputPrompt = ref('');
const layoutRef = ref<HTMLElement | null>(null);
const leftWidth = ref(0);
const rightWidth = ref(0);
const dividerSize = 8;
const isRightCollapsed = ref(false);
const selectedModel = ref<ModelInfo | null>(null);
const activeKeys = ref<Array<string | number>>([
  'model',
  'basic',
  'agent',
  'features',
]);

const gridStyle = computed(() => {
  if (leftWidth.value === 0) {
    return {
      gridTemplateColumns: isRightCollapsed.value
        ? `1fr ${dividerSize}px 1fr 0px 0px`
        : `3fr ${dividerSize}px 3fr ${dividerSize}px 2fr`,
    };
  }
  return {
    gridTemplateColumns: isRightCollapsed.value
      ? `minmax(0,1fr) ${dividerSize}px minmax(0,1fr) 0px 0px`
      : `${leftWidth.value}px ${dividerSize}px minmax(0,1fr) ${dividerSize}px ${rightWidth.value}px`,
  };
});

// 拖拽逻辑保持现状
function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max);
}

function startDrag(target: 'left' | 'right', event: MouseEvent) {
  if (!layoutRef.value) {
    return;
  }

  event.preventDefault();

  const startX = event.clientX;
  const startLeft = leftWidth.value;
  const startRight = rightWidth.value;
  const containerWidth = layoutRef.value.clientWidth;
  const minLeft = 320;
  const minCenter = 360;
  const minRight = 340;

  const handleMove = (moveEvent: MouseEvent) => {
    const deltaX = moveEvent.clientX - startX;

    // 如果是第一次拖拽，先从 3:3:2 转换为像素值
    if (leftWidth.value === 0) {
      const currentContainerWidth = layoutRef.value!.clientWidth;
      const total = currentContainerWidth - dividerSize * 2;
      const unit = total / 8;
      leftWidth.value = Math.round(unit * 3);
      rightWidth.value = Math.round(unit * 2);
    }

    if (target === 'left') {
      const maxLeft = containerWidth -
        (minCenter + (isRightCollapsed.value ? dividerSize : minRight + dividerSize * 2));
      leftWidth.value = clamp(startLeft + deltaX, minLeft, Math.max(minLeft, maxLeft));
      return;
    }

    const maxRight = containerWidth - (leftWidth.value + dividerSize * 2 + minCenter);
    rightWidth.value = clamp(startRight - deltaX, minRight, Math.max(minRight, maxRight));
  };

  const handleUp = () => {
    document.removeEventListener('mousemove', handleMove);
    document.removeEventListener('mouseup', handleUp);
    document.body.style.cursor = '';
  };

  document.body.style.cursor = 'col-resize';
  document.addEventListener('mousemove', handleMove);
  document.addEventListener('mouseup', handleUp);
}

function togglePreview() {
  const nextCollapsed = !isRightCollapsed.value;
  isRightCollapsed.value = nextCollapsed;
  if (!nextCollapsed) {
    leftWidth.value = 0; // 重置为 0 以便 gridStyle 使用 fr 比例
  }
}

onMounted(() => {
  // 初始状态使用 3:3:2 的 fr 比例，不需要手动 resetWidthsToRatio
});

async function updateSelectedModel(modelId?: number | string, fillForm = false) {
  if (modelId) {
    try {
      const model = await getModelInfo(modelId);
      selectedModel.value = model;
      if (fillForm) {
        modelFormApi.setValues({
          aiModelBaseUrl: model.baseUrl,
          aiModelApiKey: model.apiKey,
          aiModelTemperature: model.temperature,
          aiModelTopP: model.topP,
        });
      }
    } catch (error) {
      console.error('获取模型详情失败:', error);
      selectedModel.value = null;
    }
  } else {
    selectedModel.value = null;
    if (fillForm) {
      modelFormApi.setValues({
        aiModelBaseUrl: undefined,
        aiModelApiKey: undefined,
        aiModelTemperature: undefined,
        aiModelTopP: undefined,
      });
    }
  }
}

// 分组表单
const [ModelForm, modelFormApi] = useVbenForm({
  commonConfig: {
    componentProps: { class: 'w-full' },
    formItemClass: 'col-span-1',
  },
  layout: 'vertical',
  schema: modelConfigSchema((value) => updateSelectedModel(value, true)),
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
  layout: 'horizontal',
  schema: featuresSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-1 gap-3',
});

const [BasicDrawer, drawerApi] = useVbenDrawer({
  onCancel: handleCancel,
  onConfirm: handleConfirm,
  async onOpenChange(isOpen) {
    if (!isOpen) {
      return null;
    }
    try {
      const { id } = drawerApi.getData() as { id?: number | string };
      isUpdate.value = !!id;

      if (isUpdate.value && id) {
        const record = await agentInfo(id);
        formData.value = cloneDeep(record);
        systemPrompt.value = record.systemPrompt || '';
        userInputPrompt.value = record.userInputPrompt || '';

        modelFormApi.setValues({
          id: record.id,
          originalModelId: record.originalModelId,
          aiModelBaseUrl: record.aiModel?.baseUrl,
          aiModelApiKey: record.aiModel?.apiKey,
          aiModelTemperature: record.aiModel?.temperature,
          aiModelTopP: record.aiModel?.topP,
        });
        updateSelectedModel(record.originalModelId);
        basicFormApi.setValues({
          name: record.name,
          code: record.code,
          description: record.description,
        });
        agentFormApi.setValues({
          adapterType: record.adapterType,
          maxIteration: record.maxIteration,
          type: record.type,
          status: record.status,
          order: record.order,
        });
        featuresFormApi.setValues({
          withWriteTodos: record.withWriteTodos,
          withWebSearchAgent: record.withWebSearchAgent,
        });
      } else {
        formData.value = {};
        systemPrompt.value = '';
        userInputPrompt.value = '';
        modelFormApi.resetForm();
        basicFormApi.resetForm();
        agentFormApi.resetForm();
        featuresFormApi.resetForm();
      }
    } catch (error) {
      console.error(error);
    }
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
      ...basicValues,
      ...agentValues,
      ...featuresValues,
      ...modelValues,
      aiModel: {
        baseUrl: modelValues.aiModelBaseUrl,
        apiKey: modelValues.aiModelApiKey,
        temperature: modelValues.aiModelTemperature,
        topP: modelValues.aiModelTopP,
      },
      systemPrompt: systemPrompt.value,
      userInputPrompt: userInputPrompt.value,
    } as any;

    // 清理扁平化的 key
    delete data.aiModelBaseUrl;
    delete data.aiModelApiKey;
    delete data.aiModelTemperature;
    delete data.aiModelTopP;

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
  <BasicDrawer :close-on-click-modal="true" :title="title" class="w-[95vw]">
    <div class="h-[calc(100vh-120px)] flex flex-col gap-4 bg-gradient-to-br from-slate-50 via-blue-50/30 to-slate-50 p-4 rounded-xl">
      <div class="bg-white/80 backdrop-blur-sm border border-slate-200/60 shadow-sm rounded-xl px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-blue-500/25">
              <ThunderboltOutlined class="text-white text-lg" />
            </div>
            <div>
              <div class="text-lg font-semibold text-slate-800">Agent 配置工作台</div>
              <div class="text-xs text-slate-500 mt-0.5">
                同屏完成提示词、配置与预览调试
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="px-3 py-1.5 bg-blue-50 text-blue-700 text-xs font-medium rounded-full border border-blue-100">
              开发模式
            </span>
            <button
              type="button"
              class="px-4 py-2 text-sm font-medium text-slate-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200 flex items-center gap-2"
              @click="togglePreview"
            >
              <FullscreenExitOutlined v-if="!isRightCollapsed" class="text-sm" />
              <FullscreenOutlined v-else class="text-sm" />
              {{ isRightCollapsed ? '展开预览' : '收起预览' }}
            </button>
          </div>
        </div>
      </div>

      <div
        ref="layoutRef"
        class="flex-1 min-h-0 grid gap-4"
        :style="gridStyle"
      >
        <!-- 左侧：提示词编辑器 -->
        <section :class="panelClass">
          <div :class="panelHeaderClass">
            <div class="flex items-center gap-2">
              <CodeOutlined class="text-blue-600" />
              <div class="text-sm font-semibold text-slate-800">提示词编辑器</div>
            </div>
            <div class="flex items-center gap-4">
              <div class="text-xs text-slate-400 font-medium">系统/用户提示词</div>
              <div class="flex items-center gap-3">
                <span class="text-xs font-medium text-slate-500">选择模板:</span>
                <Select
                  placeholder="选择模板"
                  size="small"
                  class="w-[160px]"
                  :loading="templateLoading"
                  @select="handleTemplateSelect"
                >
                  <Select.Option v-for="item in currentTemplates" :key="String(item.id)" :value="item.id">
                    {{ item.name }}
                  </Select.Option>
                </Select>
              </div>
            </div>
          </div>
          <div :class="[panelBodyClass, 'overflow-hidden']">
            <PromptEditorWrapper
              :system-prompt="systemPrompt"
              :user-input-prompt="userInputPrompt"
              :active-tab="activePromptTab"
              @update:systemPrompt="(value) => (systemPrompt = value)"
              @update:userInputPrompt="(value) => (userInputPrompt = value)"
              @update:activeTab="(value) => (activePromptTab = value)"
              class="h-full"
            />
          </div>
        </section>

        <div
          :class="dividerClass"
          :style="{ width: `${dividerSize}px` }"
          @mousedown="startDrag('left', $event)"
        >
          <div class="w-0.5 h-full bg-slate-200 group-hover:bg-blue-400 transition-all duration-200 rounded-full" />
        </div>

        <!-- 中间：Agent 配置 -->
        <section :class="panelClass">
          <div :class="panelHeaderClass">
            <div class="flex items-center gap-2">
              <SettingOutlined class="text-blue-600" />
              <div class="text-sm font-semibold text-slate-800">Agent 配置</div>
            </div>
            <div class="text-xs text-slate-400 font-medium">模型与运行参数</div>
          </div>
          <div :class="[panelBodyClass, 'overflow-y-auto pr-2']">
            <Collapse
              :activeKey="activeKeys"
              class="bg-transparent space-y-3"
              :bordered="false"
              @update:activeKey="(val) => (activeKeys = Array.isArray(val) ? val : [val])"
            >
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
        </section>

        <div
          :class="[dividerClass, isRightCollapsed ? 'opacity-0 pointer-events-none' : 'opacity-100']"
          :style="{ width: `${dividerSize}px` }"
          @mousedown="startDrag('right', $event)"
        >
          <div class="w-0.5 h-full bg-slate-200 group-hover:bg-blue-400 transition-all duration-200 rounded-full" />
        </div>

        <!-- 右侧：聊天测试 -->
        <section
          :class="[
            panelClass,
            rightPanelTransitionClass,
            isRightCollapsed
              ? 'opacity-0 pointer-events-none translate-x-2 w-0 overflow-hidden'
              : 'opacity-100',
          ]"
        >
          <div :class="panelHeaderClass">
            <div class="flex items-center gap-2">
              <MessageOutlined class="text-blue-600" />
              <div class="text-sm font-semibold text-slate-800">预览与调试</div>
            </div>
            <div class="flex items-center gap-3">
              <div class="text-xs text-slate-400 font-medium">模拟对话</div>
              <button
                type="button"
                class="text-xs text-slate-500 hover:text-slate-700 font-medium transition-colors"
                @click="togglePreview"
              >
                收起
              </button>
            </div>
          </div>
          <div :class="[panelBodyClass, 'overflow-hidden p-0']">
            <ChatTestPanel class="h-full" variant="plain" :show-header="false" />
          </div>
        </section>
      </div>
    </div>
  </BasicDrawer>
</template>

<style scoped>
:deep(.ant-collapse-item) {
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 0.75rem;
  overflow: hidden;
  background: #ffffff;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.06);
}

:deep(.ant-collapse-header) {
  padding: 0.75rem 1rem !important;
  font-weight: 500;
  color: #334155;
  background: linear-gradient(90deg, rgba(248, 250, 252, 1), rgba(248, 250, 252, 0));
}

:deep(.ant-collapse-content) {
  border-top: 1px solid rgba(226, 232, 240, 0.7);
  background: #ffffff;
}

:deep(.ant-form-item-label > label) {
  font-weight: 600;
  color: #334155;
}

:deep(.ant-input),
:deep(.ant-input-affix-wrapper),
:deep(.ant-input-number),
:deep(.ant-input-number-input),
:deep(.ant-select-selector),
:deep(.ant-picker) {
  border-radius: 0.75rem !important;
  border-color: rgba(226, 232, 240, 0.9) !important;
}

:deep(.ant-input-affix-wrapper),
:deep(.ant-input-number),
:deep(.ant-select-selector),
:deep(.ant-picker) {
  min-height: 40px;
}

:deep(.ant-input-number-input) {
  height: 40px;
}
</style>
