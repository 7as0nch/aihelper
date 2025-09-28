<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { Button, Cell, Dialog, showToast, Field } from 'vant';
import Avatar from 'vant';
import { updateUserInfo } from '../../api/user';

const router = useRouter();
const userStore = useUserStore();

// 用户信息
const userInfo = ref<any>({});
const editing = ref(false);
const editingInfo = ref<any>({});
const isLoading = ref(false);

// 初始化
onMounted(() => {
  // 检查用户是否已登录
  if (!userStore.isLoggedIn) {
    router.push('/login');
    return;
  }
  
  // 获取用户信息
  getUserInfo();
});

// 获取用户信息
const getUserInfo = async () => {
  try {
    await userStore.getUserInfo();
    if (userStore.userInfo) {
      userInfo.value = { ...userStore.userInfo };
      // 初始化编辑信息
      editingInfo.value = { ...userStore.userInfo };
    }
  } catch (error) {
    console.error('获取用户信息失败:', error);
  }
};

// 编辑用户信息
const handleEdit = () => {
  editing.value = true;
};

// 保存用户信息
const handleSave = async () => {
  isLoading.value = true;
  try {
    const result = await updateUserInfo(editingInfo.value);
    if (result.code === 200) {
      userInfo.value = { ...result.data };
      editing.value = false;
      showToast('保存成功');
      // 更新用户store中的信息
      await userStore.updateUserInfo(userInfo.value);
    } else {
      showToast(result.message || '保存失败，请重试');
    }
  } catch (error) {
    console.error('保存用户信息失败:', error);
    showToast('保存失败，请检查网络连接');
  } finally {
    isLoading.value = false;
  }
};

// 取消编辑
const handleCancel = () => {
  editingInfo.value = { ...userInfo.value };
  editing.value = false;
};

// 退出登录
const handleLogout = () => {
  Dialog.confirm({
    title: '退出登录',
    message: '确定要退出登录吗？',
  }).then(() => {
    userStore.logout();
    router.push('/login');
  });
};

// 跳转到聊天页面
const goToChat = () => {
  router.push('/chat');
};
</script>

<template>
  <div class="profile-container">
    <!-- 头部 -->
    <header class="profile-header">
      <div class="header-content">
        <h1 class="page-title">个人中心</h1>
        <Button 
          type="primary"
          size="small"
          class="chat-button"
          @click="goToChat"
        >
          去聊天
        </Button>
      </div>
    </header>

    <!-- 个人信息卡片 -->
    <section class="profile-card">
      <div class="profile-info">
        <Avatar class="profile-avatar" :icon="'user-o'" size="80px" />
        <div class="profile-details">
          <h2 class="profile-name">{{ userInfo.username || '未设置' }}</h2>
          <p class="profile-email">{{ userInfo.email || '未设置' }}</p>
        </div>
        <Button 
          type="primary"
          size="small"
          class="edit-button"
          @click="handleEdit"
          :disabled="editing"
        >
          编辑资料
        </Button>
      </div>
    </section>

    <!-- 用户信息表单 -->
    <section v-if="editing" class="profile-form">
      <h3 class="form-title">编辑个人资料</h3>
      <Field
        v-model="editingInfo.username"
        label="用户名"
        placeholder="请输入用户名"
        class="form-field"
      />
      <Field
        v-model="editingInfo.email"
        label="邮箱"
        type="email"
        placeholder="请输入邮箱"
        class="form-field"
      />
      <Field
        v-model="editingInfo.nickname"
        label="昵称"
        placeholder="请输入昵称（选填）"
        class="form-field"
      />
      <div class="form-actions">
        <Button 
          type="default"
          class="cancel-button"
          @click="handleCancel"
        >
          取消
        </Button>
        <Button 
          type="primary"
          class="save-button"
          @click="handleSave"
          :loading="isLoading"
        >
          保存
        </Button>
      </div>
    </section>

    <!-- 功能列表 -->
    <section class="profile-features">
      <h3 class="features-title">设置</h3>
      <div class="features-list">
        <Cell title="账户安全" icon="shield-o" is-link />
        <Cell title="消息通知" icon="bell-o" is-link />
        <Cell title="隐私设置" icon="lock" is-link />
        <Cell title="帮助与反馈" icon="question-circle-o" is-link />
        <Cell title="关于我们" icon="info-o" is-link />
      </div>
    </section>

    <!-- 退出登录按钮 -->
    <section class="logout-section">
      <Button 
        type="danger"
        class="logout-button"
        @click="handleLogout"
      >
        退出登录
      </Button>
    </section>

    <!-- 页脚 -->
    <footer class="profile-footer">
      <p class="footer-text">AI Chat © 2023 版权所有</p>
    </footer>
  </div>
</template>

<style scoped>
.profile-container {
  min-height: 100vh;
  background-color: var(--bg-primary);
}

/* 头部 */
.profile-header {
  padding: 15px 20px;
  background-color: white;
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.chat-button {
  height: 32px;
  font-size: 12px;
  padding: 0 16px;
}

/* 个人信息卡片 */
.profile-card {
  margin: 20px;
  padding: 25px;
  background-color: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.profile-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.profile-avatar {
  border: 2px solid var(--border-color);
}

.profile-details {
  flex: 1;
  margin-left: 20px;
}

.profile-name {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.profile-email {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0;
}

.edit-button {
  height: 36px;
  font-size: 13px;
  padding: 0 20px;
}

/* 编辑表单 */
.profile-form {
  margin: 0 20px 20px;
  padding: 25px;
  background-color: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.form-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
}

.form-field {
  margin-bottom: 20px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.cancel-button, .save-button {
  flex: 1;
  height: 40px;
}

/* 功能列表 */
.profile-features {
  margin: 0 20px 20px;
}

.features-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text-tertiary);
  margin-bottom: 10px;
  padding: 0 25px;
}

.features-list {
  background-color: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

/* 退出登录按钮 */
.logout-section {
  margin: 0 20px 40px;
  padding: 25px;
  background-color: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.logout-button {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 500;
}

/* 页脚 */
.profile-footer {
  margin-top: auto;
  padding: 20px;
  background-color: var(--bg-secondary);
  text-align: center;
}

.footer-text {
  font-size: 0.9rem;
  color: var(--text-tertiary);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .profile-info {
    flex-direction: column;
    text-align: center;
  }
  
  .profile-avatar {
    margin-bottom: 15px;
  }
  
  .profile-details {
    margin-left: 0;
    margin-bottom: 15px;
  }
  
  .form-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .profile-header {
    padding: 12px 15px;
  }
  
  .page-title {
    font-size: 1.1rem;
  }
  
  .profile-card, .profile-form, .logout-section {
    margin: 15px;
    padding: 20px;
  }
  
  .profile-avatar {
    size: 60px;
  }
  
  .profile-name {
    font-size: 1.2rem;
  }
  
  .logout-button {
    height: 40px;
    font-size: 14px;
  }
}
</style>