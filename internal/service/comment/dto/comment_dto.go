package dto

import "adult-short-videos/internal/service/comment/model"

// AddCommentReq 发表评论请求
type AddCommentReq struct {
	VideoId  int64  `json:"video_id"`  // 视频 ID
	Content  string `json:"content"`   // 评论内容
	ParentId int64  `json:"parent_id"` // 父评论ID（0表示顶级评论）
}

// CommentListReq 评论列表请求
type CommentListReq struct {
	VideoId int64 // 视频 ID
	Page    int
	Size    int
}

// CommentListResp 评论列表响应
type CommentListResp struct {
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	Size     int                      `json:"size"`
	Comments []*model.CommentWithUser `json:"comments"`
}
