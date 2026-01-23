<script setup lang="ts">
defineProps<{
  visible: boolean;
  url: string;
  appName?: string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm'): void;
}>();

const handleCancel = () => {
  emit('close');
};

const handleConfirm = () => {
  emit('confirm');
};
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="visible" class="fixed inset-0 z-[10000] flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="handleCancel"></div>
        
        <!-- Modal -->
        <div 
          class="relative w-full max-w-sm bg-white dark:bg-[#1e1e1e] rounded-3xl shadow-2xl overflow-hidden animate-in zoom-in-95 duration-200"
        >
          <div class="p-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-gray-100 mb-3">
              你即将离开 {{ appName || 'AI Chat' }}，跳转到第三方网站
            </h3>
            
            <p class="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-4">
              {{ appName || 'AI Chat' }} 出于为您提供便利的目的向您提供第三方链接，我们不对第三方的内容负责，请您谨慎访问，保护您的信息及财产安全
            </p>

            <div class="mb-6 px-4 py-2 bg-gray-50 dark:bg-gray-800/50 rounded-xl border border-gray-100 dark:border-gray-700/50">
              <p class="text-[11px] text-gray-400 dark:text-gray-500 text-left mb-1">即将前往：</p>
              <p class="text-xs text-blue-500 dark:text-blue-400 break-all text-left font-mono line-clamp-2">
                {{ url }}
              </p>
            </div>
            
            <div class="flex gap-3 mt-2">
              <button 
                @click="handleCancel"
                class="flex-1 px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-700 text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                取消
              </button>
              <button 
                @click="handleConfirm"
                class="flex-1 px-4 py-2.5 rounded-xl bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 text-sm font-medium hover:opacity-90 transition-opacity"
              >
                继续访问
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
