<template>
  <div>
    <div v-if="!userStore.isLoggedIn"
         class="flex flex-col items-center justify-center py-28 gap-5 animate-fade-in">
      <div class="w-24 h-24 rounded-full bg-bg-surface border-2 border-border flex items-center justify-center">
        <svg class="w-12 h-12 text-text-muted" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
      </div>
      <div class="text-center">
        <p class="text-text-primary font-bold text-xl mb-1">登录 mitun69</p>
        <p class="text-text-muted text-sm">管理你的内容，查看播放数据</p>
      </div>
      <button @click="userStore.openAuth('login')" class="btn-primary px-10 py-3 text-sm font-bold">立即登录</button>
      <button @click="userStore.openAuth('register')" class="text-xs text-text-muted hover:text-primary transition-colors">没有账号？免费注册</button>
    </div>

    <template v-else>
      <div class="bg-bg-surface border-b border-border">
        <div class="max-w-screen-xl mx-auto px-4 py-6">
          <div class="flex items-start gap-4">
            <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-primary p-0.5 shrink-0">
              <div class="w-full h-full rounded-full bg-bg overflow-hidden">
                <img :src="userStore.userInfo?.avatar || '/logo.png'" :alt="userStore.userInfo?.username" class="w-full h-full object-cover"/>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <h1 class="text-lg font-black text-text-primary">{{ userStore.userInfo?.username }}</h1>
                <span :class="['px-2 py-0.5 rounded-full text-[10px] font-semibold border',
                  isCreator
                    ? 'bg-primary/10 border-primary/30 text-primary'
                    : 'bg-bg-hover border-border text-text-muted']">
                  {{ isCreator ? '创作者' : '普通用户' }}
                </span>
              </div>
              <p class="text-text-muted text-xs mt-1 line-clamp-2">
                {{ userStore.userInfo?.bio || '这个人很懒，什么都没写...' }}
              </p>
              <div v-if="isCreator" class="flex items-center gap-4 mt-3">
                <div class="text-center">
                  <p class="text-base font-black text-text-primary">{{ stats.videoCount }}</p>
                  <p class="text-[10px] text-text-muted">视频</p>
                </div>
                <div class="w-px h-8 bg-border"></div>
                <div class="text-center">
                  <p class="text-base font-black text-text-primary">{{ formatCount(stats.totalPlays) }}</p>
                  <p class="text-[10px] text-text-muted">累计播放</p>
                </div>
                <div class="w-px h-8 bg-border"></div>
                <div class="text-center">
                  <p class="text-base font-black text-text-primary">{{ stats.followers }}</p>
                  <p class="text-[10px] text-text-muted">粉丝</p>
                </div>
              </div>
            </div>
          </div>
          <div class="flex gap-2 mt-4">
            <button @click="openEdit"
                    class="flex-1 sm:flex-none flex items-center justify-center gap-1.5 px-4 py-2 rounded-lg
                     border border-border text-text-secondary text-xs font-medium
                     hover:border-border-light hover:text-text-primary transition-all">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
              编辑信息
            </button>
            <button v-if="isCreator" @click="uploadModal = true"
                    class="flex-1 sm:flex-none btn-primary flex items-center justify-center gap-1.5 px-5 py-2 text-xs font-bold">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
              </svg>
              上传视频
            </button>
          </div>
        </div>
      </div>

      <div v-if="!isCreator" class="max-w-screen-xl mx-auto px-4 py-5 space-y-2">
        <component :is="item.coming ? 'div' : 'router-link'"
                   v-for="item in menuItems" :key="item.to"
                   v-bind="item.coming ? {} : {to: item.to}"
                   :class="['flex items-center justify-between px-4 py-3.5 rounded-xl bg-bg-card border transition-all duration-200 group',
                     item.coming
                       ? 'border-border opacity-60 cursor-default'
                       : 'border-border hover:border-primary/30 hover:bg-bg-hover cursor-pointer']">
          <div class="flex items-center gap-3">
            <div :class="['w-8 h-8 rounded-lg border flex items-center justify-center transition-all duration-200',
                          item.coming
                            ? 'bg-bg-hover border-border'
                            : 'bg-primary/10 border-primary/20 group-hover:bg-primary/20 group-hover:scale-105']">
              <svg :class="['w-4 h-4', item.coming ? 'text-text-muted' : 'text-primary']"
                   fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" :d="item.icon"/>
              </svg>
            </div>
            <span class="text-sm text-text-primary font-medium">{{ item.label }}</span>
          </div>
          <span v-if="item.coming"
                class="text-[10px] px-2 py-0.5 rounded-full bg-bg-hover border border-border text-text-muted">
            即将上线
          </span>
          <svg v-else class="w-4 h-4 text-text-muted group-hover:text-primary transition-colors"
               fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
          </svg>
        </component>

        <button @click="userStore.logout()"
                class="w-full flex items-center gap-3 px-4 py-3.5 rounded-xl
                       border border-red-500/20 bg-red-500/5 text-red-400
                       hover:bg-red-500/10 hover:border-red-400/40 transition-all duration-200 group">
          <div class="w-8 h-8 rounded-lg bg-red-500/10 flex items-center justify-center
                      group-hover:bg-red-500/20 group-hover:scale-105 transition-all duration-200">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
          </div>
          <span class="text-sm font-medium">退出登录</span>
        </button>
      </div>

      <!-- 创作者：内容 Tab -->
      <div v-else class="max-w-screen-xl mx-auto px-3 sm:px-4">
        <div class="flex border-b border-border mb-4 overflow-x-auto scrollbar-none">
          <button v-for="tab in creatorTabs" :key="tab.value" @click="activeTab = tab.value"
                  :class="['flex items-center gap-1.5 px-4 py-3 text-sm font-semibold relative shrink-0 transition-colors',
              activeTab === tab.value ? 'text-text-primary' : 'text-text-muted hover:text-text-secondary']">
            {{ tab.label }}
            <span class="px-1.5 py-0.5 rounded-full bg-bg-hover text-[10px] text-text-muted">{{ tab.count }}</span>
            <span v-if="activeTab === tab.value"
                  class="absolute bottom-0 left-4 right-4 h-0.5 bg-gradient-primary rounded-full"></span>
          </button>
        </div>
        <div class="columns-2 sm:columns-3 md:columns-4 gap-3 pb-6">
          <div v-for="v in currentTabVideos" :key="v.video_id" class="break-inside-avoid mb-3">
            <VideoCard :video="v"/>
          </div>
          <template v-if="videosLoading">
            <div v-for="i in 4" :key="`sk${i}`"
                 class="break-inside-avoid mb-3 rounded-xl bg-bg-card border border-border animate-pulse">
              <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
              <div class="p-2.5 space-y-1.5">
                <div class="h-2.5 bg-bg-hover rounded"></div>
                <div class="h-2 bg-bg-hover rounded w-2/3"></div>
              </div>
            </div>
          </template>
        </div>
        <div v-if="!currentTabVideos.length && !videosLoading"
             class="text-center py-16 text-text-muted text-sm">暂无内容</div>
      </div>
    </template>

    <teleport to="body">
      <transition name="modal-backdrop">
        <div v-if="editModal" class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4"
             style="background:rgba(0,0,0,0.75);backdrop-filter:blur(6px)" @click.self="editModal=false">
          <div class="w-full sm:max-w-sm bg-bg-surface border border-border rounded-t-2xl sm:rounded-2xl shadow-2xl">
            <div class="flex justify-center pt-3 pb-1 sm:hidden">
              <div class="w-10 h-1 rounded-full bg-border-light"></div>
            </div>
            <div class="px-5 pt-4 pb-3 border-b border-border flex items-center justify-between">
              <h3 class="font-bold text-text-primary">编辑信息</h3>
              <button @click="editModal=false" class="text-text-muted hover:text-text-primary transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <div class="px-5 py-4 space-y-4">
              <!-- 头像上传 -->
              <div class="flex flex-col items-center gap-2 pb-4 border-b border-border">
                <div class="relative">
                  <div class="w-20 h-20 rounded-full bg-gradient-primary p-0.5">
                    <div class="w-full h-full rounded-full bg-bg overflow-hidden">
                      <img :src="avatarPreview || userStore.userInfo?.avatar || '/logo.png'"
                           class="w-full h-full object-cover"/>
                    </div>
                  </div>
                  <button @click="canChangeAvatar && $refs.avatarInput.click()"
                          :class="['absolute -bottom-0.5 -right-0.5 w-7 h-7 rounded-full border-2 border-bg-surface flex items-center justify-center text-white transition-all',
                            canChangeAvatar ? 'bg-primary hover:bg-primary-hover cursor-pointer' : 'bg-bg-hover cursor-not-allowed']">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round"
                            d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                  </button>
                  <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="onAvatarSelect"/>
                </div>
                <p v-if="!canChangeAvatar" class="text-[11px] text-amber-500/80">
                  头像 {{ daysUntilChange('avatar') }} 天后可再次修改
                </p>
                <p v-else class="text-[11px] text-text-muted">点击相机图标更换头像</p>
              </div>

              <!-- 昵称 -->
              <div>
                <div class="flex items-center justify-between mb-1.5">
                  <label class="text-xs text-text-muted">昵称</label>
                  <span class="text-[10px]" :class="canChangeUsername ? 'text-text-muted' : 'text-amber-500/80'">
                    {{ canChangeUsername ? '每30天可修改一次' : `${daysUntilChange('username')} 天后可修改` }}
                  </span>
                </div>
                <input v-model="editForm.username" class="input-base"
                       :class="{'opacity-50 cursor-not-allowed': !canChangeUsername}"
                       :disabled="!canChangeUsername" placeholder="请输入昵称"/>
              </div>

              <!-- 个人简介 -->
              <div>
                <label class="block text-xs text-text-muted mb-1.5">个人简介</label>
                <textarea v-model="editForm.bio" class="input-base resize-none" rows="3"
                          placeholder="介绍一下自己..."></textarea>
              </div>

              <button @click="saveProfile" class="btn-primary w-full py-3 text-sm">保存</button>
            </div>
          </div>
        </div>
      </transition>
    </teleport>
    <teleport to="body">
      <transition name="modal-backdrop">
        <div v-if="uploadModal" class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4"
             style="background:rgba(0,0,0,0.75);backdrop-filter:blur(6px)" @click.self="uploadModal=false">
          <div class="w-full sm:max-w-lg bg-bg-surface border border-border rounded-t-2xl sm:rounded-2xl shadow-2xl
                      max-h-[90vh] overflow-y-auto">
            <div class="flex justify-center pt-3 pb-1 sm:hidden">
              <div class="w-10 h-1 rounded-full bg-border-light"></div>
            </div>
            <div class="px-5 pt-4 pb-3 border-b border-border flex items-center justify-between">
              <h3 class="font-bold text-text-primary">上传视频</h3>
              <button @click="uploadModal=false" class="text-text-muted hover:text-text-primary transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <div class="px-5 py-4 space-y-4">
              <div class="border-2 border-dashed border-border rounded-xl p-8 text-center
                          hover:border-primary/50 transition-colors cursor-pointer"
                   @click="$refs.fileInput.click()">
                <input ref="fileInput" type="file" accept="video/*" class="hidden" @change="onFileSelect"/>
                <svg class="w-10 h-10 text-text-muted mx-auto mb-3" fill="none" stroke="currentColor"
                     stroke-width="1.5" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round"
                        d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                </svg>
                <p class="text-text-secondary text-sm font-medium">
                  {{ selectedFile ? selectedFile.name : '点击或拖拽上传视频' }}
                </p>
                <p class="text-text-muted text-xs mt-1">支持 MP4、AVI、MOV，最大 2GB</p>
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">视频标题 *</label>
                <input v-model="uploadForm.title" class="input-base text-sm" placeholder="请输入标题"/>
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">标签（逗号分隔）</label>
                <input v-model="uploadForm.tags" class="input-base text-sm" placeholder="如：丝袜,巨乳,御姐"/>
              </div>
              <button @click="submitUpload" class="btn-primary w-full py-3 text-sm font-bold">提交审核</button>
              <p class="text-center text-xs text-text-muted">提交后将在 24 小时内完成审核</p>
            </div>
          </div>
        </div>
      </transition>
    </teleport>
  </div>
</template>

<script setup>
import {computed, onMounted, ref, watch} from 'vue'
import VideoCard from '@/components/common/VideoCard.vue'
import {useUserStore} from '@/stores/user'
import {useToastStore} from '@/stores/toast'
import {userApi, videoApi} from '@/api'

const userStore = useUserStore()
const toastStore = useToastStore()
const activeTab = ref('published')
const editModal = ref(false), uploadModal = ref(false), videosLoading = ref(false)
const selectedFile = ref(null), avatarPreview = ref(null)
const publishedVideos = ref([]), reviewingVideos = ref([])
const uploadingVideos = ref([]), delistedVideos = ref([])
const stats = ref({videoCount: 0, totalPlays: 0, followers: 0})
const editForm = ref({username: '', bio: ''})
const uploadForm = ref({title: '', tags: ''})
const isCreator = computed(() => userStore.userInfo?.role === 'creator')

const CHANGE_INTERVAL_MS = 30 * 24 * 60 * 60 * 1000
// 头像/昵称冷却期均以服务端时间戳（Unix 秒）为准，null 表示从未修改过
const canChangeAvatar = computed(() => {
  const t = userStore.userInfo?.avatar_changed_at
  return !t || Date.now() - t * 1000 > CHANGE_INTERVAL_MS
})
const canChangeUsername = computed(() => {
  const t = userStore.userInfo?.username_changed_at
  return !t || Date.now() - t * 1000 > CHANGE_INTERVAL_MS
})
function daysUntilChange(key) {
  const t = key === 'avatar'
    ? (userStore.userInfo?.avatar_changed_at || 0)
    : (userStore.userInfo?.username_changed_at || 0)
  return Math.ceil((CHANGE_INTERVAL_MS - (Date.now() - t * 1000)) / (24 * 60 * 60 * 1000))
}

const menuItems = [
  {label: '观看记录', to: '/history', coming: true, icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z'},
  {label: '我的收藏', to: '/favorites', coming: true, icon: 'M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z'},
  {label: '我的互动', to: '/interactions', coming: true, icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z'},
]

const creatorTabs = computed(() => [
  {label: '已发布', value: 'published', count: publishedVideos.value.length},
  {label: '审核中', value: 'reviewing', count: reviewingVideos.value.length},
  {label: '上传中', value: 'uploading', count: uploadingVideos.value.length},
  {label: '已下架', value: 'delisted', count: delistedVideos.value.length},
])

const currentTabVideos = computed(() => {
  const map = {published: publishedVideos, reviewing: reviewingVideos, uploading: uploadingVideos, delisted: delistedVideos}
  return map[activeTab.value]?.value || []
})

function formatCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n / 10000).toFixed(1) + '万' : String(n)
}

function openEdit() {
  editForm.value = {username: userStore.userInfo?.username || '', bio: userStore.userInfo?.bio || ''}
  avatarPreview.value = null
  editModal.value = true
}

async function onAvatarSelect(e) {
  const file = e.target.files?.[0]
  if (!file) return
  // 本地预览（立即显示，不等上传完成）
  const reader = new FileReader()
  reader.onload = (ev) => { avatarPreview.value = ev.target.result }
  reader.readAsDataURL(file)
  // 上传到服务端，更新数据库
  try {
    await userApi.uploadAvatar(file)
    await userStore.fetchInfo() // 刷新 avatar + avatar_changed_at，触发冷却期计算
    avatarPreview.value = null
    toastStore.show('头像已更新', 'success')
  } catch (err) {
    avatarPreview.value = null
    toastStore.show(err?.response?.data?.message || '头像上传失败', 'error')
  }
}

async function saveProfile() {
  try {
    await userApi.updateProfile({
      username: canChangeUsername.value ? editForm.value.username : '',
      bio: editForm.value.bio,
    })
    await userStore.fetchInfo() // 刷新 username_changed_at 等服务端字段
    toastStore.show('保存成功', 'success')
    editModal.value = false
  } catch (e) {
    toastStore.show(e?.response?.data?.message || '保存失败', 'error')
  }
}

function onFileSelect(e) {
  selectedFile.value = e.target.files?.[0] || null
}

function submitUpload() {
  if (!uploadForm.value.title) { toastStore.show('请填写视频标题', 'error'); return }
  // TODO: POST /api/video/upload
  toastStore.show('功能即将上线，敬请期待！', 'info')
  uploadModal.value = false
}

async function loadTabData() {
  if (!userStore.isLoggedIn || !isCreator.value) return
  videosLoading.value = true
  try {
    if (activeTab.value === 'published') {
      const res = await videoApi.list({page: 1, page_size: 20, order_by: 'created_at'})
      publishedVideos.value = res.data?.videos || []
      stats.value.videoCount = publishedVideos.value.length
    }
    // 审核中、上传中、已下架待后端接口支持
  } catch {}
  videosLoading.value = false
}

watch(activeTab, loadTabData)
watch(() => userStore.isLoggedIn, (v) => { if (v) loadTabData() })

onMounted(() => { if (userStore.isLoggedIn) loadTabData() })
</script>

<style>
.modal-backdrop-enter-active, .modal-backdrop-leave-active { transition: opacity 0.25s ease }
.modal-backdrop-enter-from, .modal-backdrop-leave-to { opacity: 0 }
</style>
