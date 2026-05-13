<template>
  <div class="space-y-8">

    <!-- Banner 管理 -->
    <section>
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-sm font-semibold text-text-primary">轮播 Banner</h2>
        <button class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium"
                @click="openBanner(null)">+ 新增</button>
      </div>
      <div class="rounded-xl border border-border/60 overflow-hidden">
        <table class="w-full text-sm">
          <thead class="bg-bg-surface border-b border-border/60">
            <tr>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium">标题</th>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-12 hidden sm:table-cell">排序</th>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-14">状态</th>
              <th class="px-4 py-2.5 w-20"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in banners" :key="b.id" class="border-b border-border/40 hover:bg-bg-hover/50">
              <td class="px-4 py-2 text-xs text-text-primary" v-html="b.title"></td>
              <td class="px-4 py-2 text-xs text-text-muted hidden sm:table-cell">{{ b.sort }}</td>
              <td class="px-4 py-2">
                <span :class="b.status===1 ? 'text-emerald-400' : 'text-red-400'" class="text-[11px]">
                  {{ b.status===1 ? '上线' : '下线' }}
                </span>
              </td>
              <td class="px-4 py-2">
                <div class="flex gap-2 justify-end">
                  <button class="text-xs text-text-secondary hover:text-text-primary" @click="openBanner(b)">编辑</button>
                  <button class="text-xs text-red-400 hover:text-red-300" @click="delBanner(b)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- App 广告管理 -->
    <section>
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-sm font-semibold text-text-primary">App 广告位</h2>
        <button class="px-3 py-1.5 rounded-lg bg-primary text-white text-xs font-medium"
                @click="openApp(null)">+ 新增</button>
      </div>
      <div class="rounded-xl border border-border/60 overflow-hidden">
        <table class="w-full text-sm">
          <thead class="bg-bg-surface border-b border-border/60">
            <tr>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium">名称</th>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium hidden sm:table-cell">分类</th>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-12 hidden sm:table-cell">排序</th>
              <th class="text-left px-4 py-2.5 text-xs text-text-muted font-medium w-14">状态</th>
              <th class="px-4 py-2.5 w-20"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in apps" :key="a.id" class="border-b border-border/40 hover:bg-bg-hover/50">
              <td class="px-4 py-2 text-xs text-text-primary">{{ a.icon }} {{ a.name }}</td>
              <td class="px-4 py-2 text-xs text-text-muted hidden sm:table-cell">{{ a.categories }}</td>
              <td class="px-4 py-2 text-xs text-text-muted hidden sm:table-cell">{{ a.sort }}</td>
              <td class="px-4 py-2">
                <span :class="a.status===1 ? 'text-emerald-400' : 'text-red-400'" class="text-[11px]">
                  {{ a.status===1 ? '上线' : '下线' }}
                </span>
              </td>
              <td class="px-4 py-2">
                <div class="flex gap-2 justify-end">
                  <button class="text-xs text-text-secondary hover:text-text-primary" @click="openApp(a)">编辑</button>
                  <button class="text-xs text-red-400 hover:text-red-300" @click="delApp(a)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- 编辑弹窗 (通用) -->
    <div v-if="modal.show"
         class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
         @click.self="modal.show = false">
      <div class="w-full max-w-md bg-bg-surface rounded-2xl border border-border p-5 space-y-3 mx-4">
        <h3 class="text-sm font-semibold text-text-primary">{{ modal.title }}</h3>
        <div v-for="f in modal.fields" :key="f.key" class="space-y-1">
          <label class="text-xs text-text-muted">{{ f.label }}</label>
          <input v-model="modal.data[f.key]" :placeholder="f.placeholder || f.label"
                 class="w-full px-3 py-2 rounded-lg bg-bg border border-border text-sm text-text-primary
                        focus:outline-none focus:border-primary/60 placeholder:text-text-muted"/>
        </div>
        <!-- 状态开关 -->
        <div class="flex items-center gap-3">
          <label class="text-xs text-text-muted">状态</label>
          <button :class="['px-3 py-1 rounded-full text-xs font-medium transition-all',
                           modal.data.status===1 ? 'bg-emerald-500/15 text-emerald-400' : 'bg-red-500/15 text-red-400']"
                  @click="modal.data.status = modal.data.status===1 ? 0 : 1">
            {{ modal.data.status===1 ? '上线' : '下线' }}
          </button>
        </div>
        <div class="flex gap-2 pt-1">
          <button class="flex-1 py-2 rounded-lg bg-primary text-white text-sm font-medium" @click="saveModal">保存</button>
          <button class="px-4 py-2 rounded-lg border border-border text-sm text-text-secondary" @click="modal.show=false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { manageApi } from '@/api/manage'

const banners = ref([])
const apps = ref([])

const modal = reactive({ show: false, type: '', title: '', fields: [], data: {}, id: null })

const bannerFields = [
  { key: 'title', label: '标题（支持HTML）' },
  { key: 'sub_title', label: '副标题' },
  { key: 'bg_style', label: '背景CSS', placeholder: 'background:linear-gradient(...)' },
  { key: 'tag', label: '角标文字' },
  { key: 'tag_class', label: '角标CSS类' },
  { key: 'url', label: '跳转链接' },
  { key: 'sort', label: '排序权重（越大越靠前）' },
]
const appFields = [
  { key: 'name', label: '应用名称' },
  { key: 'icon', label: 'Emoji图标' },
  { key: 'icon_text', label: '图标文字（可选）' },
  { key: 'bg_style', label: '背景CSS', placeholder: 'background:linear-gradient(...)' },
  { key: 'categories', label: '分类（逗号分隔）', placeholder: 'all,video,chat' },
  { key: 'url', label: '跳转链接' },
  { key: 'sort', label: '排序权重' },
]

function openBanner(b) {
  modal.type = 'banner'
  modal.title = b ? '编辑 Banner' : '新增 Banner'
  modal.fields = bannerFields
  modal.id = b?.id || null
  modal.data = b ? { ...b } : { title: '', sub_title: '', bg_style: '', tag: '', tag_class: '', url: '', sort: 0, status: 1 }
  modal.show = true
}

function openApp(a) {
  modal.type = 'app'
  modal.title = a ? '编辑 App 广告' : '新增 App 广告'
  modal.fields = appFields
  modal.id = a?.id || null
  modal.data = a ? { ...a } : { name: '', icon: '', icon_text: '', bg_style: '', official: false, hot: false, categories: 'all', url: '', sort: 0, status: 1 }
  modal.show = true
}

async function saveModal() {
  if (modal.type === 'banner') {
    if (modal.id) await manageApi.bannerUpdate(modal.id, modal.data)
    else await manageApi.bannerCreate(modal.data)
    await loadBanners()
  } else {
    if (modal.id) await manageApi.appUpdate(modal.id, modal.data)
    else await manageApi.appCreate(modal.data)
    await loadApps()
  }
  modal.show = false
}

async function delBanner(b) {
  if (!confirm(`删除 Banner「${b.title.replace(/<[^>]+>/g, '')}」？`)) return
  await manageApi.bannerDelete(b.id)
  banners.value = banners.value.filter(x => x.id !== b.id)
}

async function delApp(a) {
  if (!confirm(`删除广告「${a.name}」？`)) return
  await manageApi.appDelete(a.id)
  apps.value = apps.value.filter(x => x.id !== a.id)
}

async function loadBanners() {
  const res = await manageApi.bannerList()
  banners.value = res.data?.banners || []
}
async function loadApps() {
  const res = await manageApi.appList()
  apps.value = res.data?.apps || []
}

onMounted(() => { loadBanners(); loadApps() })
</script>
