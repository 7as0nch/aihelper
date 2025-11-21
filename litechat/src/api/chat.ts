import request from '@/utils/request';

export function sendMessage(data: any) {
    return request({
        url: '/chat/send',
        method: 'post',
        data,
    });
}

export async function sendMessageStream(
    content: string,
    onChunk: (text: string) => void,
    signal?: AbortSignal
): Promise<void> {
    // Simulate network delay
    await new Promise((resolve, reject) => {
        const timeout = setTimeout(resolve, 1500);
        if (signal) {
            signal.addEventListener('abort', () => {
                clearTimeout(timeout);
                reject(new DOMException('Aborted', 'AbortError'));
            });
        }
    });

    // Mock streaming response
    const mockResponse = `Here is a markdown response for: "${content}"\n\n## Features\n- **Markdown** support\n- *Streaming* output\n- Code blocks:\n\`\`\`typescript\nconst x = 1;\n\`\`\``;

    const chars = mockResponse.split('');
    let currentText = '';

    for (const char of chars) {
        if (signal?.aborted) {
            throw new DOMException('Aborted', 'AbortError');
        }

        await new Promise(resolve => setTimeout(resolve, 30));
        currentText += char;
        onChunk(currentText);
    }
}
