import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { searchApi, gossipApi } from '@/api'

const PAGE_SIZE = 20

// 管理探索页搜索状态：短视频 + 吃瓜两个类型
export function useExploreSearch({ router, exploreStore }) {
  const searchInput = ref('')
  const currentQuery = ref('')
  const searchTab = ref('video') // 'video' | 'gossip'，下拉框控制

  // 短视频搜索状态
  const searchResults = ref([])
  const searchSort = ref('created_at')
  const searching = ref(false)
  const searchLoadingMore = ref(false)
  const searchHasMore = ref(false)
  const searchPage = ref(1)
  const searchSentinel = ref(null)

  // 吃瓜搜索状态
  const gossipResults = ref([])
  const gossipSearching = ref(false)
  const gossipLoadingMore = ref(false)
  const gossipHasMore = ref(false)
  const gossipPage = ref(1)
  const gossipSentinel = ref(null)

  const columnCount = ref(2)

  // 每次搜索递增，防止竞态覆盖
  let searchToken = 0
  let searchSentinelOb = null
  let gossipSentinelOb = null

  const searchResultColumns = computed(() => {
    const cols = Array.from({ length: columnCount.value }, () => [])
    searchResults.value.forEach((v, idx) => cols[idx % columnCount.value].push(v))
    return cols
  })

  const skeletonColumns = computed(() => {
    const cols = Array.from({ length: columnCount.value }, () => [])
    for (let i = 0; i < 10; i++) cols[i % columnCount.value].push(i)
    return cols
  })

  function syncColumnCount() {
    const w = window.innerWidth
    if (w >= 1280) columnCount.value = 5
    else if (w >= 768) columnCount.value = 4
    else if (w >= 640) columnCount.value = 3
    else columnCount.value = 2
  }

  // 切换类型（下拉框触发），对无缓存的类型懒加载
  function setSearchTab(tab) {
    if (searchTab.value === tab) return
    searchTab.value = tab
    if (!currentQuery.value) return
    if (tab === 'gossip' && !gossipResults.value.length && !gossipSearching.value) {
      doGossipSearch(currentQuery.value)
    }
    if (tab === 'video' && !searchResults.value.length && !searching.value) {
      doVideoSearch(currentQuery.value)
    }
  }

  // 提交搜索（按当前下拉类型搜索，重置两侧缓存）
  async function doSearch(keyword) {
    const q = (keyword ?? searchInput.value).trim()
    if (!q) return
    currentQuery.value = q
    searchInput.value = q
    // 重置双侧缓存
    searchResults.value = []
    searchPage.value = 1
    searchHasMore.value = false
    gossipResults.value = []
    gossipPage.value = 1
    gossipHasMore.value = false
    exploreStore.addHistory(q)
    if (searchTab.value === 'gossip') {
      await doGossipSearch(q)
    } else {
      await doVideoSearch(q)
    }
  }

  // 短视频搜索内部实现
  async function doVideoSearch(q) {
    const token = ++searchToken
    searching.value = true
    searchResults.value = []
    searchPage.value = 1
    searchHasMore.value = false
    try {
      const res = await searchApi.videos({ keyword: q, page: 1, page_size: PAGE_SIZE, order_by: searchSort.value })
      if (token !== searchToken) return
      const list = res.data?.videos || []
      searchResults.value = list
      searchPage.value = 2
      searchHasMore.value = list.length >= PAGE_SIZE
    } catch {
      if (token !== searchToken) return
      searchResults.value = []
    } finally {
      if (token === searchToken) searching.value = false
    }
    await nextTick()
    observeSearchSentinel()
  }

  async function loadMoreSearch() {
    if (searchLoadingMore.value || !searchHasMore.value || !currentQuery.value) return
    const token = searchToken
    searchLoadingMore.value = true
    try {
      const res = await searchApi.videos({
        keyword: currentQuery.value, page: searchPage.value,
        page_size: PAGE_SIZE, order_by: searchSort.value,
      })
      if (token !== searchToken) return
      const list = res.data?.videos || []
      searchResults.value.push(...list)
      searchPage.value++
      searchHasMore.value = list.length >= PAGE_SIZE
    } catch {}
    if (token === searchToken) searchLoadingMore.value = false
  }

  function observeSearchSentinel() {
    if (searchSentinelOb) { searchSentinelOb.disconnect(); searchSentinelOb = null }
    if (!searchSentinel.value) return
    searchSentinelOb = new IntersectionObserver(([e]) => {
      if (e.isIntersecting) loadMoreSearch()
    }, { threshold: 0.1 })
    searchSentinelOb.observe(searchSentinel.value)
  }

  // 吃瓜搜索内部实现
  async function doGossipSearch(q) {
    if (!q) return
    gossipResults.value = []
    gossipPage.value = 1
    gossipHasMore.value = false
    gossipSearching.value = true
    try {
      const res = await gossipApi.list({ keyword: q, page: 1, page_size: PAGE_SIZE })
      const list = res.data?.posts || []
      gossipResults.value = list
      gossipPage.value = 2
      gossipHasMore.value = list.length >= PAGE_SIZE
    } catch {
      gossipResults.value = []
    } finally {
      gossipSearching.value = false
    }
    await nextTick()
    observeGossipSentinel()
  }

  async function loadMoreGossip() {
    if (gossipLoadingMore.value || !gossipHasMore.value || !currentQuery.value) return
    gossipLoadingMore.value = true
    try {
      const res = await gossipApi.list({ keyword: currentQuery.value, page: gossipPage.value, page_size: PAGE_SIZE })
      const list = res.data?.posts || []
      gossipResults.value.push(...list)
      gossipPage.value++
      gossipHasMore.value = list.length >= PAGE_SIZE
    } catch {}
    gossipLoadingMore.value = false
  }

  function observeGossipSentinel() {
    if (gossipSentinelOb) { gossipSentinelOb.disconnect(); gossipSentinelOb = null }
    if (!gossipSentinel.value) return
    gossipSentinelOb = new IntersectionObserver(([e]) => {
      if (e.isIntersecting) loadMoreGossip()
    }, { threshold: 0.1 })
    gossipSentinelOb.observe(gossipSentinel.value)
  }

  function setSearchSort(v) {
    if (searchSort.value === v) return
    searchSort.value = v
    if (currentQuery.value && searchTab.value === 'video') doVideoSearch(currentQuery.value)
  }

  function clearSearch() {
    searchInput.value = ''
    currentQuery.value = ''
    searchResults.value = []
    searchHasMore.value = false
    gossipResults.value = []
    gossipHasMore.value = false
    if (searchSentinelOb) { searchSentinelOb.disconnect(); searchSentinelOb = null }
    if (gossipSentinelOb) { gossipSentinelOb.disconnect(); gossipSentinelOb = null }
    router.replace({ path: '/explore', query: {} })
  }

  onMounted(() => {
    syncColumnCount()
    window.addEventListener('resize', syncColumnCount)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', syncColumnCount)
    if (searchSentinelOb) searchSentinelOb.disconnect()
    if (gossipSentinelOb) gossipSentinelOb.disconnect()
  })

  return {
    searchInput, currentQuery, searchTab,
    searchResults, searchSort, searching, searchLoadingMore, searchHasMore, searchSentinel,
    gossipResults, gossipSearching, gossipLoadingMore, gossipHasMore, gossipSentinel,
    columnCount, searchResultColumns, skeletonColumns,
    setSearchTab, doSearch, setSearchSort, clearSearch,
  }
}
