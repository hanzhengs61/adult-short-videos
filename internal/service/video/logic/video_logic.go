package logic

import (
	"adult-short-videos/internal/pkg/errors"
	"adult-short-videos/internal/service/video/model"
	"adult-short-videos/internal/service/video/repository"
	"context"

	"gorm.io/gorm"
)

// VideoListLogic 视频列表逻辑
type VideoListLogic struct {
	ctx       context.Context
	videoRepo repository.VideoRepository
}

// NewVideoListLogic 创建视频列表逻辑实例
func NewVideoListLogic(ctx context.Context, videoRepo repository.VideoRepository) *VideoListLogic {
	return &VideoListLogic{
		ctx:       ctx,
		videoRepo: videoRepo,
	}
}

// VideoListReq 视频列表请求参数
type VideoListReq struct {
	Page     int    // 页码（从1开始）
	Size     int    // 每页数量
	Region   string // 地区筛选（可选）
	Category string // 分类筛选（可选）
	OrderBy  string // 排序字段（可选）
}

// VideoItem 视频列表项（精简信息）
type VideoItem struct {
	VideoId       int64  `json:"video_id"`       // 视频 ID
	Title         string `json:"title"`          // 标题
	CoverURL      string `json:"cover_url"`      //  封面图 URL
	PreviewURL    string `json:"preview_url"`    //  WebP 预览图
	Duration      int32  `json:"duration"`       // 时长（秒）
	PlayCount     int64  `json:"play_count"`     // 播放次数
	FavoriteCount int64  `json:"favorite_count"` // 收藏数
	Fanhao        string `json:"fanhao"`         // 番号
	Region        string `json:"region"`         // 地区
	Category      string `json:"category"`       // 分类
	IsVipOnly     bool   `json:"is_vip_only"`    // 是否 VIP 专属
	PublishedAt   int64  `json:"published_at"`   // Unix 时间戳
}

// VideoListResp 视频列表响应
type VideoListResp struct {
	Total  int64       `json:"total"`  // 总数
	Page   int         `json:"page"`   // 当前页
	Size   int         `json:"size"`   // 每页数量
	Videos []VideoItem `json:"videos"` // 视频列表
}

// GetVideoList 获取视频列表
func (l *VideoListLogic) GetVideoList(req *VideoListReq) (*VideoListResp, error) {
	// 参数校验
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100 // 最大100条
	}

	// 构建筛选地区和分类条件
	filters := make(map[string]interface{})
	if req.Region != "" {
		filters["region"] = req.Region
	}
	if req.Category != "" {
		filters["category"] = req.Category
	}
	if req.OrderBy != "" {
		filters["order_by"] = req.OrderBy
	}

	// 计算偏移量
	// 例如：第1页 offset=0，第2页 offset=20
	offset := (req.Page - 1) * req.Size

	videos, total, err := l.videoRepo.List(l.ctx, offset, req.Size, filters)
	if err != nil {
		return nil, errors.New(errors.CodeVideoNotFound, "查询视频列表失败")
	}

	// 转换为响应格式
	items := make([]VideoItem, 0, len(videos))
	for _, video := range videos {
		items = append(items, VideoItem{
			VideoId:       video.VideoId,
			Title:         video.Title,
			CoverURL:      video.CoverURL,
			PreviewURL:    video.PreviewURL,
			Duration:      video.Duration,
			PlayCount:     video.PlayCount,
			FavoriteCount: video.FavoriteCount,
			Fanhao:        video.Fanhao,
			Region:        video.Region,
			Category:      video.Category,
			IsVipOnly:     video.IsVipOnly,
			PublishedAt:   video.PublishedAt.Unix(),
		})
	}

	return &VideoListResp{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Videos: items,
	}, nil
}

// VideoDetailLogic 视频详情逻辑
type VideoDetailLogic struct {
	ctx       context.Context
	videoRepo repository.VideoRepository
}

// NewVideoDetailLogic 创建视频详情逻辑实例
func NewVideoDetailLogic(ctx context.Context, videoRepo repository.VideoRepository) *VideoDetailLogic {
	return &VideoDetailLogic{
		ctx:       ctx,
		videoRepo: videoRepo,
	}
}

// ActorItem 演员信息
type ActorItem struct {
	ActorId   int64  `json:"actor_id"`   // 演员 ID
	ActorName string `json:"actor_name"` // 演员名字
	AvatarURL string `json:"avatar_url"` // 演员头像 URL
}

// VideoDetail 视频详情（完整信息）
type VideoDetail struct {
	VideoId       int64       `json:"video_id"`       // 视频 ID
	Title         string      `json:"title"`          // 标题
	Description   string      `json:"description"`    // 描述
	CoverURL      string      `json:"cover_url"`      // 封面图 URL
	PreviewURL    string      `json:"preview_url"`    // WebP 预览图
	Duration      int32       `json:"duration"`       // 时长（秒）
	StorageType   string      `json:"storage_type"`   // cold/hot
	PlayURL       string      `json:"play_url"`       // 播放地址
	PlayCount     int64       `json:"play_count"`     // 播放次数
	FavoriteCount int64       `json:"favorite_count"` // 收藏数
	CommentCount  int64       `json:"comment_count"`  // 评论数
	Fanhao        string      `json:"fanhao"`         // 番号
	Region        string      `json:"region"`         // 地区
	Category      string      `json:"category"`       // 分类
	Resolution    string      `json:"resolution"`     // 分辨率
	IsVipOnly     bool        `json:"is_vip_only"`    // 是否 VIP 专属
	Actors        []ActorItem `json:"actors"`         // 演员列表
	PublishedAt   int64       `json:"published_at"`   // 发布时间
}

// GetVideoDetail 获取视频详情
func (l *VideoDetailLogic) GetVideoDetail(videoId int64) (*VideoDetail, error) {
	// 查询视频信息
	video, err := l.videoRepo.FindByID(l.ctx, videoId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(errors.CodeVideoNotFound, "查询视频失败")
	}

	// 构建播放 URL
	// 根据存储类型决定播放地址
	playURL := l.buildPlayURL(video)

	// 查询演员信息
	actors, err := l.videoRepo.GetVideoActors(l.ctx, videoId)
	if err != nil {
		// 演员信息查询失败不影响主流程
		actors = []*model.Actor{}
	}

	// 转换演员信息
	actorItems := make([]ActorItem, 0, len(actors))
	for _, actor := range actors {
		actorItems = append(actorItems, ActorItem{
			ActorId:   actor.ActorId,
			ActorName: actor.ActorName,
			AvatarURL: actor.AvatarURL,
		})
	}

	// 异步增加播放次数
	// 使用 goroutine 异步执行，不影响响应速度
	go func() {
		l.videoRepo.IncrementPlayCount(context.Background(), videoId)
	}()

	// 返回详情
	return &VideoDetail{
		VideoId:       video.VideoId,
		Title:         video.Title,
		Description:   video.Description,
		CoverURL:      video.CoverURL,
		PreviewURL:    video.PreviewURL,
		Duration:      video.Duration,
		StorageType:   video.StorageType,
		PlayURL:       playURL,
		PlayCount:     video.PlayCount,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		Fanhao:        video.Fanhao,
		Region:        video.Region,
		Category:      video.Category,
		Resolution:    video.Resolution,
		IsVipOnly:     video.IsVipOnly,
		Actors:        actorItems,
		PublishedAt:   video.PublishedAt.Unix(),
	}, nil
}

// buildPlayURL 构建播放 URL
func (l *VideoDetailLogic) buildPlayURL(video *model.Video) string {
	if video.StorageType == "hot" {
		// 热数据：返回本地 CDN 地址
		return video.LocalURL
	}

	// 冷数据：返回代理地址
	// 格式: /api/storage/proxy?url=源站URL
	// 这样可以绕过防盗链
	return "/api/storage/proxy?url=" + video.SourceURL
}
