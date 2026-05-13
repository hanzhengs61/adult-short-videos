# 前端项目说明文档 (Web Frontend)

本项目是基于 **Vue 3 (Composition API)** 和 **Vite** 构建的现代短视频前端应用。

## 🛠 技术栈
- **框架**: Vue 3
- **构建工具**: Vite
- **路由**: Vue Router (管理页面跳转)
- **状态管理**: Pinia (管理登录状态、用户信息)
- **样式**: Tailwind CSS (原子化 CSS 框架)
- **网络请求**: Axios (与后端 API 交互)
- **视频播放**: hls.js (支持 HLS 流媒体播放)

### 根目录配置文件
- `index.html`: 浏览器加载的 HTML 入口。
- `vite.config.js`: Vite 构建工具配置（如代理、插件）。
- `tailwind.config.js`: Tailwind CSS 样式定制。
- `package.json`: 项目依赖包及运行脚本管理。

## 🔄 核心业务流向
1. **用户操作**: 用户点击视频或登录。
2. **状态触发**: 调用 `src/stores/` 中的 actions。
3. **数据请求**: store 调用 `src/api/` 发起网络请求。
4. **后端交互**: 请求到达 Go 后端处理并返回 JSON 数据。
5. **界面更新**: Store 更新响应式数据，`src/views/` 或 `src/components/` 自动重新渲染界面。

## 🚀 开发命令
- 安装依赖: `npm install`
- 启动开发服务器: `npm run dev`
- 打包构建: `npm run build`

## 📁 目录结构及功能

### `src/` - 源代码核心
| 目录/文件 | 说明 |
| :--- | :--- |
| `main.js` | **入口文件**：初始化 Vue 实例并加载插件。 |
| `App.vue` | **根组件**：定义全局布局（导航栏、底部栏及路由视图出口）。 |
| **router/** | **路由层**：配置 URL 路径与页面的指向关系。 |
| **views/** | **页面层**：具体的业务页面（如 HomeView 首页, VideoDetail 视频详情）。 |
| **components/** | **组件层**：可复用的 UI 单元（如 NavBar 导航, VideoCard 视频卡片）。 |
| **stores/** | **状态层**：使用 Pinia 管理全局数据（如 user.js 存储用户信息）。 |
| **composables/** | **组合式函数**：抽取通用的逻辑代码（Hooks）。 |
| **assets/** | **资源层**：存放全局 CSS 样式和图片资源。 |
| **api/** | **接口层**：封装 axios 请求，定义所有与后端交互的函数。 |

## 📃 文件说明

### `src/main.js`
*   **业务作用**: 这是整个 Vue 应用的入口文件。它负责初始化 Vue 应用实例，并挂载所有全局插件和配置。
    *   **`createApp(App)`**: 创建 Vue 应用实例，将 `App.vue` 作为根组件。
    *   **`.use(createPinia())`**: 注册 Pinia 状态管理库，使得应用中的所有组件都可以访问和修改全局状态。
    *   **`.use(router)`**: 注册 Vue Router，启用客户端路由功能，管理页面之间的跳转和 URL 匹配。
    *   **`.use(createHead())`**: 注册 `@vueuse/head` 插件，用于动态管理页面 `<head>` 标签中的内容，如页面标题、meta 描述等，对 SEO 和用户体验有帮助。
    *   **`./assets/main.css`**: 导入全局 CSS 样式，包括 Tailwind CSS 的基础样式和项目自定义的全局样式。
    *   **`.mount('#app')`**: 将整个 Vue 应用挂载到 `index.html` 中 ID 为 `app` 的 DOM 元素上，使其在浏览器中可见和可交互。

### `src/App.vue`
*   **业务作用**: 这是 Vue 应用的根组件，负责定义整个应用的整体布局和结构。
    *   **全局布局**: 包含导航栏 (`NavBar`)、底部导航 (`BottomNav`，移动端)、桌面端页脚 (`footer`) 和全局认证弹窗 (`AuthModal`)。
    *   **路由视图**: `<router-view>` 组件是所有页面组件的渲染出口，根据当前 URL 路径显示对应的页面内容。
    *   **条件渲染**: 根据 `isFeed` 计算属性（判断当前是否是 `/feed` 路径），决定是否显示导航栏、底部导航和桌面端页脚。这是为了在沉浸式视频流页面 (`FeedView`) 隐藏这些通用 UI 元素。
    *   **页面切换动画**: 使用 `<transition name="fade">` 为页面切换提供了淡入淡出的动画效果，提升用户体验。
    *   **页脚链接**: 定义了页脚的发现、账号和关于链接，其中账号链接与用户登录/注册状态管理 (`userStore.openAuth`) 结合，点击时会触发认证弹窗。
    *   **用户状态初始化**: 在 `onMounted` 生命周期钩子中调用 `userStore.fetchInfo()`，在应用加载时尝试获取当前登录用户的信息，保持登录状态。

### `src/router/index.js`
*   **业务作用**: 负责定义和管理应用的所有前端路由，将 URL 路径映射到对应的页面组件。
    *   **`createRouter`**: 创建 Vue Router 实例。
    *   **`createWebHistory()`**: 使用 HTML5 History 模式，使 URL 看起来更干净（不带 `#`）。
    *   **`routes` 数组**: 定义了具体的路由规则：
        *   `/`: 映射到 `HomeView.vue` (首页)。
        *   `/subscribe`: 映射到 `SubscribeView.vue` (订阅页)。
        *   `/explore`: 映射到 `ExploreView.vue` (探索页)。
        *   `/rankings`: 映射到 `RankingsView.vue` (排行榜页)。
        *   `/profile`: 映射到 `ProfileView.vue` (个人主页)。
        *   `/video/:id`: 映射到 `VideoDetail.vue` (视频详情页)，`:id` 是动态参数，用于获取特定视频的详情。
        *   `/feed`: 映射到 `FeedView.vue` (沉浸式视频流页)。
        *   `/gossip`: 映射到 `GossipView.vue` (吃瓜文章列表页)。
        *   `/gossip/:id`: 映射到 `GossipDetailView.vue` (吃瓜文章详情页)。
        *   `/favorites`: 映射到 `FavoritesView.vue` (收藏列表页)。
    *   **`scrollBehavior: () => ({top: 0})`**: 配置路由切换时，页面滚动到顶部，提供更好的用户体验。

### `src/views/` (页面组件)
这些文件是构成应用各个独立页面的 Vue 组件。

*   **`FeedView.vue`**:
    *   **业务作用**: 实现沉浸式短视频上下滑动播放的体验，类似 TikTok 或 YouTube Shorts。
    *   **核心功能**:
        *   **视频列表展示**: 垂直滚动的视频列表，每个视频占据整个屏幕高度 (`snap-y snap-mandatory`)。
        *   **视频播放控制**:
            *   自动播放当前可见视频，暂停其他视频。
            *   支持 HLS (m3u8) 流媒体播放，兼容性处理 (hls.js 或原生 Safari)。
            *   点击视频区域可暂停/播放。
            *   显示播放进度条，支持拖动调整播放位置。
        *   **交互功能**:
            *   **收藏**: 用户可以收藏/取消收藏视频，并实时更新收藏数量。需要登录状态，未登录会弹出登录框。
            *   **评论**: 点击评论按钮弹出评论抽屉 (`CommentSection`)，显示和发表评论。
            *   **分享**: 提供分享功能，复制视频链接或调用原生分享接口。
        *   **数据加载**: 滚动到底部时自动加载更多视频 (`fetchMore`)。
        *   **路由跳转**: 支持通过 URL 参数直接定位到特定视频 (`route.query.id`)。
        *   **UI 元素**: 包含返回按钮、暂停图标、视频标题、播放量、收藏/评论/分享按钮及数量。

*   **`HomeView.vue`**: (根据命名推断)
    *   **业务作用**: 应用的首页，通常展示推荐视频、热门视频或用户关注的内容。
    *   **可能包含**: 视频列表、分类导航、搜索入口等。

*   **`ExploreView.vue`**:
    *   **业务作用**: 提供视频探索功能，用户可以通过搜索关键词、浏览热门标签或分类来发现新视频。
    *   **可能包含**: 搜索框、搜索历史、热门搜索词、视频分类列表、搜索结果展示。

*   **`RankingsView.vue`**:
    *   **业务作用**: 展示视频排行榜，例如热门视频榜、最新视频榜、收藏榜等。
    *   **可能包含**: 不同类型的排行榜切换、视频列表展示。

*   **`SubscribeView.vue`**:
    *   **业务作用**: 展示用户订阅的创作者或频道发布的视频内容。
    *   **可能包含**: 订阅列表、订阅创作者的最新视频流。

*   **`ProfileView.vue`**:
    *   **业务作用**: 用户的个人主页，展示用户的个人信息、已发布视频、收藏列表、观看历史等。
    *   **可能包含**: 用户头像、昵称、个人简介、视频列表（我的发布、我的收藏、观看历史）。

*   **`VideoDetail.vue`**:
    *   **业务作用**: 展示单个视频的详细信息页面。
    *   **可能包含**: 视频播放器、视频标题、描述、创作者信息、评论区、相关推荐视频。

*   **`GossipView.vue`**:
    *   **业务作用**: 吃瓜栏目的文章列表页，路由 `/gossip`。
    *   **核心功能**: 顶部标签筛选（全部/爆料/八卦/女优等）、三列卡片瀑布流、无限滚动加载，调用 `gossipApi.list`。

*   **`GossipDetailView.vue`**:
    *   **业务作用**: 吃瓜文章详情页，路由 `/gossip/:id`。
    *   **核心功能**: 封面图、标题、标签、正文（`v-html` 渲染 HTML）、浏览数，底部挂关联视频推荐，调用 `gossipApi.detail`。

*   **`FavoritesView.vue`**:
    *   **业务作用**: 用户收藏列表页，路由 `/favorites`，复用首页瀑布流布局。
    *   **核心功能**: 与 `HomeView` 相同的骨架屏 + 瀑布流 + 无限滚动，数据来源改为 `favoriteApi.list`；未登录时展示引导登录提示，已登录无收藏时引导去首页。

### `src/components/common/` (通用组件)
这些文件是可在多个页面或组件中复用的 UI 模块。

*   **`NavBar.vue`**:
    *   **业务作用**: 应用顶部的导航栏。
    *   **核心功能**: 品牌 Logo、搜索入口、用户头像/登录入口、导航链接等。

*   **`SortBar.vue`**:
    *   **业务作用**: 提供排序和筛选功能。
    *   **核心功能**: 允许用户选择不同的排序方式（如按时间、按热度）或筛选条件（如视频时长、分类）。

*   **`AuthModal.vue`**:
    *   **业务作用**: 负责处理用户的登录和注册流程。
    *   **核心功能**: 弹窗形式的登录/注册表单，与 `userStore` 交互进行用户认证。

*   **`BottomNav.vue`**:
    *   **业务作用**: 移动端底部的导航栏。
    *   **核心功能**: 提供快速切换主要页面（如首页、探索、个人中心）的入口。

*   **`VideoCard.vue`**:
    *   **业务作用**: 以卡片形式展示单个视频的缩略信息。
    *   **核心功能**: 视频封面、标题、播放量、创作者信息等，通常用于视频列表展示。

*   **`Pagination.vue`**:
    *   **业务作用**: 提供分页导航功能。
    *   **核心功能**: 显示页码、上一页/下一页按钮，用于在数据量大时分批加载内容。

*   **`CommentSection.vue`**:
    *   **业务作用**: 视频的评论区。
    *   **核心功能**: 显示视频的评论列表、发表评论、点赞评论等。

*   **`GossipCard.vue`**:
    *   **业务作用**: 吃瓜文章的卡片组件，用于列表/搜索结果中展示单篇文章。
    *   **核心功能**: 封面图（无封面时显示占位渐变背景）、最多 3 个标签浮层、标题（2 行截断）、摘要（2 行截断）、浏览数与发布时间；点击跳转 `/gossip/:id`。

### `src/composables/` (组合式函数 / Hooks)
这些文件将复杂的有状态逻辑从视图层抽离，保持组件简洁。

*   **`useExploreSearch.js`**:
    *   **业务作用**: 管理探索页的搜索交互逻辑，同时支持**短视频**和**吃瓜**两种搜索类型。
    *   **核心功能**:
        *   **双类型搜索**: `searchTab` 控制当前搜索类型（`video` / `gossip`），切换时对无缓存侧懒加载，避免重复请求。
        *   **竞态保护**: `searchToken` 递增令牌，确保旧请求的响应不会覆盖新结果。
        *   **无限滚动**: 分别为短视频、吃瓜列表维护 `IntersectionObserver`（`searchSentinelOb` / `gossipSentinelOb`），滚动到哨兵元素时自动 `loadMore`。
        *   **响应式列数**: `syncColumnCount()` 监听 `resize` 事件，按断点（640 / 768 / 1280 px）动态调整瀑布流列数（2~5 列）。
        *   **搜索历史**: 每次搜索调用 `exploreStore.addHistory` 记录关键词。
        *   **排序切换**: `setSearchSort` 仅对短视频侧生效，切换后立即重新检索。

### `src/stores/` (状态管理模块)
这些文件使用 Pinia 定义了全局状态及其操作。

*   **`user.js`**:
    *   **业务作用**: 管理用户的认证状态和个人信息。
    *   **核心状态**:
        *   `token`: 用户登录凭证 (存储在 localStorage)。
        *   `userInfo`: 用户的详细信息（ID, 用户名等）。
        *   `authModal`: 控制登录/注册弹窗的显示状态和模式。
    *   **核心操作**:
        *   `isLoggedIn`: 计算属性，判断用户是否已登录。
        *   `openAuth(mode)` / `closeAuth()` / `switchAuthMode(mode)`: 控制认证弹窗的打开、关闭和模式切换。
        *   `login(data)`: 调用 `userApi.login` 进行登录，成功后保存 token 和用户信息。
        *   `logout()`: 调用 `userApi.logout` 登出，清除 token 和用户信息。
        *   `fetchInfo()`: 调用 `userApi.info` 获取用户最新信息。

*   **`explore.js`**:
    *   **业务作用**: 管理探索页面的相关状态，特别是搜索历史。
    *   **核心状态**:
        *   `history`: 搜索历史关键词列表 (存储在 localStorage)。
    *   **核心操作**:
        *   `addHistory(keyword)`: 添加关键词到搜索历史，并去重、限制数量。
        *   `removeHistory(keyword)`: 从搜索历史中移除指定关键词。
        *   `clearHistory()`: 清空所有搜索历史。

### `src/api/` (API 请求模块)
这些文件封装了与后端 Go 服务进行数据交互的逻辑。

*   **`request.js`**:
    *   **业务作用**: 配置和封装 Axios HTTP 客户端，处理请求和响应的通用逻辑。
    *   **核心功能**:
        *   创建 Axios 实例。
        *   **请求拦截器**: 统一添加 `Authorization` 头（如果存在 token），设置 `Content-Type` 等。
        *   **响应拦截器**: 统一处理后端返回的错误码、网络错误、认证失败（如 token 过期时自动登出或提示重新登录）。
        *   简化 API 调用，提供 `get`, `post`, `put`, `delete` 等方法。

*   **`index.js`**:
    *   **业务作用**: 集中管理和导出所有具体的业务 API 接口。
    *   **核心功能**: 将不同业务模块（用户、视频、搜索、收藏、播放、评论）的 API 请求函数进行分类和封装，方便在组件或 store 中调用。
        *   **`userApi`**: 包含 `register`, `login`, `logout`, `info`, `creators` 等用户相关的接口。
        *   **`videoApi`**: 包含 `list`, `detail`, `popular` 等视频相关的接口。
        *   **`searchApi`**: 包含 `videos`, `advanced`, `byFanhao` 等搜索相关的接口。
        *   **`favoriteApi`**: 包含 `add`, `remove`, `list`, `check` 等收藏相关的接口。
        *   **`playApi`**: 包含 `record`, `history` 等播放历史相关的接口。
        *   **`commentApi`**: 包含 `list`, `add`, `like`, `delete` 等评论相关的接口。
        *   **`gossipApi`**: 包含 `list`, `detail` 等吃瓜文章相关的接口。