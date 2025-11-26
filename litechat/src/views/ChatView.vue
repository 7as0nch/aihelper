<script setup lang="ts">
import { onMounted, watch, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useChatStore } from '@/stores/chat';
import ChatArea from '@/components/chat/ChatArea.vue';

const route = useRoute();
const chatStore = useChatStore();

const isScreenshotMode = ref(false);
const selectedMessageIds = ref<Set<string>>(new Set());

defineEmits<{
  (e: 'toggleSidebar'): void;
}>();

// Load chat history when route changes
const loadChat = () => {
  const chatId = route.params.id as string;
  if (chatId) {
    chatStore.loadChatHistory(chatId);
  } else {
    chatStore.clearMessages();
  }
};

const handleToggleSelect = (id: string) => {
  if (selectedMessageIds.value.has(id)) {
    selectedMessageIds.value.delete(id);
  } else {
    selectedMessageIds.value.add(id);
  }
};

onMounted(() => {
  loadChat();
});

watch(() => route.params.id, () => {
  loadChat();
});
</script>

<template>
  <ChatArea 
    :is-screenshot-mode="isScreenshotMode"
    :selected-message-ids="selectedMessageIds"
    @update:is-screenshot-mode="isScreenshotMode = $event"
    @toggle-select="handleToggleSelect"
    @toggle-sidebar="$emit('toggleSidebar')"
  />
</template>
