<script setup lang="ts">
import type { Node } from '@vue-flow/core';

import { onMounted, ref, watch } from 'vue';

import { CloseOutlined, SettingOutlined } from '@ant-design/icons-vue';
import {
  Button,
  Checkbox,
  Collapse,
  CollapsePanel,
  Input,
  InputNumber,
  Select,
  Slider,
  Textarea,
} from 'ant-design-vue';

import { modelList } from '#/api/ai/model';

interface NodeConfigPanelProps {
  node: Node;
  showHeader?: boolean;
  showFooter?: boolean;
  autoSave?: boolean;
}

const props = withDefaults(defineProps<NodeConfigPanelProps>(), {
  showHeader: true,
  showFooter: true,
  autoSave: false,
});

const emit = defineEmits<{
  (e: 'update', data: any): void;
  (e: 'close'): void;
}>();

const formData = ref({
  name: '',
  code: '',
  description: '',
  adapterType: 1,
  originalModelId: undefined as number | undefined,
  modelName: '',
  modelType: '',
  baseUrl: '',
  apiKey: '',
  category: 0,
  temperature: 0.7,
  topP: 0.9,
  maxIteration: 10,
  systemPrompt: '',
  userInputPrompt: '',
  status: 1,
  type: 2,
  withWriteTodos: false,
  withWebSearchAgent: false,
});

const lastSavedSnapshot = ref('');
const activeKeys = ref<string[]>(['basic', 'model', 'prompt', 'features']);

// Model 选项
const modelOptions = ref<Array<{ label: string; value: number }>>([]);
const modelMap = ref(new Map<number, any>());

const normalizeModel = (model: any) => {
  return {
    id: Number(model.id ?? model.model_id),
    modelName: model.modelName ?? model.model_name ?? '',
    modelType: model.modelType ?? model.model_type ?? '',
    baseUrl: model.baseUrl ?? model.base_url ?? '',
    apiKey: model.apiKey ?? model.api_key ?? '',
    category: model.category ?? 0,
    temperature: model.temperature ?? undefined,
    topP: model.topP ?? model.top_p ?? undefined,
  };
};

const loadModelOptions = async () => {
  try {
    const result = await modelList({ pageNum: 1, pageSize: 100 });
    const list = (result.list || []).map(normalizeModel);
    modelOptions.value = list.map((model: any) => ({
      label: `${model.modelName} (${model.modelType})`,
      value: model.id,
    }));
    modelMap.value = new Map(list.map((model: any) => [model.id, model]));
  } catch (error) {
    console.error('加载 Model 列表失败:', error);
  }
};

const initFormData = () => {
  const data = props.node.data || {};
  formData.value = {
    name: data.name || data.label || '',
    code: data.code || '',
    description: data.description || '',
    adapterType: data.adapterType || 1,
    originalModelId: data.originalModelId,
    modelName: data.aiModel?.modelName || '',
    modelType: data.aiModel?.modelType || '',
    baseUrl: data.aiModel?.baseUrl || '',
    apiKey: data.aiModel?.apiKey || '',
    category: data.aiModel?.category ?? 0,
    temperature: data.aiModel?.temperature ?? data.temperature ?? 0.7,
    topP: data.aiModel?.topP ?? data.topP ?? 0.9,
    maxIteration: data.maxIteration || 10,
    systemPrompt: data.systemPrompt || '',
    userInputPrompt: data.userInputPrompt || '',
    status: data.status || 1,
    type: data.type || 2,
    withWriteTodos: data.withWriteTodos || false,
    withWebSearchAgent: data.withWebSearchAgent || false,
  };
  lastSavedSnapshot.value = JSON.stringify(formData.value);
};

watch(() => props.node, initFormData, { deep: true });
watch(
  formData,
  () => {
    if (!props.autoSave) return;
    const snapshot = JSON.stringify(formData.value);
    if (snapshot === lastSavedSnapshot.value) return;
    if (autoSaveTimer) {
      clearTimeout(autoSaveTimer);
    }
    autoSaveTimer = window.setTimeout(() => {
      handleSave();
    }, 500);
  },
  { deep: true },
);

onMounted(() => {
  initFormData();
  loadModelOptions();
});

let autoSaveTimer: number | null = null;

const handleSave = () => {
  // 按照 proto 结构整理数据
  const updateData = {
    ...formData.value,
    aiModel: {
      category: formData.value.category,
      modelType: formData.value.modelType,
      modelName: formData.value.modelName,
      apiKey: formData.value.apiKey,
      baseUrl: formData.value.baseUrl,
      temperature: formData.value.temperature,
      topP: formData.value.topP,
    },
  };
  emit('update', updateData);
  lastSavedSnapshot.value = JSON.stringify(formData.value);
};

// Type-safe update handlers
const handleAdapterTypeChange = (val: any) => {
  formData.value.adapterType = Number(val);
};

const handleTypeChange = (val: any) => {
  formData.value.type = Number(val);
};

const handleStatusChange = (val: any) => {
  formData.value.status = Number(val);
};

const handleOriginalModelIdChange = (val: any) => {
  const modelId = val ? Number(val) : undefined;
  formData.value.originalModelId = modelId;
  if (!modelId) return;
  const model = modelMap.value.get(modelId);
  if (!model) return;
  formData.value.modelName = model.modelName || '';
  formData.value.modelType = model.modelType || '';
  formData.value.baseUrl = model.baseUrl || '';
  formData.value.apiKey = model.apiKey || '';
  formData.value.category = model.category ?? 0;
  formData.value.temperature = model.temperature ?? formData.value.temperature;
  formData.value.topP = model.topP ?? formData.value.topP;
};

const handleMaxIterationChange = (val: any) => {
  formData.value.maxIteration = Number(val) || 10;
};

const handleTemperatureChange = (val: any) => {
  formData.value.temperature = Number(val);
};

const handleTopPChange = (val: any) => {
  formData.value.topP = Number(val);
};

const handleActiveKeysChange = (val: any) => {
  activeKeys.value = (Array.isArray(val) ? val : [val]).map(String);
};
</script>

<template>
  <div class="flex w-full flex-col bg-white">
    <!-- Header -->
    <div
      v-if="props.showHeader"
      class="flex items-center justify-between border-b border-slate-200 bg-gradient-to-r from-slate-50 to-white px-6 py-4"
    >
      <div class="flex items-center gap-2">
        <SettingOutlined class="text-lg text-slate-700" />
        <h2 class="m-0 font-semibold text-slate-800">节点配置</h2>
      </div>
      <button
        @click="emit('close')"
        class="cursor-pointer rounded-lg border-none bg-transparent p-1 transition-colors hover:bg-slate-100"
      >
        <CloseOutlined class="text-slate-600" />
      </button>
    </div>

    <!-- Content -->
    <div
      :class="
        props.showHeader || props.showFooter
          ? 'flex-1 overflow-y-auto p-6'
          : 'p-4'
      "
    >
      <Collapse
        :active-key="activeKeys"
        @update:active-key="handleActiveKeysChange"
        class="bg-transparent"
        :bordered="false"
      >
        <!-- 基本信息 -->
        <CollapsePanel key="basic" header="基本信息">
          <div class="space-y-4 py-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                节点名称 *
              </label>
              <Input
                :value="formData.name"
                @update:value="(val) => (formData.name = val)"
                placeholder="输入节点名称"
              />
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                编码 *
              </label>
              <Input
                :value="formData.code"
                @update:value="(val) => (formData.code = val)"
                placeholder="输入唯一编码"
              />
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                描述
              </label>
              <Textarea
                :value="formData.description"
                @update:value="(val) => (formData.description = val)"
                placeholder="输入节点描述"
                :rows="3"
              />
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                适配器类型
              </label>
              <Select
                :value="formData.adapterType"
                @update:value="handleAdapterTypeChange"
                class="w-full"
              >
                <Select.Option :value="1">Eino ADK</Select.Option>
                <Select.Option :value="2">Eino DeepADK</Select.Option>
              </Select>
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                类型
              </label>
              <Select
                :value="formData.type"
                @update:value="handleTypeChange"
                class="w-full"
              >
                <Select.Option :value="1">根 Agent</Select.Option>
                <Select.Option :value="2">子 Agent</Select.Option>
              </Select>
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                状态
              </label>
              <Select
                :value="formData.status"
                @update:value="handleStatusChange"
                class="w-full"
              >
                <Select.Option :value="1">启用</Select.Option>
                <Select.Option :value="2">禁用</Select.Option>
              </Select>
            </div>
          </div>
        </CollapsePanel>

        <!-- 模型参数 -->
        <CollapsePanel key="model" header="模型参数">
          <div class="space-y-4 py-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                引用模型
              </label>
              <Select
                :value="formData.originalModelId"
                @update:value="handleOriginalModelIdChange"
                class="w-full"
                placeholder="选择模型（可选）"
                :allow-clear="true"
                :options="modelOptions"
              />
              <p class="mt-1 text-xs text-slate-500">
                选择后将模型参数填充到下方
              </p>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                AI 模型
              </label>
              <div class="grid gap-3">
                <Input
                  :value="formData.modelName"
                  @update:value="(val) => (formData.modelName = val)"
                  placeholder="模型名称"
                />
                <Input
                  :value="formData.modelType"
                  @update:value="(val) => (formData.modelType = val)"
                  placeholder="模型类型"
                />
                <Input
                  :value="formData.baseUrl"
                  @update:value="(val) => (formData.baseUrl = val)"
                  placeholder="Base URL"
                />
                <Input.Password
                  :value="formData.apiKey"
                  @update:value="(val) => (formData.apiKey = val)"
                  placeholder="API Key"
                  :visibility-toggle="true"
                />
              </div>
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                最大迭代次数
              </label>
              <InputNumber
                :value="formData.maxIteration"
                @update:value="handleMaxIterationChange"
                class="w-full"
                :min="1"
                :max="100"
              />
            </div>

            <div>
              <div class="mb-2 flex items-center justify-between">
                <label class="text-sm font-medium text-slate-700">
                  Temperature
                </label>
                <span
                  class="rounded bg-blue-50 px-1.5 py-0.5 font-mono text-xs text-blue-600"
                >
                  {{ formData.temperature }}
                </span>
              </div>
              <Slider
                :value="formData.temperature"
                @update:value="handleTemperatureChange"
                :min="0"
                :max="2"
                :step="0.1"
                class="mb-1"
              />
              <div class="flex justify-between text-[10px] text-slate-400">
                <span>精确 (0.0)</span>
                <span>创造 (2.0)</span>
              </div>
            </div>

            <div>
              <div class="mb-2 flex items-center justify-between">
                <label class="text-sm font-medium text-slate-700">
                  Top P
                </label>
                <span
                  class="rounded bg-blue-50 px-1.5 py-0.5 font-mono text-xs text-blue-600"
                >
                  {{ formData.topP }}
                </span>
              </div>
              <Slider
                :value="formData.topP"
                @update:value="handleTopPChange"
                :min="0"
                :max="1"
                :step="0.1"
              />
            </div>
          </div>
        </CollapsePanel>

        <!-- 提示词配置 -->
        <CollapsePanel key="prompt" header="提示词配置">
          <div class="space-y-4 py-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                系统提示词
              </label>
              <Textarea
                :value="formData.systemPrompt"
                @update:value="(val) => (formData.systemPrompt = val)"
                placeholder="输入系统提示词..."
                :rows="6"
                class="font-mono text-xs"
              />
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-slate-700">
                用户输入提示词
              </label>
              <Textarea
                :value="formData.userInputPrompt"
                @update:value="(val) => (formData.userInputPrompt = val)"
                placeholder="输入用户输入提示词..."
                :rows="6"
                class="font-mono text-xs"
              />
            </div>
          </div>
        </CollapsePanel>

        <!-- 功能特性 -->
        <CollapsePanel key="features" header="功能特性">
          <div class="space-y-3 py-2">
            <div
              class="flex items-center justify-between rounded-lg border border-slate-200 bg-slate-50 px-4 py-3"
            >
              <label class="text-sm font-medium text-slate-700">
                启用待办事项
              </label>
              <Checkbox
                :checked="formData.withWriteTodos"
                @update:checked="(val) => (formData.withWriteTodos = val)"
              />
            </div>

            <div
              class="flex items-center justify-between rounded-lg border border-slate-200 bg-slate-50 px-4 py-3"
            >
              <label class="text-sm font-medium text-slate-700">
                启用网络搜索
              </label>
              <Checkbox
                :checked="formData.withWebSearchAgent"
                @update:checked="(val) => (formData.withWebSearchAgent = val)"
              />
            </div>
          </div>
        </CollapsePanel>
      </Collapse>
    </div>

    <!-- Footer -->
    <div
      v-if="props.showFooter"
      class="flex items-center justify-end gap-3 border-t border-slate-200 bg-slate-50 px-6 py-4"
    >
      <Button @click="emit('close')"> 取消 </Button>
      <Button
        type="primary"
        @click="handleSave"
        class="border-none bg-gradient-to-r from-blue-500 to-blue-600"
      >
        保存更改
      </Button>
    </div>
  </div>
</template>

<style scoped>
:deep(.ant-collapse) {
  background: transparent;
  border: none;
}

:deep(.ant-collapse-item) {
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 0.75rem;
  overflow: hidden;
  background: #ffffff;
  margin-bottom: 12px;
}

:deep(.ant-collapse-header) {
  padding: 12px 16px !important;
  font-weight: 500;
  color: #334155;
  background: linear-gradient(
    90deg,
    rgba(248, 250, 252, 1),
    rgba(248, 250, 252, 0)
  );
}

:deep(.ant-collapse-content) {
  border-top: 1px solid rgba(226, 232, 240, 0.7);
  background: #ffffff;
}

:deep(.ant-input),
:deep(.ant-input-number),
:deep(.ant-select-selector) {
  border-radius: 8px !important;
}

:deep(.ant-slider-rail) {
  background-color: #e2e8f0;
}

:deep(.ant-slider-track) {
  background-color: #3b82f6;
}

:deep(.ant-slider-handle::after) {
  box-shadow: 0 0 0 2px #3b82f6;
}

/* 自定义滚动条 */
.overflow-y-auto::-webkit-scrollbar {
  width: 4px;
}
.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>
