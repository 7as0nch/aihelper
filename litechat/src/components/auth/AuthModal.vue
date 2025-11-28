<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../stores/auth';
import { X, QrCode, Smartphone, MessageCircle, Mail, User, ArrowLeft } from 'lucide-vue-next';

const authStore = useAuthStore();
const view = ref<'login' | 'register'>('login');
const loginType = ref<'phone' | 'account'>('phone');

// Form states
const phone = ref('');
const code = ref('');
const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const isLoading = ref(false);
const errorMessage = ref('');

const resetForm = () => {
  phone.value = '';
  code.value = '';
  username.value = '';
  password.value = '';
  confirmPassword.value = '';
  errorMessage.value = '';
};

const switchView = (newView: 'login' | 'register') => {
  view.value = newView;
  resetForm();
};

const handleLogin = async () => {
  if (isLoading.value) return;
  errorMessage.value = '';
  isLoading.value = true;

  try {
    let success = false;
    if (loginType.value === 'phone') {
      if (!phone.value || !code.value) {
        errorMessage.value = '请输入手机号和验证码';
        isLoading.value = false;
        return;
      }
      success = await authStore.loginWithPhone(phone.value, code.value);
    } else {
      if (!username.value || !password.value) {
        errorMessage.value = '请输入用户名和密码';
        isLoading.value = false;
        return;
      }
      success = await authStore.loginWithPassword(username.value, password.value);
    }

    if (!success) {
      errorMessage.value = '登录失败，请检查凭证';
    }
  } catch (e) {
    errorMessage.value = '发生错误，请重试';
  } finally {
    isLoading.value = false;
  }
};

const handleRegister = async () => {
  if (isLoading.value) return;
  errorMessage.value = '';
  
  if (!username.value || !password.value || !confirmPassword.value) {
    errorMessage.value = '请填写所有字段';
    return;
  }

  if (password.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致';
    return;
  }

  isLoading.value = true;
  try {
    const success = await authStore.register(username.value, password.value);
    if (!success) {
      errorMessage.value = '注册失败，请重试';
    }
  } catch (e) {
    errorMessage.value = '发生错误，请重试';
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div v-if="authStore.showAuthModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm animate-fade-in">
    <div class="bg-white dark:bg-[#2a2a2a] rounded-2xl shadow-2xl w-full max-w-3xl overflow-hidden flex relative animate-scale-in min-h-[500px]">
      <!-- Close Button -->
      <button 
        @click="authStore.closeModal()"
        class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors z-10"
      >
        <X class="w-6 h-6" />
      </button>

      <!-- Left Side: QR Code (Only show if enabled and in login view) -->
      <div 
        v-if="authStore.config.enableQrLogin && view === 'login'" 
        class="w-1/2 bg-gray-50 dark:bg-[#242424] p-8 flex flex-col items-center justify-center border-r border-gray-100 dark:border-gray-800 hidden md:flex"
      >
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">打开 App</h3>
        <p class="text-sm text-gray-500 mb-8">在「我的页」右上角打开扫一扫</p>
        
        <div class="bg-white p-4 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 mb-4">
          <!-- Mock QR Code -->
          <div class="w-40 h-40 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center relative overflow-hidden">
            <QrCode class="w-24 h-24 text-primary opacity-80" />
            <div class="absolute inset-0 bg-gradient-to-tr from-transparent via-white/20 to-transparent animate-shimmer"></div>
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="w-10 h-10 bg-white rounded-full flex items-center justify-center shadow-md">
                <span class="text-primary font-bold text-xs">AI</span>
              </div>
            </div>
          </div>
        </div>
        
        <p class="text-sm text-gray-500 flex items-center gap-2">
          其他扫码方式：微信
        </p>
      </div>
      
      <!-- Left Side: Register/Account Info (When QR is hidden or in register view) -->
      <div 
        v-else 
        class="w-1/2 bg-primary/5 dark:bg-primary/10 p-8 flex flex-col items-center justify-center border-r border-gray-100 dark:border-gray-800 hidden md:flex text-center"
      >
        <div class="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mb-6">
          <User class="w-10 h-10 text-primary" />
        </div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">
          {{ view === 'register' ? '加入 AI Chat' : '欢迎回来' }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400 max-w-xs">
          {{ view === 'register' ? '注册账号以保存您的对话历史和偏好设置' : '登录以继续您的 AI 之旅' }}
        </p>
      </div>

      <!-- Right Side: Form -->
      <div class="w-full md:w-1/2 p-8 md:p-12 flex flex-col">
        
        <!-- Login View -->
        <template v-if="view === 'login'">
          <div class="flex gap-6 mb-8 border-b border-gray-100 dark:border-gray-800">
            <button 
              class="pb-3 text-sm font-medium transition-colors relative"
              :class="loginType === 'phone' ? 'text-primary' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
              @click="loginType = 'phone'"
            >
              验证码登录
              <div v-if="loginType === 'phone'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full"></div>
            </button>
            <button 
              class="pb-3 text-sm font-medium transition-colors relative"
              :class="loginType === 'account' ? 'text-primary' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
              @click="loginType = 'account'"
            >
              密码登录
              <div v-if="loginType === 'account'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full"></div>
            </button>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4 flex-1">
            <!-- Phone Login Form -->
            <template v-if="loginType === 'phone'">
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
                    type="text" 
                    placeholder="输入 6 位短信验证码" 
                    class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                  />
                  <button 
                    type="button"
                    class="text-xs text-primary hover:text-blue-600 font-medium"
                  >
                    获取短信验证码
                  </button>
                </div>
              </div>
            </template>

            <!-- Account Login Form -->
            <template v-else>
              <div class="space-y-4">
                <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
                  <input 
                    v-model="username"
                    type="text" 
                    placeholder="用户名 / 邮箱" 
                    class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                  />
                </div>
                
                <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
                  <input 
                    v-model="password"
                    type="password" 
                    placeholder="密码" 
                    class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                  />
                </div>
              </div>
              
              <div class="flex justify-between items-center text-xs">
                <button type="button" class="text-gray-500 hover:text-gray-700 dark:text-gray-400">忘记密码?</button>
                <button type="button" @click="switchView('register')" class="text-primary hover:underline">注册账号</button>
              </div>
            </template>

            <div v-if="errorMessage" class="text-red-500 text-xs text-center">{{ errorMessage }}</div>

            <button 
              type="submit"
              class="w-full bg-primary text-white py-2.5 rounded-lg font-medium hover:bg-blue-600 transition-colors disabled:opacity-50 mt-6"
              :disabled="isLoading"
            >
              {{ isLoading ? '登录中...' : (loginType === 'phone' ? '登录/注册' : '登录') }}
            </button>
          </form>
        </template>

        <!-- Register View -->
        <template v-else>
          <div class="mb-8">
            <button 
              @click="switchView('login')"
              class="flex items-center text-sm text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white transition-colors mb-4"
            >
              <ArrowLeft class="w-4 h-4 mr-1" />
              返回登录
            </button>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white">注册账号</h2>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-4 flex-1">
            <div class="space-y-4">
              <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
                <input 
                  v-model="username"
                  type="text" 
                  placeholder="用户名" 
                  class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                />
              </div>
              
              <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
                <input 
                  v-model="password"
                  type="password" 
                  placeholder="设置密码" 
                  class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                />
              </div>

              <div class="flex items-center border-b border-gray-200 dark:border-gray-700 py-2">
                <input 
                  v-model="confirmPassword"
                  type="password" 
                  placeholder="确认密码" 
                  class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white placeholder-gray-400"
                />
              </div>
            </div>

            <div v-if="errorMessage" class="text-red-500 text-xs text-center">{{ errorMessage }}</div>

            <button 
              type="submit"
              class="w-full bg-primary text-white py-2.5 rounded-lg font-medium hover:bg-blue-600 transition-colors disabled:opacity-50 mt-6"
              :disabled="isLoading"
            >
              {{ isLoading ? '注册中...' : '注册' }}
            </button>
          </form>
        </template>

        <!-- Social Login (Only on Login View) -->
        <div v-if="view === 'login'" class="mt-8">
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

        <div class="mt-auto pt-4 text-[10px] text-gray-400 text-center leading-relaxed">
          {{ view === 'login' ? '未注册手机验证后自动登录，' : '' }}注册即代表同意《LiteChat 协议》《隐私保护指引》
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
