package data

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/logger"
	adModel "adult-short-videos/internal/service/ad/model"
	commentModel "adult-short-videos/internal/service/comment/model"
	favoriteModel "adult-short-videos/internal/service/favorite/model"
	followModel "adult-short-videos/internal/service/follow/model"
	gossipModel "adult-short-videos/internal/service/gossip/model"
	playModel "adult-short-videos/internal/service/play/model"
	tagModel "adult-short-videos/internal/service/tag/model"
	userModel "adult-short-videos/internal/service/user/model"
	videoModel "adult-short-videos/internal/service/video/model"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := gormLogger.Info
	if cfg.Server.Mode == "release" {
		logLevel = gormLogger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info("数据库连接池配置",
		zap.Int("max_idle", cfg.Database.MaxIdleConns),
		zap.Int("max_open", cfg.Database.MaxOpenConns),
		zap.Int("max_lifetime", cfg.Database.ConnMaxLifetime),
	)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	// 在迁移时临时关闭 SQL 日志打印，避免输出大量系统元数据查询
	session := db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(gormLogger.Silent),
	})

	if err := session.AutoMigrate(
		&userModel.User{},
		&userModel.UserStatistics{},
		&userModel.LoginLog{},
		&videoModel.Video{},
		&videoModel.VideoHeatStats{},
		&favoriteModel.Favorite{},
		&playModel.PlayHistory{},
		&commentModel.Comment{},
		&commentModel.CommentLike{},
		&followModel.AuthorFollow{},
		&gossipModel.GossipPost{},
		&tagModel.Tag{},
		&adModel.AdBanner{},
		&adModel.AdApp{},
	); err != nil {
		return err
	}

	if err := seedTags(session); err != nil {
		return err
	}
	return seedAds(session)
}

// seedTags 从 seed_tags.txt 读取标签写入数据库，已存在的跳过（幂等）
func seedTags(db *gorm.DB) error {
	// seed_tags.txt 与本文件同目录，每行一个标签
	f, err := os.Open("internal/data/seed_tags.txt")
	if err != nil {
		return fmt.Errorf("打开 seed_tags.txt 失败: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}
		tag := tagModel.Tag{Name: name, Status: 1}
		if err := db.Where(tagModel.Tag{Name: name}).FirstOrCreate(&tag).Error; err != nil {
			return err
		}
	}
	return scanner.Err()
}

// seedAds 首次启动时写入默认广告数据，表不为空则跳过（幂等）
func seedAds(db *gorm.DB) error {
	var count int64
	db.Model(&adModel.AdBanner{}).Count(&count)
	if count > 0 {
		return nil
	}

	banners := []adModel.AdBanner{
		{Title: `送你 <span style="color:#FFE600">1888</span> · 1元爆大奖`, SubTitle: "PG电子 · 扶持放水 · PG72.COM", BgStyle: "background:linear-gradient(135deg,#0f0080 0%,#3a00cc 40%,#7B2FFF 100%)", Tag: "赌场", TagClass: "bg-yellow-400 text-black", URL: "#", Sort: 40, Status: 1},
		{Title: `注册送 <span style="color:#FFD700">888</span> 体验金`, SubTitle: "开元棋牌 · 官方正版 · 8319.COM", BgStyle: "background:linear-gradient(135deg,#8B0000 0%,#c0392b 50%,#e74c3c 100%)", Tag: "官方", TagClass: "bg-pink-500 text-white", URL: "#", Sort: 30, Status: 1},
		{Title: "全国寂寞小姐姐\n等你来撩", SubTitle: "约炮裸聊 · 真人在线 · 免费注册", BgStyle: "background:linear-gradient(135deg,#1a0533 0%,#6B21A8 60%,#FF2080 100%)", Tag: "直播", TagClass: "bg-blue-500 text-white", URL: "#", Sort: 20, Status: 1},
		{Title: `福利APP下载免费领 <span style="color:#FF2080">VIP</span>`, SubTitle: "海量内容 · 每日更新 · 无需翻墙", BgStyle: "background:linear-gradient(135deg,#003322 0%,#006644 50%,#00aa66 100%)", Tag: "免费", TagClass: "bg-emerald-500 text-white", URL: "#", Sort: 10, Status: 1},
	}
	if err := db.Create(&banners).Error; err != nil {
		return err
	}

	apps := []adModel.AdApp{
		{Name: "巨乳约炮", Icon: "🔞", BgStyle: "background:linear-gradient(135deg,#1a0010,#7B2FFF)", Hot: true, Categories: "all,chat", URL: "#", Sort: 100, Status: 1},
		{Name: "激情迷药", Icon: "💊", IconText: "情趣", BgStyle: "background:linear-gradient(135deg,#003366,#0066cc)", Categories: "all,chat", URL: "#", Sort: 90, Status: 1},
		{Name: "PG电子", Icon: "🎰", IconText: "PG", BgStyle: "background:linear-gradient(135deg,#0f0080,#7B2FFF)", Official: true, Hot: true, Categories: "all,game", URL: "#", Sort: 80, Status: 1},
		{Name: "棋牌电子", Icon: "🃏", IconText: "888", BgStyle: "background:linear-gradient(135deg,#003399,#0055cc)", Official: true, Categories: "all,game", URL: "#", Sort: 70, Status: 1},
		{Name: "3223棋牌", Icon: "♠️", IconText: "KY", BgStyle: "background:linear-gradient(135deg,#8B0000,#cc0000)", Official: true, Categories: "all,game", URL: "#", Sort: 60, Status: 1},
		{Name: "开元棋牌", Icon: "🎲", IconText: "KY", BgStyle: "background:linear-gradient(135deg,#333300,#999900)", Official: true, Categories: "all,game", URL: "#", Sort: 50, Status: 1},
		{Name: "PG放水", Icon: "💸", IconText: "大放水", BgStyle: "background:linear-gradient(135deg,#000000,#1a1a1a)", Hot: true, Categories: "all,game", URL: "#", Sort: 45, Status: 1},
		{Name: "澳门赌场", Icon: "🎪", IconText: "翻倍", BgStyle: "background:linear-gradient(135deg,#660000,#cc3300)", Categories: "all,game", URL: "#", Sort: 40, Status: 1},
		{Name: "77直播", Icon: "📺", IconText: "77", BgStyle: "background:linear-gradient(135deg,#8B0000,#FF2080)", Hot: true, Categories: "all,live", URL: "#", Sort: 35, Status: 1},
		{Name: "澳门新葡京", Icon: "🏯", IconText: "3A", BgStyle: "background:linear-gradient(135deg,#1a1a00,#4a4a00)", Official: true, Categories: "all,game", URL: "#", Sort: 30, Status: 1},
		{Name: "金沙直播", Icon: "🌟", IconText: "金沙", BgStyle: "background:linear-gradient(135deg,#660033,#cc0066)", Categories: "all,live", URL: "#", Sort: 25, Status: 1},
		{Name: "国产精品", Icon: "🎬", BgStyle: "background:linear-gradient(135deg,#003300,#006600)", Hot: true, Categories: "all,video", URL: "#", Sort: 20, Status: 1},
		{Name: "91约炮", Icon: "💬", IconText: "约炮", BgStyle: "background:linear-gradient(135deg,#1a0033,#4a0099)", Hot: true, Categories: "all,chat", URL: "#", Sort: 15, Status: 1},
		{Name: "男同专区", Icon: "👬", BgStyle: "background:linear-gradient(135deg,#003366,#006699)", Categories: "all,gay", URL: "#", Sort: 10, Status: 1},
		{Name: "成人漫画", Icon: "📖", IconText: "H漫", BgStyle: "background:linear-gradient(135deg,#330033,#660066)", Categories: "all,manga", URL: "#", Sort: 8, Status: 1},
		{Name: "美女直播", Icon: "💃", BgStyle: "background:linear-gradient(135deg,#4a0020,#FF2080)", Hot: true, Categories: "all,live", URL: "#", Sort: 6, Status: 1},
		{Name: "麻豆传媒", Icon: "🎭", IconText: "麻豆", BgStyle: "background:linear-gradient(135deg,#1a0000,#660000)", Official: true, Hot: true, Categories: "all,video", URL: "#", Sort: 5, Status: 1},
		{Name: "新葡京", Icon: "🎯", IconText: "024", BgStyle: "background:linear-gradient(135deg,#660000,#990000)", Categories: "all,game", URL: "#", Sort: 4, Status: 1},
		{Name: "裸聊平台", Icon: "📱", IconText: "裸聊", BgStyle: "background:linear-gradient(135deg,#001a33,#0066cc)", Categories: "all,chat", URL: "#", Sort: 3, Status: 1},
		{Name: "同志交友", Icon: "🌈", BgStyle: "background:linear-gradient(135deg,#003344,#006688)", Categories: "gay", URL: "#", Sort: 2, Status: 1},
	}
	return db.Create(&apps).Error
}
