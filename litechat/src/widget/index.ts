import { setToken } from '@/utils/cookie';
import { createWidgetShell } from '@/widget/shadow-wrapper';
import { createIframeRenderer } from '@/widget/iframe-renderer';
import { setRuntimeConfig } from '../config';
import type { InitOptions, AiChat } from './types';

export type { InitOptions, AiChat } from './types';
export { useAiChat } from './vue';

export function initAiChat(options: InitOptions = {}): AiChat {
    // 1. Set global config

    // Map LiteChat SDK options to RuntimeConfig
    const mappedConfig: Record<string, any> = {};

    if (options.setting) {
        if (options.setting.apiBaseUrl) mappedConfig.VITE_API_BASE_URL = options.setting.apiBaseUrl;
        // Language handling could be added here if app supports it
    }

    if (options.ui) {
        if (options.ui.header) {
            if (options.ui.header.title) mappedConfig.VITE_APP_TITLE = options.ui.header.title;
            if (options.ui.header.icon) mappedConfig.VITE_APP_LOGO = options.ui.header.icon;
            if (options.ui.header.isNeed !== undefined) mappedConfig.VITE_SHOW_HEADER = String(options.ui.header.isNeed);
        }

        if (options.ui.footer) {
            if (options.ui.footer.isNeed !== undefined) mappedConfig.VITE_SHOW_FOOTER = String(options.ui.footer.isNeed);
            if (options.ui.footer.expressionText) mappedConfig.VITE_FOOTER_TEXT = options.ui.footer.expressionText;
        }

        if (options.ui.chatSlot?.input) {
            if (options.ui.chatSlot.input.placeholder) mappedConfig.VITE_INPUT_PLACEHOLDER = options.ui.chatSlot.input.placeholder;
            if (options.ui.chatSlot.input.isNeedAudio !== undefined) mappedConfig.VITE_SHOW_AUDIO_BTN = String(options.ui.chatSlot.input.isNeedAudio);
            // isNeedSendButton is not currently mapped to a global config but could be if needed
        }

        if (options.ui.uploadBtn) {
            if (options.ui.uploadBtn.isNeed !== undefined) mappedConfig.VITE_SHOW_UPLOAD_BTN = String(options.ui.uploadBtn.isNeed);
        }
    }

    if (options.chat) {
        // Assuming appId maps to the model/bot ID in this system
        if (options.chat.appId) mappedConfig.VITE_OPENAI_MODEL = options.chat.appId;
    }

    // Handle Authentication
    if (options.auth) {
        if (options.auth.token) {
            mappedConfig.VITE_OPENAI_API_KEY = options.auth.token;
            // Also set token in cookie for backend requests
            setToken(options.auth.token);
        }
    }

    // Merge with legacy config (legacy takes precedence if both exist, or mapped serves as base)
    const finalConfig = { ...mappedConfig, ...options.config };

    if (Object.keys(finalConfig).length > 0) {
        setRuntimeConfig(finalConfig);
    }

    // 2. Create container (Shadow DOM host)
    const containerId = options.containerId || 'ai-chat-widget-root';
    let container = document.getElementById(containerId);

    if (!container) {
        container = document.createElement('div');
        container.id = containerId;
        document.body.appendChild(container);
    }

    // 3. Initialize Widget Shell (Shadow DOM, Button, Window Container)
    // Default show is true unless explicitly set to false
    const initialShow = options.defaultShow !== false;
    const { windowContainer, open, close, toggle } = createWidgetShell(container, options.defaultOpen);

    // 4. Initialize Iframe Renderer (The Content)
    // Mount inside the window container provided by the shell
    const { mount, unmount } = createIframeRenderer(windowContainer, options);

    // 5. Mount the app inside the iframe
    mount();

    // Ensure container visibility matches initialShow
    if (container) {
        container.style.display = initialShow ? 'block' : 'none';
    }

    return {
        unmount: () => {
            unmount();
            if (container && container.parentNode) {
                container.parentNode.removeChild(container);
            }
        },
        open,
        close,
        toggle,
        show: () => {
            if (container) container.style.display = 'block';
        },
        hide: () => {
            if (container) container.style.display = 'none';
        }
    };
}

// Auto-init if script tag has data attributes (optional, for simple usage)
// <script src="..." data-auto-init data-config-title="My Chat"></script>
if (typeof window !== 'undefined') {
    const currentScript = document.currentScript as HTMLScriptElement;
    if (currentScript && currentScript.hasAttribute('data-auto-init')) {
        const config: Record<string, string> = {};
        let containerId: string | undefined;
        let defaultOpen: boolean | undefined;
        let defaultShow: boolean | undefined;

        // Parse data-config-* attributes and others
        for (const attr of currentScript.attributes) {
            if (attr.name.startsWith('data-config-')) {
                const key = attr.name.replace('data-config-', '').toUpperCase().replace(/-/g, '_');
                config[`VITE_${key}`] = attr.value;
            } else if (attr.name === 'data-container-id') {
                containerId = attr.value;
            } else if (attr.name === 'data-default-open') {
                defaultOpen = attr.value !== 'false';
            } else if (attr.name === 'data-default-show') {
                defaultShow = attr.value !== 'false';
            }
        }

        const widget = initAiChat({ config, containerId, defaultOpen, defaultShow });

        // Expose instance for auto-init
        (window as any).AiChatWidgetInstance = widget;
    }

    // Expose to window for UMD
    (window as any).AiChatWidget = { initAiChat };
}
