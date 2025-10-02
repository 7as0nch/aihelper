<template>
  <div class="h5-home">
    <CommonHeader />
    <main class="main-content">
      <div class="banner">
        <h2>欢迎使用AI Chat</h2>
        <p>智能对话，触手可及</p>
      </div>
      <div class="features">
        <div class="feature-card">
          <i class="icon-chat"></i>
          <h3>智能对话</h3>
          <p>与AI进行自然语言对话</p>
        </div>
        <div class="feature-card">
          <i class="icon-search"></i>
          <h3>知识检索</h3>
          <p>快速获取所需信息</p>
        </div>
        <div class="feature-card">
          <i class="icon-translate"></i>
          <h3>多语言支持</h3>
          <p>支持多种语言互译</p>
        </div>
      </div>
    </main>
    <CommonFooter />
  </div>
</template>

<script setup lang="ts">
import CommonHeader from '../../components/common/Header.vue';
import CommonFooter from '../../components/common/Footer.vue';
</script>

<style scoped>
.h5-home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.banner {
  background: linear-gradient(135deg, #42b983, #42b983);
  color: white;
  padding: 2rem;
  border-radius: 8px;
  text-align: center;
  margin-bottom: 2rem;
}

.banner h2 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.banner p {
  font-size: 1rem;
  opacity: 0.9;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.feature-card {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-card i {
  font-size: 2rem;
  color: #42b983;
  margin-bottom: 1rem;
}

.feature-card h3 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: #495057;
}

.feature-card p {
  color: #6c757d;
  font-size: 0.875rem;
}
</style>