<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { useThemeStore } from '../../stores/theme';
import { X, User, Moon, Sun, Monitor, MessageSquare, Upload, Trash2 } from 'lucide-vue-next';

const props = defineProps<{
  isOpen: boolean;
}>();

// Use props to avoid lint error
console.log('Settings modal open:', props.isOpen);

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const authStore = useAuthStore();
const themeStore = useThemeStore();
const activeTab = ref<'profile' | 'theme' | 'feedback'>('profile');

// Feedback Logic
const feedbackTitle = ref('');
const feedbackDesc = ref('');
const feedbackImages = ref<{ id: string; url: string; file: File }[]>([]);
const fileInput = ref<HTMLInputElement | null>(null);

const handleImageUpload = (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (target.files) {
    const files = Array.from(target.files);
    if (feedbackImages.value.length + files.length > 8) {
      alert('最多上传8张图片');
      return;
    }
    
    files.forEach(file => {
      feedbackImages.value.push({
        id: Date.now() + Math.random().toString(),
        url: URL.createObjectURL(file),
        file
      });
    });
    target.value = '';
  }
};

const removeImage = (id: string) => {
  feedbackImages.value = feedbackImages.value.filter(img => img.id !== id);
};

const submitFeedback = async () => {
  if (!feedbackTitle.value.trim() || !feedbackDesc.value.trim()) {
    alert('请填写标题和描述');
    return;
  }
  
  if (feedbackDesc.value.length > 200) {
    alert('描述请控制在200字以内');
    return;
  }

  // Mock API call
  console.log('Submitting feedback:', {
    title: feedbackTitle.value,
    desc: feedbackDesc.value,
    images: feedbackImages.value.map(i => i.file.name)
  });

  alert('反馈提交成功！');
  feedbackTitle.value = '';
  feedbackDesc.value = '';
  feedbackImages.value = [];
  emit('close');
};

// Profile Logic
const username = ref(authStore.user?.username || '');
const updateProfile = () => {
  if (username.value.trim()) {
    // Mock update
    if (authStore.user) {
        authStore.user.username = username.value;
    }
    alert('个人信息更新成功');
  }
};
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')"></div>
    
    <!-- Modal -->
    <div class="relative w-full max-w-2xl bg-white dark:bg-[#1a1a1a] rounded-2xl shadow-2xl overflow-hidden flex flex-col h-[90vh] md:h-auto md:max-h-[90vh]">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-800 shrink-0">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">设置</h2>
        <button @click="emit('close')" class="p-1 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors">
          <X class="w-6 h-6 text-gray-500" />
        </button>
      </div>
      
      <div class="flex flex-col md:flex-row flex-1 min-h-0">
        <!-- Sidebar -->
        <div class="w-full md:w-48 border-b md:border-b-0 md:border-r border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-[#242424] p-2 flex md:flex-col gap-1 overflow-x-auto md:overflow-visible shrink-0">
          <button 
            @click="activeTab = 'profile'"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors"
            :class="activeTab === 'profile' ? 'bg-white dark:bg-[#2a2a2a] text-primary shadow-sm' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'"
          >
            <User class="w-4 h-4" />
            个人信息
          </button>
          <button 
            @click="activeTab = 'theme'"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors"
            :class="activeTab === 'theme' ? 'bg-white dark:bg-[#2a2a2a] text-primary shadow-sm' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'"
          >
            <Sun class="w-4 h-4" />
            主题设置
          </button>
          <button 
            @click="activeTab = 'feedback'"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors"
            :class="activeTab === 'feedback' ? 'bg-white dark:bg-[#2a2a2a] text-primary shadow-sm' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'"
          >
            <MessageSquare class="w-4 h-4" />
            用户反馈
          </button>
        </div>
        
        <!-- Content -->
        <div class="flex-1 p-6 overflow-y-auto">
          <!-- Profile Tab -->
          <div v-if="activeTab === 'profile'" class="space-y-6">
            <div class="flex items-center gap-4">
              <div class="relative group cursor-pointer">
                <img 
                  :src="authStore.user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'" 
                  class="w-20 h-20 rounded-full bg-gray-100"
                />
                <div class="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <Upload class="w-6 h-6 text-white" />
                </div>
              </div>
              <div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ authStore.user?.username }}</h3>
                <p class="text-sm text-gray-500">ID: {{ authStore.user?.id || 'Unknown' }}</p>
              </div>
            </div>
            
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">用户名</label>
                <input 
                  v-model="username"
                  type="text" 
                  class="w-full px-3 py-2 bg-white dark:bg-[#2a2a2a] border border-gray-300 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all"
                />
              </div>
              <button 
                @click="updateProfile"
                class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                保存修改
              </button>
            </div>
          </div>
          
          <!-- Theme Tab -->
          <div v-if="activeTab === 'theme'" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">外观设置</h3>
            <div class="grid grid-cols-3 gap-4">
              <button 
                @click="themeStore.setMode('light')"
                class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all"
                :class="themeStore.mode === 'light' ? 'border-primary bg-blue-50 dark:bg-blue-900/20' : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'"
              >
                <Sun class="w-8 h-8 text-orange-500" />
                <span class="text-sm font-medium">浅色模式</span>
              </button>
              <button 
                @click="themeStore.setMode('dark')"
                class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all"
                :class="themeStore.mode === 'dark' ? 'border-primary bg-blue-50 dark:bg-blue-900/20' : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'"
              >
                <Moon class="w-8 h-8 text-blue-500" />
                <span class="text-sm font-medium">深色模式</span>
              </button>
              <button 
                @click="themeStore.setMode('system')"
                class="flex flex-col items-center gap-3 p-4 rounded-xl border-2 transition-all"
                :class="themeStore.mode === 'system' ? 'border-primary bg-blue-50 dark:bg-blue-900/20' : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'"
              >
                <Monitor class="w-8 h-8 text-gray-500" />
                <span class="text-sm font-medium">跟随系统</span>
              </button>
            </div>
          </div>
          
          <!-- Feedback Tab -->
          <div v-if="activeTab === 'feedback'" class="space-y-6">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">标题</label>
                <input 
                  v-model="feedbackTitle"
                  type="text" 
                  placeholder="简要描述您的问题或建议"
                  class="w-full px-3 py-2 bg-white dark:bg-[#2a2a2a] border border-gray-300 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  描述 <span class="text-xs text-gray-500 font-normal">(200字以内)</span>
                </label>
                <textarea 
                  v-model="feedbackDesc"
                  rows="4"
                  maxlength="200"
                  placeholder="请详细描述您遇到的问题或建议..."
                  class="w-full px-3 py-2 bg-white dark:bg-[#2a2a2a] border border-gray-300 dark:border-gray-700 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all resize-none"
                ></textarea>
                <div class="text-right text-xs text-gray-500 mt-1">
                  {{ feedbackDesc.length }}/200
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  图片上传 <span class="text-xs text-gray-500 font-normal">(最多8张)</span>
                </label>
                <div class="flex flex-wrap gap-3">
                  <div 
                    v-for="img in feedbackImages" 
                    :key="img.id"
                    class="relative w-20 h-20 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden group"
                  >
                    <img :src="img.url" class="w-full h-full object-cover" />
                    <button 
                      @click="removeImage(img.id)"
                      class="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <Trash2 class="w-5 h-5 text-white" />
                    </button>
                  </div>
                  
                  <button 
                    v-if="feedbackImages.length < 8"
                    @click="fileInput?.click()"
                    class="w-20 h-20 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-700 flex flex-col items-center justify-center text-gray-400 hover:border-primary hover:text-primary transition-colors"
                  >
                    <Upload class="w-6 h-6 mb-1" />
                    <span class="text-xs">上传</span>
                  </button>
                  <input 
                    ref="fileInput"
                    type="file" 
                    accept="image/*" 
                    multiple 
                    class="hidden"
                    @change="handleImageUpload"
                  />
                </div>
              </div>
              
              <div class="pt-4">
                <button 
                  @click="submitFeedback"
                  class="w-full px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  提交反馈
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
