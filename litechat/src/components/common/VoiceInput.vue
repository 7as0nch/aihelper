<script setup lang="ts">
import { computed, watch } from 'vue';
import { Mic, MicOff, Languages } from 'lucide-vue-next';
import { useVoiceRecognition } from '../../composables/useVoiceRecognition';

const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
  showLanguageSwitch?: boolean;
}>(), {
  size: 'md',
  showLabel: false,
  showLanguageSwitch: false,
});

const emit = defineEmits<{
  (e: 'start'): void;
  (e: 'stop'): void;
  (e: 'result', text: string): void;
  (e: 'interim', text: string): void;
  (e: 'error', message: string): void;
}>();

const {
  isSupported,
  isRecording,
  transcript,
  interimTranscript,
  error,
  start,
  stop,
  setLanguage,
  supportedLanguages,
} = useVoiceRecognition({
  lang: 'zh-CN',
  continuous: true,
  interimResults: true,
});

// Watch for final transcripts
watch(transcript, (newValue) => {
  if (newValue) {
    emit('result', newValue);
  }
});

// Watch for interim transcripts
watch(interimTranscript, (newValue) => {
  if (newValue) {
    emit('interim', newValue);
  }
});

// Watch for errors
watch(error, (newValue) => {
  if (newValue) {
    emit('error', newValue);
  }
});

const handleClick = () => {
  if (!isSupported) {
    emit('error', '您的浏览器不支持语音识别。请使用 Chrome、Edge 或 Safari 浏览器。');
    return;
  }

  if (isRecording.value) {
    stop();
    emit('stop');
  } else {
    start();
    emit('start');
  }
};

const handleLanguageChange = (event: Event) => {
  const select = event.target as HTMLSelectElement;
  setLanguage(select.value);
};

// Size classes
const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'p-1.5';
    case 'lg':
      return 'p-3';
    default:
      return 'p-2';
  }
});

const iconSizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'w-4 h-4';
    case 'lg':
      return 'w-6 h-6';
    default:
      return 'w-5 h-5';
  }
});
</script>

<template>
  <div class="voice-input-wrapper flex items-center gap-1">
    <!-- Language Selector -->
    <div v-if="showLanguageSwitch && isSupported" class="relative">
      <select
        class="appearance-none bg-transparent text-xs text-gray-500 dark:text-gray-400 pr-5 py-1 cursor-pointer focus:outline-none"
        @change="handleLanguageChange"
      >
        <option 
          v-for="lang in supportedLanguages" 
          :key="lang.code" 
          :value="lang.code"
        >
          {{ lang.label }}
        </option>
      </select>
      <Languages class="w-3 h-3 absolute right-0 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" />
    </div>

    <!-- Voice Button -->
    <button 
      @click="handleClick"
      class="relative rounded-full transition-all duration-200 ease-in-out"
      :class="[
        sizeClasses,
        isRecording 
          ? 'text-red-500 bg-red-50 dark:bg-red-900/20 hover:bg-red-100 dark:hover:bg-red-900/30' 
          : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700',
        !isSupported && 'opacity-50 cursor-not-allowed'
      ]"
      :title="!isSupported ? '浏览器不支持语音输入' : (isRecording ? '停止录音' : '语音输入')"
      :disabled="!isSupported"
    >
      <!-- Recording Animation -->
      <span 
        v-if="isRecording" 
        class="absolute inset-0 rounded-full animate-ping bg-red-400/30"
      />
      
      <!-- Icon -->
      <Mic v-if="!isRecording" :class="iconSizeClasses" />
      <MicOff v-else :class="iconSizeClasses" />
      
      <!-- Label -->
      <span v-if="showLabel" class="ml-1 text-sm">
        {{ isRecording ? '停止' : '语音' }}
      </span>
    </button>

    <!-- Recording Indicator -->
    <span 
      v-if="isRecording" 
      class="text-xs text-red-500 animate-pulse"
    >
      录音中...
    </span>
  </div>
</template>

<style scoped>
.voice-input-wrapper {
  position: relative;
}

@keyframes ping {
  75%, 100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

.animate-ping {
  animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
}
</style>
