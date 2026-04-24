# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言
所有回复使用中文。

## Commands

```bash
make run          # 启动应用（go run ./cmd/server/）
make build        # 编译到 bin/main
make test         # 运行所有测试
make fmt          # 格式化代码
make lint         # golangci-lint 检查
make deps         # go mod download && tidy
make docker-up    # 启动 PostgreSQL + Redis + App（via docker-compose）
make docker-down  # 停止所有容器

# 单个包测试
go test -v ./internal/service/user/...
go test -v -run TestXxx ./internal/service/video/...
```

运行前需要 PostgreSQL（5432）和 Redis（6379）。本地开发用 `make docker-up` 启动依赖，或只起依赖服务：
```bash
docker-compose -f deployments/docker-compose.yml up postgres redis -d
make run
```

配置文件：`configs/config.yaml`。敏感字段可被环境变量覆盖：`DB_PASSWORD`、`JWT_SECRET`、`REDIS_PASSWORD`。

## 架构

**单体应用**，按业务模块拆分目录，每个服务遵循四层结构：

```
internal/service/<service>/
  handler/    # HTTP 入口，绑定参数、调用 Logic
  logic/      # 业务逻辑，组合 Repository
  repository/ # 数据库操作，封装 GORM 查询
  model/      # GORM 模型定义
```

`main.go` 完成所有初始化：加载配置 → 日志 → DB → Redis → AutoMigrate → 建索引 → 注册路由 → 启动热度计算器。

**中间件**（`internal/pkg/middleware/`）：
- `AuthMiddleware`：JWT Bearer Token 解析，成功后将 `user_id`/`username` 写入 `gin.Context`
- `RateLimit`：基于 IP 的令牌桶限流（10 req/s，burst 20）
- `CORS`：全局放行（生产需收紧）

**冷热分离**：视频存储分 `cold`（源站 URL）和 `hot`（本地 CDN URL）两种类型。`internal/service/storage/scheduler/heat_calculator.go` 每 5 分钟运行一次，计算热度分数（24h播放×5 + 7d播放×1 + 总播放×0.1），`is_hot=true` 的冷数据视频会被标记待热化（当前打印日志，生产应发 Kafka）。

**防盗链代理**（`internal/service/storage/proxy/referer_proxy.go`）：冷数据播放地址为 `/api/storage/proxy?url=<源站URL>`，代理层伪造 Referer/Origin/User-Agent 访问源站，M3U8 文件中的分片 URL 会被递归重写为代理地址。

**认证路由**：部分路由组需要 JWT，例如 `/api/favorite/*`、`/api/play/*`、`/api/comment/add`。公开路由：`/api/video/*`、`/api/search/*`、`GET /api/comment/list`。

## 代码规范

- 单文件不超过 400 行，超出须拆分
- 嵌套层数不超过 4 层
- 优先编辑已有文件，不轻易新建
- 输出简洁，推理详尽
- 拆小任务 — 每次只解决一个具体问题，而不是"帮我把整个搜索模块重构一遍"