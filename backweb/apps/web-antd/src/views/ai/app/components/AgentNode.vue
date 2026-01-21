<script setup lang="ts">
import type { NodeProps } from '@vue-flow/core';

import { computed } from 'vue';

import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  RobotOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons-vue';
import { Handle, Position } from '@vue-flow/core';

const props = defineProps<NodeProps>();

const agentData = computed(() => props.data || {});

const statusColor = computed(() => {
  return agentData.value.status === 1 ? 'text-green-500' : 'text-red-500';
});

const typeLabel = computed(() => {
  return agentData.value.type === 1 ? 'Root' : 'Sub';
});

const adapterTypeLabel = computed(() => {
  const types: Record<number, string> = {
    1: 'ADK',
    2: 'DeepADK',
  };
  return types[agentData.value.adapterType] || 'Unknown';
});

const descriptionPreview = computed(() => {
  const text = agentData.value.description || '';
  return text.length > 10 ? `${text.slice(0, 10)}...` : text;
});
</script>

<template>
  <div
    class="agent-node min-w-[200px] rounded-xl border-2 border-slate-200 bg-white shadow-lg transition-all hover:shadow-xl"
    :class="{
      'border-blue-400': props.selected,
    }"
  >
    <!-- Input Handle -->
    <Handle
      type="target"
      :position="Position.Left"
      class="!h-3 !w-3 !border-2 !border-white !bg-blue-500"
    />

    <!-- Header -->
    <div
      class="border-b border-slate-100 bg-gradient-to-r from-blue-50 to-transparent px-4 py-3"
    >
      <div class="flex items-center gap-2">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600"
        >
          <RobotOutlined class="text-sm text-white" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="truncate text-sm font-semibold text-slate-800">
            {{ agentData.name || '未命名节点' }}
          </div>
          <div class="text-xs text-slate-500">
            {{ agentData.code || '-' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Body -->
    <div class="space-y-2 px-4 py-3">
      <div class="flex items-center justify-between text-xs">
        <span class="text-slate-500">类型:</span>
        <span class="rounded bg-blue-50 px-2 py-0.5 font-medium text-blue-700">
          {{ typeLabel }}
        </span>
      </div>

      <div class="flex items-center justify-between text-xs">
        <span class="text-slate-500">适配器:</span>
        <span
          class="rounded bg-indigo-50 px-2 py-0.5 font-medium text-indigo-700"
        >
          {{ adapterTypeLabel }}
        </span>
      </div>

      <div class="flex items-center justify-between text-xs">
        <span class="text-slate-500">状态:</span>
        <div class="flex items-center gap-1">
          <CheckCircleOutlined
            v-if="agentData.status === 1"
            :class="statusColor"
          />
          <CloseCircleOutlined v-else :class="statusColor" />
          <span :class="statusColor">
            {{ agentData.status === 1 ? '启用' : '禁用' }}
          </span>
        </div>
      </div>

      <div
        v-if="agentData.description"
        class="truncate border-t border-slate-100 pt-2 text-xs text-slate-400"
        :title="agentData.description"
      >
        {{ descriptionPreview }}
      </div>

      <!-- Features -->
      <div
        v-if="agentData.withWriteTodos || agentData.withWebSearchAgent"
        class="flex gap-1 pt-2"
      >
        <span
          v-if="agentData.withWriteTodos"
          class="rounded bg-green-50 px-1.5 py-0.5 text-[10px] text-green-600"
        >
          待办
        </span>
        <span
          v-if="agentData.withWebSearchAgent"
          class="rounded bg-orange-50 px-1.5 py-0.5 text-[10px] text-orange-600"
        >
          搜索
        </span>
      </div>
    </div>

    <!-- Footer -->
    <div
      v-if="agentData.subAgents && agentData.subAgents.length > 0"
      class="rounded-b-xl border-t border-slate-100 bg-slate-50 px-4 py-2"
    >
      <div class="flex items-center gap-1 text-xs text-slate-500">
        <ThunderboltOutlined class="text-[10px]" />
        <span>{{ agentData.subAgents.length }} 个子 Agent</span>
      </div>
    </div>

    <!-- Output Handle -->
    <Handle
      type="source"
      :position="Position.Right"
      class="!h-3 !w-3 !border-2 !border-white !bg-green-500"
    />
  </div>
</template>

<style scoped>
.agent-node {
  cursor: pointer;
}

</style>
