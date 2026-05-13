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

// Start 启动定时任务，启动时立即计算一次，之后按 interval 周期执行
func (h *HeatCalculator) Start(ctx context.Context) {
	log.Println("🔥 热度计算器已启动")
	h.calculateHeat()

	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

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
	// 热度分 = 本平台真实播放行为（24h活跃用户*5 + 7d活跃用户*1）
	// is_hot 按热度分排名取前50，平台播放数相同时以源站播放量兜底排序
	// 不直接用源站 play_count 算热度分，避免采集视频全部变热门
	sql := `
		WITH stats AS (
			SELECT
				v.video_id,
				v.play_count AS source_play_count,
				COALESCE(COUNT(DISTINCT CASE
					WHEN h.last_play_at > NOW() - INTERVAL '24 hours'
					THEN h.id
				END), 0) AS play_count24h,
				COALESCE(COUNT(DISTINCT CASE
					WHEN h.last_play_at > NOW() - INTERVAL '7 days'
					THEN h.id
				END), 0) AS play_count7d
			FROM videos v
			LEFT JOIN play_history h ON v.video_id = h.video_id
			WHERE v.status = 1
			GROUP BY v.video_id, v.play_count
		),
		ranked AS (
			SELECT *,
				(play_count24h * 5 + play_count7d * 1 + source_play_count * 0.01) AS heat_score,
				RANK() OVER (
					ORDER BY (play_count24h * 5 + play_count7d * 1) DESC,
					         source_play_count DESC
				) AS rnk
			FROM stats
		)
		INSERT INTO video_heat_stats (video_id, play_count24h, play_count7d, heat_score, is_hot, last_calculated_at)
		SELECT
			video_id, play_count24h, play_count7d, heat_score,
			(rnk <= 50) AS is_hot,
			NOW()
		FROM ranked
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
