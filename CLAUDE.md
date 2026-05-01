# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言要求
- 所有对话使用中文
- 所有文档使用中文
- 所有代码注释使用中文

## 执行要求
- 生成说明、总结、计划、提交说明时，统一使用中文，并使用中文注释
- 

## 全栈架构

**后端**：Go 单体服务（Gin + GORM + PostgreSQL + Redis）  
**前端**：Vue 3 + Vite + Pinia + Tailwind CSS（在 `web/` 目录）  
**核心特性**：冷热分离、防盗链代理、HLS 流媒体、沉浸式短视频列表、评论/收藏/播放记录

## Commands

### 后端 (Go)
```bash
make run              # go run ./cmd/server/
make build            # 编译到 bin/main
make test             # 运行所有测试
make fmt              # 格式化代码
make lint             # golangci-lint 检查
make deps             # go mod download && tidy
make docker-up        # 启动 PostgreSQL + Redis + App
make docker-down      # 停止所有容器
make clean            # 清理 bin/ 和日志
make doc              # swag 生成 API 文档
make bench            # go test -bench=. -benchmem
make migrate          # 数据库迁移

# 单包测试
go test -v ./internal/service/user/...
go test -v -run TestXxx ./internal/service/video/...
```

### 前端 (Vue 3)
```bash
cd web
npm install           # 安装依赖
npm run dev           # 启动开发服务器（http://localhost:5173）
npm run build         # 生产构建到 dist/
npm run preview       # 预览生产构建
```

### 环境配置
依赖：PostgreSQL（5432）+ Redis（6379）。本地快速启动：
```bash
make docker-up        # 一键启动 PostgreSQL + Redis + App
```

仅启动依赖服务：
```bash
docker-compose -f deployments/docker-compose.yml up postgres redis -d
make run
```

配置文件：`configs/config.yaml`。环境变量覆盖：`DB_PASSWORD`、`JWT_SECRET`、`REDIS_PASSWORD`。

## 后端架构

**单体应用**，按业务模块拆分，每个服务遵循四层结构：
```
internal/service/<service>/
  handler/    # HTTP 入口，绑定参数、调用 Logic
  logic/      # 业务逻辑，组合 Repository
  repository/ # 数据库操作，GORM 查询
  model/      # GORM 模型定义
```

`main.go` 初始化流程：配置 → 日志 → DB → Redis → AutoMigrate → 索引 → 路由 → 热度计算器。

**中间件**（`internal/pkg/middleware/`）：
- `AuthMiddleware`：JWT Bearer Token 解析，将 `user_id`/`username` 写入 Context
- `RateLimit`：基于 IP 的令牌桶限流（10 req/s，burst 20）
- `CORS`：全局放行（生产需收紧）

**冷热分离**：视频分 `cold`（源站 URL）和 `hot`（本地 CDN URL）。`heat_calculator.go` 每 5 分钟计算热度分数（24h播放×5 + 7d播放×1 + 总播放×0.1），标记待热化视频。

**防盗链代理**（`internal/service/storage/proxy/`）：冷数据播放地址为 `/api/storage/proxy?url=<源站URL>`，代理伪造 Referer/Origin/User-Agent，M3U8 分片 URL 递归重写为代理地址。

**路由**：认证路由 `/api/favorite/*`、`/api/play/*`、`/api/comment/add`。公开路由 `/api/video/*`、`/api/search/*`、`GET /api/comment/list`。

## 前端架构

**文件组织**：
- `views/`：页面组件（HomeView、FeedView 沉浸式列表、VideoDetail、ProfileView、ExploreView 等）
- `components/`：可复用 UI（NavBar、VideoCard、CommentSection、AuthModal 等）
- `stores/`：Pinia 状态管理（user 登录状态、explore 搜索历史）
- `api/`：Axios 请求封装（request.js 拦截器、index.js 聚合业务 API）

**核心流向**：用户操作 → Store actions → Store 调用 `api/` → 后端交互 → Store 更新响应式数据 → 视图重新渲染

**技术栈**：Vue 3 Composition API、Vue Router 4、Pinia 3、Tailwind CSS、hls.js（HLS 流）、Axios

## 代码规范

- 嵌套层数不超过 4 层
- 优先编辑已有文件，不轻易新建
- 输出简洁，推理详尽
- 拆小任务 — 每次只解决一个具体问题，而不是"帮我把整个搜索模块重构一遍"