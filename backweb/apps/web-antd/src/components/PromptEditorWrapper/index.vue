<script setup lang="ts">
import { ref, watch } from 'vue';
import { Tabs } from 'ant-design-vue';
import PromptEditor from '#/components/PromptEditor/index.vue';

const props = defineProps<{
  systemPrompt?: string;
  userInputPrompt?: string;
}>();

const emit = defineEmits<{
  'update:systemPrompt': [value: string];
  'update:userInputPrompt': [value: string];
}>();

const activeTab = ref('system');
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
  <div class="h-full flex flex-col bg-white rounded-lg border border-gray-200 overflow-hidden">
    <Tabs v-model:activeKey="activeTab" class="flex-1 flex flex-col">
      <Tabs.TabPane key="system" tab="系统提示词">
        <div class="h-full">
          <PromptEditor
            :model-value="systemPromptValue"
            @update:model-value="handleSystemPromptChange"
            class="h-full"
          />
        </div>
      </Tabs.TabPane>
      <Tabs.TabPane key="user" tab="用户输入提示词">
        <div class="h-full">
          <PromptEditor
            :model-value="userInputPromptValue"
            @update:model-value="handleUserInputPromptChange"
            class="h-full"
          />
        </div>
      </Tabs.TabPane>
    </Tabs>
  </div>
</template>

<style scoped>
:deep(.ant-tabs) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

:deep(.ant-tabs-content-holder) {
  flex: 1;
  overflow: hidden;
}

:deep(.ant-tabs-tabpane) {
  height: 100%;
  overflow: hidden;
}
</style>
