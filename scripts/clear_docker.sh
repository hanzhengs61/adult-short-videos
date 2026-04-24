#!/bin/bash
set -e  # 遇错即停

echo "=========================================="
echo "  PostgreSQL 17 → 18 升级脚本"
echo "=========================================="

# ====== 配置区 ======
OLD_CONTAINER="adult-videos-postgres"
OLD_VOLUME="your_project_postgres_data"  # 根据实际情况修改
BACKUP_FILE="pg_backup_before_upgrade_$(date +%Y%m%d).sql"

# ====== Step 1: 备份 ======
echo ""
echo "📦 [Step 1/6] 备份现有数据..."
if docker ps | grep -q "$OLD_CONTAINER"; then
    docker exec "$OLD_CONTAINER" pg_dumpall -U postgres > "$BACKUP_FILE"
    echo "   ✅ 备份完成: $BACKUP_FILE ($(du -h $BACKUP_FILE | cut -f1))"
else
    echo "   ⚠️  容器未运行，跳过备份"
fi

# ====== Step 2: 停止服务 ======
echo ""
echo "🛑 [Step 2/6] 停止服务..."
docker compose down
echo "   ✅ 服务已停止"

# ====== Step 3: 处理旧 Volume ======
echo ""
echo "📁 [Step 3/6] 处理旧数据卷..."
if docker volume ls | grep -q "${OLD_VOLUME}"; then
    # 重命名旧 volume 作为备份
    docker volume create "${OLD_VOLUME}_backup_17" 2>/dev/null || true

    # 注意：Docker 不支持直接重命名 volume，这里我们创建新备份
    echo "   ℹ️  旧 volume: ${OLD_VOLUME}"
    echo "   💡 建议: 保留旧 volume 作为备份，创建新的空 volume"
    echo "   ⚠️  如果不需要旧数据，可以执行: docker volume rm ${OLD_VOLUME}"
else
    echo "   ℹ️  未找到旧 volume"
fi

# ====== Step 4: 更新配置 ======
echo ""
echo "✏️  [Step 4/6] 请确认已更新 docker-compose.yml:"
echo "   将 volumes: - postgres_data:/var/lib/postgresql/data"
echo "   改为:     volumes: - postgres_data:/var/lib/postgresql"
read -p "   已修改配置？(y/N): " config_confirm

if [[ $config_confirm != [yY] ]]; then
    echo "❌ 请先修改配置后重新运行此脚本"
    exit 1
fi

# ====== Step 5: 启动新版本 ======
echo ""
echo "🚀 [Step 5/6] 启动 PostgreSQL 18..."
docker compose up -d postgres

# 等待就绪
echo "   ⏳ 等待 PostgreSQL 就绪..."
for i in {1..30}; do
    if docker exec "$OLD_CONTAINER" pg_isready -U postgres 2>/dev/null; then
        break
    fi
    sleep 2
done

echo "   ✅ PostgreSQL 18 已启动"

# ====== Step 6: 恢复数据 ======
if [ -f "$BACKUP_FILE" ]; then
    echo ""
    "📥 [Step 6/6] 恢复数据..."
    docker exec -i "$OLD_CONTAINER" psql -U postgres < "$BACKUP_FILE"
    echo "   ✅ 数据恢复完成"

    # 验证
    echo ""
    echo "🔍 验证恢复结果:"
    docker exec "$OLD_CONTAINER" psql -U postgres -c "\l"
else
    echo ""
    echo "⏭️  [Step 6/6] 无备份数据，跳过恢复（全新安装）"
fi

echo ""
echo "=========================================="
echo "  ✅ 升级完成！"
echo "=========================================="
echo ""
echo "后续建议："
echo "1. 检查应用连接是否正常"
echo "2. 确认数据完整性"
echo "3. 一周后确认无问题再删除旧 volume"
echo "   docker volume rm ${OLD_VOLUME}"