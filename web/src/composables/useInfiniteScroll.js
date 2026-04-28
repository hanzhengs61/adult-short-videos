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
        } catch {
        }
        loading.value = false
    }

    function reset() {
        hasMore.value = true
        loading.value = false
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
