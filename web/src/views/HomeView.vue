<template>
  <div ref="homeRoot">
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

      <!-- 首次加载骨架 -->
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

      <!-- 视频 + 广告卡混合瀑布流（每 4 条视频插 1 个广告卡） -->
      <div class="flex gap-3 items-start">
        <div v-for="(col, colIdx) in videoColumns" :key="`col-${colIdx}`" class="flex-1 space-y-3">
          <template v-for="item in col" :key="item._isAd ? item._adKey : item.video_id">

            <!-- 广告卡：外观与 VideoCard 保持一致 -->
            <div v-if="item._isAd"
                 class="group block rounded-xl overflow-hidden bg-bg-card border border-border/60
                        hover:border-primary/20 transition-all duration-300 cursor-pointer">
              <!-- 封面区：渐变背景替代视频封面 -->
              <div class="relative aspect-video overflow-hidden
                          bg-gradient-to-br from-primary/25 via-purple-900/20 to-bg-card
                          flex flex-col items-center justify-center gap-1.5">
                <svg class="w-7 h-7 text-white/20" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
                </svg>
                <span class="text-white/30 text-[9px] tracking-widest uppercase">Advertisement</span>
                <!-- 右下角广告标签 -->
                <span class="absolute bottom-1.5 right-1.5 px-1.5 py-0.5 rounded
                             bg-black/60 text-white/50 text-[9px]">广告</span>
              </div>
              <!-- 卡片信息区：模仿 VideoCard 布局 -->
              <div class="p-2 space-y-1.5">
                <p class="text-xs font-medium text-text-muted line-clamp-2 leading-snug">
                  广告位招租 · 接入广告联盟
                </p>
                <div class="flex items-center gap-1.5">
                  <img src="/logo.png" alt="ad" class="w-5 h-5 rounded-full shrink-0 object-cover"/>
                  <span class="text-text-muted text-[10px] flex-1 truncate leading-none">mitun69</span>
                </div>
              </div>
            </div>

            <!-- 视频卡 -->
            <VideoCard v-else :video="item" :order-by="activeSort"/>

          </template>
        </div>
      </div>

      <div ref="sentinel" class="h-6 mt-2"></div>
      <p v-if="!hasMore && videos.length" class="text-center py-10 text-text-muted text-xs">── 已经到底了 ──</p>
      <p v-if="!loading && !videos.length" class="text-center py-20 text-text-muted">暂无内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onActivated, onMounted, onUnmounted } from 'vue'
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
const homeRoot = ref(null)
const page = ref(1)
const activeSort = ref('created_at')
const columnCount = ref(2)
const HOME_FEED_RETURN_KEY = 'home_feed_return_anchor'
const HOME_FEED_RETURN_PENDING_KEY = 'home_feed_return_pending'

// 每 AD_EVERY 条真实视频后插入一个广告卡
const AD_EVERY = 4
const mixedList = computed(() => {
  const result = []
  videos.value.forEach((v, i) => {
    result.push(v)
    if ((i + 1) % AD_EVERY === 0) {
      result.push({ _isAd: true, _adKey: `ad_${i}` })
    }
  })
  return result
})

const videoColumns = computed(() => {
  const cols = Array.from({ length: columnCount.value }, () => [])
  mixedList.value.forEach((item, idx) => {
    cols[idx % columnCount.value].push(item)
  })
  return cols
})

const skeletonColumns = computed(() => {
  const total = 10
  const cols = Array.from({ length: columnCount.value }, () => [])
  for (let i = 0; i < total; i++) cols[i % columnCount.value].push(i)
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
  sessionStorage.removeItem(HOME_FEED_RETURN_KEY)
  sessionStorage.removeItem(HOME_FEED_RETURN_PENDING_KEY)
  scrollWithoutSmooth(0)
  activeSort.value = v
  videos.value = []; page.value = 1; reset()
  nextTick(() => scrollWithoutSmooth(0))
}

function scrollWithoutSmooth(top) {
  const html = document.documentElement
  const prevScrollBehavior = html.style.scrollBehavior
  html.style.scrollBehavior = 'auto'
  window.scrollTo(0, Math.max(0, top))
  requestAnimationFrame(() => {
    html.style.scrollBehavior = prevScrollBehavior
  })
}

function findVideoCard(videoId) {
  if (!homeRoot.value || !videoId) return null
  return Array
    .from(homeRoot.value.querySelectorAll('[data-feed-video-id]'))
    .find(el => el.dataset.feedVideoId === String(videoId))
}

async function restoreFeedReturnAnchor() {
  if (sessionStorage.getItem(HOME_FEED_RETURN_PENDING_KEY) !== '1') {
    sessionStorage.removeItem(HOME_FEED_RETURN_KEY)
    return
  }
  let anchor = null
  try {
    anchor = JSON.parse(sessionStorage.getItem(HOME_FEED_RETURN_KEY) || 'null')
  } catch {}
  if (!anchor?.videoId) {
    sessionStorage.removeItem(HOME_FEED_RETURN_PENDING_KEY)
    return
  }

  await nextTick()
  const restoreOnce = () => {
    const el = findVideoCard(anchor.videoId)
    if (!el) return false
    const rect = el.getBoundingClientRect()
    const targetTop = window.scrollY + rect.top - (Number(anchor.top) || 0)
    scrollWithoutSmooth(targetTop)
    return true
  }

  if (!restoreOnce()) {
    sessionStorage.removeItem(HOME_FEED_RETURN_KEY)
    sessionStorage.removeItem(HOME_FEED_RETURN_PENDING_KEY)
    return
  }
  setTimeout(restoreOnce, 120)
  setTimeout(() => {
    restoreOnce()
    sessionStorage.removeItem(HOME_FEED_RETURN_KEY)
    sessionStorage.removeItem(HOME_FEED_RETURN_PENDING_KEY)
  }, 360)
}

onActivated(restoreFeedReturnAnchor)

onMounted(() => {
  syncColumnCount()
  window.addEventListener('resize', syncColumnCount)
})

onUnmounted(() => {
  window.removeEventListener('resize', syncColumnCount)
})
</script>
