<template>
  <div class="space-y-4">
    <div class="flex gap-2 items-center">
      <input v-model="keyword" placeholder="搜索标题…"
             class="px-3 py-1.5 rounded-lg bg-bg-surface border border-border text-sm text-text-primary
                    placeholder:text-text-muted focus:outline-none focus:border-primary/60 w-52"/>
      <button class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium" @click="load">搜索</button>
      <button class="ml-auto px-3 py-1.5 rounded-lg border border-primary/50 text-primary text-xs font-medium"
              @click="openEdit(null)">+ 新增文章</button>
    </div>

    <div class="rounded-xl border border-border/60 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-bg-surface border-b border-border/60">
          <tr>
            <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium">标题</th>
            <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-16 hidden sm:table-cell">阅读</th>
            <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-14">状态</th>
            <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-28 hidden md:table-cell">发布时间</th>
            <th class="px-4 py-2.5 w-20"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading"><td colspan="5" class="text-center py-16 text-text-muted">加载中…</td></tr>
          <tr v-else-if="!posts.length"><td colspan="5" class="text-center py-16 text-text-muted">暂无数据</td></tr>
          <tr v-for="p in posts" :key="p.id" class="border-b border-border/40 hover:bg-bg-hover/50">
            <td class="px-4 py-2 text-xs text-text-primary line-clamp-1">{{ p.title }}</td>
            <td class="px-4 py-2 text-xs text-text-muted hidden sm:table-cell">{{ p.view_count }}</td>
            <td class="px-4 py-2">
              <span :class="p.status===1 ? 'text-emerald-400' : 'text-red-400'" class="text-[11px]">
                {{ p.status===1 ? '发布' : '下架' }}
              </span>
            </td>
            <td class="px-4 py-2 text-xs text-text-muted hidden md:table-cell">{{ p.published_at }}</td>
            <td class="px-4 py-2">
              <div class="flex gap-2 justify-end">
                <button class="text-xs text-text-secondary hover:text-text-primary" @click="openEdit(p)">编辑</button>
                <button class="text-xs text-red-400 hover:text-red-300" @click="del(p)">删除</button>
              </div>
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

    <!-- 编辑弹窗 -->
    <div v-if="editing" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
         @click.self="editing = null">
      <div class="w-full max-w-lg bg-bg-surface rounded-2xl border border-border p-5 space-y-3 overflow-y-auto max-h-[90vh]">
        <h3 class="text-sm font-semibold text-text-primary">{{ editing.id ? '编辑文章' : '新增文章' }}</h3>
        <div v-for="f in editFields" :key="f.key" class="space-y-1">
          <label class="text-xs text-text-muted">{{ f.label }}</label>
          <textarea v-if="f.rows" v-model="editing[f.key]" :rows="f.rows"
                    class="w-full px-3 py-2 rounded-lg bg-bg border border-border text-sm text-text-primary
                           focus:outline-none focus:border-primary/60 resize-none"/>
          <input v-else v-model="editing[f.key]" :placeholder="f.label"
                 class="w-full px-3 py-2 rounded-lg bg-bg border border-border text-sm text-text-primary
                        focus:outline-none focus:border-primary/60"/>
        </div>
        <div class="flex items-center gap-3">
          <label class="text-xs text-text-muted">状态</label>
          <button :class="['px-3 py-1 rounded-full text-xs font-medium transition-all',
                           editing.status===1 ? 'bg-emerald-500/15 text-emerald-400' : 'bg-red-500/15 text-red-400']"
                  @click="editing.status = editing.status===1 ? 0 : 1">
            {{ editing.status===1 ? '发布' : '下架' }}
          </button>
        </div>
        <div class="flex gap-2 pt-1">
          <button class="flex-1 py-2 rounded-lg bg-primary text-white text-sm font-medium" @click="save">保存</button>
          <button class="px-4 py-2 rounded-lg border border-border text-sm text-text-secondary" @click="editing=null">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const loading = ref(true)
const posts = ref([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const editing = ref(null)

const editFields = [
  { key: 'title',   label: '标题' },
  { key: 'cover',   label: '封面图 URL' },
  { key: 'summary', label: '摘要', rows: 2 },
  { key: 'content', label: '正文（HTML）', rows: 5 },
  { key: 'tags',    label: '标签（逗号分隔）' },
  { key: 'published_at', label: '发布时间（RFC3339）' },
]

function openEdit(p) {
  editing.value = p
    ? { ...p, published_at: p.published_at || new Date().toISOString() }
    : { title: '', cover: '', summary: '', content: '', tags: '', status: 1, published_at: new Date().toISOString() }
}

async function load() {
  loading.value = true
  try {
    const res = await manageApi.gossipList({ page: page.value, page_size: 20, keyword: keyword.value })
    posts.value = res.data?.posts || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function save() {
  if (editing.value.id) await manageApi.gossipUpdate(editing.value.id, editing.value)
  else await manageApi.gossipCreate(editing.value)
  editing.value = null
  load()
}

async function del(p) {
  if (!confirm(`删除「${p.title}」？`)) return
  await manageApi.gossipDelete(p.id)
  posts.value = posts.value.filter(x => x.id !== p.id)
  total.value--
}

onMounted(load)
</script>
