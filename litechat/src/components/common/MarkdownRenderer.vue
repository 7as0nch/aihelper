<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { MdPreview, config } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';
import { useDark } from '@vueuse/core';
import mermaid from 'mermaid';
import ExternalLinkModal from './ExternalLinkModal.vue';

import 'highlight.js/styles/github-dark.css';

// Props
const props = defineProps<{
  content: string;
  loading?: boolean;
  quoteSearchLinks?: any[];
}>();

// Dark mode detection
const isDark = useDark();
const theme = computed(() => isDark.value ? 'dark' : 'light');

// MD Editor ID for isolation
const editorId = 'litechat-preview';

// Citation Tooltip State
const hoverCitation = ref<any>(null);
const tooltipPosition = ref({ top: 0, left: 0 });
const tooltipVisible = ref(false);
let hoverTimeout: any = null;

// External Link Modal State
const linkModalVisible = ref(false);
const pendingUrl = ref('');

const openExternalLink = (url: string) => {
  pendingUrl.value = url;
  linkModalVisible.value = true;
  tooltipVisible.value = false;
};

const confirmLinkJump = () => {
  if (pendingUrl.value) {
    window.open(pendingUrl.value, '_blank', 'noopener noreferrer');
  }
  linkModalVisible.value = false;
  pendingUrl.value = '';
};

// Interaction Handler
const handleMouseMove = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  const citationSup = target.closest('.citation-sup');
  
  if (citationSup) {
    const index = parseInt(citationSup.getAttribute('data-citation-index') || '0');
    const link = props.quoteSearchLinks?.[index - 1];
    
    if (link) {
      const rect = citationSup.getBoundingClientRect();
      const tooltipWidth = window.innerWidth < 640 ? window.innerWidth - 32 : 288;
      const padding = 16;
      
      let left = rect.left + (rect.width / 2);
      
      // Horizontal boundary detection
      if (left - tooltipWidth / 2 < padding) {
        left = tooltipWidth / 2 + padding;
      } else if (left + tooltipWidth / 2 > window.innerWidth - padding) {
        left = window.innerWidth - tooltipWidth / 2 - padding;
      }
      
      // Vertical check: if there's no space above, show below
      const showBelow = rect.top < 160; 
      
      tooltipPosition.value = {
        top: showBelow ? rect.bottom + 8 : rect.top - 8,
        left: left
      };
      
      hoverCitation.value = { ...link, index, showBelow };
      tooltipVisible.value = true;
      
      if (hoverTimeout) {
        clearTimeout(hoverTimeout);
        hoverTimeout = null;
      }
    }
  } else if (!target.closest('.citation-tooltip')) {
    if (!hoverTimeout && tooltipVisible.value) {
      hoverTimeout = setTimeout(() => {
        tooltipVisible.value = false;
        hoverTimeout = null;
      }, 300);
    }
  }
};

// Hide tooltip on scroll to prevent it from detaching from the element
const handleScroll = () => {
  if (tooltipVisible.value) {
    tooltipVisible.value = false;
    if (hoverTimeout) {
      clearTimeout(hoverTimeout);
      hoverTimeout = null;
    }
  }
};

onMounted(() => {
    // Watch for content changes could be handled by `MdPreview` internally rereading modelValue,
    // but our custom mermaid renderer needs to trigger on updates.
    // However, `MdPreview` rerender might replace DOM.
    // Let's rely on updated hook or simply re-run renderMermaid when content changes watching props.
    window.addEventListener('scroll', handleScroll, true); // Use capture to catch internal container scrolls
});

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll, true);
});

const handleMouseLeave = () => {
  if (!hoverTimeout) {
    hoverTimeout = setTimeout(() => {
      tooltipVisible.value = false;
      hoverTimeout = null;
    }, 300);
  }
};

const handleTooltipEnter = () => {
  if (hoverTimeout) {
    clearTimeout(hoverTimeout);
    hoverTimeout = null;
  }
};

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

    // 3. Handle [quote:idx] citations
    md.inline.ruler.after('emphasis', 'quote', (state: any, silent: boolean) => {
      const regex = /^\[quote:(\d+)\]/;
      const str = state.src.slice(state.pos);
      const match = str.match(regex);
      if (!match) return false;

      if (!silent) {
        const token = state.push('quote', '', 0);
        token.meta = { index: parseInt(match[1]) + 1 };
      }

      state.pos += match[0].length;
      return true;
    });

    md.renderer.rules.quote = (tokens: any[], idx: number) => {
      const index = tokens[idx].meta.index;
      return `<sup class="citation-sup inline-flex items-center justify-center w-4 h-4 text-[10px] font-bold text-gray-400 dark:text-gray-500 bg-gray-100 dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 mx-0.5 select-none cursor-pointer hover:bg-primary hover:text-white transition-colors" title="点击查看引用 [${index}]" data-citation-index="${index}">${index}</sup>`;
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

  // 1. Handle Citation Clicks
  const citationSup = target.closest('.citation-sup');
  if (citationSup) {
    const index = parseInt(citationSup.getAttribute('data-citation-index') || '0');
    const link = props.quoteSearchLinks?.[index - 1];
    if (link?.url) {
      openExternalLink(link.url);
    }
    return;
  }

  // Handle generic links in markdown
  const anchor = target.closest('a');
  if (anchor && anchor.href && !anchor.href.startsWith(window.location.origin) && !anchor.classList.contains('export-table-btn') && !anchor.closest('.md-editor-toolbar')) {
    // Only intercept if it's a real external link and not an internal action
    if (anchor.target === '_blank' || anchor.hostname !== window.location.hostname) {
      event.preventDefault();
      openExternalLink(anchor.href);
      return;
    }
  }

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
  <div 
    ref="rootRef" 
    class="relative" 
    @click="handleClick"
    @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave"
  >
    <MdPreview 
      :editorId="editorId" 
      :modelValue="props.content" 
      :theme="theme"
      :codeFoldable="true"
      previewTheme="default" 
    />

    <!-- External Link Confirmation Modal -->
    <ExternalLinkModal 
      :visible="linkModalVisible" 
      :url="pendingUrl"
      appName="LiteChat"
      @close="linkModalVisible = false"
      @confirm="confirmLinkJump"
    />

    <!-- Citation Hover Tooltip -->
    <Teleport to="body">
      <div 
        v-if="tooltipVisible && hoverCitation"
        class="citation-tooltip fixed z-[9999] w-[calc(100vw-32px)] sm:w-72 animate-in fade-in zoom-in duration-200"
        :style="{ 
          top: tooltipPosition.top + 'px', 
          left: tooltipPosition.left + 'px',
          transform: hoverCitation.showBelow ? 'translateX(-50%)' : 'translate(-50%, -100%)' 
        }"
        @mouseenter="handleTooltipEnter"
        @mouseleave="handleMouseLeave"
      >
        <div class="bg-white dark:bg-[#1e1e1e] rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700/50 p-3 overflow-hidden">
          <div class="flex items-center gap-2 mb-2 pb-2 border-b border-gray-100 dark:border-gray-800">
            <span class="shrink-0 text-[10px] font-bold text-gray-400 dark:text-gray-500 w-4 h-4 flex items-center justify-center rounded bg-gray-100 dark:bg-gray-700">
              {{ hoverCitation.index }}
            </span>
            <div class="text-xs font-bold text-gray-800 dark:text-gray-200 truncate pr-4">
              {{ hoverCitation.title }}
            </div>
          </div>
          <div class="text-[11px] text-gray-500 dark:text-gray-400 line-clamp-3 leading-relaxed italic mb-2">
            {{ hoverCitation.content }}
          </div>
          <div class="flex items-center justify-between mt-2 pt-2 border-t border-gray-50 dark:border-gray-800/50">
            <div class="text-[9px] text-gray-400 truncate max-w-[180px]">
              {{ hoverCitation.url }}
            </div>
            <div 
              @click="openExternalLink(hoverCitation.url)"
              class="text-[10px] text-primary hover:underline font-medium flex items-center gap-0.5 cursor-pointer"
            >
              访问来源
              <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-external-link"><path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/></svg>
            </div>
          </div>
        </div>
        <!-- Arrow -->
        <div 
          v-if="!hoverCitation.showBelow"
          class="absolute -bottom-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-white dark:bg-[#1e1e1e] border-r border-b border-gray-200 dark:border-gray-700/50 rotate-45"
        ></div>
        <div 
          v-else
          class="absolute -top-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-white dark:bg-[#1e1e1e] border-l border-t border-gray-200 dark:border-gray-700/50 rotate-45"
        ></div>
      </div>
    </Teleport>
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

/* Citation styles */
.citation-sup {
  vertical-align: baseline;
  position: relative;
  top: -6px;
  line-height: 0;
  text-decoration: none !important;
}

.citation-sup:hover {
  border-color: var(--primary-color, #3b82f6);
}
</style>
