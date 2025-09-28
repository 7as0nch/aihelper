<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { Button, Swipe, SwipeItem, Card, Icon } from 'vant';

const router = useRouter();
const userStore = useUserStore();
const isLoggedIn = ref(false);

// 轮播图数据
const banners = [
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNjY2MiIHN0cm9rZS13aWR0aD0iMiIvPjxyZWN0IHg9IjEiIHk9IjEiIHdpZHRoPSI5OCIgaGVpZ2h0PSI0OCIgcng9IjIiIGZpbGw9IiNmZmYiIHN0cm9rZT0iI2NjYyIgc3Ryb2tlLXdpZHRoPSIxIi8+PC9zdmc+',
    title: '智能对话',
    desc: 'AI助手随时随地为您解答问题'
  },
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNjY2MiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik0yNSAxMGgyMHYyMGgtMjBWMXBiLTIuNzYgMC01IDIuMjQtNSA1czIuMjQgNSA1IDVjMi43NiAwIDUtMi4yNCA1LTVzLTIuMjQtNS01LTV6IiBzdHJva2U9IiM2NmUiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik01NSAxMGgyMHYyMGgtMjBWMXBiLTIuNzYgMC01IDIuMjQtNSA1czIuMjQgNSA1IDVjMi43NiAwIDUtMi4yNCA1LTVzLTIuMjQtNS01LTV6IiBzdHJva2U9IiM2NmUiIHN0cm9rZS13aWR0aD0iMiIvPjwvc3ZnPg==',
    title: '多功能工具',
    desc: '强大的AI工具提升您的工作效率'
  },
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNjY2MiIHN0cm9rZS13aWR0aD0iMiIvPjxjaXJjbGUgY3g9IjUwIiBjeT0iMjUiIHI9IjIwIiBmaWxsPSIjNjZlIiBzdHJva2U9IiM2NmUiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik0yMCAzMHYtMjBoMjB2MjBIMjB6TTYwIDMwdjIwSDQwdjIwSDYweiIgc3Ryb2tlPSIjNjZlIiBzdHJva2Utd2lkdGg9IjIiIGZpbGw9Im5vbmUiLz48L3N2Zz4=',
    title: '个性化体验',
    desc: '根据您的喜好定制AI服务'
  }
];

// 功能卡片数据
const features = [
  {
    icon: 'chat-o',
    title: '智能聊天',
    desc: '与AI进行自然语言对话'
  },
  {
    icon: 'edit-outline',
    title: '文本生成',
    desc: '快速生成各类文本内容'
  },
  {
    icon: 'image-o',
    title: '图像识别',
    desc: '分析图像内容获取信息'
  },
  {
    icon: 'code-o',
    title: '代码辅助',
    desc: '编写和优化编程代码'
  }
];

onMounted(() => {
  // 检查用户是否已登录
  isLoggedIn.value = userStore.isLoggedIn;
});

// 跳转到登录页面
const goToLogin = () => {
  router.push('/login');
};

// 跳转到聊天页面
const goToChat = () => {
  router.push('/chat');
};
</script>

<template>
  <div class="home-container">
    <!-- 头部 -->
    <header class="home-header">
      <div class="header-content">
        <h1 class="app-name">AI Chat</h1>
        <Button 
          v-if="!isLoggedIn"
          type="primary"
          class="login-button"
          @click="goToLogin"
        >
          登录/注册
        </Button>
        <Button 
          v-else
          type="primary"
          class="chat-button"
          @click="goToChat"
        >
          开始聊天
        </Button>
      </div>
    </header>

    <!-- 轮播图 -->
    <section class="banner-section">
      <Swipe class="banner-swipe" :autoplay="3000">
        <SwipeItem v-for="(item, index) in banners" :key="index">
          <div class="banner-item">
            <img :src="item.image" alt="Banner" class="banner-image" />
            <div class="banner-text">
              <h2 class="banner-title">{{ item.title }}</h2>
              <p class="banner-desc">{{ item.desc }}</p>
            </div>
          </div>
        </SwipeItem>
      </Swipe>
    </section>

    <!-- 功能特点 -->
    <section class="features-section">
      <h2 class="section-title">强大功能</h2>
      <div class="features-grid">
        <Card class="feature-card" v-for="(feature, index) in features" :key="index">
          <div class="feature-content">
            <div class="feature-icon">
              <Icon :name="feature.icon" size="36px" color="#667eea" />
            </div>
            <div class="feature-text">
              <h3 class="feature-title">{{ feature.title }}</h3>
              <p class="feature-desc">{{ feature.desc }}</p>
            </div>
          </div>
        </Card>
      </div>
    </section>

    <!-- 使用说明 -->
    <section class="guide-section">
      <h2 class="section-title">使用流程</h2>
      <div class="guide-steps">
        <div class="step-item">
          <div class="step-number">1</div>
          <div class="step-text">注册账号</div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step-item">
          <div class="step-number">2</div>
          <div class="step-text">登录系统</div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step-item">
          <div class="step-number">3</div>
          <div class="step-text">开始聊天</div>
        </div>
        <div class="step-arrow">→</div>
        <div class="step-item">
          <div class="step-number">4</div>
          <div class="step-text">获取结果</div>
        </div>
      </div>
    </section>

    <!-- 底部按钮 -->
    <section class="action-section">
      <Button 
        v-if="!isLoggedIn"
        type="primary"
        class="action-button"
        @click="goToLogin"
      >
        立即体验
      </Button>
      <Button 
        v-else
        type="primary"
        class="action-button"
        @click="goToChat"
      >
        继续聊天
      </Button>
    </section>

    <!-- 页脚 -->
    <footer class="home-footer">
      <p class="footer-text">AI Chat © 2023 版权所有</p>
    </footer>
  </div>
</template>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: var(--bg-primary);
}

/* 头部 */
.home-header {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-name {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0;
}

.login-button, .chat-button {
  height: 40px;
  font-size: 14px;
  padding: 0 20px;
  background-color: white;
  color: #667eea;
  border: none;
  border-radius: var(--radius-md);
  font-weight: 600;
  transition: all 0.3s ease;
}

.login-button:hover, .chat-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.3);
}

/* 轮播图 */
.banner-section {
  max-width: 1200px;
  margin: 30px auto;
  padding: 0 20px;
}

.banner-swipe {
  height: 200px;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-md);
}

.banner-item {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.banner-image {
  width: 100px;
  height: 100px;
}

.banner-text {
  flex: 1;
  margin-left: 30px;
}

.banner-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.banner-desc {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 功能特点 */
.features-section {
  max-width: 1200px;
  margin: 40px auto;
  padding: 0 20px;
}

.section-title {
  font-size: 1.8rem;
  font-weight: 600;
  text-align: center;
  color: var(--text-primary);
  margin-bottom: 30px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.feature-card {
  border: none;
  box-shadow: var(--shadow-sm);
  border-radius: var(--radius-lg);
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
}

.feature-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  text-align: center;
}

.feature-icon {
  margin-bottom: 15px;
}

.feature-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.feature-desc {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 使用说明 */
.guide-section {
  max-width: 1200px;
  margin: 40px auto;
  padding: 0 20px;
}

.guide-steps {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  margin-bottom: 8px;
}

.step-text {
  font-size: 0.9rem;
  color: var(--text-primary);
}

.step-arrow {
  color: var(--text-secondary);
  font-size: 1.2rem;
}

/* 底部按钮 */
.action-section {
  max-width: 1200px;
  margin: 40px auto;
  padding: 0 20px;
  text-align: center;
}

.action-button {
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  padding: 0 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: var(--radius-md);
  transition: all 0.3s ease;
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

/* 页脚 */
.home-footer {
  margin-top: 60px;
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
  .banner-item {
    flex-direction: column;
    text-align: center;
    padding: 20px;
  }
  
  .banner-image {
    margin-bottom: 20px;
  }
  
  .banner-text {
    margin-left: 0;
  }
  
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .guide-steps {
    flex-direction: column;
  }
  
  .step-arrow {
    transform: rotate(90deg);
  }
}

@media (max-width: 480px) {
  .app-name {
    font-size: 1.5rem;
  }
  
  .login-button, .chat-button {
    height: 36px;
    font-size: 13px;
    padding: 0 16px;
  }
  
  .banner-swipe {
    height: 180px;
  }
  
  .section-title {
    font-size: 1.5rem;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .action-button {
    height: 45px;
    font-size: 15px;
    padding: 0 30px;
  }
}
</style>