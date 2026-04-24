#!/bin/bash

# ========================================
# 应用停止脚本
# ========================================

echo "🛑 停止服务..."
docker-compose -f deployments/docker-compose.yml down

echo "✅ 服务已停止"