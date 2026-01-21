<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import type { Node } from '@vue-flow/core';
import {
  CloseOutlined,
  SettingOutlined,
  ControlOutlined,
  MessageOutlined,
  CodeOutlined,
} from '@ant-design/icons-vue';
import {
  Input,
  Textarea,
  Select,
  Slider,
  InputNumber,
  Button,
} from 'ant-design-vue';

interface NodeConfigPanelProps {
  node: Node;
}

const props = defineProps<NodeConfigPanelProps>();

const emit = defineEmits<{
  (e: 'update', data: any): void;
  (e: 'close'): void;
}>();

const formData = ref({
  name: '',
  description: '',
  temperature: 0.7,
  topP: 0.9,
  maxTokens: 2000,
  systemPrompt: '',
  status: 1,
});

const initFormData = () => {
  const data = props.node.data || {};
  formData.value = {
    name: data.name || data.label || '',
    description: data.description || '',
    temperature: data.aiModel?.temperature ?? data.temperature ?? 0.7,
    topP: data.aiModel?.topP ?? data.topP ?? 0.9,
    maxTokens: data.maxTokens || 2000,
    systemPrompt: data.systemPrompt || '',
    status: data.status || 1,
  };
};

watch(() => props.node, initFormData, { deep: true });

onMounted(initFormData);

const handleSave = () => {
  // 按照 proto 结构整理数据
  const updateData = {
    ...formData.value,
    aiModel: {
      temperature: formData.value.temperature,
      topP: formData.value.topP,
    },
  };
  emit('update', updateData);
};
</script>

<template>
  <div
    class="absolute top-0 right-0 h-full w-80 bg-white border-l border-slate-200 shadow-2xl z-10 flex flex-col"
  >
    <!-- Header -->
    <div
      class="px-6 py-4 border-b border-slate-200 flex items-center justify-between bg-gradient-to-r from-slate-50 to-white"
    >
      <div class="flex items-center gap-2">
        <SettingOutlined class="text-slate-700 text-lg" />
        <h2 class="font-semibold text-slate-800 m-0">节点配置</h2>
      </div>
      <button
        @click="emit('close')"
        class="p-1 hover:bg-slate-100 rounded-lg transition-colors border-none bg-transparent cursor-pointer"
      >
        <CloseOutlined class="text-slate-600" />
      </button>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-6 space-y-6">
      <!-- 基本信息 -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 text-sm font-medium text-slate-700">
          <MessageOutlined class="text-blue-500" />
          基本信息
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">
            节点名称
          </label>
          <Input
            :value="formData.name"
            @update:value="(val) => (formData.name = val)"
            placeholder="输入节点名称"
            class="rounded-lg"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">
            描述
          </label>
          <Textarea
            :value="formData.description"
            @update:value="(val) => (formData.description = val)"
            placeholder="输入节点描述"
            :rows="3"
            class="rounded-lg resize-none"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">
            状态
          </label>
          <Select
            :value="formData.status"
            @update:value="(val) => (formData.status = val)"
            class="w-full"
          >
            <Select.Option :value="1">启用</Select.Option>
            <Select.Option :value="2">禁用</Select.Option>
          </Select>
        </div>
      </div>

      <!-- 模型参数 -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 text-sm font-medium text-slate-700">
          <ControlOutlined class="text-indigo-500" />
          模型参数
        </div>

        <div>
          <div class="flex justify-between items-center mb-2">
            <label class="text-sm font-medium text-slate-700">
              Temperature
            </label>
            <span class="text-xs font-mono text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded">
              {{ formData.temperature }}
            </span>
          </div>
          <Slider
            :value="formData.temperature"
            @update:value="(val) => (formData.temperature = val)"
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
          <div class="flex justify-between items-center mb-2">
            <label class="text-sm font-medium text-slate-700">
              Top P
            </label>
            <span class="text-xs font-mono text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded">
              {{ formData.topP }}
            </span>
          </div>
          <Slider
            :value="formData.topP"
            @update:value="(val) => (formData.topP = val)"
            :min="0"
            :max="1"
            :step="0.1"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">
            最大 Tokens
          </label>
          <InputNumber
            :value="formData.maxTokens"
            @update:value="(val) => (formData.maxTokens = val)"
            class="w-full rounded-lg"
            :min="1"
            :max="128000"
          />
        </div>
      </div>

      <!-- 提示词 -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 text-sm font-medium text-slate-700">
          <CodeOutlined class="text-orange-500" />
          系统提示词
        </div>

        <div>
          <Textarea
            :value="formData.systemPrompt"
            @update:value="(val) => (formData.systemPrompt = val)"
            placeholder="输入系统提示词..."
            :rows="8"
            class="rounded-lg font-mono text-xs"
          />
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div
      class="px-6 py-4 border-t border-slate-200 bg-slate-50 flex items-center justify-end gap-3"
    >
      <Button @click="emit('close')" class="rounded-lg">
        取消
      </Button>
      <Button
        type="primary"
        @click="handleSave"
        class="rounded-lg bg-gradient-to-r from-blue-500 to-blue-600 border-none shadow-md shadow-blue-500/20"
      >
        保存更改
      </Button>
    </div>
  </div>
</template>

<style scoped>
:deep(.ant-input),
:deep(.ant-input-number),
:deep(.ant-select-selector),
:deep(.ant-btn) {
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
