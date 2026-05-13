<template>
  <div>
    <!-- 吸顶 Tab 栏 -->
    <div class="sticky top-14 z-30 border-b border-border/60"
         style="background:rgba(13,13,15,0.92);backdrop-filter:blur(16px)">
      <div class="max-w-screen-2xl mx-auto px-3 sm:px-4">
        <div class="flex items-center gap-1 py-2 overflow-x-auto scrollbar-none">
          <button v-for="t in tags" :key="t.value"
                  @click="setTag(t.value)"
                  :class="['shrink-0 px-3 py-1.5 rounded-full text-sm font-medium transition-all',
                    activeTag === t.value
                      ? 'bg-primary text-white'
                      : 'text-text-primary/60 hover:text-text-primary hover:bg-bg-hover']">
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
        <GossipCard v-for="post in posts" :key="post.id" :post="post"/>
      </div>

      <div ref="sentinel" class="h-6 mt-2"></div>
      <p v-if="!hasMore && posts.length" class="text-center py-10 text-text-muted text-xs">── 已经到底了 ──</p>
      <p v-if="!loading && !posts.length" class="text-center py-20 text-text-muted">暂无内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useHead } from '@vueuse/head'
import { gossipApi } from '@/api'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import GossipCard from '@/components/common/GossipCard.vue'

useHead({
  title: '吃瓜 - mitun69',
  meta: [{ name: 'description', content: '娱乐爆料、八卦吃瓜，每日更新。' }],
})

// 标签筛选配置
const tags = [
  { label: '全部', value: '' },
  { label: '今日吃瓜', value: '今日吃瓜' },
  { label: '最高点击', value: '最高点击' },
  { label: '必吃大瓜', value: '必吃大瓜' },
  { label: '学生校园', value: '学生校园' },
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

</script>
