import { createRouter, createWebHistory } from 'vue-router';
import ChatView from '@/views/ChatView.vue';
import KnowledgeBaseView from '@/views/KnowledgeBaseView.vue';
import CollectionsView from '@/views/CollectionsView.vue';
import HistoryView from '@/views/HistoryView.vue';

const routes = [
    {
        path: '/',
        name: 'Chat',
        component: ChatView,
    },
    {
        path: '/chat/:id',
        name: 'ChatHistory',
        component: ChatView,
    },
    {
        path: '/knowledge',
        name: 'KnowledgeBase',
        component: KnowledgeBaseView,
    },
    {
        path: '/collections',
        name: 'Collections',
        component: CollectionsView,
    },
    {
        path: '/history',
        name: 'History',
        component: HistoryView,
    },
    {
        path: '/history',
        name: 'History',
        component: HistoryView,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, _from, next) => {
    // Dynamic import to avoid circular dependency issues if any, 
    // though mainly to ensure pinia is active.
    // In a real app, we might import useAuthStore at top level if pinia instance is ready.
    // But here we can just import it.
    // Actually, we need to be careful about Pinia initialization. 
    // Ideally, we should import the store hook inside the guard.

    // We need to import the store here to avoid "Pinia not installed" error if imported at top level before app.use(pinia)
    // However, since this file is imported in main.ts, we can't use the hook at top level.
    // We'll use a helper or just import it inside.

    // Check if route requires auth
    // The user requirement: "Except the first one (Home), others need validation"
    const publicRoutes = ['Chat'];

    // If it's not a public route
    if (!publicRoutes.includes(to.name as string)) {
        // We need to check auth status. 
        // Since we can't easily access the store instance here without importing it,
        // and we want to avoid circular deps, let's try importing.
        // But wait, we can access localStorage directly for a quick check or import store.
        // Let's try importing store inside.

        // Note: We need to handle the case where authStore isn't ready yet? 
        // No, router is used after pinia is installed in main.ts usually.
        // But let's check main.ts.
        // main.ts: app.use(pinia).use(router)
        // So router is installed after pinia.

        // However, `router` is defined here.
        // Let's just use localStorage for a simple check first, or better, use the store.

        const token = localStorage.getItem('token');
        if (!token) {
            // Redirect to home and show auth modal? 
            // The user said "verify login". If not logged in, we should probably block access.
            // But simply redirecting to '/' might be confusing if we don't show the modal.
            // We can use a query param or just rely on the user clicking login.
            // Or we can trigger the modal if we have access to the store.

            // For now, let's redirect to '/' if not authenticated.
            // And maybe we can trigger the modal via a global event bus or store if we can access it.

            // Let's try to access the store.
            import('@/stores/auth').then(({ useAuthStore }) => {
                const authStore = useAuthStore();
                if (!authStore.isAuthenticated) {
                    authStore.openModal();
                    next('/');
                } else {
                    next();
                }
            });
            return;
        }
    }
    next();
});

export default router;
