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

async function streamOpenAI(
    data: SendMessageParams,
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
) {
    const apiKey = import.meta.env.VITE_OPENAI_API_KEY;
    const baseURL = import.meta.env.VITE_OPENAI_BASE_URL;
    const model = import.meta.env.VITE_OPENAI_MODEL || 'gpt-3.5-turbo';

    if (!apiKey) {
        throw new Error('OpenAI API Key is missing');
    }

    // Convert messages to OpenAI format
    const messages = [
        ...data.history.map(m => ({
            role: m.role,
            content: m.content
        })),
        {
            role: data.curMessage.role,
            content: data.curMessage.content
        }
    ];

    const response = await fetch(`${baseURL}/chat/completions`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${apiKey}`
        },
        body: JSON.stringify({
            model,
            messages,
            stream: true
        }),
        signal
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.error?.message || 'OpenAI API request failed');
    }

    if (!response.body) {
        throw new Error('Response body is empty');
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';

        for (const line of lines) {
            if (line.trim() === '') continue;
            if (line.trim() === 'data: [DONE]') return;

            if (line.startsWith('data: ')) {
                try {
                    const json = JSON.parse(line.slice(6));
                    const delta = json.choices[0]?.delta;
                    if (delta) {
                        onChunk({
                            content: delta.content || undefined,
                            reasoning_content: delta.reasoning_content || undefined
                        });
                    }
                } catch (e) {
                    console.warn('Failed to parse SSE message:', line);
                }
            }
        }
    }
}

async function streamBackend(
    data: SendMessageParams,
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
) {
    const baseURL = import.meta.env.VITE_BASE_URL || '/';
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

export async function sendMessageStream(
    data: SendMessageParams,
    onChunk: (data: { content?: string; reasoning_content?: string }) => void,
    signal?: AbortSignal
): Promise<void> {
    const aiType = import.meta.env.VITE_AI_TYPE;

    if (aiType === 'frontend') {
        await streamOpenAI(data, onChunk, signal);
    } else {
        await streamBackend(data, onChunk, signal);
    }
}
