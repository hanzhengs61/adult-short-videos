<template>
  <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">

    <!-- 顶部搜索框 -->
    <form class="mb-5" @submit.prevent="doSearch()">
      <div class="relative">
        <input v-model="searchInput" type="text" placeholder="搜索视频、创作者..."
               class="w-full h-11 pl-11 pr-10 bg-bg-surface border border-border rounded-full text-text-primary
                 text-sm placeholder-text-muted outline-none transition-all
                 focus:border-primary/60 focus:ring-1 focus:ring-primary/20"/>
        <svg class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted pointer-events-none"
             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
             viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
        <button v-if="searchInput" type="button" @click="clearSearch"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </form>

    <!-- 搜索结果（有搜索词时始终展示） -->
    <div v-if="currentQuery" class="animate-slide-up">
      <div class="flex items-center justify-between mb-3">
        <h2 class="section-title">
          <span class="w-1 h-4 bg-gradient-primary rounded-full inline-block"></span>
          搜索：<span class="text-primary">{{ currentQuery }}</span>
          <span class="text-text-muted text-xs font-normal ml-1">({{ searchResults.length }} 个结果)</span>
        </h2>
        <button @click="clearSearch" class="text-xs text-text-muted hover:text-primary transition-colors">清除</button>
      </div>

      <!-- 有结果时显示排序 + 瀑布流 -->
      <template v-if="searchResults.length">
        <div class="mb-3">
          <SortBar :modelValue="searchSort" @update:modelValue="setSearchSort"/>
        </div>
        <div class="flex gap-3 items-start mb-6">
          <div v-for="(col, colIdx) in searchResultColumns" :key="`result-col-${colIdx}`" class="flex-1 space-y-3">
            <div v-for="v in col" :key="v.video_id">
              <VideoCard :video="v" :order-by="searchSort"/>
            </div>
          </div>
        </div>
      </template>

      <!-- 无结果空状态 -->
      <div v-else-if="!searching" class="flex flex-col items-center py-20 gap-4 text-center">
        <svg class="w-16 h-16 text-text-muted opacity-40" fill="none" stroke="currentColor" stroke-width="1.2" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"/><path stroke-linecap="round" d="M21 21l-4.35-4.35"/>
        </svg>
        <p class="text-text-primary font-medium">没有找到「{{ currentQuery }}」相关内容</p>
        <p class="text-text-muted text-sm">换个关键词试试？</p>
        <button @click="clearSearch" class="btn-ghost text-xs px-5 py-2 mt-1">清除搜索</button>
      </div>

      <!-- 搜索中骨架 -->
      <div v-else class="flex gap-3 items-start">
        <div v-for="(col, colIdx) in skeletonColumns" :key="`search-sk-col-${colIdx}`" class="flex-1 space-y-3">
          <div v-for="i in col" :key="`search-sk-${colIdx}-${i}`"
               class="rounded-xl bg-bg-card border border-border animate-pulse">
            <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
            <div class="p-2.5 space-y-1.5">
              <div class="h-2.5 bg-bg-hover rounded w-full"></div>
              <div class="h-2 bg-bg-hover rounded w-2/3"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 探索内容（无搜索词时） -->
    <template v-else>
      <!-- 最近搜索 -->
      <div v-if="exploreStore.history.length" class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h2 class="section-title text-sm">
            <svg class="w-4 h-4 text-text-muted" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            最近搜索
          </h2>
          <button @click="exploreStore.clearHistory()"
                  class="text-xs text-text-muted hover:text-red-400 transition-colors">全部清除
          </button>
        </div>
        <div class="flex flex-wrap gap-2">
          <button v-for="h in exploreStore.history" :key="h"
                  class="group flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-bg-surface border border-border
                   text-text-primary/60 text-sm hover:border-primary/50 hover:text-primary transition-all">
            <span @click="quickSearch(h)">{{ h }}</span>
            <span @click.stop="exploreStore.removeHistory(h)"
                  class="text-text-muted group-hover:text-primary/70 hover:text-red-400! transition-colors leading-none">×</span>
          </button>
        </div>
      </div>

      <!-- 热搜榜 -->
      <div class="mb-6">
        <h2 class="section-title text-sm mb-3">
          <svg class="w-4 h-4 text-primary" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z"/>
          </svg>
          热搜关键词
        </h2>
        <div class="flex flex-wrap gap-2">
          <button v-for="(k, i) in hotKeywords" :key="k" @click="quickSearch(k)"
                  :class="['flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm border transition-all',
              i < 3
                ? 'bg-primary/10 border-primary/30 text-primary hover:bg-primary/20'
                : 'bg-bg-surface border-border text-text-primary/60 hover:border-border-light hover:text-text-primary']">
            <span :class="['font-bold text-[10px] w-3.5 text-center',
              i === 0 ? 'text-red-400' : i === 1 ? 'text-orange-400' : i === 2 ? 'text-yellow-400' : 'text-text-muted']">
              {{ i + 1 }}
            </span>
            {{ k }}
            <span v-if="i < 3" class="text-[9px] text-primary font-semibold">HOT</span>
          </button>
        </div>
      </div>

      <!-- 广告位 -->
      <div class="mb-6 rounded-xl border border-border/40 bg-bg-surface/60 h-12
                  flex items-center justify-center text-text-muted text-xs tracking-widest">
        AD · 广告位
      </div>

      <!-- 推荐创作者 -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h2 class="section-title text-sm">
            <svg class="w-4 h-4 text-text-muted" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
            </svg>
            推荐创作者
          </h2>
        </div>

        <!-- 横滑创作者卡片 -->
        <div class="flex gap-3 overflow-x-auto scrollbar-none pb-3 snap-x snap-mandatory -mx-3 px-3">
          <div v-for="actor in actors" :key="actor.actor_id"
               class="snap-start shrink-0 w-44 bg-bg-surface border border-border rounded-xl overflow-hidden
                   hover:border-primary/30 hover:shadow-card-hover transition-all duration-300 cursor-pointer group"
               @click="quickSearch(actor.name)">

            <!-- 创作者头像 -->
            <div class="p-3 flex items-center gap-2.5 border-b border-border/50">
              <div class="w-9 h-9 rounded-full bg-gradient-primary p-0.5 shrink-0 group-hover:shadow-glow transition-shadow">
                <div class="w-full h-full rounded-full bg-bg overflow-hidden flex items-center justify-center">
                  <img v-if="actor.avatar" :src="actor.avatar" :alt="actor.name"
                       class="w-full h-full object-cover"
                       @error="e => e.target.style.display='none'"/>
                  <span v-if="!actor.avatar" class="text-white font-bold text-sm">{{ actor.name[0] }}</span>
                </div>
              </div>
              <div class="min-w-0">
                <p class="text-sm font-semibold text-text-primary truncate group-hover:text-primary transition-colors">
                  {{ actor.name }}
                </p>
                <p class="text-xs text-text-primary/60">创作者</p>
              </div>
            </div>

            <!-- 视频预览缩略图（3宫格） -->
            <div class="grid grid-cols-3 gap-0.5 p-0.5">
              <div v-for="v in (actor.videos || []).slice(0,3)" :key="v?.video_id"
                   class="aspect-square bg-bg-hover overflow-hidden">
                <img v-if="v?.cover_url" :src="v.cover_url"
                     class="w-full h-full object-cover opacity-80 group-hover:opacity-100 transition-opacity"
                     loading="lazy" @error="e => e.target.style.display='none'"/>
                <div v-else class="w-full h-full bg-bg-hover flex items-center justify-center">
                  <svg class="w-4 h-4 text-text-muted" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z"/>
                  </svg>
                </div>
              </div>
              <div v-for="i in Math.max(0, 3-(actor.videos||[]).length)" :key="`empty${i}`"
                   class="aspect-square bg-bg-hover"></div>
            </div>

            <div class="px-3 py-2">
              <p class="text-xs text-text-primary/60">点击查看全部作品</p>
            </div>
          </div>

          <!-- 骨架 -->
          <template v-if="actorsLoading">
            <div v-for="i in 5" :key="`ask${i}`"
                 class="snap-start shrink-0 w-44 bg-bg-surface border border-border rounded-xl overflow-hidden animate-pulse">
              <div class="h-14 bg-bg-hover"></div>
              <div class="grid grid-cols-3 gap-0.5 p-0.5">
                <div v-for="j in 3" :key="j" class="aspect-square bg-bg-hover"></div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useHead } from '@vueuse/head'
import VideoCard from '@/components/common/VideoCard.vue'
import SortBar from '@/components/common/SortBar.vue'
import { useExploreStore } from '@/stores/explore'
import { searchApi, userApi } from '@/api'

const route = useRoute()
const router = useRouter()
const exploreStore = useExploreStore()
const searchInput = ref('')
const currentQuery = ref('')
const searchResults = ref([])
const searchSort = ref('created_at')
const searching = ref(false)
const actors = ref([])
const actorsLoading = ref(true)
const columnCount = ref(2)

useHead(computed(() => ({
  title: currentQuery.value
    ? `搜索「${currentQuery.value}」- MITUN69`
    : 'MITUN69 - 探索发现',
  meta: [
    { name: 'description', content: currentQuery.value
        ? `MITUN69 搜索「${currentQuery.value}」相关视频，共 ${searchResults.value.length} 个结果。`
        : '探索热门关键词、热搜榜、推荐创作者，发现更多精彩内容。' },
  ],
})))

const searchResultColumns = computed(() => {
  const cols = Array.from({ length: columnCount.value }, () => [])
  searchResults.value.forEach((video, idx) => {
    cols[idx % columnCount.value].push(video)
  })
  return cols
})

const hotKeywords = ['中文字幕', '无码', '人妻', '美乳', '国产自拍', '制服诱惑', '美少女', '素人', '4K高清', '巨乳']

// 瀑布流骨架
const skeletonColumns = computed(() => {
  const total = 10
  const cols = Array.from({ length: columnCount.value }, () => [])
  for (let i = 0; i < total; i++) {
    cols[i % columnCount.value].push(i)
  }
  return cols
})

// 监听窗口大小
function syncColumnCount() {
  const w = window.innerWidth
  if (w >= 1280) columnCount.value = 5
  else if (w >= 768) columnCount.value = 4
  else if (w >= 640) columnCount.value = 3
  else columnCount.value = 2
}

// 搜索
async function doSearch(keyword) {
  const q = (keyword ?? searchInput.value).trim()
  if (!q) return
  currentQuery.value = q
  searchInput.value = q
  searchResults.value = []
  searching.value = true
  exploreStore.addHistory(q)
  try {
    const res = await searchApi.videos({
      keyword: q,
      page: 1,
      page_size: 10,
      order_by: searchSort.value,
    })
    searchResults.value = res.data?.videos || []
  } catch {
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

// 快速搜索
function quickSearch(q) {
  doSearch(q)
}

// 排序
function setSearchSort(v) {
  if (searchSort.value === v) return
  searchSort.value = v
  if (currentQuery.value) doSearch(currentQuery.value)
}

// 清空搜索
function clearSearch() {
  searchInput.value = ''
  currentQuery.value = ''
  searchResults.value = []
  router.replace({path: '/explore', query: {}})
}

// 获取创作者
async function fetchActors() {
  actorsLoading.value = true
  try {
    const res = await userApi.creators(10)
    actors.value = (res.data?.creators || []).map(c => ({
      actor_id: c.user_id,
      name: c.username,
      avatar: c.avatar || '',
      videos: [],
    }))
    if (!actors.value.length) {
      actors.value = hotKeywords.slice(0, 8).map((name, i) => ({actor_id: i, name, videos: []}))
    }
  } catch {
    actors.value = hotKeywords.slice(0, 8).map((name, i) => ({actor_id: i, name, videos: []}))
  }
  actorsLoading.value = false
}

// 响应 NavBar 搜索跳转过来的 query
watch(() => route.query.q, (newQ) => {
  if (newQ) {
    searchInput.value = String(newQ)
    doSearch()
  }
}, {immediate: true})

onMounted(fetchActors)
onMounted(() => {
  syncColumnCount()
  window.addEventListener('resize', syncColumnCount)
})

onUnmounted(() => {
  window.removeEventListener('resize', syncColumnCount)
})
</script>
