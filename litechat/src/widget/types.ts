/*
 * @Author: chengjiang
 * @Date: 2025-11-27 22:22:06
 * @Description: 
 */
import { Ref, ShallowRef } from 'vue';
import type { RuntimeConfig } from '../config';

export interface InitOptions {
    config?: Partial<RuntimeConfig>;
    defaultOpen?: boolean;
    defaultShow?: boolean;
    containerId?: string;
}
export interface AiChat {
    isMounted?: Ref<boolean, boolean>;
    widgetInstance?: ShallowRef<AiChat | null, AiChat | null>;
    unmount: () => void; // 卸载浮窗
    open: () => void; // 打开浮窗
    close: () => void; // 关闭浮窗
    toggle: () => void;
    show: () => void; // 显示浮球
    hide: () => void; // 隐藏浮窗
}
