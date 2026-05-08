import {onMounted, onUnmounted, ref} from 'vue'

// 通用无限滚动 hook
// 用法：const { sentinel, loading, hasMore, reset } = useInfiniteScroll(fetchFn)
export function useInfiniteScroll(fetchFn) {
    const sentinel = ref(null)
    const loading = ref(false)
    const hasMore = ref(true)
    let observer = null

    async function load() {
        if (loading.value || !hasMore.value) return
        loading.value = true
        try {
            const more = await fetchFn()
            if (more === false) hasMore.value = false
        } catch (error) { // 捕获错误
            console.error('useInfiniteScroll fetchFn failed:', error) // 打印错误
            // 即使发生错误，也应该重置 loading 状态，以便可以再次尝试加载
        } finally {
            loading.value = false // 确保无论成功或失败都重置 loading 状态
        }
    }

    function reset() {
        hasMore.value = true
        loading.value = false
        // IntersectionObserver 只在交叉状态变化时触发；
        // 若 sentinel 仍在视口内（如上一次请求失败、列表为空），切换排序后 observer 不会重新触发。
        // 所以这里主动发起一次加载。
        load()
    }

    onMounted(() => {
        observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) load()
            },
            {rootMargin: '200px'}
        )
        if (sentinel.value) observer.observe(sentinel.value)
    })

    onUnmounted(() => observer?.disconnect())

    return {sentinel, loading, hasMore, reset}
}
