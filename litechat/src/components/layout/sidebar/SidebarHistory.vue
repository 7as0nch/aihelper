<script setup lang="ts">
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { useChatStore } from '../../../stores/chat';
import { Clock, ChevronDown, ChevronRight, MoreHorizontal, Edit2, Trash2 } from 'lucide-vue-next';
import { vClickOutside } from './composables/useClickOutside';
import { useHistoryMenu } from './composables/useHistoryMenu';

defineProps<{
  isCollapsed: boolean;
}>();

const route = useRoute();
const store = useChatStore();
const isHistoryOpen = ref(true);

const emit = defineEmits<{
  (e: 'scroll-end'): void;
}>();

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  if (target.scrollTop + target.clientHeight >= target.scrollHeight - 50) {
    emit('scroll-end');
  }
};

// Use history menu composable
const {
  activeDropdownId,
  editingId,
  editTitle,
  deleteConfirmId,
  contextMenu,
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
} = useHistoryMenu();
</script>

<template>
  <div 
    class="flex flex-col min-h-0 pt-4" 
    :class="{ 'border-t border-gray-200 dark:border-gray-800 mt-2': !isCollapsed }"
  >
    <!-- Collapsed State: Icon Only -->
    <div v-if="isCollapsed" class="px-3">
      <router-link 
        to="/"
        class="flex items-center justify-center px-3 py-3 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        :class="{ 'bg-blue-50 text-primary dark:bg-blue-900/20 dark:text-blue-400': route.path === '/history' }"
        title="历史记录"
      >
        <Clock class="w-5 h-5 shrink-0" />
      </router-link>
    </div>

    <!-- Expanded State: Accordion -->
    <div v-else class="flex-1 flex flex-col min-h-0">
      <div class="px-3">
        <button 
          @click="isHistoryOpen = !isHistoryOpen"
          class="w-full flex items-center gap-3 px-3 py-3 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 ease-in-out group"
        >
          <Clock class="w-5 h-5 shrink-0" />
          <span class="flex-1 text-left whitespace-nowrap overflow-hidden">历史记录</span>
          <component 
            :is="isHistoryOpen ? ChevronDown : ChevronRight" 
            class="w-4 h-4 transition-transform duration-200" 
          />
        </button>
      </div>

      <div 
        class="flex-1 min-h-0 grid transition-all duration-300 ease-in-out"
        :class="isHistoryOpen ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'"
      >
        <div class="overflow-hidden flex flex-col">
          <div 
            class="flex-1 overflow-y-auto px-3 py-1 space-y-0.5 custom-scrollbar"
            @scroll="handleScroll"
          >
            <div
              v-for="history in store.historyItems" 
              :key="history.id"
              class="relative group"
              @contextmenu.prevent="handleContextMenu($event, history.id)"
            >
              <!-- Editing Mode -->
              <div 
                v-if="editingId === history.id" 
                class="px-3 py-2 pl-9"
                v-click-outside="cancelRename"
              >
                <input 
                  :ref="setEditInputRef"
                  v-model="editTitle"
                  @keydown.enter="confirmRename"
                  @keydown.esc="cancelRename"
                  class="w-full bg-white dark:bg-gray-800 border border-primary rounded px-2 py-1 text-sm focus:outline-none"
                />
              </div>

              <!-- Normal Mode -->
              <router-link 
                v-else
                :to="`/chat/${history.id}`"
                class="block px-3 py-2 pl-9 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg truncate transition-colors pr-8"
                :class="{ 'bg-blue-50 text-primary dark:bg-blue-900/20 dark:text-blue-400': route.path === `/chat/${history.id}` }"
                :title="history.title"
              >
                {{ history.title }}
              </router-link>

              <!-- Dropdown Trigger (Hover) -->
              <div 
                v-if="editingId !== history.id"
                class="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <button 
                  @click.prevent.stop="handleContextMenu($event, history.id)"
                  class="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  <MoreHorizontal class="w-4 h-4" />
                </button>
              </div>

              <!-- Dropdown Menu -->
              <div 
                v-if="activeDropdownId === history.id"
                class="absolute right-0 top-full mt-1 w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 z-50 overflow-hidden"
                v-click-outside="closeDropdown"
              >
                <!-- Normal Menu -->
                <template v-if="deleteConfirmId !== history.id">
                  <button 
                    @click.stop="startRename(history.id, history.title)"
                    class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left transition-colors"
                  >
                    <Edit2 class="w-4 h-4" />
                    重命名
                  </button>
                  <button 
                    @click.stop="showDeleteConfirm(history.id)"
                    class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left transition-colors"
                  >
                    <Trash2 class="w-4 h-4" />
                    删除
                  </button>
                </template>
                
                <!-- Delete Confirmation -->
                <template v-else>
                  <div class="p-3">
                    <div class="text-sm font-medium text-gray-900 dark:text-gray-100 mb-1">确认删除？</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400 mb-3">
                      此操作无法撤销，将永久删除此对话。
                    </div>
                    <div class="flex gap-2">
                      <button 
                        @click.stop="confirmDelete(history.id)"
                        class="flex-1 px-3 py-2 text-xs font-medium bg-red-600 text-white hover:bg-red-700 rounded-lg transition-colors"
                      >
                        确认删除
                      </button>
                      <button 
                        @click.stop="cancelDelete"
                        class="flex-1 px-3 py-2 text-xs font-medium bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 rounded-lg transition-colors"
                      >
                        取消
                      </button>
                    </div>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Context Menu -->
      <div 
        v-if="contextMenu.show && contextMenu.id"
        class="fixed z-[100] w-48 bg-white dark:bg-[#2a2a2a] rounded-xl shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden"
        :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
        v-click-outside="closeContextMenu"
      >
        <!-- Normal Menu -->
        <template v-if="deleteConfirmId !== contextMenu.id">
          <button 
            @click.stop="startRename(contextMenu.id, store.historyItems.find(h => h.id === contextMenu.id)?.title || '')"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left transition-colors"
          >
            <Edit2 class="w-4 h-4" />
            重命名
          </button>
          <button 
            @click.stop="showDeleteConfirm(contextMenu.id)"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left transition-colors"
          >
            <Trash2 class="w-4 h-4" />
            删除
          </button>
        </template>
        
        <!-- Delete Confirmation -->
        <template v-else>
          <div class="p-3">
            <div class="text-sm font-medium text-gray-900 dark:text-gray-100 mb-1">确认删除？</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mb-3">
              此操作无法撤销，将永久删除此对话。
            </div>
            <div class="flex gap-2">
              <button 
                @click.stop="confirmDelete(contextMenu.id)"
                class="flex-1 px-3 py-2 text-xs font-medium bg-red-600 text-white hover:bg-red-700 rounded-lg transition-colors"
              >
                确认删除
              </button>
              <button 
                @click.stop="cancelDelete"
                class="flex-1 px-3 py-2 text-xs font-medium bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 rounded-lg transition-colors"
              >
                取消
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 滚动条样式已移至全局 style.css */
</style>
