<template>
  <div class="max-w-screen-xl mx-auto px-4 py-6">
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="aspect-video bg-bg-card rounded-xl"></div>
      <div class="h-6 bg-bg-card rounded w-3/4"></div>
    </div>

    <template v-else-if="video">
      <div class="flex flex-col lg:flex-row gap-6">
        <!-- 左侧主体 -->
        <div class="flex-1 min-w-0">
          <!-- 播放器 -->
          <div class="relative aspect-video bg-black rounded-xl overflow-hidden mb-4 shadow-2xl">
            <video v-if="playUrl" ref="videoEl" controls class="w-full h-full" autoplay
              :src="playUrl" @play="recordPlay">
            </video>
            <div v-else class="absolute inset-0 flex items-center justify-center cursor-pointer"
              @click="loadPlayer">
              <img :src="video.cover_url" class="w-full h-full object-cover opacity-60" />
              <div class="absolute inset-0 bg-black/40 flex items-center justify-center">
                <div class="w-20 h-20 rounded-full bg-primary/90 flex items-center justify-center shadow-glow
                            hover:scale-105 transition-transform">
                  <svg class="w-8 h-8 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>

          <!-- 标题行 -->
          <div class="flex items-start justify-between gap-4 mb-4">
            <h1 class="text-lg font-bold text-text-primary flex-1">{{ video.title }}</h1>
            <button @click="toggleFavorite"
              :class="['shrink-0 flex items-center gap-1.5 px-3 py-2 rounded-lg border text-sm font-medium transition-all',
                isFavorited
                  ? 'bg-primary/10 border-primary/40 text-primary'
                  : 'border-border text-text-secondary hover:border-primary/40 hover:text-primary']">
              <svg class="w-4 h-4" :fill="isFavorited ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
              </svg>
              {{ isFavorited ? '已收藏' : '收藏' }}
            </button>
          </div>

          <!-- 元数据 -->
          <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-text-muted mb-4 pb-4 border-b border-border">
            <span class="flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5z"/>
              </svg>
              {{ video.play_count?.toLocaleString() }} 次播放
            </span>
            <span v-if="video.duration">时长 {{ formatDuration(video.duration) }}</span>
            <span v-if="video.region">{{ video.region }}</span>
            <span v-if="video.category">{{ video.category }}</span>
          </div>

          <!-- 演员 -->
          <div v-if="video.actors?.length" class="mb-4">
            <p class="text-xs text-text-muted mb-2">演员</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="a in video.actors" :key="a.actor_id" class="tag">{{ a.name }}</span>
            </div>
          </div>

          <!-- 标签 -->
          <div v-if="video.tags" class="mb-6">
            <div class="flex flex-wrap gap-2">
              <span v-for="t in video.tags.split(',')" :key="t" class="tag">{{ t.trim() }}</span>
            </div>
          </div>

          <!-- 评论 -->
          <CommentSection :video-id="videoId" />
        </div>

        <!-- 右侧推荐 -->
        <aside class="w-full lg:w-72 xl:w-80 shrink-0">
          <h2 class="section-title mb-4">
            <span class="w-1 h-4 bg-gradient-primary rounded-full"></span>
            相关视频
          </h2>
          <div class="space-y-3">
            <router-link v-for="v in related" :key="v.video_id" :to="`/video/${v.video_id}`"
              class="flex gap-3 p-2 rounded-lg hover:bg-bg-hover transition-colors group">
              <div class="relative w-32 aspect-video rounded-lg overflow-hidden bg-bg-surface shrink-0">
                <img :src="v.cover_url" class="w-full h-full object-cover group-hover:scale-105 transition-transform" loading="lazy"
                  @error="e => e.target.src = '/placeholder.svg'" />
                <span v-if="v.duration" class="absolute bottom-1 right-1 text-xs bg-black/80 text-white px-1 rounded font-mono">
                  {{ formatDuration(v.duration) }}
                </span>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-text-primary line-clamp-2 leading-snug group-hover:text-primary transition-colors">
                  {{ v.title }}
                </p>
                <p class="text-xs text-text-muted mt-1">{{ formatCount(v.play_count) }} 次</p>
              </div>
            </router-link>
          </div>
        </aside>
      </div>
    </template>

    <div v-else class="text-center py-20 text-text-muted">视频不存在或已被删除</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { videoApi, favoriteApi, playApi } from '@/api'
import { useUserStore } from '@/stores/user'
import CommentSection from '@/components/common/CommentSection.vue'

const route = useRoute()
const userStore = useUserStore()

const videoId = computed(() => Number(route.params.id))
const video = ref(null)
const related = ref([])
const loading = ref(true)
const playUrl = ref('')
const isFavorited = ref(false)
const played = ref(false)

async function fetchVideo() {
  loading.value = true
  playUrl.value = ''
  played.value = false
  try {
    const res = await videoApi.detail(videoId.value)
    video.value = res.data
    if (userStore.isLoggedIn) {
      const fav = await favoriteApi.check(videoId.value).catch(() => null)
      isFavorited.value = fav?.data?.is_favorited || false
    }
    const relRes = await videoApi.list({ page: 1, page_size: 8, category: res.data?.category })
    related.value = (relRes.data?.list || []).filter(v => v.video_id !== videoId.value).slice(0, 6)
  } catch {}
  loading.value = false
}

function loadPlayer() {
  if (!video.value) return
  const url = video.value.play_url || video.value.hot_url || video.value.cold_url
  playUrl.value = url?.startsWith('http')
    ? `/api/storage/proxy?url=${encodeURIComponent(url)}`
    : url
}

function recordPlay() {
  if (played.value || !userStore.isLoggedIn) return
  played.value = true
  playApi.record(videoId.value).catch(() => {})
}

async function toggleFavorite() {
  if (!userStore.isLoggedIn) { alert('请先登录'); return }
  try {
    if (isFavorited.value) {
      await favoriteApi.remove(videoId.value)
    } else {
      await favoriteApi.add(videoId.value)
    }
    isFavorited.value = !isFavorited.value
  } catch {}
}

function formatDuration(s) {
  if (!s) return ''
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

function formatCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

watch(videoId, fetchVideo)
onMounted(fetchVideo)
</script>
