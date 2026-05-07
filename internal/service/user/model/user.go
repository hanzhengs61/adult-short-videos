package model

import "time"

type User struct {
	UserId            int64      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"user_id"`
	Username          string     `gorm:"uniqueIndex;type:varchar(50);not null;comment:用户名" json:"username"`
	Password          string     `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email             string     `gorm:"uniqueIndex;type:varchar(100);comment:邮箱" json:"email"`
	Avatar            string     `gorm:"type:varchar(255);comment:头像 URL" json:"avatar"`
	Bio               string     `gorm:"type:varchar(500);comment:个人简介" json:"bio"`
	Role              string     `gorm:"type:varchar(20);default:user;comment:角色：user/creator" json:"role"`
	Status            int32      `gorm:"default:1;comment:账号状态：0=禁用，1=正常" json:"status"`
	UsernameChangedAt *time.Time `gorm:"type:timestamp;comment:昵称最后修改时间" json:"username_changed_at"`
	AvatarChangedAt   *time.Time `gorm:"type:timestamp;comment:头像最后修改时间" json:"avatar_changed_at"`
	LastLoginAt       time.Time  `gorm:"type:timestamp;comment:最后登录时间" json:"last_login_at"`
	LastLoginIp       string     `gorm:"type:varchar(50);comment:最后登录 IP" json:"last_login_ip"`
	CreatedAt         time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserStatistics 用户统计表
// 用于记录用户的行为统计数据
type UserStatistics struct {
	UserId         int64     `gorm:"primaryKey;comment:用户ID" json:"user_id"`
	PlayCount      int64     `gorm:"default:0;comment:播放次数" json:"play_count"`
	FavoriteCount  int64     `gorm:"default:0;comment:收藏数" json:"favorite_count"`
	SubscribeCount int64     `gorm:"default:0;comment:订阅数" json:"subscribe_count"`
	CommentCount   int64     `gorm:"default:0;comment:评论数" json:"comment_count"`
	TotalWatchTime int64     `gorm:"default:0;comment:总观看时长(秒)" json:"total_watch_time"`
	LastActiveAt   time.Time `gorm:"type:timestamp;comment:最后活跃" json:"last_active_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (UserStatistics) TableName() string {
	return "user_statistics"
}

// LoginLog 登录日志表
// 记录每次登录的详细信息
type LoginLog struct {
	LogId      int64     `gorm:"primaryKey;autoIncrement;comment:登录日志ID" json:"log_id"`
	UserId     int64     `gorm:"index;not null;comment:用户ID" json:"user_id"`
	LoginType  string    `gorm:"type:varchar(20);comment:登录类型：password,sms,oauth" json:"login_type"`
	LoginIp    string    `gorm:"type:varchar(50);comment:登录IP" json:"login_ip"`
	UserAgent  string    `gorm:"type:varchar(255);comment:用户代理" json:"user_agent"`
	Status     int32     `gorm:"default:1;comment:状态：0=失败,1=成功" json:"status"`
	FailReason string    `gorm:"type:varchar(255);comment:失败原因" json:"fail_reason"`
	CreatedAt  time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}
