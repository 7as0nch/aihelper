import { onMounted, onUnmounted, ref, shallowRef } from 'vue';
import { initAiChat, type InitOptions } from './index';
import type { AiChat } from './types';

export function useAiChat(options: InitOptions = {}): AiChat {
    const widgetInstance = shallowRef<ReturnType<typeof initAiChat> | null>(null);
    const isMounted = ref(false);

    onMounted(() => {
        widgetInstance.value = initAiChat(options);
        isMounted.value = true;
    });

    onUnmounted(() => {
        if (widgetInstance.value) {
            widgetInstance.value.unmount();
            widgetInstance.value = null;
            isMounted.value = false;
        }
    });

    const open = () => widgetInstance.value?.open();
    const close = () => widgetInstance.value?.close();
    const toggle = () => widgetInstance.value?.toggle();
    const show = () => widgetInstance.value?.show();
    const hide = () => widgetInstance.value?.hide();

    return {
        widgetInstance,
        isMounted,
        unmount: () => {
            if (widgetInstance.value) {
                widgetInstance.value.unmount();
                widgetInstance.value = null;
                isMounted.value = false;
            }
        },
        open,
        close,
        toggle,
        show,
        hide
    };
}
