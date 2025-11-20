<script setup lang="ts">
import { ref, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useChatStore } from '../../stores/chat';
import { useAuthStore } from '../../stores/auth';
import ThemeToggle from '../ThemeToggle.vue';
import { 
  Search, 
  BookOpen, 
  Bookmark, 
  Clock, 
  LogOut,
  X,
  PanelLeftClose,
  PanelLeftOpen,
  ChevronDown,
  ChevronRight,
  MoreHorizontal,
  Edit2,
  Trash2,
  User
} from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
  isCollapsed: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'toggleCollapse'): void;
}>();

const route = useRoute();
const router = useRouter();
const store = useChatStore();
const authStore = useAuthStore();
const isHistoryOpen = ref(true);

const handleLogout = () => {
  authStore.logout();
};

const navItems = [
  { icon: Search, label: '搜索', active: false, path: '/' },
  { icon: BookOpen, label: '知识库', active: false, path: '/knowledge' },
  { icon: Bookmark, label: '收藏', active: false, path: '/collections' },
];

// History Management
const activeDropdownId = ref<string | null>(null);
const editingId = ref<string | null>(null);
const editTitle = ref('');
const editInputRef = ref<HTMLInputElement | null>(null);
const contextMenu = ref<{ show: boolean; x: number; y: number; id: string | null }>({
  show: false,
  x: 0,
  y: 0,
  id: null
});

const toggleDropdown = (id: string) => {
  if (activeDropdownId.value === id) {
    activeDropdownId.value = null;
  } else {
    activeDropdownId.value = id;
    contextMenu.value.show = false;
  }
};

const closeDropdown = () => {
  activeDropdownId.value = null;
};

const handleContextMenu = (e: MouseEvent, id: string) => {
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

const startRename = (id: string, currentTitle: string) => {
  editingId.value = id;
  editTitle.value = currentTitle;
  activeDropdownId.value = null;
  contextMenu.value.show = false;
  nextTick(() => {
    editInputRef.value?.focus();
  });
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

const deleteHistory = (id: string) => {
  if (confirm('确定要删除这条对话吗？')) {
    store.deleteChat(id);
    if (route.path === `/chat/${id}`) {
      router.push('/');
    }
  }
  activeDropdownId.value = null;
  contextMenu.value.show = false;
};

// Click outside directive
const vClickOutside = {
  mounted(el: any, binding: any) {
    el.clickOutsideEvent = (event: Event) => {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value(event);
      }
    };
    document.body.addEventListener('click', el.clickOutsideEvent);
  },
  unmounted(el: any) {
    document.body.removeEventListener('click', el.clickOutsideEvent);
  },
};
</script>

<template>
  <!-- Mobile Overlay -->
  <div 
    v-if="isOpen" 
    class="fixed inset-0 bg-black/50 z-40 md:hidden"
    @click="emit('close')"
  ></div>

  <!-- Sidebar -->
  <aside 
    class="fixed md:static inset-y-0 left-0 z-50 bg-[#f9f9f9] dark:bg-[#1a1a1a] border-r border-gray-200 dark:border-gray-800 transform transition-all duration-300 ease-in-out flex flex-col"
    :class="[
      isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
      isCollapsed ? 'md:w-20' : 'w-64'
    ]"
  >
    <!-- Header -->
    <div 
      class="flex items-center transition-all duration-300"
      :class="[
        isCollapsed 
          ? 'flex-col justify-center py-4 gap-4 h-auto' 
          : 'h-16 justify-between px-4 border-b border-transparent'
      ]"
    >
      <!-- Desktop Collapse Button -->
      <button 
        class="hidden md:flex text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-800 p-1 rounded-lg transition-colors shrink-0"
        @click="emit('toggleCollapse')"
        :title="isCollapsed ? '展开侧边栏' : '收起侧边栏'"
      >
        <PanelLeftOpen v-if="isCollapsed" class="w-5 h-5" />
        <PanelLeftClose v-else class="w-5 h-5" />
      </button>

      <div class="flex items-center gap-2 text-primary font-bold text-xl overflow-hidden whitespace-nowrap">
        <span class="bg-primary text-white p-1 rounded shrink-0">知</span>
        <span v-show="!isCollapsed">知乎直答</span>
      </div>

      <!-- Mobile Close Button -->
      <button 
        class="md:hidden text-gray-500"
        @click="emit('close')"
      >
        <X class="w-6 h-6" />
      </button>
    </div>

    <!-- Nav Items -->
    <nav class="flex-1 px-3 py-4 space-y-2 overflow-y-auto">
    <router-link 
        v-for="item in navItems" 
        :key="item.label"
        :to="item.path" 
        class="flex items-center gap-3 px-3 py-3 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 ease-in-out hover:translate-x-1"
        :class="[
          { 'bg-blue-50 text-primary dark:bg-blue-900/20 dark:text-blue-400': route.path === item.path },
          isCollapsed ? 'justify-center' : ''
        ]"
        :title="isCollapsed ? item.label : ''"
      >
        <component :is="item.icon" class="w-5 h-5 shrink-0" />
        <span v-show="!isCollapsed" class="whitespace-nowrap overflow-hidden">{{ item.label }}</span>
      </router-link>

      <!-- History Section -->
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
                  <div v-if="editingId === history.id" class="px-3 py-2 pl-9">
                    <input 
                      ref="editInputRef"
                      v-model="editTitle"
                      @blur="confirmRename"
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
                      @click.prevent.stop="toggleDropdown(history.id)"
                      class="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    >
                      <MoreHorizontal class="w-4 h-4" />
                    </button>
                  </div>

                  <!-- Dropdown Menu -->
                  <div 
                    v-if="activeDropdownId === history.id"
                    class="absolute right-0 top-full mt-1 w-32 bg-white dark:bg-[#2a2a2a] rounded-lg shadow-xl border border-gray-100 dark:border-gray-700 z-50 py-1 overflow-hidden"
                    v-click-outside="closeDropdown"
                  >
                    <button 
                      @click="startRename(history.id, history.title)"
                      class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left"
                    >
                      <Edit2 class="w-3 h-3" />
                      重命名
                    </button>
                    <button 
                      @click="deleteHistory(history.id)"
                      class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left"
                    >
                      <Trash2 class="w-3 h-3" />
                      删除
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

    <!-- Context Menu -->
    <div 
      v-if="contextMenu.show"
      class="fixed z-[100] w-32 bg-white dark:bg-[#2a2a2a] rounded-lg shadow-xl border border-gray-100 dark:border-gray-700 py-1 overflow-hidden"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      v-click-outside="closeContextMenu"
    >
      <button 
        @click="startRename(contextMenu.id!, store.historyItems.find(h => h.id === contextMenu.id)?.title || '')"
        class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 text-left"
      >
        <Edit2 class="w-3 h-3" />
        重命名
      </button>
      <button 
        @click="deleteHistory(contextMenu.id!)"
        class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 text-left"
      >
        <Trash2 class="w-3 h-3" />
        删除
      </button>
    </div>
        </div>
      </div>
    </nav>

    <!-- Footer -->
    <div class="p-3 border-t border-gray-200 dark:border-gray-800">
      <div v-if="!isCollapsed" class="flex items-center justify-between gap-2">
        <ThemeToggle />
        
        <div v-if="authStore.isAuthenticated" class="flex items-center gap-2 flex-1 min-w-0">
          <img 
            :src="authStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'" 
            class="w-8 h-8 rounded-full bg-gray-100"
          />
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ authStore.user?.username }}</div>
          </div>
          <button 
            @click="handleLogout"
            class="p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
            title="退出登录"
          >
            <LogOut class="w-4 h-4" />
          </button>
        </div>

        <button 
          v-else
          @click="authStore.openModal()"
          class="flex items-center justify-center gap-2 px-3 py-2 rounded-lg bg-primary text-white hover:bg-blue-600 transition-colors flex-1 text-sm font-medium"
        >
          登录
        </button>
      </div>

      <div v-else class="flex flex-col gap-2 items-center">
        <ThemeToggle />
        
        <div v-if="authStore.isAuthenticated" class="relative group">
          <img 
            :src="authStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'" 
            class="w-8 h-8 rounded-full bg-gray-100 cursor-pointer"
          />
          <div class="absolute left-full bottom-0 ml-2 mb-[-4px] hidden group-hover:block z-50">
            <button 
              @click="handleLogout"
              class="flex items-center gap-2 px-3 py-2 bg-white dark:bg-[#2a2a2a] border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg text-sm text-red-500 hover:bg-gray-50 dark:hover:bg-gray-800 whitespace-nowrap"
            >
              <LogOut class="w-4 h-4" />
              退出
            </button>
          </div>
        </div>

        <button 
          v-else
          @click="authStore.openModal()"
          class="w-8 h-8 flex items-center justify-center rounded-full bg-primary text-white hover:bg-blue-600 transition-colors"
          title="登录"
        >
          <User class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>
```
