<template>
  <div class="relative bg-black" style="height:100vh;overflow:hidden">

    <!-- 滚动容器 -->
    <div ref="container" class="w-full h-full overflow-y-scroll snap-y snap-mandatory"
         style="-webkit-overflow-scrolling:touch">

      <div v-for="(v, idx) in videos" :key="v.video_id"
           :data-idx="idx"
           class="relative w-full snap-start snap-always shrink-0"
           :style="{ height: screenH + 'px' }">

        <!-- 封面图（始终可见，视频就绪后被覆盖；object-contain 与视频保持一致比例） -->
        <img :src="v.cover_url || '/placeholder.svg'"
             class="absolute inset-0 w-full h-full object-contain pointer-events-none"/>

        <!-- 视频元素（object-contain 保持原始比例，横屏视频不裁剪，黑边填充） -->
        <video
          :ref="el => { if (el) videoRefs[idx] = el }"
          class="absolute inset-0 w-full h-full object-contain"
          style="background:#000"
          loop playsinline preload="none"
          @click="togglePlay(idx)"
        />

        <!-- 渐变遮罩 -->
        <div class="absolute inset-0 pointer-events-none"
             style="background:linear-gradient(to top,rgba(0,0,0,0.72) 0%,transparent 45%,transparent 70%,rgba(0,0,0,0.25) 100%)"/>

        <!-- 暂停图标 -->
        <transition name="fade-icon">
          <div v-if="pausedIdx === idx"
               class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="w-16 h-16 rounded-full bg-black/40 flex items-center justify-center">
              <svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
              </svg>
            </div>
          </div>
        </transition>

        <!-- 右侧操作栏 -->
        <div class="absolute right-3 bottom-28 flex flex-col items-center gap-5">
          <button @click.stop="toggleFavorite(v)" class="flex flex-col items-center gap-1">
            <div :class="['w-11 h-11 rounded-full bg-black/30 backdrop-blur flex items-center justify-center transition-all',
              v._favorited ? 'text-primary' : 'text-white']">
              <svg class="w-6 h-6" :fill="v._favorited ? 'currentColor' : 'none'"
                   stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
              </svg>
            </div>
            <span class="text-white text-[10px] font-medium drop-shadow">{{ fmtN(v.favorite_count) }}</span>
          </button>

          <button @click.stop="openComment(v)" class="flex flex-col items-center gap-1">
            <div class="w-11 h-11 rounded-full bg-black/30 backdrop-blur flex items-center justify-center text-white">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
            </div>
            <span class="text-white text-[10px] font-medium drop-shadow">{{ fmtN(v.comment_count) }}</span>
          </button>

          <button @click.stop="share(v)" class="flex flex-col items-center gap-1">
            <div class="w-11 h-11 rounded-full bg-black/30 backdrop-blur flex items-center justify-center text-white">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
              </svg>
            </div>
            <span class="text-white text-[10px] font-medium drop-shadow">分享</span>
          </button>
        </div>

        <!-- 底部信息 -->
        <div class="absolute left-3 right-16 bottom-14">
          <p class="text-white font-semibold text-sm leading-snug line-clamp-2 drop-shadow">{{ v.title }}</p>
          <div class="flex items-center gap-1 text-white/60 text-xs mt-1.5">
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
            {{ fmtN(v.play_count) }} 次播放
          </div>
        </div>

        <!-- 进度条 -->
        <div class="absolute bottom-0 left-0 right-0 z-10 pb-2 pt-1" @click.stop @touchstart.stop.passive>
          <div class="flex items-center justify-between px-3 mb-1">
            <span class="text-white/60 text-[10px]">{{ fmtTime(videoTimes[idx]?.current) }}</span>
            <span class="text-white/60 text-[10px]">{{ fmtTime(videoTimes[idx]?.duration || v.duration) }}</span>
          </div>
          <div class="relative h-1 mx-3 bg-white/25 rounded-full">
            <div class="absolute inset-y-0 left-0 bg-white rounded-full pointer-events-none transition-none"
                 :style="{width: getProgress(idx)+'%'}"/>
            <input type="range" min="0" max="100" step="0.1"
                   :value="getProgress(idx)"
                   @input="seekTo(idx, Number($event.target.value))"
                   class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                   style="touch-action:none"/>
          </div>
        </div>
      </div>

      <!-- 底部加载指示 -->
      <div v-if="loading || hasMore" ref="sentinel"
           class="snap-start shrink-0 flex items-center justify-center"
           :style="{ height: screenH + 'px' }">
        <div v-if="loading" class="w-8 h-8 border-2 border-white/30 border-t-white rounded-full animate-spin"/>
      </div>
    </div>

    <!-- 左上角：返回 -->
    <div class="absolute top-4 left-3 z-10">
      <button @click="goBack"
              class="w-9 h-9 rounded-full bg-black/40 backdrop-blur
                     flex items-center justify-center text-white active:scale-95 transition-transform">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>
    </div>

    <!-- 评论抽屉 -->
    <transition name="slide-up">
      <div v-if="commentVideo"
           class="absolute inset-0 z-20 flex flex-col justify-end bg-black/50"
           @click.self="commentVideo = null">
        <div class="bg-bg-surface rounded-t-2xl border-t border-border max-h-[70vh] flex flex-col">
          <div class="flex items-center justify-between px-4 py-3 border-b border-border shrink-0">
            <span class="text-sm font-semibold text-text-primary">评论</span>
            <button @click="commentVideo = null" class="text-text-muted hover:text-text-primary">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto">
            <CommentSection :video-id="commentVideo.video_id"/>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Hls from 'hls.js'
import { videoApi, favoriteApi, playApi } from '@/api'
import { useUserStore } from '@/stores/user'
import CommentSection from '@/components/common/CommentSection.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const container = ref(null)
const videos = ref([])
const videoRefs = reactive([])   // 用 reactive 避免 ref[idx] 赋值到包装对象上
const pausedIdx = ref(-1)
const currentPlayingIdx = ref(-1)
const commentVideo = ref(null)
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
const screenH = ref(window.innerHeight)
const startIdRaw = route.query.id == null ? '' : String(route.query.id)
const startId = Number(startIdRaw) || 0
const sentinel = ref(null)

const orderByAllowed = ['created_at', 'play_count', 'comment_count', 'favorite_count']
const orderByRaw = route.query.order_by == null ? '' : String(route.query.order_by)
const orderBy = orderByAllowed.includes(orderByRaw) ? orderByRaw : 'created_at'
// idx -> { current: number, duration: number }
const videoTimes = reactive({})

// idx -> Hls 实例
const hlsMap = {}
// 已完成初始化的 idx
const hlsReady = new Set()
// 待播放的 idx（等 MANIFEST_PARSED 后播放）
let pendingPlay = -1

let sentinelObserver = null

function observeSentinel() {
  sentinelObserver?.disconnect()
  sentinelObserver = null
  if (!sentinel.value || !container.value) return
  sentinelObserver = new IntersectionObserver(([entry]) => {
    if (entry.isIntersecting) fetchMore()
  }, { root: container.value, rootMargin: '200px', threshold: 0.01 })
  sentinelObserver.observe(sentinel.value)
}

function fmtN(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function stopAll() {
  videoRefs.forEach(el => { if (el) el.pause() })
}

function initVideo(idx) {
  if (hlsReady.has(idx) || hlsMap[idx]) return
  const v = videos.value[idx]
  const src = v?.play_url
  if (!src) return
  const el = videoRefs[idx]
  if (!el) return

  el.addEventListener('timeupdate', () => {
    videoTimes[idx] = { current: el.currentTime, duration: el.duration || 0 }
  })

  if (src.includes('.m3u8') && Hls.isSupported()) {
    // Chrome / Firefox / Edge：用 hls.js
    const hls = new Hls({ enableWorker: true, maxBufferLength: 20 })
    hls.loadSource(src)
    hls.attachMedia(el)
    hlsMap[idx] = hls

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      hlsReady.add(idx)
      if (pendingPlay === idx) {
        el.muted = false
        el.play().catch(() => { el.muted = true; el.play().catch(() => {}) })
      }
    })

    hls.on(Hls.Events.ERROR, (_, data) => {
      if (data.fatal) hls.destroy()
    })
    return
  }

  // Safari 原生 HLS 或 MP4
  el.src = src
  hlsReady.add(idx)
}

function playAt(idx) {
  const el = videoRefs[idx]
  if (!el) return

  // 防止 IntersectionObserver 在列表追加/重新观察时重复触发同一条视频，
  // 导致 stopAll() 后重头播放。
  if (idx === pausedIdx.value) return
  if (idx === currentPlayingIdx.value && !el.paused) return

  stopAll()
  pausedIdx.value = -1
  currentPlayingIdx.value = idx
  pendingPlay = idx
  initVideo(idx)

  const v = videos.value[idx]
  // 观看视频不应要求登录；仅记录播放历史需要登录，避免 401 导致重定向。
  if (v && userStore.isLoggedIn) playApi.record(v.video_id).catch(() => {})

  if (hlsReady.has(idx)) {
    el.muted = false
    el.play().catch(() => { el.muted = true; el.play().catch(() => {}) })
  }
  // 否则等 MANIFEST_PARSED 触发
}

function togglePlay(idx) {
  const el = videoRefs[idx]
  if (!el) return
  if (el.paused) { el.play().catch(() => {}); pausedIdx.value = -1 }
  else { el.pause(); pausedIdx.value = idx }
  if (!el.paused) currentPlayingIdx.value = idx
}


function getProgress(idx) {
  const t = videoTimes[idx]
  if (!t || !t.duration) return 0
  return Math.min((t.current / t.duration) * 100, 100)
}

function seekTo(idx, pct) {
  const el = videoRefs[idx]
  if (!el || !el.duration) return
  el.currentTime = (pct / 100) * el.duration
}

function fmtTime(s) {
  if (!s || isNaN(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = String(Math.floor(s % 60)).padStart(2, '0')
  return `${m}:${sec}`
}

let itemObservers = []
function observeItems() {
  itemObservers.forEach(o => o.disconnect())
  itemObservers = []
  const els = container.value?.querySelectorAll('[data-idx]')
  els?.forEach(el => {
    const idx = Number(el.dataset.idx)
    const io = new IntersectionObserver(([entry]) => {
      if (!entry.isIntersecting || entry.intersectionRatio < 0.6) return

      // 滑动过程中可能出现相邻两条同时满足交叉阈值，
      // 用 scrollTop 推断“当前应该播放的 idx”，避免回切/重头播放。
      if (!container.value) return
      const expectedRaw = (container.value.scrollTop + screenH.value * 0.5) / screenH.value
      const expectedIdx = Math.floor(expectedRaw)
      const maxIdx = Math.max(0, videos.value.length - 1)
      const safeExpectedIdx = Math.min(Math.max(0, expectedIdx), maxIdx)
      if (idx !== safeExpectedIdx) return

      playAt(idx)
    }, { threshold: 0.6, root: container.value })
    io.observe(el)
    itemObservers.push(io)
  })
}

async function fetchMore() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  try {
    const pageSize = 20
    const res = await videoApi.list({ page: page.value, page_size: pageSize, order_by: orderBy })
    const rawList = res.data?.videos || []
    const list = rawList
      .filter(v => !videos.value.some(e => e.video_id === v.video_id))
      .map(v => ({ ...v, _favorited: false }))
    videos.value.push(...list)
    page.value++
    // 注意：这里用 rawList 判断是否到底，而不是用 filter 后长度。
    // 否则当目标视频重复被过滤时，会误判“已到底”，导致无法继续滑动。
    if (rawList.length < pageSize) hasMore.value = false
    await nextTick()
    observeItems()
    observeSentinel()
  } catch {}
  loading.value = false
}

const onResize = () => { screenH.value = window.innerHeight }

function goBack() {
  stopAll()
  router.back()
}

function onContainerScroll() {
  const el = container.value
  if (!el) return
  if (el.scrollHeight - el.scrollTop - el.clientHeight < screenH.value * 1.5) fetchMore()
}

onMounted(async () => {
  const pageSize = 10

  if (startIdRaw) {
    const maxSearchPages = 50
    let foundIdx = -1
    let currentPage = 1
    let lastRawLen = 0
    videos.value = []
    hlsReady.clear()
    Object.keys(hlsMap).forEach(k => delete hlsMap[k])

    while (currentPage <= maxSearchPages && foundIdx < 0) {
      const res = await videoApi.list({ page: currentPage, page_size: pageSize, order_by: orderBy })
      const rawList = res.data?.videos || []
      lastRawLen = rawList.length

      const newOnes = rawList
        .filter(v => !videos.value.some(e => e.video_id === v.video_id))
        .map(v => ({ ...v, _favorited: false }))

      videos.value.push(...newOnes)

      foundIdx = videos.value.findIndex(v => String(v.video_id) === startIdRaw)
      if (rawList.length < pageSize) break
      currentPage++
    }

    // 下一次 fetchMore 从哪一页继续
    page.value = currentPage + 1
    hasMore.value = lastRawLen >= pageSize

    if (foundIdx < 0) {
      return
    }

    await nextTick()
    if (container.value) container.value.scrollTop = foundIdx * screenH.value
    await nextTick()
    observeItems()
    playAt(foundIdx)
  } else {
    // 不带目标 ID：正常加载第一页并从 0 开始播放
    await fetchMore()
    await nextTick()
    observeItems()
    playAt(0)
  }

  // 用 sentinel 触底加载，避免 snap/阈值导致“第9条后不再触发加载”
  observeSentinel()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  stopAll()
  Object.values(hlsMap).forEach(h => h.destroy())
  itemObservers.forEach(o => o.disconnect())
  sentinelObserver?.disconnect()
  window.removeEventListener('resize', onResize)
})

async function toggleFavorite(v) {
  if (!userStore.isLoggedIn) { userStore.openAuth('login'); return }
  try {
    if (v._favorited) {
      await favoriteApi.remove(v.video_id)
      v._favorited = false
      v.favorite_count = Math.max(0, (v.favorite_count || 0) - 1)
    } else {
      await favoriteApi.add(v.video_id)
      v._favorited = true
      v.favorite_count = (v.favorite_count || 0) + 1
    }
  } catch {}
}

function openComment(v) { commentVideo.value = v }

async function share(v) {
  const url = `${location.origin}/feed?id=${v.video_id}&order_by=${orderBy}`
  if (navigator.share) {
    navigator.share({ title: v.title, url }).catch(() => {})
  } else {
    await navigator.clipboard.writeText(url).catch(() => {})
    alert('链接已复制')
  }
}
</script>

<style scoped>
.fade-icon-enter-active { transition: opacity 0.2s; }
.fade-icon-leave-active { transition: opacity 0.4s; }
.fade-icon-enter-from, .fade-icon-leave-to { opacity: 0; }

.slide-up-enter-active { transition: transform 0.3s ease; }
.slide-up-leave-active { transition: transform 0.25s ease; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); }
</style>
