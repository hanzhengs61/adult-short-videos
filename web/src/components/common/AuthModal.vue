<template>
  <teleport to="body">
    <transition name="modal-backdrop">
      <div v-if="userStore.authModal.show"
           class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4"
           style="background: rgba(0,0,0,0.75); backdrop-filter: blur(6px);"
           @click.self="userStore.closeAuth()">

        <transition name="modal-panel">
          <div v-if="userStore.authModal.show"
               class="w-full sm:max-w-sm bg-bg-surface border border-border
                   rounded-t-2xl sm:rounded-2xl shadow-2xl overflow-hidden">

            <!-- 顶部拖拽指示（移动端） -->
            <div class="flex justify-center pt-3 pb-1 sm:hidden">
              <div class="w-10 h-1 rounded-full bg-border-light"></div>
            </div>

            <!-- Tab 切换 -->
            <div class="flex border-b border-border mx-6 mt-2">
              <button @click="userStore.switchAuthMode('login')"
                      :class="['flex-1 py-3 text-sm font-semibold transition-all relative',
                  userStore.authModal.mode === 'login' ? 'text-text-primary' : 'text-text-muted hover:text-text-secondary']">
                登录
                <span v-if="userStore.authModal.mode === 'login'"
                      class="absolute bottom-0 left-1/4 right-1/4 h-0.5 bg-gradient-primary rounded-full"></span>
              </button>
              <button @click="userStore.switchAuthMode('register')"
                      :class="['flex-1 py-3 text-sm font-semibold transition-all relative',
                  userStore.authModal.mode === 'register' ? 'text-text-primary' : 'text-text-muted hover:text-text-secondary']">
                注册
                <span v-if="userStore.authModal.mode === 'register'"
                      class="absolute bottom-0 left-1/4 right-1/4 h-0.5 bg-gradient-primary rounded-full"></span>
              </button>
            </div>

            <!-- 内容区 -->
            <div class="px-6 py-5">
              <!-- 登录 -->
              <form v-if="userStore.authModal.mode === 'login'" @submit.prevent="handleLogin" class="space-y-3">
                <div>
                  <input v-model="loginForm.username" type="text" placeholder="用户名"
                         class="input-base" autocomplete="username" required/>
                </div>
                <div class="relative">
                  <input v-model="loginForm.password" :type="showPwd ? 'text' : 'password'"
                         placeholder="密码" class="input-base pr-10" autocomplete="current-password" required/>
                  <button type="button" @click="showPwd = !showPwd"
                          class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                </div>
                <div v-if="errorMsg"
                     class="text-xs text-red-400 bg-red-400/10 border border-red-400/20 rounded-lg px-3 py-2">
                  {{ errorMsg }}
                </div>
                <button type="submit" :disabled="submitting"
                        class="btn-primary w-full py-3 disabled:opacity-60 disabled:cursor-not-allowed
                         disabled:hover:scale-100 disabled:hover:shadow-none">
                  <span v-if="submitting"
                        class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                  {{ submitting ? '登录中...' : '立即登录' }}
                </button>
              </form>

              <!-- 注册 -->
              <form v-else @submit.prevent="handleRegister" class="space-y-3">
                <input v-model="regForm.username" type="text" placeholder="用户名（3-20字符）"
                       class="input-base" autocomplete="username" required/>
                <input v-model="regForm.email" type="email" placeholder="邮箱"
                       class="input-base" autocomplete="email" required/>
                <div class="relative">
                  <input v-model="regForm.password" :type="showPwd ? 'text' : 'password'"
                         placeholder="密码（至少8位）" class="input-base pr-10"
                         autocomplete="new-password" required minlength="8"/>
                  <button type="button" @click="showPwd = !showPwd"
                          class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                </div>
                <!-- 密码强度 -->
                <div class="flex gap-1">
                  <div v-for="i in 4" :key="i"
                       :class="['h-1 flex-1 rounded-full transition-colors duration-300',
                      pwdStrength >= i ? strengthColor : 'bg-border']"></div>
                </div>
                <div v-if="errorMsg"
                     class="text-xs text-red-400 bg-red-400/10 border border-red-400/20 rounded-lg px-3 py-2">
                  {{ errorMsg }}
                </div>
                <button type="submit" :disabled="submitting"
                        class="btn-primary w-full py-3 disabled:opacity-60 disabled:cursor-not-allowed
                         disabled:hover:scale-100 disabled:hover:shadow-none">
                  <span v-if="submitting"
                        class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                  {{ submitting ? '注册中...' : '免费注册' }}
                </button>
              </form>

              <!-- 社会认同 -->
              <p class="text-center text-xs text-text-muted mt-4">
                已有 <span class="text-primary font-semibold">10万+</span> 用户选择 mitun69
              </p>
            </div>
          </div>
        </transition>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import {computed, ref} from 'vue'
import {useUserStore} from '@/stores/user'
import {userApi} from '@/api'

const userStore = useUserStore()
const showPwd = ref(false)
const submitting = ref(false)
const errorMsg = ref('')

const loginForm = ref({username: '', password: ''})
const regForm = ref({username: '', email: '', password: ''})

const pwdStrength = computed(() => {
  const p = regForm.value.password
  if (!p) return 0
  let s = 0
  if (p.length >= 8) s++
  if (/[A-Z]/.test(p)) s++
  if (/[0-9]/.test(p)) s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return s
})
const strengthColor = computed(() => {
  if (pwdStrength.value <= 1) return 'bg-red-500'
  if (pwdStrength.value <= 2) return 'bg-yellow-500'
  if (pwdStrength.value <= 3) return 'bg-blue-500'
  return 'bg-green-500'
})

async function handleLogin() {
  errorMsg.value = ''
  submitting.value = true
  try {
    await userStore.login(loginForm.value)
    loginForm.value = {username: '', password: ''}
  } catch (err) {
    errorMsg.value = err?.message || '用户名或密码错误'
  }
  submitting.value = false
}

async function handleRegister() {
  errorMsg.value = ''
  submitting.value = true
  try {
    await userApi.register(regForm.value)
    await userStore.login({username: regForm.value.username, password: regForm.value.password})
    regForm.value = {username: '', email: '', password: ''}
  } catch (err) {
    errorMsg.value = err?.message || '注册失败，请稍后重试'
  }
  submitting.value = false
}
</script>

<style scoped>
.modal-backdrop-enter-active, .modal-backdrop-leave-active {
  transition: opacity 0.25s ease;
}

.modal-backdrop-enter-from, .modal-backdrop-leave-to {
  opacity: 0;
}

.modal-panel-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-panel-leave-active {
  transition: all 0.2s ease;
}

.modal-panel-enter-from {
  opacity: 0;
  transform: translateY(40px) scale(0.97);
}

.modal-panel-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

@media (max-width: 640px) {
  .modal-panel-enter-from {
    transform: translateY(100%);
  }

  .modal-panel-leave-to {
    transform: translateY(100%);
  }
}
</style>
