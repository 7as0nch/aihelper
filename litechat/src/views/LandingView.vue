<script setup lang="ts">
import { computed, ref, type Component } from 'vue';
import { RouterLink } from 'vue-router';
import {
  ArrowRight,
  Bot,
  BrainCircuit,
  Camera,
  ExternalLink,
  Github,
  Globe,
  LayoutDashboard,
  Mail,
  MessageSquare,
  Mic,
  Moon,
  Search,
  Sparkles,
  Sun,
  Workflow,
} from 'lucide-vue-next';

interface FeatureCard {
  title: string;
  description: string;
  icon: Component;
}

interface StackCard {
  title: string;
  description: string;
  icon: Component;
}

const isDark = ref(true);
const currentYear = new Date().getFullYear();

const companySiteLink = 'https://aihelper.chat/';
const githubLink = 'https://github.com/7as0nch/litechat';
const icp = '蜀ICP备2026011723号-1';
const icpLink = 'https://beian.miit.gov.cn';
const contactEmail = '7as0nch@gmail.com';

const featureCards: FeatureCard[] = [
  {
    title: '普通聊天与深度搜索切换',
    description: '面向不同复杂度问题，在同一个助手入口中完成轻量对话与更深入的搜索式推理。',
    icon: MessageSquare,
  },
  {
    title: '深度搜索聊天',
    description: '围绕更复杂的问题理解、资料整合与多步检索，提供更适合复杂任务的对话方式。',
    icon: Search,
  },
  {
    title: '截图理解',
    description: '支持将屏幕截图作为上下文输入，让 AI 助手更容易理解页面问题和视觉信息。',
    icon: Camera,
  },
  {
    title: '网页搜索',
    description: '支持联网搜索，让回答可以结合当前网页信息与在线资料进一步补充。',
    icon: Globe,
  },
  {
    title: '语音输入',
    description: '支持语音转文字输入，减少手动输入负担，提升日常提问与记录效率。',
    icon: Mic,
  },
  {
    title: '多 Agent 协作扩展',
    description: '后端基于 Go 与 Eino 继续推进多 Agent 协作，为更复杂业务场景预留空间。',
    icon: Workflow,
  },
];

const stackCards: StackCard[] = [
  {
    title: 'Eino 编排引擎',
    description: '基于 Golang 生态的 Eino AI 框架，承接模型调用、工作流和多 Agent 协作逻辑。',
    icon: BrainCircuit,
  },
  {
    title: 'Vben 管理后台',
    description: 'BackWeb 基于 Vben 二次开发，已支持模型配置，并适合作为统一的后台治理入口。',
    icon: LayoutDashboard,
  },
  {
    title: 'LiteChat Web + Widget',
    description: '官网前台与聊天页面可以独立运行，也支持作为插件形式植入到其他任意网页中。',
    icon: Sparkles,
  },
  {
    title: '免费开源与持续扩展',
    description: '当前项目面向所有用户免费开源，目标是减少后续使用者在 AI 助手前端上的重复开发。',
    icon: Bot,
  },
];

const themeClass = computed(() => (isDark.value ? 'theme-dark' : 'theme-light'));

const toggleTheme = () => {
  isDark.value = !isDark.value;
};
</script>

<template>
  <div class="landing-page" :class="themeClass">
    <div class="background-layer">
      <div class="glow glow-left" />
      <div class="glow glow-right" />
      <div class="grid-overlay" />
    </div>

    <nav class="top-nav">
      <div class="brand">
        <img class="brand-logo" src="/logo.png" alt="LiteChat" />
        <span class="brand-text">LiteChat</span>
      </div>

      <div class="nav-links">
        <a href="#features">Features</a>
        <a href="#tech">Eino &amp; Vben Tech</a>
        <a href="#open-source">Open Source</a>
        <a :href="githubLink" target="_blank" rel="noreferrer">Docs</a>
      </div>

      <div class="nav-actions">
        <button class="theme-toggle" type="button" @click="toggleTheme" :aria-label="isDark ? '切换亮色主题' : '切换暗色主题'">
          <Sun v-if="isDark" :size="18" />
          <Moon v-else :size="18" />
        </button>
        <RouterLink to="/app" class="primary-nav-button">
          Get Started
        </RouterLink>
      </div>
    </nav>

    <main>
      <section class="hero-section">
        <div class="hero-copy">
          <div class="status-pill">
            <span class="status-dot" />
            <span>Singularity v2.4 Online</span>
          </div>

          <h1>
            你的专属 <br>
            <span>助手</span>
          </h1>

          <p>
            LiteChat 是一套面向所有人的开源 AI 助手前端方案，目标是把官网、聊天入口、插件植入、后台配置与工程复用统一起来。
            当前已支持普通聊天、深度搜索聊天、截图、网页搜索、语音输入，并在 BackWeb 中支持模型配置。
          </p>

          <div class="hero-buttons">
            <RouterLink to="/app" class="hero-primary">
              Launch LiteChat
              <ArrowRight :size="18" />
            </RouterLink>
            <a class="hero-secondary" :href="githubLink" target="_blank" rel="noreferrer">
              <Github :size="18" />
              View on GitHub
            </a>
          </div>

          <div class="hero-meta">
            <div>
              <span>CORE ENGINE</span>
              <strong>Eino Orchestrator</strong>
            </div>
            <div class="divider" />
            <div>
              <span>UI FRAMEWORK</span>
              <strong>Vben Admin</strong>
            </div>
          </div>
        </div>

        <div class="hero-visual">
          <div class="hero-orb-panel glass-panel neon-glow">
            <div class="orb-shell">
              <div class="orb-ring outer" />
              <div class="orb-ring inner" />
              <div class="orb-core">
                <img src="/logo.png" alt="LiteChat Orb" />
              </div>
            </div>
          </div>

          <div class="floating-card top-card glass-panel">
            <BrainCircuit :size="30" />
            <span>Eino</span>
          </div>

          <div class="floating-card bottom-card glass-panel">
            <LayoutDashboard :size="30" />
            <span>Vben</span>
          </div>
        </div>
      </section>

      <section id="features" class="section-block section-surface">
        <div class="section-heading center">
          <h2>Evolved Intelligence</h2>
          <p>围绕真实 AI 助手场景设计的高频能力，优先解决接入、使用与复用问题。</p>
        </div>

        <div class="feature-grid">
          <article v-for="card in featureCards" :key="card.title" class="feature-card">
            <div class="feature-icon">
              <component :is="card.icon" :size="28" />
            </div>
            <h3>{{ card.title }}</h3>
            <p>{{ card.description }}</p>
          </article>
        </div>
      </section>

      <section id="tech" class="section-block orchestration-section">
        <div class="orchestration-copy">
          <h2>Neural Orchestration</h2>
          <p>
            在后端侧，<strong>Eino</strong> 负责模型协同、工作流编排和多 Agent 能力延展。
            在前端侧，LiteChat 则把聊天、搜索、截图、语音输入和插件植入整合为一个清晰的产品入口。
          </p>

          <div class="stack-list">
            <div class="stack-item">
              <Workflow :size="20" />
              <div>
                <h4>Smart Routing</h4>
                <p>根据任务复杂度与产品场景，组织更适合的模型能力和调用链路。</p>
              </div>
            </div>

            <div class="stack-item secondary">
              <Sparkles :size="20" />
              <div>
                <h4>Reusable Delivery Layer</h4>
                <p>让官网、聊天页面和嵌入式组件共享统一的交付思路，减少重复建设。</p>
              </div>
            </div>
          </div>
        </div>

        <div class="orchestration-visual glass-panel">
          <div class="orchestration-core">
            <div class="core-circle outer">
              <div class="core-circle inner">
                <div class="core-node">EINO</div>
              </div>
            </div>
            <div class="orbit-node orbit-top">AI</div>
            <div class="orbit-node orbit-left">Web</div>
            <div class="orbit-node orbit-right">Ops</div>
          </div>
        </div>
      </section>

      <section class="section-block section-surface admin-section">
        <div class="section-heading center">
          <h2>Enterprise Customization</h2>
          <p>BackWeb 负责模型配置与后续治理能力承接，让前台官网不再是脱离后台的单独页面。</p>
        </div>

        <div class="admin-mockup">
          <div class="mockup-topbar">
            <span />
            <span />
            <span />
            <div class="mockup-title">litechat-admin-panel.vben</div>
          </div>

          <div class="mockup-body">
            <aside class="mockup-sidebar">
              <div class="sidebar-line" />
              <div class="sidebar-line short" />
              <div class="sidebar-line active" />
              <div class="sidebar-line medium" />
            </aside>

            <div class="mockup-main">
              <div class="mockup-stats">
                <div v-for="card in stackCards.slice(0, 3)" :key="card.title" class="stat-card glass-panel">
                  <component :is="card.icon" :size="20" />
                  <span>{{ card.title }}</span>
                </div>
              </div>

              <div class="dashboard-preview glass-panel">
                <div class="preview-grid">
                  <div class="preview-panel large" />
                  <div class="preview-panel" />
                  <div class="preview-panel" />
                  <div class="preview-panel" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="open-source" class="section-block mission-section">
        <div class="mission-visual">
          <div class="future-card glass-panel">
            <div class="future-label">CRAWFISH PROTOCOL</div>
            <div class="future-content">
              <Bot :size="88" />
              <h3>怎么养一个可靠有趣的小龙虾</h3>
              <p>这会作为后续可期待的新服务方向之一，与“你的专属贾维斯”一起继续拓展产品边界。</p>
            </div>
          </div>
        </div>

        <div class="mission-copy">
          <h2>More Than Just an <span>AI.</span></h2>
          <p>
            LiteChat 的核心目标，是减少 AI 产品落地过程中的重复前端开发。
            它不是一次性展示页，而是一个面向官网、聊天入口、嵌入式交付和后续业务扩展的开源基础层。
          </p>

          <ul>
            <li>当前面向所有用户免费开源。</li>
            <li>原生支持 Web 与插件式接入。</li>
            <li>支持基于 Go + Eino 持续扩展多 Agent 协作。</li>
            <li>可继续延展到“专属 Jarvis”等更多长期服务。</li>
          </ul>

          <a :href="githubLink" target="_blank" rel="noreferrer" class="roadmap-link">
            Explore our Roadmap
            <ExternalLink :size="18" />
          </a>
        </div>
      </section>
    </main>

    <footer class="site-footer">
      <div class="footer-grid">
        <div class="footer-brand">
          <span class="footer-title">LiteChat</span>
          <p>
            面向 AI 助手产品的开源官网与前端交付方案，基于 Golang、Eino 与 Vben 持续演进。
          </p>
          <div class="footer-badge">BUILT WITH EINO</div>
        </div>

        <div class="footer-column">
          <h4>Ecosystem</h4>
          <a href="#features">Features</a>
          <a href="#tech">Eino &amp; Vben Tech</a>
          <a href="#open-source">Open Source</a>
          <a :href="githubLink" target="_blank" rel="noreferrer">Docs</a>
        </div>

        <div class="footer-column">
          <h4>Legacy</h4>
          <a :href="icpLink" target="_blank" rel="noreferrer">{{ icp }}</a>
          <a :href="companySiteLink" target="_blank" rel="noreferrer">Official Site</a>
          <a :href="githubLink" target="_blank" rel="noreferrer">GitHub</a>
          <a :href="`mailto:${contactEmail}`">Contact</a>
        </div>

        <div class="footer-column">
          <h4>Status</h4>
          <div class="status-card">
            <div class="status-row">
              <span>Mainnet Status</span>
              <span class="status-indicator" />
            </div>
            <strong>All Systems Nominal</strong>
          </div>
        </div>
      </div>

      <div class="footer-bottom">
        <span>© {{ currentYear }} LiteChat. Built with Eino.</span>
        <div class="footer-icons">
          <a :href="companySiteLink" target="_blank" rel="noreferrer" aria-label="official site">
            <Globe :size="18" />
          </a>
          <a :href="`mailto:${contactEmail}`" aria-label="email">
            <Mail :size="18" />
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Space+Grotesk:wght@400;500;700&display=swap');

.landing-page {
  --bg: #131316;
  --bg-soft: rgba(42, 41, 45, 0.76);
  --bg-panel: rgba(32, 31, 35, 0.82);
  --line: rgba(66, 70, 86, 0.55);
  --line-soft: rgba(255, 255, 255, 0.08);
  --text: #e5e1e6;
  --text-soft: #c2c6d8;
  --text-muted: #8c90a1;
  --primary: #00daf8;
  --primary-strong: #007d8f;
  --secondary: #d1bcff;
  --tertiary: #ecb2ff;
  --shadow: 0 0 24px rgba(0, 218, 248, 0.12);
  position: relative;
  min-height: 100vh;
  color: var(--text);
  background: var(--bg);
  font-family: 'Inter', sans-serif;
  overflow-x: hidden;
}

.theme-light {
  --bg: #f4f7fb;
  --bg-soft: rgba(255, 255, 255, 0.84);
  --bg-panel: rgba(255, 255, 255, 0.9);
  --line: rgba(148, 163, 184, 0.28);
  --line-soft: rgba(15, 23, 42, 0.08);
  --text: #0f172a;
  --text-soft: #475569;
  --text-muted: #64748b;
  --primary: #0ea5e9;
  --primary-strong: #0284c7;
  --secondary: #7c3aed;
  --tertiary: #9333ea;
  --shadow: 0 16px 40px rgba(14, 165, 233, 0.12);
}

.background-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 36px 36px;
  opacity: 0.4;
}

.theme-light .grid-overlay {
  background-image:
    linear-gradient(rgba(15, 23, 42, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(15, 23, 42, 0.05) 1px, transparent 1px);
}

.glow {
  position: absolute;
  border-radius: 999px;
  filter: blur(120px);
}

.glow-left {
  top: 10%;
  left: -140px;
  width: 460px;
  height: 460px;
  background: rgba(0, 218, 248, 0.12);
}

.glow-right {
  right: -160px;
  bottom: 14%;
  width: 460px;
  height: 460px;
  background: rgba(183, 0, 247, 0.12);
}

.top-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  height: 80px;
  padding: 0 24px;
  border-bottom: 1px solid var(--line-soft);
  background: rgba(2, 6, 23, 0.82);
  backdrop-filter: blur(28px);
  box-shadow: 0 4px 30px rgba(0, 218, 248, 0.03);
}

.theme-light .top-nav {
  background: rgba(255, 255, 255, 0.84);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  object-fit: cover;
}

.brand-text,
.footer-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.04em;
  color: var(--primary);
}

.nav-links,
.nav-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.nav-links a,
.footer-column a {
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.2s ease, transform 0.2s ease;
}

.nav-links a:hover,
.footer-column a:hover {
  color: var(--text);
  transform: translateX(2px);
}

.theme-toggle,
.primary-nav-button,
.hero-primary,
.hero-secondary {
  border: 0;
  cursor: pointer;
  text-decoration: none;
}

.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  color: var(--text);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--line-soft);
}

.primary-nav-button,
.hero-primary {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  border-radius: 14px;
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 700;
  color: #04232c;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  box-shadow: 0 0 20px rgba(0, 218, 248, 0.24);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.primary-nav-button:hover,
.hero-primary:hover,
.hero-secondary:hover {
  transform: translateY(-1px);
}

.hero-secondary {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  border-radius: 14px;
  color: var(--text);
  background: rgba(112, 0, 255, 0.05);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(0, 218, 248, 0.1);
}

main {
  padding: 96px 24px 0;
}

.hero-section,
.orchestration-section,
.mission-section {
  max-width: 1280px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 64px;
  align-items: center;
  min-height: calc(100vh - 80px);
}

.hero-copy h1,
.section-heading h2,
.orchestration-copy h2,
.mission-copy h2 {
  margin: 0;
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 700;
  letter-spacing: -0.05em;
  line-height: 0.96;
}

.hero-copy h1 {
  font-size: clamp(3.5rem, 9vw, 6.8rem);
}

.hero-copy h1 span,
.mission-copy h2 span {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-copy p,
.section-heading p,
.orchestration-copy p,
.mission-copy p,
.feature-card p,
.stack-item p,
.future-content p,
.footer-brand p {
  color: var(--text-soft);
  line-height: 1.8;
}

.status-pill,
.footer-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  border-radius: 999px;
  background: var(--bg-soft);
  border: 1px solid var(--line);
  font-family: 'Space Grotesk', sans-serif;
  font-size: 0.78rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--tertiary);
  animation: pulse-ai 4s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.hero-copy,
.hero-visual,
.orchestration-copy,
.mission-copy {
  position: relative;
  z-index: 1;
}

.hero-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 28px;
}

.hero-meta {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-top: 36px;
  color: var(--text-muted);
}

.hero-meta span {
  display: block;
  font-size: 0.74rem;
  letter-spacing: 0.16em;
}

.hero-meta strong {
  display: block;
  margin-top: 6px;
  color: var(--text);
  font-size: 0.92rem;
}

.divider {
  width: 1px;
  height: 42px;
  background: var(--line);
}

.hero-visual {
  position: relative;
}

.glass-panel,
.feature-card,
.future-card,
.admin-mockup,
.status-card {
  background: rgba(112, 0, 255, 0.05);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(0, 218, 248, 0.1);
}

.theme-light .glass-panel,
.theme-light .feature-card,
.theme-light .future-card,
.theme-light .admin-mockup,
.theme-light .status-card {
  background: rgba(255, 255, 255, 0.76);
  border: 1px solid rgba(14, 165, 233, 0.16);
}

.neon-glow {
  box-shadow: var(--shadow);
}

.hero-orb-panel {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 999px;
  padding: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.orb-shell {
  position: relative;
  width: 100%;
  height: 100%;
}

.orb-ring {
  position: absolute;
  inset: 0;
  border-radius: 999px;
  border: 1px solid rgba(0, 218, 248, 0.24);
}

.orb-ring.inner {
  inset: 12%;
  border-style: dashed;
  border-color: rgba(209, 188, 255, 0.24);
}

.orb-core {
  position: absolute;
  inset: 28%;
  border-radius: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(32, 31, 35, 0.9);
  box-shadow: 0 0 24px rgba(0, 218, 248, 0.18);
}

.theme-light .orb-core {
  background: rgba(255, 255, 255, 0.92);
}

.orb-core img {
  width: 72%;
  height: 72%;
  object-fit: contain;
  border-radius: 999px;
  opacity: 0.86;
}

.floating-card {
  position: absolute;
  width: 132px;
  height: 132px;
  border-radius: 18px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  font-family: 'Space Grotesk', sans-serif;
  color: var(--primary);
}

.top-card {
  top: -18px;
  right: -18px;
  transform: rotate(10deg);
}

.bottom-card {
  bottom: 34px;
  left: -24px;
}

.section-block {
  max-width: 1280px;
  margin: 0 auto;
  padding: 120px 0;
}

.section-surface {
  background: rgba(255, 255, 255, 0.02);
}

.theme-light .section-surface {
  background: rgba(255, 255, 255, 0.4);
}

.section-heading.center {
  max-width: 800px;
  margin: 0 auto 72px;
  text-align: center;
}

.section-heading h2 {
  font-size: clamp(2.6rem, 6vw, 4rem);
}

.feature-grid,
.footer-grid {
  display: grid;
  gap: 28px;
}

.feature-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.feature-card {
  position: relative;
  overflow: hidden;
  padding: 32px;
  border-radius: 18px;
  transition: transform 0.3s ease, border-color 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  border-color: rgba(0, 218, 248, 0.24);
}

.feature-icon {
  color: var(--primary);
  margin-bottom: 24px;
}

.feature-card h3,
.stack-item h4,
.future-content h3,
.footer-column h4,
.status-card strong {
  font-family: 'Space Grotesk', sans-serif;
  margin: 0;
}

.feature-card h3 {
  font-size: 1.25rem;
  margin-bottom: 12px;
}

.stack-list {
  display: grid;
  gap: 18px;
  margin-top: 28px;
}

.stack-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  border-left: 4px solid var(--primary);
  border-radius: 14px;
  background: var(--bg-soft);
}

.stack-item.secondary {
  border-left-color: var(--secondary);
}

.stack-item svg {
  color: var(--primary);
  margin-top: 2px;
}

.orchestration-visual {
  border-radius: 22px;
  min-height: 520px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 36px;
}

.orchestration-core {
  position: relative;
  width: min(100%, 520px);
  aspect-ratio: 16 / 10;
  display: flex;
  align-items: center;
  justify-content: center;
}

.core-circle {
  border-radius: 999px;
  border: 1px solid rgba(0, 218, 248, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.core-circle.outer {
  width: 100%;
  aspect-ratio: 1 / 1;
}

.core-circle.inner {
  width: 76%;
  aspect-ratio: 1 / 1;
  border-style: dashed;
  border-color: rgba(209, 188, 255, 0.2);
}

.core-node {
  width: 42%;
  aspect-ratio: 1 / 1;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Space Grotesk', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary);
  background: var(--bg-soft);
  box-shadow: var(--shadow);
}

.orbit-node {
  position: absolute;
  width: 64px;
  height: 64px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-soft);
  border: 1px solid rgba(0, 218, 248, 0.18);
  color: var(--primary);
  font-family: 'Space Grotesk', sans-serif;
  font-size: 0.88rem;
  font-weight: 600;
}

.orbit-top {
  top: 10px;
  left: calc(50% - 32px);
}

.orbit-left {
  left: 12px;
  bottom: 72px;
}

.orbit-right {
  right: 12px;
  bottom: 72px;
}

.admin-mockup {
  border-radius: 22px;
  overflow: hidden;
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.24);
}

.mockup-topbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 18px;
  background: rgba(255, 255, 255, 0.03);
}

.theme-light .mockup-topbar {
  background: rgba(15, 23, 42, 0.04);
}

.mockup-topbar > span {
  width: 12px;
  height: 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.24);
}

.mockup-title {
  margin-left: 10px;
  font-size: 0.7rem;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.mockup-body {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 32px;
  padding: 28px;
}

.mockup-sidebar {
  display: grid;
  gap: 16px;
  align-content: start;
}

.sidebar-line {
  height: 42px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.06);
}

.theme-light .sidebar-line {
  background: rgba(15, 23, 42, 0.08);
}

.sidebar-line.short {
  width: 82%;
}

.sidebar-line.medium {
  width: 68%;
}

.sidebar-line.active {
  background: rgba(209, 188, 255, 0.2);
  border: 1px solid rgba(209, 188, 255, 0.2);
}

.mockup-main {
  display: grid;
  gap: 28px;
}

.mockup-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.stat-card {
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  color: var(--primary);
}

.stat-card span {
  color: var(--text);
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 600;
}

.dashboard-preview {
  padding: 18px;
  border-radius: 18px;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.preview-panel {
  min-height: 120px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(0, 218, 248, 0.14), rgba(209, 188, 255, 0.08));
}

.preview-panel.large {
  grid-column: span 3;
  min-height: 260px;
}

.mission-copy h2 {
  font-size: clamp(2.8rem, 7vw, 5rem);
}

.mission-copy ul {
  margin: 24px 0 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 14px;
}

.mission-copy li {
  position: relative;
  padding-left: 20px;
  color: var(--text);
}

.mission-copy li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 10px;
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: var(--primary);
}

.roadmap-link {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-top: 28px;
  color: var(--primary);
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 700;
  text-decoration: none;
}

.future-card {
  border-radius: 32px;
  padding: 18px;
}

.future-label {
  text-align: right;
  font-size: 0.72rem;
  letter-spacing: 0.28em;
  color: var(--tertiary);
}

.future-content {
  min-height: 540px;
  border-radius: 28px;
  padding: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 22px;
  border: 1px solid var(--line-soft);
  background: rgba(14, 14, 17, 0.6);
}

.theme-light .future-content {
  background: rgba(255, 255, 255, 0.74);
}

.future-content svg {
  color: var(--tertiary);
}

.site-footer {
  margin-top: 40px;
  padding: 80px 24px 32px;
  background: linear-gradient(to top, rgba(2, 6, 23, 1), transparent);
}

.theme-light .site-footer {
  background: linear-gradient(to top, rgba(226, 232, 240, 0.8), transparent);
}

.footer-grid,
.footer-bottom {
  max-width: 1280px;
  margin: 0 auto;
}

.footer-grid {
  grid-template-columns: 1.2fr repeat(3, minmax(0, 1fr));
}

.footer-brand p {
  max-width: 320px;
  margin-top: 18px;
}

.footer-badge {
  margin-top: 18px;
  color: var(--primary);
  font-size: 0.66rem;
}

.footer-column {
  display: grid;
  gap: 12px;
  align-content: start;
}

.status-card {
  border-radius: 16px;
  padding: 18px;
}

.status-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  color: var(--text-muted);
  font-size: 0.82rem;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #22c55e;
  box-shadow: 0 0 8px rgba(34, 197, 94, 0.5);
}

.footer-bottom {
  margin-top: 56px;
  padding-top: 26px;
  border-top: 1px solid var(--line-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  color: var(--text-muted);
  font-size: 0.92rem;
}

.footer-icons {
  display: flex;
  align-items: center;
  gap: 18px;
}

.footer-icons a {
  color: var(--text-muted);
}

@keyframes pulse-ai {
  0%,
  100% {
    opacity: 0.1;
    transform: scale(1);
  }

  50% {
    opacity: 0.3;
    transform: scale(1.2);
  }
}

@media (max-width: 1100px) {
  .hero-section,
  .orchestration-section,
  .mission-section {
    grid-template-columns: 1fr;
    min-height: auto;
    gap: 48px;
  }

  .feature-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .mockup-body,
  .footer-grid {
    grid-template-columns: 1fr;
  }

  .mockup-stats {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 820px) {
  .top-nav {
    height: auto;
    padding: 16px;
    flex-wrap: wrap;
  }

  .nav-links {
    display: none;
  }

  main {
    padding: 110px 16px 0;
  }

  .section-block {
    padding: 84px 0;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .hero-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .divider {
    width: 100%;
    height: 1px;
  }

  .floating-card {
    position: static;
    margin-top: 18px;
    width: 100%;
    height: 88px;
    flex-direction: row;
  }

  .hero-visual {
    display: grid;
    gap: 16px;
  }

  .preview-grid {
    grid-template-columns: 1fr;
  }

  .preview-panel.large {
    grid-column: span 1;
  }

  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
