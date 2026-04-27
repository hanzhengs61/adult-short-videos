<template>
  <div>
    <!-- 吸顶排序栏 -->
    <div class="sticky top-14 z-30 border-b border-border/60"
         style="background:rgba(13,13,15,0.92);backdrop-filter:blur(16px)">
      <div class="max-w-screen-2xl mx-auto px-3 sm:px-4">
        <!-- 排序 -->
        <div class="flex items-center gap-1.5 flex-wrap py-2.5">
          <div class="flex items-center gap-1.5 shrink-0">
            <button v-for="s in sorts" :key="s.value" @click="setSort(s.value)"
                    :class="['shrink-0 flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all',
                activeSort===s.value
                  ? 'bg-gradient-primary border-transparent text-white shadow-glow'
                  : 'border-border bg-bg-surface text-text-muted hover:text-text-primary hover:border-border-light']">
              <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2"
                   stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
                <path :d="s.iconPath"/>
              </svg>
              {{ s.label }}
            </button>
          </div>
          <button v-for="r in regions" :key="r.value" @click="setRegion(r.value)"
                  :class="['shrink-0 px-3 py-1.5 rounded-full text-xs border transition-all',
              activeRegion===r.value
                ? 'border-primary/60 text-primary bg-primary/10'
                : 'border-border text-text-muted hover:text-text-primary']">
            {{ r.label }}
          </button>
        </div>
      </div>
    </div>

    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">
      <!-- 广告位（首屏顶部）-->
      <div class="mb-4 rounded-xl border border-border/40 bg-bg-surface/60 h-14
                  flex items-center justify-center text-text-muted text-xs tracking-widest">
        AD · mitun69 广告位
      </div>

      <!-- 视频网格 -->
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5 gap-3">
        <template v-for="(v, idx) in videos" :key="v.video_id">
          <VideoCard :video="v" />
          <!-- 每 20 条插入广告 -->
          <div v-if="(idx+1) % 20 === 0"
            class="col-span-full rounded-xl border border-border/40 bg-bg-surface/60 h-12
                   flex items-center justify-center text-text-muted text-xs tracking-widest">
            AD · 广告位
          </div>
        </template>

        <!-- 骨架屏 -->
        <template v-if="loading">
          <div v-for="i in 10" :key="`sk${i}`"
            class="rounded-xl bg-bg-card border border-border animate-pulse">
            <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
            <div class="p-2.5 space-y-1.5">
              <div class="h-2.5 bg-bg-hover rounded w-full"></div>
              <div class="h-2.5 bg-bg-hover rounded w-3/4"></div>
              <div class="h-2 bg-bg-hover rounded w-1/2 mt-2"></div>
            </div>
          </div>
        </template>
      </div>

      <!-- 无限滚动哨兵 -->
      <div ref="sentinel" class="h-6 mt-2"></div>

      <p v-if="!hasMore && videos.length" class="text-center py-10 text-text-muted text-xs">
        ── 已经到底了 ──
      </p>
      <p v-if="!loading && !videos.length" class="text-center py-20 text-text-muted">暂无内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import VideoCard from '@/components/common/VideoCard.vue'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import { videoApi } from '@/api'

const videos = ref([])
const page = ref(1)
const activeSort = ref('created_at')
const activeRegion = ref('')

const sorts = [
  { label: '最新', value: 'created_at', iconPath: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' },
  { label: '播放最多', value: 'play_count', iconPath: 'M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z' },
  { label: '评论最多', value: 'comment_count', iconPath: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z' },
  { label: '收藏最多', value: 'favorite_count', iconPath: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z' },
]

async function fetchMore() {
  const res = await videoApi.list({ page: page.value, page_size: 20,
    order_by: activeSort.value, region: activeRegion.value || undefined })
  const list = res.data?.list || []
  videos.value.push(...list)
  page.value++
  if (list.length < 20) return false
}

const { sentinel, loading, hasMore, reset } = useInfiniteScroll(fetchMore)

function setSort(v) {
  if (activeSort.value === v) return
  activeSort.value = v; videos.value = []; page.value = 1; reset()
}
function setRegion(v) {
  if (activeRegion.value === v) return
  activeRegion.value = v; videos.value = []; page.value = 1; reset()
}
</script>
