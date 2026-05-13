<template>
  <div class="space-y-4">
    <div class="flex flex-wrap gap-2 items-center">
      <input v-model="keyword" placeholder="搜索评论内容…"
             class="px-3 py-1.5 rounded-lg bg-bg-surface border border-border text-sm text-text-primary
                    placeholder:text-text-muted focus:outline-none focus:border-primary/60 w-52"/>
      <button class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium" @click="load">搜索</button>
      <button v-if="selected.length" class="px-3 py-1.5 rounded-lg bg-red-500/80 text-white text-xs font-medium"
              @click="batchDel">批量删除 ({{ selected.length }})</button>
    </div>

    <div class="rounded-xl border border-border/60 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-bg-surface border-b border-border/60">
          <tr>
            <th class="px-4 py-3 w-8">
              <input type="checkbox" :checked="allSelected" @change="toggleAll"
                     class="accent-primary"/>
            </th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium">用户</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium">内容</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium hidden md:table-cell">视频</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-28 hidden md:table-cell">时间</th>
            <th class="px-4 py-3 w-16"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="6" class="text-center py-16 text-text-muted">加载中…</td>
          </tr>
          <tr v-else-if="!comments.length">
            <td colspan="6" class="text-center py-16 text-text-muted">暂无数据</td>
          </tr>
          <tr v-for="c in comments" :key="c.comment_id"
              class="border-b border-border/40 hover:bg-bg-hover/50 transition-colors">
            <td class="px-4 py-2.5">
              <input type="checkbox" :value="c.comment_id" v-model="selected" class="accent-primary"/>
            </td>
            <td class="px-4 py-2.5 text-xs text-text-secondary">{{ c.username }}</td>
            <td class="px-4 py-2.5 text-xs text-text-primary max-w-xs">
              <p class="line-clamp-2">{{ c.content }}</p>
            </td>
            <td class="px-4 py-2.5 text-xs text-text-muted hidden md:table-cell line-clamp-1 max-w-[160px]">
              {{ c.video_title }}
            </td>
            <td class="px-4 py-2.5 text-xs text-text-muted hidden md:table-cell">{{ c.created_at }}</td>
            <td class="px-4 py-2.5">
              <button class="text-xs text-red-400 hover:text-red-300 transition-colors" @click="del(c)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex items-center justify-between text-xs text-text-muted">
      <span>共 {{ total }} 条</span>
      <div class="flex gap-1">
        <button :disabled="page <= 1" class="px-2.5 py-1 rounded border border-border disabled:opacity-40"
                @click="page--; load()">上一页</button>
        <span class="px-3 py-1">{{ page }}</span>
        <button :disabled="page * 20 >= total" class="px-2.5 py-1 rounded border border-border disabled:opacity-40"
                @click="page++; load()">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const loading = ref(true)
const comments = ref([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const selected = ref([])

const allSelected = computed(() =>
  comments.value.length > 0 && selected.value.length === comments.value.length
)
function toggleAll(e) {
  selected.value = e.target.checked ? comments.value.map(c => c.comment_id) : []
}

async function load() {
  loading.value = true
  selected.value = []
  try {
    const res = await manageApi.commentList({ page: page.value, page_size: 20, keyword: keyword.value })
    comments.value = res.data?.comments || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function del(c) {
  await manageApi.commentDelete(c.comment_id)
  comments.value = comments.value.filter(x => x.comment_id !== c.comment_id)
  total.value--
}

async function batchDel() {
  if (!confirm(`确认删除选中的 ${selected.value.length} 条评论？`)) return
  await manageApi.commentBatchDel(selected.value)
  comments.value = comments.value.filter(c => !selected.value.includes(c.comment_id))
  total.value -= selected.value.length
  selected.value = []
}

onMounted(load)
</script>
