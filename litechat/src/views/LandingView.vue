<template>
  <div :class="[
    'min-h-screen overflow-hidden relative font-sans transition-colors duration-300',
    isDark ? 'bg-black text-white selection:bg-purple-500 selection:text-white' : 'bg-gray-50 text-gray-900 selection:bg-purple-200 selection:text-purple-900'
  ]" @mousemove="handleMouseMove">
    
    <!-- Canvas Background for Particles -->
    <canvas ref="canvasRef" class="absolute top-0 left-0 w-full h-full z-0 pointer-events-none"></canvas>

    <!-- Background Effects (Parallax) -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden z-0 pointer-events-none">
      <div :class="['absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full blur-[120px] animate-pulse-slow transition-transform duration-100 ease-out', isDark ? 'bg-purple-600/20' : 'bg-purple-300/40']" :style="{ transform: `translate(${mouseX * 0.02}px, ${mouseY * 0.02}px)` }"></div>
      <div :class="['absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] rounded-full blur-[120px] animate-pulse-slow delay-1000 transition-transform duration-100 ease-out', isDark ? 'bg-blue-600/20' : 'bg-blue-300/40']" :style="{ transform: `translate(${mouseX * -0.02}px, ${mouseY * -0.02}px)` }"></div>
    </div>

    <!-- Grid Pattern Overlay -->
    <div class="absolute inset-0 opacity-20 z-0" :style="{ backgroundImage: `url('data:image/svg+xml;base64,${isDark ? 'PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImdyaWQiIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PHBhdGggZD0iTTAgNDBWMGg0MCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSIjZmZmIiBzdHJva2Utb3BhY2l0eT0iMC4wMyIvPjwvcGF0dGVybj48L2RlZnM+PHJlY3Qgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIgZmlsbD0idXJsKCNncmlkKSIvPjwvc3ZnPg==' : 'PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImdyaWQiIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PHBhdGggZD0iTTAgNDBWMGg0MCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSIjMDAwIiBzdHJva2Utb3BhY2l0eT0iMC4wMyIvPjwvcGF0dGVybj48L2RlZnM+PHJlY3Qgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIgZmlsbD0idXJsKCNncmlkKSIvPjwvc3ZnPg=='}')` }"></div>

    <!-- Navbar -->
    <nav :class="['relative z-50 flex items-center justify-between px-6 py-6 max-w-7xl mx-auto backdrop-blur-sm rounded-b-2xl border-b transition-all duration-300', isDark ? 'bg-black/20 border-white/5 hover:bg-black/40' : 'bg-white/50 border-black/5 hover:bg-white/80']">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 bg-gradient-to-br from-purple-500 to-blue-600 rounded-lg flex items-center justify-center shadow-lg shadow-purple-500/20">
          <span class="font-bold text-lg text-white">A</span>
        </div>
        <span :class="['text-xl font-bold tracking-tight bg-clip-text text-transparent', isDark ? 'bg-gradient-to-r from-white to-gray-400' : 'bg-gradient-to-r from-gray-900 to-gray-600']">{{ t.brand }}</span>
      </div>
      <div class="flex items-center gap-4 sm:gap-6">
        <!-- Theme Toggle -->
        <button @click="toggleTheme" :class="['p-2 rounded-full transition-colors', isDark ? 'text-gray-400 hover:text-white hover:bg-white/10' : 'text-gray-600 hover:text-black hover:bg-black/5']">
          <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"/></svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
        </button>

        <!-- Language Toggle -->
        <button @click="toggleLang" :class="['text-sm font-medium px-3 py-1 rounded-full transition-colors', isDark ? 'text-gray-400 hover:text-white hover:bg-white/10' : 'text-gray-600 hover:text-black hover:bg-black/5']">
          {{ lang === 'en' ? '中文' : 'English' }}
        </button>
        
        <a href="http://81.69.160.8:5666/" target="_blank" :class="['text-sm transition-colors hidden sm:block', isDark ? 'text-gray-400 hover:text-white' : 'text-gray-600 hover:text-black']">{{ t.adminConsole }}</a>
        
        <router-link to="/chat" class="group relative px-6 py-2 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-full font-medium text-sm overflow-hidden transition-all hover:scale-105 active:scale-95 shadow-lg shadow-purple-500/20">
          <span class="relative z-10">{{ t.launchApp }}</span>
          <div class="absolute inset-0 bg-white/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        </router-link>
      </div>
    </nav>

    <!-- Hero Section -->
    <main class="relative z-10 max-w-7xl mx-auto px-6 pt-20 pb-32">
      <div class="flex flex-col items-center text-center mb-32">
        <div :class="['inline-flex items-center gap-2 px-3 py-1 rounded-full border text-xs mb-8 backdrop-blur-sm animate-fade-in-up hover:scale-105 transition-transform cursor-default', isDark ? 'bg-white/5 border-white/10 text-purple-300' : 'bg-black/5 border-black/10 text-purple-700']">
          <span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
          {{ t.systemStatus }}
        </div>
        
        <h1 class="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight mb-8 animate-fade-in-up delay-100">
          <span :class="['block text-transparent bg-clip-text bg-gradient-to-b transition-colors duration-500', isDark ? 'from-white via-white to-gray-500 hover:text-white' : 'from-gray-900 via-gray-900 to-gray-500 hover:text-black']">{{ t.heroTitle1 }}</span>
          <span class="block text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-blue-400 to-purple-400 animate-gradient-x pb-2">{{ t.heroTitle2 }}</span>
        </h1>
        
        <p :class="['max-w-2xl text-lg mb-12 leading-relaxed animate-fade-in-up delay-200', isDark ? 'text-gray-400' : 'text-gray-600']">
          {{ t.heroDesc }}
        </p>

        <div class="flex flex-col sm:flex-row items-center gap-4 animate-fade-in-up delay-300">
          <router-link to="/chat" :class="['px-8 py-4 rounded-full font-bold text-lg transition-all hover:scale-105 shadow-lg', isDark ? 'bg-white text-black hover:bg-gray-100 hover:shadow-[0_0_40px_-10px_rgba(255,255,255,0.3)]' : 'bg-black text-white hover:bg-gray-800 hover:shadow-[0_0_40px_-10px_rgba(0,0,0,0.3)]']">
            {{ t.getStarted }}
          </router-link>
          <a href="#features" :class="['px-8 py-4 rounded-full border transition-all backdrop-blur-sm', isDark ? 'border-white/10 text-gray-300 hover:bg-white/5 hover:text-white' : 'border-black/10 text-gray-600 hover:bg-black/5 hover:text-black']">
            {{ t.exploreFeatures }}
          </a>
        </div>
      </div>

      <!-- Features List (Vertical) -->
      <div id="features" class="flex flex-col gap-32">
        
        <!-- Admin & DevOps -->
        <div class="flex flex-col md:flex-row items-center gap-12 group scroll-reveal" ref="feature1">
          <div class="w-full md:w-1/2 order-2 md:order-1">
            <div :class="['relative aspect-video rounded-3xl overflow-hidden border transition-all duration-500 group-hover:shadow-[0_0_30px_-5px_rgba(168,85,247,0.2)]', isDark ? 'bg-white/5 border-white/10 hover:border-purple-500/30' : 'bg-black/5 border-black/10 hover:border-purple-500/30']">
               <!-- Carousel -->
               <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide1 * 100}%)` }">
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-900/20 to-transparent">
                    <span class="text-purple-500/30 text-6xl font-bold">ADMIN UI 1</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-800/20 to-transparent">
                    <span class="text-purple-500/30 text-6xl font-bold">LOGS VIEW</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-purple-700/20 to-transparent">
                    <span class="text-purple-500/30 text-6xl font-bold">SETTINGS</span>
                 </div>
               </div>
               <!-- Carousel Controls -->
               <div class="absolute bottom-4 left-0 w-full flex justify-center gap-2">
                 <button v-for="i in 3" :key="i" @click="activeSlide1 = i - 1" :class="['w-2 h-2 rounded-full transition-all', activeSlide1 === i - 1 ? (isDark ? 'bg-white w-4' : 'bg-black w-4') : (isDark ? 'bg-white/30 hover:bg-white/50' : 'bg-black/30 hover:bg-black/50')]"></button>
               </div>
            </div>
          </div>
          <div class="w-full md:w-1/2 order-1 md:order-2">
            <div class="w-12 h-12 rounded-2xl bg-purple-500/20 flex items-center justify-center mb-6 text-purple-500 group-hover:scale-110 transition-transform duration-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.1a2 2 0 0 1-1-1.72v-.51a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>
            </div>
            <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.adminTitle }}</h3>
            <p :class="['text-lg leading-relaxed mb-6', isDark ? 'text-gray-400' : 'text-gray-600']">
              {{ t.adminDesc }}
            </p>
            <ul :class="['space-y-3', isDark ? 'text-gray-500' : 'text-gray-600']">
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-purple-500 rounded-full"></span>{{ t.adminFeat1 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-purple-500 rounded-full"></span>{{ t.adminFeat2 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-purple-500 rounded-full"></span>{{ t.adminFeat3 }}</li>
            </ul>
          </div>
        </div>

        <!-- Business Intelligence -->
        <div class="flex flex-col md:flex-row items-center gap-12 group scroll-reveal" ref="feature2">
          <div class="w-full md:w-1/2 order-1">
            <div class="w-12 h-12 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6 text-blue-500 group-hover:scale-110 transition-transform duration-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
            </div>
            <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.businessTitle }}</h3>
            <p :class="['text-lg leading-relaxed mb-6', isDark ? 'text-gray-400' : 'text-gray-600']">
              {{ t.businessDesc }}
            </p>
            <ul :class="['space-y-3', isDark ? 'text-gray-500' : 'text-gray-600']">
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>{{ t.businessFeat1 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>{{ t.businessFeat2 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>{{ t.businessFeat3 }}</li>
            </ul>
          </div>
          <div class="w-full md:w-1/2 order-2">
            <div :class="['relative aspect-video rounded-3xl overflow-hidden border transition-all duration-500 group-hover:shadow-[0_0_30px_-5px_rgba(59,130,246,0.2)]', isDark ? 'bg-white/5 border-white/10 hover:border-blue-500/30' : 'bg-black/5 border-black/10 hover:border-blue-500/30']">
               <!-- Carousel -->
               <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide2 * 100}%)` }">
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-bl from-blue-900/20 to-transparent">
                    <span class="text-blue-500/30 text-6xl font-bold">DATA UI 1</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-bl from-blue-800/20 to-transparent">
                    <span class="text-blue-500/30 text-6xl font-bold">CHARTS</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-bl from-blue-700/20 to-transparent">
                    <span class="text-blue-500/30 text-6xl font-bold">REPORTS</span>
                 </div>
               </div>
               <!-- Carousel Controls -->
               <div class="absolute bottom-4 left-0 w-full flex justify-center gap-2">
                 <button v-for="i in 3" :key="i" @click="activeSlide2 = i - 1" :class="['w-2 h-2 rounded-full transition-all', activeSlide2 === i - 1 ? (isDark ? 'bg-white w-4' : 'bg-black w-4') : (isDark ? 'bg-white/30 hover:bg-white/50' : 'bg-black/30 hover:bg-black/50')]"></button>
               </div>
            </div>
          </div>
        </div>

        <!-- Client Services -->
        <div class="flex flex-col md:flex-row items-center gap-12 group scroll-reveal" ref="feature3">
          <div class="w-full md:w-1/2 order-2 md:order-1">
            <div :class="['relative aspect-video rounded-3xl overflow-hidden border transition-all duration-500 group-hover:shadow-[0_0_30px_-5px_rgba(34,197,94,0.2)]', isDark ? 'bg-white/5 border-white/10 hover:border-green-500/30' : 'bg-black/5 border-black/10 hover:border-green-500/30']">
               <!-- Carousel -->
               <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide3 * 100}%)` }">
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-green-900/20 to-transparent">
                    <span class="text-green-500/30 text-6xl font-bold">CLIENT UI 1</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-green-800/20 to-transparent">
                    <span class="text-green-500/30 text-6xl font-bold">CHATBOT</span>
                 </div>
                 <div class="min-w-full h-full flex items-center justify-center bg-gradient-to-br from-green-700/20 to-transparent">
                    <span class="text-green-500/30 text-6xl font-bold">SUPPORT</span>
                 </div>
               </div>
               <!-- Carousel Controls -->
               <div class="absolute bottom-4 left-0 w-full flex justify-center gap-2">
                 <button v-for="i in 3" :key="i" @click="activeSlide3 = i - 1" :class="['w-2 h-2 rounded-full transition-all', activeSlide3 === i - 1 ? (isDark ? 'bg-white w-4' : 'bg-black w-4') : (isDark ? 'bg-white/30 hover:bg-white/50' : 'bg-black/30 hover:bg-black/50')]"></button>
               </div>
            </div>
          </div>
          <div class="w-full md:w-1/2 order-1 md:order-2">
            <div class="w-12 h-12 rounded-2xl bg-green-500/20 flex items-center justify-center mb-6 text-green-500 group-hover:scale-110 transition-transform duration-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
            </div>
            <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.clientTitle }}</h3>
            <p :class="['text-lg leading-relaxed mb-6', isDark ? 'text-gray-400' : 'text-gray-600']">
              {{ t.clientDesc }}
            </p>
            <ul :class="['space-y-3', isDark ? 'text-gray-500' : 'text-gray-600']">
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-green-500 rounded-full"></span>{{ t.clientFeat1 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-green-500 rounded-full"></span>{{ t.clientFeat2 }}</li>
              <li class="flex items-center gap-3"><span class="w-1.5 h-1.5 bg-green-500 rounded-full"></span>{{ t.clientFeat3 }}</li>
            </ul>
          </div>
        </div>

      </div>
    </main>

    <!-- Footer -->
    <footer :class="['relative z-10 border-t py-12 mt-32 backdrop-blur-md', isDark ? 'bg-black/50 border-white/5' : 'bg-white/50 border-black/5']">
      <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-6">
        <div :class="['text-sm', isDark ? 'text-gray-500' : 'text-gray-600']">
          {{ t.copyright }}
        </div>
        <div class="flex gap-6">
          <a href="#" :class="['text-sm transition-colors', isDark ? 'text-gray-500 hover:text-white' : 'text-gray-600 hover:text-black']">{{ t.privacy }}</a>
          <a href="#" :class="['text-sm transition-colors', isDark ? 'text-gray-500 hover:text-white' : 'text-gray-600 hover:text-black']">{{ t.terms }}</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useIntersectionObserver } from '@vueuse/core';

const lang = ref<'en' | 'zh'>('zh');
const isDark = ref(false); // Default to light mode
const canvasRef = ref<HTMLCanvasElement | null>(null);

// Mouse Parallax
const mouseX = ref(0);
const mouseY = ref(0);
const handleMouseMove = (e: MouseEvent) => {
  mouseX.value = e.clientX - window.innerWidth / 2;
  mouseY.value = e.clientY - window.innerHeight / 2;
};

// Carousel State
const activeSlide1 = ref(0);
const activeSlide2 = ref(0);
const activeSlide3 = ref(0);

// Auto-play Carousel
let carouselInterval: any;
const startCarousel = () => {
  carouselInterval = setInterval(() => {
    activeSlide1.value = (activeSlide1.value + 1) % 3;
    activeSlide2.value = (activeSlide2.value + 1) % 3;
    activeSlide3.value = (activeSlide3.value + 1) % 3;
  }, 5000);
};

// Scroll Reveal
const feature1 = ref(null);
const feature2 = ref(null);
const feature3 = ref(null);

const setupScrollReveal = (target: any) => {
  useIntersectionObserver(target, ([entry]) => {
    if (entry.isIntersecting) {
      entry.target.classList.add('revealed');
    }
  }, { threshold: 0.2 });
};

const translations = {
  en: {
    brand: 'AI Assistant',
    adminConsole: 'Admin Console',
    launchApp: 'Launch App',
    systemStatus: 'System Operational v2.0',
    heroTitle1: 'Intelligent Agent',
    heroTitle2: 'Orchestration',
    heroDesc: 'Empower your workflow with advanced AI agents. From automated reporting to intelligent customer service, experience the future of business automation.',
    getStarted: 'Get Started',
    exploreFeatures: 'Explore Features',
    adminTitle: 'Admin & DevOps',
    adminDesc: 'Comprehensive monitoring and code auditing. Ensure your system runs smoothly with AI-powered insights.',
    adminFeat1: 'Log Monitoring & Analysis',
    adminFeat2: 'Code Audit & Security',
    adminFeat3: 'Knowledge Base Generation',
    businessTitle: 'Business Intelligence',
    businessDesc: 'Data-driven insights and automated reporting. Transform raw data into actionable strategies effortlessly.',
    businessFeat1: 'Automated Excel Reports',
    businessFeat2: 'Data Analysis & Stats',
    businessFeat3: 'Business Flow Management',
    clientTitle: 'Client Services',
    clientDesc: 'Seamless onboarding and smart support. Guide your users with intelligent, rule-based assistance.',
    clientFeat1: 'Smart Customer Support',
    clientFeat2: 'System Onboarding',
    clientFeat3: 'Rule-based Guidance',
    copyright: '© 2024 AI Assistant. All rights reserved.',
    privacy: 'Privacy Policy',
    terms: 'Terms of Service',
  },
  zh: {
    brand: 'AI 智能助手',
    adminConsole: '管理控制台',
    launchApp: '启动应用',
    systemStatus: '系统运行中 v2.0',
    heroTitle1: '智能体',
    heroTitle2: '编排与服务',
    heroDesc: '利用先进的 AI 智能体赋能您的工作流程。从自动化报表到智能客户服务，体验业务自动化的未来。',
    getStarted: '立即开始',
    exploreFeatures: '探索功能',
    adminTitle: '管理与运维',
    adminDesc: '全面的监控与代码审计功能。通过 AI 驱动的洞察力确保您的系统平稳运行。',
    adminFeat1: '日志监控与分析',
    adminFeat2: '代码审计与安全',
    adminFeat3: '知识库生成',
    businessTitle: '商业智能',
    businessDesc: '数据驱动的洞察与自动化报表。轻松将原始数据转化为可执行的策略。',
    businessFeat1: '自动化 Excel 报表',
    businessFeat2: '数据分析与统计',
    businessFeat3: '业务流程管理',
    clientTitle: '客户服务',
    clientDesc: '无缝入驻与智能支持。通过智能、基于规则的协助引导您的用户。',
    clientFeat1: '智能客户支持',
    clientFeat2: '系统入驻引导',
    clientFeat3: '规则导向指引',
    copyright: '© 2024 AI 智能助手. 保留所有权利.',
    privacy: '隐私政策',
    terms: '服务条款',
  }
};

const t = computed(() => translations[lang.value]);

const toggleLang = () => {
  lang.value = lang.value === 'en' ? 'zh' : 'en';
};

const toggleTheme = () => {
  isDark.value = !isDark.value;
  // Re-init particles to change color
  initParticles();
};

// Particle System
let animationId: number;
interface Particle {
  x: number;
  y: number;
  size: number;
  speedX: number;
  speedY: number;
  opacity: number;
}

const initParticles = () => {
  const canvas = canvasRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;

  const particles: Particle[] = [];
  const particleCount = 50; // Number of particles
  const particleColor = '34, 197, 94'; // Green-500

  for (let i = 0; i < particleCount; i++) {
    particles.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      size: Math.random() * 4 + 2, // Slightly larger squares
      speedX: (Math.random() - 0.5) * 0.5,
      speedY: (Math.random() - 0.5) * 0.5,
      opacity: Math.random() * 0.5 + 0.1
    });
  }

  const animate = () => {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    particles.forEach(p => {
      p.x += p.speedX;
      p.y += p.speedY;

      // Wrap around screen
      if (p.x < 0) p.x = canvas.width;
      if (p.x > canvas.width) p.x = 0;
      if (p.y < 0) p.y = canvas.height;
      if (p.y > canvas.height) p.y = 0;

      ctx.fillStyle = `rgba(${particleColor}, ${p.opacity})`;
      ctx.fillRect(p.x, p.y, p.size, p.size);
    });

    animationId = requestAnimationFrame(animate);
  };

  if (animationId) cancelAnimationFrame(animationId);
  animate();
};

const handleResize = () => {
  if (canvasRef.value) {
    canvasRef.value.width = window.innerWidth;
    canvasRef.value.height = window.innerHeight;
    initParticles(); // Re-init to handle resize properly
  }
};

onMounted(() => {
  initParticles();
  window.addEventListener('resize', handleResize);
  startCarousel();
  
  setupScrollReveal(feature1);
  setupScrollReveal(feature2);
  setupScrollReveal(feature3);
});

onUnmounted(() => {
  cancelAnimationFrame(animationId);
  window.removeEventListener('resize', handleResize);
  clearInterval(carouselInterval);
});
</script>

<style scoped>
@keyframes pulse-slow {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes gradient-x {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.animate-pulse-slow {
  animation: pulse-slow 8s infinite ease-in-out;
}

.animate-float {
  animation: float 6s infinite ease-in-out;
}

.animate-fade-in-up {
  animation: fade-in-up 0.8s ease-out forwards;
  opacity: 0; /* Start hidden */
}

.animate-gradient-x {
  background-size: 200% 200%;
  animation: gradient-x 3s ease infinite;
}

.delay-100 { animation-delay: 100ms; }
.delay-200 { animation-delay: 200ms; }
.delay-300 { animation-delay: 300ms; }
.delay-500 { animation-delay: 500ms; }
.delay-1000 { animation-delay: 1000ms; }

/* Scroll Reveal Styles */
.scroll-reveal {
  opacity: 0;
  transform: translateY(50px);
  transition: all 1s ease-out;
}

.scroll-reveal.revealed {
  opacity: 1;
  transform: translateY(0);
}
</style>
