<script setup lang="ts">
import { ref, nextTick } from 'vue';
import { X, Download, Copy, Check } from 'lucide-vue-next';

const props = defineProps<{
  isScreenshotMode: boolean;
  selectedIds: Set<string>;
  targetElement?: HTMLElement | null;
}>();

const emit = defineEmits<{
  (e: 'update:isScreenshotMode', value: boolean): void;
  (e: 'clearSelection'): void;
}>();

const isGeneratingScreenshot = ref(false);
const showFlash = ref(false);
const showScreenshotPreview = ref(false);
const screenshotPreviewUrl = ref<string | null>(null);
const isCopying = ref(false);

const cancelScreenshotMode = () => {
  emit('update:isScreenshotMode', false);
  emit('clearSelection');
};

const generateScreenshot = async () => {
  if (props.selectedIds.size === 0) {
    alert('请至少选择一条消息');
    return;
  }
  
  isGeneratingScreenshot.value = true;
  
  try {
    // Wait for UI to update if needed
    await nextTick();
    
    const messageListEl = props.targetElement;
    
    if (!messageListEl) {
      console.error('Message list element not provided or not found');
      throw new Error('Message list not found');
    }

    // Create a temporary container for the screenshot
    const container = document.createElement('div');
    container.style.position = 'absolute';
    container.style.left = '-9999px';
    container.style.top = '0'; // Reset to 0 to ensure full capture
    container.style.width = '1000px'; // Increased width for better layout
    container.style.backgroundColor = document.documentElement.classList.contains('dark') ? '#1f2937' : '#ffffff';
    container.style.padding = '2rem';
    container.style.height = 'auto';
    container.style.overflow = 'visible';
    document.body.appendChild(container);

    // Clone selected messages
    const selectedElements = Array.from(messageListEl.querySelectorAll('.message-item'))
      .filter(el => props.selectedIds.has(el.getAttribute('data-message-id') || ''));
      
    // Sort elements by their position in DOM to maintain order
    selectedElements.sort((a, b) => {
      return (a as HTMLElement).offsetTop - (b as HTMLElement).offsetTop;
    });

    for (const el of selectedElements) {
      const clone = el.cloneNode(true) as HTMLElement;
      // Remove any selection UI from clone if present
      const selectionUI = clone.querySelector('.selection-checkbox');
      if (selectionUI) selectionUI.remove();
      
      // Ensure clone is visible and has proper styling
      clone.style.marginBottom = '1rem';
      clone.style.height = 'auto';
      clone.style.overflow = 'visible';
      
      // Fix User Message Alignment (ensure they stay on the right)
      if (clone.classList.contains('flex-row-reverse')) {
        clone.style.display = 'flex';
        clone.style.flexDirection = 'row-reverse';
        clone.style.justifyContent = 'flex-start'; // Changed from flex-end to flex-start because of row-reverse
      }
      
      // Force text wrapping and visibility on all children
      const allChildren = clone.querySelectorAll('*');
      allChildren.forEach((child) => {
        const el = child as HTMLElement;
        el.style.height = 'auto';
        el.style.overflow = 'visible';
        
        // Fix User Message Bubble
        if (el.classList.contains('user-message-bubble')) {
             el.style.whiteSpace = 'pre-wrap';
             el.style.color = '#ffffff'; // Force white color
             el.style.backgroundColor = '#3b82f6'; // Primary blue
             el.style.display = 'inline-block';
             el.style.width = 'auto';
             el.style.maxWidth = '100%';
             el.style.borderRadius = '1rem';
             el.style.borderTopRightRadius = '0.125rem';
             el.style.padding = '0.5rem 1rem';
             
             // Ensure children also have white color
             const children = el.querySelectorAll('*');
             children.forEach(child => {
               (child as HTMLElement).style.color = '#ffffff';
             });
        }
        
        // Fix Avatar Stretching
        if (el.classList.contains('w-8') && el.classList.contains('h-8')) {
            el.style.width = '32px';
            el.style.height = '32px';
            el.style.flexShrink = '0';
        }
        
        // Fix Paragraphs and Divs text wrapping
        if (el.tagName === 'P' || el.tagName === 'DIV') {
             if (!el.classList.contains('user-message-bubble')) {
                 el.style.whiteSpace = 'pre-wrap';
             }
        }
        
        // Fix Image max width and remove overflow: visible for images
        if (el.tagName === 'IMG') {
            el.style.maxWidth = '100%';
            el.style.height = 'auto';
            el.style.overflow = 'hidden'; // Fix warning
        }

        // Fix Markdown Content (Assistant)
        if (el.classList.contains('prose') || el.closest('.prose')) {
           el.style.color = document.documentElement.classList.contains('dark') ? '#e5e7eb' : '#374151';
           
           // Fix Code Blocks
           if (el.tagName === 'PRE') {
             el.style.backgroundColor = '#1e1e1e'; // Dark background for code blocks
             el.style.color = '#d4d4d4';
             el.style.padding = '1rem';
             el.style.borderRadius = '0.5rem';
             el.style.overflowX = 'hidden'; // Avoid scrollbars in screenshot
             el.style.whiteSpace = 'pre-wrap'; // Wrap code lines
           }
           
           // Fix Mermaid Diagrams
           if (el.classList.contains('mermaid')) {
             el.style.backgroundColor = 'white'; // Ensure background for visibility
             if (document.documentElement.classList.contains('dark')) {
                el.style.backgroundColor = '#1e1e1e';
             }
           }
        }
      });

      container.appendChild(clone);
    }
    
    // Add a footer/watermark
    const footer = document.createElement('div');
    footer.style.textAlign = 'center';
    footer.style.marginTop = '2rem';
    footer.style.color = '#9ca3af';
    footer.style.fontSize = '0.875rem';
    footer.innerText = 'Generated by LiteChat';
    container.appendChild(footer);

    const html2canvasModule = await import('html2canvas');
    const html2canvas = html2canvasModule.default || html2canvasModule;
    
    // Use a slight delay to ensure rendering is complete
    await new Promise(resolve => setTimeout(resolve, 100));
    
    const canvas = await html2canvas(container, {
      useCORS: true,
      backgroundColor: null,
      scale: 2,
      logging: false,
      windowHeight: container.scrollHeight + 100, // Add some buffer
      height: container.scrollHeight + 50,
      x: 0,
      y: 0, // Capture from top of container
      scrollX: 0,
      scrollY: 0
    });
    
    screenshotPreviewUrl.value = canvas.toDataURL('image/png');
    showScreenshotPreview.value = true;
    
    document.body.removeChild(container);
    
    // Exit mode logic is handled by the parent or after user action in preview
    cancelScreenshotMode();
    
  } catch (err) {
    console.error('Screenshot failed:', err);
    alert('生成截图失败，请重试');
  } finally {
    isGeneratingScreenshot.value = false;
  }
};

const downloadScreenshot = () => {
  if (!screenshotPreviewUrl.value) return;
  const link = document.createElement('a');
  link.download = `chat-selection-${Date.now()}.png`;
  link.href = screenshotPreviewUrl.value;
  link.click();
  showScreenshotPreview.value = false;
  // Show flash effect after successful action
  triggerFlash();
};

const copyScreenshot = async () => {
  if (!screenshotPreviewUrl.value) return;
  
  try {
    isCopying.value = true;
    const response = await fetch(screenshotPreviewUrl.value);
    const blob = await response.blob();
    await navigator.clipboard.write([
      new ClipboardItem({
        [blob.type]: blob
      })
    ]);
    
    // Show success feedback
    setTimeout(() => {
      isCopying.value = false;
      showScreenshotPreview.value = false;
      triggerFlash();
    }, 1000);
  } catch (err) {
    console.error('Failed to copy image:', err);
    isCopying.value = false;
    alert('复制失败，请重试');
  }
};

const triggerFlash = () => {
  showFlash.value = true;
  setTimeout(() => {
    showFlash.value = false;
  }, 300);
};
</script>

<template>
  <div>
    <!-- Screenshot Mode Floating Bar -->
    <div 
      v-if="isScreenshotMode"
      class="fixed bottom-32 left-0 right-0 mx-auto z-50 flex items-center gap-3 md:gap-4 bg-white dark:bg-[#2a2a2a] px-4 md:px-6 py-2.5 md:py-3 rounded-full shadow-xl border border-gray-200 dark:border-gray-700 animate-fade-in-up w-fit max-w-[90%] justify-between"
    >
      <span class="text-sm font-medium text-gray-700 dark:text-gray-200 whitespace-nowrap">
        <span class="md:hidden">已选 {{ selectedIds.size }}</span>
        <span class="hidden md:inline">已选择 {{ selectedIds.size }} 条消息</span>
      </span>
      <div class="h-4 w-px bg-gray-200 dark:bg-gray-700"></div>
      <div class="flex items-center gap-2">
        <button 
          @click="cancelScreenshotMode"
          class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors p-1.5 md:p-0 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 md:hover:bg-transparent md:dark:hover:bg-transparent"
          title="取消"
        >
          <X class="w-5 h-5 md:hidden" />
          <span class="hidden md:inline">取消</span>
        </button>
        <button 
          @click="generateScreenshot"
          class="bg-primary text-white text-sm font-medium rounded-full hover:bg-blue-600 transition-colors shadow-sm flex items-center justify-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed p-1.5 md:px-4 md:py-1.5"
          :disabled="isGeneratingScreenshot"
          title="生成长截图"
        >
          <span v-if="isGeneratingScreenshot" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <template v-else>
            <Check class="w-5 h-5 md:hidden" />
            <span class="hidden md:inline">生成长截图</span>
          </template>
        </button>
      </div>
    </div>

    <!-- Screenshot Flash Effect -->
    <div 
      v-if="showFlash"
      class="fixed inset-0 bg-white z-[100] pointer-events-none animate-flash"
    ></div>

    <!-- Screenshot Preview Modal -->
    <div v-if="showScreenshotPreview" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in p-4">
      <div class="bg-white dark:bg-[#2a2a2a] rounded-2xl shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden animate-scale-in">
        <div class="flex items-center justify-between p-4 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">截屏预览</h3>
          <button 
            @click="showScreenshotPreview = false"
            class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
        
        <div class="flex-1 overflow-auto p-8 bg-gray-50 dark:bg-[#1a1a1a] flex items-center justify-center">
          <img :src="screenshotPreviewUrl || ''" alt="Screenshot Preview" class="max-w-full shadow-lg rounded-lg" />
        </div>
        
        <div class="p-4 border-t border-gray-100 dark:border-gray-700 flex justify-end gap-3 bg-white dark:bg-[#2a2a2a]">
          <button 
            @click="showScreenshotPreview = false"
            class="px-4 py-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
          >
            取消
          </button>
          <button 
            @click="copyScreenshot"
            class="px-4 py-2 flex items-center gap-2 border border-gray-200 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg transition-colors"
            :disabled="isCopying"
          >
            <Check v-if="isCopying" class="w-4 h-4 text-green-500" />
            <Copy v-else class="w-4 h-4" />
            {{ isCopying ? '已复制' : '复制图片' }}
          </button>
          <button 
            @click="downloadScreenshot"
            class="px-4 py-2 flex items-center gap-2 bg-primary text-white hover:bg-blue-600 rounded-lg transition-colors shadow-sm"
          >
            <Download class="w-4 h-4" />
            下载图片
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.3s ease-out;
}

.animate-flash {
  animation: flash 0.3s ease-out forwards;
}

.animate-scale-in {
  animation: scaleIn 0.2s ease-out;
}

.animate-fade-in {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes flash {
  0% { opacity: 0.8; }
  100% { opacity: 0; }
}
</style>
