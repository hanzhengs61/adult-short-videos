package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adult-short-videos/internal/config"
	"adult-short-videos/internal/data"
	"adult-short-videos/internal/pkg/cache"
	"adult-short-videos/internal/pkg/logger"
	"adult-short-videos/internal/server"
	"adult-short-videos/internal/service/storage/scheduler"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 持有所有依赖，负责应用的完整生命周期。
type App struct {
	cfg  *config.Config
	db   *gorm.DB
	srv  *http.Server
	heat *scheduler.HeatCalculator
}

// New 按顺序初始化所有依赖，任何步骤失败均提前返回错误。
func New(cfg *config.Config) (*App, error) {
	db, err := data.NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}
	logger.Info("✅ 数据库连接成功")

	if err := data.Migrate(db); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}
	logger.Info("✅ 数据库迁移完成")

	if err := cache.Init(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB); err != nil {
		logger.Warn("Redis 连接失败，缓存功能将不可用", zap.Error(err))
	} else {
		logger.Info("✅ Redis 连接成功")
	}

	return &App{
		cfg:  cfg,
		db:   db,
		srv:  server.New(cfg, db),
		heat: scheduler.NewHeatCalculator(db, 5*time.Minute),
	}, nil
}

// Run 启动后台任务与 HTTP 服务器，阻塞直到收到退出信号，然后优雅关闭。
func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.heat.Start(ctx)
	logger.Info("🔥 热度计算器已启动")

	go func() {
		logger.Info("🚀 HTTP 服务器启动", zap.String("addr", a.srv.Addr))
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器异常退出", zap.Error(err))
		}
	}()

	a.printBanner()
	a.waitForShutdown(cancel)
}

func (a *App) waitForShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 正在关闭服务器...")
	cancel()

	// 给正在处理的请求最多 10s 完成
	shutdownCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("HTTP 服务器强制关闭", zap.Error(err))
	}

	if cache.Client != nil {
		_ = cache.Client.Close()
	}

	logger.Info("✅ 服务器已安全关闭")
}

func (a *App) printBanner() {
	fmt.Printf(`
╔════════════════════════════════════════╗
║   🎬 Adult Short Videos Platform      ║
╚════════════════════════════════════════╝
`)
	fmt.Printf("🚀 地址: http://localhost%s\n", a.cfg.Server.Port)
	fmt.Printf("📝 模式: %s\n", a.cfg.Server.Mode)
	fmt.Printf("📊 日志: %s\n", a.cfg.Log.Level)
}
