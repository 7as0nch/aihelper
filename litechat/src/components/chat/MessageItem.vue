<script setup lang="ts">
import { computed, onMounted, nextTick, watch } from 'vue';
import MarkdownIt from 'markdown-it';
import hljs from 'highlight.js';
import mermaid from 'mermaid';
import type { Message, Attachment } from '../../stores/chat';
import { User, Bot, Quote, FileText } from 'lucide-vue-next';
import MessageActions from './MessageActions.vue';

const props = defineProps<{
  message: Message
}>();

const emit = defineEmits<{
  (e: 'quote', messageId: string, content: string): void;
  (e: 'regenerate', messageId: string): void;
  (e: 'previewImage', url: string): void;
}>();

// Initialize mermaid
mermaid.initialize({
  startOnLoad: false,
  theme: 'default',
  securityLevel: 'loose',
});

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
    return `<div class="mermaid" id="${id}">${token.content}</div>`;
  }
  
  const lang = info ? info.split(/\s+/)[0] : '';
  const languageLabel = lang ? lang.toUpperCase() : 'TEXT';
  const code = options.highlight ? options.highlight(token.content, lang) : md.utils.escapeHtml(token.content);
  
  return `
    <div class="code-block-wrapper my-4 rounded-lg overflow-hidden bg-[#f3f4f6] dark:bg-[#1f2937] border border-gray-200 dark:border-gray-700">
      <div class="code-block-header flex items-center justify-between px-4 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
        <span class="text-xs font-mono font-medium text-gray-500 dark:text-gray-400">${languageLabel}</span>
        <button class="copy-btn flex items-center gap-1 text-xs text-gray-500 hover:text-primary transition-colors" data-code="${encodeURIComponent(token.content)}">
          <span class="copy-icon">📋</span>
          <span class="copy-text">复制</span>
        </button>
      </div>
      <div class="overflow-x-auto">
        <pre class="!m-0 !p-4 !bg-transparent !border-0"><code class="hljs language-${lang}">${code}</code></pre>
      </div>
    </div>
  `;
};

// Custom table rules for horizontal scrolling
md.renderer.rules.table_open = () => '<div class="table-wrapper overflow-x-auto my-4 w-full border border-gray-200 dark:border-gray-700 rounded-lg"><table class="w-full text-left text-sm">';
md.renderer.rules.table_close = () => '</table></div>';

const renderedContent = computed(() => {
  return md.render(props.message.content);
});

const renderMermaid = async () => {
  await nextTick();
  const mermaidDivs = document.querySelectorAll('.mermaid');
  mermaidDivs.forEach(async (div) => {
    const id = div.id;
    const content = div.textContent || '';
    if (div.getAttribute('data-processed')) return;
    
    try {
      const { svg } = await mermaid.render(id + '-svg', content);
      div.innerHTML = svg;
      div.setAttribute('data-processed', 'true');
    } catch (error) {
      console.error('Mermaid render error:', error);
      div.innerHTML = `<div class="text-red-500 bg-red-50 p-2 rounded">Diagram render error</div>`;
    }
  });
};

// Handle copy button clicks via delegation
const handleMessageClick = async (event: MouseEvent) => {
  const target = (event.target as HTMLElement).closest('.copy-btn');
  if (!target) return;
  
  const btn = target as HTMLElement;
  const code = decodeURIComponent(btn.getAttribute('data-code') || '');
  
  if (code) {
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(code);
      } else {
        // Fallback for older browsers or non-secure contexts (often needed on Android WebViews)
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
      // Optional: Show a toast or alert if copy fails
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
  a.download = file.name;
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
        class="inline-block text-left max-w-full"
        :class="[
          message.role === 'user' 
            ? 'bg-primary text-white rounded-2xl rounded-tr-sm px-4 py-2' 
            : 'prose dark:prose-invert max-w-none'
        ]"
      >
        <div v-if="message.role === 'user'">
          {{ message.content }}
        </div>
        <div 
          v-else 
          v-html="renderedContent"
          class="markdown-body max-w-full overflow-x-auto break-words"
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
            @click="handleImageClick(img.url!)"
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
        
        <!-- Streaming Cursor -->
        <span 
          v-if="message.isStreaming" 
          class="inline-block w-2 h-4 ml-1 bg-primary align-middle animate-pulse"
        ></span>
      </div>
      
      <!-- Message Actions (for assistant messages only) -->
      <MessageActions 
        v-if="message.role === 'assistant' && !message.isStreaming"
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
</style>
