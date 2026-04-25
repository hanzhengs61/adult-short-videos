package main

import (
	"fmt"
	"os"

	"adult-short-videos/internal/app"
	"adult-short-videos/internal/config"
	"adult-short-videos/internal/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(
		cfg.Log.Level,
		cfg.Log.FilePath,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
	); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("🚀 应用启动", zap.String("mode", cfg.Server.Mode))

	a, err := app.New(cfg)
	if err != nil {
		logger.Fatal("应用初始化失败", zap.Error(err))
	}

	a.Run()
}
