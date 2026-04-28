<template>
  <!-- 单根元素，height:100vh 不用 fixed，避免与 Vue transition 的 transitionend 冲突 -->
  <div class="relative bg-black" style="height:100vh;overflow:hidden">

    <!-- 滚动容器 -->
    <div ref="container" class="w-full h-full overflow-y-scroll snap-y snap-mandatory"
         style="-webkit-overflow-scrolling:touch">

      <div v-for="(v, idx) in videos" :key="v.video_id"
           :data-idx="idx"
           class="relative w-full snap-start snap-always shrink-0"
           :style="{ height: screenH + 'px' }">

        <!-- 封面图兜底：视频加载前/失败时不显示黑屏 -->
        <img v-if="v.cover_url" :src="v.cover_url"
             class="absolute inset-0 w-full h-full object-cover opacity-60 pointer-events-none"/>
        <video
          :ref="el => { if (el) videoRefs[idx] = el }"
          class="absolute inset-0 w-full h-full object-cover"
          :src="v.play_url || v.source_url"
          loop playsinline preload="metadata"
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
        <div class="absolute left-3 right-16 bottom-8">
          <p class="text-white font-semibold text-sm leading-snug line-clamp-2 drop-shadow mb-2">{{ v.title }}</p>
          <div class="flex flex-wrap gap-1.5 mb-2">
            <span v-for="t in tagsOf(v)" :key="t"
                  class="text-[11px] text-white/80 bg-white/10 backdrop-blur rounded-full px-2 py-0.5">
              #{{ t }}
            </span>
          </div>
          <div class="flex items-center gap-1 text-white/60 text-xs">
            <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
            {{ fmtN(v.play_count) }} 次播放
          </div>
        </div>
      </div>

      <!-- 加载哨兵 -->
      <div ref="sentinel" class="snap-start shrink-0 flex items-center justify-center"
           :style="{ height: screenH + 'px' }">
        <div v-if="loading" class="w-8 h-8 border-2 border-white/30 border-t-white rounded-full animate-spin"/>
      </div>
    </div>

    <!-- 左上角返回按钮 -->
    <div class="absolute top-4 left-3 z-10">
      <button @click="goBack"
              class="w-9 h-9 rounded-full bg-black/40 backdrop-blur
                     flex items-center justify-center text-white active:scale-95 transition-transform">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>
    </div>

    <!-- 评论抽屉（absolute 在 fixed 容器内，无需 teleport） -->
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
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { videoApi, favoriteApi, playApi } from '@/api'
import { useUserStore } from '@/stores/user'
import CommentSection from '@/components/common/CommentSection.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const container = ref(null)
const sentinel = ref(null)
const videos = ref([])
const videoRefs = ref([])
const pausedIdx = ref(-1)
const commentVideo = ref(null)
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
const screenH = ref(window.innerHeight)
const startId = Number(route.query.id) || 0

function fmtN(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function tagsOf(v) {
  return (v.tags || '').split(',').map(t => t.trim()).filter(Boolean).slice(0, 3)
}

function stopAll() {
  videoRefs.value.forEach(v => { if (v) { v.pause(); v.currentTime = 0 } })
}

function playAt(idx) {
  const el = videoRefs.value[idx]
  if (!el) return
  stopAll()
  pausedIdx.value = -1
  el.play().catch(() => {})
  const v = videos.value[idx]
  if (v) playApi.record(v.video_id).catch(() => {})
}

function togglePlay(idx) {
  const el = videoRefs.value[idx]
  if (!el) return
  if (el.paused) { el.play().catch(() => {}); pausedIdx.value = -1 }
  else { el.pause(); pausedIdx.value = idx }
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

let sentinelIo = null
const onResize = () => { screenH.value = window.innerHeight }

function goBack() {
  stopAll()
  router.back()
}

onMounted(async () => {
  if (startId) {
    try {
      const res = await videoApi.detail(startId)
      if (res.data) videos.value.push({ ...res.data, _favorited: false })
    } catch {}
  }
  await fetchMore()
  await nextTick()
  playAt(0)

  sentinelIo = new IntersectionObserver(([entry]) => {
    if (entry.isIntersecting) fetchMore()
  }, { rootMargin: '600px', root: container.value })
  if (sentinel.value) sentinelIo.observe(sentinel.value)

  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  stopAll()
  itemObservers.forEach(o => o.disconnect())
  sentinelIo?.disconnect()
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
  const url = `${location.origin}/video/${v.video_id}`
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
