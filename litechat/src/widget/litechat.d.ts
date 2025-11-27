/*
 * @Author: chengjiang
 * @Date: 2025-11-26 17:00:50
 * @Description:
 */
import { RuntimeConfig } from '../config';
import { AiChat } from './types';
import { initAiChat } from './index';
import { useAiChat } from './vue';

export { initAiChat, useAiChat };
export interface InitOptions {
    config?: Partial<RuntimeConfig>;
    defaultOpen?: boolean;
    defaultShow?: boolean;
    containerId?: string;
}
export type { InitOptions };
export type { AiChat };
