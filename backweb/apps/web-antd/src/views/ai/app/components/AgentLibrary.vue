<script setup lang="ts">
import type { AgentInfo } from '#/api/ai/agent/model';
import type { AppAgent } from '#/api/ai/application/model';

import { computed, onMounted, ref } from 'vue';

import {
  AppstoreAddOutlined,
  PlusOutlined,
  RobotOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue';
import { Input, Spin } from 'ant-design-vue';

import { agentList } from '#/api/ai/agent';

const emit = defineEmits<{
  (e: 'addNode', agent: AppAgent): void;
}>();

const agents = ref<AgentInfo[]>([]);
const loading = ref(false);
const searchKeyword = ref('');

const loadAgents = async () => {
  try {
    loading.value = true;
    const result = await agentList({ pageNum: 1, pageSize: 100 });
    agents.value = result.list || [];
  } catch (error) {
    console.error('加载 Agent 列表失败:', error);
  } finally {
    loading.value = false;
  }
};

const filteredAgents = computed(() => {
  if (!searchKeyword.value) return agents.value;
  const keyword = searchKeyword.value.toLowerCase();
  return agents.value.filter(
    (agent) =>
      agent.name.toLowerCase().includes(keyword) ||
      agent.code.toLowerCase().includes(keyword) ||
      agent.description?.toLowerCase().includes(keyword),
  );
});

const handleAddAgent = (agent: AgentInfo) => {
  // 转换 AgentInfo 为 AppAgent 格式
  const appAgent: AppAgent = {
    id: agent.id,
    name: agent.name,
    code: agent.code,
    description: agent.description,
    adapterType: agent.adapterType,
    originalModelId: agent.originalModelId,
    aiModel: agent.aiModel,
    maxIteration: agent.maxIteration,
    systemPrompt: agent.systemPrompt,
    userInputPrompt: agent.userInputPrompt,
    status: agent.status,
    type: agent.type,
    withWriteTodos: agent.withWriteTodos,
    withWebSearchAgent: agent.withWebSearchAgent,
    systemType: agent.systemType,
    createdAt: agent.createdAt,
  };

  emit('addNode', appAgent);
};

const handleAddEmptyNode = () => {
  // 添加空白节点
  const emptyAgent: AppAgent = {
    name: '新节点',
    code: `node_${Date.now()}`,
    description: '',
    adapterType: 1,
    status: 1,
    type: 2,
  };

  emit('addNode', emptyAgent);
};

onMounted(() => {
  loadAgents();
});
</script>

<template>
  <div class="flex h-full flex-col bg-white">
    <!-- Header -->
    <div class="border-b border-slate-200 px-4 py-4">
      <div class="mb-3 flex items-center gap-2">
        <AppstoreAddOutlined class="text-lg text-blue-600" />
        <h3 class="m-0 font-semibold text-slate-800">节点库</h3>
      </div>

      <!-- Search -->
      <Input
        :model-value="searchKeyword"
        placeholder="搜索 Agent..."
        size="small"
      >
        <template #prefix>
          <SearchOutlined class="text-slate-400" />
        </template>
      </Input>
    </div>

    <!-- Add Empty Node Button -->
    <div class="border-b border-slate-200 px-4 py-3">
      <button
        @click="handleAddEmptyNode"
        class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg border-2 border-dashed border-slate-300 bg-white px-3 py-2 text-sm text-slate-600 transition-all hover:border-blue-400 hover:bg-blue-50 hover:text-blue-600"
      >
        <PlusOutlined />
        <span>添加空白节点</span>
      </button>
    </div>

    <!-- Agent List -->
    <div class="flex-1 space-y-2 overflow-y-auto px-4 py-3">
      <Spin :spinning="loading">
        <div
          v-if="filteredAgents.length === 0"
          class="py-8 text-center text-sm text-slate-400"
        >
          {{ searchKeyword ? '未找到匹配的 Agent' : '暂无可用 Agent' }}
        </div>

        <div
          v-for="agent in filteredAgents"
          :key="agent.id"
          @click="handleAddAgent(agent)"
          class="agent-card group cursor-pointer rounded-lg border border-slate-200 bg-slate-50 p-3 transition-all hover:border-blue-300 hover:bg-blue-50"
        >
          <div class="flex items-start gap-2">
            <div
              class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600"
            >
              <RobotOutlined class="text-xs text-white" />
            </div>
            <div class="min-w-0 flex-1">
              <div
                class="truncate text-sm font-medium text-slate-800 group-hover:text-blue-600"
              >
                {{ agent.name }}
              </div>
              <div class="truncate text-xs text-slate-500">
                {{ agent.code }}
              </div>
              <div
                v-if="agent.description"
                class="mt-1 line-clamp-2 text-xs text-slate-400"
              >
                {{ agent.description }}
              </div>
              <div class="mt-2 flex gap-1">
                <span
                  v-if="agent.type === 1"
                  class="rounded bg-blue-100 px-1.5 py-0.5 text-[10px] text-blue-700"
                >
                  根 Agent
                </span>
                <span
                  v-if="agent.status === 1"
                  class="rounded bg-green-100 px-1.5 py-0.5 text-[10px] text-green-700"
                >
                  启用
                </span>
              </div>
            </div>
          </div>
        </div>
      </Spin>
    </div>

    <!-- Footer Tips -->
    <div class="border-t border-slate-200 bg-slate-50 px-4 py-3">
      <p class="m-0 text-xs text-slate-500">💡 点击 Agent 卡片添加到画布</p>
    </div>
  </div>
</template>

<style scoped>
.agent-card {
  transition: all 0.2s ease;
}

.agent-card:active {
  transform: scale(0.98);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
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
