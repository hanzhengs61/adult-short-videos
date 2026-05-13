package model

import "time"

// AdBanner 广告轮播图表
type AdBanner struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"type:varchar(300);not null;comment:主标题，支持简单 HTML 标签（如 <span>）" json:"title"`
	SubTitle string `gorm:"type:varchar(500);comment:副标题" json:"sub_title"`
	// CSS background 完整字符串，如 "background:linear-gradient(...)"
	BgStyle  string `gorm:"type:varchar(500);comment:CSS background 样式字符串" json:"bg_style"`
	Tag      string `gorm:"type:varchar(50);comment:左上角角标文字，空则不显示" json:"tag"`
	TagClass string `gorm:"type:varchar(200);comment:角标 Tailwind CSS 类" json:"tag_class"`
	URL      string `gorm:"type:varchar(500);comment:点击跳转链接" json:"url"`
	Sort     int32  `gorm:"default:0;comment:排序权重，越大越靠前" json:"sort"`
	Status   int32  `gorm:"default:1;index;comment:0:下架 1:正常" json:"status"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AdBanner) TableName() string { return "ad_banners" }

// AdApp 广告应用（图标网格中每一项）
type AdApp struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null;comment:应用名称" json:"name"`
	Icon     string `gorm:"type:varchar(50);comment:图标 emoji" json:"icon"`
	IconText string `gorm:"type:varchar(50);comment:图标区域内的辅助文字" json:"icon_text"`
	// CSS background 完整字符串
	BgStyle  string `gorm:"type:varchar(500);comment:图标背景 CSS background 样式" json:"bg_style"`
	Official bool   `gorm:"default:false;comment:是否显示官方角标" json:"official"`
	Hot      bool   `gorm:"default:false;comment:是否显示 HOT 角标" json:"hot"`
	// 逗号分隔的分类标识，固定值: all / video / chat / gay / live / manga / game
	Categories string `gorm:"type:varchar(200);default:'all';comment:所属分类，逗号分隔" json:"categories"`
	URL        string `gorm:"type:varchar(500);comment:点击跳转链接" json:"url"`
	Sort       int32  `gorm:"default:0;comment:排序权重，越大越靠前" json:"sort"`
	Status     int32  `gorm:"default:1;index;comment:0:下架 1:正常" json:"status"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AdApp) TableName() string { return "ad_apps" }
