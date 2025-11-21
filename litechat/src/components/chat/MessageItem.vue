<script setup lang="ts">
import { computed, onMounted, nextTick, watch } from 'vue';
import MarkdownIt from 'markdown-it';
import hljs from 'highlight.js';
import mermaid from 'mermaid';
import type { Message, Attachment } from '../../stores/chat';
import { User, Bot, Quote, FileText } from 'lucide-vue-next';
import MessageActions from './MessageActions.vue';
import 'highlight.js/styles/github-dark.css'; // Import highlight.js style

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
        return '<pre class="hljs"><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }

    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
});

// Custom fence rule for mermaid
const defaultFence = md.renderer.rules.fence || function(tokens: any[], idx: number, options: any, _env: any, self: any) {
  return self.renderToken(tokens, idx, options);
};

md.renderer.rules.fence = (tokens: any[], idx: number, options: any, env: any, self: any) => {
  const token = tokens[idx];
  const info = token.info.trim();
  
  if (info === 'mermaid') {
    const id = `mermaid-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    return `<div class="mermaid" id="${id}">${token.content}</div>`;
  }
  
  return defaultFence(tokens, idx, options, env, self);
};

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

// Watch for content changes to re-render mermaid
watch(() => props.message.content, () => {
  renderMermaid();
  addCopyButtonsToCodeBlocks();
});

// Add copy buttons to code blocks
onMounted(() => {
  renderMermaid();
  nextTick(() => {
    addCopyButtonsToCodeBlocks();
  });
});

const addCopyButtonsToCodeBlocks = () => {
  // Wait for DOM update
  setTimeout(() => {
    const codeBlocks = document.querySelectorAll('.markdown-body pre');
    codeBlocks.forEach((block) => {
      if (block.querySelector('.copy-code-button')) return; // Already added
      
      const button = document.createElement('button');
      button.className = 'copy-code-button';
      button.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>';
      button.title = '复制代码';
      
      button.addEventListener('click', async () => {
        const code = block.querySelector('code')?.textContent || '';
        try {
          await navigator.clipboard.writeText(code);
          button.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>';
          button.classList.add('copied');
          setTimeout(() => {
            button.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>';
            button.classList.remove('copied');
          }, 2000);
        } catch (err) {
          console.error('Failed to copy code:', err);
        }
      });
      
      (block as HTMLElement).style.position = 'relative';
      block.appendChild(button);
    });
  }, 100);
};

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

.markdown-body pre {
  background-color: #f3f4f6;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  position: relative;
  max-width: 100%;
  margin-bottom: 1em;
  -webkit-overflow-scrolling: touch; /* Smooth scrolling on iOS */
}

.markdown-body code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.875em;
  word-break: break-word;
}

.markdown-body pre code {
  word-break: normal;
  white-space: pre;
}

.dark .markdown-body pre {
  background-color: #1f2937;
}

.markdown-body table {
  display: block;
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border-collapse: collapse;
  margin-bottom: 1em;
}

.markdown-body th,
.markdown-body td {
  padding: 0.5rem;
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

/* Code copy button */
.copy-code-button {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  padding: 0.375rem;
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  cursor: pointer;
  opacity: 0; /* Hidden by default on desktop */
  transition: opacity 0.2s, background-color 0.2s;
  color: #6b7280;
  z-index: 10;
}

/* Always show copy button on mobile/touch devices */
@media (hover: none) {
  .copy-code-button {
    opacity: 1;
    background-color: rgba(255, 255, 255, 0.95);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }
  
  .dark .copy-code-button {
    background-color: rgba(31, 41, 55, 0.95);
  }
}

.dark .copy-code-button {
  background-color: rgba(31, 41, 55, 0.9);
  border-color: #374151;
  color: #9ca3af;
}

.markdown-body pre:hover .copy-code-button {
  opacity: 1;
}

.copy-code-button:hover {
  background-color: #f9fafb;
  color: #374151;
}

.dark .copy-code-button:hover {
  background-color: #374151;
  color: #d1d5db;
}

.copy-code-button.copied {
  color: #10b981;
  opacity: 1;
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
