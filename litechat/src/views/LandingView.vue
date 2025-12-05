<template>
  <div :class="[
    'min-h-screen overflow-x-hidden relative font-sans transition-colors duration-500',
    isDark ? 'bg-[#0a0a0a] text-white selection:bg-cyan-500/30 selection:text-cyan-100' : 'bg-gray-50 text-gray-900 selection:bg-cyan-200 selection:text-cyan-900'
  ]">
    
    <!-- Tech Background (Grid + Snakes) -->
    <div class="fixed inset-0 z-0 overflow-hidden pointer-events-none">
        <!-- Grid Overlay -->
        <div class="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]"></div>
        <div class="absolute inset-0 bg-[radial-gradient(circle_800px_at_50%_200px,#3b82f615,transparent)]"></div>
    </div>

    <!-- Canvas Background for Multiple Snakes -->
    <canvas ref="canvasRef" class="fixed top-0 left-0 w-full h-full z-0 pointer-events-none opacity-60"></canvas>

    <!-- Navbar -->
    <nav :class="['fixed top-0 left-0 right-0 z-50 transition-all duration-500', isScrolled ? 'bg-[#0a0a0a]/80 backdrop-blur-md border-b border-white/10 shadow-lg shadow-cyan-500/5 py-4' : 'py-6 bg-transparent']">
      <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
        <div class="flex items-center gap-3 cursor-pointer" @click="$router.push('/')">
          <div class="w-10 h-10 bg-gradient-to-br from-cyan-500/20 to-blue-600/20 backdrop-blur-xl rounded-2xl flex items-center justify-center border border-cyan-500/30 shadow-[0_0_15px_rgba(6,182,212,0.15)]">
            <span class="font-bold text-xl text-cyan-400">A</span>
          </div>
          <span :class="['text-xl font-bold tracking-tight', isDark ? 'text-white' : 'text-gray-900']">{{ t.brand }}</span>
        </div>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center gap-4">
          <button @click="toggleTheme" :class="['p-2.5 rounded-full transition-all hover:scale-110 active:scale-90 border border-white/5', isDark ? 'bg-white/5 text-gray-400 hover:text-white hover:bg-white/10' : 'bg-black/5 text-gray-600 hover:text-black hover:bg-black/10']">
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"/></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
          </button>

          <button @click="toggleLang" :class="['text-sm font-medium px-5 py-2.5 rounded-full transition-all hover:scale-110 active:scale-90 border border-white/5', isDark ? 'bg-white/5 text-gray-300 hover:text-white hover:bg-white/10' : 'bg-black/5 text-gray-700 hover:text-black hover:bg-black/10']">
            {{ lang === 'en' ? '中文' : 'English' }}
          </button>
          
          <a href="http://81.69.160.8:5666/" target="_blank" :class="['text-sm font-medium px-5 py-2.5 rounded-full transition-all hover:scale-110 active:scale-90 border border-white/5', isDark ? 'bg-white/5 text-gray-300 hover:text-white hover:bg-white/10' : 'bg-black/5 text-gray-700 hover:text-black hover:bg-black/10']">{{ t.adminConsole }}</a>
          
          <router-link to="/chat" class="px-8 py-3 bg-gradient-to-r from-cyan-400 to-blue-500 hover:from-cyan-300 hover:to-blue-400 text-white rounded-full font-bold text-sm shadow-lg shadow-cyan-500/30 hover:shadow-cyan-500/50 hover:scale-110 active:scale-90 hover:-rotate-1 transition-all duration-300">
            {{ t.launchApp }}
          </router-link>
        </div>

        <!-- Mobile Menu Button -->
        <button @click="isMobileMenuOpen = !isMobileMenuOpen" :class="['md:hidden p-2 rounded-xl border border-white/10', isDark ? 'bg-white/5 text-white' : 'bg-black/5 text-black']">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
        </button>
      </div>

      <!-- Mobile Menu Dropdown -->
      <div v-show="isMobileMenuOpen" :class="['md:hidden absolute top-full left-0 right-0 border-b backdrop-blur-xl p-4 flex flex-col gap-3 shadow-2xl transition-all duration-300', isDark ? 'bg-[#0a0a0a]/95 border-white/10' : 'bg-white/95 border-black/10']">
          <a href="http://81.69.160.8:5666/" target="_blank" :class="['p-3 rounded-2xl text-sm font-medium', isDark ? 'bg-white/5 text-white' : 'bg-black/5 text-black']">{{ t.adminConsole }}</a>
          <div class="flex items-center gap-2">
             <button @click="toggleTheme" :class="['flex-1 p-3 rounded-2xl flex items-center justify-center', isDark ? 'bg-white/5 text-white' : 'bg-black/5 text-black']">
               <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"/></svg>
               <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
             </button>
             <button @click="toggleLang" :class="['flex-1 p-3 rounded-2xl text-sm font-medium', isDark ? 'bg-white/5 text-white' : 'bg-black/5 text-black']">{{ lang === 'en' ? '中文' : 'English' }}</button>
          </div>
          <router-link to="/chat" class="w-full py-4 bg-gradient-to-r from-cyan-400 to-blue-500 text-white rounded-2xl font-bold text-center shadow-lg shadow-cyan-500/30">
            {{ t.launchApp }}
          </router-link>
      </div>
    </nav>

    <!-- Hero Section -->
    <main class="relative z-10 max-w-7xl mx-auto px-6 pt-40 pb-20">
      <div class="flex flex-col items-center text-center mb-32 relative">
        
        <div :class="['inline-flex items-center gap-2 px-4 py-1.5 rounded-full border backdrop-blur-md text-xs font-mono mb-8 shadow-lg transition-transform hover:scale-105 cursor-default', isDark ? 'bg-cyan-500/10 border-cyan-500/20 text-cyan-400' : 'bg-cyan-500/10 border-cyan-500/20 text-cyan-700']">
          <span class="relative flex h-2.5 w-2.5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-cyan-500"></span>
          </span>
          {{ t.systemStatus }}
        </div>
        
        <h1 class="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight mb-8 animate-fade-in-up delay-100 relative z-10 drop-shadow-2xl">
          <span :class="['block bg-clip-text text-transparent bg-gradient-to-b', isDark ? 'from-white to-white/50' : 'from-black to-black/50']">{{ t.heroTitle1 }}</span>
          <span class="block bg-clip-text text-transparent bg-gradient-to-r from-cyan-400 via-blue-500 to-purple-600 animate-gradient-x pb-4">{{ t.heroTitle2 }}</span>
        </h1>
        
        <p :class="['max-w-2xl text-lg md:text-xl mb-12 leading-relaxed animate-fade-in-up delay-200 font-medium', isDark ? 'text-gray-400' : 'text-gray-600']">
          {{ t.heroDesc }}
        </p>

        <div class="flex flex-col sm:flex-row items-center gap-6 animate-fade-in-up delay-300 relative z-10">
          <router-link to="/chat" style="padding-left: 2.5rem; padding-right: 2.5rem;" class="py-4 rounded-full font-bold text-lg transition-all duration-300 hover:scale-110 active:scale-95 hover:-rotate-1 shadow-xl shadow-cyan-500/30 flex items-center gap-2 bg-gradient-to-r from-cyan-400 to-blue-500 text-white hover:shadow-cyan-500/50">
            {{ t.getStarted }}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path></svg>
          </router-link>
          <a href="#features" style="padding-left: 2.5rem; padding-right: 2.5rem;" :class="['py-4 rounded-full font-semibold border transition-all duration-300 backdrop-blur-md hover:scale-110 active:scale-95 hover:rotate-1', isDark ? 'bg-white/5 border-white/10 text-white hover:bg-white/10' : 'bg-black/5 border-black/10 text-gray-900 hover:bg-black/10']">
            {{ t.exploreFeatures }}
          </a>
        </div>
      </div>

      <!-- Online Free Trial Section -->
      <div class="mb-40 animate-fade-in-up delay-500">
        <div class="text-center mb-16">
           <h2 :class="['text-3xl md:text-4xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.freeTrialTitle }}</h2>
           <p :class="['text-lg font-medium', isDark ? 'text-gray-400' : 'text-gray-500']">{{ t.freeTrialDesc }}</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div v-for="(mode, index) in playModes" :key="index" :class="['group relative rounded-[32px] p-[1px] transition-all duration-500 hover:-translate-y-2', isDark ? 'bg-gradient-to-b from-white/10 to-transparent hover:from-cyan-500/50' : 'bg-gradient-to-b from-black/10 to-transparent hover:from-cyan-500/50']">
            <div :class="['relative h-full rounded-[31px] p-8 overflow-hidden backdrop-blur-xl flex flex-col border border-white/5 shadow-2xl', isDark ? 'bg-[#0a0a0a]/80' : 'bg-white/80']">
              
              <!-- Background Image or Icon -->
              <div class="absolute top-0 right-0 p-6 opacity-10 group-hover:opacity-20 transition-opacity duration-500 transform group-hover:scale-110 grayscale group-hover:grayscale-0">
                <img v-if="mode.image" :src="mode.image" class="w-32 h-32 object-cover rounded-2xl" alt="mode bg" />
                <component v-else :is="mode.icon" class="w-32 h-32" />
              </div>
              
              <div :class="['w-16 h-16 rounded-2xl flex items-center justify-center mb-6 transition-all duration-500 shadow-lg border border-white/10', isDark ? 'bg-white/5 text-cyan-400 group-hover:bg-cyan-500 group-hover:text-black group-hover:shadow-[0_0_20px_rgba(6,182,212,0.4)]' : 'bg-black/5 text-cyan-600 group-hover:bg-cyan-500 group-hover:text-white']">
                <img v-if="mode.image" :src="mode.image" class="w-full h-full object-cover rounded-2xl" alt="mode icon" />
                <component v-else :is="mode.icon" class="w-8 h-8" />
              </div>

              <h3 :class="['text-2xl font-bold mb-3', isDark ? 'text-white' : 'text-gray-900']">{{ mode.title }}</h3>
              <p :class="['text-sm mb-8 flex-grow leading-relaxed font-medium', isDark ? 'text-gray-400' : 'text-gray-600']">{{ mode.desc }}</p>
              
              <button :class="['w-full py-4 rounded-full font-bold transition-all duration-300 flex items-center justify-center gap-2 backdrop-blur-md shadow-lg hover:scale-105 active:scale-95', isDark ? 'bg-gradient-to-r from-cyan-500/20 to-blue-500/20 hover:from-cyan-400 hover:to-blue-500 text-white hover:shadow-cyan-500/40 border border-white/10' : 'bg-gradient-to-r from-gray-100 to-gray-200 hover:from-cyan-400 hover:to-blue-500 text-gray-900 hover:text-white hover:shadow-cyan-500/30 border border-black/5']">
                {{ t.tryNow }}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Features Timeline Section -->
      <div id="features" class="relative">
        <!-- Vertical Timeline Line -->
        <div class="absolute left-1/2 transform -translate-x-1/2 h-full w-px bg-gradient-to-b from-transparent via-cyan-500/50 to-transparent hidden md:block"></div>
        
        <div class="flex flex-col gap-40">
          
          <!-- Feature 1: Admin -->
          <div class="relative flex flex-col md:flex-row items-center gap-16 group scroll-reveal" ref="feature1">
            <!-- Timeline Dot -->
            <div class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 w-4 h-4 rounded-full bg-cyan-500 shadow-[0_0_15px_rgba(6,182,212,0.8)] hidden md:block z-20 ring-4 ring-[#0a0a0a]"></div>

            <!-- Content Left -->
            <div class="w-full md:w-1/2 flex justify-end md:pr-20">
               <div :class="['w-full max-w-lg relative aspect-video rounded-3xl overflow-hidden border transition-all duration-700 group-hover:scale-105 shadow-2xl', isDark ? 'bg-white/5 border-white/10 shadow-cyan-500/10' : 'bg-black/5 border-black/10 shadow-cyan-500/10']">
                  <div class="absolute inset-0 backdrop-blur-sm"></div>
                  <!-- Carousel -->
                  <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide1 * 100}%)` }">
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-cyan-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">ADMIN</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-cyan-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">LOGS</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-cyan-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">CONFIG</span>
                    </div>
                  </div>
                  <!-- Controls -->
                  <div class="absolute bottom-6 left-0 w-full flex justify-center gap-3">
                    <button v-for="i in 3" :key="i" @click="activeSlide1 = i - 1" :class="['h-1.5 rounded-full transition-all duration-300', activeSlide1 === i - 1 ? 'w-8 bg-cyan-500' : 'w-2 bg-cyan-500/30 hover:bg-cyan-500/60']"></button>
                  </div>
               </div>
            </div>

            <!-- Text Right -->
            <div class="w-full md:w-1/2 md:pl-20 text-left">
              <div class="w-14 h-14 rounded-2xl bg-cyan-500/10 flex items-center justify-center mb-6 text-cyan-400 shadow-[0_0_15px_rgba(6,182,212,0.2)] border border-cyan-500/20">
                <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.1a2 2 0 0 1-1-1.72v-.51a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>
              </div>
              <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.adminTitle }}</h3>
              <p :class="['text-lg leading-relaxed mb-6 font-medium', isDark ? 'text-gray-400' : 'text-gray-600']">
                {{ t.adminDesc }}
              </p>
              <ul :class="['space-y-3', isDark ? 'text-gray-300' : 'text-gray-700']">
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-cyan-500 rounded-full shadow-[0_0_8px_rgba(6,182,212,0.8)]"></span>{{ t.adminFeat1 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-cyan-500 rounded-full shadow-[0_0_8px_rgba(6,182,212,0.8)]"></span>{{ t.adminFeat2 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-cyan-500 rounded-full shadow-[0_0_8px_rgba(6,182,212,0.8)]"></span>{{ t.adminFeat3 }}</li>
              </ul>
            </div>
          </div>

          <!-- Feature 2: Business (Reversed) -->
          <div class="relative flex flex-col md:flex-row-reverse items-center gap-16 group scroll-reveal" ref="feature2">
            <!-- Timeline Dot -->
            <div class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 w-4 h-4 rounded-full bg-blue-500 shadow-[0_0_15px_rgba(59,130,246,0.8)] hidden md:block z-20 ring-4 ring-[#0a0a0a]"></div>

            <!-- Content Right (Image) -->
            <div class="w-full md:w-1/2 flex justify-start md:pl-20">
               <div :class="['w-full max-w-lg relative aspect-video rounded-3xl overflow-hidden border transition-all duration-700 group-hover:scale-105 shadow-2xl', isDark ? 'bg-white/5 border-white/10 shadow-blue-500/10' : 'bg-black/5 border-black/10 shadow-blue-500/10']">
                  <div class="absolute inset-0 backdrop-blur-sm"></div>
                  <!-- Carousel -->
                  <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide2 * 100}%)` }">
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-blue-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">DATA</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-blue-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">CHARTS</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-blue-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">REPORTS</span>
                    </div>
                  </div>
                  <!-- Controls -->
                  <div class="absolute bottom-6 left-0 w-full flex justify-center gap-3">
                    <button v-for="i in 3" :key="i" @click="activeSlide2 = i - 1" :class="['h-1.5 rounded-full transition-all duration-300', activeSlide2 === i - 1 ? 'w-8 bg-blue-500' : 'w-2 bg-blue-500/30 hover:bg-blue-500/60']"></button>
                  </div>
               </div>
            </div>

            <!-- Text Left -->
            <div class="w-full md:w-1/2 md:pr-20 text-left md:text-right">
              <div class="w-14 h-14 rounded-2xl bg-blue-500/10 flex items-center justify-center mb-6 text-blue-500 shadow-[0_0_15px_rgba(59,130,246,0.2)] border border-blue-500/20 md:ml-auto">
                <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
              </div>
              <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.businessTitle }}</h3>
              <p :class="['text-lg leading-relaxed mb-6 font-medium', isDark ? 'text-gray-400' : 'text-gray-600']">
                {{ t.businessDesc }}
              </p>
              <ul :class="['space-y-3 flex flex-col md:items-end', isDark ? 'text-gray-300' : 'text-gray-700']">
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-blue-500 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.8)]"></span>{{ t.businessFeat1 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-blue-500 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.8)]"></span>{{ t.businessFeat2 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-blue-500 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.8)]"></span>{{ t.businessFeat3 }}</li>
              </ul>
            </div>
          </div>

          <!-- Feature 3: Client -->
          <div class="relative flex flex-col md:flex-row items-center gap-16 group scroll-reveal" ref="feature3">
            <!-- Timeline Dot -->
            <div class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 w-4 h-4 rounded-full bg-purple-500 shadow-[0_0_15px_rgba(168,85,247,0.8)] hidden md:block z-20 ring-4 ring-[#0a0a0a]"></div>

            <!-- Content Left -->
            <div class="w-full md:w-1/2 flex justify-end md:pr-20">
               <div :class="['w-full max-w-lg relative aspect-video rounded-3xl overflow-hidden border transition-all duration-700 group-hover:scale-105 shadow-2xl', isDark ? 'bg-white/5 border-white/10 shadow-purple-500/10' : 'bg-black/5 border-black/10 shadow-purple-500/10']">
                  <div class="absolute inset-0 backdrop-blur-sm"></div>
                  <!-- Carousel -->
                  <div class="absolute inset-0 flex transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${activeSlide3 * 100}%)` }">
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-purple-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">CLIENT</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-purple-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">CHATBOT</span>
                    </div>
                    <div class="min-w-full h-full flex items-center justify-center bg-black/20">
                       <span class="text-purple-500/20 text-5xl md:text-6xl font-bold tracking-widest font-mono">SUPPORT</span>
                    </div>
                  </div>
                  <!-- Controls -->
                  <div class="absolute bottom-6 left-0 w-full flex justify-center gap-3">
                    <button v-for="i in 3" :key="i" @click="activeSlide3 = i - 1" :class="['h-1.5 rounded-full transition-all duration-300', activeSlide3 === i - 1 ? 'w-8 bg-purple-500' : 'w-2 bg-purple-500/30 hover:bg-purple-500/60']"></button>
                  </div>
               </div>
            </div>

            <!-- Text Right -->
            <div class="w-full md:w-1/2 md:pl-20 text-left">
              <div class="w-14 h-14 rounded-2xl bg-purple-500/10 flex items-center justify-center mb-6 text-purple-500 shadow-[0_0_15px_rgba(168,85,247,0.2)] border border-purple-500/20">
                <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
              </div>
              <h3 :class="['text-3xl font-bold mb-4', isDark ? 'text-white' : 'text-gray-900']">{{ t.clientTitle }}</h3>
              <p :class="['text-lg leading-relaxed mb-6 font-medium', isDark ? 'text-gray-400' : 'text-gray-600']">
                {{ t.clientDesc }}
              </p>
              <ul :class="['space-y-3', isDark ? 'text-gray-300' : 'text-gray-700']">
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-purple-500 rounded-full shadow-[0_0_8px_rgba(168,85,247,0.8)]"></span>{{ t.clientFeat1 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-purple-500 rounded-full shadow-[0_0_8px_rgba(168,85,247,0.8)]"></span>{{ t.clientFeat2 }}</li>
                <li class="flex items-center gap-3"><span class="w-2 h-2 bg-purple-500 rounded-full shadow-[0_0_8px_rgba(168,85,247,0.8)]"></span>{{ t.clientFeat3 }}</li>
              </ul>
            </div>
          </div>

        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer :class="['relative z-10 border-t py-12 mt-40 backdrop-blur-xl', isDark ? 'bg-black/30 border-white/5' : 'bg-white/30 border-black/5']">
      <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-6">
        <div :class="['text-sm font-medium', isDark ? 'text-gray-500' : 'text-gray-500']">
          {{ t.copyright }}
        </div>
        <div class="flex gap-8">
          <a href="#" :class="['text-sm font-medium transition-colors', isDark ? 'text-gray-500 hover:text-white' : 'text-gray-500 hover:text-black']">{{ t.privacy }}</a>
          <a href="#" :class="['text-sm font-medium transition-colors', isDark ? 'text-gray-500 hover:text-white' : 'text-gray-500 hover:text-black']">{{ t.terms }}</a>
        </div>
      </div>
    </footer>
    
    <!-- Easter Egg Hint (Bottom Right) -->
    <div v-if="showSnakeHint" :class="['fixed bottom-8 right-8 z-50 px-6 py-3 rounded-full border backdrop-blur-xl text-sm font-medium shadow-2xl transition-all duration-500 flex items-center gap-3 max-w-xs', isDark ? 'bg-cyan-500/10 border-cyan-500/20 text-cyan-400' : 'bg-cyan-500/20 border-cyan-500/30 text-cyan-700']" style="animation: slideInRight 0.5s ease-out, bounce-slow 2s ease-in-out 0.5s infinite;">
      <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
      <span>{{ t.snakeHint }}</span>
      <button @click="showSnakeHint = false" class="ml-2 p-1 rounded-full hover:bg-white/10 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue';
import { useIntersectionObserver, useWindowScroll } from '@vueuse/core';

const lang = ref<'en' | 'zh'>('zh');
const isDark = ref(true);
const canvasRef = ref<HTMLCanvasElement | null>(null);
const isMobileMenuOpen = ref(false);
const showSnakeHint = ref(true);

// Scroll State for Header
const { y } = useWindowScroll();
const isScrolled = computed(() => y.value > 20);

// Mouse Position for Hero Snake
const mouseX = ref(0);
const mouseY = ref(0);
const isMouseMoving = ref(false);
const isMouseInWindow = ref(false);
let mouseTimeout: any;
let hintTimeout: any;

const handleMouseMove = (e: MouseEvent) => {
  mouseX.value = e.clientX;
  mouseY.value = e.clientY;
  isMouseMoving.value = true;
  isMouseInWindow.value = true;
  
  clearTimeout(mouseTimeout);
  mouseTimeout = setTimeout(() => {
    isMouseMoving.value = false;
  }, 2000);
  
  // Hide hint after user interacts
  if (showSnakeHint.value) {
    clearTimeout(hintTimeout);
    hintTimeout = setTimeout(() => {
      showSnakeHint.value = false;
    }, 5000);
  }
};

const handleMouseLeave = () => {
  isMouseInWindow.value = false;
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
  }, { threshold: 0.15 });
};

// Icons for Play Modes
const UserIcon = {
  render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': 2 }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' })
  ])
};
const MicIcon = {
  render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': 2 }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z' })
  ])
};
const ImageIcon = {
  render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': 2 }, [
    h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', d: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z' })
  ])
};

const translations = {
  en: {
    snakeHint: '🎮 Move your mouse to control the cyan snake!',
    brand: 'LiteChat AI',
    adminConsole: 'Admin Console',
    launchApp: 'Launch App',
    systemStatus: 'System Online v2.0',
    heroTitle1: 'Enterprise-Grade',
    heroTitle2: 'AI Platform',
    heroDesc: 'Flexible AI solution for individuals and enterprises. Support SaaS subscription and private deployment with workflow automation, RAG, and MCP integration.',
    getStarted: 'Try Free',
    exploreFeatures: 'Learn More',
    freeTrialTitle: 'Flexible Deployment',
    freeTrialDesc: 'Choose the mode that fits your needs - standalone, cloud-based, or demo.',
    tryNow: 'Try Now',
    adminTitle: 'Workflow & Automation',
    adminDesc: 'Visual workflow designer with drag-and-drop interface. Configure agents, prompts, and integrate MCP tools seamlessly.',
    adminFeat1: 'Visual Workflow Designer',
    adminFeat2: 'Agent & Prompt Management',
    adminFeat3: 'MCP Tool Integration',
    businessTitle: 'Enterprise Management',
    businessDesc: 'Comprehensive admin panel for enterprises. Customize client UI, monitor logs with AI analysis, and manage knowledge bases.',
    businessFeat1: 'Visual UI Customization',
    businessFeat2: 'AI-Powered Log Analysis',
    businessFeat3: 'RAG Knowledge Base',
    clientTitle: 'Smart Client',
    clientDesc: 'Rich features including voice-to-text, web search, long-term memory, meeting summaries, and multi-modal interactions.',
    clientFeat1: 'Voice & Meeting Summaries',
    clientFeat2: 'Long-term Memory (RAG)',
    clientFeat3: 'Multi-modal Input Support',
    copyright: '© 2024 LiteChat AI. All rights reserved.',
    privacy: 'Privacy Policy',
    terms: 'Terms of Service',
  },
  zh: {
    brand: 'LiteChat AI',
    adminConsole: '管理控制台',
    launchApp: '启动应用',
    systemStatus: '系统运行中 v2.0',
    snakeHint: '🎮 移动鼠标控制青色贪食蛇！',
    heroTitle1: '企业级',
    heroTitle2: 'AI 智能平台',
    heroDesc: '面向个人与企业的灵活 AI 解决方案。支持 SaaS 订阅与私有化部署，提供工作流自动化、RAG 知识库与 MCP 工具集成。',
    getStarted: '在线使用',
    exploreFeatures: '了解更多',
    freeTrialTitle: '灵活部署模式',
    freeTrialDesc: '选择适合您的模式 - 纯前端、云端后台或演示样例。',
    tryNow: '立即试用',
    adminTitle: '工作流与自动化',
    adminDesc: '可视化工作流设计器，拖拽式配置。轻松管理智能体、提示词模板，并集成 MCP 工具。',
    adminFeat1: '可视化工作流设计',
    adminFeat2: '智能体与提示词管理',
    adminFeat3: 'MCP 工具集成',
    businessTitle: '企业管理平台',
    businessDesc: '企业级管理后台。支持客户端 UI 可视化定制、AI 日志分析与知识库管理。',
    businessFeat1: '客户端 UI 可视化定制',
    businessFeat2: 'AI 驱动的日志分析',
    businessFeat3: 'RAG 知识库对接',
    clientTitle: '智能客户端',
    clientDesc: '丰富的功能特性：语音转文字、联网搜索、长期记忆、会议总结、多模态交互。',
    clientFeat1: '语音与会议总结',
    clientFeat2: '长期记忆（RAG）',
    clientFeat3: '多模态输入支持',
    copyright: '© 2024 LiteChat AI. 保留所有权利.',
    privacy: '隐私政策',
    terms: '服务条款',
  }
};

const t = computed(() => translations[lang.value]);

const playModes = computed(() => {
    const isZh = lang.value === 'zh';
    return [
        {
            title: isZh ? '纯前端模式' : 'Backend Mode',
            desc: isZh ? '无需登录，仅需 API Key 即可使用，数据本地存储。适合个人快速体验。' : 'No login required, just need API Key. Data stored locally. Perfect for personal use.',
            icon: UserIcon,
            image: '' // Reserved for backend image URL
        },
        {
            title: isZh ? '云端后台模式' : 'Frontend Mode',
            desc: isZh ? '需登录，对接完整管理后台。支持单点登录与企业级功能。' : 'Login required, full admin panel integration. Supports SSO and enterprise features.',
            icon: MicIcon,
            image: '' // Reserved for backend image URL
        },
        {
            title: isZh ? '演示样例模式' : 'Demo Mode',
            desc: isZh ? '包含各类 Demo 示例，使用 Mock 数据。快速了解系统功能。' : 'Contains various demo examples with mock data. Quick system overview.',
            icon: ImageIcon,
            image: '' // Reserved for backend image URL
        }
    ];
});

const toggleLang = () => {
  lang.value = lang.value === 'en' ? 'zh' : 'en';
};

const toggleTheme = () => {
  isDark.value = !isDark.value;
};

// Multiple Snakes Logic
interface Snake {
    body: { x: number; y: number }[];
    direction: { x: number; y: number };
    color: string;
    isHero: boolean;
    circleAngle?: number;
}

interface Food {
    x: number;
    y: number;
}

let animationId: number;
const cellSize = 20;
const snakes = ref<Snake[]>([]);
const foods = ref<Food[]>([]);
const snakeCount = ref(1); // Configurable snake count
let lastUpdate = 0;
const updateInterval = 80;
const foodCount = 3;

const initSnakes = (width: number, height: number) => {
    snakes.value = [];
    
    // Hero Snake (Cyan, follows mouse)
    snakes.value.push(createSnake(width, height, isDark.value ? '#06b6d4' : '#0891b2', true)); // Cyan
    
    // Ambient Snakes
    const colors = isDark.value ? ['#8b5cf6', '#10b981', '#ec4899', '#3b82f6'] : ['#7c3aed', '#059669', '#db2777', '#2563eb'];
    for (let i = 0; i < snakeCount.value; i++) {
        snakes.value.push(createSnake(width, height, colors[i % colors.length], false));
    }
    
    // Spawn initial foods
    foods.value = [];
    for (let i = 0; i < foodCount; i++) {
        spawnFood(width, height);
    }
};

const createSnake = (width: number, height: number, color: string, isHero: boolean): Snake => {
    const cols = Math.floor(width / cellSize);
    const rows = Math.floor(height / cellSize);
    const startX = Math.floor(Math.random() * cols);
    const startY = Math.floor(Math.random() * rows);
    const body = [];
    for (let i = 0; i < 5; i++) {
        body.push({ x: startX, y: startY });
    }
    return {
        body,
        direction: { x: 1, y: 0 },
        color,
        isHero,
        circleAngle: 0
    };
};

const spawnFood = (width: number, height: number) => {
    const cols = Math.floor(width / cellSize);
    const rows = Math.floor(height / cellSize);
    foods.value.push({
        x: Math.floor(Math.random() * cols),
        y: Math.floor(Math.random() * rows)
    });
};

const updateSnakes = (width: number, height: number) => {
    snakes.value.forEach(snake => {
        if (snake.isHero && isMouseInWindow.value) {
            if (isMouseMoving.value) {
                // Hero follows mouse
                const head = snake.body[0];
                const headPixelX = head.x * cellSize + cellSize / 2;
                const headPixelY = head.y * cellSize + cellSize / 2;
                const dx = mouseX.value - headPixelX;
                const dy = mouseY.value - headPixelY;
                if (Math.abs(dx) > Math.abs(dy)) {
                    const newDirX = dx > 0 ? 1 : -1;
                    if (newDirX !== -snake.direction.x) snake.direction = { x: newDirX, y: 0 };
                } else {
                    const newDirY = dy > 0 ? 1 : -1;
                    if (newDirY !== -snake.direction.y) snake.direction = { x: 0, y: newDirY };
                }
            } else {
                // Circle around mouse when idle
                if (snake.circleAngle === undefined) snake.circleAngle = 0;
                snake.circleAngle += 0.05;
                
                const radius = 8; // cells radius
                const centerX = Math.floor(mouseX.value / cellSize);
                const centerY = Math.floor(mouseY.value / cellSize);
                const targetX = centerX + Math.round(Math.cos(snake.circleAngle) * radius);
                const targetY = centerY + Math.round(Math.sin(snake.circleAngle) * radius);
                
                const head = snake.body[0];
                const dx = targetX - head.x;
                const dy = targetY - head.y;
                
                if (Math.abs(dx) > Math.abs(dy)) {
                    const newDirX = dx > 0 ? 1 : -1;
                    if (newDirX !== -snake.direction.x) snake.direction = { x: newDirX, y: 0 };
                } else {
                    const newDirY = dy > 0 ? 1 : -1;
                    if (newDirY !== -snake.direction.y) snake.direction = { x: 0, y: newDirY };
                }
            }
        } else if (!snake.isHero) {
            // Random movement for ambient snakes
            if (Math.random() < 0.05) {
                 const moves = [{x:1,y:0}, {x:-1,y:0}, {x:0,y:1}, {x:0,y:-1}];
                 const validMoves = moves.filter(m => m.x !== -snake.direction.x || m.y !== -snake.direction.y);
                 snake.direction = validMoves[Math.floor(Math.random() * validMoves.length)];
            }
        } else {
            // Hero snake random movement when mouse not in window
            if (Math.random() < 0.05) {
                 const moves = [{x:1,y:0}, {x:-1,y:0}, {x:0,y:1}, {x:0,y:-1}];
                 const validMoves = moves.filter(m => m.x !== -snake.direction.x || m.y !== -snake.direction.y);
                 snake.direction = validMoves[Math.floor(Math.random() * validMoves.length)];
            }
        }

        const newHead = {
            x: snake.body[0].x + snake.direction.x,
            y: snake.body[0].y + snake.direction.y
        };

        // Wrap around logic
        const cols = Math.floor(width / cellSize);
        const rows = Math.floor(height / cellSize);
        if (newHead.x < 0) newHead.x = cols - 1;
        if (newHead.x >= cols) newHead.x = 0;
        if (newHead.y < 0) newHead.y = rows - 1;
        if (newHead.y >= rows) newHead.y = 0;

        snake.body.unshift(newHead);
        
        // Check food collision (only hero snake can eat)
        let ate = false;
        if (snake.isHero) {
            foods.value = foods.value.filter(food => {
                if (newHead.x === food.x && newHead.y === food.y) {
                    ate = true;
                    spawnFood(width, height);
                    return false;
                }
                return true;
            });
        }
        
        if (!ate) {
            snake.body.pop();
        }
    });
};

const drawGame = (ctx: CanvasRenderingContext2D, width: number, height: number) => {
    ctx.clearRect(0, 0, width, height);
    
    // Draw foods with pulsing effect
    const time = Date.now() / 200;
    foods.value.forEach(food => {
        const pulse = Math.sin(time) * 0.3 + 0.7;
        const foodSize = cellSize * 0.8;
        const offset = (cellSize - foodSize * pulse) / 2;
        
        ctx.fillStyle = isDark.value ? '#fbbf24' : '#f59e0b'; // Amber/Gold
        ctx.shadowBlur = 20 * pulse;
        ctx.shadowColor = '#fbbf24';
        ctx.globalAlpha = pulse;
        
        ctx.beginPath();
        ctx.arc(
            food.x * cellSize + cellSize / 2,
            food.y * cellSize + cellSize / 2,
            (foodSize * pulse) / 2,
            0,
            Math.PI * 2
        );
        ctx.fill();
        
        ctx.globalAlpha = 1;
        ctx.shadowBlur = 0;
    });
    
    // Draw snakes
    snakes.value.forEach(snake => {
        ctx.fillStyle = snake.color;
        ctx.shadowBlur = snake.isHero ? 15 : 10;
        ctx.shadowColor = snake.color;
        
        snake.body.forEach((segment, index) => {
            const alpha = 1 - index / snake.body.length;
            ctx.globalAlpha = Math.max(0.2, alpha * (snake.isHero ? 0.9 : 0.8));
            
            // Hero snake has rounded corners and bigger size
            if (snake.isHero) {
                const size = cellSize - 1;
                const radius = 3;
                ctx.beginPath();
                ctx.roundRect(segment.x * cellSize, segment.y * cellSize, size, size, radius);
                ctx.fill();
            } else {
                ctx.fillRect(segment.x * cellSize, segment.y * cellSize, cellSize - 2, cellSize - 2);
            }
        });
        ctx.globalAlpha = 1;
        ctx.shadowBlur = 0;
    });
};

const gameLoop = (timestamp: number) => {
    if (!canvasRef.value) return;
    const canvas = canvasRef.value;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    if (timestamp - lastUpdate > updateInterval) {
        updateSnakes(canvas.width, canvas.height);
        lastUpdate = timestamp;
    }
    
    drawGame(ctx, canvas.width, canvas.height);
    animationId = requestAnimationFrame(gameLoop);
};

const handleResize = () => {
  if (canvasRef.value) {
    canvasRef.value.width = window.innerWidth;
    canvasRef.value.height = window.innerHeight;
    initSnakes(canvasRef.value.width, canvasRef.value.height);
  }
};

onMounted(() => {
  if (canvasRef.value) {
      canvasRef.value.width = window.innerWidth;
      canvasRef.value.height = window.innerHeight;
      initSnakes(canvasRef.value.width, canvasRef.value.height);
      animationId = requestAnimationFrame(gameLoop);
  }
  
  window.addEventListener('mousemove', handleMouseMove);
  window.addEventListener('mouseleave', handleMouseLeave);
  window.addEventListener('resize', handleResize);
  
  startCarousel();
  
  setupScrollReveal(feature1);
  setupScrollReveal(feature2);
  setupScrollReveal(feature3);
});

onUnmounted(() => {
  cancelAnimationFrame(animationId);
  window.removeEventListener('mousemove', handleMouseMove);
  window.removeEventListener('mouseleave', handleMouseLeave);
  window.removeEventListener('resize', handleResize);
  clearInterval(carouselInterval);
  clearTimeout(mouseTimeout);
});
</script>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes gradient-x {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.animate-fade-in-up {
  animation: fade-in-up 1s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
}

.animate-gradient-x {
  background-size: 200% 200%;
  animation: gradient-x 5s ease infinite;
}

.delay-100 { animation-delay: 100ms; }
.delay-200 { animation-delay: 200ms; }
.delay-300 { animation-delay: 300ms; }
.delay-500 { animation-delay: 500ms; }

/* Scroll Reveal Styles */
.scroll-reveal {
  opacity: 0;
  transform: translateY(60px);
  transition: all 1.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.scroll-reveal.revealed {
  opacity: 1;
  transform: translateY(0);
}

@keyframes bounce-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.animate-bounce-slow {
  animation: bounce-slow 2s ease-in-out infinite;
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(100%);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
