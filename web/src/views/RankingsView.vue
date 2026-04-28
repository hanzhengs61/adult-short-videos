<template>
  <div class="max-w-screen-xl mx-auto px-3 sm:px-4 py-4">

    <!-- 控制栏 -->
    <div class="flex flex-col gap-2 mb-6">
      <div class="flex bg-bg-surface border border-border rounded-xl p-1 gap-1">
        <button v-for="t in types" :key="t.value" @click="activeType = t.value"
                :class="['px-4 py-1.5 rounded-lg text-xs font-semibold transition-all',
            activeType===t.value ? 'bg-gradient-primary text-white' : 'text-text-muted hover:text-text-primary']">
          {{ t.label }}
        </button>
      </div>

      <!-- 视频榜才显示时间筛选 -->
      <div v-if="activeType !== 'creator'" class="flex gap-1.5">
        <button v-for="p in periods" :key="p.value" @click="activePeriod = p.value"
                :class="['px-3.5 py-1.5 rounded-full text-xs border font-medium transition-all',
            activePeriod===p.value
              ? 'border-primary/60 text-primary bg-primary/10'
              : 'border-border text-text-muted hover:border-border-light hover:text-text-primary']">
          {{ p.label }}
        </button>
      </div>

      <div class="ml-auto hidden sm:flex items-center gap-1.5 text-xs text-text-muted">
        <svg class="w-3.5 h-3.5 text-primary" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
        </svg>
        更新于今日 00:00
      </div>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 10" :key="i" class="h-20 bg-bg-surface rounded-xl animate-pulse border border-border/50"></div>
    </div>

    <!-- 创作者榜 -->
    <template v-else-if="activeType === 'creator'">
      <div v-if="creators.length" class="space-y-3">
        <div v-for="(c, i) in creators" :key="c.user_id"
             class="flex items-center gap-3 p-3 rounded-xl bg-bg-surface border border-border/50
                    hover:border-border hover:bg-bg-hover transition-all cursor-pointer group">
          <!-- 排名 -->
          <div :class="['w-7 text-center font-black text-sm shrink-0',
            i===0 ? 'text-yellow-500' : i===1 ? 'text-gray-400' : i===2 ? 'text-amber-600' : 'text-text-muted']">
            {{ i + 1 }}
          </div>
          <!-- 头像 -->
          <div :class="['w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shrink-0',
            i===0 ? 'bg-yellow-500' : i===1 ? 'bg-gray-400' : i===2 ? 'bg-amber-600' : 'bg-gradient-primary']">
            {{ (c.username || '?')[0].toUpperCase() }}
          </div>
          <!-- 信息 -->
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-text-primary group-hover:text-primary transition-colors truncate">
              {{ c.username }}
            </p>
            <div class="flex items-center gap-3 mt-0.5 text-xs text-text-muted">
              <span>{{ c.video_count }} 个视频</span>
              <span>{{ formatCount(c.total_plays) }} 播放</span>
            </div>
          </div>
          <!-- 奖牌 -->
          <span v-if="i===0" class="text-lg shrink-0">🥇</span>
          <span v-else-if="i===1" class="text-lg shrink-0">🥈</span>
          <span v-else-if="i===2" class="text-lg shrink-0">🥉</span>
        </div>
      </div>
      <div v-else class="text-center py-20 text-text-muted">暂无创作者数据</div>
    </template>

    <!-- 视频榜 -->
    <template v-else-if="videos.length">
      <!-- TOP 3 大卡片 -->
      <div class="grid grid-cols-3 gap-3 mb-6">
        <div v-for="(v, i) in videos.slice(0,3)" :key="v.video_id"
             :class="['relative rounded-xl overflow-hidden border cursor-pointer group transition-all hover:-translate-y-0.5',
            i===0 ? 'border-yellow-500/40 ring-1 ring-yellow-500/20' :
            i===1 ? 'border-gray-400/30' : 'border-amber-700/30']"
             @click="$router.push(`/video/${v.video_id}`)">
          <div class="relative aspect-video overflow-hidden">
            <img :src="v.cover_url" :alt="v.title"
                 class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                 @error="e => e.target.src='/placeholder.svg'" loading="lazy"/>
            <div class="absolute inset-0 bg-gradient-card"></div>
            <div :class="['absolute top-2 left-2 w-7 h-7 rounded-lg flex items-center justify-center font-black text-sm',
              i===0 ? 'bg-yellow-500 text-black' : i===1 ? 'bg-gray-300 text-black' : 'bg-amber-700 text-white']">
              {{ i + 1 }}
            </div>
            <div class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
              <div class="w-10 h-10 rounded-full bg-white/20 backdrop-blur flex items-center justify-center">
                <svg class="w-5 h-5 text-white ml-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
              </div>
            </div>
          </div>
          <div class="p-2.5">
            <p class="text-xs font-medium text-text-primary line-clamp-2 leading-snug mb-1">{{ v.title }}</p>
            <div class="flex items-center justify-between text-[11px] text-text-muted">
              <span>{{ formatCount(v.play_count) }} 播放</span>
              <span v-if="i===0" class="text-yellow-500 font-bold text-[10px]">🥇 冠军</span>
              <span v-else-if="i===1" class="text-gray-400 font-bold text-[10px]">🥈 亚军</span>
              <span v-else class="text-amber-600 font-bold text-[10px]">🥉 季军</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 4-20 名列表 -->
      <div class="space-y-2">
        <div v-for="(v, i) in videos.slice(3, 20)" :key="v.video_id"
             class="flex items-center gap-3 p-3 rounded-xl bg-bg-surface border border-border/50
                    hover:border-border hover:bg-bg-hover cursor-pointer group transition-all"
             @click="$router.push(`/video/${v.video_id}`)">
          <div :class="['w-7 text-center font-black text-sm shrink-0',
            i < 3 ? 'text-primary' : 'text-text-muted']">{{ i + 4 }}</div>
          <div class="relative w-24 sm:w-28 aspect-video rounded-lg overflow-hidden bg-bg-card shrink-0">
            <img :src="v.cover_url" class="w-full h-full object-cover group-hover:scale-105 transition-transform"
                 loading="lazy" @error="e => e.target.src='/placeholder.svg'"/>
            <span v-if="v.duration"
                  class="absolute bottom-1 right-1 px-1 rounded bg-black/80 text-white text-[10px] font-mono">
              {{ formatDuration(v.duration) }}
            </span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm text-text-primary line-clamp-2 leading-snug group-hover:text-primary transition-colors">
              {{ v.title }}
            </p>
            <div class="flex items-center gap-3 mt-1.5 text-xs text-text-muted flex-wrap">
              <span class="flex items-center gap-1">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5C21.27 7.61 17 4.5 12 4.5z"/>
                </svg>
                {{ formatCount(v.play_count) }}
              </span>
              <span v-if="v.region">{{ v.region }}</span>
            </div>
          </div>
          <svg class="w-4 h-4 text-text-muted shrink-0 hidden sm:block" fill="none" stroke="currentColor"
               stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
          </svg>
        </div>
      </div>
    </template>

    <div v-else-if="!loading" class="text-center py-20 text-text-muted">暂无数据</div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { videoApi, userApi } from '@/api'

const activeType = ref('short')
const activePeriod = ref('day')
const videos = ref([])
const creators = ref([])
const loading = ref(true)

const types = [
  { label: '短视频', value: 'short' },
  { label: '长视频', value: 'long' },
  { label: '创作者', value: 'creator' },
]
const periods = [
  { label: '日榜', value: 'day' },
  { label: '周榜', value: 'week' },
  { label: '月榜', value: 'month' },
]

async function fetchRankings() {
  loading.value = true
  try {
    if (activeType.value === 'creator') {
      const res = await userApi.creators(20)
      creators.value = res.data?.creators || []
    } else {
      const res = await videoApi.popular({
        page: 1, page_size: 20,
        period: activePeriod.value, type: activeType.value
      })
      videos.value = res.data?.videos || []
    }
  } catch {
    videos.value = []
    creators.value = []
  }
  loading.value = false
}

function formatCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function formatDuration(s) {
  if (!s) return ''
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

watch([activeType, activePeriod], fetchRankings)
onMounted(fetchRankings)
</script>
