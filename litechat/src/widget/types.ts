import type { RuntimeConfig } from '../config';

export interface InitOptions {
    config?: Partial<RuntimeConfig>;
    defaultOpen?: boolean;
    defaultShow?: boolean;
    containerId?: string;
}
