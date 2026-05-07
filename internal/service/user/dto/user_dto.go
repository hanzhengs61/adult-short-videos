package dto

// LoginReq 登录请求参数
type LoginReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

// LoginResp 登录响应
type LoginResp struct {
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RegisterReq 注册请求（添加验证标签）
type RegisterReq struct {
	Username   string `json:"username" validate:"required,username,min=3,max=20"`
	Password   string `json:"password" validate:"required,min=6,max=50"`
	Email      string `json:"email" validate:"required,email"`
	InviteCode string `json:"invite_code"`
}

// RegisterResp 注册响应
type RegisterResp struct {
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Role     string `json:"role"`
	Status   int32  `json:"status"`
}

// UpdateProfileReq 更新用户资料请求
type UpdateProfileReq struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// UpdateProfileResp 更新用户资料响应
type UpdateProfileResp struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// CreatorItem 创作者榜单项
type CreatorItem struct {
	UserId     int64   `json:"user_id"`
	Username   string  `json:"username"`
	Avatar     string  `json:"avatar"`
	VideoCount int64   `json:"video_count"`
	TotalPlays int64   `json:"total_plays"`
	Score      float64 `json:"score"`
}
