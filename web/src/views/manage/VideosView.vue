<template>
  <div class="space-y-4">
    <!-- 过滤栏 -->
    <div class="flex flex-wrap gap-2 items-center">
      <input v-model="keyword" placeholder="搜索标题…"
             class="px-3 py-1.5 rounded-lg bg-bg-surface border border-border text-sm text-text-primary
                    placeholder:text-text-muted focus:outline-none focus:border-primary/60 w-48"/>
      <div class="flex gap-1">
        <button v-for="s in statusOpts" :key="s.value"
                :class="['px-3 py-1.5 rounded-lg text-xs font-medium border transition-all',
                         statusFilter === s.value
                           ? 'bg-primary/15 border-primary/40 text-primary'
                           : 'border-border text-text-secondary hover:text-text-primary']"
                @click="statusFilter = s.value; load()">
          {{ s.label }}
        </button>
      </div>
      <button class="ml-auto px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium"
              @click="load">刷新</button>
    </div>

    <!-- 表格 -->
    <div class="rounded-xl border border-border/60 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-bg-surface border-b border-border/60">
          <tr>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-16">封面</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium">标题</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-20 hidden sm:table-cell">播放</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-16">状态</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-28 hidden md:table-cell">时间</th>
            <th class="px-4 py-3 w-24"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="6" class="text-center py-16 text-text-muted">加载中…</td>
          </tr>
          <tr v-else-if="!videos.length">
            <td colspan="6" class="text-center py-16 text-text-muted">暂无数据</td>
          </tr>
          <tr v-for="v in videos" :key="v.video_id"
              class="border-b border-border/40 hover:bg-bg-hover/50 transition-colors">
            <td class="px-4 py-2.5">
              <img :src="v.cover_url" class="w-12 h-8 object-cover rounded"
                   @error="e => e.target.src='/placeholder.svg'"/>
            </td>
            <td class="px-4 py-2.5">
              <p class="text-text-primary line-clamp-1 text-xs">{{ v.title }}</p>
            </td>
            <td class="px-4 py-2.5 text-text-secondary text-xs hidden sm:table-cell">
              {{ fmtCount(v.play_count) }}
            </td>
            <td class="px-4 py-2.5">
              <span :class="['text-[11px] px-2 py-0.5 rounded-full font-medium',
                             v.status === 1 ? 'bg-emerald-500/15 text-emerald-400' : 'bg-red-500/15 text-red-400']">
                {{ v.status === 1 ? '正常' : '下架' }}
              </span>
            </td>
            <td class="px-4 py-2.5 text-text-muted text-xs hidden md:table-cell">{{ v.created_at }}</td>
            <td class="px-4 py-2.5">
              <div class="flex gap-1 justify-end">
                <button class="px-2 py-1 rounded text-xs border border-border hover:border-border-light transition-all"
                        @click="toggleStatus(v)">
                  {{ v.status === 1 ? '下架' : '上架' }}
                </button>
                <button class="px-2 py-1 rounded text-xs border border-red-500/30 text-red-400 hover:bg-red-500/10 transition-all"
                        @click="del(v)">删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex items-center justify-between text-xs text-text-muted">
      <span>共 {{ total }} 条</span>
      <div class="flex gap-1">
        <button :disabled="page <= 1" class="px-2.5 py-1 rounded border border-border disabled:opacity-40"
                @click="page--; load()">上一页</button>
        <span class="px-3 py-1">{{ page }}</span>
        <button :disabled="page * pageSize >= total" class="px-2.5 py-1 rounded border border-border disabled:opacity-40"
                @click="page++; load()">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const loading = ref(true)
const videos = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const keyword = ref('')
const statusFilter = ref(null)

const statusOpts = [
  { label: '全部', value: null },
  { label: '正常', value: 1 },
  { label: '下架', value: 0 },
]

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize, keyword: keyword.value }
    if (statusFilter.value !== null) params.status = statusFilter.value
    const res = await manageApi.videoList(params)
    videos.value = res.data?.videos || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function toggleStatus(v) {
  await manageApi.videoStatus(v.video_id, v.status === 1 ? 0 : 1)
  v.status = v.status === 1 ? 0 : 1
}

async function del(v) {
  if (!confirm(`确认删除「${v.title}」？`)) return
  await manageApi.videoDelete(v.video_id)
  videos.value = videos.value.filter(x => x.video_id !== v.video_id)
  total.value--
}

function fmtCount(n) {
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n || 0)
}

onMounted(load)
</script>
