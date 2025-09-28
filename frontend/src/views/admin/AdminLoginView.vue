<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage, ElLoading } from 'element-plus';
import 'element-plus/es/components/message/style/css';
import 'element-plus/es/components/loading/style/css';

const router = useRouter();
const userStore = useUserStore();

// 表单数据
const form = ref({
  username: '',
  password: ''
});

// 表单规则
const rules = ref({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
});

// 表单引用
const formRef = ref<InstanceType<typeof ElForm>>();

// 登录处理
const handleLogin = async () => {
  // 表单验证
  if (!formRef.value) return;
  try {
    await formRef.value.validate();
  } catch (error) {
    return;
  }

  // 显示加载状态
  const loading = ElLoading.service({
    lock: true,
    text: '登录中...',
    background: 'rgba(0, 0, 0, 0.7)'
  });

  try {
    const result = await userStore.login(form.value);
    if (result.success && userStore.userInfo?.role === 'admin') {
      ElMessage.success('登录成功');
      router.push('/admin/dashboard');
    } else {
      ElMessage.error(result.message || '用户名或密码错误，或无管理员权限');
    }
  } catch (error) {
    console.error('登录失败:', error);
    ElMessage.error('登录失败，请检查网络连接');
  } finally {
    // 关闭加载状态
    loading.close();
  }
};

// 跳转到用户登录
const goToUserLogin = () => {
  router.push('/login');
};
</script>

<template>
  <div class="admin-login-container">
    <div class="login-wrapper">
      <div class="login-form">
        <div class="login-header">
          <h1 class="admin-title">管理后台</h1>
          <p class="login-subtitle">管理员登录</p>
        </div>
        
        <ElForm 
          ref="formRef"
          :model="form" 
          :rules="rules"
          class="form-container"
        >
          <ElFormItem label="用户名" prop="username">
            <ElInput 
              v-model="form.username" 
              placeholder="请输入管理员用户名"
              prefix-icon="User" 
            />
          </ElFormItem>
          
          <ElFormItem label="密码" prop="password">
            <ElInput 
              v-model="form.password" 
              type="password" 
              placeholder="请输入管理员密码"
              prefix-icon="Lock" 
              show-password
            />
          </ElFormItem>
          
          <ElFormItem>
            <ElButton 
              type="primary" 
              class="login-button"
              @click="handleLogin"
              :loading="false"
              native-type="button"
            >
              登录
            </ElButton>
          </ElFormItem>
        </ElForm>
        
        <div class="login-footer">
          <button class="link-button" @click="goToUserLogin">
            返回用户登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-login-container {
  width: 100%;
  height: 100vh;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.admin-login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url('data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"%3E%3Cg fill-rule="evenodd"%3E%3Cg fill="%23ffffff" fill-opacity="0.05"%3E%3Cpath opacity=".5" d="M96 95h4v1h-4v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4h-9v4h-1v-4H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15v-9H0v-1h15V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h9V0h1v15h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9h4v1h-4v9zm-1 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm9-10v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-10 0v-9h-9v9h9zm-9-10h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9zm10 0h9v-9h-9v9z"/%3E%3Cpath d="M6 5V0H5v5H0v1h5v94h1V6h94V5H6z"/%3E%3C/g%3E%3C/g%3E%3C/svg%3E');
}

.login-wrapper {
  width: 100%;
  max-width: 450px;
  padding: 20px;
}

.login-form {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 10px;
  padding: 40px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.admin-title {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 10px;
}

.login-subtitle {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.form-container {
  width: 100%;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  border: none;
}

.login-button:hover {
  background: linear-gradient(135deg, #2a5298 0%, #1e3c72 100%);
}

.login-footer {
  text-align: center;
  margin-top: 20px;
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
  .login-form {
    padding: 30px 20px;
  }
  
  .admin-title {
    font-size: 1.8rem;
  }
  
  .login-button {
    height: 40px;
    font-size: 15px;
  }
}
</style>