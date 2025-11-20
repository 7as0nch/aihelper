<script setup lang="ts">
import { onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useChatStore } from '@/stores/chat';
import ChatArea from '@/components/chat/ChatArea.vue';

const route = useRoute();
const chatStore = useChatStore();

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

onMounted(() => {
  loadChat();
});

watch(() => route.params.id, () => {
  loadChat();
});
</script>

<template>
  <ChatArea 
    @toggle-sidebar="$emit('toggleSidebar')"
  />
</template>
