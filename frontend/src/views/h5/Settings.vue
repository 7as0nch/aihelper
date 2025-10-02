<template>
  <div class="settings">
    <CommonHeader />
    <main class="main-content">
      <div class="settings-section">
        <h2>账户设置</h2>
        <div class="setting-item">
          <span>修改密码</span>
          <i class="icon-arrow-right"></i>
        </div>
        <div class="setting-item">
          <span>绑定手机</span>
          <i class="icon-arrow-right"></i>
        </div>
        <div class="setting-item">
          <span>绑定邮箱</span>
          <i class="icon-arrow-right"></i>
        </div>
      </div>
      
      <div class="settings-section">
        <h2>隐私设置</h2>
        <div class="setting-item">
          <span>隐私保护</span>
          <a-switch :checked="privacyProtection" @change="(val) => privacyProtection = Boolean(val)" />
        </div>
        <div class="setting-item">
          <span>消息推送</span>
          <a-switch :checked="messagePush" @change="(val) => messagePush = Boolean(val)" />
        </div>
      </div>
      
      <div class="settings-section">
        <h2>其他设置</h2>
        <div class="setting-item">
          <span>清除缓存</span>
          <i class="icon-arrow-right"></i>
        </div>
        <div class="setting-item">
          <span>关于我们</span>
          <i class="icon-arrow-right"></i>
        </div>
        <div class="setting-item">
          <span>检查更新</span>
          <i class="icon-arrow-right"></i>
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
import { Switch as ASwitch } from 'ant-design-vue';

const privacyProtection = ref<boolean>(true);
const messagePush = ref<boolean>(false);
</script>

<style scoped>
.settings {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.settings-section {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin-bottom: 1.5rem;
}

.settings-section h2 {
  font-size: 1.125rem;
  color: #495057;
  margin-bottom: 1rem;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e9ecef;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item span {
  font-size: 0.875rem;
  color: #495057;
}

.setting-item i {
  color: #adb5bd;
}
</style>