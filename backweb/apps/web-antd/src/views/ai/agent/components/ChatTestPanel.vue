<script setup lang="ts">
import { ref, h, computed } from 'vue';
import { Input, Button } from 'ant-design-vue';
import { SendOutlined } from '@ant-design/icons-vue';

const props = withDefaults(
  defineProps<{
    variant?: 'card' | 'plain';
    showHeader?: boolean;
  }>(),
  {
    variant: 'card',
    showHeader: true,
  },
);

const containerClass = computed(() =>
  props.variant === 'card'
    ? 'h-full flex flex-col bg-white rounded-xl border border-slate-200/60 shadow-lg shadow-slate-200/50 overflow-hidden'
    : 'h-full flex flex-col',
);

const messages = ref<Array<{ role: 'user' | 'assistant'; content: string }>>([]);
const inputValue = ref('');
const loading = ref(false);

const handleSend = async () => {
  if (!inputValue.value.trim()) {
    return;
  }

  const userMessage = inputValue.value.trim();
  messages.value.push({ role: 'user', content: userMessage });
  inputValue.value = '';
  loading.value = true;

  // TODO: 调用实际的 Agent API 进行测试
  // 这里暂时模拟响应
  setTimeout(() => {
    messages.value.push({
      role: 'assistant',
      content: '这是一个测试响应。实际功能需要连接后端 Agent API。',
    });
    loading.value = false;
  }, 1000);
};

const handleKeyPress = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    handleSend();
  }
};
</script>

<template>
  <div :class="containerClass">
    <div v-if="props.showHeader" class="px-4 py-3 border-b border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-sm font-medium text-gray-700">Agent 预览测试</h3>
          <p class="text-xs text-gray-500 mt-1">测试当前 Agent 的对话能力</p>
        </div>
        <span class="text-xs text-gray-400">仅本地模拟</span>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto p-4 space-y-4 bg-gradient-to-b from-slate-50/30 to-white">
      <div
        v-if="messages.length === 0"
        class="flex items-center justify-center h-full text-gray-400 text-sm"
      >
        暂无对话记录，开始测试吧
      </div>
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="[
          'flex',
          msg.role === 'user' ? 'justify-end' : 'justify-start',
        ]"
      >
        <div
          :class="[
            'max-w-[85%] rounded-2xl px-4 py-2.5 shadow-sm',
            msg.role === 'user'
              ? 'bg-gradient-to-br from-blue-500 to-blue-600 text-white'
              : 'bg-white text-slate-800 border border-slate-200',
          ]"
        >
          <div class="text-sm whitespace-pre-wrap break-words leading-relaxed">
            {{ msg.content }}
          </div>
        </div>
      </div>
      <div v-if="loading" class="flex justify-start">
        <div class="bg-white rounded-2xl px-4 py-2 shadow-sm border border-slate-200">
          <div class="text-sm text-slate-500">思考中...</div>
        </div>
      </div>
    </div>

    <div class="px-4 py-4 border-t border-slate-100 bg-white">
      <div class="flex gap-2">
        <Input
          :value="inputValue"
          placeholder="输入消息进行测试..."
          @keypress="handleKeyPress"
          @update:value="(val) => (inputValue = val)"
          :disabled="loading"
          class="flex-1"
        />
        <Button
          type="primary"
          :icon="h(SendOutlined)"
          :loading="loading"
          @click="handleSend"
          :disabled="!inputValue.trim()"
          class="shadow-lg shadow-blue-500/25"
        >
          发送
        </Button>
      </div>
    </div>
  </div>
</template>
