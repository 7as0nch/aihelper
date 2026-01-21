<script setup lang="ts">
import type { Edge, Node } from '@vue-flow/core';

import type { AppAgent, ApplicationInfo } from '#/api/ai/application/model';

import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import {
  ArrowLeftOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  SaveOutlined,
} from '@ant-design/icons-vue';
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import { useVueFlow, VueFlow } from '@vue-flow/core';
import { message } from 'ant-design-vue';

import {
  applicationAdd,
  applicationInfo,
  applicationUpdate,
} from '#/api/ai/application';

import AgentLibrary from './components/AgentLibrary.vue';
import AgentNode from './components/AgentNode.vue';
import NodeConfigPanel from './components/NodeConfigPanel.vue';

const route = useRoute();
const router = useRouter();

// VueFlow 实例
const { onConnect, addEdges, onNodeClick, getSelectedNodes, fitView } =
  useVueFlow();

// 应用数据
const applicationData = ref<Partial<ApplicationInfo>>({
  name: '新应用',
  code: '',
  description: '',
  version: '1.0.0',
  mode: 1,
  status: 1,
  type: 2,
  scope: 1,
});

// 流程图数据
const nodes = ref<Node[]>([]);
const edges = ref<Edge[]>([]);
const selectedNode = ref<Node | null>(null);

// UI 状态
const loading = ref(false);
const saving = ref(false);
const isRightPanelCollapsed = ref(false);
const isLeftPanelCollapsed = ref(false);

// 是否为编辑模式
const isEditMode = computed(() => !!route.params.id);

// 加载应用数据
const loadApplication = async () => {
  if (!route.params.id) return;

  try {
    loading.value = true;
    const data = await applicationInfo(route.params.id as string);
    applicationData.value = data;

    // 解析 schema
    if (data.schema) {
      try {
        const schema = JSON.parse(data.schema);
        nodes.value = schema.nodes || [];
        edges.value = schema.edges || [];

        // 延迟执行 fitView 以确保画布已渲染
        setTimeout(() => {
          fitView({ padding: 0.2 });
        }, 100);
      } catch (error) {
        console.error('解析 schema 失败:', error);
        message.error('加载流程图失败');
      }
    }
  } catch (error) {
    console.error('加载应用失败:', error);
    message.error('加载应用失败');
  } finally {
    loading.value = false;
  }
};

// 保存应用
const handleSave = async () => {
  try {
    saving.value = true;

    // 序列化流程图为 schema
    const schema = JSON.stringify({
      nodes: nodes.value.map((node) => ({
        id: node.id,
        type: node.type,
        position: node.position,
        data: node.data,
      })),
      edges: edges.value.map((edge) => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        type: edge.type,
      })),
    });

    const saveData = {
      ...applicationData.value,
      schema,
    };

    if (isEditMode.value) {
      await applicationUpdate(saveData as any);
      message.success('保存成功');
    } else {
      const result = await applicationAdd(saveData as any);
      message.success('创建成功');
      // 创建成功后跳转到编辑模式
      router.replace(`/ai/app/editor/${result.id}`);
    }
  } catch (error) {
    console.error('保存失败:', error);
    message.error('保存失败');
  } finally {
    saving.value = false;
  }
};

// 返回列表
const handleBack = () => {
  router.push('/ai/app');
};

// 连接节点
onConnect((params) => {
  addEdges([params]);
});

// 节点点击事件
onNodeClick(({ node }) => {
  selectedNode.value = node;
  isRightPanelCollapsed.value = false;
});

// 更新节点数据
const handleUpdateNode = (data: any) => {
  if (!selectedNode.value) return;

  const nodeIndex = nodes.value.findIndex(
    (n) => n.id === selectedNode.value!.id,
  );
  if (nodeIndex !== -1) {
    const node = nodes.value[nodeIndex];
    if (node) {
      node.data = {
        ...node.data,
        ...data,
      };
    }
  }

  message.success('节点配置已更新');
};

// 关闭配置面板
const handleClosePanel = () => {
  selectedNode.value = null;
};

// 从节点库添加节点
const handleAddNode = (agentData: AppAgent) => {
  const newNode: Node = {
    id: `node-${Date.now()}`,
    type: 'agentNode',
    position: {
      x: Math.random() * 400 + 100,
      y: Math.random() * 300 + 100,
    },
    data: agentData,
  };

  nodes.value.push(newNode);
  message.success(`已添加节点: ${agentData.name}`);
};

// 删除选中节点
const handleDeleteSelected = () => {
  const selected = getSelectedNodes.value;
  if (selected.length === 0) {
    message.warning('请先选择要删除的节点');
    return;
  }

  nodes.value = nodes.value.filter((n) => !selected.some((s) => s.id === n.id));
  edges.value = edges.value.filter(
    (e) => !selected.some((s) => s.id === e.source || s.id === e.target),
  );

  selectedNode.value = null;
  message.success(`已删除 ${selected.length} 个节点`);
};

onMounted(() => {
  loadApplication();
});
</script>

<template>
  <div class="flex h-screen flex-col bg-slate-50">
    <!-- Header -->
    <div
      class="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-6 shadow-sm"
    >
      <div class="flex items-center gap-4">
        <button
          @click="handleBack"
          class="cursor-pointer rounded-lg border-none bg-transparent p-2 transition-colors hover:bg-slate-100"
        >
          <ArrowLeftOutlined class="text-slate-600" />
        </button>
        <div>
          <h1 class="m-0 text-lg font-semibold text-slate-800">
            {{ applicationData.name }}
          </h1>
          <p class="m-0 text-xs text-slate-500">
            {{ isEditMode ? '编辑应用流程' : '创建新应用' }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <a-button @click="handleDeleteSelected" danger> 删除选中节点 </a-button>
        <a-button type="primary" :loading="saving" @click="handleSave">
          <SaveOutlined />
          保存
        </a-button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex min-h-0 flex-1">
      <!-- Left Panel - Agent Library -->
      <div
        class="border-r border-slate-200 bg-white transition-all duration-300"
        :class="[isLeftPanelCollapsed ? 'w-0 overflow-hidden' : 'w-64']"
      >
        <AgentLibrary v-if="!isLeftPanelCollapsed" @add-node="handleAddNode" />
      </div>

      <!-- Center - VueFlow Canvas -->
      <div class="relative flex-1">
        <VueFlow
          :nodes="nodes"
          :edges="edges"
          :default-zoom="1"
          :min-zoom="0.2"
          :max-zoom="4"
          class="bg-slate-50"
          @update:nodes="(val) => (nodes = val)"
          @update:edges="(val) => (edges = val)"
        >
          <Background pattern-color="#e2e8f0" :gap="16" />
          <Controls />

          <!-- Custom Node Templates -->
          <template #node-agentNode="nodeProps">
            <AgentNode v-bind="nodeProps" />
          </template>
        </VueFlow>

        <!-- Toggle Left Panel Button -->
        <button
          @click="isLeftPanelCollapsed = !isLeftPanelCollapsed"
          class="absolute left-4 top-4 z-10 rounded-lg border border-slate-200 bg-white p-2 shadow-sm transition-colors hover:bg-slate-50"
        >
          <FullscreenExitOutlined v-if="!isLeftPanelCollapsed" />
          <FullscreenOutlined v-else />
        </button>
      </div>

      <!-- Right Panel - Node Config -->
      <div
        class="border-l border-slate-200 bg-white transition-all duration-300"
        :class="[isRightPanelCollapsed ? 'w-0 overflow-hidden' : 'w-96']"
      >
        <NodeConfigPanel
          v-if="selectedNode && !isRightPanelCollapsed"
          :node="selectedNode"
          @update="handleUpdateNode"
          @close="handleClosePanel"
        />
        <div
          v-else-if="!isRightPanelCollapsed"
          class="flex h-full items-center justify-center p-6 text-center text-sm text-slate-400"
        >
          点击节点以配置其属性
        </div>
      </div>
    </div>
  </div>
</template>

<style>
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';

.vue-flow__node {
  cursor: pointer;
}

.vue-flow__node.selected {
  box-shadow: 0 0 0 2px #3b82f6;
}
</style>
