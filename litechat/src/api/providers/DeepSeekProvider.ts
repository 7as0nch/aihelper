import { OpenAIProvider } from './OpenAIProvider';

export class DeepSeekProvider extends OpenAIProvider {
    // DeepSeek is fully compatible with OpenAI API, including reasoningContent in delta.
    // We inherit the implementation for now.
    // If DeepSeek adds specific parameters or different behavior, we override here.
}
