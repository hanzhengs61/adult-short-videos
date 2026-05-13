<template>
  <header class="fixed top-0 left-0 right-0 z-50 border-b border-border/50"
          style="background:rgba(13,13,15,0.9);backdrop-filter:blur(20px)">
    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 h-14 flex items-center gap-2">

      <!-- 汉堡菜单（预留分类侧栏入口） -->
      <button @click="sidebarOpen = !sidebarOpen"
              class="w-8 h-8 rounded-lg text-text-primary/60 flex items-center justify-center shrink-0">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
        </svg>
      </button>

      <!-- Logo -->
      <router-link to="/" class="flex items-center shrink-0 group">
        <img src="/logo.png" alt="MITUN69"
             class="h-10 w-auto transition-all duration-300 group-hover:drop-shadow-[0_0_10px_rgba(255,100,150,0.5)]"/>
      </router-link>

      <!-- 桌面导航（md+）-->
      <nav class="hidden md:flex items-center gap-0.5 ml-1">
        <router-link v-for="n in desktopNav" :key="n.to" :to="n.to"
                     :class="['px-3 py-1.5 rounded-lg text-base font-medium transition-all',
            isActive(n.to)
              ? 'text-primary bg-primary/10'
              : 'text-text-primary/60 hover:text-text-primary hover:bg-bg-hover']">
          {{ n.label }}
        </router-link>
      </nav>

      <div class="flex-1"></div>

      <!-- 搜索图标（跳转探索页） -->
      <router-link to="/explore"
                   class="w-8 h-8 rounded-lg text-text-primary/60 flex items-center justify-center shrink-0">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2"
             stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
      </router-link>

      <!-- 添加到主屏幕（预留升级为下载 App 入口） -->
      <button @click="addToHomescreen"
              class="flex items-center gap-1 px-2 h-8 rounded-lg border border-border text-text-primary/60 shrink-0">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
        </svg>
        <span class="text-base font-bold">添加主屏幕</span>
      </button>
    </div>
  </header>
  <div class="h-14"></div>
</template>

<script setup>
import {onMounted, ref} from 'vue'
import {useRoute} from 'vue-router'

const route = useRoute()
const sidebarOpen = ref(false)  // 预留侧栏开关
const deferredPrompt = ref(null)

const desktopNav = [
  {label: '首页', to: '/'},
  {label: '吃瓜', to: '/gossip'},
  {label: '探索', to: '/explore'},
  {label: '福利', to: '/rankings'},
]

function isActive(to) {
  return to === '/' ? route.path === '/' : route.path.startsWith(to)
}

async function addToHomescreen() {
  if (deferredPrompt.value) {
    deferredPrompt.value.prompt()
    await deferredPrompt.value.userChoice
    deferredPrompt.value = null
  } else {
    alert('请在浏览器菜单（分享/设置）中选择「添加到主屏幕」')
  }
}

onMounted(() => {
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt.value = e
  })
})
</script>
