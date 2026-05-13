package handler

import (
	"strconv"
	"time"

	"adult-short-videos/internal/pkg/response"
	adModel "adult-short-videos/internal/service/ad/model"
	gossipModel "adult-short-videos/internal/service/gossip/model"
	"adult-short-videos/internal/service/manage/dto"

	"github.com/gin-gonic/gin"
)

// ── 用户管理 ─────────────────────────────────────────

// UserList 用户列表
// GET /api/manage/users
func (h *ManageHandler) UserList(c *gin.Context) {
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	type row struct {
		UserId    int64     `gorm:"column:user_id"`
		Username  string    `gorm:"column:username"`
		Email     string    `gorm:"column:email"`
		Role      string    `gorm:"column:role"`
		Status    int32     `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	q := h.db.Table("users")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		q = q.Where("username LIKE ? OR email LIKE ?", like, like)
	}
	if req.Status != nil {
		q = q.Where("status = ?", *req.Status)
	}

	var total int64
	q.Count(&total)

	var rows []row
	q.Select("user_id,username,email,role,status,created_at").
		Order("created_at DESC").
		Offset(req.Offset()).Limit(req.PageSize).
		Scan(&rows)

	items := make([]dto.UserItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.UserItem{
			UserId:    r.UserId,
			Username:  r.Username,
			Email:     r.Email,
			Role:      r.Role,
			Status:    r.Status,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	response.Success(c, dto.UserListResp{Total: total, Users: items})
}

// UserUpdate 修改用户状态或角色
// PUT /api/manage/users/:id
func (h *ManageHandler) UserUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	updates := map[string]any{}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if len(updates) > 0 {
		h.db.Table("users").Where("user_id = ?", id).Updates(updates)
	}
	response.Success(c, nil)
}

// ── 广告 Banner ──────────────────────────────────────

func (h *ManageHandler) BannerList(c *gin.Context) {
	var list []adModel.AdBanner
	h.db.Order("sort DESC, id ASC").Find(&list)
	response.Success(c, gin.H{"banners": list})
}

func (h *ManageHandler) BannerCreate(c *gin.Context) {
	var req dto.BannerSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	b := adModel.AdBanner{Title: req.Title, SubTitle: req.SubTitle, BgStyle: req.BgStyle, Tag: req.Tag, TagClass: req.TagClass, URL: req.URL, Sort: req.Sort, Status: req.Status}
	h.db.Create(&b)
	response.Success(c, b)
}

func (h *ManageHandler) BannerUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.BannerSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	h.db.Model(&adModel.AdBanner{}).Where("id = ?", id).Updates(map[string]any{
		"title": req.Title, "sub_title": req.SubTitle, "bg_style": req.BgStyle,
		"tag": req.Tag, "tag_class": req.TagClass, "url": req.URL,
		"sort": req.Sort, "status": req.Status,
	})
	response.Success(c, nil)
}

func (h *ManageHandler) BannerDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.db.Delete(&adModel.AdBanner{}, id)
	response.Success(c, nil)
}

// ── 广告 App ─────────────────────────────────────────

func (h *ManageHandler) AppList(c *gin.Context) {
	var list []adModel.AdApp
	h.db.Order("sort DESC, id ASC").Find(&list)
	response.Success(c, gin.H{"apps": list})
}

func (h *ManageHandler) AppCreate(c *gin.Context) {
	var req dto.AppSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	a := adModel.AdApp{Name: req.Name, Icon: req.Icon, IconText: req.IconText, BgStyle: req.BgStyle, Official: req.Official, Hot: req.Hot, Categories: req.Categories, URL: req.URL, Sort: req.Sort, Status: req.Status}
	h.db.Create(&a)
	response.Success(c, a)
}

func (h *ManageHandler) AppUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.AppSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	h.db.Model(&adModel.AdApp{}).Where("id = ?", id).Updates(map[string]any{
		"name": req.Name, "icon": req.Icon, "icon_text": req.IconText,
		"bg_style": req.BgStyle, "official": req.Official, "hot": req.Hot,
		"categories": req.Categories, "url": req.URL, "sort": req.Sort, "status": req.Status,
	})
	response.Success(c, nil)
}

func (h *ManageHandler) AppDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.db.Delete(&adModel.AdApp{}, id)
	response.Success(c, nil)
}

// ── 吃瓜文章 ─────────────────────────────────────────

func (h *ManageHandler) GossipList(c *gin.Context) {
	var req dto.GossipListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	q := h.db.Model(&gossipModel.GossipPost{})
	if req.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		q = q.Where("status = ?", *req.Status)
	}

	var total int64
	q.Count(&total)

	var posts []gossipModel.GossipPost
	q.Order("published_at DESC").Offset(req.Offset()).Limit(req.PageSize).Find(&posts)

	items := make([]dto.GossipItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, dto.GossipItem{
			ID:          p.ID,
			Title:       p.Title,
			Cover:       p.Cover,
			ViewCount:   p.ViewCount,
			Status:      p.Status,
			PublishedAt: p.PublishedAt.Format("2006-01-02 15:04"),
		})
	}
	response.Success(c, dto.GossipListResp{Total: total, Posts: items})
}

func (h *ManageHandler) GossipCreate(c *gin.Context) {
	var req dto.GossipSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	pub, _ := time.Parse(time.RFC3339, req.PublishedAt)
	if pub.IsZero() {
		pub = time.Now()
	}
	p := gossipModel.GossipPost{Title: req.Title, Cover: req.Cover, Summary: req.Summary, Content: req.Content, Tags: req.Tags, Status: req.Status, PublishedAt: pub}
	h.db.Create(&p)
	response.Success(c, p)
}

func (h *ManageHandler) GossipUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req dto.GossipSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	pub, _ := time.Parse(time.RFC3339, req.PublishedAt)
	if pub.IsZero() {
		pub = time.Now()
	}
	h.db.Model(&gossipModel.GossipPost{}).Where("id = ?", id).Updates(map[string]any{
		"title": req.Title, "cover": req.Cover, "summary": req.Summary,
		"content": req.Content, "tags": req.Tags, "status": req.Status, "published_at": pub,
	})
	response.Success(c, nil)
}

func (h *ManageHandler) GossipDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.db.Table("gossip_posts").Where("id = ?", id).Update("status", 0)
	response.Success(c, nil)
}
