import type { SendMessageParams } from '../chat';

export interface IChatProvider {
    streamChat(
        params: SendMessageParams,
        onChunk: (data: { content?: string; reasoning_content?: string }) => void,
        signal?: AbortSignal
    ): Promise<void>;
}
