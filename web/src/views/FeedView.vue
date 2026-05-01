<template>
  <div ref="feedRoot" class="relative bg-black" style="height:100vh;overflow:hidden">

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
            playsinline preload="none"
            @click="togglePlay(idx)"
            @play="syncPlayState(idx, true)"
            @pause="syncPlayState(idx, false)"
            @ended="handleEnded(idx)"
            @loadedmetadata="syncVideoTime(idx)"
        />

        <!-- 渐变遮罩 -->
        <div v-if="!cleanMode" class="absolute inset-0 pointer-events-none"
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
        <div v-if="!cleanMode" class="absolute right-3 bottom-36 flex flex-col items-center gap-5">
          <button @click.stop="toggleFavorite(v)" class="flex flex-col items-center gap-1">
            <div :class="['w-11 h-11 rounded-full bg-black/30 backdrop-blur flex items-center justify-center transition-all',
              v._favorited ? 'text-yellow-400' : 'text-white']">
              <svg class="w-6 h-6" :fill="v._favorited ? 'currentColor' : 'none'"
                   stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z"/>
              </svg>
            </div>
            <span class="text-white text-[10px] font-medium drop-shadow">{{ fmtN(v.favorite_count) }}</span>
          </button>

          <!-- 评论按钮 -->
          <button @click.stop="openComment(v)" class="flex flex-col items-center gap-1">
            <div class="w-11 h-11 rounded-full bg-black/30 backdrop-blur flex items-center justify-center text-white">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
            </div>
            <span class="text-white text-[10px] font-medium drop-shadow">{{ fmtN(v.comment_count) }}</span>
          </button>

          <!-- 分享按钮 -->
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
        <div v-if="!cleanMode" class="absolute left-3 right-16 bottom-24">
          <p class="text-white font-semibold text-sm leading-snug line-clamp-2 drop-shadow">{{ v.title }}</p>
          <div class="flex items-center gap-1 text-white/60 text-xs mt-1.5">
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z"/>
            </svg>
            {{ fmtN(v.play_count) }} 次播放
          </div>
        </div>

        <!-- 底部控制栏 -->
        <div class="absolute bottom-0 left-0 right-0 z-10 bg-black/90 text-white" @click.stop @touchstart.stop.passive>
          <div class="relative h-2 bg-white/20">
            <div class="absolute inset-y-0 left-0 bg-white rounded-full pointer-events-none transition-none"
                 :style="{width: getProgress(idx)+'%'}"/>
            <input type="range" min="0" max="100" step="0.1"
                   :value="getProgress(idx)"
                   @input="seekTo(idx, Number($event.target.value))"
                   class="feed-progress-range absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                   style="touch-action:none"/>
          </div>
          <div class="flex min-h-12 items-center gap-3 px-3">
            <button class="control-icon-btn shrink-0" @click="togglePlay(idx)" :aria-label="isPlaying(idx) ? '暂停' : '播放'">
              <svg v-if="isPlaying(idx)" class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z"/>
              </svg>
            </button>
            <div class="text-xs tabular-nums text-white/90 shrink-0">
              {{ fmtTime(videoTimes[idx]?.current) }} / {{ fmtTime(videoTimes[idx]?.duration || v.duration) }}
            </div>
            <div class="flex-1"></div>
            <div class="flex items-center gap-2 overflow-x-auto controls-scroll">
              <button :class="['control-toggle', continuousPlay ? 'is-on' : '']" @click="toggleContinuousPlay">
                <span class="toggle-dot"></span>
                <span>连播</span>
              </button>
              <button :class="['control-toggle', cleanMode ? 'is-on' : '']" @click="cleanMode = !cleanMode">
                <span class="toggle-dot"></span>
                <span>清屏</span>
              </button>
              <button class="control-text-btn" @click="cycleSpeed">{{ playbackRateLabel }}</button>
              <button class="control-icon-btn" @click="toggleMuted" :aria-label="muted ? '打开声音' : '静音'">
                <svg v-if="muted" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16 9l5 5m0-5l-5 5M5 9v6h4l5 4V5L9 9H5z"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 9v6h4l5 4V5L9 9H5z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" d="M17 9.5a4 4 0 010 5M19.5 7a7 7 0 010 10"/>
                </svg>
              </button>
              <button class="control-icon-btn" @click="toggleFullscreen" :aria-label="isFullscreen ? '退出全屏' : '全屏'">
                <svg v-if="isFullscreen" class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v6H3M15 3v6h6M9 21v-6H3M15 21v-6h6"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8 3H3v5M16 3h5v5M8 21H3v-5M16 21h5v-5"/>
                </svg>
              </button>
            </div>
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
    <div v-if="!cleanMode" class="absolute top-4 left-3 z-10">
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
            <button @click="commentVideo = null" class="text-sm font-semibold text-text-primary">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto">
            <CommentSection
                :video-id="commentVideo.video_id"
                @comment-count-updated="updateCommentCount"
            />
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import {computed, nextTick, onMounted, onUnmounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import Hls from 'hls.js'
import { videoApi, favoriteApi, playApi } from '@/api'
import { useUserStore } from '@/stores/user'
import CommentSection from '@/components/common/CommentSection.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const feedRoot = ref(null)
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
const cleanMode = ref(false)
const continuousPlay = ref(true)
const muted = ref(false)
const isFullscreen = ref(false)
const playbackRates = [1, 1.25, 1.5, 2]
const playbackRateIdx = ref(0)
const playStates = reactive({})
const playbackRate = computed(() => playbackRates[playbackRateIdx.value])
const playbackRateLabel = computed(() => `${playbackRate.value === 1 ? '倍速' : playbackRate.value + 'x'}`)

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
  }, {root: container.value, rootMargin: '200px', threshold: 0.01})
  sentinelObserver.observe(sentinel.value)
}

function fmtN(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function stopAll() {
  videoRefs.forEach(el => {
    if (el) el.pause()
  })
}

function applyVideoSettings(el) {
  if (!el) return
  el.muted = muted.value
  el.playbackRate = playbackRate.value
}

function syncVideoTime(idx) {
  const el = videoRefs[idx]
  if (!el) return
  videoTimes[idx] = {current: el.currentTime || 0, duration: el.duration || 0}
  applyVideoSettings(el)
}

function syncPlayState(idx, playing) {
  playStates[idx] = playing
  if (playing) {
    pausedIdx.value = -1
    currentPlayingIdx.value = idx
  } else if (currentPlayingIdx.value === idx) {
    playStates[idx] = false
  }
}

function initVideo(idx) {
  if (hlsReady.has(idx) || hlsMap[idx]) return
  const v = videos.value[idx]
  const src = v?.play_url
  if (!src) return
  const el = videoRefs[idx]
  if (!el) return
  applyVideoSettings(el)

  el.addEventListener('timeupdate', () => {
    videoTimes[idx] = {current: el.currentTime, duration: el.duration || 0}
  })

  if (src.includes('.m3u8') && Hls.isSupported()) {
    // Chrome / Firefox / Edge：用 hls.js
    const hls = new Hls({enableWorker: true, maxBufferLength: 20})
    hls.loadSource(src)
    hls.attachMedia(el)
    hlsMap[idx] = hls

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      hlsReady.add(idx)
      if (pendingPlay === idx) {
        applyVideoSettings(el)
        el.play().catch(() => {
          muted.value = true
          applyVideoSettings(el)
          el.play().catch(() => {
          })
        })
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
  applyVideoSettings(el)
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
  if (v && userStore.isLoggedIn) playApi.record(v.video_id).catch(() => {
  })

  if (hlsReady.has(idx)) {
    applyVideoSettings(el)
    el.play().catch(() => {
      muted.value = true
      applyVideoSettings(el)
      el.play().catch(() => {
      })
    })
  }
  // 否则等 MANIFEST_PARSED 触发
}

function togglePlay(idx) {
  const el = videoRefs[idx]
  if (!el) return
  if (el.paused) {
    applyVideoSettings(el)
    el.play().catch(() => {
    });
    pausedIdx.value = -1
    playStates[idx] = true
  } else {
    el.pause();
    pausedIdx.value = idx
    playStates[idx] = false
  }
  if (!el.paused) currentPlayingIdx.value = idx
}

function isPlaying(idx) {
  return !!playStates[idx]
}

async function handleEnded(idx) {
  playStates[idx] = false
  if (!continuousPlay.value) {
    pausedIdx.value = idx
    return
  }

  const nextIdx = idx + 1
  if (nextIdx >= videos.value.length) {
    await fetchMore()
  }
  if (nextIdx >= videos.value.length) {
    pausedIdx.value = idx
    return
  }
  pausedIdx.value = -1
  if (container.value) {
    container.value.scrollTo({top: nextIdx * screenH.value, behavior: 'smooth'})
  }
  window.setTimeout(() => playAt(nextIdx), 220)
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
  syncVideoTime(idx)
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
    }, {threshold: 0.6, root: container.value})
    io.observe(el)
    itemObservers.push(io)
  })
}

async function checkFavorites(list) {
  if (!userStore.isLoggedIn || !list.length) return
  await Promise.allSettled(
      list.map(v =>
          favoriteApi.check(v.video_id)
              .then(res => {
                v._favorited = !!res.data?.is_favorited
              })
              .catch(() => {
              })
      )
  )
}

async function fetchMore() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  try {
    const pageSize = 20
    const res = await videoApi.list({page: page.value, page_size: pageSize, order_by: orderBy})
    const rawList = res.data?.videos || []
    const list = rawList
        .filter(v => !videos.value.some(e => e.video_id === v.video_id))
        .map(v => ({...v, _favorited: false}))
    videos.value.push(...list)
    page.value++
    // 注意：这里用 rawList 判断是否到底，而不是用 filter 后长度。
    // 否则当目标视频重复被过滤时，会误判”已到底”，导致无法继续滑动。
    if (rawList.length < pageSize) hasMore.value = false
    await checkFavorites(list)
    await nextTick()
    observeItems()
    observeSentinel()
  } catch {
  }
  loading.value = false
}

const onResize = () => {
  screenH.value = window.innerHeight
}

function toggleContinuousPlay() {
  continuousPlay.value = !continuousPlay.value
}

function toggleMuted() {
  muted.value = !muted.value
  videoRefs.forEach(applyVideoSettings)
}

function cycleSpeed() {
  playbackRateIdx.value = (playbackRateIdx.value + 1) % playbackRates.length
  videoRefs.forEach(applyVideoSettings)
}

function syncFullscreenState() {
  isFullscreen.value = !!(document.fullscreenElement || document.webkitFullscreenElement)
}

async function toggleFullscreen() {
  const target = feedRoot.value || document.documentElement
  if (!isFullscreen.value) {
    const request = target.requestFullscreen || target.webkitRequestFullscreen
    if (request) await request.call(target).catch(() => {})
  } else {
    const exit = document.exitFullscreen || document.webkitExitFullscreen
    if (exit) await exit.call(document).catch(() => {})
  }
  syncFullscreenState()
}

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
      const res = await videoApi.list({page: currentPage, page_size: pageSize, order_by: orderBy})
      const rawList = res.data?.videos || []
      lastRawLen = rawList.length

      const newOnes = rawList
          .filter(v => !videos.value.some(e => e.video_id === v.video_id))
          .map(v => ({...v, _favorited: false}))

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

    await checkFavorites(videos.value)
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
  document.addEventListener('fullscreenchange', syncFullscreenState)
  document.addEventListener('webkitfullscreenchange', syncFullscreenState)
})

onUnmounted(() => {
  stopAll()
  Object.values(hlsMap).forEach(h => h.destroy())
  itemObservers.forEach(o => o.disconnect())
  sentinelObserver?.disconnect()
  window.removeEventListener('resize', onResize)
  document.removeEventListener('fullscreenchange', syncFullscreenState)
  document.removeEventListener('webkitfullscreenchange', syncFullscreenState)
})

// 收藏
async function toggleFavorite(v) {
  if (!userStore.isLoggedIn) {
    userStore.openAuth('login')
    return
  }

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
  } catch (err) {
    // 可以在这里提示用户
  }
}

function openComment(v) { commentVideo.value = v }

function updateCommentCount(count) {
  if (!commentVideo.value) return
  commentVideo.value.comment_count = count
  const item = videos.value.find(v => v.video_id === commentVideo.value.video_id)
  if (item) item.comment_count = count
}

async function share(v) {
  const url = `${location.origin}/feed?id=${v.video_id}&order_by=${orderBy}`
  if (navigator.share) {
    navigator.share({title: v.title, url}).catch(() => {
    })
  } else {
    await navigator.clipboard.writeText(url).catch(() => {
    })
    alert('链接已复制')
  }
}
</script>

<style scoped>
.fade-icon-enter-active {
  transition: opacity 0.2s;
}

.fade-icon-leave-active {
  transition: opacity 0.4s;
}

.fade-icon-enter-from, .fade-icon-leave-to {
  opacity: 0;
}

.slide-up-enter-active {
  transition: transform 0.3s ease;
}

.slide-up-leave-active {
  transition: transform 0.25s ease;
}

.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}

.control-icon-btn {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.92);
  display: inline-flex;
  height: 32px;
  justify-content: center;
  min-width: 32px;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.control-icon-btn:active,
.control-icon-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.control-text-btn {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.92);
  display: inline-flex;
  font-size: 13px;
  font-weight: 600;
  height: 32px;
  justify-content: center;
  min-width: 44px;
  padding: 0 8px;
  white-space: nowrap;
}

.control-text-btn:active,
.control-text-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.control-toggle {
  align-items: center;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.86);
  display: inline-flex;
  font-size: 13px;
  font-weight: 600;
  gap: 5px;
  height: 32px;
  padding: 0 6px;
  white-space: nowrap;
}

.control-toggle:active,
.control-toggle:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.toggle-dot {
  background: rgba(255, 255, 255, 0.66);
  border-radius: 999px;
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.16);
  height: 8px;
  width: 8px;
}

.control-toggle.is-on .toggle-dot {
  background: #ff2b75;
  box-shadow: 0 0 0 3px rgba(255, 43, 117, 0.24);
}

.feed-progress-range {
  appearance: none;
}

.feed-progress-range::-webkit-slider-thumb {
  appearance: none;
  height: 18px;
  width: 18px;
}

.feed-progress-range::-moz-range-thumb {
  border: 0;
  height: 18px;
  width: 18px;
}

.controls-scroll {
  scrollbar-width: none;
}

.controls-scroll::-webkit-scrollbar {
  display: none;
}
</style>
