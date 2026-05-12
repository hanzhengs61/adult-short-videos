package data

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/logger"
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
	); err != nil {
		return err
	}

	return seedTags(session)
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
