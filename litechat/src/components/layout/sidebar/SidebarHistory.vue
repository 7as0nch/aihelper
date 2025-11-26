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
  <div class="pt-4" :class="{ 'border-t border-gray-200 dark:border-gray-800 mt-2': !isCollapsed }">
    <!-- Collapsed State: Icon Only -->
    <router-link 
      v-if="isCollapsed"
      to="/"
      class="flex items-center justify-center px-3 py-3 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
      :class="{ 'bg-blue-50 text-primary dark:bg-blue-900/20 dark:text-blue-400': route.path === '/history' }"
      title="历史记录"
    >
      <Clock class="w-5 h-5 shrink-0" />
    </router-link>

    <!-- Expanded State: Accordion -->
    <div v-else>
      <button 
        @click="isHistoryOpen = !isHistoryOpen"
        class="w-full flex items-center justify-between px-3 py-2 text-xs font-medium text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors group"
      >
        <span class="flex items-center gap-2">
          <Clock class="w-4 h-4" />
          历史记录
        </span>
        <component 
          :is="isHistoryOpen ? ChevronDown : ChevronRight" 
          class="w-3 h-3 transition-transform duration-200" 
        />
      </button>

      <div 
        class="grid transition-all duration-300 ease-in-out"
        :class="isHistoryOpen ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'"
      >
        <div class="overflow-hidden">
          <div class="mt-1 space-y-0.5">
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
                class="absolute right-0 top-full mt-1 w-32 bg-white dark:bg-[#2a2a2a] rounded-lg shadow-xl border border-gray-100 dark:border-gray-700 z-50 overflow-hidden"
                v-click-outside="closeDropdown"
              >
                <!-- Normal Menu -->
                <template v-if="deleteConfirmId !== history.id">
                  <button 
                    @click.stop="startRename(history.id, history.title)"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left"
                  >
                    <Edit2 class="w-3 h-3" />
                    重命名
                  </button>
                  <button 
                    @click.stop="showDeleteConfirm(history.id)"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left"
                  >
                    <Trash2 class="w-3 h-3" />
                    删除
                  </button>
                </template>
                
                <!-- Delete Confirmation -->
                <template v-else>
                  <div class="p-2">
                    <div class="text-xs text-gray-500 dark:text-gray-400 mb-2 px-1">
                      确认删除此对话？
                    </div>
                    <div class="flex gap-1">
                      <button 
                        @click.stop="confirmDelete(history.id)"
                        class="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 text-xs bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/30 rounded transition-colors"
                      >
                        <Trash2 class="w-3 h-3" />
                        删除
                      </button>
                      <button 
                        @click.stop="cancelDelete"
                        class="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 text-xs bg-gray-50 text-gray-600 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 rounded transition-colors"
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
        v-if="contextMenu.show"
        class="fixed z-[100] w-32 bg-white dark:bg-[#2a2a2a] rounded-lg shadow-xl border border-gray-100 dark:border-gray-700 overflow-hidden"
        :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
        v-click-outside="closeContextMenu"
      >
        <!-- Normal Menu -->
        <template v-if="deleteConfirmId !== contextMenu.id">
          <button 
            @click.stop="startRename(contextMenu.id!, store.historyItems.find(h => h.id === contextMenu.id)?.title || '')"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left"
          >
            <Edit2 class="w-3 h-3" />
            重命名
          </button>
          <button 
            @click.stop="showDeleteConfirm(contextMenu.id!)"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left"
          >
            <Trash2 class="w-3 h-3" />
            删除
          </button>
        </template>
        
        <!-- Delete Confirmation -->
        <template v-else>
          <div class="p-2">
            <div class="text-xs text-gray-500 dark:text-gray-400 mb-2 px-1">
              确认删除此对话？
            </div>
            <div class="flex gap-1">
              <button 
                @click.stop="confirmDelete(contextMenu.id!)"
                class="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 text-xs bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/30 rounded transition-colors"
              >
                <Trash2 class="w-3 h-3" />
                删除
              </button>
              <button 
                @click.stop="cancelDelete"
                class="flex-1 flex items-center justify-center gap-1 px-2 py-1.5 text-xs bg-gray-50 text-gray-600 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 rounded transition-colors"
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
