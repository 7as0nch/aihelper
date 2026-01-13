import { defineStore } from 'pinia';
import { ref } from 'vue';
import { sendMessageStream, type Attachment, type Message, chatApi } from '../api/chat';
import { getConfig } from '@/config';
import { useStorage } from '@vueuse/core';
export type { Attachment, Message };




export const useChatStore = defineStore('chat', () => {
    const messages = ref<Message[]>([]);
    const isLoading = ref(false);
    const isThinking = ref(false);
    const thinkingMode = useStorage<'smart' | 'deep' | 'quick'>('litechat_thinking_mode', 'smart');
    const searchByWeb = useStorage<boolean>('litechat_search_by_web', false);

    const currentChatId = ref<string | null>(null);

    const historyItems = ref<{ id: string; title: string }[]>([]);

    const isHistoryLoading = ref(false);
    const historyPage = ref(1);
    const historyTotal = ref(0);
    const historyPageSize = ref(20);

    const fetchHistoryList = async (loadMore = false) => {
        if (isHistoryLoading.value) return;

        // If loading more and we've reached the end, stop.
        if (loadMore && historyItems.value.length >= historyTotal.value) {
            return;
        }

        if (!loadMore) {
            historyPage.value = 1;
        } else {
            historyPage.value++;
        }

        isHistoryLoading.value = true;
        try {
            const res = await chatApi.getHistoryList({
                page: historyPage.value,
                pageSize: historyPageSize.value
            });

            historyTotal.value = res.total;

            if (loadMore) {
                // Filter out duplicates just in case
                const newItems = res.sessions.filter(s => !historyItems.value.some(args => args.id === s.id));
                historyItems.value.push(...newItems);
            } else {
                historyItems.value = res.sessions;
            }
        } catch (e) {
            console.error('Failed to fetch history list', e);
            if (loadMore) historyPage.value--; // Revert
        } finally {
            isHistoryLoading.value = false;
        }
    };

    const clearHistoryList = () => {
        historyItems.value = [];
        historyPage.value = 1;
        historyTotal.value = 0;
    };

    const addMessage = (message: Message) => {
        messages.value.push(message);
    };

    const updateLastMessage = (data: Partial<Message>) => {
        const lastMsg = messages.value[messages.value.length - 1];
        if (lastMsg && lastMsg.role === 'assistant') {
            if (data.content) {
                lastMsg.content += data.content;
            }
            if (data.reasoningContent) {
                lastMsg.reasoningContent = (lastMsg.reasoningContent || '') + data.reasoningContent;
            }
            // Merge other fields if present
            if (data.attachments) {
                lastMsg.attachments = [...(lastMsg.attachments || []), ...data.attachments];
            }
            if (data.tokenUsage) {
                lastMsg.tokenUsage = data.tokenUsage;
            }
            if (data.quoteId) lastMsg.quoteId = data.quoteId;
            if (data.quoteContent) lastMsg.quoteContent = data.quoteContent;
            if (data.quoteSearchLinks) lastMsg.quoteSearchLinks = data.quoteSearchLinks;
            if (data.callingTools) lastMsg.callingTools = data.callingTools;
            if (data.aiModel) lastMsg.aiModel = { ...lastMsg.aiModel, ...data.aiModel };

            // Persist to LocalStorage in frontend mode
            if (currentChatId.value && getConfig('VITE_AI_TYPE') === 'frontend') {
                const currentHistoryItem = historyItems.value.find(h => h.id === currentChatId.value);
                if (currentHistoryItem) {
                    // We need to cast chatApi to any or extend interface to include saveChat
                    // For now, let's assume we can call it.
                    // Ideally we should update the interface.
                    (chatApi as any).saveChat(currentChatId.value, currentHistoryItem.title, messages.value);
                }
            }
        }
    };

    const loadChatHistory = async (chatId: string) => {
        currentChatId.value = chatId;
        try {
            messages.value = await chatApi.getHistoryMsg(chatId);
        } catch (e) {
            console.error('Failed to load chat history', e);
            messages.value = [];
        }
    };

    const clearMessages = () => {
        messages.value = [];
        currentChatId.value = null;
    };

    const deleteChat = async (id: string) => {
        try {
            await chatApi.deleteChat(id);
            // Immediately remove from local state for instant feedback
            historyItems.value = historyItems.value.filter(item => item.id !== id);
            historyTotal.value = Math.max(0, historyTotal.value - 1);
            
            if (currentChatId.value === id) {
                clearMessages();
            }
            // Optionally refetch to ensure pagination state is correct
            // fetchHistoryList(); 
        } catch (e) {
            console.error('Failed to delete chat', e);
        }
    };

    const renameChat = async (id: string, newTitle: string) => {
        try {
            await chatApi.renameChat(id, newTitle);
            const item = historyItems.value.find(h => h.id === id);
            if (item) {
                item.title = newTitle;
            }
        } catch (e) {
            console.error('Failed to rename chat', e);
        }
    };

    const abortController = ref<AbortController | null>(null);

    const stopGeneration = () => {
        if (abortController.value) {
            abortController.value.abort();
            abortController.value = null;
            isLoading.value = false;
            isThinking.value = false;

            // Update last message status
            const lastMsg = messages.value[messages.value.length - 1];
            if (lastMsg && lastMsg.role === 'assistant') {
                lastMsg.content += ' [已停止生成]';

                // Persist to LocalStorage/Backend immediately
                if (currentChatId.value && getConfig('VITE_AI_TYPE') === 'frontend') {
                    const currentHistoryItem = historyItems.value.find(h => h.id === currentChatId.value);
                    if (currentHistoryItem) {
                        (chatApi as any).saveChat(currentChatId.value, currentHistoryItem.title, messages.value);
                    }
                }
            }
        }
    };

    const sendMessage = async (content: string, attachments: Attachment[] = [], quote?: { quoteId: string; quoteContent: string }) => {
        // User message
        const userMessage: Message = {
            id: Date.now().toString(),
            role: 'user',
            content,
            timestamp: Date.now(),
            attachments,
            quoteId: quote?.quoteId,
            quoteContent: quote?.quoteContent,
            aiModel: {
                id: 'deepseek', // Default mock model
                modelName: 'deepseek',
                thinkingMode: thinkingMode.value,
                searchByWeb: searchByWeb.value,
            }
        };
        // Auto-create session if not exists
        if (!currentChatId.value) {
            // Generate title from first 20 chars of content
            const title = content.slice(0, 20) + (content.length > 20 ? '...' : '');

            const session = await chatApi.saveChat("0", title, [userMessage]);
            if (session.id && session.id !== "0") {
                currentChatId.value = session.id;
                // Add to history.
                historyItems.value.unshift({
                    id: session.id,
                    title: title
                });
            }
        }

        // User message
        addMessage(userMessage);

        // Persist immediately
        if (currentChatId.value && getConfig('VITE_AI_TYPE') === 'frontend') {
            console.log('Persisting chat', currentChatId.value);
            const currentHistoryItem = historyItems.value.find(h => h.id === currentChatId.value);
            if (currentHistoryItem) {
                chatApi.saveChat(currentChatId.value, currentHistoryItem.title, messages.value);
            }
        }

        isLoading.value = true;
        isThinking.value = true;

        // Create new AbortController
        abortController.value = new AbortController();
        const signal = abortController.value.signal;

        try {
            isThinking.value = false;

            // Simulate AI response
            const responseId = (Date.now() + 1).toString();
            addMessage({
                id: responseId,
                role: 'assistant',
                content: '',
                timestamp: Date.now(),
                aiModel: {
                    id: 'gpt-4',
                    modelName: 'GPT-4',
                    thinkingMode: thinkingMode.value,
                }
            });

            await sendMessageStream(
                {
                    history: messages.value.slice(0, -2), // Exclude current user msg and pending assistant msg
                    curMessage: userMessage,
                    curSessionID: currentChatId.value || '', // 当前会话ID
                },
                (data) => {
                    updateLastMessage(data);
                },
                signal
            );


        } catch (error) {
            if (error instanceof DOMException && error.name === 'AbortError') {
                console.log('Generation stopped by user');
            } else {
                console.error('Generation error:', error);
                const lastMsg = messages.value[messages.value.length - 1];
                if (lastMsg && lastMsg.role === 'assistant') {
                    lastMsg.content += '\n[生成出错]';
                }
            }
        } finally {
            isLoading.value = false;
            isThinking.value = false;
            abortController.value = null;


        }
    };

    const favorites = ref<string[]>([]);

    const toggleFavorite = (chatId: string) => {
        const index = favorites.value.indexOf(chatId);
        if (index === -1) {
            favorites.value.push(chatId);
        } else {
            favorites.value.splice(index, 1);
        }
    };

    const isFavorite = (chatId: string) => {
        return favorites.value.includes(chatId);
    };

    return {
        messages,
        currentChatId,
        historyItems,
        favorites,
        isLoading,
        isThinking,
        thinkingMode,
        searchByWeb,
        sendMessage,
        loadChatHistory,
        clearMessages,
        deleteChat,
        renameChat,
        fetchHistoryList,
        clearHistoryList,
        toggleFavorite,
        isFavorite,
        stopGeneration,
        isHistoryLoading,
        historyPage,
        historyTotal,
    };
});
