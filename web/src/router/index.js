import {createRouter, createWebHistory} from 'vue-router'

const manageRoutes = {
    path: '/manage',
    component: () => import('@/views/manage/ManageLayout.vue'),
    meta: { requiresAdmin: true },
    redirect: '/manage/dashboard',
    children: [
        { path: 'dashboard', component: () => import('@/views/manage/DashboardView.vue') },
        { path: 'videos',    component: () => import('@/views/manage/VideosView.vue') },
        { path: 'comments',  component: () => import('@/views/manage/CommentsView.vue') },
        { path: 'users',     component: () => import('@/views/manage/UsersView.vue') },
        { path: 'banners',   component: () => import('@/views/manage/AdsView.vue') },
        { path: 'apps',      component: () => import('@/views/manage/AdsView.vue') },
        { path: 'gossip',    component: () => import('@/views/manage/GossipView.vue') },
    ],
}

const routes = [
    {path: '/', component: () => import('@/views/HomeView.vue')},
    {path: '/subscribe', component: () => import('@/views/SubscribeView.vue')},
    {path: '/gossip', component: () => import('@/views/GossipView.vue')},
    {path: '/gossip/:id', component: () => import('@/views/GossipDetailView.vue')},
    {path: '/favorites', component: () => import('@/views/FavoritesView.vue')},
    {path: '/explore', component: () => import('@/views/ExploreView.vue')},
    {path: '/rankings', component: () => import('@/views/RankingsView.vue')},
    {path: '/profile', component: () => import('@/views/ProfileView.vue')},
    {path: '/video/:id', component: () => import('@/views/VideoDetail.vue')},
    {path: '/feed', component: () => import('@/views/FeedView.vue')},
    manageRoutes,
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    // savedPosition 在浏览器后退/前进时由浏览器提供，直接还原；新导航滚顶
    scrollBehavior: (to, from, savedPosition) => {
        if (savedPosition) {
            // 全局 CSS 开了 smooth scroll，浏览器后退恢复滚动时需要临时禁用，
            // 否则会从顶部快速滑到原位置，看起来像 Feed 返回时在“刷视频”。
            const html = document.documentElement
            const prevScrollBehavior = html.style.scrollBehavior
            html.style.scrollBehavior = 'auto'
            window.scrollTo(savedPosition.left ?? 0, savedPosition.top ?? 0)
            requestAnimationFrame(() => {
                html.style.scrollBehavior = prevScrollBehavior
            })
            return false
        }
        return {top: 0}
    },
})

// 管理后台守卫：检查 JWT 中的 role 字段
router.beforeEach((to) => {
    if (!to.meta.requiresAdmin) return true
    try {
        const token = localStorage.getItem('access_token')
        if (!token) return '/'
        // JWT payload 是 base64，无需验证签名，只看 role 字段
        const payload = JSON.parse(atob(token.split('.')[1]))
        if (payload.role !== 'admin') return '/'
    } catch {
        return '/'
    }
    return true
})

export default router
