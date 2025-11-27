import { createWidgetShell } from '@/widget/shadow-wrapper';
import { createIframeRenderer } from '@/widget/iframe-renderer';
import { setRuntimeConfig } from '../config';
import type { InitOptions, AiChat } from './types';

export type { InitOptions } from './types';

export function initAiChat(options: InitOptions = {}): AiChat {
    // 1. Set global config
    if (options.config) {
        setRuntimeConfig(options.config);
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
