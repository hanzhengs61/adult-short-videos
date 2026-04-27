<template>
  <header class="fixed top-0 left-0 right-0 z-50 border-b border-border/50"
          style="background:rgba(13,13,15,0.9);backdrop-filter:blur(20px)">
    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 h-14 flex items-center gap-3">

      <!-- Logo -->
      <router-link to="/" class="flex items-center shrink-0 group">
        <img src="/logo.png" alt="MITUN69"
             class="h-10 w-auto transition-all duration-300 group-hover:drop-shadow-[0_0_10px_rgba(255,100,150,0.5)]" />
      </router-link>

      <!-- 桌面导航（md+）-->
      <nav class="hidden md:flex items-center gap-0.5 ml-1">
        <router-link v-for="n in desktopNav" :key="n.to" :to="n.to"
          :class="['px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
            isActive(n.to)
              ? 'text-primary bg-primary/10'
              : 'text-text-muted hover:text-text-primary hover:bg-bg-hover']">
          {{ n.label }}
        </router-link>
      </nav>

      <!-- 搜索框 -->
      <form class="flex-1 max-w-sm mx-auto" @submit.prevent="handleSearch">
        <div class="relative">
          <input v-model="q" type="text" placeholder="搜索…"
            class="w-full h-9 pl-9 pr-4 bg-bg-surface border border-border rounded-full text-text-primary
                   text-xs placeholder-text-muted outline-none transition-all
                   focus:border-primary/60 focus:ring-1 focus:ring-primary/20" />
          <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none"
               fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>
          </svg>
        </div>
      </form>

      <!-- 用户区 -->
      <div class="flex items-center gap-2 shrink-0">
        <template v-if="userStore.isLoggedIn">
          <!-- 上传按钮 -->
          <router-link to="/profile" class="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-lg
            border border-border text-text-muted text-xs hover:border-primary/40 hover:text-primary
            transition-all">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
            </svg>
            上传
          </router-link>

          <!-- 头像下拉 -->
          <div class="relative" v-click-outside="() => dropdownOpen = false">
            <button @click="dropdownOpen = !dropdownOpen"
              class="w-8 h-8 rounded-full bg-gradient-primary flex items-center justify-center
                     text-white text-sm font-bold hover:shadow-glow transition-all">
              {{ userStore.userInfo?.username?.[0]?.toUpperCase() || 'U' }}
            </button>
            <transition name="dropdown">
              <div v-if="dropdownOpen"
                class="absolute right-0 top-10 w-44 bg-bg-surface border border-border rounded-xl shadow-2xl py-1 z-50">
                <div class="px-3 py-2.5 border-b border-border">
                  <p class="text-sm font-bold text-text-primary truncate">{{ userStore.userInfo?.username }}</p>
                  <p class="text-xs text-text-muted mt-0.5">免费用户</p>
                </div>
                <router-link to="/profile" @click="dropdownOpen=false"
                  class="flex items-center gap-2 px-3 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-bg-hover">
                  我的主页
                </router-link>
                <router-link to="/subscribe" @click="dropdownOpen=false"
                  class="flex items-center gap-2 px-3 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-bg-hover">
                  我的订阅
                </router-link>
                <hr class="border-border my-1" />
                <button @click="handleLogout"
                  class="w-full flex items-center gap-2 px-3 py-2.5 text-sm text-red-400 hover:bg-bg-hover transition-colors">
                  退出登录
                </button>
              </div>
            </transition>
          </div>
        </template>

        <template v-else>
          <button @click="userStore.openAuth('login')" class="btn-ghost hidden sm:inline-flex text-xs py-1.5 px-3">
            登录
          </button>
          <button @click="userStore.openAuth('register')" class="btn-primary text-xs px-4 py-1.5">
            注册
          </button>
        </template>

        <!-- 更多选项 -->
        <div class="relative" v-click-outside="() => moreOpen = false">
          <button @click="moreOpen = !moreOpen"
            class="w-8 h-8 rounded-lg border border-border text-text-muted flex items-center justify-center
                   hover:border-border-light hover:text-text-primary transition-all">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <circle cx="5" cy="12" r="1.5"/><circle cx="12" cy="12" r="1.5"/><circle cx="19" cy="12" r="1.5"/>
            </svg>
          </button>
          <transition name="dropdown">
            <div v-if="moreOpen"
              class="absolute right-0 top-10 w-48 bg-bg-surface border border-border rounded-xl shadow-2xl py-1 z-50">
              <button @click="addToHomescreen"
                class="w-full flex items-center gap-2.5 px-3 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-bg-hover transition-colors">
                <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"/>
                </svg>
                添加到主屏幕
              </button>
              <button @click="shareApp"
                class="w-full flex items-center gap-2.5 px-3 py-2.5 text-sm text-text-secondary hover:text-text-primary hover:bg-bg-hover transition-colors">
                <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
                </svg>
                分享
              </button>
              <hr class="border-border my-1" />
              <div class="px-3 py-2">
                <p class="text-[10px] text-text-muted leading-relaxed">mitun69.com · 高清视频平台</p>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
  <div class="h-14"></div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const q = ref('')
const dropdownOpen = ref(false)
const moreOpen = ref(false)
const deferredPrompt = ref(null)

const desktopNav = [
  { label: '首页', to: '/' },
  { label: '订阅', to: '/subscribe' },
  { label: '探索', to: '/explore' },
  { label: '榜单', to: '/rankings' },
]

const vClickOutside = {
  mounted(el, b) { el._co = e => { if (!el.contains(e.target)) b.value(e) }; document.addEventListener('click', el._co) },
  unmounted(el) { document.removeEventListener('click', el._co) },
}

function isActive(to) {
  return to === '/' ? route.path === '/' : route.path.startsWith(to)
}

function handleSearch() {
  const kw = q.value.trim()
  if (!kw) return
  router.push({ path: '/explore', query: { q: kw } })
  q.value = ''
}

async function handleLogout() {
  dropdownOpen.value = false
  await userStore.logout()
  router.push('/')
}

async function addToHomescreen() {
  moreOpen.value = false
  if (deferredPrompt.value) {
    deferredPrompt.value.prompt()
    await deferredPrompt.value.userChoice
    deferredPrompt.value = null
  } else {
    alert('请在浏览器菜单（分享/设置）中选择「添加到主屏幕」')
  }
}

async function shareApp() {
  moreOpen.value = false
  const shareData = {
    title: 'mitun69 - 高清视频平台',
    text: '高清视频，极致体验 — mitun69.com',
    url: window.location.href,
  }
  if (navigator.share) {
    try { await navigator.share(shareData) } catch {}
  } else {
    await navigator.clipboard.writeText(window.location.href).catch(() => {})
    alert('链接已复制到剪贴板')
  }
}

onMounted(() => {
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt.value = e
  })
})
</script>

<style scoped>
.dropdown-enter-active { transition: all 0.15s ease-out; }
.dropdown-leave-active { transition: all 0.1s ease; }
.dropdown-enter-from  { opacity: 0; transform: translateY(-6px) scale(0.97); }
.dropdown-leave-to    { opacity: 0; }
</style>
