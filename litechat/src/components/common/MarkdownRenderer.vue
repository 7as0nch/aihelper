<script setup lang="ts">
import { computed, ref, onMounted, nextTick, watch } from 'vue';
import { MdPreview, config } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';
import { useDark } from '@vueuse/core';
import mermaid from 'mermaid';

import 'highlight.js/styles/github-dark.css';

// Props
const props = defineProps<{
  content: string;
  loading?: boolean;
}>();

// Dark mode detection
const isDark = useDark();
const theme = computed(() => isDark.value ? 'dark' : 'light');

// MD Editor ID for isolation
const editorId = 'litechat-preview';

// Initialize Mermaid
mermaid.initialize({
  startOnLoad: false,
  theme: 'default',
  securityLevel: 'loose',
});

// Custom Config to restore features
config({
  markdownItConfig(md: any) {
    // 1. Restore Table Export
    md.renderer.rules.table_open = () => '<div class="table-wrapper my-4 w-full border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 shadow-sm"><table class="w-full text-left text-sm">';
    md.renderer.rules.table_close = () => `
      </table>
      <div class="flex justify-end p-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
        <button class="export-table-btn flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-table"><path d="M12 3v18"/><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M3 15h18"/></svg>
          Export to Sheets
        </button>
      </div>
    </div>`;

    // 2. Customize Code Blocks & Mermaid
    // We override 'fence' to handle both:
    // - Wrapper for Mermaid to verify it renders and enable image export
    // - Fallback to default for Code Blocks (md-editor-v3 default style)
    const defaultFence = md.renderer.rules.fence;

    md.renderer.rules.fence = (tokens: any[], idx: number, options: any, env: any, self: any) => {
      const token = tokens[idx];
      const info = token.info ? token.info.trim() : '';

      // Handle Mermaid specifically
      if (info === 'mermaid') {
        const uniqueId = `mermaid-${idx}-${Math.random().toString(36).substr(2, 9)}`;
        const encodedCode = encodeURIComponent(token.content);
        return `
          <div class="mermaid-container relative my-4 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
             <div id="${uniqueId}" class="mermaid-view flex justify-center overflow-x-auto" data-code="${encodedCode}">
                ${token.content}
             </div>
             <div class="absolute top-2 right-2 opacity-0 hover:opacity-100 transition-opacity">
                <button class="export-mermaid-btn p-1 bg-white/80 dark:bg-gray-700/80 rounded shadow hover:bg-white dark:hover:bg-gray-600">
                   <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-image-down"><path d="M10.3 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10l-3.1-3.1a2 2 0 0 0-2.814.014L6 21"/><path d="m14 19 3 3v-5.5"/><path d="m14 22 3-3"/><circle cx="9" cy="9" r="2"/></svg>
                </button>
             </div>
          </div>
        `;
      }
      
      // Fallback to default styling for everything else
      return defaultFence(tokens, idx, options, env, self);
    };
  }
});

const rootRef = ref<HTMLElement | null>(null);

// Rerender Mermaid manually because we bypassed standard logic
const renderMermaid = async () => {
    await nextTick();
    if (!rootRef.value) return;
    const mermaidDivs = rootRef.value.querySelectorAll('.mermaid-view');
    
    mermaidDivs.forEach(async (div: any) => {
        if (div.getAttribute('data-processed')) return;
        const content = decodeURIComponent(div.getAttribute('data-code') || '');
        const id = div.id;
        try {
            const { svg } = await mermaid.render(`${id}-svg`, content);
            div.innerHTML = svg;
            div.setAttribute('data-processed', 'true');
        } catch (e) {
            div.innerHTML = `<div class="text-red-500 text-xs p-2">流程图渲染失败</div>`;
            console.error(e);
        }
    });
};

// Interaction Handler
const handleClick = async (event: MouseEvent) => {
  const target = event.target as HTMLElement;

  // 2. Export Table
  const exportTableBtn = target.closest('.export-table-btn');
  if (exportTableBtn) {
    const wrapper = exportTableBtn.closest('.table-wrapper');
    const table = wrapper?.querySelector('table');
    if (table) {
      try {
        const XLSX = await import('xlsx');
        const wb = XLSX.utils.table_to_book(table);
        XLSX.writeFile(wb, `table-export-${Date.now()}.xlsx`);
      } catch (err) {
        console.error('Export failed', err);
      }
    }
    return;
  }

  // 3. Export Mermaid Image
  const exportMermaidBtn = target.closest('.export-mermaid-btn');
  if (exportMermaidBtn) {
    const container = exportMermaidBtn.closest('.mermaid-container');
    const mermaidView = container?.querySelector('.mermaid-view');
    if (mermaidView) {
      try {
        const html2canvas = (await import('html2canvas')).default;
        // Need to temporarily set background to white for transparent SVGs if needed, though usually fine
        const canvas = await html2canvas(mermaidView as HTMLElement, { backgroundColor: '#ffffff' });
        const link = document.createElement('a');
        link.download = `mermaid-${Date.now()}.png`;
        link.href = canvas.toDataURL();
        link.click();
      } catch (err) {
         console.error('Mermaid export failed', err);
      }
    }
  }
};

onMounted(() => {
    // Watch for content changes could be handled by `MdPreview` internally rereading modelValue,
    // but our custom mermaid renderer needs to trigger on updates.
    // However, `MdPreview` rerender might replace DOM.
    // Let's rely on updated hook or simply re-run renderMermaid when content changes watching props.
});

// Since MdPreview is reactive, we need to ensure our manual mermaid rendering runs after format updates.
// We can use a watcher or simply rely on the fact that `MdPreview` recreates DOM.
watch(() => props.content, () => {
    setTimeout(renderMermaid, 100); // Small delay to allow DOM update
});
onMounted(renderMermaid);

</script>

<template>
  <div ref="rootRef" @click="handleClick">
    <MdPreview 
      :editorId="editorId" 
      :modelValue="props.content" 
      :theme="theme"
      :codeFoldable="true"
      previewTheme="default" 
    />
  </div>
</template>

<style scoped>
:deep(.md-editor-preview-wrapper) {
  padding: 0;
}
</style>

<style>

/* Ensure sticky header works */
.md-editor-code-head {
  position: sticky;
  top: 0;
  z-index: 20;
  background-color: #f6f8fa; /* github-light code header color */
}

.dark .md-editor-code-head {
  background-color: #161b22; /* github-dark code header color */
}

/* Ensure no parent container blocks sticky behavior */
.md-editor-code,
.md-editor-preview-wrapper,
.md-editor-preview,
#litechat-preview,
#litechat-preview-preview {
  overflow: visible !important;
}

/* Force dark backgorund for code blocks as requested */
.chat-code-block.mac-style {
  background-color: #1e1e1e !important;
}

/* Mermaid hover effect */
.mermaid-container:hover .absolute {
  opacity: 1;
}
</style>
