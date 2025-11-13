<template>
  <div class="quick-prompts-container">
    <div class="quick-prompts-wrapper">
      <div
        v-for="(prompt, index) in prompts"
        :key="index"
        class="quick-prompt-item"
        @click="sendPrompt(prompt)"
      >
        {{ prompt }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// 定义快捷发送语
const prompts = ref([
	'根据渠道ID和时间范围，查询渠道ROI报表',
  '帮我生成前三天的用户报表',
  '帮我生成所有未成年的男生报表',
  '帮我生成最近一年的收益报表'
])

// 发送快捷语的事件
const emit = defineEmits<{
  (e: 'send', prompt: string): void
}>()

// 发送快捷语
const sendPrompt = (prompt: string) => {
  emit('send', prompt)
}
</script>

<style scoped>
.quick-prompts-container {
  width: 100%;
  padding: 10px 0;
  background-color: transparent;
  border-bottom: 1px solid var(--border-color, #e4e4e4);
}

.quick-prompts-wrapper {
  display: flex;
  overflow-x: auto;
  padding: 0 16px;
  gap: 12px;
  /* 隐藏滚动条 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.quick-prompts-wrapper::-webkit-scrollbar {
  display: none; /* Chrome Safari */
}

.quick-prompt-item {
  flex: 0 0 auto;
  padding: 8px 16px;
  background-color: rgba(32, 128, 240, 0.1);
  color: #2080f0;
  border-radius: 16px;
  font-size: 14px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid var(--border-color, #e4e4e4);
}

.quick-prompt-item:hover {
  background-color: rgba(32, 128, 240, 0.2);
  border-color: #2080f0;
  transform: translateY(-1px);
}

/* 暗黑模式适配 */
html.dark .quick-prompt-item {
  background-color: rgba(64, 158, 255, 0.15);
  color: #409eff;
  border-color: #444;
}

html.dark .quick-prompt-item:hover {
  background-color: rgba(64, 158, 255, 0.25);
  border-color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .quick-prompts-wrapper {
    padding: 0 12px;
    gap: 10px;
  }

  .quick-prompt-item {
    padding: 6px 14px;
    font-size: 13px;
  }
}
</style>
