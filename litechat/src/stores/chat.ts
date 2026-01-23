import { defineStore } from 'pinia';
import { ref } from 'vue';
import { sendMessageStream, type Attachment, type Message, type QuoteSearchLink, chatApi } from '../api/chat';
import { getConfig } from '@/config';
import { useStorage } from '@vueuse/core';
export type { Attachment, Message, QuoteSearchLink };




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
            if (data.quoteSearchLinks) {
                const existingLinks = lastMsg.quoteSearchLinks || [];
                const newLinks = data.quoteSearchLinks.filter(nl => 
                    !existingLinks.some(el => el.url === nl.url && el.url !== '')
                );
                if (newLinks.length > 0) {
                    lastMsg.quoteSearchLinks = [...existingLinks, ...newLinks];
                }
            }
            if (data.callingTools) {
                lastMsg.callingTools = data.callingTools;
            }
            if (data.aiModel) lastMsg.aiModel = { ...lastMsg.aiModel, ...data.aiModel };
            if (data.extra) lastMsg.extra = { ...(lastMsg.extra || {}), ...data.extra };

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

    const lastLoadRequestId = ref(0);

    const loadChatHistory = async (chatId: string) => {
        // 先中止之前的生成（如果正在进行）
        stopGeneration();
        
        const requestId = ++lastLoadRequestId.value;
        currentChatId.value = chatId;
        messages.value = []; // 先清空，防止看到旧会话消息
        
        try {
            const history = await chatApi.getHistoryMsg(chatId);
            
            // 检查：如果在请求过程中 chatId 已经变了，或者有更新的请求进来了，则丢弃这次结果
            if (currentChatId.value !== chatId || requestId !== lastLoadRequestId.value) return;
            
            messages.value = history;
            
            // 【自动续传】检查最后一条消息是否正在流式传输中
            const lastMsg = messages.value[messages.value.length - 1];
            if (lastMsg && lastMsg.role === 'assistant' && lastMsg.isStreaming) {
                console.log('Detecting active stream, resuming...');
                const userMsg = messages.value[messages.value.length - 2];
                if (userMsg && userMsg.role === 'user') {
                    resumeStream(userMsg);
                }
            }
        } catch (e) {
            console.error('Failed to load chat history', e);
            if (currentChatId.value === chatId && requestId === lastLoadRequestId.value) {
                messages.value = [];
            }
        }
    };

    // 新增：恢复流式传输
    const resumeStream = async (userMessage: Message) => {
        // 确保先中止之前的生成
        stopGeneration();
        
        isLoading.value = true;
        isThinking.value = false; // 既然是续传，说明已经思考过了
        
        abortController.value = new AbortController();
        const signal = abortController.value.signal;

        try {
            await sendMessageStream(
                {
                    history: messages.value.slice(0, messages.value.length - 2), 
                    curMessage: userMessage,
                    curSessionID: currentChatId.value || '',
                },
                (data) => {
                    // 更新最后一条消息（即那条 marked 为 isStreaming 的消息）
                    updateLastMessage(data);
                },
                signal
            );
        } catch (error) {
            console.error('Resume stream error:', error);
        } finally {
            isLoading.value = false;
            abortController.value = null;
        }
    };

    const clearMessages = () => {
        messages.value = [];
        currentChatId.value = null;
    };

    const deleteChat = async (id: string) => {
        try {
            await chatApi.deleteChat(id);
            // 简单处理：删除后直接清空列表并重新获取第一页，避免分页逻辑错乱
            clearHistoryList();
            fetchHistoryList(false);
            
            if (currentChatId.value === id) {
                clearMessages();
            }
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
        // 先中止之前的生成
        stopGeneration();
        
        // User message
        const userMessage: Message = {
            id: '0',
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
        resumeStream,
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
