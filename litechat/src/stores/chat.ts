import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface Attachment {
    id: string;
    type: 'image' | 'file';
    name: string;
    url?: string;
}

export interface Message {
    id: string;
    role: 'user' | 'assistant';
    content: string;
    timestamp: number;
    isStreaming?: boolean;
    quoteId?: string;
    quoteContent?: string;
    attachments?: Attachment[];
}

// Fake chat history data
const fakeChats: Record<string, Message[]> = {
    '1': [
        {
            id: '1-1',
            role: 'user',
            content: '是否',
            timestamp: Date.now() - 120000,
        },
        {
            id: '1-2',
            role: 'assistant',
            content: `Here is a markdown response for: "是否"

## Features
- **Markdown** support
- *Streaming* output
- Code blocks:

\`\`\`typescript
const x = 1;
\`\`\``,
            timestamp: Date.now() - 118000,
        },
        {
            id: '1-3',
            role: 'user',
            content: '能展示一些图片和文件吗？',
            timestamp: Date.now() - 100000,
            attachments: [
                {
                    id: 'img-1',
                    type: 'image',
                    name: 'screenshot.png',
                    url: 'https://picsum.photos/400/300',
                },
            ],
        },
        {
            id: '1-4',
            role: 'assistant',
            content: `当然可以！这里有一些示例：

## 图片展示
我可以展示各种图片。

## 代码示例
\`\`\`python
def hello_world():
    print("Hello, World!")
    return True
\`\`\`

## Markdown 功能
- **粗体文本**
- *斜体文本*
- \`代码片段\`
- [链接示例](https://example.com)

## 列表示例
1. 第一项
2. 第二项
3. 第三项`,
            timestamp: Date.now() - 95000,
            attachments: [
                {
                    id: 'img-2',
                    type: 'image',
                    name: 'diagram.png',
                    url: 'https://picsum.photos/500/350',
                },
                {
                    id: 'file-1',
                    type: 'file',
                    name: 'document.pdf',
                },
            ],
        },
        {
            id: '1-5',
            role: 'user',
            content: '展示一些复杂的代码块',
            timestamp: Date.now() - 60000,
        },
        {
            id: '1-6',
            role: 'assistant',
            content: `好的，这里有一些复杂的代码示例：

## Vue 3 组件示例
\`\`\`vue
<script setup lang="ts">
import { ref, computed } from 'vue';

const count = ref(0);
const doubled = computed(() => count.value * 2);

function increment() {
  count.value++;
}
</script>

<template>
  <div>
    <p>Count: {{ count }}</p>
    <p>Doubled: {{ doubled }}</p>
    <button @click="increment">Increment</button>
  </div>
</template>
\`\`\`

## TypeScript 高级类型
\`\`\`typescript
interface User {
  id: number;
  name: string;
  email: string;
}

type PartialUser = Partial<User>;
type RequiredUser = Required<User>;

function updateUser<T extends User>(user: T, updates: Partial<T>): T {
  return { ...user, ...updates };
}
\`\`\`

## JavaScript 异步处理
\`\`\`javascript
async function fetchData(url) {
  try {
    const response = await fetch(url);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching data:', error);
    throw error;
  }
}
\`\`\``,
            timestamp: Date.now() - 55000,
        },
    ],
};

export const useChatStore = defineStore('chat', () => {
    const messages = ref<Message[]>([]);
    const isLoading = ref(false);
    const isThinking = ref(false);
    const thinkingMode = ref<'smart' | 'deep' | 'quick'>('smart');

    const historyItems = ref<{ id: string; title: string }[]>([]);

    // Mock history data
    const mockHistoryItems = [
        { id: '1', title: '如何使用 Vue3 Composition API' },
        { id: '2', title: 'TailwindCSS 最佳实践' },
        { id: '3', title: '2025年 AI 发展趋势报告' },
        { id: '4', title: 'TypeScript 高级类型解析' },
        { id: '5', title: 'Vite 构建性能优化指南' },
        { id: '6', title: 'Pinia 状态管理最佳实践' },
        { id: '7', title: 'Vue Router 路由守卫详解' },
        { id: '8', title: '前端工程化落地实践' },
    ];

    const fetchHistoryList = async () => {
        // In a real app, we would check auth here or rely on API to return 401
        // But since we are mocking, we'll just simulate a fetch
        await new Promise(resolve => setTimeout(resolve, 500));
        historyItems.value = mockHistoryItems;
    };

    const clearHistoryList = () => {
        historyItems.value = [];
    };

    const addMessage = (message: Message) => {
        messages.value.push(message);
    };

    const updateLastMessage = (content: string) => {
        const lastMsg = messages.value[messages.value.length - 1];
        if (lastMsg && lastMsg.role === 'assistant') {
            lastMsg.content = content;
        }
    };

    const loadChatHistory = (chatId: string) => {
        const history = fakeChats[chatId];
        if (history) {
            messages.value = JSON.parse(JSON.stringify(history)); // Deep clone
        } else {
            messages.value = [];
        }
    };

    const clearMessages = () => {
        messages.value = [];
    };

    const deleteChat = (id: string) => {
        historyItems.value = historyItems.value.filter(item => item.id !== id);
        if (messages.value.length > 0) {
            // Ideally we should check if the current chat is the one being deleted
            // For now, we just clear messages if we are "in" a chat context that gets deleted
            // But since we don't track currentChatId in store yet, we'll leave this simple
        }
    };

    const renameChat = (id: string, newTitle: string) => {
        const item = historyItems.value.find(item => item.id === id);
        if (item) {
            item.title = newTitle;
        }
    };

    const sendMessage = async (content: string, attachments: Attachment[] = [], quote?: { quoteId: string; quoteContent: string }) => {
        // User message
        addMessage({
            id: Date.now().toString(),
            role: 'user',
            content,
            timestamp: Date.now(),
            attachments,
            quoteId: quote?.quoteId,
            quoteContent: quote?.quoteContent,
        });

        isLoading.value = true;
        isThinking.value = true;

        // Simulate Thinking delay
        await new Promise(resolve => setTimeout(resolve, 1500));

        isThinking.value = false;

        // Simulate AI response
        const responseId = (Date.now() + 1).toString();
        addMessage({
            id: responseId,
            role: 'assistant',
            content: '',
            timestamp: Date.now(),
            isStreaming: true,
        });

        // Mock streaming
        const mockResponse = `Here is a markdown response for: "${content}"\n\n## Features\n- **Markdown** support\n- *Streaming* output\n- Code blocks:\n\`\`\`typescript\nconst x = 1;\n\`\`\``;

        let currentText = '';
        const chars = mockResponse.split('');

        for (const char of chars) {
            await new Promise(resolve => setTimeout(resolve, 30));
            currentText += char;
            updateLastMessage(currentText);
        }

        const lastMsg = messages.value[messages.value.length - 1];
        if (lastMsg) {
            lastMsg.isStreaming = false;
        }
        isLoading.value = false;
    };

    return {
        messages,
        historyItems,
        isLoading,
        isThinking,
        thinkingMode,
        sendMessage,
        loadChatHistory,
        clearMessages,
        deleteChat,
        renameChat,
        fetchHistoryList,
        clearHistoryList,
    };
});
