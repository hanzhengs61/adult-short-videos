import {defineStore} from 'pinia'
import {ref} from 'vue'

// 探索页：搜索历史 + 热搜词存储
export const useExploreStore = defineStore('explore', () => {
    const MAX_HISTORY = 12
    const history = ref(JSON.parse(localStorage.getItem('search_history') || '[]'))

    function addHistory(keyword) {
        const k = keyword.trim()
        if (!k) return
        const arr = [k, ...history.value.filter(h => h !== k)].slice(0, MAX_HISTORY)
        history.value = arr
        localStorage.setItem('search_history', JSON.stringify(arr))
    }

    function removeHistory(keyword) {
        history.value = history.value.filter(h => h !== keyword)
        localStorage.setItem('search_history', JSON.stringify(history.value))
    }

    function clearHistory() {
        history.value = []
        localStorage.removeItem('search_history')
    }

    return {history, addHistory, removeHistory, clearHistory}
})
