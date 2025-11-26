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
    isStreaming?: boolean;
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
                { id: '1', title: 'Demo Chat 1', updateTime: Date.now() },
                { id: '2', title: 'Demo Chat 2', updateTime: Date.now() - 86400000 }
            ];
        }
        if (aiType === 'backend') {
            // Pure frontend mode might store in local storage
            const saved = localStorage.getItem('litechat_history');
            return saved ? JSON.parse(saved) : [];
        }
        return request({ url: '/chat/history', method: 'get' });
    },

    async getHistoryMsg(id: string): Promise<Message[]> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') {
            return [
                { id: 'msg1', role: 'user', content: 'Hello', timestamp: Date.now() - 10000 },
                { id: 'msg2', role: 'assistant', content: 'Hi! This is a demo.', timestamp: Date.now() }
            ];
        }
        if (aiType === 'backend') {
            // Pure frontend mode might store in local storage
            const saved = localStorage.getItem(`litechat_msg_${id}`);
            return saved ? JSON.parse(saved) : [];
        }
        return request({ url: `/chat/history/${id}`, method: 'get' });
    },

    async deleteChat(id: string): Promise<void> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') return;
        if (aiType === 'backend') {
            const saved = localStorage.getItem('litechat_history');
            if (saved) {
                const list = JSON.parse(saved).filter((item: any) => item.id !== id);
                localStorage.setItem('litechat_history', JSON.stringify(list));
            }
            localStorage.removeItem(`litechat_msg_${id}`);
            return;
        }
        return request({ url: `/chat/history/${id}`, method: 'delete' });
    },

    async renameChat(id: string, title: string): Promise<void> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') return;
        if (aiType === 'backend') {
            const saved = localStorage.getItem('litechat_history');
            if (saved) {
                const list = JSON.parse(saved).map((item: any) => item.id === id ? { ...item, title } : item);
                localStorage.setItem('litechat_history', JSON.stringify(list));
            }
            return;
        }
        return request({ url: `/chat/history/${id}/rename`, method: 'post', data: { title } });
    },

    async sendMessage(messages: Message[], onProgress?: (content: string) => void): Promise<string> {
        const apiKey = getConfig('VITE_OPENAI_API_KEY');
        const baseURL = getConfig('VITE_OPENAI_BASE_URL');
        const model = getConfig('VITE_OPENAI_MODEL', 'gpt-3.5-turbo');

        if (!apiKey) {
            throw new Error('OpenAI API Key is missing');
        }

        // Silence unused variable warnings for now
        void messages;
        void onProgress;
        void baseURL;
        void model;

        // Placeholder
        return "Message sent";
    }
};



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
            // Add auth token if needed
            // 'Authorization': `Bearer ${token}`
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
