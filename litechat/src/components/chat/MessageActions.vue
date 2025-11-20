<script setup lang="ts">
import { ref } from 'vue';
import { Copy, Quote, RotateCcw, ThumbsUp, ThumbsDown, Share2 } from 'lucide-vue-next';

import { useAuthStore } from '../../stores/auth';

const props = defineProps<{
  messageContent: string;
  messageId: string;
}>();

const emit = defineEmits<{
  (e: 'quote', messageId: string, content: string): void;
  (e: 'regenerate', id: string): void;
}>();

const authStore = useAuthStore();
const liked = ref(false);
const disliked = ref(false);
const copied = ref(false);

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(props.messageContent);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (err) {
    console.error('Failed to copy:', err);
  }
};

const handleQuote = () => {
  if (!authStore.checkAuth()) return;
  emit('quote', props.messageId, props.messageContent);
};

const handleRegenerate = () => {
  if (!authStore.checkAuth()) return;
  emit('regenerate', props.messageId);
};

const handleLike = () => {
  liked.value = !liked.value;
  if (liked.value) disliked.value = false;
};

const handleDislike = () => {
  disliked.value = !disliked.value;
  if (disliked.value) liked.value = false;
};

const handleShare = async () => {
  if (navigator.share) {
    try {
      await navigator.share({
        title: 'AI Chat',
        text: props.messageContent,
      });
    } catch (err) {
      // User cancelled or error
      copyToClipboard();
    }
  } else {
    copyToClipboard();
  }
};
</script>

<template>
  <div class="flex items-center gap-1 mt-2 text-gray-500">
    <!-- Copy -->
    <button
      @click="copyToClipboard"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      :title="copied ? '已复制' : '复制'"
    >
      <Copy class="w-4 h-4" :class="{ 'text-green-500': copied }" />
    </button>

    <!-- Quote -->
    <button
      @click="handleQuote"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      title="引用"
    >
      <Quote class="w-4 h-4" />
    </button>

    <!-- Regenerate -->
    <button
      @click="handleRegenerate"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      title="重新生成"
    >
      <RotateCcw class="w-4 h-4" />
    </button>

    <!-- Like -->
    <button
      @click="handleLike"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      :class="{ 'text-blue-500': liked }"
      title="点赞"
    >
      <ThumbsUp class="w-4 h-4" :class="{ 'fill-current': liked }" />
    </button>

    <!-- Dislike -->
    <button
      @click="handleDislike"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      :class="{ 'text-red-500': disliked }"
      title="不满意"
    >
      <ThumbsDown class="w-4 h-4" :class="{ 'fill-current': disliked }" />
    </button>

    <!-- Share -->
    <button
      @click="handleShare"
      class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
      title="分享"
    >
      <Share2 class="w-4 h-4" />
    </button>

    <!-- Copied indicator -->
    <span v-if="copied" class="text-xs text-green-500 ml-2">已复制</span>
  </div>
</template>
