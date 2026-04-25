package data

import (
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/logger"
	commentModel "adult-short-videos/internal/service/comment/model"
	favoriteModel "adult-short-videos/internal/service/favorite/model"
	playModel "adult-short-videos/internal/service/play/model"
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

	return session.AutoMigrate(
		&userModel.User{},
		&userModel.UserStatistics{},
		&userModel.LoginLog{},
		&videoModel.Video{},
		&videoModel.VideoHeatStats{},
		&videoModel.Actor{},
		&videoModel.VideoActor{},
		&favoriteModel.Favorite{},
		&playModel.PlayHistory{},
		&commentModel.Comment{},
		&commentModel.CommentLike{},
	)
}
