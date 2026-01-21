<script setup lang="ts">
import { onMounted, ref, watch, h, nextTick } from 'vue';
import { onClickOutside } from '@vueuse/core';
import { 
  Button, 
  Tooltip, 
  message 
} from 'ant-design-vue';
import { 
  HighlightOutlined, 
  CopyOutlined,
  ClearOutlined
} from '@ant-design/icons-vue';

const props = defineProps({
  value: { type: String, default: '' },
  modelValue: { type: String, default: '' },
});

const emit = defineEmits(['update:value', 'update:modelValue', 'change']);

// State
const editorRef = ref<HTMLElement | null>(null);
const containerRef = ref<HTMLElement | null>(null);

// Floating Menu State
const showFloatingMenu = ref(false);
const floatingPosition = ref({ x: 0, y: 0 });
const currentSelectionRange = ref<Range | null>(null);
const selectedText = ref('');

// --- Serialization & Parsing ---

const DEFAULT_PLACEHOLDER = '此为编辑块，输入内容替换此内容';

// ... (previous code)

const serializeContent = () => {
  if (!editorRef.value) return '';

  const processNode = (node: Node): string => {
     if (node.nodeType === Node.TEXT_NODE) return node.nodeValue || '';
     if (node.nodeType !== Node.ELEMENT_NODE) return '';
     
     const el = node as HTMLElement;
     
     if (el.classList.contains('prompt-slot-tag')) {
        const input = el.querySelector('input');
        const format = el.dataset.format || 'tag'; 
        
        // If input exists, use its value. If not, check innerText (legacy)
        let val = input ? input.value : el.innerText;
        
        if (format === 'handlebars') {
            // Variable Block: {{ val }}
            // If empty, what do we save? {{.}} or {{.placeholder}}?
            // "content" usually means the variable name.
            return `{{ ${val} }}`;
        }
        
        if (format === 'parens') {
            // Edit Block: "val"
            return `"${val}"`;
        }
        
        // Default Legacy Tag Format convert to Edit Block "val" if has value
        const placeholder = el.dataset.placeholder || DEFAULT_PLACEHOLDER;
        if (val) {
             return `"${val}"`;
        }
        return `{#InputSlot placeholder="${placeholder}" mode="input"#}${val}{#/InputSlot#}`;
      }
     
     if (el.tagName === 'BR') return '\n';
     
     let content = '';
     for (const child of Array.from(el.childNodes)) {
         content += processNode(child);
     }
     
     if (el.tagName === 'DIV' || el.tagName === 'P') {
         // Block elements imply a newline separation relative to siblings
         // Standardize to \n
         return content + '\n';
     }
     
     return content;
  };
  
  let result = '';
  editorRef.value.childNodes.forEach(node => {
      result += processNode(node);
  });
  
  return result.replace(/\n+$/, '');
};

const escapeHtml = (unsafe: string) => {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
};

const parseContent = (text: string) => {
  if (!text) return '';
  
  // 1. Escape HTML first
  let html = escapeHtml(text);

  // 2. Restore Slots
  
  // A. Standard Tag Format
  html = html.replace(/{#InputSlot placeholder=&quot;(.*?)&quot; mode=&quot;input&quot;#}([\s\S]*?){#\/InputSlot#}/g, (_, placeholder, content) => {
     return createSlotHtml({ placeholder, value: content, format: 'tag' });
  });

  // B. Handlebars Variable Block {{ val }}
  html = html.replace(/{{(.*?)\s*}}/g, (_, content) => {
      const val = content.trim();
      return createSlotHtml({ placeholder: val || DEFAULT_PLACEHOLDER, value: val, format: 'handlebars' });
  });
  
  // C. Edit Block "val" 
  html = html.replace(/&quot;(.*?)&quot;/g, (_, value) => {
      return createSlotHtml({ placeholder: DEFAULT_PLACEHOLDER, value: value, format: 'parens' }); 
  });
  
  // 3. Markdown Support
  // Headers
  html = html.replace(/(^|\n)(#+ )(.+?)(\n|$)/g, '$1<span class="md-heading">$2$3</span>$4');
  
  // Bold (**text**)
  html = html.replace(/(\*\*|__)(.*?)\1/g, '<span class="md-bold">$1$2$1</span>');
  
  // Italic (*text*) - Careful with conflict with bold, usually match bold first.
  // Using simple lookahead/behind is hard in JS regex without lookbehind support in all browsers?
  // Let's stick to simple *text* but avoid matching within **text**.
  // Actually, standard regex for *text* is `/\*([^*]+)\*/`.
  html = html.replace(/(\*)([^*]+)(\*)/g, '<span class="md-italic">$1$2$3</span>');

  // Lists (- item)
  html = html.replace(/(^|\n)(- )(.+?)(\n|$)/g, '$1<span class="md-list-item">$2$3</span>$4');
  
  // Numbered Lists (1. item)
  html = html.replace(/(^|\n)(\d+\. )(.+?)(\n|$)/g, '$1<span class="md-list-item">$2$3</span>$4');

  // Code Block (inline `code`)
  // html = html.replace(/`([^`]+)`/g, '<span class="md-code">`$1`</span>');

  // 4. Newlines
  html = html.replace(/\n/g, '<br/>');

  return html;
};

const createSlotHtml = ({ placeholder, value, format }: { placeholder: string; value: string; format: string }) => {
    const safePh = placeholder; 
    let displayValue = value;
    let extraClass = '';
    
    if (format === 'handlebars') {
        extraClass = 'slot-variable';
    }
    
    // Construct HTML. 
    return `<span class="prompt-slot-tag" data-type="input-slot" data-format="${format}" data-placeholder="${safePh}" contenteditable="false"><input type="text" class="slot-input ${extraClass}" placeholder="${safePh}" value="${displayValue}"></span>`;
};

const triggerUpdate = () => {
    const newVal = serializeContent();
    emit('update:value', newVal);
    emit('update:modelValue', newVal);
    emit('change', newVal);
};

// Auto-Width Helper
const measureTextWidth = (text: string, font: string) => {
  const canvas = document.createElement('canvas');
  const context = canvas.getContext('2d');
  if (context) {
    context.font = font;
    return context.measureText(text).width;
  }
  return 0;
};

const updateInputWidth = (input: HTMLInputElement) => {
   const val = input.value || input.placeholder;
   const computedStyle = window.getComputedStyle(input);
   const font = `${computedStyle.fontWeight} ${computedStyle.fontSize} ${computedStyle.fontFamily}`;
   const width = measureTextWidth(val, font);
   // Add padding (px-2 = ~8-12px total? CSS says 6px*2?)
   // CSS: padding: 2px 6px;
   input.style.width = `${width + 16}px`; 
};

// --- Editor Core Logic ---

const handleSelectionChange = () => {
  const selection = window.getSelection();
  if (!selection || selection.rangeCount === 0 || selection.isCollapsed) {
    showFloatingMenu.value = false;
    return;
  }

  const text = selection.toString().trim();
  if (!text) {
    showFloatingMenu.value = false;
    return;
  }

  if (editorRef.value && !editorRef.value.contains(selection.anchorNode)) {
    return;
  }

  selectedText.value = text;
  currentSelectionRange.value = selection.getRangeAt(0);
  
  // Calculate position absolute for Teleport
  const rect = currentSelectionRange.value.getBoundingClientRect();
  // We use viewport coordinates + scrollY for absolute positioning in body
  floatingPosition.value = { 
      x: rect.left + (rect.width / 2), 
      y: rect.top + window.scrollY 
  };
  showFloatingMenu.value = true;
};

const convertToSlot = (type: 'edit' | 'variable') => {
  if (!currentSelectionRange.value) return;

  const range = currentSelectionRange.value;
  const text = range.toString();
  
  // Determine Type
  const format = type === 'variable' ? 'handlebars' : 'parens';
  const placeholder = DEFAULT_PLACEHOLDER;
  
  // Create Slot manually
  const span = document.createElement('span');
  span.className = 'prompt-slot-tag animate-in';
  span.dataset.type = 'input-slot';
  span.dataset.placeholder = placeholder;
  span.dataset.format = format;
  span.contentEditable = 'false';
  
  const input = document.createElement('input');
  input.type = 'text';
  input.className = `slot-input ${format === 'handlebars' ? 'slot-variable' : ''}`;
  input.placeholder = placeholder;
  input.value = text; // Initial value from selection
  
  // Bind input event for auto-width/updates
  input.addEventListener('input', (e) => {
     triggerUpdate();
     updateInputWidth(e.target as HTMLInputElement);
  });

  span.appendChild(input);

  range.deleteContents();
  range.insertNode(span);
  
  window.getSelection()?.removeAllRanges();
  showFloatingMenu.value = false;
  currentSelectionRange.value = null;
  
  // Initial width adjustment
  updateInputWidth(input);
  
  triggerUpdate();
};



const handleInputClick = (e: Event) => {
    // Delegate input events from dynamic HTML (from parseContent)
    const target = e.target as HTMLInputElement;
    if (target && target.classList.contains('slot-input')) {
         // Manually update the value attribute in DOM for serialization 
         // because 'value' attribute doesn't auto-sync with property for innerHTML
         target.setAttribute('value', target.value);
         
         updateInputWidth(target);
    }
    // Always trigger update for any input event (text changes or slot changes)
    triggerUpdate();
};

const handleClick = (e: MouseEvent) => {
  // Focus helper?
};

const handleMouseUp = () => {
  setTimeout(handleSelectionChange, 10);
};

const handleCopy = () => {
  const result = serializeContent();
  navigator.clipboard.writeText(result).then(() => {
    message.success('已复制模版内容');
  });
};

const handleClear = () => {
    if (editorRef.value) {
        editorRef.value.innerHTML = '';
        triggerUpdate();
    }
}

// --- Lifecycle ---
onMounted(() => {
  if (editorRef.value) {
      const initialVal = props.modelValue || props.value || '';
      editorRef.value.innerHTML = parseContent(initialVal);
      
      // Initialize widths
      editorRef.value.querySelectorAll('.slot-input').forEach((input: any) => {
          updateInputWidth(input);
      });
  }
});

// Watchers
watch(() => props.modelValue, (val) => {
    if (editorRef.value && val !== serializeContent()) {
         editorRef.value.innerHTML = parseContent(val || '');
         nextTick(() => {
            editorRef.value?.querySelectorAll('.slot-input').forEach((input: any) => {
                updateInputWidth(input);
            });
         });
    }
});
watch(() => props.value, (val) => {
    if (editorRef.value && val !== serializeContent()) {
         editorRef.value.innerHTML = parseContent(val || '');
         nextTick(() => {
            editorRef.value?.querySelectorAll('.slot-input').forEach((input: any) => {
                updateInputWidth(input);
            });
         });
    }
});

onClickOutside(containerRef, () => {
  showFloatingMenu.value = false;
});

</script>

<template>
  <div class="h-full flex flex-col w-full min-h-0" ref="containerRef">
      <!-- Editor Area -->
      <div class="flex-1 min-h-0 flex flex-col bg-white rounded-lg border border-gray-300 overflow-hidden relative">
        <!-- Toolbar -->
        <div class="h-10 border-b border-gray-200 flex items-center px-4 justify-between bg-gray-50 z-10 w-full">
          <div class="flex items-center gap-4">
             <div class="font-medium text-gray-600 text-sm">提示词编辑器</div>
             <div class="flex items-center gap-2 text-xs text-gray-400">
                <span>选中文字可设置为:</span>
                <span class="flex items-center gap-1"><span class="w-3 h-3 rounded bg-[#f3f0ff] border border-[#d3adf7]"></span>编辑块</span>
                <span class="flex items-center gap-1"><span class="w-3 h-3 rounded bg-[#e6f7ff] border border-[#91caff]"></span>变量块</span>
             </div>
          </div>
          <div class="flex gap-1">
            <Tooltip title="清空"><Button type="text" size="small" shape="circle" :icon="h(ClearOutlined)" @click="handleClear" /></Tooltip>
            <Tooltip title="复制结果"><Button type="text" size="small" shape="circle" :icon="h(CopyOutlined)" @click="handleCopy" /></Tooltip>
          </div>
        </div>

        <!-- Editor Content -->
        <div class="flex-1 min-h-0 overflow-y-auto p-4 cursor-text bg-white" @click="handleClick">
           <div 
             ref="editorRef"
             class="prose max-w-none focus:outline-none text-gray-800 leading-relaxed outline-none min-h-0"
             contenteditable="true"
             spellcheck="false"
             @input="handleInputClick"
             @mouseup="handleMouseUp"
           ></div>
        </div>

        <!-- Floating Popover Menu (Teleported) -->
        <Teleport to="body">
            <div 
               v-if="showFloatingMenu"
               class="absolute z-[9999] transform -translate-x-1/2 -translate-y-full pb-2 transition-all duration-200"
               :style="{ top: `${floatingPosition.y}px`, left: `${floatingPosition.x}px` }"
            >
               <div class="bg-gray-800 text-white rounded shadow-lg p-1 flex gap-1 animate-fade-in-up items-center">
                  <div 
                    class="cursor-pointer hover:bg-gray-700 px-3 py-1.5 rounded text-xs flex items-center gap-1 transition-colors"
                    @click="convertToSlot('edit')"
                  >
                      <HighlightOutlined />
                      <span>设为编辑块</span>
                  </div>
                  <div class="w-[1px] h-4 bg-gray-600 mx-1"></div>
                  <div 
                    class="cursor-pointer hover:bg-gray-700 px-3 py-1.5 rounded text-xs flex items-center gap-1 transition-colors"
                    @click="convertToSlot('variable')"
                  >
                      <HighlightOutlined />
                      <span>新增变量块</span>
                  </div>
               </div>
            </div>
        </Teleport>
      </div>
    
    </div>
</template>

<style scoped>
/* Editor Styles */
.prose {
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap; 
  word-break: break-word;
}

/* Markdown Styles */
:deep(.md-heading) {
  color: #00b96b; 
  font-weight: 600;
}

:deep(.md-bold) {
  font-weight: 700;
  color: #262626;
}

:deep(.md-italic) {
  font-style: italic;
}

:deep(.md-list-item) {
  color: #00b96b;
}

:deep(.md-code) {
  background-color: #f5f5f5;
  border-radius: 4px;
  padding: 0 4px;
  color: #eb2f96;
  font-family: monospace;
}

/* Slot Tag Style */
:deep(.prompt-slot-tag) {
  display: inline-flex;
  vertical-align: baseline;
  align-items: center;
  margin: 0 2px;
}

:deep(.slot-input) {
  border: 1px solid #d3adf7;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 13px;
  background-color: #f3f0ff;
  color: #722ed1;
  outline: none;
  min-width: 60px; /* Base width */
  max-width: 300px;
  transition: all 0.2s;
  font-weight: 500;
}

:deep(.slot-input:focus) {
  border-color: #722ed1;
  background-color: #fff;
  box-shadow: 0 0 0 2px rgba(114, 46, 209, 0.2);
}

/* Variable Block Distinct Style (Blue) */
:deep(.slot-input.slot-variable) {
  border-color: #91caff;
  background-color: #e6f7ff;
  color: #1890ff;
}

:deep(.slot-input.slot-variable:focus) {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

:deep(.slot-input::placeholder) {
  color: #b37feb;
  opacity: 0.8;
}

.animate-fade-in-up {
  animation: fadeInUp 0.2s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
