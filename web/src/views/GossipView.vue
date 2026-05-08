<template>
  <div>
    <!-- 吸顶 Tab 栏 -->
    <div class="sticky top-14 z-30 border-b border-border/60"
         style="background:rgba(13,13,15,0.92);backdrop-filter:blur(16px)">
      <div class="max-w-screen-2xl mx-auto px-3 sm:px-4">
        <div class="flex items-center gap-1 py-2 overflow-x-auto scrollbar-none">
          <button v-for="t in tags" :key="t.value"
                  @click="setTag(t.value)"
                  :class="['shrink-0 px-3 py-1.5 rounded-full text-xs font-medium transition-all',
                    activeTag === t.value
                      ? 'bg-primary text-white'
                      : 'text-text-muted hover:text-text-primary hover:bg-bg-hover']">
            {{ t.label }}
          </button>
        </div>
      </div>
    </div>

    <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">

      <!-- 骨架屏 -->
      <div v-if="loading && !posts.length" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <div v-for="i in 6" :key="i"
             class="rounded-xl bg-bg-card border border-border animate-pulse overflow-hidden">
          <div class="aspect-video bg-bg-hover"></div>
          <div class="p-3 space-y-2">
            <div class="h-3 bg-bg-hover rounded w-full"></div>
            <div class="h-3 bg-bg-hover rounded w-4/5"></div>
            <div class="h-2.5 bg-bg-hover rounded w-1/3 mt-3"></div>
          </div>
        </div>
      </div>

      <!-- 文章卡片列表 -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <router-link v-for="post in posts" :key="post.id"
                     :to="`/gossip/${post.id}`"
                     class="group block rounded-xl overflow-hidden bg-bg-card border border-border/60
                            hover:border-primary/30 transition-all duration-300">
          <!-- 封面 -->
          <div class="relative aspect-video overflow-hidden bg-bg-hover">
            <img v-if="post.cover" :src="post.cover" :alt="post.title"
                 class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                 loading="lazy"/>
            <!-- 无封面时的占位 -->
            <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-primary/20 to-bg-card">
              <svg class="w-10 h-10 text-white/20" fill="currentColor" viewBox="0 0 24 24">
                <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
              </svg>
            </div>
            <!-- 标签浮层 -->
            <div v-if="post.tags?.length" class="absolute bottom-2 left-2 flex gap-1 flex-wrap">
              <span v-for="tag in post.tags.slice(0, 3)" :key="tag"
                    class="px-1.5 py-0.5 rounded bg-black/60 text-white/80 text-[9px]">
                {{ tag }}
              </span>
            </div>
          </div>

          <!-- 信息区 -->
          <div class="p-3 space-y-1.5">
            <h3 class="text-sm font-semibold text-text-primary leading-snug line-clamp-2
                       group-hover:text-primary transition-colors">
              {{ post.title }}
            </h3>
            <p v-if="post.summary" class="text-xs text-text-muted leading-relaxed line-clamp-2">
              {{ post.summary }}
            </p>
            <div class="flex items-center gap-3 pt-1">
              <!-- 浏览数 -->
              <span class="flex items-center gap-1 text-[10px] text-text-muted">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7
                           -1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
                {{ formatCount(post.view_count) }}
              </span>
              <!-- 发布时间 -->
              <span class="text-[10px] text-text-muted ml-auto">{{ formatDate(post.published_at) }}</span>
            </div>
          </div>
        </router-link>
      </div>

      <div ref="sentinel" class="h-6 mt-2"></div>
      <p v-if="!hasMore && posts.length" class="text-center py-10 text-text-muted text-xs">── 已经到底了 ──</p>
      <p v-if="!loading && !posts.length" class="text-center py-20 text-text-muted">暂无内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useHead } from '@vueuse/head'
import { gossipApi } from '@/api'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'

useHead({
  title: '吃瓜 - mitun69',
  meta: [{ name: 'description', content: '娱乐爆料、八卦吃瓜，每日更新。' }],
})

// 标签筛选配置
const tags = [
  { label: '全部', value: '' },
  { label: '爆料', value: '爆料' },
  { label: '八卦', value: '八卦' },
  { label: '女优', value: '女优' },
  { label: '新人', value: '新人' },
  { label: '事件', value: '事件' },
]

const posts = ref([])
const page = ref(1)
const activeTag = ref('')

async function fetchMore() {
  const res = await gossipApi.list({ page: page.value, page_size: 18, tag: activeTag.value })
  const list = res.data?.posts || []
  posts.value.push(...list)
  page.value++
  if (list.length < 18) return false
}

const { sentinel, loading, hasMore, reset } = useInfiniteScroll(fetchMore)

function setTag(v) {
  if (activeTag.value === v) return
  activeTag.value = v
  posts.value = []
  page.value = 1
  reset()
}

// 浏览数格式化：超 1 万显示 x.xw
function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

// 时间格式化：今天/昨天/x天前/日期
function formatDate(unix) {
  if (!unix) return ''
  const now = Date.now()
  const diff = now - unix * 1000
  const day = 86400000
  if (diff < day) return '今天'
  if (diff < 2 * day) return '昨天'
  if (diff < 7 * day) return Math.floor(diff / day) + '天前'
  const d = new Date(unix * 1000)
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>
