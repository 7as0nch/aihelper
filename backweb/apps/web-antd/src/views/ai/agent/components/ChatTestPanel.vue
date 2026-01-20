<script setup lang="ts">
import { ref, h } from 'vue';
import { Input, Button, message } from 'ant-design-vue';
import { SendOutlined } from '@ant-design/icons-vue';

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
  <div class="h-full flex flex-col bg-white rounded-lg border border-gray-200 overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
      <h3 class="text-sm font-medium text-gray-700">Agent 预览测试</h3>
      <p class="text-xs text-gray-500 mt-1">测试当前 Agent 的对话能力</p>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-if="messages.length === 0" class="flex items-center justify-center h-full text-gray-400 text-sm">
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
            'max-w-[80%] rounded-lg px-4 py-2',
            msg.role === 'user'
              ? 'bg-blue-500 text-white'
              : 'bg-gray-100 text-gray-800',
          ]"
        >
          <div class="text-sm whitespace-pre-wrap break-words">
            {{ msg.content }}
          </div>
        </div>
      </div>
      <div v-if="loading" class="flex justify-start">
        <div class="bg-gray-100 rounded-lg px-4 py-2">
          <div class="text-sm text-gray-500">思考中...</div>
        </div>
      </div>
    </div>

    <div class="px-4 py-3 border-t border-gray-200 bg-gray-50">
      <div class="flex gap-2">
        <Input.TextArea
          v-model:value="inputValue"
          :rows="2"
          placeholder="输入消息进行测试..."
          @keypress="handleKeyPress"
          :disabled="loading"
        />
        <Button
          type="primary"
          :icon="h(SendOutlined)"
          :loading="loading"
          @click="handleSend"
          :disabled="!inputValue.trim()"
        >
          发送
        </Button>
      </div>
    </div>
  </div>
</template>
