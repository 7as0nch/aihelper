import { getConfig } from '@/config';
import type { IChatProvider } from './types';
import { OpenAIProvider } from './OpenAIProvider';
import { DeepSeekProvider } from './DeepSeekProvider';

export function createChatProvider(model?: string): IChatProvider {
    const currentModel = model || getConfig('VITE_OPENAI_MODEL', 'gpt-3.5-turbo');

    if (currentModel.includes('deepseek')) {
        return new DeepSeekProvider();
    }

    // Default to OpenAI provider (works for GPT, Claude via proxy, Gemini via proxy)
    return new OpenAIProvider();
}
