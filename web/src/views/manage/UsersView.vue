<template>
  <div class="space-y-4">
    <div class="flex flex-wrap gap-2">
      <input v-model="keyword" placeholder="搜索用户名/邮箱…"
             class="px-3 py-1.5 rounded-lg bg-bg-surface border border-border text-sm text-text-primary
                    placeholder:text-text-muted focus:outline-none focus:border-primary/60 w-52"/>
      <button class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium" @click="load">搜索</button>
    </div>

    <div class="rounded-xl border border-border/60 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-bg-surface border-b border-border/60">
          <tr>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium">用户名</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium hidden md:table-cell">邮箱</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-20">角色</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-16">状态</th>
            <th class="text-left px-4 py-3 text-xs text-text-muted font-medium w-28 hidden lg:table-cell">注册时间</th>
            <th class="px-4 py-3 w-36"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="6" class="text-center py-16 text-text-muted">加载中…</td>
          </tr>
          <tr v-else-if="!users.length">
            <td colspan="6" class="text-center py-16 text-text-muted">暂无数据</td>
          </tr>
          <tr v-for="u in users" :key="u.user_id"
              class="border-b border-border/40 hover:bg-bg-hover/50 transition-colors">
            <td class="px-4 py-2.5 text-xs text-text-primary font-medium">{{ u.username }}</td>
            <td class="px-4 py-2.5 text-xs text-text-secondary hidden md:table-cell">{{ u.email || '—' }}</td>
            <td class="px-4 py-2.5">
              <span :class="['text-[11px] px-2 py-0.5 rounded-full font-medium',
                             u.role === 'admin'   ? 'bg-primary/15 text-primary' :
                             u.role === 'creator' ? 'bg-blue-500/15 text-blue-400' :
                                                    'bg-bg-hover text-text-muted']">
                {{ roleLabel(u.role) }}
              </span>
            </td>
            <td class="px-4 py-2.5">
              <span :class="['text-[11px] px-2 py-0.5 rounded-full font-medium',
                             u.status === 1 ? 'bg-emerald-500/15 text-emerald-400' : 'bg-red-500/15 text-red-400']">
                {{ u.status === 1 ? '正常' : '封禁' }}
              </span>
            </td>
            <td class="px-4 py-2.5 text-xs text-text-muted hidden lg:table-cell">{{ u.created_at }}</td>
            <td class="px-4 py-2.5">
              <div class="flex gap-1 justify-end flex-wrap">
                <button class="px-2 py-0.5 rounded text-[11px] border border-border hover:border-border-light transition-all"
                        @click="toggleStatus(u)">
                  {{ u.status === 1 ? '封禁' : '解封' }}
                </button>
                <button class="px-2 py-0.5 rounded text-[11px] border border-border hover:border-border-light transition-all"
                        @click="toggleAdmin(u)">
                  {{ u.role === 'admin' ? '取消管理员' : '设管理员' }}
                </button>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const loading = ref(true)
const users = ref([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')

async function load() {
  loading.value = true
  try {
    const res = await manageApi.userList({ page: page.value, page_size: 20, keyword: keyword.value })
    users.value = res.data?.users || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function toggleStatus(u) {
  const next = u.status === 1 ? 0 : 1
  if (!confirm(`确认${next === 0 ? '封禁' : '解封'}用户「${u.username}」？`)) return
  await manageApi.userUpdate(u.user_id, { status: next })
  u.status = next
}

async function toggleAdmin(u) {
  const next = u.role === 'admin' ? 'user' : 'admin'
  if (!confirm(`确认${next === 'admin' ? '设置' : '取消'}「${u.username}」的管理员权限？`)) return
  await manageApi.userUpdate(u.user_id, { role: next })
  u.role = next
}

function roleLabel(r) {
  return { admin: '管理员', creator: '创作者', user: '用户' }[r] || r
}

onMounted(load)
</script>
