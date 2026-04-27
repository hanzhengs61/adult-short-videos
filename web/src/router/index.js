import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/',          component: () => import('@/views/HomeView.vue') },
  { path: '/subscribe', component: () => import('@/views/SubscribeView.vue') },
  { path: '/explore',   component: () => import('@/views/ExploreView.vue') },
  { path: '/rankings',  component: () => import('@/views/RankingsView.vue') },
  { path: '/profile',   component: () => import('@/views/ProfileView.vue') },
  { path: '/video/:id', component: () => import('@/views/VideoDetail.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

export default router
