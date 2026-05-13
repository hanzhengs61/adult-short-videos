<template>
  <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">

    <!-- 顶部搜索框 -->
    <form class="mb-5" @submit.prevent="doSearch()">
      <div class="flex items-center h-11 bg-bg-surface border border-border rounded-full
                  transition-all focus-within:border-primary/60 focus-within:ring-1 focus-within:ring-primary/20">
        <!-- 搜索图标 -->
        <svg class="ml-3.5 shrink-0 w-4 h-4 text-text-muted pointer-events-none"
             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
             viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8"/>
          <path d="M21 21l-4.35-4.35"/>
        </svg>
        <!-- 类型自定义下拉选择器 -->
        <div class="relative ml-2 shrink-0">
          <button type="button" @click.stop="toggleCategoryMenu"
                  class="flex items-center gap-1 text-xs font-medium text-text-primary/80 outline-none cursor-pointer select-none">
            {{ searchTab === 'video' ? '短视频' : '吃瓜' }}
            <svg class="w-3 h-3 text-text-muted transition-transform duration-200"
                 :class="showCategoryMenu ? 'rotate-180' : ''"
                 fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>
          <!-- 下拉面板 -->
          <div v-if="showCategoryMenu"
               class="absolute top-full left-0 mt-2 w-24 bg-bg-surface border border-border
                      rounded-xl shadow-lg z-50 py-1 overflow-hidden">
            <button type="button"
                    v-for="opt in [{value:'video',label:'短视频'},{value:'gossip',label:'吃瓜'}]"
                    :key="opt.value"
                    @click="selectCategory(opt.value)"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm hover:bg-bg-hover transition-colors">
              <span :class="['w-3.5 h-3.5 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors',
                searchTab === opt.value ? 'border-primary' : 'border-border']">
                <span v-if="searchTab === opt.value" class="w-1.5 h-1.5 rounded-full bg-primary"/>
              </span>
              <span :class="searchTab === opt.value ? 'text-primary font-medium' : 'text-text-primary/70'">
                {{ opt.label }}
              </span>
            </button>
          </div>
        </div>
        <!-- 分隔线 -->
        <div class="w-px h-4 bg-border mx-2 shrink-0"/>
        <!-- 输入框 -->
        <input v-model="searchInput" type="text" placeholder="搜索..."
               class="flex-1 min-w-0 bg-transparent text-sm text-text-primary placeholder-text-primary/50 outline-none"/>
        <!-- 清除按钮 -->
        <button v-if="searchInput" type="button" @click="clearSearch"
                class="mr-3 shrink-0 text-text-muted hover:text-text-primary transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </form>

    <!-- 搜索结果（有搜索词时始终展示） -->
    <div v-if="currentQuery" class="animate-slide-up">
      <div class="flex items-center justify-between mb-4">
        <h2 class="section-title">
          <span class="w-1 h-4 bg-gradient-primary rounded-full inline-block"></span>
          搜索：<span class="text-primary">{{ currentQuery }}</span>
          <span class="text-text-muted text-xs font-normal ml-1">
            ({{ searchTab === 'video' ? searchResults.length : gossipResults.length }} 个结果)
          </span>
        </h2>
        <button @click="clearSearch" class="text-xs text-text-muted hover:text-primary transition-colors">清除</button>
      </div>

      <!-- 短视频结果 -->
      <template v-if="searchTab === 'video'">
        <!-- 有结果时显示排序 + 瀑布流 -->
        <template v-if="searchResults.length">
          <div class="mb-3">
            <SortBar :modelValue="searchSort" @update:modelValue="setSearchSort"/>
          </div>
          <div class="flex gap-3 items-start mb-6">
            <div v-for="(col, colIdx) in searchResultColumns" :key="`result-col-${colIdx}`" class="flex-1 space-y-3">
              <div v-for="v in col" :key="v.video_id">
                <VideoCard :video="v" :order-by="searchSort" :keyword="currentQuery"/>
              </div>
            </div>
          </div>
          <!-- 底部加载更多 -->
          <div ref="searchSentinel" class="flex justify-center py-6">
            <div v-if="searchLoadingMore" class="w-6 h-6 border-2 border-primary/30 border-t-primary rounded-full animate-spin"/>
            <p v-else-if="!searchHasMore" class="text-text-muted text-sm">已加载全部结果</p>
          </div>
        </template>

        <!-- 无结果空状态 -->
        <div v-else-if="!searching" class="flex flex-col items-center py-20 gap-4 text-center">
          <svg class="w-16 h-16 text-text-muted opacity-40" fill="none" stroke="currentColor" stroke-width="1.2" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8"/><path stroke-linecap="round" d="M21 21l-4.35-4.35"/>
          </svg>
          <p class="text-text-primary font-medium">没有找到「{{ currentQuery }}」相关短视频</p>
          <p class="text-text-muted text-sm">换个关键词试试？</p>
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
      </template>

      <!-- 吃瓜结果 -->
      <template v-else>
        <!-- 有结果时显示卡片列表 -->
        <template v-if="gossipResults.length">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 mb-6">
            <GossipCard v-for="post in gossipResults" :key="post.id" :post="post"/>
          </div>
          <!-- 底部加载更多 -->
          <div ref="gossipSentinel" class="flex justify-center py-6">
            <div v-if="gossipLoadingMore" class="w-6 h-6 border-2 border-primary/30 border-t-primary rounded-full animate-spin"/>
            <p v-else-if="!gossipHasMore" class="text-text-muted text-sm">已加载全部结果</p>
          </div>
        </template>

        <!-- 无结果空状态 -->
        <div v-else-if="!gossipSearching" class="flex flex-col items-center py-20 gap-4 text-center">
          <svg class="w-16 h-16 text-text-muted opacity-40" fill="none" stroke="currentColor" stroke-width="1.2" viewBox="0 0 24 24">
            <circle cx="11" cy="11" r="8"/><path stroke-linecap="round" d="M21 21l-4.35-4.35"/>
          </svg>
          <p class="text-text-primary font-medium">没有找到「{{ currentQuery }}」相关吃瓜内容</p>
          <p class="text-text-muted text-sm">换个关键词试试？</p>
        </div>

        <!-- 搜索中骨架 -->
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <div v-for="i in 6" :key="i"
               class="rounded-xl bg-bg-card border border-border animate-pulse overflow-hidden">
            <div class="aspect-video bg-bg-hover"></div>
            <div class="p-3 space-y-2">
              <div class="h-3 bg-bg-hover rounded w-full"></div>
              <div class="h-3 bg-bg-hover rounded w-4/5"></div>
            </div>
          </div>
        </div>
      </template>
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
                  class="text-xs text-text-primary/60 hover:text-red-400 transition-colors">全部清除
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
          热搜榜
        </h2>
        <!-- 骨架 -->
        <div v-if="hotTagsLoading" class="flex flex-wrap gap-2">
          <div v-for="i in 8" :key="i"
               class="h-8 rounded-full bg-bg-surface border border-border animate-pulse"
               :style="`width:${56 + (i % 3) * 16}px`"/>
        </div>
        <!-- 标签列表 -->
        <div v-else class="flex flex-wrap gap-2">
          <button v-for="(tag, i) in hotTags" :key="tag.name"
                  @click="clickHotTag(tag.name)"
                  :class="['flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm border transition-all',
                    i < 3
                      ? 'bg-primary/10 border-primary/30 text-primary hover:bg-primary/20'
                      : 'bg-bg-surface border-border text-text-primary/60 hover:border-border-light hover:text-text-primary']">
            <span :class="['font-bold text-[10px] w-3.5 text-center shrink-0',
              i === 0 ? 'text-red-400' : i === 1 ? 'text-orange-400' : i === 2 ? 'text-yellow-400' : 'text-text-muted']">
              {{ i + 1 }}
            </span>
            {{ tag.name }}
            <span v-if="tag.click_count > 0"
                  class="text-[9px] text-text-muted/60 shrink-0">
              {{ fmtClick(tag.click_count) }}
            </span>
            <span v-if="i < 3" class="text-[9px] text-primary font-semibold shrink-0">HOT</span>
          </button>
        </div>
      </div>

      <!-- 广告位 -->
      <div class="mb-6 rounded-xl border border-border/40 bg-bg-surface/60 h-12
                  flex items-center justify-center text-text-muted text-xs tracking-widest">
        AD · 广告位
      </div>

      <!-- 猜你喜欢 -->
      <div class="mb-6">
        <RecommendSection/>
      </div>

    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useHead } from '@vueuse/head'
import VideoCard from '@/components/common/VideoCard.vue'
import SortBar from '@/components/common/SortBar.vue'
import GossipCard from '@/components/common/GossipCard.vue'
import RecommendSection from '@/components/common/RecommendSection.vue'
import { useExploreStore } from '@/stores/explore'
import { tagApi } from '@/api'
import { useExploreSearch } from '@/composables/useExploreSearch'

const route = useRoute()
const router = useRouter()
const exploreStore = useExploreStore()
const showCategoryMenu = ref(false)
// 热搜标签：从后端拉取，含点击数
const hotTags = ref([])
const hotTagsLoading = ref(true)

const {
  searchInput, currentQuery, searchTab,
  searchResults, searchSort, searching, searchLoadingMore, searchHasMore, searchSentinel,
  gossipResults, gossipSearching, gossipLoadingMore, gossipHasMore, gossipSentinel,
  searchResultColumns, skeletonColumns,
  setSearchTab, doSearch, setSearchSort, clearSearch,
} = useExploreSearch({ router, exploreStore })

useHead(computed(() => ({
  title: currentQuery.value ? `搜索「${currentQuery.value}」- MITUN69` : 'MITUN69 - 探索发现',
  meta: [{ name: 'description', content: currentQuery.value
      ? `MITUN69 搜索「${currentQuery.value}」相关内容。`
      : '探索热门关键词、热搜榜、推荐创作者，发现更多精彩内容。' }],
})))

function toggleCategoryMenu() {
  showCategoryMenu.value = !showCategoryMenu.value
  if (showCategoryMenu.value) {
    // 下一帧挂全局监听，点击其他区域关闭
    setTimeout(() => document.addEventListener('click', closeCategoryMenu, { once: true }), 0)
  }
}

function closeCategoryMenu() {
  showCategoryMenu.value = false
}

function selectCategory(val) {
  setSearchTab(val)
  showCategoryMenu.value = false
}

// 点击热搜标签：上报点击数（fire-and-forget）并触发搜索
function clickHotTag(name) {
  tagApi.click(name).catch(() => {})
  // 乐观更新本地计数，避免等待接口
  const tag = hotTags.value.find(t => t.name === name)
  if (tag) tag.click_count++
  doSearch(name)
}

function quickSearch(q) { doSearch(q) }

// 点击数格式化：1000 → 1k，10000 → 1w
function fmtClick(n) {
  if (n >= 10000) return (n / 10000).toFixed(1).replace(/\.0$/, '') + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
  return String(n)
}

async function fetchHotTags() {
  hotTagsLoading.value = true
  try {
    const res = await tagApi.hot(10)
    hotTags.value = res.data?.tags || []
  } catch {
    hotTags.value = []
  }
  hotTagsLoading.value = false
}

watch(() => route.query.q, (newQ) => {
  if (newQ) { searchInput.value = String(newQ); doSearch() }
}, { immediate: true })

onMounted(fetchHotTags)
</script>
