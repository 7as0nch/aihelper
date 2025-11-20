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

export default router;
