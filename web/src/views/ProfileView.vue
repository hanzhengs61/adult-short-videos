<template>
  <div>
    <!-- 未登录 -->
    <div v-if="!userStore.isLoggedIn"
      class="flex flex-col items-center justify-center py-28 gap-5 animate-fade-in">
      <div class="w-24 h-24 rounded-full bg-bg-surface border-2 border-border flex items-center justify-center">
        <svg class="w-12 h-12 text-text-muted" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
      </div>
      <div class="text-center">
        <p class="text-text-primary font-bold text-xl mb-1">登录 mitun69</p>
        <p class="text-text-muted text-sm">管理你的内容，查看播放数据</p>
      </div>
      <button @click="userStore.openAuth('login')" class="btn-primary px-10 py-3 text-sm font-bold">立即登录</button>
      <button @click="userStore.openAuth('register')" class="text-xs text-text-muted hover:text-primary transition-colors">
        没有账号？免费注册
      </button>
    </div>

    <template v-else>
      <!-- 用户信息卡 -->
      <div class="bg-bg-surface border-b border-border">
        <div class="max-w-screen-xl mx-auto px-4 py-6">
          <div class="flex items-start gap-4">
            <!-- 头像 -->
            <div class="relative shrink-0 cursor-pointer" @click="editProfile = true; editTab = 'profile'">
              <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-primary p-0.5">
                <div class="w-full h-full rounded-full bg-bg overflow-hidden flex items-center justify-center">
                  <img v-if="userStore.userInfo?.avatar" :src="userStore.userInfo.avatar"
                    class="w-full h-full object-cover" />
                  <span v-else class="font-black text-2xl text-white">
                    {{ userStore.userInfo?.username?.[0]?.toUpperCase() }}
                  </span>
                </div>
              </div>
              <div class="absolute bottom-0 right-0 w-5 h-5 rounded-full bg-bg-surface border border-border
                          flex items-center justify-center">
                <svg class="w-3 h-3 text-text-muted" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
              </div>
            </div>

            <!-- 用户名 & Bio -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <h1 class="text-lg font-black text-text-primary">
                  {{ userStore.userInfo?.username }}
                </h1>
                <span class="px-2 py-0.5 rounded-full bg-bg-hover border border-border text-[10px] text-text-muted">
                  普通用户
                </span>
              </div>
              <p class="text-text-muted text-xs mt-1 line-clamp-2">
                {{ userStore.userInfo?.bio || '这个人很懒，什么都没写...' }}
              </p>

              <!-- 统计数据 -->
              <div class="flex items-center gap-4 mt-3">
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
                <div class="w-px h-8 bg-border"></div>
                <div class="text-center">
                  <p class="text-base font-black text-text-primary">{{ stats.following }}</p>
                  <p class="text-[10px] text-text-muted">关注</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex gap-2 mt-4">
            <button @click="editProfile = true; editTab = 'profile'"
              class="flex-1 sm:flex-none flex items-center justify-center gap-1.5 px-4 py-2
                     rounded-lg border border-border text-text-secondary text-xs font-medium
                     hover:border-border-light hover:text-text-primary transition-all">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
              编辑资料
            </button>
            <button @click="uploadModal = true"
              class="flex-1 sm:flex-none btn-primary flex items-center justify-center gap-1.5 px-5 py-2 text-xs font-bold">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
              </svg>
              上传视频
            </button>
            <button @click="showHistory"
              class="flex items-center justify-center gap-1.5 px-4 py-2 rounded-lg border border-border
                     text-text-secondary text-xs font-medium hover:border-border-light hover:text-text-primary transition-all">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              历史
            </button>
          </div>
        </div>
      </div>

      <!-- 内容 Tab -->
      <div class="max-w-screen-xl mx-auto px-3 sm:px-4">
        <div class="flex border-b border-border mb-4">
          <button v-for="tab in tabs" :key="tab.value" @click="activeTab = tab.value"
            :class="['flex items-center gap-1.5 px-4 py-3 text-sm font-semibold relative transition-colors',
              activeTab===tab.value ? 'text-text-primary' : 'text-text-muted hover:text-text-secondary']">
            {{ tab.label }}
            <span v-if="tab.count !== undefined"
              class="px-1.5 py-0.5 rounded-full bg-bg-hover text-[10px] text-text-muted">
              {{ tab.count }}
            </span>
            <span v-if="activeTab===tab.value"
              class="absolute bottom-0 left-4 right-4 h-0.5 bg-gradient-primary rounded-full"></span>
          </button>
        </div>

        <!-- 已发布 -->
        <div v-if="activeTab==='published'">
          <div v-if="!publishedVideos.length && !videosLoading"
            class="text-center py-16 text-text-muted">
            <p class="text-text-primary font-medium mb-2">还没有发布视频</p>
            <p class="text-sm mb-4">上传你的第一个视频，开始创作之旅</p>
            <button @click="uploadModal=true" class="btn-primary px-6">上传视频</button>
          </div>
          <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 pb-6">
            <VideoCard v-for="v in publishedVideos" :key="v.video_id" :video="v" />
            <div v-if="videosLoading" v-for="i in 4" :key="`sk${i}`"
              class="rounded-xl bg-bg-card border border-border animate-pulse">
              <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
              <div class="p-2.5 space-y-1.5">
                <div class="h-2.5 bg-bg-hover rounded"></div>
                <div class="h-2 bg-bg-hover rounded w-2/3"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 审核中 -->
        <div v-if="activeTab==='reviewing'">
          <div v-if="!reviewingVideos.length"
            class="text-center py-16 text-text-muted text-sm">
            没有审核中的视频
          </div>
          <div v-else class="space-y-3 pb-6">
            <div v-for="v in reviewingVideos" :key="v.video_id"
              class="flex items-center gap-3 p-3 rounded-xl bg-bg-surface border border-border">
              <div class="w-24 aspect-video rounded-lg overflow-hidden bg-bg-hover shrink-0">
                <img :src="v.cover_url" class="w-full h-full object-cover opacity-60"
                  @error="e => e.target.src='/placeholder.svg'" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-text-primary font-medium line-clamp-1">{{ v.title }}</p>
                <div class="flex items-center gap-2 mt-1.5">
                  <span class="flex items-center gap-1 px-2 py-0.5 rounded-full bg-yellow-500/10
                               border border-yellow-500/30 text-yellow-400 text-[10px] font-semibold">
                    <span class="w-1.5 h-1.5 rounded-full bg-yellow-400 animate-pulse"></span>
                    审核中
                  </span>
                  <span class="text-xs text-text-muted">预计 24h 内完成</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 收藏 -->
        <div v-if="activeTab==='favorites'">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 pb-6">
            <VideoCard v-for="v in favoriteVideos" :key="v.video_id" :video="v" />
            <div v-if="videosLoading" v-for="i in 4" :key="`fsk${i}`"
              class="rounded-xl bg-bg-card border border-border animate-pulse">
              <div class="aspect-video bg-bg-hover rounded-t-xl"></div>
              <div class="p-2.5 space-y-1.5">
                <div class="h-2.5 bg-bg-hover rounded"></div>
                <div class="h-2 bg-bg-hover rounded w-2/3"></div>
              </div>
            </div>
          </div>
          <div v-if="!favoriteVideos.length && !videosLoading" class="text-center py-16 text-text-muted text-sm">
            还没有收藏任何视频
          </div>
          <div ref="favSentinel" class="h-4"></div>
        </div>
      </div>
    </template>

    <!-- 编辑资料弹窗 -->
    <teleport to="body">
      <transition name="modal-backdrop">
        <div v-if="editProfile" class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4"
          style="background:rgba(0,0,0,0.75);backdrop-filter:blur(6px)" @click.self="editProfile=false">
          <div class="w-full sm:max-w-sm bg-bg-surface border border-border rounded-t-2xl sm:rounded-2xl shadow-2xl">
            <div class="flex justify-center pt-3 pb-1 sm:hidden">
              <div class="w-10 h-1 rounded-full bg-border-light"></div>
            </div>
            <div class="px-5 pt-4 pb-2 border-b border-border flex items-center justify-between">
              <h3 class="font-bold text-text-primary">编辑资料</h3>
              <button @click="editProfile=false" class="text-text-muted hover:text-text-primary transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>

            <!-- 编辑 Tab -->
            <div class="flex border-b border-border">
              <button v-for="t in editTabs" :key="t.value" @click="editTab = t.value"
                :class="['flex-1 py-2.5 text-xs font-semibold transition-colors relative',
                  editTab===t.value ? 'text-text-primary' : 'text-text-muted hover:text-text-secondary']">
                {{ t.label }}
                <span v-if="editTab===t.value"
                  class="absolute bottom-0 left-4 right-4 h-0.5 bg-gradient-primary rounded-full"></span>
              </button>
            </div>

            <!-- 基本资料 -->
            <div v-if="editTab==='profile'" class="px-5 py-4 space-y-4">
              <div>
                <label class="block text-xs text-text-muted mb-1.5">用户名</label>
                <input v-model="editForm.username" class="input-base" placeholder="用户名" />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">个人简介</label>
                <textarea v-model="editForm.bio" class="input-base resize-none" rows="3" placeholder="介绍一下自己..."></textarea>
              </div>
              <button @click="saveProfile" class="btn-primary w-full py-3 text-sm">保存</button>
            </div>

            <!-- 修改密码 -->
            <div v-else-if="editTab==='password'" class="px-5 py-4 space-y-4">
              <div>
                <label class="block text-xs text-text-muted mb-1.5">当前密码</label>
                <input v-model="passwordForm.current" type="password" class="input-base" placeholder="请输入当前密码" />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">新密码</label>
                <input v-model="passwordForm.newPwd" type="password" class="input-base" placeholder="至少 8 位" />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">确认新密码</label>
                <input v-model="passwordForm.confirm" type="password" class="input-base" placeholder="再次输入新密码" />
                <p v-if="passwordForm.confirm && passwordForm.newPwd !== passwordForm.confirm"
                  class="text-red-400 text-xs mt-1">两次密码不一致</p>
              </div>
              <button @click="changePassword" class="btn-primary w-full py-3 text-sm">修改密码</button>
            </div>

            <!-- 修改邮箱 -->
            <div v-else-if="editTab==='email'" class="px-5 py-4 space-y-4">
              <div v-if="userStore.userInfo?.email">
                <label class="block text-xs text-text-muted mb-1.5">当前邮箱</label>
                <input :value="userStore.userInfo.email" class="input-base opacity-60 cursor-not-allowed" disabled />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">新邮箱</label>
                <input v-model="emailForm.newEmail" type="email" class="input-base" placeholder="请输入新邮箱地址" />
              </div>
              <div>
                <label class="block text-xs text-text-muted mb-1.5">当前密码（验证身份）</label>
                <input v-model="emailForm.password" type="password" class="input-base" placeholder="请输入当前密码" />
              </div>
              <button @click="changeEmail" class="btn-primary w-full py-3 text-sm">修改邮箱</button>
            </div>
          </div>
        </div>
      </transition>
    </teleport>

    <!-- 上传视频弹窗 -->
    <teleport to="body">
      <transition name="modal-backdrop">
        <div v-if="uploadModal" class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4"
          style="background:rgba(0,0,0,0.75);backdrop-filter:blur(6px)" @click.self="uploadModal=false">
          <div class="w-full sm:max-w-lg bg-bg-surface border border-border rounded-t-2xl sm:rounded-2xl shadow-2xl max-h-[90vh] overflow-y-auto">
            <div class="flex justify-center pt-3 pb-1 sm:hidden">
              <div class="w-10 h-1 rounded-full bg-border-light"></div>
            </div>
            <div class="px-5 pt-4 pb-2 border-b border-border flex items-center justify-between">
              <h3 class="font-bold text-text-primary">上传视频</h3>
              <button @click="uploadModal=false" class="text-text-muted hover:text-text-primary transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <div class="px-5 py-4 space-y-4">
              <!-- 文件上传区 -->
              <div class="border-2 border-dashed border-border rounded-xl p-8 text-center
                          hover:border-primary/50 transition-colors cursor-pointer"
                   @click="$refs.fileInput.click()">
                <input ref="fileInput" type="file" accept="video/*" class="hidden" @change="onFileSelect" />
                <svg class="w-10 h-10 text-text-muted mx-auto mb-3" fill="none" stroke="currentColor"
                     stroke-width="1.5" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                </svg>
                <p class="text-text-secondary text-sm font-medium">{{ selectedFile ? selectedFile.name : '点击或拖拽上传视频' }}</p>
                <p class="text-text-muted text-xs mt-1">支持 MP4、AVI、MOV，最大 2GB</p>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs text-text-muted mb-1.5">视频标题 *</label>
                  <input v-model="uploadForm.title" class="input-base text-sm" placeholder="请输入标题" />
                </div>
                <div>
                  <label class="block text-xs text-text-muted mb-1.5">分类</label>
                  <select v-model="uploadForm.category" class="input-base text-sm cursor-pointer">
                    <option value="">请选择</option>
                    <option value="japan">日本</option>
                    <option value="western">欧美</option>
                    <option value="chinese">国产</option>
                    <option value="asian">亚洲</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs text-text-muted mb-1.5">标签（逗号分隔）</label>
                <input v-model="uploadForm.tags" class="input-base text-sm" placeholder="如：中文字幕, 无码, 美乳" />
              </div>

              <button @click="submitUpload" class="btn-primary w-full py-3 text-sm font-bold">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                </svg>
                提交审核
              </button>
              <p class="text-center text-xs text-text-muted">提交后将在 24 小时内完成审核</p>
            </div>
          </div>
        </div>
      </transition>
    </teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import VideoCard from '@/components/common/VideoCard.vue'
import { useUserStore } from '@/stores/user'
import { videoApi, favoriteApi, playApi } from '@/api'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('published')
const editProfile = ref(false)
const editTab = ref('profile')
const uploadModal = ref(false)
const videosLoading = ref(false)
const publishedVideos = ref([])
const reviewingVideos = ref([])
const favoriteVideos = ref([])
const selectedFile = ref(null)

const stats = ref({ videoCount: 0, totalPlays: 0, followers: 0, following: 0 })
const editForm = ref({ username: '', bio: '' })
const passwordForm = ref({ current: '', newPwd: '', confirm: '' })
const emailForm = ref({ newEmail: '', password: '' })
const uploadForm = ref({ title: '', category: '', tags: '' })

const editTabs = [
  { label: '基本资料', value: 'profile' },
  { label: '修改密码', value: 'password' },
  { label: '修改邮箱', value: 'email' },
]

const tabs = computed(() => [
  { label: '已发布', value: 'published', count: stats.value.videoCount },
  { label: '审核中', value: 'reviewing' },
  { label: '我的收藏', value: 'favorites' },
])

function formatCount(n) {
  if (!n) return '0'
  return n >= 10000 ? (n/10000).toFixed(1) + '万' : String(n)
}

async function loadTabData() {
  if (!userStore.isLoggedIn) return
  videosLoading.value = true
  try {
    if (activeTab.value === 'published') {
      const res = await videoApi.list({ page: 1, page_size: 20, order_by: 'created_at' })
      publishedVideos.value = res.data?.list || []
      stats.value.videoCount = publishedVideos.value.length
    } else if (activeTab.value === 'favorites') {
      const res = await favoriteApi.list({ page: 1, page_size: 20 })
      favoriteVideos.value = res.data?.list || []
    }
  } catch {}
  videosLoading.value = false
}

function saveProfile() {
  // TODO: PUT /api/user/profile
  if (editForm.value.username) {
    userStore.userInfo.username = editForm.value.username
  }
  editProfile.value = false
}

function changePassword() {
  const { current, newPwd, confirm } = passwordForm.value
  if (!current || !newPwd || !confirm) { alert('请填写所有密码字段'); return }
  if (newPwd.length < 8) { alert('新密码至少 8 位'); return }
  if (newPwd !== confirm) { alert('两次密码不一致'); return }
  // TODO: PUT /api/user/password
  alert('功能即将上线，敬请期待！')
  passwordForm.value = { current: '', newPwd: '', confirm: '' }
}

function changeEmail() {
  const { newEmail, password } = emailForm.value
  if (!newEmail || !password) { alert('请填写新邮箱和当前密码'); return }
  // TODO: PUT /api/user/email
  alert('功能即将上线，敬请期待！')
  emailForm.value = { newEmail: '', password: '' }
}

function onFileSelect(e) {
  selectedFile.value = e.target.files?.[0] || null
}

function submitUpload() {
  if (!uploadForm.value.title) { alert('请填写视频标题'); return }
  // TODO: POST /api/video/upload
  alert('功能即将上线，敬请期待！')
  uploadModal.value = false
}

function showHistory() {
  activeTab.value = 'favorites'
}

watch(activeTab, loadTabData)
watch(() => userStore.isLoggedIn, (v) => { if (v) loadTabData() })

onMounted(() => {
  if (userStore.isLoggedIn) {
    editForm.value.username = userStore.userInfo?.username || ''
    loadTabData()
  }
})
</script>

<style>
.modal-backdrop-enter-active, .modal-backdrop-leave-active { transition: opacity 0.25s ease; }
.modal-backdrop-enter-from, .modal-backdrop-leave-to { opacity: 0; }
</style>
