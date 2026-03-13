<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const redirectUrl = ref('');

onMounted(() => {
  // Get redirect URL from query parameters
  redirectUrl.value = route.query.redirect as string || '/';
});

const handleMockLogin = (type: 'success' | 'fail') => {
  if (!redirectUrl.value) return;

  const url = new URL(redirectUrl.value);
  
  if (type === 'success') {
    // Generate a mock token
    const mockToken = 'mock_qq_token_' + Math.random().toString(36).substring(7);
    url.searchParams.set('qq_token', mockToken);
  } else {
    url.searchParams.set('qq_error', 'mock_login_failed');
  }

  // Redirect back
  window.location.href = url.toString();
};
</script>

<template>
  <div class="min-h-screen bg-[#f5f6fa] flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-lg w-full max-w-md overflow-hidden">
      <!-- QQ Header -->
      <div class="bg-[#12b7f5] p-6 text-center">
        <h1 class="text-white text-2xl font-bold tracking-wider">QQ 登录</h1>
        <p class="text-blue-100 mt-2 text-sm">模拟授权页面</p>
      </div>

      <!-- Content -->
      <div class="p-8">
        <div class="flex flex-col items-center mb-8">
          <div class="w-20 h-20 bg-gray-100 rounded-full mb-4 flex items-center justify-center">
            <span class="text-4xl">🐧</span>
          </div>
          <h2 class="text-xl font-medium text-gray-800">AI Chat</h2>
          <p class="text-gray-500 text-sm mt-2 text-center">
            该应用将获取您的公开信息（昵称、头像等）
          </p>
        </div>

        <div class="space-y-4">
          <button 
            @click="handleMockLogin('success')"
            class="w-full bg-[#12b7f5] hover:bg-[#0ea4db] text-white py-3 rounded-xl font-medium transition-colors shadow-sm"
          >
            同意授权并登录
          </button>
          
          <button 
            @click="handleMockLogin('fail')"
            class="w-full bg-gray-100 hover:bg-gray-200 text-gray-700 py-3 rounded-xl font-medium transition-colors"
          >
            拒绝授权 (模拟失败)
          </button>
        </div>
      </div>
      
      <!-- Footer -->
      <div class="bg-gray-50 p-4 text-center text-xs text-gray-400">
        此页面仅用于本地开发和演示环境的 QQ 登录模拟
      </div>
    </div>
  </div>
</template>
