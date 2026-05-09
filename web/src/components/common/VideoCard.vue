<template>
  <div
    ref="cardRoot"
    :data-feed-video-id="video.video_id"
    class="group block rounded-xl overflow-hidden bg-bg-card border border-border
           hover:border-primary/30 hover:shadow-card-hover transition-all duration-300
           hover:-translate-y-0.5 cursor-pointer select-none"
    @click="handleClick"
    @touchstart.passive="onTouchStart"
    @touchend="onTouchEnd"
    @touchmove.passive="cancelLongPress"
    @touchcancel.passive="cancelLongPress">

    <!-- 封面区域：图片自然比例，不强制宽高比 -->
    <div class="relative overflow-hidden bg-bg-surface">

      <img :src="video.cover_url" :alt="video.title"
           class="w-full block transition-transform duration-500 group-hover:scale-105"
           loading="lazy"
           @error="e => { e.target.style.minHeight='120px'; e.target.src=''; }"/>

      <!-- 底部渐变遮罩 -->
      <div class="absolute inset-0 pointer-events-none"
           style="background:linear-gradient(to top,rgba(0,0,0,0.6) 0%,transparent 50%)"/>

      <!-- 左下：播放次数 -->
      <div class="absolute bottom-1.5 left-1.5 flex items-center gap-0.5 text-white/90 text-[10px]">
        <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
        </svg>
        {{ fmtCount(video.play_count) }}
      </div>

      <!-- 右下：时长 -->
      <div v-if="video.duration"
           class="absolute bottom-1.5 right-1.5 px-1.5 py-0.5 rounded bg-black/75 text-white text-[10px] font-mono">
        {{ fmtDuration(video.duration) }}
      </div>

      <!-- 长按预览视频 -->
      <Transition name="preview-fade">
        <video v-if="previewing" ref="previewEl"
               class="absolute inset-0 w-full h-full object-cover"
               playsinline loop preload="none"
               style="background:transparent"/>
      </Transition>

      <!-- 长按预览时：右上角静音切换 -->
      <!-- @touchend.stop.prevent：阻止 touchend 冒泡到卡片（否则会触发 cancelLongPress 关掉预览）
           .prevent 同时阻止浏览器把 touchend 转化为 click，避免 togglePreviewMute 被调用两次 -->
      <Transition name="preview-fade">
        <button v-if="previewing"
                @touchstart.stop
                @touchend.stop.prevent="togglePreviewMute"
                @click.stop
                class="absolute top-2 right-2 z-10 w-7 h-7 rounded-full bg-black/60
                       flex items-center justify-center text-white">
          <svg v-if="previewMuted" class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"/>
          </svg>
          <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M15.536 8.464a5 5 0 010 7.072M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
          </svg>
        </button>
      </Transition>
    </div>

    <!-- 标题 + 作者 + 时间 -->
    <div class="p-2.5 space-y-2">
      <h3 class="text-sm font-medium text-text-primary line-clamp-2 leading-snug
                 group-hover:text-primary transition-colors">
        {{ video.title }}
      </h3>
      <div class="flex items-center gap-2">
        <!-- 头像：有则显示用户头像，无则默认 logo -->
        <img :src="video.author_avatar || '/logo.png'" :alt="video.author_name"
             class="w-6 h-6 rounded-full shrink-0 object-cover"
             @error="e => e.target.src='/logo.png'"/>
        <!-- 用户名 -->
        <span class="text-text-primary/60 text-xs flex-1 truncate leading-none">
          {{ video.author_name || '投稿者' }}
        </span>
        <!-- 发布时间 -->
        <span class="text-text-primary/60 text-xs shrink-0 leading-none">
          {{ fmtDate(video.published_at) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Hls from 'hls.js'

const props = defineProps({
  video: { type: Object, required: true },
  // 用于把当前列表排序传到播放页，保证“点哪个视频进去就按哪个顺序播放+继续加载”
  orderBy: { type: String, default: 'created_at' },
})
const router = useRouter()
const route = useRoute()

const previewing = ref(false)
const cardRoot = ref(null)
const previewEl = ref(null)
const previewMuted = ref(true)

let longPressTimer = null
let longPressed = false
let previewHls = null
const HOME_FEED_RETURN_KEY = 'home_feed_return_anchor'

function handleClick() {
  if (longPressed) { longPressed = false; return }
  if (route.path === '/' && cardRoot.value) {
    const rect = cardRoot.value.getBoundingClientRect()
    sessionStorage.setItem(HOME_FEED_RETURN_KEY, JSON.stringify({
      videoId: String(props.video.video_id),
      top: rect.top,
      ts: Date.now(),
    }))
  }
  router.push({ path: '/feed', query: { id: props.video.video_id, order_by: props.orderBy } })
}

function onTouchStart() {
  longPressed = false
  longPressTimer = setTimeout(async () => {
    longPressed = true
    previewing.value = true
    await nextTick()
    const el = previewEl.value
    if (!el) return
    const src = props.video.play_url
    if (!src) return
    el.muted = previewMuted.value
    if (src.includes('.m3u8') && Hls.isSupported()) {
      previewHls = new Hls({ maxBufferLength: 8 })
      previewHls.loadSource(src)
      previewHls.attachMedia(el)
      previewHls.on(Hls.Events.MANIFEST_PARSED, () => el.play().catch(() => {}))
    } else {
      el.src = src
      el.play().catch(() => {})
    }
  }, 550)
}

function onTouchEnd(e) {
  cancelLongPress()
  if (longPressed) e.preventDefault()
}

function cancelLongPress() {
  clearTimeout(longPressTimer)
  if (previewing.value) {
    previewing.value = false
    previewMuted.value = true
    if (previewHls) { previewHls.destroy(); previewHls = null }
    if (previewEl.value) { previewEl.value.pause(); previewEl.value.src = '' }
  }
}

function togglePreviewMute() {
  previewMuted.value = !previewMuted.value
  if (previewEl.value) previewEl.value.muted = previewMuted.value
}

// ── 格式化工具 ────────────────────────────────────────

function fmtDuration(s) {
  if (!s) return ''
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

function fmtCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function fmtDate(t) {
  if (!t) return ''
  const diff = (Date.now() - t * 1000) / 1000
  if (diff < 3600)              return '刚刚上传'
  if (diff < 86400)             return `${Math.floor(diff / 3600)}小时前`
  if (diff < 86400 * 30)        return `${Math.floor(diff / 86400)}天前`
  if (diff < 86400 * 365)       return `${Math.floor(diff / 86400 / 30)}个月前`
  return `${Math.floor(diff / 86400 / 365)}年前`
}

const AVATAR_COLORS = [
  '#ff6b6b', '#ffa94d', '#ffd43b', '#a9e34b',
  '#4dabf7', '#cc5de8', '#f783ac', '#63e6be',
]

function avatarColor(id) {
  return AVATAR_COLORS[Number(id) % AVATAR_COLORS.length]
}

function avatarLetter(author) {
  if (!author) return '?'
  // 中文取第一个字，英文取首字母大写
  const c = [...author][0]
  return c ? c.toUpperCase() : '?'
}
</script>

<style scoped>
.preview-fade-enter-active, .preview-fade-leave-active { transition: opacity 0.2s; }
.preview-fade-enter-from, .preview-fade-leave-to { opacity: 0; }
</style>
