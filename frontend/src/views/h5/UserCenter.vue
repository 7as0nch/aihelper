<template>
  <div class="user-center">
    <CommonHeader />
    <main class="main-content">
      <div class="user-profile">
        <div class="avatar">
          <img src="@/assets/avatar.png" alt="用户头像" />
        </div>
        <h2 class="username">张三</h2>
        <p class="user-id">ID: 123456789</p>
      </div>
      
      <div class="user-actions">
        <div class="action-card" @click="goToProfile">
          <i class="icon-profile"></i>
          <span>个人资料</span>
        </div>
        <div class="action-card" @click="goToSettings">
          <i class="icon-settings"></i>
          <span>设置</span>
        </div>
        <div class="action-card" @click="goToHistory">
          <i class="icon-history"></i>
          <span>对话历史</span>
        </div>
        <div class="action-card" @click="logout">
          <i class="icon-logout"></i>
          <span>退出登录</span>
        </div>
      </div>
    </main>
    <CommonFooter />
  </div>
</template>

<script setup lang="ts">
import CommonHeader from '../../components/common/Header.vue';
import CommonFooter from '../../components/common/Footer.vue';

const goToProfile = () => {
  // 跳转到个人资料页面
  console.log('跳转到个人资料页面');
};

const goToSettings = () => {
  // 跳转到设置页面
  console.log('跳转到设置页面');
};

const goToHistory = () => {
  // 跳转到对话历史页面
  console.log('跳转到对话历史页面');
};

const logout = () => {
  // 退出登录
  console.log('退出登录');
};
</script>

<style scoped>
.user-center {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.user-profile {
  text-align: center;
  padding: 2rem 1rem;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.avatar img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 1rem;
}

.username {
  font-size: 1.5rem;
  color: #495057;
  margin-bottom: 0.5rem;
}

.user-id {
  color: #6c757d;
  font-size: 0.875rem;
}

.user-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.action-card {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1.5rem 1rem;
  text-align: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s;
}

.action-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
}

.action-card i {
  font-size: 1.5rem;
  color: #42b983;
  margin-bottom: 0.5rem;
}

.action-card span {
  color: #495057;
  font-size: 0.875rem;
}
</style>