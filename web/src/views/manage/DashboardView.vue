<template>
  <div class="space-y-5">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="card in statCards" :key="card.label"
           class="rounded-xl p-4 border border-border/60 bg-bg-surface">
        <p class="text-xs text-text-muted mb-1.5">{{ card.label }}</p>
        <p class="text-2xl font-bold text-text-primary">
          {{ loading ? '—' : fmtNum(card.value) }}
        </p>
        <p class="text-xs mt-1" :class="card.subClass">{{ card.sub }}</p>
      </div>
    </div>

    <!-- 趋势图（简单条形图） -->
    <div class="rounded-xl border border-border/60 bg-bg-surface p-5">
      <p class="text-sm font-semibold text-text-primary mb-4">近 7 天播放量</p>
      <div v-if="loading" class="h-32 flex items-end gap-2">
        <div v-for="i in 7" :key="i" class="flex-1 bg-bg-hover rounded animate-pulse" :style="`height:${30+i*10}px`"></div>
      </div>
      <div v-else-if="stats.trend?.length" class="flex items-end gap-2 h-32">
        <div v-for="item in stats.trend" :key="item.date"
             class="flex-1 flex flex-col items-center gap-1 group">
          <span class="text-[9px] text-text-muted opacity-0 group-hover:opacity-100 transition-opacity">
            {{ fmtNum(item.plays) }}
          </span>
          <div class="w-full rounded-t transition-all"
               style="background:linear-gradient(180deg,#FF2080,#8B2FC9)"
               :style="`height:${barH(item.plays)}px`"></div>
          <span class="text-[9px] text-text-muted">{{ item.date }}</span>
        </div>
      </div>
      <p v-else class="text-center text-text-muted text-sm py-10">暂无数据</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const loading = ref(true)
const stats = ref({})

const statCards = computed(() => [
  { label: '视频总数',  value: stats.value.total_videos,   sub: '全部状态',      subClass: 'text-text-muted' },
  { label: '注册用户',  value: stats.value.total_users,    sub: '累计注册',      subClass: 'text-text-muted' },
  { label: '评论总数',  value: stats.value.total_comments, sub: '正常状态',      subClass: 'text-text-muted' },
  { label: '今日播放',  value: stats.value.today_plays,    sub: '今天 00:00 起', subClass: 'text-emerald-500' },
])

function fmtNum(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function barH(plays) {
  const max = Math.max(...(stats.value.trend || []).map(t => t.plays), 1)
  return Math.max(4, Math.round((plays / max) * 100))
}

onMounted(async () => {
  try {
    const res = await manageApi.dashboard()
    stats.value = res.data || {}
  } finally {
    loading.value = false
  }
})
</script>
