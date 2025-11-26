import { ref, nextTick, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useChatStore } from '../../../../stores/chat';

/**
 * Composable for managing history menu functionality
 * Handles rename, delete, and context menu interactions
 */
export function useHistoryMenu() {
    const route = useRoute();
    const router = useRouter();
    const store = useChatStore();

    // State
    const activeDropdownId = ref<string | null>(null);
    const editingId = ref<string | null>(null);
    const editTitle = ref('');
    const deleteConfirmId = ref<string | null>(null);
    const contextMenu = ref<{ show: boolean; x: number; y: number; id: string | null }>({
        show: false,
        x: 0,
        y: 0,
        id: null
    });

    // Ref callback for edit input
    const setEditInputRef = (el: any) => {
        if (el) {
            nextTick(() => {
                el.focus();
                el.select();
            });
        }
    };

    // Dropdown actions
    const closeDropdown = () => {
        activeDropdownId.value = null;
    };

    // Context menu actions
    const handleContextMenu = (e: MouseEvent, id: string) => {
        e.preventDefault();
        e.stopPropagation();
        contextMenu.value = {
            show: true,
            x: e.clientX,
            y: e.clientY,
            id
        };
        activeDropdownId.value = null;
    };

    const closeContextMenu = () => {
        contextMenu.value.show = false;
    };

    // Rename actions
    const startRename = (id: string, currentTitle: string) => {
        editingId.value = id;
        editTitle.value = currentTitle;
        activeDropdownId.value = null;
        contextMenu.value.show = false;
    };

    const confirmRename = () => {
        if (editingId.value && editTitle.value.trim()) {
            store.renameChat(editingId.value, editTitle.value.trim());
        }
        cancelRename();
    };

    const cancelRename = () => {
        editingId.value = null;
        editTitle.value = '';
    };

    // Delete actions
    const showDeleteConfirm = (id: string) => {
        deleteConfirmId.value = id;
    };

    const confirmDelete = async (id: string) => {
        await store.deleteChat(id);
        if (route.path === `/chat/${id}`) {
            router.push('/');
        }
        deleteConfirmId.value = null;
        activeDropdownId.value = null;
        contextMenu.value.show = false;
    };

    const cancelDelete = () => {
        deleteConfirmId.value = null;
    };

    // Handle window blur (e.g. clicking into iframe)
    const handleWindowBlur = () => {
        activeDropdownId.value = null;
        contextMenu.value.show = false;
        // Only cancel rename if we're not focused on the input (which we wouldn't be if window blurred)
        if (editingId.value) {
            cancelRename();
        }
    };

    onMounted(() => {
        window.addEventListener('blur', handleWindowBlur);
    });

    onUnmounted(() => {
        window.removeEventListener('blur', handleWindowBlur);
    });

    return {
        // State
        activeDropdownId,
        editingId,
        editTitle,
        deleteConfirmId,
        contextMenu,

        // Methods
        setEditInputRef,
        closeDropdown,
        handleContextMenu,
        closeContextMenu,
        startRename,
        confirmRename,
        cancelRename,
        showDeleteConfirm,
        confirmDelete,
        cancelDelete,
    };
}
