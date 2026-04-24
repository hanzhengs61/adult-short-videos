#!/bin/bash

# ========================================
# 应用启动脚本
# ========================================

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}短视频平台 - 启动脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安装${NC}"
    exit 1
fi

# 创建必要目录
echo -e "${YELLOW}📁 创建目录...${NC}"
mkdir -p logs
mkdir -p data/postgres
mkdir -p data/redis

# 启动服务
echo -e "${YELLOW}🚀 启动服务...${NC}"
docker-compose up -d

# 等待服务就绪
echo -e "${YELLOW}⏳ 等待服务就绪...${NC}"
sleep 10

# 检查服务状态
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ 服务启动完成${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 显示服务状态
docker-compose ps

echo ""
echo -e "${GREEN}📝 访问地址:${NC}"
echo -e "  应用: http://localhost:8080"
echo -e "  健康检查: http://localhost:8080/health"
echo ""
echo -e "${YELLOW}💡 提示:${NC}"
echo -e "  查看日志: docker-compose logs -f app"
echo -e "  停止服务: docker-compose down"
echo -e "  重启服务: docker-compose restart"
echo ""