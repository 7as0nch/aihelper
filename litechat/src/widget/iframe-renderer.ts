import { mountApp } from '../mount';
import type { InitOptions } from './types';
// @ts-ignore
import styleContent from '../style.css?inline';
// import markdownLightCss from 'github-markdown-css/github-markdown-light.css?inline';
// import markdownDarkCss from 'github-markdown-css/github-markdown-dark.css?inline';

export function createIframeRenderer(container: HTMLElement, options: InitOptions) {
    const iframe = document.createElement('iframe');

    // Basic Iframe styles
    iframe.style.border = 'none';
    iframe.style.width = '100%';
    iframe.style.height = '100%';
    iframe.style.display = 'block';

    // CRITICAL: Allow downloads in iframe
    // Without this, file downloads will fail in iframe environments
    iframe.setAttribute('allow', 'downloads');

    // Note: Positioning is now handled by the shell (window-container)

    container.appendChild(iframe);

    let appInstance: any = null;

    const mount = () => {
        const doc = iframe.contentDocument;
        if (!doc) return;

        // 1. Write initial HTML structure
        doc.open();
        doc.write(`
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <title>AI Chat Widget</title>
                <style>
                    body { margin: 0; padding: 0; overflow: hidden; background: transparent; }
                    #app { height: 100vh; width: 100vw; }
                    /* Scrollbar styling */
                    ::-webkit-scrollbar { width: 6px; height: 6px; }
                    ::-webkit-scrollbar-track { background: transparent; }
                    ::-webkit-scrollbar-thumb { background: rgba(0, 0, 0, 0.2); border-radius: 3px; }
                    ::-webkit-scrollbar-thumb:hover { background: rgba(0, 0, 0, 0.3); }
                </style>
            </head>
            <body>
                <div id="app"></div>
            </body>
            </html>
        `);
        doc.close();

        // 2. Mount Vue App when iframe is ready
        iframe.onload = () => {
            const doc = iframe.contentDocument;
            if (!doc) return;

            // Inject compiled styles
            const styleEl = doc.createElement('style');
            styleEl.textContent = styleContent;
            doc.head.appendChild(styleEl);

            // // Inject markdown styles for code blocks and tables
            // const markdownStyleEl = doc.createElement('style');
            // markdownStyleEl.textContent = markdownLightCss + '\n' + markdownDarkCss;
            // doc.head.appendChild(markdownStyleEl);

            // Sync initial dark mode state from parent
            const syncTheme = () => {
                const isDark = document.documentElement.classList.contains('dark');
                if (isDark) {
                    doc.documentElement.classList.add('dark');
                } else {
                    doc.documentElement.classList.remove('dark');
                }
            };

            // Apply initial theme
            syncTheme();

            // Watch for theme changes in parent and sync to iframe
            const observer = new MutationObserver((mutations) => {
                mutations.forEach((mutation) => {
                    if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
                        syncTheme();
                    }
                });
            });

            observer.observe(document.documentElement, {
                attributes: true,
                attributeFilter: ['class']
            });

            // Mount App
            const appContainer = doc.getElementById('app');
            if (appContainer) {
                appInstance = mountApp({
                    config: options.config,
                    container: appContainer,
                    routerMode: 'memory',
                    initialPath: '/chat',
                    styleContainer: doc.head
                });
            }
        };
    };

    const unmount = () => {
        if (appInstance) {
            appInstance.app.unmount();
        }
        iframe.remove();
    };

    return { mount, unmount };
}
