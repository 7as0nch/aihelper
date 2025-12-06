<script setup lang="ts">
import { ref } from 'vue';
import { ArrowUp, Square } from 'lucide-vue-next';
import VoiceInput from '../../common/VoiceInput.vue';

defineProps<{
  modelValue: string;
  isLoading: boolean;
  isRecording: boolean;
  hasAttachments: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'send'): void;
  (e: 'stop'): void;
  (e: 'keydown', event: KeyboardEvent): void;
  (e: 'paste', event: ClipboardEvent): void;
  (e: 'focus', event: FocusEvent): void;
  (e: 'resize', event: Event): void;
  (e: 'voice-start'): void;
  (e: 'voice-stop'): void;
  (e: 'voice-result', text: string): void;
  (e: 'voice-interim', text: string): void;
  (e: 'voice-error', message: string): void;
}>();

const textareaRef = ref<HTMLTextAreaElement | null>(null);

const handleInput = (e: Event) => {
  const target = e.target as HTMLTextAreaElement;
  emit('update:modelValue', target.value);
  emit('resize', e);
};

// Expose focus method
const focus = () => {
  textareaRef.value?.focus();
};

const setSelectionRange = (start: number, end: number) => {
  textareaRef.value?.setSelectionRange(start, end);
};

defineExpose({
  focus,
  setSelectionRange,
  textarea: textareaRef
});
</script>

<template>
  <div class="flex items-end gap-2 p-2 md:block md:p-0">
    <textarea
      ref="textareaRef"
      :value="modelValue"
      rows="1"
      class="flex-1 w-full bg-transparent border-0 focus:ring-0 focus:outline-none resize-none py-3 px-3 md:py-4 md:pl-4 md:pr-12 max-h-48 text-gray-900 dark:text-gray-100 placeholder-gray-400 overflow-y-hidden"
      placeholder="输入你的问题..."
      @input="handleInput"
      @keydown="emit('keydown', $event)"
      @paste="emit('paste', $event)"
      @focus="emit('focus', $event)"
    ></textarea>

    <!-- Mobile Right Actions (Voice & Send) -->
    <div class="flex md:hidden items-center gap-2 pb-1">
      <VoiceInput
        size="md"
        @start="emit('voice-start')"
        @stop="emit('voice-stop')"
        @result="emit('voice-result', $event)"
        @interim="emit('voice-interim', $event)"
        @error="emit('voice-error', $event)"
      />

      <button 
        class="p-2 rounded-full transition-all duration-200 ease-in-out"
        :class="[
          (modelValue.trim() || hasAttachments) || isLoading
            ? 'bg-primary text-white hover:bg-blue-600 shadow-md' 
            : 'bg-gray-100 text-gray-400 dark:bg-gray-700 dark:text-gray-500 cursor-not-allowed',
          isLoading ? 'bg-red-500 hover:bg-red-600' : ''
        ]"
        :disabled="(!modelValue.trim() && !hasAttachments) && !isLoading"
        @click="isLoading ? emit('stop') : emit('send')"
        v-tracker="{ type: 'click', name: '发送消息' }"
      >
        <div v-if="isLoading" class="flex items-center justify-center w-5 h-5">
          <Square class="w-3 h-3 fill-current" />
        </div>
        <ArrowUp v-else class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>
