<template>
  <div class="notifications">
    <CommonHeader />
    <main class="main-content">
      <div class="notification-list">
        <div class="notification-item" v-for="notification in notifications" :key="notification.id">
          <div class="notification-icon">
            <i :class="notification.icon"></i>
          </div>
          <div class="notification-content">
            <h3>{{ notification.title }}</h3>
            <p>{{ notification.message }}</p>
            <span class="notification-time">{{ notification.time }}</span>
          </div>
        </div>
      </div>
    </main>
    <CommonFooter />
  </div>
</template>

<script setup lang="ts">
import CommonHeader from '../../components/common/Header.vue';
import CommonFooter from '../../components/common/Footer.vue';
import { ref } from 'vue';

interface Notification {
  id: number;
  title: string;
  message: string;
  time: string;
  icon: string;
}

const notifications = ref<Notification[]>([
  {
    id: 1,
    title: '系统通知',
    message: '您的账户已在新设备上登录',
    time: '2小时前',
    icon: 'icon-security'
  },
  {
    id: 2,
    title: '活动提醒',
    message: '您参与的活动即将开始',
    time: '1天前',
    icon: 'icon-calendar'
  },
  {
    id: 3,
    title: '消息通知',
    message: '您有一条新的私信',
    time: '2天前',
    icon: 'icon-message'
  }
]);
</script>

<style scoped>
.notifications {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.notification-item {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: flex-start;
}

.notification-icon {
  background-color: #e8f4f0;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 1rem;
}

.notification-icon i {
  color: #42b983;
  font-size: 1.25rem;
}

.notification-content {
  flex: 1;
}

.notification-content h3 {
  font-size: 1rem;
  color: #495057;
  margin-bottom: 0.25rem;
}

.notification-content p {
  font-size: 0.875rem;
  color: #6c757d;
  margin-bottom: 0.5rem;
}

.notification-time {
  font-size: 0.75rem;
  color: #adb5bd;
}
</style>