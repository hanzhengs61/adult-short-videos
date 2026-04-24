# ========================================
# Makefile for Adult Short Videos Platform
# ========================================

.PHONY: help build run test clean docker-up docker-down logs

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help:
	@echo "Adult Short Videos Platform - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build         编译应用"
	@echo "  run           运行应用"
	@echo "  test          运行测试"
	@echo "  fmt           格式化代码"
	@echo "  lint          代码检查"
	@echo "  clean         清理构建文件"
	@echo "  docker-build  构建 Docker 镜像"
	@echo "  docker-up     启动 Docker 服务"
	@echo "  docker-down   停止 Docker 服务"
	@echo "  logs          查看日志"
	@echo ""

# 编译应用
build:
	@echo "🔨 编译应用..."
	@go build -o bin/main ./cmd/server/
	@echo "✅ 编译完成"

# 运行应用
run:
	@echo "🚀 运行应用..."
	@go run ./cmd/server/

# 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "📝 格式化代码..."
	@go fmt ./...
	@echo "✅ 格式化完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@golangci-lint run
	@echo "✅ 检查完成"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf bin/
	@rm -rf logs/*.log
	@echo "✅ 清理完成"

# 构建 Docker 镜像
docker-build:
	@echo "🐳 构建 Docker 镜像..."
	@docker build -f build/Dockerfile -t adult-videos-api:latest .
	@echo "✅ 镜像构建完成"

# 启动 Docker 服务
docker-up:
	@echo "🚀 启动 Docker 服务..."
	@chmod +x scripts/start.sh
	@./scripts/start.sh

# 停止 Docker 服务
docker-down:
	@echo "🛑 停止 Docker 服务..."
	@docker-compose -f deployments/docker-compose.yml down
	@echo "✅ 服务已停止"

# 查看日志
logs:
	@docker-compose -f deployments/docker-compose.yml logs -f app

# 数据库迁移
migrate:
	@echo "📊 数据库迁移..."
	@go run ./cmd/server/ migrate
	@echo "✅ 迁移完成"

# 生成 API 文档
doc:
	@echo "📚 生成 API 文档..."
	@swag init
	@echo "✅ 文档生成完成"

# 性能测试
bench:
	@echo "⚡ 性能测试..."
	@go test -bench=. -benchmem ./...

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖安装完成"