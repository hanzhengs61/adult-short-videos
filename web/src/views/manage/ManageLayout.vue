<template>
  <div class="min-h-screen flex bg-bg" style="font-family:Inter,system-ui,sans-serif">

    <!-- 侧边栏 -->
    <aside :class="['fixed inset-y-0 left-0 z-50 flex flex-col border-r border-border transition-all duration-300',
                    sideOpen ? 'w-52' : 'w-14']"
           style="background:#0f0f12">
      <!-- Logo -->
      <div class="flex items-center gap-2.5 px-3.5 h-14 border-b border-border shrink-0">
        <div class="w-7 h-7 rounded-lg bg-gradient-primary flex items-center justify-center shrink-0">
          <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
        </div>
        <span v-if="sideOpen" class="text-sm font-bold text-text-primary whitespace-nowrap">管理后台</span>
      </div>

      <!-- 菜单 -->
      <nav class="flex-1 overflow-y-auto py-3 space-y-0.5 px-2">
        <router-link v-for="item in menu" :key="item.to" :to="item.to"
                     :class="['flex items-center gap-3 px-2.5 py-2 rounded-lg text-sm transition-all group',
                              isActive(item.to)
                                ? 'bg-primary/15 text-primary'
                                : 'text-text-secondary hover:text-text-primary hover:bg-bg-hover']">
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" :d="item.icon"/>
          </svg>
          <span v-if="sideOpen" class="truncate">{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- 收起/展开 -->
      <button class="flex items-center justify-center h-10 border-t border-border text-text-muted hover:text-text-primary transition-colors shrink-0"
              @click="sideOpen = !sideOpen">
        <svg class="w-4 h-4 transition-transform" :class="sideOpen ? '' : 'rotate-180'"
             fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>
    </aside>

    <!-- 主内容区 -->
    <div :class="['flex-1 flex flex-col min-w-0 transition-all duration-300', sideOpen ? 'ml-52' : 'ml-14']">
      <!-- 顶部栏 -->
      <header class="sticky top-0 z-40 flex items-center justify-between px-5 h-14 border-b border-border/60 shrink-0"
              style="background:rgba(13,13,15,0.92);backdrop-filter:blur(16px)">
        <h1 class="text-sm font-semibold text-text-primary">{{ currentTitle }}</h1>
        <div class="flex items-center gap-2 text-xs text-text-muted">
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
          管理员
          <button class="ml-3 px-3 py-1 rounded-lg border border-border hover:border-border-light hover:text-text-primary transition-all"
                  @click="$router.push('/')">
            返回前台
          </button>
        </div>
      </header>

      <!-- 页面内容 -->
      <main class="flex-1 p-5 overflow-auto">
        <router-view/>
      </main>
    </div>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const sideOpen = ref(true)

const menu = [
  { label: '数据概览',  to: '/manage/dashboard', icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z' },
  { label: '视频管理',  to: '/manage/videos',    icon: 'M15 10l4.553-2.069A1 1 0 0121 8.82v6.36a1 1 0 01-1.447.89L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z' },
  { label: '评论管理',  to: '/manage/comments',  icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z' },
  { label: '用户管理',  to: '/manage/users',     icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
  { label: 'Banner',   to: '/manage/banners',   icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z' },
  { label: 'App广告',  to: '/manage/apps',      icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z' },
  { label: '吃瓜文章', to: '/manage/gossip',    icon: 'M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z' },
]

const titleMap = Object.fromEntries(menu.map(m => [m.to, m.label]))
const currentTitle = computed(() => titleMap[route.path] || '管理后台')

function isActive(to) {
  return route.path === to || route.path.startsWith(to + '/')
}
</script>
