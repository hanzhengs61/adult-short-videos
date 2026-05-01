<template>
  <div>
    <h2 class="section-title mb-4">
      <span class="w-1 h-4 bg-gradient-primary rounded-full"></span>
      评论 <span class="text-text-muted font-normal text-sm ml-1">{{ total }}</span>
    </h2>

    <!-- 评论列表 -->
    <div class="space-y-4">
      <div v-for="c in comments" :key="c.comment_id"
           class="flex gap-3 p-4 rounded-xl bg-bg-surface border border-border/50">
        <div class="w-8 h-8 rounded-full bg-gradient-primary flex items-center justify-center
                    text-white text-sm font-bold shrink-0">
          {{ c.username?.[0]?.toUpperCase() || 'U' }}
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 mb-1">
            <span class="text-sm font-medium text-text-primary">{{ c.username }}</span>
            <span class="text-xs text-text-muted">{{ formatDate(c.created_at) }}</span>
          </div>
          <p class="text-sm text-text-secondary leading-relaxed">{{ c.content }}</p>
          <div class="flex items-center gap-3 mt-2">
            <button type="button" @click="toggleLikeComment(c)" :disabled="pendingLikeIds.has(c.comment_id)"
                    :class="[c.is_liked ? 'text-primary' : 'text-text-muted hover:text-primary',
                 'flex items-center gap-1 text-xs transition-colors disabled:opacity-60 disabled:cursor-wait']">
              <svg class="w-3.5 h-3.5" :fill="c.is_liked ? 'currentColor' : 'none'"
                   stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905a3.61 3.61 0 01-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"/>
              </svg>
              {{ c.like_count || 0 }}
            </button>
          </div>
        </div>
      </div>
      <div v-if="!comments.length && !loading" class="text-center py-8 text-text-muted text-sm">
        还没有评论，来抢沙发吧
      </div>
    </div>

    <!-- 发表评论 -->
    <div v-if="userStore.isLoggedIn" class="mb-6">
      <textarea v-model="content" rows="3" placeholder="发表你的看法..."
                class="input-base resize-none mb-2"></textarea>
      <div class="flex justify-end">
        <button @click="submitComment" :disabled="!content.trim()" class="btn-primary text-sm px-4 py-2
          disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:scale-100 disabled:hover:shadow-none">
          发表评论
        </button>
      </div>
    </div>
    <div v-else class="mb-6 p-4 rounded-lg bg-bg-surface border border-border text-center text-sm text-text-muted">
      <button @click="userStore.openAuth('login')" class="text-primary hover:underline">登录</button>
      后参与评论
    </div>

    <Pagination v-if="totalPages > 1" :current="page" :total-pages="totalPages" @change="changePage"/>
  </div>
</template>

<script setup>
import {onMounted, ref, watch} from 'vue'
import {commentApi} from '@/api'
import {useUserStore} from '@/stores/user'
import Pagination from './Pagination.vue'

const props = defineProps({videoId: {type: Number, required: true}})
const emits = defineEmits(['comment-count-updated']) //
const userStore = useUserStore()

const comments = ref([])
const content = ref('')
const page = ref(1)
const total = ref(0)
const totalPages = ref(1)
const loading = ref(false)
const pendingLikeIds = ref(new Set())
const localLikedCommentIds = ref(new Set())

function setPendingLike(commentId, pending) {
  const next = new Set(pendingLikeIds.value)
  if (pending) {
    next.add(commentId)
  } else {
    next.delete(commentId)
  }
  pendingLikeIds.value = next
}

function setLocalLiked(commentId, liked) {
  const next = new Set(localLikedCommentIds.value)
  if (liked) {
    next.add(commentId)
  } else {
    next.delete(commentId)
  }
  localLikedCommentIds.value = next
}

function applyLikedState(comment, liked) {
  comment.is_liked = liked
  setLocalLiked(comment.comment_id, liked)
}

async function fetchComments() {
  loading.value = true
  try {
    const res = await commentApi.list({
      video_id: props.videoId,
      page: page.value,
      page_size: 20
    })

    comments.value = (res.data?.comments || []).map(comment => ({
      ...comment,
      is_liked: Boolean(comment.is_liked) || localLikedCommentIds.value.has(comment.comment_id)
    }))

    total.value = res.data?.total || 0
    totalPages.value = Math.ceil((res.data?.total || 0) / 20) || 1
    emits('comment-count-updated', total.value)
  } catch (err) {
    console.error('fetch comments failed', err)
  }
  loading.value = false
}

async function submitComment() {
  if (!content.value.trim()) return
  try {
    await commentApi.add({video_id: props.videoId, content: content.value.trim()})
    content.value = ''
    page.value = 1
    await fetchComments() // 重新获取评论列表以更新UI和点赞状态
  } catch (err) {
    console.error('submit comment failed', err)
    // 可以显示错误提示
  }
}

async function toggleLikeComment(comment) {
  if (!userStore.isLoggedIn) {
    userStore.openAuth('login')
    return
  }
  if (pendingLikeIds.value.has(comment.comment_id)) return

  setPendingLike(comment.comment_id, true)
  try {
    if (comment.is_liked) {
      // 取消点赞
      await commentApi.unlike(comment.comment_id)
      applyLikedState(comment, false)
      comment.like_count = Math.max(0, (comment.like_count || 0) - 1)
    } else {
      // 点赞
      await commentApi.like(comment.comment_id)
      applyLikedState(comment, true)
      comment.like_count = (comment.like_count || 0) + 1
    }
  } catch (err) {
    console.error('toggle like comment failed', err)
    if (err.message?.includes('已经点赞')) {
      applyLikedState(comment, true)
    } else if (err.message?.includes('未点赞')) {
      applyLikedState(comment, false)
    }
    // 可以在这里添加用户友好的错误提示，例如使用一个全局的 Toast 组件
    // alert(err.response?.data?.msg || '操作失败，请稍后重试')
  } finally {
    setPendingLike(comment.comment_id, false)
  }
}

function changePage(p) {
  page.value = p
}

function formatDate(t) {
  if (!t) return ''
  const d = new Date(t * 1000)
  const diff = (Date.now() - d.getTime()) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

watch(() => props.videoId, () => {
  page.value = 1;
  fetchComments()
})
watch(page, fetchComments)
onMounted(fetchComments)
</script>

<style scoped>
/* 样式保持不变 */
</style>
