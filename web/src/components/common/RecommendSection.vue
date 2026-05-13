<template>
  <div v-if="current.length || loading">
    <div class="flex items-center justify-between mb-3">
      <h2 class="section-title text-sm">
        <svg class="w-4 h-4 text-red-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
        猜你喜欢
      </h2>
      <button @click="refresh" :disabled="refreshing"
              class="flex items-center gap-1 text-xs text-text-primary/60 hover:text-primary transition-colors disabled:opacity-40">
        <svg class="w-3.5 h-3.5 transition-transform duration-500"
             :class="refreshing ? 'animate-spin' : ''"
             fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
        </svg>
        换一批
      </button>
    </div>

    <!-- 骨架屏 -->
    <div v-if="loading" class="flex gap-3 items-start">
      <div v-for="colIdx in 2" :key="colIdx" class="flex-1 space-y-3">
        <div v-for="i in 3" :key="i" class="rounded-xl bg-bg-card border border-border animate-pulse">
          <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
          <div class="p-2.5 space-y-1.5">
            <div class="h-2.5 bg-bg-hover rounded w-full"></div>
            <div class="h-2 bg-bg-hover rounded w-2/3"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 瀑布流：与首页相同结构 -->
    <div v-else class="flex gap-3 items-start">
      <div v-for="(col, colIdx) in columns" :key="colIdx" class="flex-1 space-y-3">
        <VideoCard v-for="v in col" :key="v.video_id" :video="v"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import VideoCard from '@/components/common/VideoCard.vue'
import { videoApi } from '@/api'

const PAGE_SIZE = 6

const pool = ref([])       // 预拉取的视频池
const offset = ref(0)      // 当前展示窗口的起始位置
const loading = ref(true)
const refreshing = ref(false)

// 当前展示的视频：确保偶数个，避免瀑布流底部出现空缺
const current = computed(() => {
  const items = pool.value.slice(offset.value, offset.value + PAGE_SIZE)
  return items.length % 2 === 0 ? items : items.slice(0, -1)
})

const columns = computed(() => {
  const cols = [[], []]
  current.value.forEach((v, i) => cols[i % 2].push(v))
  return cols
})

async function fetchPool(append = false) {
  try {
    // 多拉一些，保证换一批有内容
    const res = await videoApi.recommend(30)
    const videos = res.data?.videos || []
    const valid = videos.filter(v => v?.video_id && v?.cover_url)
    pool.value = append ? [...pool.value, ...valid] : valid
  } catch {
    if (!append) pool.value = []
  }
}

async function refresh() {
  refreshing.value = true
  const next = offset.value + PAGE_SIZE
  // 池子剩余不足一批时重新拉取并追加
  if (next + PAGE_SIZE > pool.value.length) {
    await fetchPool(true)
  }
  offset.value = next < pool.value.length ? next : 0
  refreshing.value = false
}

onMounted(async () => {
  await fetchPool()
  loading.value = false
})
</script>
