package dto

// BannerResp 轮播广告响应
type BannerResp struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	SubTitle string `json:"sub_title"`
	BgStyle  string `json:"bg_style"`
	Tag      string `json:"tag"`
	TagClass string `json:"tag_class"`
	URL      string `json:"url"`
}

// AppListReq 广告 App 列表请求
type AppListReq struct {
	// 分类筛选，空或 "all" 表示全部
	Category string `form:"category"`
}

// AppResp 广告 App 响应
type AppResp struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IconText   string `json:"icon_text"`
	BgStyle    string `json:"bg_style"`
	Official   bool   `json:"official"`
	Hot        bool   `json:"hot"`
	Categories string `json:"categories"`
	URL        string `json:"url"`
}
