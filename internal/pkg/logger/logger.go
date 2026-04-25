package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Init 初始化日志系统
func Init(level, filePath string, maxSize, maxBackups, maxAge int) error {
	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                        // 时间字段
		LevelKey:       "level",                       // 日志级别字段
		NameKey:        "logger",                      // 日志名称字段
		CallerKey:      "caller",                      // 调用者字段
		MessageKey:     "msg",                         // 消息字段
		StacktraceKey:  "stacktrace",                  // 堆栈字段
		LineEnding:     zapcore.DefaultLineEnding,     // 换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写输出 (如 info)
		EncodeTime:     customTimeEncoder,             // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 调用者编码
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

	// 文件输出
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 创建 Core
	// 统一使用 JSONEncoder 以满足用户要求的 JSON 格式
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, zapLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleWriter, zapLevel),
	)

	// 创建 Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
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
	err := Logger.Sync()
	if err != nil {
		fmt.Println("同步日志失败:", err)
	}
}
