package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// HeatCalculator 热度计算器
// 定期计算视频热度，决定是否需要从冷数据转为热数据
type HeatCalculator struct {
	db       *gorm.DB
	interval time.Duration // 计算间隔
}

// NewHeatCalculator 创建热度计算器
func NewHeatCalculator(db *gorm.DB, interval time.Duration) *HeatCalculator {
	return &HeatCalculator{
		db:       db,
		interval: interval,
	}
}

// Start 启动定时任务
func (h *HeatCalculator) Start(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	log.Println("🔥 热度计算器已启动")

	for {
		select {
		case <-ticker.C:
			h.calculateHeat()
		case <-ctx.Done():
			log.Println("⏹️ 热度计算器已停止")
			return
		}
	}
}

// calculateHeat 计算热度
func (h *HeatCalculator) calculateHeat() {
	log.Println("📊 开始计算视频热度...")

	// ========== 更新热度统计 ==========
	// 计算最近24小时和7天的播放次数
	sql := `
		INSERT INTO video_heat_stats (video_id, play_count24h, play_count7d, heat_score, is_hot, last_calculated_at)
		SELECT 
			v.video_id,
			COALESCE(COUNT(DISTINCT CASE 
				WHEN h.last_play_at > NOW() - INTERVAL '24 hours' 
				THEN h.id 
			END), 0) as play_count24h,
			COALESCE(COUNT(DISTINCT CASE 
				WHEN h.last_play_at > NOW() - INTERVAL '7 days' 
				THEN h.id 
			END), 0) as play_count7d,
			-- 热度分数 = 24小时播放*5 + 7天播放*1 + 总播放*0.1
			(
				COALESCE(COUNT(DISTINCT CASE WHEN h.last_play_at > NOW() - INTERVAL '24 hours' THEN h.id END), 0) * 5 +
				COALESCE(COUNT(DISTINCT CASE WHEN h.last_play_at > NOW() - INTERVAL '7 days' THEN h.id END), 0) * 1 +
				v.play_count * 0.1
			) as heat_score,
			-- 判断是否为热门：24小时播放>200 或 热度分数>100
			(
				COALESCE(COUNT(DISTINCT CASE WHEN h.last_play_at > NOW() - INTERVAL '24 hours' THEN h.id END), 0) > 200
				OR 
				(
					COALESCE(COUNT(DISTINCT CASE WHEN h.last_play_at > NOW() - INTERVAL '24 hours' THEN h.id END), 0) * 5 +
					COALESCE(COUNT(DISTINCT CASE WHEN h.last_play_at > NOW() - INTERVAL '7 days' THEN h.id END), 0) * 1 +
					v.play_count * 0.1
				) > 100
			) as is_hot,
			NOW() as last_calculated_at
		FROM videos v
		LEFT JOIN play_history h ON v.video_id = h.video_id
		WHERE v.status = 1
		GROUP BY v.video_id
		ON CONFLICT (video_id) 
		DO UPDATE SET
			play_count24h = EXCLUDED.play_count24h,
			play_count7d = EXCLUDED.play_count7d,
			heat_score = EXCLUDED.heat_score,
			is_hot = EXCLUDED.is_hot,
			last_calculated_at = EXCLUDED.last_calculated_at
	`

	result := h.db.Exec(sql)
	if result.Error != nil {
		log.Printf("❌ 热度计算失败: %v\n", result.Error)
		return
	}

	log.Printf("✅ 热度计算完成，更新了 %d 个视频\n", result.RowsAffected)

	// 识别需要热化的视频
	h.identifyHotVideos()
}

// identifyHotVideos 识别需要热化的视频
func (h *HeatCalculator) identifyHotVideos() {
	// 查找冷数据中的热门视频
	var videos []struct {
		VideoId int64
		Title   string
	}

	err := h.db.Table("videos v").
		Select("v.video_id, v.title").
		Joins("INNER JOIN video_heat_stats h ON v.video_id = h.video_id").
		Where("v.storage_type = ? AND h.is_hot = ?", "cold", true).
		Limit(10). // 每次最多处理10个
		Scan(&videos).Error

	if err != nil {
		log.Printf("❌ 查询待热化视频失败: %v\n", err)
		return
	}

	if len(videos) == 0 {
		return
	}

	log.Printf("🔥 发现 %d 个冷数据视频需要热化\n", len(videos))

	// TODO: 这里应该发送到消息队列（Kafka）让异步任务处理
	// 异步任务会：下载视频 → 去水印 → 加水印 → 转码 → 上传 MinIO
	for _, video := range videos {
		fmt.Printf("  - 视频 %d: %s\n", video.VideoId, video.Title)
		// 实际生产环境：
		// kafka.Produce("video-hot-process", video.VideoId)
	}
}
