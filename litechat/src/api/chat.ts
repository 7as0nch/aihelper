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
    searchByWeb?: boolean; // 是否联网搜索。
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
    role: 'user' | 'assistant' | 'human'; // 用户，ai，人工客服
    aiModel?: AIModel; // 选择的AI模型。
    content: string;
    reasoningContent?: string; // 深度思考内容。
    timestamp: string | number;

    quoteId?: string;
    quoteContent?: string;
    quoteSearchLinks?: QuoteSearchLink[]; // 联网搜索引用文章链接和内容。
    tokenUsage?: { // token消耗累计。
        currentTokens: number; // 当前消耗
        totalTokens: number; // 累计消耗
    };
    callingTools?: CallingTool[]; // agent工具调用情况。
    attachments?: Attachment[]; // 生成或者发送的附件/图片。
    isStreaming?: boolean; // 是否正在流式输出
}

export interface SendMessageParams {
    history: Message[]; // 历史消息
    curMessage: Message; // 当前消息
    curSessionID: string; // 当前会话ID
    needTODOPlan?: 'smart' | 'need' | 'no'; // 是否需要生成待办计划。
}

export function sendMessage(data: SendMessageParams) {
    return request({
        url: '/chat/send',
        method: 'post',
        data,
    });
}

import { getConfig } from '@/config';

export interface Session {
    id: string;
    title: string;
    updateTime: number;
}

export interface HistoryListReply {
    sessions: Session[];
    total: number;
}

export interface HistoryRequest {
    // id: string;
    page: number;
    pageSize: number;
}

export interface MessageListReply {
    messages: Message[];
    total: number;
}

export const chatApi = {
    async getHistoryList(params: HistoryRequest): Promise<HistoryListReply> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') {
            return { sessions: mockHistoryList, total: mockHistoryList.length };
        }
        if (aiType === 'frontend') {
            // Frontend mode: Load from local storage
            const saved = localStorage.getItem('litechat_history');
            const allSessions = saved ? JSON.parse(saved) : [];
            // Implement basic pagination for frontend mode
            const start = (params.page - 1) * params.pageSize;
            const end = start + params.pageSize;
            return {
                sessions: allSessions.slice(start, end),
                total: allSessions.length
            };
        }
        // 检查是否登录，登录后才能请求。
        const token = getToken();
        if (!token) {
            return { sessions: [], total: 0 };
        }
        // Backend mode: Call API
        return await request<HistoryListReply>({
            url: '/chat/history',
            method: 'get',
            params
        });
    },

    async getHistoryMsg(id: string): Promise<Message[]> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo') {
            return mockChats[id] || [];
        }
        if (aiType === 'frontend') {
            // Frontend mode: Load from local storage
            const saved = localStorage.getItem(`litechat_msg_${id}`);
            return saved ? JSON.parse(saved) : [];
        }
        // Backend mode: Call API
        const resp = await request<MessageListReply>({ url: `/chat/history/${id}`, method: 'get' });
        return resp.messages;
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
    async saveChat(id: string, title: string, messages: Message[]): Promise<Session> {
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
        } else if (aiType === 'backend') {
            // Backend mode: Call API
            return request<Session>({ url: `/chat/session`, method: 'post', data: { title, messages } });
        }
        return {
            id,
            title,
            updateTime: Date.now()
        };
    },

    async getExtButtons(): Promise<{ id: string; name: string; api: string; desc: string; icon: string }[]> {
        const aiType = getConfig('VITE_AI_TYPE');
        if (aiType === 'demo' || aiType === 'frontend') {
            // return [
            //     { id: 'summary', name: '语音总结', api: '/api/summary', desc: '生成语音总结', icon: 'Sparkles' },
            //     { id: 'image', name: '生成图片', api: '/api/image', desc: '根据内容生成图片', icon: 'Image' }
            // ];
        }
        // Backend mode
        // return request({ url: '/chat/ext', method: 'get' });
        return [
            { id: 'summary', name: '语音总结', api: '/api/summary', desc: '生成语音总结', icon: 'Sparkles' },
            { id: 'image', name: '生成图片', api: '/api/image', desc: '根据内容生成图片', icon: 'Image' }
        ];
    }
};



import { getToken } from '@/utils/cookie';
import { mockChats, mockHistoryList, streamDemo } from './demo-data';

async function streamBackend(
    data: SendMessageParams,
    onChunk: (data: Partial<Message>) => void,
    signal?: AbortSignal
) {
    const baseURL = getConfig('VITE_BASE_URL', '/');
    // Remove trailing slash if present to avoid double slashes
    const cleanBaseURL = baseURL.endsWith('/') ? baseURL.slice(0, -1) : baseURL;
    console.log('baseURL', cleanBaseURL);
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
    let buffer = '';
    let currentEvent = 'delta'; // 默认事件类型

    while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        
        const lines = buffer.split('\n');
        // Keep the last partial line in the buffer
        buffer = lines.pop() || '';

        for (const line of lines) {
            const trimmedLine = line.trim();
            if (!trimmedLine) continue;

            // 处理事件类型行
            if (trimmedLine.startsWith('event: ')) {
                currentEvent = trimmedLine.slice(7).trim();
                continue;
            }

            // 处理数据行
            if (trimmedLine.startsWith('data: ')) {
                const dataStr = trimmedLine.slice(6).trim();
                
                // 处理完成事件
                if (currentEvent === 'done' || dataStr === '[DONE]') return;

                try {
                    const json = JSON.parse(dataStr);
                    
                    // 如果是错误事件，抛出异常或进行错误处理
                    if (currentEvent === 'error') {
                        throw new Error(json.message || 'Backend stream error');
                    }

                    // 正常增量数据处理
                    onChunk(json);
                } catch (e) {
                    console.error('Failed to parse SSE data:', dataStr, e);
                    // 容错处理：如果解析失败且不是已知结构，尝试当作纯文本
                    if (currentEvent !== 'error') {
                        onChunk({ content: dataStr });
                    }
                }
                
                // 注意：根据规范，event 会持续到下一个 event 出现，这里可以不重置
            }
        }
    }
}

import { createChatProvider } from './providers/factory';

export async function sendMessageStream(
    data: SendMessageParams,
    onChunk: (data: Partial<Message>) => void,
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
