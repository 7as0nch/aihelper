import request from '@/utils/request';

export interface Attachment {
    id: string;
    type: 'image' | 'file';
    name: string;
    url?: string;
}

export interface AIModel {
    id: string;
    modelName: string;
    thinkingMode: 'smart' | 'deep' | 'quick'; // 思考类型，深度思考，快速思考，智能思考。
}

export interface CallingTool {
    name: string;
    description: string;
    functionName: string;
}

export interface QuoteSearchLink {
    url: string;
    title: string;
    content: string;
    highlight: string[];
}

export interface Message {
    id: string;
    role: 'user' | 'assistant';
    aiModel?: AIModel; // 选择的AI模型。
    content: string;
    reasoning_content?: string; // 深度思考内容。
    timestamp: number;

    quoteId?: string;
    quoteContent?: string;
    quoteSearchLinks?: QuoteSearchLink[]; // 联网搜索引用文章链接和内容。
    tokenUsage?: { // token消耗累计。
        currentTokens: number;
        totalTokens: number;
    };
    callingTools?: CallingTool[]; // agent工具调用情况。
    attachments?: Attachment[]; // 生成或者发送的附件/图片。
}

export interface SendMessageParams {
    history: Message[]; // 历史消息
    curMessage: Message; // 当前消息
}

export function sendMessage(data: SendMessageParams) {
    return request({
        url: '/chat/send',
        method: 'post',
        data,
    });
}

import { getConfig } from '@/config';

export const chatApi = {
    async getHistoryList(): Promise<any[]> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') {
            return [
                { id: '1', title: '如何使用 Vue3 Composition API', updateTime: Date.now() },
                { id: '2', title: 'Demo Chat 2', updateTime: Date.now() - 86400000 }
            ];
        }
        if (aiType === 'frontend') {
            // Frontend mode: Load from local storage
            const saved = localStorage.getItem('litechat_history');
            return saved ? JSON.parse(saved) : [];
        }
        // 检查是否登录，登录后才能请求。
        const token = getToken();
        if (!token) {
            return [];
        }
        // Backend mode: Call API
        return request<any[]>({ url: '/chat/history', method: 'get' });
    },

    async getHistoryMsg(id: string): Promise<Message[]> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') {
            const mockChats: Record<string, Message[]> = {
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
            return mockChats[id] || [];
        }
        if (aiType === 'frontend') {
            // Frontend mode: Load from local storage
            const saved = localStorage.getItem(`litechat_msg_${id}`);
            return saved ? JSON.parse(saved) : [];
        }
        // Backend mode: Call API
        return request<Message[]>({ url: `/chat/history/${id}`, method: 'get' });
    },

    async deleteChat(id: string): Promise<void> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') return;
        if (aiType === 'frontend') {
            // Frontend mode: Delete from local storage
            const saved = localStorage.getItem('litechat_history');
            if (saved) {
                const list = JSON.parse(saved).filter((item: any) => item.id !== id);
                localStorage.setItem('litechat_history', JSON.stringify(list));
            }
            localStorage.removeItem(`litechat_msg_${id}`);
            return;
        }
        // Backend mode: Call API
        return request<void>({ url: `/chat/history/${id}`, method: 'delete' });
    },

    async renameChat(id: string, title: string): Promise<void> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') return;
        if (aiType === 'frontend') {
            // Frontend mode: Update local storage
            const saved = localStorage.getItem('litechat_history');
            if (saved) {
                const list = JSON.parse(saved).map((item: any) => item.id === id ? { ...item, title } : item);
                localStorage.setItem('litechat_history', JSON.stringify(list));
            }
            return;
        }
        // Backend mode: Call API
        return request<void>({ url: `/chat/history/${id}/rename`, method: 'post', data: { title } });
    },

    // Helper to save chat history and messages in frontend mode
    async saveChat(id: string, title: string, messages: Message[]): Promise<void> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'frontend') {
            // 1. Update history list
            const historyStr = localStorage.getItem('litechat_history');
            let history = historyStr ? JSON.parse(historyStr) : [];

            const existingIndex = history.findIndex((h: any) => h.id === id);
            if (existingIndex > -1) {
                history[existingIndex].updateTime = Date.now();
                // Update title if it's the first message or explicitly changed (handled by rename)
                // Here we just update time. Title is usually set on creation.
            } else {
                history.unshift({
                    id,
                    title,
                    updateTime: Date.now()
                });
            }
            localStorage.setItem('litechat_history', JSON.stringify(history));

            // 2. Save messages
            localStorage.setItem(`litechat_msg_${id}`, JSON.stringify(messages));
        }
    }
};



import { getToken } from '@/utils/cookie';

async function streamBackend(
    data: SendMessageParams,
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
) {
    const baseURL = getConfig('VITE_BASE_URL', '/');
    // Remove trailing slash if present to avoid double slashes
    const cleanBaseURL = baseURL.endsWith('/') ? baseURL.slice(0, -1) : baseURL;

    const response = await fetch(`${cleanBaseURL}/chat/send`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            ...(getToken() ? { 'Authorization': `Bearer ${getToken()}` } : {})
        },
        body: JSON.stringify(data),
        signal
    });

    if (!response.ok) {
        throw new Error('Backend API request failed');
    }

    if (!response.body) {
        throw new Error('Response body is empty');
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();

    while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        // Assuming backend returns raw text chunks or SSE. 
        // If SSE, we need parsing logic similar to OpenAI.
        // For now, let's assume raw text stream or simple SSE.
        // If it's SSE, we should parse 'data: '.
        // Let's assume it's a simple text stream for now, or adapt if user specified.
        // User said "backend mode calls api interface".
        // Let's assume standard SSE for backend too as it's common for chat.

        // Simple pass-through for now, but if it's SSE, we might need to parse.
        // Let's try to parse as SSE if it starts with 'data:', otherwise pass through.

        if (chunk.startsWith('data: ')) {
            const lines = chunk.split('\n');
            for (const line of lines) {
                if (line.startsWith('data: ')) {
                    const dataStr = line.slice(6).trim();
                    if (dataStr === '[DONE]') return;
                    try {
                        // Try to parse as JSON if possible, or just use string
                        // If backend sends JSON: data: {"content": "..."}
                        // If backend sends string: data: ...
                        if (dataStr.startsWith('{')) {
                            const json = JSON.parse(dataStr);
                            onChunk({
                                content: json.content,
                                reasoning_content: json.reasoning_content
                            });
                        } else {
                            onChunk({ content: dataStr });
                        }
                    } catch {
                        onChunk({ content: dataStr });
                    }
                }
            }
        } else {
            onChunk({ content: chunk });
        }
    }
}

async function streamDemo(
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
) {
    const mockResponse = "This is a mock response from Demo Mode. I am simulating a streaming response.";
    const mockReasoning = "I am thinking about how to simulate this response...";

    // Simulate reasoning first
    onChunk({ reasoning_content: "" });
    const reasoningChars = mockReasoning.split('');
    for (const char of reasoningChars) {
        if (signal?.aborted) return;
        await new Promise(resolve => setTimeout(resolve, 50));
        onChunk({ reasoning_content: char });
    }
    onChunk({ reasoning_content: "\n" }); // End reasoning

    // Simulate content
    const chars = mockResponse.split('');
    for (const char of chars) {
        if (signal?.aborted) return;
        await new Promise(resolve => setTimeout(resolve, 30));
        onChunk({ content: char });
    }
}

import { createChatProvider } from './providers/factory';

export async function sendMessageStream(
    data: SendMessageParams,
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
): Promise<void> {
    const aiType = getConfig('VITE_AI_TYPE');

    if (aiType === 'demo') {
        await streamDemo(onChunk, signal);
    } else if (aiType === 'frontend') {
        // Pure Frontend Mode (Direct OpenAI)
        // Use the Provider Factory to select the correct provider based on model
        const provider = createChatProvider();
        await provider.streamChat(data, onChunk, signal);
    } else {
        // Frontend Mode (Backend API) - Default
        await streamBackend(data, onChunk, signal);
    }
}
