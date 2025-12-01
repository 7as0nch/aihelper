import { getConfig } from "@/config";

export interface WidgetShell {
    shadowRoot: ShadowRoot;
    windowContainer: HTMLElement;
    toggleBtn: HTMLElement;
}

export function createWidgetShell(container: HTMLElement, initialOpen: boolean = false): WidgetShell & { open: () => void; close: () => void; toggle: () => void } {
    const shadow = container.attachShadow({ mode: 'open' });

    // Inject styles
    const style = document.createElement('style');
    style.textContent = `
        :host {
            position: fixed;
            z-index: 9999;
            top: 0;
            left: 0;
            width: 0;
            height: 0;
            overflow: visible;
        }
        
        .widget-shell {
            position: fixed;
            z-index: 10000;
            transition: transform 0.1s ease-out; /* Smooth drag */
            touch-action: none; /* Prevent scroll on mobile */
        }

        /* Floating Button */
        .floating-btn {
            width: 60px;
            height: 60px;
            border-radius: 50%;
            background: linear-gradient(135deg, #3b82f6, #2563eb);
            box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            user-select: none;
            transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.2s ease;
            position: relative;
        }
        .floating-btn:hover { transform: scale(1.1); }
        .floating-btn:active { transform: scale(0.95); }
        .floating-btn svg { width: 32px; height: 32px; color: white; }
        .floating-btn img { width: 100%; height: 100%; object-fit: cover; border-radius: 50%; }
        
        .floating-btn.hidden {
            opacity: 0;
            pointer-events: none;
            transform: scale(0.8);
        }

        /* Welcome Message Bubble */
        .welcome-bubble {
            position: absolute;
            bottom: 70px;
            right: 0;
            background: white;
            padding: 12px 16px;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            font-size: 14px;
            color: #334155;
            white-space: nowrap;
            opacity: 0;
            transform: translateY(10px) scale(0.9);
            transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
            pointer-events: none;
        }
        .welcome-bubble.visible {
            opacity: 1;
            transform: translateY(0) scale(1);
        }
        /* Triangle for bubble */
        .welcome-bubble::after {
            content: '';
            position: absolute;
            bottom: -6px;
            right: 20px;
            width: 12px;
            height: 12px;
            background: white;
            transform: rotate(45deg);
            border-radius: 2px;
        }

        /* Window Container */
        .window-container {
            width: 380px;
            height: 700px;
            max-width: 90vw;
            max-height: 85vh;
            background: white;
            border-radius: 16px;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
            display: flex;
            flex-direction: column;
            overflow: hidden;
            opacity: 0;
            pointer-events: none;
            transform: scale(0.95);
            transition: opacity 0.2s ease, transform 0.2s ease;
            position: absolute; /* Relative to .widget-shell */
            bottom: 0;
            right: 0;
        }

        .window-container.visible {
            opacity: 1;
            pointer-events: auto;
            transform: scale(1);
        }

        /* Drag Handle for Window */
        .window-header {
            height: 32px;
            background: #f8fafc;
            cursor: grab;
            display: flex;
            align-items: center;
            justify-content: center;
            border-bottom: 1px solid #e2e8f0;
            flex-shrink: 0;
        }
        .window-header:active { cursor: grabbing; }
        .drag-indicator {
            width: 40px;
            height: 4px;
            background: #cbd5e1;
            border-radius: 2px;
        }
    `;
    shadow.appendChild(style);

    // Shell Element (The moving part)
    const shell = document.createElement('div');
    shell.className = 'widget-shell';
    // Initial Position (Bottom Right)
    const initialRight = 20;
    const initialBottom = 20;
    shell.style.right = `${initialRight}px`;
    shell.style.bottom = `${initialBottom}px`;
    shadow.appendChild(shell);

    // Floating Button
    const btn = document.createElement('div');
    btn.className = 'floating-btn';

    // Check for custom image
    const customImage = getConfig('VITE_FLOAT_BALL_IMAGE');
    if (customImage) {
        btn.style.background = 'transparent'; // Remove default gradient if image exists
        btn.style.boxShadow = 'none'; // Optional: remove shadow if image has its own or to look cleaner
        btn.innerHTML = `<img src="${customImage}" alt="Chat" draggable="false" />`;
    } else {
        btn.innerHTML = `
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
        `;
    }
    shell.appendChild(btn);

    // Welcome Message Bubble
    const welcomeBubble = document.createElement('div');
    welcomeBubble.className = 'welcome-bubble';
    shell.appendChild(welcomeBubble);

    // Welcome Message Logic
    const welcomeContentsStr = getConfig('VITE_FLOAT_BALL_WELCOME_CONTENTS');
    let welcomeContents: string[] = [];
    try {
        if (welcomeContentsStr) {
            welcomeContents = JSON.parse(welcomeContentsStr);
        }
    } catch (e) {
        console.warn('Failed to parse VITE_FLOAT_BALL_WELCOME_CONTENTS', e);
    }

    if (welcomeContents.length > 0) {
        btn.addEventListener('mouseenter', () => {
            if (isOpen) return; // Don't show if window is open
            const randomMsg = welcomeContents[Math.floor(Math.random() * welcomeContents.length)];
            welcomeBubble.textContent = randomMsg;
            welcomeBubble.classList.add('visible');
        });

        btn.addEventListener('mouseleave', () => {
            welcomeBubble.classList.remove('visible');
        });
    }

    // Window Container
    const winContainer = document.createElement('div');
    winContainer.className = 'window-container';

    // Add Drag Handle and Close Button
    const header = document.createElement('div');
    header.className = 'window-header';

    // Drag Indicator
    const indicator = document.createElement('div');
    indicator.className = 'drag-indicator';
    header.appendChild(indicator);

    // Close Button
    const closeBtn = document.createElement('div');
    closeBtn.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18M6 6l12 12"></path>
        </svg>
    `;
    closeBtn.style.cssText = `
        position: absolute;
        right: 12px;
        cursor: pointer;
        color: #64748b;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        transition: background 0.2s;
    `;
    closeBtn.onmouseover = () => { closeBtn.style.background = '#e2e8f0'; };
    closeBtn.onmouseout = () => { closeBtn.style.background = 'transparent'; };
    closeBtn.onclick = (e) => {
        e.stopPropagation(); // Prevent drag start
        toggle();
    };
    header.appendChild(closeBtn);

    winContainer.appendChild(header);

    // Content Area (for iframe)
    const content = document.createElement('div');
    content.style.flex = '1';
    content.style.position = 'relative';
    winContainer.appendChild(content);

    shell.appendChild(winContainer);

    // --- Logic ---

    let isOpen = initialOpen;

    const updateState = () => {
        if (isOpen) {
            btn.classList.add('hidden');
            welcomeBubble.classList.remove('visible'); // Hide bubble when opening
            winContainer.classList.add('visible');
        } else {
            btn.classList.remove('hidden');
            winContainer.classList.remove('visible');
        }
        container.dispatchEvent(new CustomEvent('litechat-toggle', { detail: { isOpen } }));
    };

    // Initial state
    updateState();

    const toggle = () => {
        isOpen = !isOpen;
        updateState();
    };

    const open = () => {
        if (!isOpen) {
            isOpen = true;
            updateState();
        }
    };

    const close = () => {
        if (isOpen) {
            isOpen = false;
            updateState();
        }
    };

    // Toggle on click (only if not dragging)
    let isDragging = false;
    btn.addEventListener('click', () => {
        if (!isDragging) toggle();
    });

    // --- Drag Implementation (Unified for Shell) ---
    // We drag the SHELL. Both button and window move together.

    const makeDraggable = (handle: HTMLElement) => {
        let startX = 0, startY = 0;
        let startRight = 0, startBottom = 0;

        const onStart = (clientX: number, clientY: number) => {
            isDragging = false;
            startX = clientX;
            startY = clientY;

            const rect = shell.getBoundingClientRect();
            startRight = window.innerWidth - rect.right;
            startBottom = window.innerHeight - rect.bottom;

            document.addEventListener('mousemove', onMove);
            document.addEventListener('mouseup', onEnd);
            document.addEventListener('touchmove', onTouchMove, { passive: false });
            document.addEventListener('touchend', onEnd);
        };

        const onMove = (e: MouseEvent) => handleMove(e.clientX, e.clientY);

        const onTouchMove = (e: TouchEvent) => {
            e.preventDefault();
            handleMove(e.touches[0].clientX, e.touches[0].clientY);
        };

        const handleMove = (clientX: number, clientY: number) => {
            const dx = startX - clientX;
            const dy = startY - clientY;

            if (Math.abs(dx) > 5 || Math.abs(dy) > 5) isDragging = true;

            shell.style.right = `${startRight + dx}px`;
            shell.style.bottom = `${startBottom + dy}px`;
            shell.style.transition = 'none';
        };

        const onEnd = () => {
            document.removeEventListener('mousemove', onMove);
            document.removeEventListener('mouseup', onEnd);
            document.removeEventListener('touchmove', onTouchMove);
            document.removeEventListener('touchend', onEnd);

            // Rebound / Boundary Check
            const activeEl = isOpen ? winContainer : btn;
            const rect = activeEl.getBoundingClientRect();
            const winWidth = window.innerWidth;
            const winHeight = window.innerHeight;

            let finalRight = winWidth - rect.right;
            let finalBottom = winHeight - rect.bottom;

            const margin = 20;

            if (rect.left < margin) {
                finalRight = winWidth - margin - rect.width;
            }
            if (rect.right > winWidth - margin) {
                finalRight = margin;
            }
            if (rect.top < margin) {
                finalBottom = winHeight - margin - rect.height;
            }
            if (rect.bottom > winHeight - margin) {
                finalBottom = margin;
            }

            shell.style.transition = 'all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)';
            shell.style.right = `${finalRight}px`;
            shell.style.bottom = `${finalBottom}px`;

            setTimeout(() => { isDragging = false; }, 50);
        };

        handle.addEventListener('mousedown', (e) => onStart(e.clientX, e.clientY));
        handle.addEventListener('touchstart', (e) => onStart(e.touches[0].clientX, e.touches[0].clientY), { passive: false });
    };

    makeDraggable(btn);
    makeDraggable(header);

    // Initial Boundary Check (in case window resize pushed it out)
    const checkBounds = () => {
        const activeEl = isOpen ? winContainer : btn;
        const rect = activeEl.getBoundingClientRect();

        // If hidden (display: none), rect will be all zeros. Do not update position.
        if (rect.width === 0 && rect.height === 0) return;

        let needsUpdate = false;
        let finalRight = window.innerWidth - rect.right;
        let finalBottom = window.innerHeight - rect.bottom;
        const margin = 20;
        const winWidth = window.innerWidth;
        const winHeight = window.innerHeight;

        if (rect.left < margin) {
            finalRight = winWidth - margin - rect.width;
            needsUpdate = true;
        }
        if (rect.right > winWidth - margin) {
            finalRight = margin;
            needsUpdate = true;
        }
        if (rect.top < margin) {
            finalBottom = winHeight - margin - rect.height;
            needsUpdate = true;
        }
        if (rect.bottom > winHeight - margin) {
            finalBottom = margin;
            needsUpdate = true;
        }

        if (needsUpdate) {
            shell.style.transition = 'all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)';
            shell.style.right = `${finalRight}px`;
            shell.style.bottom = `${finalBottom}px`;
        }
    };
    window.addEventListener('resize', checkBounds);

    // Also check bounds when toggling, because window is bigger than button
    container.addEventListener('litechat-toggle', () => {
        setTimeout(checkBounds, 50);
    });

    setTimeout(checkBounds, 100);

    return { shadowRoot: shadow, windowContainer: content, toggleBtn: btn, open, close, toggle };
}
