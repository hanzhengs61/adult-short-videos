<template>
  <div class="max-w-screen-2xl mx-auto px-3 sm:px-4 py-4">

    <!-- 未登录引导 -->
    <div v-if="!userStore.isLoggedIn"
         class="flex flex-col items-center justify-center py-28 gap-5 animate-fade-in">
      <div class="w-24 h-24 rounded-full bg-bg-surface border-2 border-border flex items-center justify-center">
        <svg class="w-12 h-12 text-text-muted" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
      </div>
      <div class="text-center">
        <p class="text-text-primary font-bold text-xl mb-1">关注喜欢的创作者</p>
        <p class="text-text-muted text-sm">登录后可关注，第一时间收到最新视频</p>
      </div>
      <button @click="userStore.openAuth('login')" class="btn-primary px-10 py-3 text-sm font-bold">
        立即登录
      </button>
      <router-link to="/explore" class="text-xs text-text-muted hover:text-primary transition-colors">
        或先去探索内容 →
      </router-link>
    </div>

    <template v-else>
      <!-- 已关注创作者 -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h2 class="section-title">
            <span class="w-1 h-4 bg-gradient-primary rounded-full inline-block"></span>
            我的关注
          </h2>
          <router-link to="/explore" class="text-xs text-primary hover:underline">发现更多</router-link>
        </div>

        <!-- 无关注状态 -->
        <div v-if="!following.length"
             class="flex items-center gap-4 p-4 rounded-xl bg-bg-surface border border-border">
          <div class="w-12 h-12 rounded-full bg-bg-hover flex items-center justify-center shrink-0">
            <svg class="w-6 h-6 text-text-muted" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
            </svg>
          </div>
          <div>
            <p class="text-text-primary text-sm font-medium">还没有关注任何人</p>
            <p class="text-text-muted text-xs mt-0.5">去探索页找到感兴趣的创作者</p>
          </div>
          <router-link to="/explore" class="ml-auto btn-ghost text-xs py-1.5 px-3 shrink-0">去探索</router-link>
        </div>

        <!-- 创作者头像横滑 -->
        <div v-else class="flex gap-3 overflow-x-auto scrollbar-none pb-1">
          <div v-for="u in following" :key="u.id"
               class="shrink-0 flex flex-col items-center gap-1.5 cursor-pointer group"
               @click="scrollToCreator(u.id)">
            <div class="relative">
              <div class="w-14 h-14 rounded-full p-0.5"
                   :class="u.hasNew ? 'bg-gradient-primary' : 'bg-border'">
                <div class="w-full h-full rounded-full bg-bg overflow-hidden flex items-center justify-center">
                  <img v-if="u.avatar" :src="u.avatar" :alt="u.name"
                       class="w-full h-full object-cover"
                       @error="e => e.target.style.display='none'"/>
                  <span v-if="!u.avatar" class="font-bold text-lg text-text-primary">{{ u.name[0] }}</span>
                </div>
              </div>
              <span v-if="u.hasNew"
                    class="absolute -top-0.5 -right-0.5 w-3.5 h-3.5 rounded-full bg-primary border-2 border-bg"></span>
            </div>
            <span class="text-[11px] text-text-secondary group-hover:text-primary transition-colors
                         max-w-[56px] truncate text-center">{{ u.name }}</span>
          </div>
        </div>
      </div>

      <!-- 动态 Feed：仅在有关注时展示 -->
      <template v-if="following.length">
        <div class="mb-3">
          <h2 class="section-title">
            <span class="w-1 h-4 bg-gradient-primary rounded-full inline-block"></span>
            最新动态
          </h2>
        </div>

        <div class="columns-2 sm:columns-3 md:columns-4 xl:columns-5 gap-3">
          <div v-for="v in feedVideos" :key="v.video_id" class="break-inside-avoid mb-3">
            <VideoCard :video="v"/>
          </div>
          <template v-if="loading">
            <div v-for="i in 6" :key="`sk${i}`" class="break-inside-avoid mb-3
                 rounded-xl bg-bg-card border border-border animate-pulse">
              <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
              <div class="p-2.5 space-y-1.5">
                <div class="h-2.5 bg-bg-hover rounded"></div>
                <div class="h-2 bg-bg-hover rounded w-2/3"></div>
              </div>
            </div>
          </template>
        </div>
        <div ref="sentinel" class="h-4 mt-2"></div>
        <p v-if="!hasMore && feedVideos.length" class="text-center py-8 text-text-muted text-xs">── 已加载全部 ──</p>
      </template>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import VideoCard from '@/components/common/VideoCard.vue'
import { useUserStore } from '@/stores/user'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import { followApi } from '@/api'

const userStore = useUserStore()
const following = ref([])
const feedVideos = ref([])
const page = ref(1)

onMounted(async () => {
  if (!userStore.isLoggedIn) return
  try {
    const res = await followApi.list()
    following.value = (res.data?.authors || []).map(u => ({
      name: u.name,
      avatar: u.avatar || '',
      hasNew: false,
    }))
  } catch {}
})

async function fetchFeed() {
  const res = await followApi.feed({ page: page.value, page_size: 20 })
  const list = res.data?.videos || []
  feedVideos.value.push(...list)
  page.value++
  if (!res.data?.has_more) return false
}

const { sentinel, loading, hasMore } = useInfiniteScroll(fetchFeed)

function scrollToCreator() {}
</script>
