<template>
  <nav class="fixed bottom-0 left-0 right-0 z-40 md:hidden border-t border-border/50"
       style="background:rgba(13,13,15,0.96);backdrop-filter:blur(20px)">
    <div class="flex items-stretch" style="padding-bottom:env(safe-area-inset-bottom)">
      <button v-for="item in navItems" :key="item.to"
              @click="handleNav(item)"
              :class="['flex-1 flex flex-col items-center justify-center gap-0.5 py-2.5 transition-all relative',
          isActive(item) ? 'text-primary' : 'text-text-primary/60']">

        <!-- 活跃背景光晕 -->
        <span v-if="isActive(item)"
              class="absolute top-0 left-1/2 -translate-x-1/2 w-8 h-0.5 bg-gradient-primary rounded-b-full"></span>

        <!-- 图标（活跃时 filled，否则 outline） -->
        <span class="relative">
          <svg class="w-6 h-6 transition-transform" :class="isActive(item) ? 'scale-110' : ''"
               :fill="isActive(item) ? 'currentColor' : 'none'"
               stroke="currentColor" :stroke-width="isActive(item) ? '0' : '1.8'"
               stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
            <path :d="isActive(item) ? item.iconFill : item.iconOutline"/>
          </svg>
          <!-- 角标（订阅有新内容）-->
          <span v-if="item.badge"
                class="absolute -top-1 -right-1.5 min-w-[14px] h-3.5 px-0.5 rounded-full
                   bg-primary text-white text-[9px] font-bold flex items-center justify-center">
            {{ item.badge }}
          </span>
        </span>

        <span :class="['text-base font-semibold leading-none transition-colors',
          isActive(item) ? 'text-primary' : 'text-text-primary/60']">
          {{ item.label }}
        </span>
      </button>
    </div>
  </nav>
  <!-- 占位 -->
  <div class="h-16 md:hidden" style="padding-bottom:env(safe-area-inset-bottom)"></div>
</template>

<script setup>
import {computed} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 每个 Tab 的 filled 和 outline SVG path
const navItems = computed(() => [
  {
    label: '首页', to: '/',
    iconOutline: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
    iconFill: 'M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z',
  },
  {
    label: '吃瓜', to: '/gossip',
    iconOutline: 'M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z',
    iconFill: 'M9 4.804A7.968 7.968 0 005.5 4c-1.255 0-2.443.29-3.5.804v10A7.969 7.969 0 015.5 14c1.669 0 3.218.51 4.5 1.385A7.962 7.962 0 0114.5 14c1.255 0 2.443.29 3.5.804v-10A7.968 7.968 0 0014.5 4c-1.669 0-3.218.51-4.5 1.385V12a1 1 0 11-2 0V4.804z',
  },
  {
    label: '探索', to: '/explore',
    iconOutline: 'M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z',
    iconFill: 'M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z',
  },
  {
    label: '福利', to: '/rankings',
    iconOutline: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7',
    iconFill: 'M5 3a2 2 0 000 4h5.5A2.5 2.5 0 0010 2.5v0A2.5 2.5 0 007.5 0 2.5 2.5 0 005 2.5V3zm8.5 4H19a2 2 0 000-4h-5.5A2.5 2.5 0 0014 5.5v0A2.5 2.5 0 0016.5 3h-3v.5A2.5 2.5 0 0013.5 6V7zM3 8h7v13H5a2 2 0 01-2-2V8zm11 0h7v11a2 2 0 01-2 2h-5V8z',
  },
  {
    label: '我的', to: '/profile',
    iconOutline: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
    iconFill: 'M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z',
  },
])

function isActive(item) {
  if (item.to === '/') return route.path === '/'
  return route.path.startsWith(item.to)
}

function handleNav(item) {
  if (item.to === '/profile' && !userStore.isLoggedIn) {
    userStore.openAuth('login')
    return
  }
  router.push(item.to)
}
</script>
