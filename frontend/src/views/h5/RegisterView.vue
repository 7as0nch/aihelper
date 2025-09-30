<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { Field, Button, showToast } from 'vant';

const router = useRouter();
const userStore = useUserStore();

// 表单数据
const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const email = ref('');
const loading = ref(false);

// 注册处理
const handleRegister = async () => {
  // 表单验证
  if (!username.value.trim()) {
    showToast('请输入用户名');
    return;
  }
  if (username.value.trim().length < 3) {
    showToast('用户名至少3个字符');
    return;
  }
  if (!password.value) {
    showToast('请输入密码');
    return;
  }
  if (password.value.length < 6) {
    showToast('密码至少6个字符');
    return;
  }
  if (password.value !== confirmPassword.value) {
    showToast('两次输入的密码不一致');
    return;
  }
  if (!email.value.trim()) {
    showToast('请输入邮箱');
    return;
  }
  // 简单的邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email.value.trim())) {
    showToast('请输入有效的邮箱地址');
    return;
  }

  loading.value = true;
  try {
    const result = await userStore.register({
      username: username.value.trim(),
      password: password.value,
      confirmPassword: confirmPassword.value,
      email: email.value.trim()
    });
    if (result.success) {
      showToast('注册成功，请登录');
      router.push('/login');
    } else {
      showToast(result.message || '注册失败，请重试');
    }
  } catch (error) {
    console.error('注册失败:', error);
    showToast('注册失败，请检查网络连接');
  } finally {
    loading.value = false;
  }
};

// 跳转到登录页面
const goToLogin = () => {
  router.push('/login');
};
</script>

<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h1 class="app-title">AI Chat</h1>
        <p class="register-subtitle">创建账号</p>
      </div>
      
      <div class="form-wrapper">
        <Field
          v-model="username"
          placeholder="请输入用户名"
          clearable
          class="form-field"
        />
        
        <Field
          v-model="email"
          placeholder="请输入邮箱"
          clearable
          type="email"
          class="form-field"
        />
        
        <Field
          v-model="password"
          type="password"
          placeholder="请输入密码"
          clearable
          password-visibility-toggle
          class="form-field"
        />
        
        <Field
          v-model="confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          clearable
          password-visibility-toggle
          class="form-field"
        />
        
        <Button
          type="primary"
          :loading="loading"
          class="register-button"
          @click="handleRegister"
        >
          注册
        </Button>
        
        <div class="register-footer">
          <p class="register-text">
            已有账号？
            <button class="link-button" @click="goToLogin">
              立即登录
            </button>
          </p>
        </div>
      </div>
    </div>
    
    <!-- <div class="register-bg" /> -->
  </div>
</template>

<style scoped>
.register-container {
  position: relative;
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #93d3fb 0%, #5791f5 100%);
}

.register-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"%3E%3Cg fill-rule="evenodd"%3E%3Cg fill="%23ffffff" fill-opacity="0.05"%3E%3Cpath opacity=".5" d="M96 95h4v1h-4v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9zm-1 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9z"/%3E%3Cpath d="M6 5V0H5v5H0v1h5v94h1V6h94V5H6z"/%3E%3C/g%3E%3C/g%3E%3C/svg%3E');
}

.register-card {
  position: relative;
  width: 100%;
  max-width: 400px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: var(--radius-lg);
  padding: 30px;
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(10px);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.app-title {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #93d3fb 0%, #5791f5 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 10px;
}

.register-subtitle {
  font-size: 1rem;
  color: var(--text-secondary);
}

.form-wrapper {
  width: 100%;
}

.form-field {
  margin-bottom: 20px;
}

.register-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #93d3fb 0%, #5791f5 100%);
  border: none;
  border-radius: var(--radius-md);
  transition: all 0.3s ease;
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(147, 211, 251, 0.4);
}

.register-footer {
  text-align: center;
  margin-top: 20px;
}

.register-text {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.link-button {
  background: none;
  border: none;
  color: var(--primary-color);
  font-size: 14px;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: var(--radius-sm);
  transition: background-color 0.3s ease;
}

.link-button:hover {
  background-color: rgba(59, 130, 246, 0.1);
}

/* 响应式设计 */
@media (max-width: 480px) {
  .register-card {
    padding: 20px;
  }
  
  .app-title {
    font-size: 1.8rem;
  }
  
  .register-button {
    height: 44px;
    font-size: 15px;
  }
}
</style>