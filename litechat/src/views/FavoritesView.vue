<script setup lang="ts">
import { ref, computed } from 'vue';
import { useChatStore } from '../stores/chat';
import { Search, Star, MessageSquare, Calendar } from 'lucide-vue-next';
import { useRouter } from 'vue-router';

const store = useChatStore();
const router = useRouter();
const searchQuery = ref('');

const favoriteChats = computed(() => {
  return store.historyItems.filter(item => store.isFavorite(item.id));
});

const filteredFavorites = computed(() => {
  if (!searchQuery.value) return favoriteChats.value;
  const query = searchQuery.value.toLowerCase();
  return favoriteChats.value.filter(item => 
    item.title.toLowerCase().includes(query)
  );
});

const navigateToChat = (id: string) => {
  router.push(`/chat/${id}`);
};

const unstar = (e: Event, id: string) => {
  e.stopPropagation();
  store.toggleFavorite(id);
};
</script>

<template>
  <div class="flex-1 flex flex-col h-full bg-gray-50 dark:bg-[#1a1a1a] overflow-hidden">
    <!-- Header -->
    <div class="px-8 py-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">收藏列表</h1>
      
      <!-- Search -->
      <div class="relative max-w-xl">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="搜索你收藏的内容"
          class="w-full pl-10 pr-4 py-3 bg-white dark:bg-[#2a2a2a] border-none rounded-xl shadow-sm focus:ring-2 focus:ring-primary/50 outline-none transition-all text-gray-900 dark:text-white placeholder-gray-400"
        />
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-8 pb-8">
      <div v-if="filteredFavorites.length > 0" class="grid gap-4">
        <div 
          v-for="item in filteredFavorites" 
          :key="item.id"
          @click="navigateToChat(item.id)"
          class="group bg-white dark:bg-[#2a2a2a] p-6 rounded-xl shadow-sm hover:shadow-md transition-all cursor-pointer border border-transparent hover:border-primary/20"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2 truncate group-hover:text-primary transition-colors">
                {{ item.title }}
              </h3>
              <p class="text-gray-500 dark:text-gray-400 text-sm line-clamp-2 mb-4">
                {{ '这里是对话的摘要内容，实际应用中应该从 store 获取最后一条消息的内容...' }}
              </p>
              <div class="flex items-center gap-4 text-xs text-gray-400">
                <span class="flex items-center gap-1">
                  <Calendar class="w-3 h-3" />
                  2025.04.08
                </span>
                <span class="flex items-center gap-1">
                  <MessageSquare class="w-3 h-3" />
                  {{ Math.floor(Math.random() * 20) + 1 }} 条对话
                </span>
              </div>
            </div>
            
            <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button 
                @click="unstar($event, item.id)"
                class="p-2 text-yellow-500 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 rounded-lg transition-colors"
                title="取消收藏"
              >
                <Star class="w-5 h-5 fill-current" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="flex flex-col items-center justify-center h-64 text-gray-400">
        <div class="w-16 h-16 bg-gray-100 dark:bg-[#2a2a2a] rounded-full flex items-center justify-center mb-4">
          <Star class="w-8 h-8 text-gray-300 dark:text-gray-600" />
        </div>
        <p class="text-lg font-medium mb-1">暂无收藏内容</p>
        <p class="text-sm">点击对话右上角的星号即可收藏</p>
      </div>
    </div>
  </div>
</template>
