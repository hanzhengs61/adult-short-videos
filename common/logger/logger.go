package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// Init 初始化日志
// Init 初始化日志系统
func Init(level, filePath string, maxSize, maxBackups, maxAge int) error {
	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                           // 时间字段
		LevelKey:       "level",                          // 日志级别字段
		NameKey:        "logger",                         // 日志名称字段
		CallerKey:      "caller",                         // 调用者字段
		MessageKey:     "msg",                            // 消息字段
		StacktraceKey:  "stacktrace",                     // 堆栈字段
		LineEnding:     zapcore.DefaultLineEnding,        // 换行符
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色输出
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间编码
		EncodeDuration: zapcore.SecondsDurationEncoder,   // 持续时间编码
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 调用者编码
	}

	// 日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 文件输出（使用 lumberjack 实现日志轮转）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize, // MB
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // 天
		Compress:   true,   // 压缩
	})

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 创建 Core
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, zapLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, zapLevel),
	)

	// 创建 Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// getLogWriter 获取日志输出位置
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}

// Debug Info Warn Error Fatal 封装常用方法
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync 同步
func Sync() {
	Logger.Sync()
}
