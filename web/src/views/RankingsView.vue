<template>
  <div class="pb-8">

    <!-- ① 顶部标题 -->
    <div class="text-center py-3 px-4">
      <p class="text-sm font-semibold text-text-primary/80 tracking-wide">
        福利APP下载免费领VIP
      </p>
    </div>

    <!-- ② Banner 轮播 -->
    <div class="relative overflow-hidden mx-3 rounded-xl"
         style="height:160px"
         @touchstart.passive="onTouchStart"
         @touchend="onTouchEnd">
      <div class="flex h-full transition-transform duration-500 ease-in-out"
           :style="`transform:translateX(-${bannerIdx * 100}%)`">
        <div v-for="(b, i) in banners" :key="i"
             class="shrink-0 w-full h-full relative cursor-pointer rounded-xl overflow-hidden"
             :style="b.bg"
             @click="b.url && openUrl(b.url)">
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
          <!-- 主文案 -->
          <div class="absolute inset-0 flex flex-col items-center justify-center gap-1 px-4 text-center">
            <p class="text-2xl sm:text-3xl font-black text-white drop-shadow-lg leading-tight"
               v-html="b.title"></p>
            <p class="text-xs text-white/75">{{ b.sub }}</p>
          </div>
          <!-- 角标 -->
          <span v-if="b.tag"
                class="absolute top-2 left-2 text-[10px] font-bold px-2 py-0.5 rounded"
                :class="b.tagClass">{{ b.tag }}</span>
        </div>
      </div>
      <!-- 指示点 -->
      <div class="absolute bottom-2 left-1/2 -translate-x-1/2 flex items-center gap-1.5">
        <span v-for="(_, i) in banners" :key="i"
              class="rounded-full transition-all duration-300 cursor-pointer"
              :class="bannerIdx===i ? 'w-5 h-1.5 bg-primary' : 'w-1.5 h-1.5 bg-white/30'"
              @click="bannerIdx = i"></span>
      </div>
    </div>

    <!-- ③ 分类 Tab（横向可滚动） -->
    <div class="mt-4 px-3">
      <div class="flex gap-2 overflow-x-auto scrollbar-none pb-1">
        <button v-for="c in categories" :key="c.value"
                class="shrink-0 px-4 py-1.5 rounded-full text-sm font-semibold transition-all"
                :class="activeCategory === c.value
                  ? 'bg-primary text-white shadow-glow'
                  : 'text-text-secondary hover:text-text-primary border border-border hover:border-border-light'"
                @click="activeCategory = c.value">
          {{ c.label }}
        </button>
      </div>
    </div>

    <!-- ④ App 图标网格 -->
    <div class="mt-3 px-3 grid grid-cols-4 gap-x-3 gap-y-5">
      <div v-for="app in filteredApps" :key="app.id"
           class="flex flex-col items-center gap-2 cursor-pointer group"
           @click="openUrl(app.url)">
        <!-- App 图标 -->
        <div class="w-full aspect-square rounded-2xl overflow-hidden relative shadow-lg
                    border border-white/10 transition-transform duration-200 active:scale-90 group-hover:scale-95"
             :style="app.bg">
          <!-- 内容区 -->
          <div class="absolute inset-0 flex flex-col items-center justify-center gap-0.5 px-1">
            <span class="text-3xl leading-none">{{ app.icon }}</span>
            <span v-if="app.iconText"
                  class="text-[10px] font-black text-white/90 text-center leading-tight mt-0.5"
                  style="text-shadow:0 1px 2px rgba(0,0,0,0.8)">{{ app.iconText }}</span>
          </div>
          <!-- 官方角标 -->
          <span v-if="app.official"
                class="absolute top-1 left-1 text-[9px] font-black px-1 py-0.5 rounded
                       bg-primary text-white leading-none">官方</span>
          <!-- 热门角标 -->
          <span v-if="app.hot"
                class="absolute top-1 right-1 text-[9px] font-black px-1 py-0.5 rounded
                       bg-red-500 text-white leading-none">HOT</span>
        </div>
        <!-- 名称 -->
        <p class="text-xs text-text-secondary text-center leading-tight w-full truncate group-hover:text-text-primary transition-colors">
          {{ app.name }}
        </p>
        <!-- 下载按钮 -->
        <button class="w-full text-xs py-1 rounded-full border border-primary/60 text-primary
                       hover:bg-primary hover:text-white transition-all font-medium active:scale-95">
          立即下载
        </button>
      </div>
    </div>

    <!-- ⑤ 加载更多占位 -->
    <p class="text-center text-xs text-text-muted mt-8">── 已加载全部应用 ──</p>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

// ── Banner ────────────────────────────────────────────
const bannerIdx = ref(0)
let bannerTimer = null

const banners = [
  {
    title: '送你 <span style="color:#FFE600">1888</span> · 1元爆大奖',
    sub: 'PG电子 · 扶持放水 · PG72.COM',
    bg: 'background:linear-gradient(135deg,#0f0080 0%,#3a00cc 40%,#7B2FFF 100%)',
    tag: '赌场', tagClass: 'bg-yellow-400 text-black',
    url: '#',
  },
  {
    title: '注册送 <span style="color:#FFD700">888</span> 体验金',
    sub: '开元棋牌 · 官方正版 · 8319.COM',
    bg: 'background:linear-gradient(135deg,#8B0000 0%,#c0392b 50%,#e74c3c 100%)',
    tag: '官方', tagClass: 'bg-primary text-white',
    url: '#',
  },
  {
    title: '全国寂寞小姐姐\n等你来撩',
    sub: '约炮裸聊 · 真人在线 · 免费注册',
    bg: 'background:linear-gradient(135deg,#1a0533 0%,#6B21A8 60%,#FF2080 100%)',
    tag: '直播', tagClass: 'bg-blue-500 text-white',
    url: '#',
  },
  {
    title: '福利APP下载\n免费领 <span style="color:#FF2080">VIP</span>',
    sub: '海量内容 · 每日更新 · 无需翻墙',
    bg: 'background:linear-gradient(135deg,#003322 0%,#006644 50%,#00aa66 100%)',
    tag: '免费', tagClass: 'bg-emerald-500 text-white',
    url: '#',
  },
]

function startBanner() {
  bannerTimer = setInterval(() => {
    bannerIdx.value = (bannerIdx.value + 1) % banners.length
  }, 4000)
}

let touchX = 0
function onTouchStart(e) { touchX = e.touches[0].clientX }
function onTouchEnd(e) {
  const dx = e.changedTouches[0].clientX - touchX
  if (Math.abs(dx) < 40) return
  bannerIdx.value = dx < 0
    ? (bannerIdx.value + 1) % banners.length
    : (bannerIdx.value - 1 + banners.length) % banners.length
}

// ── 分类 Tab ──────────────────────────────────────────
const activeCategory = ref('all')

const categories = [
  { label: '推荐',   value: 'all'   },
  { label: '国产看片', value: 'video' },
  { label: '约炮裸聊', value: 'chat'  },
  { label: '男同',   value: 'gay'   },
  { label: '美女直播', value: 'live'  },
  { label: '成人漫画', value: 'manga' },
  { label: '电子游戏', value: 'game'  },
]

// ── App 数据 ──────────────────────────────────────────
// 每条数据: id / name / icon(emoji) / iconText / bg(gradient) / official / hot / categories / url
const apps = [
  {
    id: 1, name: '巨乳约炮',
    icon: '🔞', iconText: null,
    bg: 'background:linear-gradient(135deg,#1a0010,#7B2FFF)',
    official: false, hot: true,
    cats: ['all', 'chat'],
    url: '#',
  },
  {
    id: 2, name: '激情迷药',
    icon: '💊', iconText: '情趣',
    bg: 'background:linear-gradient(135deg,#003366,#0066cc)',
    official: false, hot: false,
    cats: ['all', 'chat'],
    url: '#',
  },
  {
    id: 3, name: 'PG电子',
    icon: '🎰', iconText: 'PG',
    bg: 'background:linear-gradient(135deg,#0f0080,#7B2FFF)',
    official: true, hot: true,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 4, name: '棋牌电子',
    icon: '🃏', iconText: '888',
    bg: 'background:linear-gradient(135deg,#003399,#0055cc)',
    official: true, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 5, name: '3223棋牌',
    icon: '♠️', iconText: 'KY',
    bg: 'background:linear-gradient(135deg,#8B0000,#cc0000)',
    official: true, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 6, name: '开元棋牌',
    icon: '🎲', iconText: 'KY',
    bg: 'background:linear-gradient(135deg,#333300,#999900)',
    official: true, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 7, name: 'PG放水',
    icon: '💸', iconText: '大放水',
    bg: 'background:linear-gradient(135deg,#000000,#1a1a1a)',
    official: false, hot: true,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 8, name: '澳门赌场',
    icon: '🎪', iconText: '翻倍',
    bg: 'background:linear-gradient(135deg,#660000,#cc3300)',
    official: false, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 9, name: '77直播',
    icon: '📺', iconText: '77',
    bg: 'background:linear-gradient(135deg,#8B0000,#FF2080)',
    official: false, hot: true,
    cats: ['all', 'live'],
    url: '#',
  },
  {
    id: 10, name: '澳门新葡京',
    icon: '🏯', iconText: '3A',
    bg: 'background:linear-gradient(135deg,#1a1a00,#4a4a00)',
    official: true, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 11, name: '金沙直播',
    icon: '🌟', iconText: '金沙',
    bg: 'background:linear-gradient(135deg,#660033,#cc0066)',
    official: false, hot: false,
    cats: ['all', 'live'],
    url: '#',
  },
  {
    id: 12, name: '国产精品',
    icon: '🎬', iconText: null,
    bg: 'background:linear-gradient(135deg,#003300,#006600)',
    official: false, hot: true,
    cats: ['all', 'video'],
    url: '#',
  },
  {
    id: 13, name: '91约炮',
    icon: '💬', iconText: '约炮',
    bg: 'background:linear-gradient(135deg,#1a0033,#4a0099)',
    official: false, hot: true,
    cats: ['all', 'chat'],
    url: '#',
  },
  {
    id: 14, name: '男同专区',
    icon: '👬', iconText: null,
    bg: 'background:linear-gradient(135deg,#003366,#006699)',
    official: false, hot: false,
    cats: ['all', 'gay'],
    url: '#',
  },
  {
    id: 15, name: '成人漫画',
    icon: '📖', iconText: 'H漫',
    bg: 'background:linear-gradient(135deg,#330033,#660066)',
    official: false, hot: false,
    cats: ['all', 'manga'],
    url: '#',
  },
  {
    id: 16, name: '美女直播',
    icon: '💃', iconText: null,
    bg: 'background:linear-gradient(135deg,#4a0020,#FF2080)',
    official: false, hot: true,
    cats: ['all', 'live'],
    url: '#',
  },
  {
    id: 17, name: '麻豆传媒',
    icon: '🎭', iconText: '麻豆',
    bg: 'background:linear-gradient(135deg,#1a0000,#660000)',
    official: true, hot: true,
    cats: ['all', 'video'],
    url: '#',
  },
  {
    id: 18, name: '新葡京',
    icon: '🎯', iconText: '024',
    bg: 'background:linear-gradient(135deg,#660000,#990000)',
    official: false, hot: false,
    cats: ['all', 'game'],
    url: '#',
  },
  {
    id: 19, name: '裸聊平台',
    icon: '📱', iconText: '裸聊',
    bg: 'background:linear-gradient(135deg,#001a33,#0066cc)',
    official: false, hot: false,
    cats: ['all', 'chat'],
    url: '#',
  },
  {
    id: 20, name: '同志交友',
    icon: '🌈', iconText: null,
    bg: 'background:linear-gradient(135deg,#003344,#006688)',
    official: false, hot: false,
    cats: ['gay'],
    url: '#',
  },
]

const filteredApps = computed(() =>
  activeCategory.value === 'all'
    ? apps
    : apps.filter(a => a.cats.includes(activeCategory.value))
)

function openUrl(url) {
  // 实际接入时替换为真实链接
  if (url && url !== '#') window.open(url, '_blank')
}

onMounted(startBanner)
onUnmounted(() => clearInterval(bannerTimer))
</script>
