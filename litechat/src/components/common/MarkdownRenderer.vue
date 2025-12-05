<script setup lang="ts">
import { onMounted, nextTick, watch, ref } from 'vue';
import Vue3MarkdownIt from 'vue3-markdown-it';
import 'github-markdown-css/github-markdown.css';
import 'highlight.js/styles/github-dark.css'; // Using dark theme as default or based on preference
import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import typescript from 'highlight.js/lib/languages/typescript';
import python from 'highlight.js/lib/languages/python';
import java from 'highlight.js/lib/languages/java';
import go from 'highlight.js/lib/languages/go';
import xml from 'highlight.js/lib/languages/xml';
import css from 'highlight.js/lib/languages/css';
import json from 'highlight.js/lib/languages/json';
import bash from 'highlight.js/lib/languages/bash';
import sql from 'highlight.js/lib/languages/sql';
import markdown from 'highlight.js/lib/languages/markdown';
import yaml from 'highlight.js/lib/languages/yaml';
import shell from 'highlight.js/lib/languages/shell';
import c from 'highlight.js/lib/languages/c';
import cpp from 'highlight.js/lib/languages/cpp';
import csharp from 'highlight.js/lib/languages/csharp';
import rust from 'highlight.js/lib/languages/rust';
import php from 'highlight.js/lib/languages/php';
import ruby from 'highlight.js/lib/languages/ruby';
import swift from 'highlight.js/lib/languages/swift';
import kotlin from 'highlight.js/lib/languages/kotlin';
import diff from 'highlight.js/lib/languages/diff';
import ini from 'highlight.js/lib/languages/ini';
import dockerfile from 'highlight.js/lib/languages/dockerfile';
import graphql from 'highlight.js/lib/languages/graphql';

// Register languages
hljs.registerLanguage('javascript', javascript);
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('python', python);
hljs.registerLanguage('java', java);
hljs.registerLanguage('go', go);
hljs.registerLanguage('xml', xml);
hljs.registerLanguage('html', xml);
hljs.registerLanguage('css', css);
hljs.registerLanguage('json', json);
hljs.registerLanguage('bash', bash);
hljs.registerLanguage('sql', sql);
hljs.registerLanguage('markdown', markdown);
hljs.registerLanguage('yaml', yaml);
hljs.registerLanguage('shell', shell);
hljs.registerLanguage('c', c);
hljs.registerLanguage('cpp', cpp);
hljs.registerLanguage('csharp', csharp);
hljs.registerLanguage('rust', rust);
hljs.registerLanguage('php', php);
hljs.registerLanguage('ruby', ruby);
hljs.registerLanguage('swift', swift);
hljs.registerLanguage('kotlin', kotlin);
hljs.registerLanguage('diff', diff);
hljs.registerLanguage('ini', ini);
hljs.registerLanguage('dockerfile', dockerfile);
hljs.registerLanguage('graphql', graphql);

const props = defineProps<{
  content: string;
  loading?: boolean;
}>();

const rootRef = ref<HTMLElement | null>(null);
const instanceId = Math.random().toString(36).substr(2, 9);

// Cache for Mermaid SVGs to prevent flickering
// Map<CodeContent, SVGString>
const mermaidCache = new Map<string, string>();
// Map<BlockIndex, LastValidSVG> - to show last valid state during streaming error
const lastValidMermaidCache = new Map<string, string>();

// Lazy load Mermaid
let mermaidInstance: any = null;
const initMermaid = async () => {
  if (!mermaidInstance) {
    const mermaidModule = await import('mermaid');
    mermaidInstance = mermaidModule.default;
    mermaidInstance.initialize({
      startOnLoad: false,
      theme: 'default',
      securityLevel: 'loose',
      suppressErrorRendering: true,
    });
    mermaidInstance.parseError = () => {};
  }
  return mermaidInstance;
};

// Custom plugin to inject our rules
const customRulesPlugin = (md: any) => {
  // Custom fence rule for mermaid and code blocks
  md.renderer.rules.fence = (tokens: any[], idx: number, options: any, _env: any, _self: any) => {
    const token = tokens[idx];
    const info = token.info.trim();
    
    if (info === 'mermaid') {
      // Use deterministic ID based on token index
      // This is stable during streaming as long as the block position doesn't change
      const uniqueId = `mermaid-${instanceId}-${idx}`;
      const encodedCode = encodeURIComponent(token.content);
      
      // Check if we have a last valid render for this specific block
      const lastValidSvg = lastValidMermaidCache.get(uniqueId);
      
      // If we have a valid SVG, use it to prevent flicker
      // We wrap it in a div that renderMermaid will find and potentially update
      let innerContent = '';
      let processedAttr = '';
      let hashAttr = '';
      
      if (lastValidSvg) {
        innerContent = lastValidSvg;
        // Optimization: Mark as processed so renderMermaid skips redundant update
        processedAttr = 'data-processed="true"';
        hashAttr = `data-content-hash="${token.content.length}"`;
      } else {
        // Show placeholder if no cache
        innerContent = `
          <div class="mermaid-placeholder flex flex-col items-center justify-center p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg text-gray-500 text-sm animate-pulse border border-gray-100 dark:border-gray-700 min-h-[100px]">
             <svg class="animate-spin mb-2 h-6 w-6 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
               <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
               <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
             </svg>
             <span>流程图生成中...</span>
          </div>
        `;
      }

      return `
        <div id="${uniqueId}" class="mermaid-view my-4" data-code="${encodedCode}" ${processedAttr} ${hashAttr}>
          ${innerContent}
        </div>
      `;
    }
    
    const lang = info ? info.split(/\s+/)[0] : '';
    const languageLabel = lang ? lang.toUpperCase() : 'TEXT';
    const code = options.highlight ? options.highlight(token.content, lang) : md.utils.escapeHtml(token.content);
    
    return `
      <div class="chat-code-block mac-style">
        <div class="chat-code-header">
          <div class="mac-buttons">
            <span class="mac-button red"></span>
            <span class="mac-button yellow"></span>
            <span class="mac-button green"></span>
          </div>
          <div class="code-actions">
            <span class="chat-code-lang">${languageLabel}</span>
            <button class="copy-btn" data-code="${encodeURIComponent(token.content)}">
              <span class="copy-icon">📋</span>
              <span class="copy-text">复制</span>
            </button>
          </div>
        </div>
        <div class="chat-code-content">
          <pre><code class="hljs language-${lang}">${code}</code></pre>
        </div>
      </div>
    `;
  };

  // Custom table rules
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
};

const plugins = [
  {
    plugin: customRulesPlugin
  }
];

const mdOptions = {
  html: false,
  linkify: true,
  typographer: true,
  highlight: {
    highlighter: (str: string, lang: string): string => {
      if (lang && hljs.getLanguage(lang)) {
        try {
          return hljs.highlight(str, { language: lang, ignoreIllegals: true }).value;
        } catch (__) {}
      }
      return ''; 
    }
  }
};

const renderMermaid = async () => {
  await nextTick();
  if (!rootRef.value) return;
  
  const mermaidDivs = rootRef.value.querySelectorAll('.mermaid-view');
  
  if (mermaidDivs.length === 0) return;
  
  const mermaid = await initMermaid();
  
  mermaidDivs.forEach(async (div) => {
    // ID is already set by fence rule
    const uniqueId = div.id; 
    const content = decodeURIComponent(div.getAttribute('data-code') || '');
    
    // Check cache first (exact match)
    if (mermaidCache.has(content)) {
      if (div.getAttribute('data-processed') === 'true') return;
      
      div.innerHTML = mermaidCache.get(content)!;
      div.setAttribute('data-processed', 'true');
      div.setAttribute('data-content-hash', String(content.length));
      return;
    }

    // If already processed with this exact content hash, skip
    if (div.getAttribute('data-processed') === 'true' && div.getAttribute('data-content-hash') === String(content.length)) {
      return;
    }

    // Improved sanitize function implementation
    const improvedSanitize = (code: string) => {
      let s = code;
      
      s = s.replace(/:::end\b/g, ':::end_class');
      s = s.replace(/classDef\s+end\b/g, 'classDef end_class');
      
      // Handle [[ ]] - Subroutine
      s = s.replace(/([A-Za-z0-9_]+)\s*(\[\[)([\s\S]*?)(\]\])/g, (m, id, open, c, close) => {
         if (c.startsWith('"')) return m;
         return `${id}${open}"${c.replace(/"/g, '#quot;')}"${close}`;
      });
      
      // Handle [] but exclude [[
      s = s.replace(/([A-Za-z0-9_]+)\s*(\[)(?!\[)([\s\S]*?)(\])/g, (m, id, _o, c, _cl) => {
         if (c.startsWith('"')) return m;
         return `${id}["${c.replace(/"/g, '#quot;')}"]`;
      });
      
      // Handle {} but exclude {{
      s = s.replace(/([A-Za-z0-9_]+)\s*(\{)(?!\{)([\s\S]*?)(\})/g, (m, id, _o, c, _cl) => {
         if (c.startsWith('"')) return m;
         return `${id}{"${c.replace(/"/g, '#quot;')}"}`;
      });
      
      // Handle () but exclude ((
      s = s.replace(/([A-Za-z0-9_]+)\s*(\()(?! \()(.*?)(\))/g, (m, id, _o, c, _cl) => {
         if (c.startsWith('(')) return m; 
         if (c.startsWith('"')) return m;
         return `${id}("${c.replace(/"/g, '#quot;')}")`;
      });
      
      return s;
    };
    
    try {
      const sanitizedContent = improvedSanitize(content);
      
      if (await mermaid.parse(sanitizedContent)) {
        const { svg } = await mermaid.render(uniqueId + '-svg', sanitizedContent);
        
        // Update caches
        mermaidCache.set(content, svg);
        lastValidMermaidCache.set(uniqueId, svg);
        
        div.innerHTML = svg;
        div.setAttribute('data-processed', 'true');
        div.setAttribute('data-content-hash', String(content.length));
      }
    } catch (error: any) {
      if (props.loading) {
        const lastValid = lastValidMermaidCache.get(uniqueId);
        if (lastValid) {
          if (div.innerHTML.includes('mermaid-placeholder')) {
             div.innerHTML = lastValid;
          }
          return;
        }
        return;
      }
      
      console.error('Mermaid render error:', error);
      div.innerHTML = `<div class="text-red-500 bg-red-50 p-2 rounded text-sm overflow-auto max-h-32">
        <div>流程图生成失败</div>
        <div class="text-xs mt-1 text-red-400">${error.message || String(error)}</div>
      </div>`;
    }
  });
};

// Handle interactions
const handleClick = async (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  
  // Copy Code
  const copyBtn = target.closest('.copy-btn');
  if (copyBtn) {
    const btn = copyBtn as HTMLElement;
    const code = decodeURIComponent(btn.getAttribute('data-code') || '');
    if (code) {
      try {
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(code);
        } else {
          // Fallback
          const textArea = document.createElement('textarea');
          textArea.value = code;
          textArea.style.position = 'fixed';
          textArea.style.left = '-9999px';
          document.body.appendChild(textArea);
          textArea.focus();
          textArea.select();
          document.execCommand('copy');
          document.body.removeChild(textArea);
        }
        
        const originalHtml = btn.innerHTML;
        btn.innerHTML = '<span class="text-green-500">✓</span><span class="text-green-500">已复制</span>';
        setTimeout(() => btn.innerHTML = originalHtml, 2000);
      } catch (err) {
        console.error('Failed to copy:', err);
      }
    }
    return;
  }

  // Export Table
  const exportBtn = target.closest('.export-table-btn');
  if (exportBtn) {
    const wrapper = exportBtn.closest('.table-wrapper');
    const table = wrapper?.querySelector('table');
    if (table) {
      try {
        const XLSX = await import('xlsx');
        const wb = XLSX.utils.table_to_book(table);
        XLSX.writeFile(wb, `table-export-${Date.now()}.xlsx`);
      } catch (err) {
        console.error('Failed to export table:', err);
        alert('导出表格失败');
      }
    }
  }
};

watch(() => props.content, () => {
  renderMermaid();
});

onMounted(() => {
  renderMermaid();
});
</script>

<template>
  <div ref="rootRef" @click="handleClick" class="markdown-body max-w-full break-words">
    <Vue3MarkdownIt 
      :source="content" 
      :plugins="plugins"
      :html="mdOptions.html"
      :linkify="mdOptions.linkify"
      :typographer="mdOptions.typographer"
      :highlight="mdOptions.highlight"
    />
  </div>
</template>

<style>
/* Basic Markdown Styles are now handled by github-markdown-css */
/* We just need to ensure some overrides or custom styles if needed */

.markdown-body {
  background-color: transparent !important;
  font-family: inherit !important;
}

.dark .markdown-body {
  color: #e5e7eb;
}

/* Mermaid styles */
.mermaid {
  background: white;
  padding: 1rem;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  overflow-x: auto;
  display: flex;
  justify-content: center;
}

.dark .mermaid {
  background: #1f2937;
}

/* Custom Code Block Styles - Mac Window Style */
.chat-code-block {
  margin: 1rem 0;
  border-radius: 0.5rem;
  overflow: hidden;
  background-color: #1e1e1e; /* Dark background for all themes */
  border: 1px solid #333;
  max-width: 100%;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.dark .chat-code-block {
  background-color: #1e1e1e;
  border-color: #333;
}

.chat-code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  background-color: #2d2d2d; /* Dark header */
  border-bottom: 1px solid #333;
  position: sticky;
  top: 0;
  z-index: 10;
}

.mac-buttons {
  display: flex;
  gap: 6px;
}

.mac-button {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.mac-button.red { background-color: #ff5f56; }
.mac-button.yellow { background-color: #ffbd2e; }
.mac-button.green { background-color: #27c93f; }

.code-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.chat-code-lang {
  font-size: 0.75rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-weight: 500;
  color: #9ca3af; /* Light gray text */
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #9ca3af;
  transition: color 0.2s;
  background: transparent;
  border: none;
  cursor: pointer;
}

.copy-btn:hover {
  color: #ffffff;
}

.chat-code-content {
  overflow-x: auto;
}

.chat-code-content pre {
  margin: 0 !important;
  padding: 1rem !important;
  background-color: transparent !important;
  border: 0 !important;
  color: #e5e7eb !important; /* Force light text */
}

/* Table styles */
.table-wrapper {
  overflow-x: auto;
  width: 100%;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  display: table; /* Ensure it behaves like a table */
  background-color: transparent; /* Allow wrapper background to show */
}

.markdown-body th,
.markdown-body td {
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
}

.dark .markdown-body th,
.dark .markdown-body td {
  border-color: #374151;
}

.markdown-body th {
  background-color: #f9fafb;
  font-weight: 600;
}

.dark .markdown-body th {
  background-color: #1f2937;
}

.markdown-body tr {
  background-color: white;
}

.dark .markdown-body tr {
  background-color: #1f2937;
}

.markdown-body tr:nth-child(2n) {
  background-color: #f9fafb;
}

.dark .markdown-body tr:nth-child(2n) {
  background-color: #111827;
}

/* Override github-markdown-css specificities */
.markdown-body pre {
  background-color: transparent !important;
}
</style>
