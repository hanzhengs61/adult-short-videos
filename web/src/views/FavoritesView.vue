<template>
  <div>
    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">

      <!-- 骨架屏 -->
      <div v-if="loading && !videos.length" class="flex gap-3 items-start">
        <div v-for="(col, colIdx) in skeletonColumns" :key="`sk-col-${colIdx}`" class="flex-1 space-y-3">
          <div v-for="i in col" :key="`sk${colIdx}-${i}`"
               class="rounded-xl bg-bg-card border border-border animate-pulse">
            <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
            <div class="p-2.5 space-y-1.5">
              <div class="h-2.5 bg-bg-hover rounded w-full"></div>
              <div class="h-2.5 bg-bg-hover rounded w-3/4"></div>
              <div class="h-2 bg-bg-hover rounded w-1/2 mt-2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 瀑布流 -->
      <div class="flex gap-3 items-start">
        <div v-for="(col, colIdx) in videoColumns" :key="`col-${colIdx}`" class="flex-1 space-y-3">
          <VideoCard v-for="v in col" :key="v.video_id" :video="v" order-by="created_at"/>
        </div>
      </div>

      <div ref="sentinel" class="h-6 mt-2"></div>
      <p v-if="!hasMore && videos.length" class="text-center py-10 text-text-muted text-xs">── 已经到底了 ──</p>

      <!-- 未登录提示 -->
      <div v-if="!loading && !videos.length && !userStore.isLoggedIn"
           class="flex flex-col items-center py-24 gap-4 text-center">
        <svg class="w-16 h-16 text-text-muted opacity-30" fill="none" stroke="currentColor" stroke-width="1.2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
        </svg>
        <p class="text-text-primary font-medium">登录后查看喜欢</p>
        <button @click="userStore.openAuth('login')"
                class="px-5 py-2 rounded-full bg-primary text-white text-sm hover:bg-primary/80 transition-colors">
          去登录
        </button>
      </div>

      <!-- 已登录但无收藏 -->
      <div v-else-if="!loading && !videos.length && userStore.isLoggedIn"
           class="flex flex-col items-center py-24 gap-3 text-center">
        <svg class="w-16 h-16 text-text-muted opacity-30" fill="none" stroke="currentColor" stroke-width="1.2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
        </svg>
        <p class="text-text-primary font-medium">还没有喜欢的内容</p>
        <router-link to="/" class="text-primary text-sm hover:underline">去首页发现内容</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useHead } from '@vueuse/head'
import VideoCard from '@/components/common/VideoCard.vue'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import { favoriteApi } from '@/api'
import { useUserStore } from '@/stores/user'

useHead({ title: '我的喜欢 - mitun69' })

const userStore = useUserStore()
const videos = ref([])
const page = ref(1)
const columnCount = ref(2)

const videoColumns = computed(() => {
  const cols = Array.from({ length: columnCount.value }, () => [])
  videos.value.forEach((v, i) => cols[i % columnCount.value].push(v))
  return cols
})

const skeletonColumns = computed(() => {
  const cols = Array.from({ length: columnCount.value }, () => [])
  for (let i = 0; i < 10; i++) cols[i % columnCount.value].push(i)
  return cols
})

async function fetchMore() {
  // 未登录不请求
  if (!userStore.isLoggedIn) return false
  const res = await favoriteApi.list({ page: page.value, page_size: 20 })
  // 收藏接口返回 { favorites: [{ video_id, video: {...} }] }，取内层 video 对象
  const raw = res.data?.favorites || res.data?.videos || []
  const list = raw.map(item => item.video || item)
  videos.value.push(...list)
  page.value++
  if (list.length < 20) return false
}

const { sentinel, loading, hasMore } = useInfiniteScroll(fetchMore)

function syncColumnCount() {
  const w = window.innerWidth
  if (w >= 1280) columnCount.value = 5
  else if (w >= 768) columnCount.value = 4
  else if (w >= 640) columnCount.value = 3
  else columnCount.value = 2
}

onMounted(() => {
  syncColumnCount()
  window.addEventListener('resize', syncColumnCount)
})

onUnmounted(() => {
  window.removeEventListener('resize', syncColumnCount)
})
</script>
