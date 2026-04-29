<template>
  <div class="relative bg-black" style="height:100vh;overflow:hidden">

    <!-- 滚动容器 -->
    <div ref="container" class="w-full h-full overflow-y-scroll snap-y snap-mandatory"
         style="-webkit-overflow-scrolling:touch">

      <div v-for="(v, idx) in videos" :key="v.video_id"
           :data-idx="idx"
           class="relative w-full snap-start snap-always shrink-0"
           :style="{ height: screenH + 'px' }">

        <!-- 封面图（始终可见，视频就绪后被覆盖） -->
        <img :src="v.cover_url || '/placeholder.svg'"
             class="absolute inset-0 w-full h-full object-cover pointer-events-none"/>

        <!-- 视频元素（背景透明，加载失败时露出封面） -->
        <video
          :ref="el => { if (el) videoRefs[idx] = el }"
          class="absolute inset-0 w-full h-full object-cover"
          style="background:transparent"
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
      <div v-if="loading || hasMore" class="snap-start shrink-0 flex items-center justify-center"
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
const commentVideo = ref(null)
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
const screenH = ref(window.innerHeight)
const startIdRaw = route.query.id == null ? '' : String(route.query.id)
const startId = Number(startIdRaw) || 0
// idx -> { current: number, duration: number }
const videoTimes = reactive({})

// idx -> Hls 实例
const hlsMap = {}
// 已完成初始化的 idx
const hlsReady = new Set()
// 待播放的 idx（等 MANIFEST_PARSED 后播放）
let pendingPlay = -1

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
  stopAll()
  pausedIdx.value = -1
  pendingPlay = idx
  initVideo(idx)

  const v = videos.value[idx]
  if (v) playApi.record(v.video_id).catch(() => {})

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
    const io = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && entry.intersectionRatio >= 0.6) {
        playAt(Number(el.dataset.idx))
      }
    }, { threshold: 0.6, root: container.value })
    io.observe(el)
    itemObservers.push(io)
  })
}

async function fetchMore() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  try {
    const res = await videoApi.list({ page: page.value, page_size: 10, order_by: 'created_at' })
    const list = (res.data?.videos || [])
      .filter(v => !videos.value.some(e => e.video_id === v.video_id))
      .map(v => ({ ...v, _favorited: false }))
    videos.value.push(...list)
    page.value++
    if (list.length < 10) hasMore.value = false
    await nextTick()
    observeItems()
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
  // ── 步骤 1：若携带目标视频 ID，先把它加载到列表第 0 位 ──────────────
  // 必须先于 fetchMore()，否则 fetchMore() 会为第一页的随机视频初始化
  // hlsMap[0]，导致 unshift 之后 index 0 的 HLS 状态错位、播错视频。
  if (startId) {
    try {
      const res = await videoApi.detail(startId)
      if (res?.data) {
        // 目标视频放到数组首位，_favorited 默认未收藏
        videos.value = [{ ...res.data, _favorited: false }]
      }
    } catch {
      // 拉取失败时降级：直接走下面的 fetchMore() 播第一页第一个视频
    }
  }

  // ── 步骤 2：追加更多视频（fetchMore 内部有去重过滤，目标视频不会重复） ──
  await fetchMore()
  await nextTick()

  // 首屏由我们显式控制目标播放，避免初始 IO 回调抢占导致播错视频。
  const targetIdx = startIdRaw
    ? Math.max(0, videos.value.findIndex(v => String(v.video_id) === startIdRaw))
    : 0
  if (container.value) {
    container.value.scrollTop = targetIdx * screenH.value
  }
  playAt(targetIdx)

  container.value?.addEventListener('scroll', onContainerScroll, { passive: true })
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  stopAll()
  Object.values(hlsMap).forEach(h => h.destroy())
  itemObservers.forEach(o => o.disconnect())
  container.value?.removeEventListener('scroll', onContainerScroll)
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
  const url = `${location.origin}/feed?id=${v.video_id}`
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
