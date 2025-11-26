import { defineStore } from 'pinia';
import { ref } from 'vue';
import { sendMessageStream, type Attachment, type Message } from '../api/chat';
export type { Attachment, Message };


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
3. 第三项

## 表格示例
| 功能 | 状态 | 说明 |
| :--- | :---: | ---: |
| Markdown | ✅ | 支持标准语法 |
| 代码高亮 | ✅ | 支持多种语言 |
| 移动端适配 | ✅ | 响应式设计 |`,
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
            aiModel: {
                id: 'gpt-4',
                modelName: 'GPT-4',
                thinkingMode: 'smart',
            },
            tokenUsage: {
                currentTokens: 15,
                totalTokens: 150,
            },
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
            aiModel: {
                id: 'gpt-4',
                modelName: 'GPT-4',
                thinkingMode: 'smart',
            },
            tokenUsage: {
                currentTokens: 500,
                totalTokens: 650,
            },
            callingTools: [
                {
                    name: 'Code Search',
                    description: 'Searching for code examples',
                    functionName: 'search_code',
                }
            ],
            quoteSearchLinks: [
                {
                    url: 'https://vuejs.org/guide/introduction.html',
                    title: 'Vue.js Documentation',
                    content: 'Vue.js is a progressive framework for building user interfaces.',
                    highlight: ['Vue.js', 'progressive framework'],
                }
            ]
        },
    ],
};

export const useChatStore = defineStore('chat', () => {
    const messages = ref<Message[]>([]);
    const isLoading = ref(false);
    const isThinking = ref(false);
    const thinkingMode = ref<'smart' | 'deep' | 'quick'>('smart');

    const currentChatId = ref<string | null>(null);

    const historyItems = ref<{ id: string; title: string }[]>([]);

    // Mock history data
    const mockHistoryItems = [
        { id: '1', title: '如何使用 Vue3 Composition API' },
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

    const updateLastMessage = (data: { content?: string; reasoning_content?: string }) => {
        const lastMsg = messages.value[messages.value.length - 1];
        if (lastMsg && lastMsg.role === 'assistant') {
            if (data.content) {
                lastMsg.content += data.content;
            }
            if (data.reasoning_content) {
                lastMsg.reasoning_content = (lastMsg.reasoning_content || '') + data.reasoning_content;
            }
        }
    };

    const loadChatHistory = (chatId: string) => {
        currentChatId.value = chatId;
        const history = fakeChats[chatId];
        if (history) {
            messages.value = JSON.parse(JSON.stringify(history)); // Deep clone
        } else {
            messages.value = [];
        }
    };

    const clearMessages = () => {
        messages.value = [];
        currentChatId.value = null;
    };

    const deleteChat = (id: string) => {
        historyItems.value = historyItems.value.filter(item => item.id !== id);
        if (currentChatId.value === id) {
            clearMessages();
        }
    };

    const renameChat = (id: string, newTitle: string) => {
        const item = historyItems.value.find(item => item.id === id);
        if (item) {
            item.title = newTitle;
        }
    };

    const abortController = ref<AbortController | null>(null);

    const stopGeneration = () => {
        if (abortController.value) {
            abortController.value.abort();
            abortController.value = null;
            isLoading.value = false;
            isThinking.value = false;

            // Update last message status
            const lastMsg = messages.value[messages.value.length - 1];
            if (lastMsg && lastMsg.role === 'assistant') {
                lastMsg.isStreaming = false;
                lastMsg.content += ' [已停止生成]';
            }
        }
    };

    const sendMessage = async (content: string, attachments: Attachment[] = [], quote?: { quoteId: string; quoteContent: string }) => {
        // Auto-create session if not exists
        if (!currentChatId.value) {
            const newId = Date.now().toString();
            currentChatId.value = newId;

            // Generate title from first 20 chars of content
            const title = content.slice(0, 20) + (content.length > 20 ? '...' : '');

            // Add to history
            historyItems.value.unshift({
                id: newId,
                title: title
            });
        }

        // User message
        const userMessage: Message = {
            id: Date.now().toString(),
            role: 'user',
            content,
            timestamp: Date.now(),
            attachments,
            quoteId: quote?.quoteId,
            quoteContent: quote?.quoteContent,
            aiModel: {
                id: 'gpt-4', // Default mock model
                modelName: 'GPT-4',
                thinkingMode: thinkingMode.value,
            }
        };

        // User message
        addMessage(userMessage);

        isLoading.value = true;
        isThinking.value = true;

        // Create new AbortController
        abortController.value = new AbortController();
        const signal = abortController.value.signal;

        try {
            isThinking.value = false;

            // Simulate AI response
            const responseId = (Date.now() + 1).toString();
            addMessage({
                id: responseId,
                role: 'assistant',
                content: '',
                timestamp: Date.now(),
                isStreaming: true,
                aiModel: {
                    id: 'gpt-4',
                    modelName: 'GPT-4',
                    thinkingMode: thinkingMode.value,
                }
            });

            await sendMessageStream(
                {
                    history: messages.value.slice(0, -2), // Exclude current user msg and pending assistant msg
                    curMessage: userMessage,
                },
                (data) => {
                    updateLastMessage(data);
                },
                signal
            );

            const lastMsg = messages.value[messages.value.length - 1];
            if (lastMsg) {
                lastMsg.isStreaming = false;
            }
        } catch (error) {
            if (error instanceof DOMException && error.name === 'AbortError') {
                console.log('Generation stopped by user');
            } else {
                console.error('Generation error:', error);
                const lastMsg = messages.value[messages.value.length - 1];
                if (lastMsg && lastMsg.role === 'assistant') {
                    lastMsg.isStreaming = false;
                    lastMsg.content += '\n[生成出错]';
                }
            }
        } finally {
            isLoading.value = false;
            isThinking.value = false;
            abortController.value = null;
        }
    };

    const favorites = ref<string[]>([]);

    const toggleFavorite = (chatId: string) => {
        const index = favorites.value.indexOf(chatId);
        if (index === -1) {
            favorites.value.push(chatId);
        } else {
            favorites.value.splice(index, 1);
        }
    };

    const isFavorite = (chatId: string) => {
        return favorites.value.includes(chatId);
    };

    return {
        messages,
        currentChatId,
        historyItems,
        favorites,
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
        toggleFavorite,
        isFavorite,
        stopGeneration,
    };
});
