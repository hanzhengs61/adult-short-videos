<template>
  <div class="max-w-screen-lg mx-auto px-3 sm:px-4 py-6">

    <!-- 骨架屏 -->
    <template v-if="loading">
      <div class="animate-pulse space-y-4">
        <div class="aspect-video rounded-xl bg-bg-card"></div>
        <div class="h-6 bg-bg-card rounded w-3/4"></div>
        <div class="h-6 bg-bg-card rounded w-1/2"></div>
        <div class="space-y-2 mt-6">
          <div v-for="i in 8" :key="i" class="h-3 bg-bg-card rounded" :style="{ width: (60 + Math.random()*40) + '%' }"></div>
        </div>
      </div>
    </template>

    <template v-else-if="post">
      <!-- 封面图 -->
      <div v-if="post.cover" class="rounded-xl overflow-hidden mb-5 aspect-video">
        <img :src="post.cover" :alt="post.title" class="w-full h-full object-cover"/>
      </div>

      <!-- 标题 -->
      <h1 class="text-xl sm:text-2xl font-bold text-text-primary leading-snug mb-3">
        {{ post.title }}
      </h1>

      <!-- 元信息行 -->
      <div class="flex flex-wrap items-center gap-3 mb-4 pb-4 border-b border-border/50">
        <span class="flex items-center gap-1 text-xs text-text-muted">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
          {{ formatDate(post.published_at) }}
        </span>
        <span class="flex items-center gap-1 text-xs text-text-muted">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7
                     -1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
          </svg>
          {{ formatCount(post.view_count) }} 次浏览
        </span>
        <!-- 标签 -->
        <div class="flex gap-1.5 flex-wrap ml-auto">
          <span v-for="tag in post.tags" :key="tag"
                @click="goTag(tag)"
                class="px-2 py-0.5 rounded-full bg-primary/15 text-primary text-[10px] cursor-pointer
                       hover:bg-primary/25 transition-colors">
            # {{ tag }}
          </span>
        </div>
      </div>

      <!-- 正文 -->
      <article class="prose prose-invert prose-sm max-w-none
                      prose-headings:text-text-primary prose-p:text-text-secondary
                      prose-a:text-primary prose-img:rounded-lg prose-img:w-full"
               v-html="post.content">
      </article>

      <!-- 分割线 -->
      <div v-if="post.related_videos?.length" class="mt-10 pt-6 border-t border-border/50">
        <h2 class="text-base font-semibold text-text-primary mb-4 flex items-center gap-2">
          <svg class="w-4 h-4 text-primary" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
          相关视频
        </h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
          <router-link v-for="v in post.related_videos" :key="v.video_id"
                       :to="`/video/${v.video_id}`"
                       class="group block rounded-xl overflow-hidden bg-bg-card border border-border/60
                              hover:border-primary/30 transition-all duration-300">
            <div class="relative aspect-video overflow-hidden bg-bg-hover">
              <img :src="v.cover_url" :alt="v.title"
                   class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                   loading="lazy"/>
              <!-- 时长 -->
              <span v-if="v.duration" class="absolute bottom-1.5 right-1.5 px-1 py-0.5 rounded
                           bg-black/70 text-white text-[9px]">
                {{ formatDuration(v.duration) }}
              </span>
            </div>
            <div class="p-2">
              <p class="text-xs font-medium text-text-primary line-clamp-2 leading-snug
                         group-hover:text-primary transition-colors">
                {{ v.title }}
              </p>
            </div>
          </router-link>
        </div>
      </div>

      <!-- 返回按钮 -->
      <div class="mt-8">
        <router-link to="/gossip"
                     class="inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-text-primary transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/>
          </svg>
          返回吃瓜
        </router-link>
      </div>
    </template>

    <!-- 404 -->
    <div v-else class="text-center py-24 text-text-muted">文章不存在或已下架</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useHead } from '@vueuse/head'
import { gossipApi } from '@/api'

const route = useRoute()
const router = useRouter()
const post = ref(null)
const loading = ref(true)

useHead({
  title: () => post.value ? post.value.title + ' - 吃瓜' : '吃瓜详情',
})

onMounted(async () => {
  try {
    const res = await gossipApi.detail(route.params.id)
    post.value = res.data
  } catch {
    post.value = null
  } finally {
    loading.value = false
  }
})

function goTag(tag) {
  router.push({ path: '/gossip', query: { tag } })
}

function formatDate(unix) {
  if (!unix) return ''
  const d = new Date(unix * 1000)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

function formatDuration(s) {
  const m = Math.floor(s / 60)
  const sec = String(s % 60).padStart(2, '0')
  return `${m}:${sec}`
}
</script>
