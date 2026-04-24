package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adult-short-videos/internal/pkg/logger"

	"go.uber.org/zap"
)

// ShutdownHook 关闭钩子
type ShutdownHook struct {
	Name string
	Fn   func(ctx context.Context) error
}

// Option 配置选项
type Option func(*App)

// WithTimeout 设置整体关闭超时（默认 30s）
func WithTimeout(d time.Duration) Option {
	return func(a *App) { a.timeout = d }
}

// App 应用生命周期管理器
type App struct {
	hooks   []ShutdownHook
	timeout time.Duration
}

// New 创建生命周期管理器
func New(opts ...Option) *App {
	a := &App{timeout: 30 * time.Second}
	for _, o := range opts {
		o(a)
	}
	return a
}

// RegisterShutdownHook 注册关闭钩子，按注册顺序串行执行
func (a *App) RegisterShutdownHook(name string, fn func(ctx context.Context) error) {
	a.hooks = append(a.hooks, ShutdownHook{Name: name, Fn: fn})
}

// Run 阻塞直到收到 SIGINT/SIGTERM，然后串行执行所有关闭钩子
func (a *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 收到关闭信号，开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	for _, hook := range a.hooks {
		logger.Info("正在关闭", zap.String("component", hook.Name))
		if err := hook.Fn(ctx); err != nil {
			logger.Error("关闭失败", zap.String("component", hook.Name), zap.Error(err))
		} else {
			logger.Info("已关闭", zap.String("component", hook.Name))
		}
	}

	logger.Info("✅ 服务器已优雅关闭")
	logger.Sync()
}
