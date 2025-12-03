import { getConfig } from '@/config';
import type { IChatProvider } from './types';
import type { SendMessageParams } from '../chat';

export class OpenAIProvider implements IChatProvider {
    async streamChat(
        params: SendMessageParams,
        onChunk: (data: { content?: string; reasoningContent?: string }) => void,
        signal?: AbortSignal
    ): Promise<void> {
        const apiKey = getConfig('VITE_OPENAI_API_KEY');
        const baseURL = getConfig('VITE_OPENAI_BASE_URL');
        const model = getConfig('VITE_OPENAI_MODEL', 'gpt-3.5-turbo');

        if (!apiKey) {
            throw new Error('OpenAI API Key is missing');
        }

        // Convert messages to OpenAI format
        const messages = [
            ...params.history.map(m => ({
                role: m.role,
                content: m.content
            })),
            {
                role: params.curMessage.role,
                content: params.curMessage.content
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
                                reasoningContent: delta.reasoningContent || undefined
                            });
                        }
                    } catch (e) {
                        console.warn('Failed to parse SSE message:', line);
                    }
                }
            }
        }
    }
}
