<script setup lang="ts">
import { ref, watch } from 'vue';
import { Tabs } from 'ant-design-vue';
import PromptEditor from '#/components/PromptEditor/index.vue';

const props = defineProps<{
  systemPrompt?: string;
  userInputPrompt?: string;
  activeTab?: string;
}>();

const emit = defineEmits<{
  'update:systemPrompt': [value: string];
  'update:userInputPrompt': [value: string];
  'update:activeTab': [value: string];
}>();

const systemPromptValue = ref(props.systemPrompt || '');
const userInputPromptValue = ref(props.userInputPrompt || '');

watch(
  () => props.systemPrompt,
  (val) => {
    systemPromptValue.value = val || '';
  },
);

watch(
  () => props.userInputPrompt,
  (val) => {
    userInputPromptValue.value = val || '';
  },
);

const handleSystemPromptChange = (value: string) => {
  systemPromptValue.value = value;
  emit('update:systemPrompt', value);
};

const handleUserInputPromptChange = (value: string) => {
  userInputPromptValue.value = value;
  emit('update:userInputPrompt', value);
};
</script>

<template>
  <div class="h-full flex flex-col bg-white rounded-xl border border-slate-200/60 overflow-hidden">
    <Tabs
      :activeKey="props.activeTab"
      class="flex-1 flex flex-col"
      @update:activeKey="(val) => emit('update:activeTab', String(val))"
    >
      <Tabs.TabPane key="system" tab="系统提示词">
        <div class="h-full flex flex-col">
          <div class="flex-1 min-h-0 p-4">
            <PromptEditor
              :model-value="systemPromptValue"
              @update:model-value="handleSystemPromptChange"
              class="h-full"
            />
          </div>
        </div>
      </Tabs.TabPane>
      <Tabs.TabPane key="user" tab="用户输入提示词">
        <div class="h-full flex flex-col">
          <div class="flex-1 min-h-0 p-4">
            <PromptEditor
              :model-value="userInputPromptValue"
              @update:model-value="handleUserInputPromptChange"
              class="h-full"
            />
          </div>
        </div>
      </Tabs.TabPane>
    </Tabs>
  </div>
</template>

<style scoped>
:deep(.ant-tabs-content-holder) {
  flex: 1;
  overflow: hidden;
}

:deep(.ant-tabs-tabpane) {
  height: 100%;
  overflow: hidden;
}

:deep(.ant-tabs-content) {
  height: 100%;
}

:deep(.ant-tabs-nav) {
  margin: 0;
  background: rgba(248, 250, 252, 0.5);
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
}

:deep(.ant-tabs-tab) {
  padding: 0.75rem 1rem;
  font-weight: 500;
}

:deep(.ant-tabs-tab-active .ant-tabs-tab-btn) {
  color: #2563eb;
}

:deep(.ant-tabs-ink-bar) {
  background: linear-gradient(90deg, #3b82f6, #6366f1);
}
</style>
