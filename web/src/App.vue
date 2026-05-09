<template>
  <NavBar v-if="!isFeed"/>
  <main class="min-h-screen relative">
    <router-view v-slot="{ Component }">
      <transition name="fade">
        <!-- FeedView 不缓存（自身管理复杂播放状态），其余页面缓存以支持后退免请求 -->
        <keep-alive :exclude="['FeedView']">
          <component :is="Component" :key="$route.path"/>
        </keep-alive>
      </transition>
    </router-view>
  </main>

  <!-- 底部 Footer（桌面端，Feed 页隐藏） -->
  <footer v-if="!isFeed" class="hidden md:block border-t border-border mt-16">
    <div class="max-w-screen-2xl mx-auto px-4 py-10">
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-8 mb-8">
        <!-- 品牌 -->
        <div>
          <div class="text-xl font-black mb-3">
            <span class="bg-gradient-primary bg-clip-text text-transparent">mitun</span>
            <span class="text-white">69</span>
          </div>
          <p class="text-text-muted text-xs leading-relaxed">高清视频，极致体验。<br/>全球优质内容聚合平台。</p>
        </div>
        <!-- 导航 -->
        <div>
          <p class="text-text-secondary text-xs font-semibold mb-3 uppercase tracking-widest">发现</p>
          <div class="space-y-2">
            <router-link v-for="l in footerLinks.discover" :key="l.to" :to="l.to"
                         class="block text-text-muted text-xs hover:text-primary transition-colors">{{ l.label }}
            </router-link>
          </div>
        </div>
        <div>
          <p class="text-text-secondary text-xs font-semibold mb-3 uppercase tracking-widest">账号</p>
          <div class="space-y-2">
            <a v-for="l in footerLinks.account" :key="l.label" href="#"
               @click.prevent="l.action && l.action()"
               class="block text-text-muted text-xs hover:text-primary transition-colors cursor-pointer">{{
                l.label
              }}</a>
          </div>
        </div>
        <div>
          <p class="text-text-secondary text-xs font-semibold mb-3 uppercase tracking-widest">关于</p>
          <div class="space-y-2">
            <a v-for="l in footerLinks.about" :key="l.label"
               class="block text-text-muted text-xs hover:text-primary transition-colors">{{ l.label }}</a>
          </div>
        </div>
      </div>
      <div class="border-t border-border pt-6 flex flex-col sm:flex-row items-center justify-between gap-2">
        <p class="text-text-muted text-xs">© 2025 mitun69.com · 仅限 18 岁以上成年人访问</p>
        <p class="text-text-muted text-xs">域名：mitun69.com</p>
      </div>
    </div>
  </footer>

  <!-- 移动端底部导航（Feed 页隐藏） -->
  <BottomNav v-if="!isFeed"/>

  <!-- 全局登录/注册弹窗 -->
  <AuthModal/>

  <!-- 全局 Toast 通知 -->
  <Toast/>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import BottomNav from '@/components/common/BottomNav.vue'
import AuthModal from '@/components/common/AuthModal.vue'
import Toast from '@/components/common/Toast.vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const route = useRoute()
const isFeed = computed(() => route.path === '/feed')

const footerLinks = {
  discover: [
    { label: '首页', to: '/' },
    { label: '探索', to: '/explore' },
    { label: '榜单', to: '/rankings' },
    { label: '订阅', to: '/subscribe' },
  ],
  account: [
    { label: '登录', action: () => userStore.openAuth('login') },
    { label: '注册', action: () => userStore.openAuth('register') },
    { label: '我的喜欢', action: () => userStore.openAuth('login') },
    { label: '个人主页', action: () => userStore.openAuth('login') },
  ],
  about: ['关于我们', '联系客服', '隐私政策', '使用条款'],
}

onMounted(() => userStore.fetchInfo())
</script>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
/* 离场时脱离文档流（并发模式下不撑高布局）并禁用点击 */
.fade-leave-active {
  position: absolute;
  top: 0; left: 0; right: 0;
  pointer-events: none;
}
</style>
