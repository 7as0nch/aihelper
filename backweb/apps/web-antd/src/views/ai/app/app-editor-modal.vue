<script setup lang="ts">
import type { Edge, Node } from '@vue-flow/core';

import type { AppAgent, ApplicationInfo } from '#/api/ai/application/model';

import { computed, onMounted, ref, watch } from 'vue';

import {
  ApartmentOutlined,
  CloseOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SaveOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue';
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import { useVueFlow, VueFlow } from '@vue-flow/core';
import { message, Modal } from 'ant-design-vue';

import {
  applicationAdd,
  applicationInfo,
  applicationUpdate,
} from '#/api/ai/application';

import AgentLibrary from './components/AgentLibrary.vue';
import AgentNode from './components/AgentNode.vue';
import NodeConfigPanel from './components/NodeConfigPanel.vue';

interface AppEditorModalProps {
  visible?: boolean;
  applicationId?: number | string;
}

const props = withDefaults(defineProps<AppEditorModalProps>(), {
  visible: false,
});

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

// VueFlow 实例
const { onConnect, addEdges, onNodeClick, getSelectedNodes, fitView } =
  useVueFlow();

// 弹框状态
const isDragging = ref(false);

// 初始状态下使用 transform 居中
const modalStyle = ref<any>({
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  margin: '0',
});

// 应用数据初始化
const applicationData = ref<Partial<ApplicationInfo>>({
  name: '新应用',
  code: '',
  description: '',
  version: '1.0.0',
  mode: 1,
  status: 1,
  type: 2,
  scope: 1,
  selfAgent: {
    name: '默认 Agent',
    code: 'default',
    adapterType: 1,
    status: 1,
    type: 1,
  },
});

// 流程图数据
const nodes = ref<Node[]>([]);
const edges = ref<Edge[]>([]);
const selectedNode = ref<Node | null>(null);
const showAppSettings = ref(false);

// UI 状态
const loading = ref(false);
const saving = ref(false);
const isRightPanelVisible = ref(false);
const isLeftPanelCollapsed = ref(false);

// 是否为编辑模式
const isEditMode = computed(() => !!props.applicationId);

// --- 实用工具 ---

// 递归转换对象键名为 snake_case
const toSnakeCase = (obj: any): any => {
  if (Array.isArray(obj)) return obj.map(toSnakeCase);
  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj).reduce((acc: any, key) => {
      const snakeKey = key.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);
      acc[snakeKey] = toSnakeCase(obj[key]);
      return acc;
    }, {});
  }
  return obj;
};

// 递归转换对象键名为 camelCase
const toCamelCase = (obj: any): any => {
  if (Array.isArray(obj)) return obj.map(toCamelCase);
  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj).reduce((acc: any, key) => {
      const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
      acc[camelKey] = toCamelCase(obj[key]);
      return acc;
    }, {});
  }
  return obj;
};

// 从 schema 构建自定义 Agent 树（subAgents 维持父子关系）
const buildAgentTreeFromSchema = (): AppAgent | undefined => {
  if (nodes.value.length === 0) return undefined;

  const nodeDataMap = new Map<string, any>();
  nodes.value.forEach((node) => {
    nodeDataMap.set(node.id, { ...node.data });
  });

  const childrenMap = new Map<string, string[]>();
  const incoming = new Set<string>();
  edges.value.forEach((edge) => {
    if (!childrenMap.has(edge.source)) {
      childrenMap.set(edge.source, []);
    }
    childrenMap.get(edge.source)!.push(edge.target);
    incoming.add(edge.target);
  });

  const rootIds = nodes.value
    .map((node) => node.id)
    .filter((id) => !incoming.has(id));

  const rootId =
    rootIds.find((id) => nodeDataMap.get(id)?.type === 1) ?? rootIds[0];
  if (!rootId) return undefined;

  const build = (id: string): AppAgent => {
    const data = nodeDataMap.get(id) || {};
    const children = childrenMap.get(id) || [];
    const subAgents = children.map(build);
    const agent = { ...data } as AppAgent;
    if (subAgents.length > 0) {
      agent.subAgents = subAgents;
    }
    return agent;
  };

  const rootAgent = build(rootId);

  // 多根场景下，将其它根节点作为 subAgents 挂载
  if (rootIds.length > 1) {
    const extraRoots = rootIds.filter((id) => id !== rootId).map(build);
    if (extraRoots.length > 0) {
      rootAgent.subAgents = [...(rootAgent.subAgents || []), ...extraRoots];
    }
  }

  return rootAgent;
};

// --- 拖拽功能实现 ---
const dragState = ref({
  startX: 0,
  startY: 0,
  startLeft: 0,
  startTop: 0,
});

const handleDragStart = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (target.closest('button') || target.closest('.ant-btn')) return;

  isDragging.value = true;

  // 捕获当前像素位置并清除 transform
  const modalElement = document.querySelector('.draggable-modal .ant-modal') as HTMLElement;
  if (modalElement) {
    const rect = modalElement.getBoundingClientRect();
    modalStyle.value = {
      top: `${rect.top}px`,
      left: `${rect.left}px`,
      transform: 'none',
      margin: '0',
      paddingBottom: '0',
      position: 'absolute',
    };

    dragState.value = {
      startX: e.clientX,
      startY: e.clientY,
      startLeft: rect.left,
      startTop: rect.top,
    };
  }

  document.addEventListener('mousemove', handleDragMove);
  document.addEventListener('mouseup', handleDragEnd);
};

const handleDragMove = (e: MouseEvent) => {
  if (!isDragging.value) return;
  const deltaX = e.clientX - dragState.value.startX;
  const deltaY = e.clientY - dragState.value.startY;
  modalStyle.value.left = `${dragState.value.startLeft + deltaX}px`;
  modalStyle.value.top = `${dragState.value.startTop + deltaY}px`;
};

const handleDragEnd = () => {
  isDragging.value = false;
  document.removeEventListener('mousemove', handleDragMove);
  document.removeEventListener('mouseup', handleDragEnd);
};

// --- 功能函数 ---

const handleClose = () => {
  emit('update:visible', false);
};

const loadApplication = async () => {
  if (!props.applicationId) {
    applicationData.value = {
      name: '新应用', code: '', description: '', version: '1.0.0',
      mode: 1, status: 1, type: 2, scope: 1,
      selfAgent: { name: '默认 Agent', code: 'default', adapterType: 1, status: 1, type: 1 }
    };
    nodes.value = [];
    edges.value = [];
    selectedNode.value = null;
    return;
  }

  try {
    loading.value = true;
    const data = await applicationInfo(props.applicationId as string);
    applicationData.value = toCamelCase(data);
    if (applicationData.value.schema) {
      const schema = JSON.parse(applicationData.value.schema);
      nodes.value = schema.nodes || [];
      edges.value = schema.edges || [];
      setTimeout(() => fitView({ padding: 0.2 }), 100);
    }
  } catch (error) {
    console.error('加载失败:', error);
    message.error('加载应用失败');
  } finally {
    loading.value = false;
  }
};

const handleSave = async () => {
  try {
    saving.value = true;
    const schema = JSON.stringify({
      nodes: nodes.value.map((n) => ({
        id: n.id,
        type: n.type,
        position: n.position,
        data: n.data,
      })),
      edges: edges.value.map((e) => ({
        id: e.id,
        source: e.source,
        target: e.target,
        type: e.type,
      })),
    });

    // schema 的 data 拆出为 selfAgent（通过 subAgents 保持层级关系）
    const derivedSelfAgent = buildAgentTreeFromSchema();
    applicationData.value.selfAgent = derivedSelfAgent;

    // 将数据转换为 snake_case 以符合后端要求
    const saveData = toSnakeCase({ ...applicationData.value, schema });

    if (isEditMode.value) {
      await applicationUpdate(saveData as any);
      message.success('保存成功');
    } else {
      await applicationAdd(saveData as any);
      message.success('创建成功');
    }
    emit('success');
    handleClose();
  } catch (error) {
    console.error('保存失败:', error);
    message.error('保存失败');
  } finally {
    saving.value = false;
  }
};

// --- 事件监听 ---

onConnect((params) => addEdges([params]));

onNodeClick(({ node }) => {
  selectedNode.value = node;
  showAppSettings.value = false;
  isRightPanelVisible.value = true;
});

const onPaneClick = () => {
  selectedNode.value = null;
  showAppSettings.value = false;
  isRightPanelVisible.value = false;
};

const handleOpenAppSettings = () => {
  selectedNode.value = null;
  showAppSettings.value = true;
  isRightPanelVisible.value = true;
};

const handleUpdateNode = (data: any) => {
  if (!selectedNode.value) return;
  const node = nodes.value.find(n => n.id === selectedNode.value!.id);
  if (node) {
    node.data = { ...node.data, ...data };
  }
};

const handleAddNode = (agentData: AppAgent) => {
  const newNode: Node = {
    id: `node-${Date.now()}`,
    type: 'agentNode',
    position: { x: Math.random() * 400 + 100, y: Math.random() * 300 + 100 },
    data: agentData,
  };
  nodes.value.push(newNode);
  message.success(`已添加节点: ${agentData.name}`);
};

const handleDeleteSelected = () => {
  const selected = getSelectedNodes.value;
  if (selected.length === 0) return message.warning('请先选择节点');
  nodes.value = nodes.value.filter(n => !selected.some(s => s.id === n.id));
  edges.value = edges.value.filter(e => !selected.some(s => s.id === e.source || s.id === e.target));
  selectedNode.value = null;
  message.success(`已删除 ${selected.length} 个节点`);
};

watch(() => props.visible, (newVal) => { if (newVal) loadApplication(); });
onMounted(() => { if (props.visible) loadApplication(); });
</script>

<template>
  <Modal
    :open="props.visible"
    :footer="null"
    :closable="false"
    :mask-closable="false"
    width="90%"
    wrap-class-name="draggable-modal"
    :body-style="{ padding: 0, height: '85vh' }"
    :style="modalStyle"
    @cancel="handleClose"
  >
    <div class="flex h-full flex-col bg-slate-50">
      <!-- Header -->
      <div
        class="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-6 shadow-sm cursor-move"
        @mousedown="handleDragStart"
      >
        <div class="flex items-center gap-4">
          <ApartmentOutlined
            class="cursor-move text-slate-400 text-lg"
          />
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
          <a-button @click="handleOpenAppSettings" size="small">
            <template #icon><SettingOutlined /></template>
            应用设置
          </a-button>
          <a-button @click="handleDeleteSelected" danger size="small">
            删除选中节点
          </a-button>
          <a-button
            type="primary"
            size="small"
            :loading="saving"
            @click="handleSave"
          >
            <template #icon><SaveOutlined /></template>
            保存
          </a-button>
          <a-button size="small" @click="handleClose">
            <template #icon><CloseOutlined /></template>
          </a-button>
        </div>
      </div>

      <!-- Main Content -->
      <div class="relative flex min-h-0 flex-1 overflow-x-hidden">
        <!-- Left Panel - Agent Library -->
        <div
          class="border-r border-slate-200 bg-white transition-all duration-300"
          :class="[isLeftPanelCollapsed ? 'w-0 overflow-hidden' : 'w-64']"
        >
          <AgentLibrary
            v-if="!isLeftPanelCollapsed"
            @addNode="handleAddNode"
          />
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
            @pane-click="onPaneClick"
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
            <MenuFoldOutlined v-if="!isLeftPanelCollapsed" />
            <MenuUnfoldOutlined v-else />
          </button>
        </div>

        <!-- Right Panel - Config (Floating) -->
        <div
          class="absolute top-20 right-4 z-20 w-[400px] max-h-[calc(100%-6rem)] rounded-2xl border border-slate-200 bg-white shadow-2xl transition-all duration-300 flex flex-col"
          :class="[isRightPanelVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0 pointer-events-none']"
        >
          <div v-if="isRightPanelVisible" class="flex flex-col h-full min-h-0">
            <div
              class="flex items-center justify-between border-b border-slate-100 p-4 shrink-0"
            >
              <h3 class="m-0 text-base font-semibold text-slate-800">
                {{ selectedNode ? '节点配置' : '应用设置' }}
              </h3>
              <a-button type="text" size="small" @click="onPaneClick">
                <template #icon><CloseOutlined /></template>
              </a-button>
            </div>

            <div class="flex-1 overflow-y-auto px-4 py-4 custom-scrollbar">
              <!-- 节点配置面板 -->
              <NodeConfigPanel
                v-if="selectedNode"
                :node="selectedNode"
                :show-header="false"
                :show-footer="false"
                :auto-save="true"
                @update="handleUpdateNode"
                @close="onPaneClick"
              />

              <!-- 应用全局配置面板 -->
              <div v-else-if="showAppSettings">
                <a-form layout="vertical">
                  <a-form-item label="应用名称">
                    <a-input
                      :value="applicationData.name"
                      @update:value="(val) => (applicationData.name = val)"
                    />
                  </a-form-item>
                  <a-form-item label="应用编码">
                    <a-input
                      :value="applicationData.code"
                      @update:value="(val) => (applicationData.code = val)"
                    />
                  </a-form-item>
                  <a-form-item label="描述">
                    <a-textarea
                      :value="applicationData.description"
                      @update:value="(val) => (applicationData.description = val)"
                    />
                  </a-form-item>
                <a-form-item label="版本">
                  <a-input
                    :value="applicationData.version"
                    @update:value="(val) => (applicationData.version = val)"
                  />
                </a-form-item>
                  <a-form-item label="模式">
                    <a-select
                      :value="applicationData.mode"
                      @update:value="(val) => (applicationData.mode = Number(val))"
                    >
                      <a-select-option :value="1">单 Agent 模式</a-select-option>
                      <a-select-option :value="2">多 Agent 模式</a-select-option>
                    </a-select>
                  </a-form-item>
                </a-form>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Modal>
</template>

<style>
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';
</style>

<style scoped>
:deep(.draggable-modal .ant-modal) {
  padding: 0;
  position: absolute;
}

:deep(.draggable-modal) {
  overflow-x: hidden;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}

.vue-flow__node {
  cursor: pointer;
}

:deep(.vue-flow__node.selected) {
  box-shadow: 0 0 0 2px #3b82f6;
}
</style>
