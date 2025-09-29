<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { Button, Swipe, SwipeItem, Card, Icon } from 'vant';

const router = useRouter();
const userStore = useUserStore();
const isLoggedIn = ref(false);

// 轮播图数据
const banners = [
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjxyZWN0IHg9IjEiIHk9IjEiIHdpZHRoPSI5OCIgaGVpZ2h0PSI0OCIgcng9IjIiIGZpbGw9IiNjY2MiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxIi8+PC9zdmc+',
    title: '智能对话',
    desc: 'AI助手随时随地为您解答问题'
  },
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik0yNSAxMGgyMHYyMGgtMjBWMXBiLTIuNzYgMC01IDIuMjQtNSA1czIuMjQgNSA1IDVjMi43NiAwIDUtMi4yNCA1LTVzLTIuMjQtNS01LTV6IiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik01NSAxMGgyMHYyMGgtMjBWMXBiLTIuNzYgMC01IDIuMjQtNSA1czIuMjQgNSA1IDVjMi43NiAwIDUtMi4yNCA1LTVzLTIuMjQtNS01LTV6IiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjwvc3ZnPg==',
    title: '多功能工具',
    desc: '强大的AI工具提升您的工作效率'
  },
  {
    image: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0wIDAgMTAwIDUwIiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjxjaXJjbGUgY3g9IjUwIiBjeT0iMjUiIHI9IjIwIiBmaWxsPSIjZmZmIiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMiIvPjxwYXRoIGQ9Ik0yMCAzMHYtMjBoMjB2MjBIMjB6TTYwIDMwdjIwSDQwdjIwSDYweiIgc3Ryb2tlPSIjZmZmIiBzdHJva2Utd2lkdGg9IjIiIGZpbGw9Im5vbmUiLz48L3N2Zz4=',
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
  
  // 初始化粒子背景
  initParticleBackground();
});

// 初始化粒子背景
const initParticleBackground = () => {
  nextTick(() => {
    const canvas = document.getElementById('particle-bg') as HTMLCanvasElement;
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    
    // 设置canvas尺寸
    const resizeCanvas = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };
    
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);
    
    // 创建粒子
    const particlesArray: any[] = [];
    const numberOfParticles = Math.floor(canvas.width * canvas.height / 5000);
    
    class Particle {
      x: number;
      y: number;
      size: number;
      speedX: number;
      speedY: number;
      
      constructor() {
        this.x = Math.random() * canvas.width;
        this.y = Math.random() * canvas.height;
        this.size = Math.random() * 2 + 0.5;
        this.speedX = Math.random() * 0.5 - 0.25;
        this.speedY = Math.random() * 0.5 - 0.25;
      }
      
      update() {
        this.x += this.speedX;
        this.y += this.speedY;
        
        if (this.size > 0.2) this.size -= 0.001;
        
        if (this.x < 0 || this.x > canvas.width) this.speedX *= -1;
        if (this.y < 0 || this.y > canvas.height) this.speedY *= -1;
      }
      
      draw() {
        if (!ctx) return;
        ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
        ctx.beginPath();
        ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
        ctx.fill();
      }
    }
    
    function init() {
      for (let i = 0; i < numberOfParticles; i++) {
        particlesArray.push(new Particle());
      }
    }
    
    function animate() {
      if (!ctx) return;
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      for (let i = 0; i < particlesArray.length; i++) {
        particlesArray[i].update();
        particlesArray[i].draw();
      }
      connect();
      requestAnimationFrame(animate);
    }
    
    function connect() {
      if (!ctx) return;
      for (let a = 0; a < particlesArray.length; a++) {
        for (let b = a; b < particlesArray.length; b++) {
          const dx = particlesArray[a].x - particlesArray[b].x;
          const dy = particlesArray[a].y - particlesArray[b].y;
          const distance = Math.sqrt(dx * dx + dy * dy);
          
          if (distance < 100) {
            ctx.beginPath();
            ctx.strokeStyle = `rgba(255, 255, 255, ${0.2 - distance/500})`;
            ctx.lineWidth = 0.5;
            ctx.moveTo(particlesArray[a].x, particlesArray[a].y);
            ctx.lineTo(particlesArray[b].x, particlesArray[b].y);
            ctx.stroke();
          }
        }
      }
    }
    
    init();
    animate();
  });
};

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
    <!-- 粒子背景 -->
    <canvas id="particle-bg" class="particle-background"></canvas>
    
    <!-- 头部 -->
    <header class="home-header">
      <div class="header-content">
        <h1 class="app-name">
          <span class="text-gradient">AI Chat</span>
        </h1>
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

    <!-- 主内容 -->
    <main class="main-content">
      <!-- 轮播图 -->
      <section class="banner-section">
        <Swipe class="banner-swipe" :autoplay="3000">
          <SwipeItem v-for="(item, index) in banners" :key="index" class="banner-swipe-item">
            <div class="banner-item">
              <div class="banner-image-container">
                <img :src="item.image" alt="Banner" class="banner-image" />
                <div class="image-glow"></div>
              </div>
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
        <h2 class="section-title text-center">强大功能</h2>
        <div class="features-grid">
          <Card class="feature-card" v-for="(feature, index) in features" :key="index">
            <div class="feature-content">
              <div class="feature-icon">
                <Icon :name="feature.icon" size="36px" color="#ffffff" />
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
        <h2 class="section-title text-center">使用流程</h2>
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
    </main>

    <!-- 页脚 -->
    <footer class="home-footer">
      <p class="footer-text">AI Chat © 2023 版权所有</p>
    </footer>
  </div>
</template>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-glow: rgba(102, 126, 234, 0.5);
  --secondary-glow: rgba(118, 75, 162, 0.5);
  --bg-dark: #0f1216;
  --bg-darker: #0a0d11;
  --text-primary: #ffffff;
  --text-secondary: #a0aec0;
}

.home-container {
  min-height: 100vh;
  background-color: var(--bg-dark);
  position: relative;
  overflow: hidden;
}

/* 粒子背景 */
.particle-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  background: radial-gradient(ellipse at center, #1a1d23 0%, var(--bg-darker) 100%);
}

/* 主内容区域 */
.main-content {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 头部 */
.home-header {
  position: relative;
  z-index: 2;
  padding: 30px 20px;
  backdrop-filter: blur(10px);
  background-color: rgba(15, 18, 22, 0.8);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-name {
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0;
  letter-spacing: -0.05em;
}

.text-gradient {
  background: linear-gradient(90deg, #ffffff 0%, #a0aec0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-fill-color: transparent;
  position: relative;
}

.text-gradient::after {
  content: attr(data-text);
  position: absolute;
  left: 2px;
  top: 2px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-fill-color: transparent;
  z-index: -1;
  opacity: 0.5;
}

.login-button, .chat-button {
  height: 48px;
  font-size: 16px;
  padding: 0 24px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 24px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.3);
  position: relative;
  overflow: hidden;
}

.login-button::before, .chat-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: all 0.6s ease;
}

.login-button:hover::before, .chat-button:hover::before {
  left: 100%;
}

.login-button:hover, .chat-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.5);
  letter-spacing: 0.5px;
}

/* 轮播图 */
.banner-section {
  margin: 60px auto;
  position: relative;
}

.banner-swipe {
  height: 240px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
}

.banner-swipe-item {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  backdrop-filter: blur(5px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.banner-item {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 60px;
  position: relative;
}

.banner-image-container {
  position: relative;
  z-index: 1;
}

.banner-image {
  width: 120px;
  height: 120px;
  animation: pulse 4s infinite ease-in-out;
}

.image-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, var(--primary-glow) 0%, transparent 70%);
  border-radius: 50%;
  z-index: -1;
  animation: glowPulse 3s infinite ease-in-out;
}

.banner-text {
  flex: 1;
  margin-left: 60px;
  max-width: 400px;
}

.banner-title {
  font-size: 2.2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 15px;
  line-height: 1.2;
  animation: slideUp 0.5s ease-out;
}

.banner-desc {
  font-size: 1.1rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
  animation: slideUp 0.5s ease-out 0.2s both;
}

/* 功能特点 */
.features-section {
  margin: 80px auto;
}

.section-title {
  font-size: 2.2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 40px;
  position: relative;
  display: inline-block;
}

.section-title::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 80px;
  height: 4px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
}

.text-center {
  text-align: center;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
}

.feature-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
  border-radius: 16px;
  transition: all 0.4s ease;
  overflow: hidden;
  position: relative;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transform: scaleX(0);
  transition: transform 0.4s ease;
}

.feature-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(102, 126, 234, 0.3);
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.15);
}

.feature-card:hover::before {
  transform: scaleX(1);
}

.feature-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px;
  text-align: center;
}

.feature-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
}

.feature-card:hover .feature-icon {
  transform: scale(1.1) rotate(5deg);
  box-shadow: 0 12px 25px rgba(102, 126, 234, 0.6);
}

.feature-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
  transition: color 0.3s ease;
}

.feature-card:hover .feature-title {
  color: #667eea;
}

.feature-desc {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
}

/* 使用说明 */
.guide-section {
  margin: 80px auto;
}

.guide-steps {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 30px;
  margin-top: 40px;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  transition: transform 0.3s ease;
}

.step-item:hover {
  transform: translateY(-5px);
}

.step-number {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 1.2rem;
  margin-bottom: 15px;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
  position: relative;
  z-index: 1;
  transition: all 0.3s ease;
}

.step-item:hover .step-number {
  transform: scale(1.1);
  box-shadow: 0 12px 25px rgba(102, 126, 234, 0.6);
}

.step-text {
  font-size: 1rem;
  color: var(--text-primary);
  font-weight: 500;
  text-align: center;
  padding: 0 10px;
}

.step-arrow {
  color: var(--text-secondary);
  font-size: 1.5rem;
  font-weight: 300;
  opacity: 0.6;
}

/* 底部按钮 */
.action-section {
  margin: 80px auto;
  text-align: center;
  position: relative;
}

.action-button {
  height: 60px;
  font-size: 1.1rem;
  font-weight: 600;
  padding: 0 60px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 30px;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
  position: relative;
  overflow: hidden;
  letter-spacing: 0.5px;
}

.action-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: all 0.8s ease;
}

.action-button:hover::before {
  left: 100%;
}

.action-button:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(102, 126, 234, 0.6);
  letter-spacing: 1px;
}

/* 页脚 */
.home-footer {
  margin-top: 100px;
  padding: 30px 20px;
  background-color: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  position: relative;
  z-index: 1;
}

.footer-text {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0;
  opacity: 0.7;
}

/* 动画效果 */
@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes glowPulse {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.6;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.1);
    opacity: 0.8;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

/* 响应式设计 */
@media (max-width: 992px) {
  .banner-item {
    padding: 0 40px;
  }
  
  .banner-text {
    margin-left: 40px;
  }
  
  .banner-title {
    font-size: 1.8rem;
  }
}

@media (max-width: 768px) {
  .home-header {
    padding: 20px 15px;
  }
  
  .app-name {
    font-size: 2rem;
  }
  
  .login-button, .chat-button {
    height: 40px;
    font-size: 14px;
    padding: 0 20px;
  }
  
  .banner-swipe {
    height: 200px;
  }
  
  .banner-item {
    flex-direction: column;
    text-align: center;
    padding: 30px 20px;
  }
  
  .banner-image {
    width: 100px;
    height: 100px;
    margin-bottom: 20px;
  }
  
  .banner-text {
    margin-left: 0;
  }
  
  .banner-title {
    font-size: 1.6rem;
  }
  
  .section-title {
    font-size: 1.8rem;
  }
  
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }
  
  .guide-steps {
    flex-direction: column;
    gap: 20px;
  }
  
  .step-arrow {
    transform: rotate(90deg);
    margin: 10px 0;
  }
  
  .action-button {
    height: 50px;
    font-size: 1rem;
    padding: 0 40px;
  }
}

@media (max-width: 480px) {
  .app-name {
    font-size: 1.8rem;
  }
  
  .banner-swipe {
    height: 180px;
  }
  
  .section-title {
    font-size: 1.6rem;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .feature-card {
    padding: 20px;
  }
  
  .feature-icon {
    width: 60px;
    height: 60px;
  }
}
</style>