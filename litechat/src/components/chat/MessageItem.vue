<script setup lang="ts">
import { computed, onMounted, nextTick, watch, ref } from 'vue';
import MarkdownIt from 'markdown-it';
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
import type { Message, Attachment } from '../../stores/chat';
import { useChatStore } from '../../stores/chat';
import { User, Bot, Quote, FileText, Wrench, Link as LinkIcon, ChevronDown, ChevronRight } from 'lucide-vue-next';
import MessageActions from './MessageActions.vue';

const props = defineProps<{
  message: Message;
  isLastMessage?: boolean;
}>();

const chatStore = useChatStore();
const isReasoningCollapsed = ref(false);

// Auto-collapse reasoning when content starts generating
watch(() => props.message.content, (newContent, oldContent) => {
  if (newContent && !oldContent && props.message.reasoning_content) {
    isReasoningCollapsed.value = true;
  }
});

const emit = defineEmits<{
  (e: 'quote', messageId: string, content: string): void;
  (e: 'regenerate', messageId: string): void;
  (e: 'previewImage', url: string): void;
}>();

// Lazy load and initialize Mermaid only when needed
let mermaidInstance: any = null;

const initMermaid = async () => {
  if (!mermaidInstance) {
    const mermaidModule = await import('mermaid');
    mermaidInstance = mermaidModule.default;
    
    // Initialize mermaid
    mermaidInstance.initialize({
      startOnLoad: false,
      theme: 'default',
      securityLevel: 'loose',
      suppressErrorRendering: true, // Suppress default error rendering
    });
    
    // Override parseError to prevent console spam
    mermaidInstance.parseError = () => {};
  }
  
  return mermaidInstance;
};

const md: MarkdownIt = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight: function (str, lang): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang, ignoreIllegals: true }).value;
      } catch (__) {}
    }
    return md.utils.escapeHtml(str);
  }
});

// Custom fence rule for mermaid and code blocks
md.renderer.rules.fence = (tokens: any[], idx: number, options: any, _env: any, _self: any) => {
  const token = tokens[idx];
  const info = token.info.trim();
  
  if (info === 'mermaid') {
    const id = `mermaid-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const encodedCode = encodeURIComponent(token.content);
    return `
      <div class="mermaid-view my-4" id="${id}" data-code="${encodedCode}">
        <div class="flex items-center justify-center p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg text-gray-500 text-sm animate-pulse border border-gray-100 dark:border-gray-700">
           <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
             <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
             <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
           </svg>
           流程图生成中...
        </div>
      </div>
    `;
  }
  
  const lang = info ? info.split(/\s+/)[0] : '';
  const languageLabel = lang ? lang.toUpperCase() : 'TEXT';
  const code = options.highlight ? options.highlight(token.content, lang) : md.utils.escapeHtml(token.content);
  
  return `
    <div class="chat-code-block">
      <div class="chat-code-header">
        <span class="chat-code-lang">${languageLabel}</span>
        <button class="copy-btn" data-code="${encodeURIComponent(token.content)}">
          <span class="copy-icon">📋</span>
          <span class="copy-text">复制</span>
        </button>
      </div>
      <div class="chat-code-content">
        <pre><code class="hljs language-${lang}">${code}</code></pre>
      </div>
    </div>
  `;
};

// Custom table rules for horizontal scrolling and export button
md.renderer.rules.table_open = () => '<div class="table-wrapper overflow-x-auto my-4 w-full border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800"><table class="w-full text-left text-sm">';
md.renderer.rules.table_close = () => `
  </table>
  <div class="flex justify-end p-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
    <button class="export-table-btn flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md transition-colors">
      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-table"><path d="M12 3v18"/><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M3 15h18"/></svg>
      Export to Sheets
    </button>
  </div>
</div>`;

const renderedContent = computed(() => {
  try {
    return md.render(props.message.content);
  } catch (e) {
    console.error('Markdown render error:', e);
    return '<div class="text-gray-500 italic p-2">当前内容暂不支持渲染</div>';
  }
});

const renderMermaid = async () => {
  await nextTick();
  const mermaidDivs = document.querySelectorAll('.mermaid-view');
  
  // Only load mermaid if there are mermaid diagrams to render
  if (mermaidDivs.length === 0) return;
  
  // Dynamically load mermaid when needed
  const mermaid = await initMermaid();
  
  mermaidDivs.forEach(async (div) => {
    const id = div.id;
    const content = decodeURIComponent(div.getAttribute('data-code') || '');
    if (div.getAttribute('data-processed')) return;
    
    const sanitizeMermaid = (code: string) => {
      let sanitized = code;
      
      // Fix unquoted brackets in node labels: id[label with [brackets]]
      // Exclude shapes starting with [[ or [{
      sanitized = sanitized.replace(/([A-Za-z0-9_]+)\s*\[(?!\[|\{)(.*?)\]/g, (match, id, content) => {
        if (content.includes('[') && !content.startsWith('"')) {
          return `${id}["${content}"]`;
        }
        return match;
      });
      
      // Fix unquoted parentheses in node labels: id(label with (parens))
      // Exclude shapes starting with ((
      sanitized = sanitized.replace(/([A-Za-z0-9_]+)\s*\((?!\()(.*?)\)/g, (match, id, content) => {
        if (content.includes('(') && !content.startsWith('"')) {
          return `${id}("${content}")`;
        }
        return match;
      });
    
      return sanitized;
    };

    try {
      const sanitizedContent = sanitizeMermaid(content);
      
      // Validate content before rendering
      if (await mermaid.parse(sanitizedContent)) {
        const { svg } = await mermaid.render(id + '-svg', sanitizedContent);
        div.innerHTML = svg;
        div.setAttribute('data-processed', 'true');
      }
    } catch (error) {
      // If streaming, suppress error as it might be incomplete code
      if (props.isLastMessage && chatStore.isLoading) {
        // Keep the loading state
        return;
      }
      console.error('Mermaid render error:', error);
      div.innerHTML = `<div class="text-red-500 bg-red-50 p-2 rounded text-sm">流程图生成失败</div>`;
    }
  });
};

// Handle click delegation
const handleMessageClick = async (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  
  // Handle Copy Code
  const copyBtn = target.closest('.copy-btn');
  if (copyBtn) {
    const btn = copyBtn as HTMLElement;
    const code = decodeURIComponent(btn.getAttribute('data-code') || '');
    if (code) {
      try {
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(code);
        } else {
          const textArea = document.createElement('textarea');
          textArea.value = code;
          textArea.style.position = 'fixed';
          textArea.style.left = '-9999px';
          textArea.style.top = '0';
          document.body.appendChild(textArea);
          textArea.focus();
          textArea.select();
          document.execCommand('copy');
          document.body.removeChild(textArea);
        }
        
        const originalHtml = btn.innerHTML;
        btn.innerHTML = '<span class="text-green-500">✓</span><span class="text-green-500">已复制</span>';
        setTimeout(() => {
          btn.innerHTML = originalHtml;
        }, 2000);
      } catch (err) {
        console.error('Failed to copy:', err);
      }
    }
    return;
  }

  // Handle Table Export
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

// Watch for content changes to re-render mermaid
watch(() => props.message.content, () => {
  renderMermaid();
});

onMounted(() => {
  renderMermaid();
});

const handleQuote = (messageId: string, content: string) => {
  emit('quote', messageId, content);
};

const handleRegenerate = (id: string) => {
  emit('regenerate', id);
};

const handleImageClick = (url: string) => {
  emit('previewImage', url);
};

const handleFileClick = (file: Attachment) => {
  // Create a temporary anchor element to trigger download
  const a = document.createElement('a');
  a.href = file.url || '#';
  console.log(file, a);
  // Ensure filename has extension and timestamp
  let filename = file.name;
  const timestamp = new Date().getTime();
  
  if (file.type && !filename.includes('.')) {
    const ext = file.type.split('/')[1];
    if (ext) {
      filename = `${filename}-${timestamp}.${ext}`;
    } else {
      filename = `${filename}-${timestamp}`;
    }
  } else {
    // Insert timestamp before extension
    const parts = filename.split('.');
    if (parts.length > 1) {
      const ext = parts.pop();
      filename = `${parts.join('.')}-${timestamp}.${ext}`;
    } else {
      filename = `${filename}-${timestamp}`;
    }
  }
  
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
};
</script>

<template>
  <div 
    class="flex gap-4 max-w-3xl mx-auto w-full p-4 animate-fade-in"
    :class="[message.role === 'user' ? 'flex-row-reverse' : 'flex-row']"
  >
    <!-- Avatar -->
    <div 
      class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
      :class="[
        message.role === 'user' 
          ? 'bg-blue-100 text-primary dark:bg-blue-900/30' 
          : 'bg-orange-100 text-orange-600 dark:bg-orange-900/30'
      ]"
    >
      <User v-if="message.role === 'user'" class="w-5 h-5" />
      <Bot v-else class="w-5 h-5" />
    </div>

    <!-- Content -->
    <div 
      class="flex-1 min-w-0 overflow-hidden"
      :class="[message.role === 'user' ? 'text-right' : 'text-left']"
    >
      <div 
        class="text-left max-w-full"
        :class="[
          message.role === 'user' 
            ? 'inline-block bg-primary text-white rounded-2xl rounded-tr-sm px-4 py-2 user-message-bubble' 
            : 'block w-full prose dark:prose-invert max-w-none'
        ]"
        :style="message.role === 'user' ? { backgroundColor: '#3b82f6', color: '#ffffff' } : {}"
      >
        <!-- Reasoning Content -->
        <div v-if="message.reasoning_content" class="mb-4">
          <div 
            class="bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
          >
            <button 
              @click="isReasoningCollapsed = !isReasoningCollapsed"
              class="w-full flex items-center justify-between px-3 py-2 text-xs font-medium text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded-full bg-purple-100 dark:bg-purple-900/30 text-purple-600 flex items-center justify-center">
                  <span class="text-[10px]">R</span>
                </div>
                <span>深度思考过程</span>
              </div>
              <component :is="isReasoningCollapsed ? ChevronRight : ChevronDown" class="w-4 h-4" />
            </button>
            
            <div v-show="!isReasoningCollapsed" class="px-3 py-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-800/30">
              <div class="prose dark:prose-invert max-w-none text-xs text-gray-600 dark:text-gray-400 leading-relaxed whitespace-pre-wrap font-mono">
                {{ message.reasoning_content }}
              </div>
            </div>
          </div>
        </div>

        <div v-if="message.role === 'user'">
          {{ message.content }}
        </div>
        <div 
          v-else 
          v-html="renderedContent"
          class="markdown-body max-w-full break-words"
          @click="handleMessageClick"
        ></div>

        <!-- Quote Display (for user messages with quotes) -->
        <div v-if="message.role === 'user' && message.quoteContent" class="mt-2 pt-2 border-t border-white/20">
          <div class="flex items-start gap-1.5 opacity-80 bg-white/10 p-2 rounded text-sm">
            <Quote class="w-3 h-3 shrink-0 mt-0.5" />
            <p class="line-clamp-2">{{ message.quoteContent }}</p>
          </div>
        </div>

        <!-- Timestamp -->
        <div class="mt-1 text-[10px] opacity-60 text-right">
          {{ new Date(message.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
        </div>

        <!-- Image Attachments -->
        <div v-if="message.attachments?.some(a => a.type === 'image')" class="flex gap-2 mb-2 overflow-x-auto pb-2 no-scrollbar">
          <div 
            v-for="img in message.attachments.filter(a => a.type === 'image')" 
            :key="img.id"
            class="relative group cursor-zoom-in shrink-0"
            @click="handleImageClick(img.url || '')"
          >
            <img 
              :src="img.url" 
              class="h-48 w-auto rounded-lg border border-gray-200 dark:border-gray-700 hover:opacity-90 transition-opacity object-cover"
              :alt="img.name"
            />
          </div>
        </div>

        <!-- File Attachments -->
        <div v-if="message.attachments?.some(a => a.type === 'file')" class="flex flex-col gap-2 mb-2">
          <div 
            v-for="file in message.attachments.filter(a => a.type === 'file')" 
            :key="file.id"
            class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors cursor-pointer group"
            @click="handleFileClick(file)"
          >
            <div class="p-2 bg-white dark:bg-gray-700 rounded-lg text-gray-500 group-hover:text-primary transition-colors">
              <FileText class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ file.name }}</div>
              <div class="text-xs text-gray-500">点击下载</div>
            </div>
          </div>
        </div>
        
        <!-- Calling Tools -->
        <div v-if="message.callingTools && message.callingTools.length > 0" class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700/50">
          <div class="text-xs font-medium text-gray-500 mb-2 flex items-center gap-1.5">
            <Wrench class="w-3.5 h-3.5" />
            <span>使用的工具</span>
          </div>
          <div class="flex flex-wrap gap-2">
            <div 
              v-for="(tool, index) in message.callingTools" 
              :key="index"
              class="flex items-center gap-1.5 px-2.5 py-1.5 bg-gray-50 dark:bg-gray-800/50 border border-gray-100 dark:border-gray-700 rounded-md text-xs text-gray-600 dark:text-gray-300"
            >
              <span class="font-medium">{{ tool.name }}</span>
            </div>
          </div>
        </div>

        <!-- Quote Search Links -->
        <div v-if="message.quoteSearchLinks && message.quoteSearchLinks.length > 0" class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700/50">
          <div class="text-xs font-medium text-gray-500 mb-2 flex items-center gap-1.5">
            <LinkIcon class="w-3.5 h-3.5" />
            <span>引用来源</span>
          </div>
          <div class="flex flex-col gap-1.5">
            <a 
              v-for="(link, index) in message.quoteSearchLinks" 
              :key="index"
              :href="link.url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-start gap-2 p-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors group border border-transparent hover:border-gray-100 dark:hover:border-gray-700"
            >
              <div class="shrink-0 mt-0.5 text-gray-400 group-hover:text-primary transition-colors">
                <span class="text-[10px] font-mono border border-current rounded px-1 flex items-center justify-center min-w-[18px] h-[18px]">{{ index + 1 }}</span>
              </div>
              <div class="min-w-0">
                <div class="text-xs font-medium text-gray-700 dark:text-gray-200 group-hover:text-primary truncate transition-colors">
                  {{ link.title }}
                </div>
                <div class="text-[10px] text-gray-400 truncate mt-0.5">
                  {{ link.content }}
                </div>
              </div>
            </a>
          </div>
        </div>

        <!-- Streaming Cursor -->
        <span 
          v-if="isLastMessage && chatStore.isLoading" 
          class="inline-block w-2 h-4 ml-1 bg-primary align-middle animate-pulse"
        ></span>
      </div>
      
      <!-- Message Actions (for assistant messages only) -->
      <MessageActions 
        v-if="message.role === 'assistant' && !(isLastMessage && chatStore.isLoading)"
        :message-content="message.content"
        :message-id="message.id"
        @quote="handleQuote"
        @regenerate="handleRegenerate"
      />
    </div>
  </div>
</template>

<style>
/* Basic Markdown Styles for AI response */
.markdown-body {
  font-size: 0.95rem;
  line-height: 1.6;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.markdown-body p {
  margin-bottom: 0.75em;
  white-space: pre-wrap;
}

.markdown-body ul {
  list-style-type: disc;
  padding-left: 1.5em;
  margin-bottom: 1em;
}

.markdown-body ol {
  list-style-type: decimal;
  padding-left: 1.5em;
  margin-bottom: 1em;
}

.markdown-body li {
  margin-bottom: 0.25em;
}

/* Code block styles are now handled by utility classes in the render function */
/* But we keep some resets */
.markdown-body pre code {
  word-break: normal;
  white-space: pre;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875em;
}

/* Table styles */
.markdown-body table {
  width: 100%;
  border-collapse: collapse;
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

/* Custom Code Block Styles */
.chat-code-block {
  margin: 1rem 0;
  border-radius: 0.5rem;
  overflow: hidden;
  background-color: #f3f4f6;
  border: 1px solid #e5e7eb;
  max-width: 100%;
}

.dark .chat-code-block {
  background-color: #1f2937;
  border-color: #374151;
}

.chat-code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  background-color: #f3f4f6; /* Match block bg for seamless look or slightly darker */
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 10;
}

.dark .chat-code-header {
  background-color: #1f2937;
  border-color: #374151;
}

.chat-code-lang {
  font-size: 0.75rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-weight: 500;
  color: #6b7280;
}

.dark .chat-code-lang {
  color: #9ca3af;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
  transition: color 0.2s;
  background: transparent;
  border: none;
  cursor: pointer;
}

.dark .copy-btn {
  color: #9ca3af;
}

.copy-btn:hover {
  color: #3b82f6; /* Primary color */
}

.chat-code-content {
  overflow-x: auto;
}

.chat-code-content pre {
  margin: 0 !important;
  padding: 1rem !important;
  background-color: transparent !important;
  border: 0 !important;
}
</style>
