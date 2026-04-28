<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-1.5 py-8">
    <button @click="emit('change', current - 1)" :disabled="current <= 1"
            class="w-9 h-9 rounded-lg border border-border text-text-secondary hover:border-primary hover:text-primary
             disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center justify-center">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
      </svg>
    </button>

    <template v-for="p in pages" :key="p">
      <span v-if="p === '...'" class="w-9 h-9 flex items-center justify-center text-text-muted text-sm">…</span>
      <button v-else @click="emit('change', p)"
              :class="['w-9 h-9 rounded-lg text-sm font-medium border transition-all',
          p === current
            ? 'bg-gradient-primary border-transparent text-white'
            : 'border-border text-text-secondary hover:border-primary hover:text-primary']">
        {{ p }}
      </button>
    </template>

    <button @click="emit('change', current + 1)" :disabled="current >= totalPages"
            class="w-9 h-9 rounded-lg border border-border text-text-secondary hover:border-primary hover:text-primary
             disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center justify-center">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
      </svg>
    </button>
  </div>
</template>

<script setup>
import {computed} from 'vue'

const props = defineProps({
  current: {type: Number, default: 1},
  totalPages: {type: Number, default: 1},
})
const emit = defineEmits(['change'])

const pages = computed(() => {
  const total = props.totalPages, cur = props.current
  if (total <= 7) return Array.from({length: total}, (_, i) => i + 1)
  const arr = [1]
  if (cur > 3) arr.push('...')
  for (let i = Math.max(2, cur - 1); i <= Math.min(total - 1, cur + 1); i++) arr.push(i)
  if (cur < total - 2) arr.push('...')
  arr.push(total)
  return arr
})
</script>
