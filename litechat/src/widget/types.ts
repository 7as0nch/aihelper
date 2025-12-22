/*
 * @Author: chengjiang
 * @Date: 2025-11-27 22:22:06
 * @Description: 
 */
import { Ref, ShallowRef } from 'vue';
import type { RuntimeConfig } from '../config';

// Enums
export enum ChatType {
    Bot = 'bot',
    App = 'app'
}

export enum Language {
    EN = 'en',
    CN = 'cn',
    JP = 'jp'
}

// Interfaces based on LiteChat SDK
export interface LiteChatChat {
    appId?: string; // Bot ID or Conversation ID
    type?: ChatType;
    repoId?: string; // For App type
}

export interface LiteChatSetting {
    apiBaseUrl?: string; // e.g. https://api.litechat.cn
    cdnBaseUrlPath?: string;
    language?: Language;
    requestHeader?: Record<string, string>;
    logLevel?: 'debug' | 'info' | 'warn' | 'error';
}

export interface LiteChatAuth {
    type?: 'external' | 'token';
    token?: string;
    onRefreshToken?: (oldToken?: string) => string | Promise<string>;
    refreshToken?: () => string | Promise<string>; // For App type compatibility
}

export interface LiteChatUser {
    id: string;
    name: string;
    avatar: string;
}

export interface LiteChatHeader {
    isNeed?: boolean;
    icon?: string;
    title?: string;
}

export interface LiteChatFooter {
    isNeed?: boolean;
    expressionText?: string;
    linkvars?: any;
}

export interface LiteChatInput {
    isNeed?: boolean;
    placeholder?: string;
    isNeedAudio?: boolean; // default false
    isNeedSendButton?: boolean; // default true
}

export interface LiteChatChatSlot {
    input?: LiteChatInput;
}

export interface LiteChatUI {
    layout?: 'pc' | 'mobile'; // default 'pc'
    isReadonly?: boolean;
    header?: LiteChatHeader;
    footer?: LiteChatFooter;
    chatSlot?: LiteChatChatSlot;
    uploadBtn?: { isNeed: boolean };
}

// Extended InitOptions
export interface InitOptions {
    // Legacy support
    config?: Partial<RuntimeConfig>;
    defaultOpen?: boolean;
    defaultShow?: boolean;
    containerId?: string;

    // LiteChat SDK Style
    chat?: LiteChatChat;
    setting?: LiteChatSetting;
    auth?: LiteChatAuth;
    user?: LiteChatUser;
    ui?: LiteChatUI;

    // Callbacks (simplified for now)
    eventCallbacks?: {
        onInitSuccess?: () => void;
        onThemeChange?: (theme: string) => void;
        // ... others can be added as needed
    };
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
