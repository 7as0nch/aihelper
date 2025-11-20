<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { X, QrCode, Smartphone, MessageCircle, Mail } from 'lucide-vue-next';

const authStore = useAuthStore();
const activeTab = ref<'sms' | 'password'>('sms');
const phone = ref('');
const code = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  isLoading.value = true;
  // Simulate login delay
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  // For demo purposes, we just use the phone/username as the username
  const username = phone.value || 'User';
  await authStore.login(username, 'password'); // Mock login
  isLoading.value = false;
};
</script>

<template>
  <div v-if="authStore.showAuthModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-[#2a2a2a] rounded-2xl shadow-2xl w-full max-w-3xl overflow-hidden flex relative animate-scale-in">
      <!-- Close Button -->
      <button 
        @click="authStore.closeModal()"
        class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors z-10"
      >
        <X class="w-6 h-6" />
      </button>

      <!-- Left Side: QR Code -->
      <div class="w-1/2 bg-gray-50 dark:bg-[#242424] p-8 flex flex-col items-center justify-center border-r border-gray-100 dark:border-gray-800 hidden md:flex">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">打开知乎 App</h3>
        <p class="text-sm text-gray-500 mb-8">在「我的页」右上角打开扫一扫</p>
        
        <div class="bg-white p-4 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 mb-4">
          <!-- Mock QR Code -->
          <div class="w-40 h-40 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center relative overflow-hidden">
            <QrCode class="w-24 h-24 text-primary opacity-80" />
            <div class="absolute inset-0 bg-gradient-to-tr from-transparent via-white/20 to-transparent animate-shimmer"></div>
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-md">
                <span class="text-primary font-bold text-xs">知</span>
              </div>
            </div>
          </div>
        </div>
        
        <p class="text-sm text-gray-500 flex items-center gap-2">
          其他扫码方式：微信
        </p>
        
        <div class="mt-auto flex gap-4">
          <button class="px-4 py-1.5 rounded-full border border-gray-200 dark:border-gray-700 text-xs text-gray-600 dark:text-gray-400 hover:bg-white dark:hover:bg-gray-800 transition-colors">
            下载知乎 App
          </button>
          <button class="px-4 py-1.5 rounded-full border border-gray-200 dark:border-gray-700 text-xs text-gray-600 dark:text-gray-400 hover:bg-white dark:hover:bg-gray-800 transition-colors">
            无障碍模式
          </button>
        </div>
      </div>

      <!-- Right Side: Login Form -->
      <div class="w-full md:w-1/2 p-8 md:p-12 flex flex-col">
        <div class="flex gap-6 mb-8 border-b border-gray-100 dark:border-gray-800">
          <button 
            class="pb-3 text-sm font-medium transition-colors relative"
            :class="activeTab === 'sms' ? 'text-primary' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="activeTab = 'sms'"
          >
            验证码登录
            <div v-if="activeTab === 'sms'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full"></div>
          </button>
          <button 
            class="pb-3 text-sm font-medium transition-colors relative"
            :class="activeTab === 'password' ? 'text-primary' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="activeTab = 'password'"
          >
            密码登录
            <div v-if="activeTab === 'password'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full"></div>
          </button>
          <div class="ml-auto">
            <button class="text-xs text-primary hover:underline">开通机构号</button>
          </div>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4 flex-1">
          <div class="space-y-4">
            <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
              <span class="text-gray-500 mr-3 text-sm">中国 +86</span>
              <input 
                v-model="phone"
                type="text" 
                placeholder="手机号" 
                class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
              />
            </div>
            
            <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
              <input 
                v-model="code"
                :type="activeTab === 'sms' ? 'text' : 'password'" 
                :placeholder="activeTab === 'sms' ? '输入 6 位短信验证码' : '输入密码'" 
                class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
              />
              <button 
                v-if="activeTab === 'sms'"
                type="button"
                class="text-xs text-primary hover:text-blue-600 font-medium"
              >
                获取短信验证码
              </button>
            </div>
          </div>


          <button 
            type="submit"
            class="w-full bg-primary text-white py-2.5 rounded-lg font-medium hover:bg-blue-600 transition-colors disabled:opacity-50 mt-6"
            :disabled="isLoading"
          >
            {{ isLoading ? '登录中...' : '登录/注册' }}
          </button>

          <div class="mt-8">
            <div class="relative flex justify-center text-xs mb-4">
              <span class="bg-white dark:bg-[#2a2a2a] px-2 text-gray-400">其他方式登录</span>
            </div>
            <div class="flex justify-center gap-6">
              <button type="button" class="p-2 rounded-full bg-gray-50 dark:bg-gray-800 text-green-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                <MessageCircle class="w-5 h-5" />
              </button>
              <button type="button" class="p-2 rounded-full bg-gray-50 dark:bg-gray-800 text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                <Smartphone class="w-5 h-5" />
              </button>
              <button type="button" class="p-2 rounded-full bg-gray-50 dark:bg-gray-800 text-red-500 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                <Mail class="w-5 h-5" />
              </button>
            </div>
          </div>
        </form>

        <div class="mt-4 text-[10px] text-gray-400 text-center leading-relaxed">
          未注册手机验证后自动登录，注册即代表同意《知乎协议》《隐私保护指引》
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-scale-in {
  animation: scaleIn 0.2s ease-out;
}

.animate-shimmer {
  animation: shimmer 2s infinite linear;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
  }
}
</style>
