<template>
  <div>
    <!-- 吸顶排序栏 -->
    <div class="sticky top-14 z-30 border-b border-border/60"
         style="background:rgba(13,13,15,0.92);backdrop-filter:blur(16px)">
      <div class="max-w-screen-2xl mx-auto px-3 sm:px-4">
        <div class="py-2.5">
          <SortBar :modelValue="activeSort" @update:modelValue="setSort"/>
        </div>
      </div>
    </div>

    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">
      <!-- 广告位 -->
      <div class="mb-4 rounded-xl border border-border/40 bg-bg-surface/60 h-14
                  flex items-center justify-center text-text-muted text-xs tracking-widest">
        AD · mitun69 广告位
      </div>

      <!-- 视频瀑布流（按顺序分发到各列：左→右→下一行） -->
      <div class="flex gap-3 items-start">
        <div v-for="(col, colIdx) in videoColumns" :key="`col-${colIdx}`" class="flex-1 space-y-3">
          <div v-for="v in col" :key="v.video_id">
            <VideoCard :video="v" :order-by="activeSort"/>
          </div>
        </div>

        <template v-if="loading">
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
        </template>
      </div>

      <div ref="sentinel" class="h-6 mt-2"></div>
      <p v-if="!hasMore && videos.length" class="text-center py-10 text-text-muted text-xs">── 已经到底了 ──</p>
      <p v-if="!loading && !videos.length" class="text-center py-20 text-text-muted">暂无内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useHead } from '@vueuse/head'
import VideoCard from '@/components/common/VideoCard.vue'
import SortBar from '@/components/common/SortBar.vue'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import { videoApi } from '@/api'

useHead({
  title: 'mitun69 - 免费短视频在线观看',
  meta: [
    { name: 'description', content: '海量高清短视频，国产自拍、日韩、欧美、无码精品，每日更新，免费在线观看。' },
    { property: 'og:title', content: 'mitun69 - 免费短视频在线观看' },
    { property: 'og:description', content: '海量高清短视频，国产自拍、日韩、欧美、无码精品，每日更新，免费在线观看。' },
    { property: 'og:type', content: 'website' },
  ],
})

const videos = ref([])
const page = ref(1)
const activeSort = ref('created_at')
const columnCount = ref(2)

const videoColumns = computed(() => {
  const cols = Array.from({ length: columnCount.value }, () => [])
  videos.value.forEach((video, idx) => {
    cols[idx % columnCount.value].push(video)
  })
  return cols
})

const skeletonColumns = computed(() => {
  const total = 10
  const cols = Array.from({ length: columnCount.value }, () => [])
  for (let i = 0; i < total; i++) {
    cols[i % columnCount.value].push(i)
  }
  return cols
})

function syncColumnCount() {
  const w = window.innerWidth
  if (w >= 1280) columnCount.value = 5
  else if (w >= 768) columnCount.value = 4
  else if (w >= 640) columnCount.value = 3
  else columnCount.value = 2
}

async function fetchMore() {
  const res = await videoApi.list({ page: page.value, page_size: 20, order_by: activeSort.value })
  const list = res.data?.videos || []
  videos.value.push(...list)
  page.value++
  if (list.length < 20) return false
}

const { sentinel, loading, hasMore, reset } = useInfiniteScroll(fetchMore)

function setSort(v) {
  if (activeSort.value === v) return
  activeSort.value = v
  videos.value = []; page.value = 1; reset()
}

onMounted(() => {
  syncColumnCount()
  window.addEventListener('resize', syncColumnCount)
})

onUnmounted(() => {
  window.removeEventListener('resize', syncColumnCount)
})
</script>
