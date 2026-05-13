package handler

import (
	"strconv"
	"time"

	"adult-short-videos/internal/pkg/response"
	"adult-short-videos/internal/service/manage/dto"

	"github.com/gin-gonic/gin"
)

// ── 视频管理 ─────────────────────────────────────────

// VideoList 视频列表（含搜索、状态筛选）
// GET /api/manage/videos
func (h *ManageHandler) VideoList(c *gin.Context) {
	var req dto.VideoListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	q := h.db.Table("videos")
	if req.Status != nil {
		q = q.Where("status = ?", *req.Status)
	}
	if req.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	var total int64
	q.Count(&total)

	type row struct {
		VideoId   int64     `gorm:"column:video_id"`
		Title     string    `gorm:"column:title"`
		CoverURL  string    `gorm:"column:cover_url"`
		Duration  int32     `gorm:"column:duration"`
		PlayCount int64     `gorm:"column:play_count"`
		Status    int32     `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var rows []row
	q.Select("video_id,title,cover_url,duration,play_count,status,created_at").
		Order("created_at DESC").
		Offset(req.Offset()).Limit(req.PageSize).
		Scan(&rows)

	items := make([]dto.VideoItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.VideoItem{
			VideoId:   r.VideoId,
			Title:     r.Title,
			CoverURL:  r.CoverURL,
			Duration:  r.Duration,
			PlayCount: r.PlayCount,
			Status:    r.Status,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	response.Success(c, dto.VideoListResp{Total: total, Videos: items})
}

// VideoUpdateStatus 修改视频状态（上架/下架）
// PUT /api/manage/videos/:id/status
func (h *ManageHandler) VideoUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的视频 ID")
		return
	}
	var req dto.VideoStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	h.db.Table("videos").Where("video_id = ?", id).Update("status", req.Status)
	response.Success(c, nil)
}

// VideoDelete 删除视频
// DELETE /api/manage/videos/:id
func (h *ManageHandler) VideoDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.db.Table("videos").Where("video_id = ?", id).Delete(nil)
	response.Success(c, nil)
}

// ── 评论管理 ─────────────────────────────────────────

// CommentList 评论列表
// GET /api/manage/comments
func (h *ManageHandler) CommentList(c *gin.Context) {
	var req dto.CommentListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	type row struct {
		CommentId  int64     `gorm:"column:comment_id"`
		VideoId    int64     `gorm:"column:video_id"`
		VideoTitle string    `gorm:"column:video_title"`
		UserId     int64     `gorm:"column:user_id"`
		Username   string    `gorm:"column:username"`
		Content    string    `gorm:"column:content"`
		LikeCount  int64     `gorm:"column:like_count"`
		CreatedAt  time.Time `gorm:"column:created_at"`
	}

	q := h.db.Table("comments c").
		Select("c.comment_id,c.video_id,v.title AS video_title,c.user_id,u.username,c.content,c.like_count,c.created_at").
		Joins("LEFT JOIN videos v ON v.video_id = c.video_id").
		Joins("LEFT JOIN users u ON u.user_id = c.user_id").
		Where("c.status = 1")

	if req.VideoId > 0 {
		q = q.Where("c.video_id = ?", req.VideoId)
	}
	if req.Keyword != "" {
		q = q.Where("c.content LIKE ?", "%"+req.Keyword+"%")
	}

	var total int64
	h.db.Table("comments").Where("status = 1").Count(&total)

	var rows []row
	q.Order("c.created_at DESC").Offset(req.Offset()).Limit(req.PageSize).Scan(&rows)

	items := make([]dto.CommentItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.CommentItem{
			CommentId:  r.CommentId,
			VideoId:    r.VideoId,
			VideoTitle: r.VideoTitle,
			UserId:     r.UserId,
			Username:   r.Username,
			Content:    r.Content,
			LikeCount:  r.LikeCount,
			CreatedAt:  r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	response.Success(c, dto.CommentListResp{Total: total, Comments: items})
}

// CommentDelete 删除单条评论
// DELETE /api/manage/comments/:id
func (h *ManageHandler) CommentDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.db.Table("comments").Where("comment_id = ?", id).Update("status", 0)
	response.Success(c, nil)
}

// CommentBatchDelete 批量删除评论
// DELETE /api/manage/comments
func (h *ManageHandler) CommentBatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	h.db.Table("comments").Where("comment_id IN ?", req.IDs).Update("status", 0)
	response.Success(c, nil)
}
