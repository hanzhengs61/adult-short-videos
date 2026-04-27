import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const userInfo = ref(null)
  const authModal = ref({ show: false, mode: 'login' }) // 全局弹窗状态

  const isLoggedIn = computed(() => !!token.value)

  function openAuth(mode = 'login') {
    authModal.value = { show: true, mode }
  }
  function closeAuth() {
    authModal.value.show = false
  }
  function switchAuthMode(mode) {
    authModal.value.mode = mode
  }

  async function login(data) {
    const res = await userApi.login(data)
    token.value = res.data.access_token
    localStorage.setItem('access_token', res.data.access_token)
    userInfo.value = { user_id: res.data.user_id, username: res.data.username }
    closeAuth()
    return res
  }

  async function logout() {
    try { await userApi.logout() } catch {}
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('access_token')
  }

  async function fetchInfo() {
    if (!token.value) return
    try {
      const res = await userApi.info()
      userInfo.value = res.data
    } catch {}
  }

  return { token, userInfo, isLoggedIn, authModal, openAuth, closeAuth, switchAuthMode, login, logout, fetchInfo }
})
