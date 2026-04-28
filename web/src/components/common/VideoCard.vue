<template>
  <div @click="goFeed"
       class="group block rounded-xl overflow-hidden bg-bg-card border border-border
           hover:border-primary/30 hover:shadow-card-hover transition-all duration-300 hover:-translate-y-0.5 cursor-pointer">
    <!-- 缩略图 -->
    <div class="relative aspect-video overflow-hidden bg-bg-surface">
      <img :src="video.cover_url || video.thumbnail" :alt="video.title"
           class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
           loading="lazy" @error="e => e.target.src = '/placeholder.svg'"/>
      <div class="absolute inset-0 bg-gradient-card opacity-70"></div>

      <!-- 播放按钮 -->
      <div class="absolute inset-0 flex items-center justify-center
                  opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <div class="w-11 h-11 rounded-full bg-primary/90 flex items-center justify-center shadow-glow
                    scale-90 group-hover:scale-100 transition-transform duration-200">
          <svg class="w-5 h-5 text-white ml-0.5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
      </div>

      <!-- 右下：时长 -->
      <div v-if="video.duration"
           class="absolute bottom-1.5 right-1.5 px-1.5 py-0.5 rounded bg-black/80 text-white text-xs font-mono">
        {{ formatDuration(video.duration) }}
      </div>

      <!-- 左上：HOT 徽章 -->
      <div v-if="video.is_hot"
           class="absolute top-1.5 left-1.5 px-1.5 py-0.5 rounded text-[10px] font-bold bg-gradient-primary text-white">
        HOT
      </div>
    </div>

    <!-- 信息 -->
    <div class="p-2.5">
      <h3 class="text-xs font-medium text-text-primary line-clamp-2 leading-snug mb-1.5
                 group-hover:text-primary transition-colors">
        {{ video.title }}
      </h3>
      <div class="flex items-center justify-between text-text-muted text-[11px]">
        <span class="flex items-center gap-1">
          <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
            <path
                d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5C21.27 7.61 17 4.5 12 4.5zm0 12.5c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
          </svg>
          {{ formatCount(video.play_count) }}
        </span>
        <span>{{ formatDate(video.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const props = defineProps({video: {type: Object, required: true}})
const router = useRouter()
function goFeed() {
  router.push({ path: '/feed', query: { id: props.video.video_id } })
}

function formatDuration(s) {
  if (!s) return ''
  const m = Math.floor(s / 60)
  return `${m}:${String(s % 60).padStart(2, '0')}`
}

function formatCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function formatDate(t) {
  if (!t) return ''
  const diff = (Date.now() - new Date(t * 1000).getTime()) / 1000
  if (diff < 86400) return '今天'
  if (diff < 86400 * 7) return `${Math.floor(diff / 86400)}天前`
  const d = new Date(t * 1000)
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>
