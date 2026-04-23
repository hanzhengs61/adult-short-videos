-- ========================================
-- 数据库索引优化脚本
-- 用于提升查询和搜索性能
-- ========================================

-- ========== 视频表索引 ==========

-- 1. 全文搜索索引（PostgreSQL GIN 索引）
-- 用于标题和番号的快速搜索
CREATE INDEX IF NOT EXISTS idx_videos_title_gin
    ON videos USING gin(to_tsvector('simple', title));

CREATE INDEX IF NOT EXISTS idx_videos_fanhao_gin
    ON videos USING gin(to_tsvector('simple', fanhao));

-- 2. 复合索引：状态 + 地区 + 分类
-- 用于筛选查询
CREATE INDEX IF NOT EXISTS idx_videos_status_region_category
    ON videos(status, region, category);

-- 3. 复合索引：状态 + 创建时间
-- 用于最新视频列表
CREATE INDEX IF NOT EXISTS idx_videos_status_created
    ON videos(status, created_at DESC);

-- 4. 复合索引：状态 + 播放次数
-- 用于热门视频排序
CREATE INDEX IF NOT EXISTS idx_videos_status_play_count
    ON videos(status, play_count DESC);

-- 5. 存储类型索引
-- 用于冷热分离查询
CREATE INDEX IF NOT EXISTS idx_videos_storage_type
    ON videos(storage_type) WHERE status = 1;

-- ========== 演员表索引 ==========

-- 1. 演员名全文搜索
CREATE INDEX IF NOT EXISTS idx_actors_name_gin
    ON actors USING gin(to_tsvector('simple', actor_name));

-- 2. 粉丝数排序索引
CREATE INDEX IF NOT EXISTS idx_actors_follower_count
    ON actors(follower_count DESC);

-- ========== 视频-演员关联表索引 ==========

-- 1. 复合索引：视频ID + 演员ID
CREATE INDEX IF NOT EXISTS idx_video_actors_video_actor
    ON video_actors(video_id, actor_id);

-- 2. 反向索引：演员ID + 视频ID
-- 用于查询演员的所有作品
CREATE INDEX IF NOT EXISTS idx_video_actors_actor_video
    ON video_actors(actor_id, video_id);

-- ========== 收藏表索引 ==========

-- 1. 复合唯一索引：用户ID + 视频ID
-- 防止重复收藏
CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_video
    ON favorites(user_id, video_id);

-- 2. 创建时间索引（用于排序）
CREATE INDEX IF NOT EXISTS idx_favorites_created_at
    ON favorites(created_at DESC);

-- ========== 播放历史表索引 ==========

-- 1. 复合唯一索引：用户ID + 视频ID
-- 同一用户同一视频只保留一条记录
CREATE UNIQUE INDEX IF NOT EXISTS idx_play_history_user_video
    ON play_history(user_id, video_id);

-- 2. 最后播放时间索引
-- 用于按时间排序
CREATE INDEX IF NOT EXISTS idx_play_history_last_play
    ON play_history(last_play_at DESC);

-- 3. 复合索引：用户ID + 最后播放时间
-- 用于用户播放历史查询
CREATE INDEX IF NOT EXISTS idx_play_history_user_last_play
    ON play_history(user_id, last_play_at DESC);

-- ========== 热度统计表索引 ==========

-- 1. 热度分数索引
CREATE INDEX IF NOT EXISTS idx_video_heat_stats_heat_score
    ON video_heat_stats(heat_score DESC);

-- 2. 是否热门标记索引
CREATE INDEX IF NOT EXISTS idx_video_heat_stats_is_hot
    ON video_heat_stats(is_hot) WHERE is_hot = true;

-- ========== 用户统计表索引 ==========

-- 1. 最后活跃时间索引
CREATE INDEX IF NOT EXISTS idx_user_statistics_last_active
    ON user_statistics(last_active_at DESC);

-- ========================================
-- 创建数据库约束
-- ========================================

-- 确保播放进度在 0-100 之间
ALTER TABLE play_history
    ADD CONSTRAINT check_play_progress
        CHECK (play_progress >= 0 AND play_progress <= 100);

-- 确保视频时长大于 0
ALTER TABLE videos
    ADD CONSTRAINT check_duration
        CHECK (duration >= 0);

-- ========================================
-- 分析表统计信息（提升查询优化器性能）
-- ========================================

ANALYZE videos;
ANALYZE actors;
ANALYZE video_actors;
ANALYZE favorites;
ANALYZE play_history;
ANALYZE video_heat_stats;
ANALYZE user_statistics;

-- ========================================
-- 完成
-- ========================================
SELECT '✅ 索引优化完成' AS status;