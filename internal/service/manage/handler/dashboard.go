package handler

import (
	"time"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/manage/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ManageHandler struct {
	db *gorm.DB
}

func NewManageHandler(db *gorm.DB) *ManageHandler {
	return &ManageHandler{db: db}
}

// Dashboard 统计概览
// GET /api/manage/dashboard
func (h *ManageHandler) Dashboard(c *gin.Context) {
	var (
		totalVideos   int64
		totalUsers    int64
		totalComments int64
		todayPlays    int64
	)

	today := time.Now().Format("2006-01-02")
	h.db.Table("videos").Where("status = 1").Count(&totalVideos)
	h.db.Table("users").Count(&totalUsers)
	h.db.Table("comments").Where("status = 1").Count(&totalComments)
	h.db.Table("play_histories").Where("DATE(played_at) = ?", today).Count(&todayPlays)

	// 近 7 天每日播放量趋势
	type row struct {
		Date  string `gorm:"column:date"`
		Plays int64  `gorm:"column:plays"`
	}
	var rows []row
	h.db.Raw(`
		SELECT TO_CHAR(played_at, 'MM-DD') AS date, COUNT(*) AS plays
		FROM play_histories
		WHERE played_at >= NOW() - INTERVAL '7 days'
		GROUP BY TO_CHAR(played_at, 'MM-DD')
		ORDER BY date
	`).Scan(&rows)

	trend := make([]dto.TrendItem, 0, len(rows))
	for _, r := range rows {
		trend = append(trend, dto.TrendItem{Date: r.Date, Plays: r.Plays})
	}

	response.Success(c, dto.DashboardResp{
		TotalVideos:   totalVideos,
		TotalUsers:    totalUsers,
		TotalComments: totalComments,
		TodayPlays:    todayPlays,
		Trend:         trend,
	})
}
