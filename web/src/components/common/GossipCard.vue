<template>
  <router-link :to="`/gossip/${post.id}`"
               class="group block rounded-xl overflow-hidden bg-bg-card border border-border/60
                      hover:border-primary/30 transition-all duration-300">
    <div class="relative aspect-video overflow-hidden bg-bg-hover">
      <img v-if="post.cover" :src="post.cover" :alt="post.title"
           class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
           loading="lazy"/>
      <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-primary/20 to-bg-card">
        <svg class="w-10 h-10 text-white/20" fill="currentColor" viewBox="0 0 24 24">
          <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
        </svg>
      </div>
      <div v-if="post.tags?.length" class="absolute bottom-2 left-2 flex gap-1 flex-wrap">
        <span v-for="tag in post.tags.slice(0, 3)" :key="tag"
              class="px-1.5 py-0.5 rounded bg-black/60 text-white/80 text-[9px]">{{ tag }}</span>
      </div>
    </div>
    <div class="p-3 space-y-1.5">
      <h3 class="text-base font-semibold text-text-primary leading-snug line-clamp-2
                 group-hover:text-primary transition-colors">{{ post.title }}</h3>
      <p v-if="post.summary" class="text-sm text-text-primary/60 leading-relaxed line-clamp-2">{{ post.summary }}</p>
      <div class="flex items-center gap-3 pt-1">
        <span class="flex items-center gap-1 text-xs text-text-primary/60">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7
                     -1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
          </svg>
          {{ formatCount(post.view_count) }}
        </span>
        <span class="text-xs text-text-primary/60 ml-auto">{{ formatDate(post.published_at) }}</span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
defineProps({
  post: { type: Object, required: true },
})

function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

function formatDate(unix) {
  if (!unix) return ''
  const now = Date.now()
  const diff = now - unix * 1000
  const day = 86400000
  if (diff < day) return '今天'
  if (diff < 2 * day) return '昨天'
  if (diff < 7 * day) return Math.floor(diff / day) + '天前'
  const d = new Date(unix * 1000)
  return `${d.getMonth() + 1}/${d.getDate()}`
}
</script>
