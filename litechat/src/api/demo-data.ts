import type { Message, Session } from './chat';

export const mockHistoryList: Session[] = [
    { id: '1', title: '如何使用 Vue3 Composition API', updateTime: Date.now() },
    { id: '2', title: 'Demo Chat 2', updateTime: Date.now() - 86400000 }
];

export const mockChats: Record<string, Message[]> = {
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

export async function streamDemo(
    onChunk: (data: Partial<Message>) => void,
    signal?: AbortSignal
) {
    const mockResponse = "This is a mock response from Demo Mode. I am simulating a streaming response.";
    const mockReasoning = "I am thinking about how to simulate this response...";

    // Simulate reasoning first
    onChunk({ reasoningContent: "" });
    const reasoningChars = mockReasoning.split('');
    for (const char of reasoningChars) {
        if (signal?.aborted) return;
        await new Promise(resolve => setTimeout(resolve, 50));
        onChunk({ reasoningContent: char });
    }
    onChunk({ reasoningContent: "\n" }); // End reasoning

    // Simulate content
    const chars = mockResponse.split('');
    for (const char of chars) {
        if (signal?.aborted) return;
        await new Promise(resolve => setTimeout(resolve, 30));
        onChunk({ content: char });
    }
}
