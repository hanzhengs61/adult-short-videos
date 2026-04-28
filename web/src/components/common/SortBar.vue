<template>
  <div class="flex items-center gap-1.5 flex-wrap">
    <button v-for="s in sorts" :key="s.value" @click="emit('update:modelValue', s.value)"
            :class="['shrink-0 flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all',
              modelValue === s.value
                ? 'bg-gradient-primary border-transparent text-white shadow-glow'
                : 'border-border bg-bg-surface text-text-muted hover:text-text-primary hover:border-border-light']">
      <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2"
           stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
        <path :d="s.iconPath"/>
      </svg>
      {{ s.label }}
    </button>
  </div>
</template>

<script setup>
defineProps({ modelValue: { type: String, default: 'created_at' } })
const emit = defineEmits(['update:modelValue'])

const sorts = [
  { label: '最新',   value: 'created_at',    iconPath: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' },
  { label: '播放最多', value: 'play_count',    iconPath: 'M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z' },
  { label: '评论最多', value: 'comment_count', iconPath: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z' },
  { label: '收藏最多', value: 'favorite_count', iconPath: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z' },
]
</script>
