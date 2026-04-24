package cache

import (
	"adult-short-videos/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Client *redis.Client

// Init 初始化 Redis
func Init(host string, port int, password string, db int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     10,              // 连接池大小
		MinIdleConns: 5,               // 最小空闲连接
		MaxRetries:   3,               // 最大重试次数
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	logger.Info("Redis 连接成功",
		zap.String("addr", fmt.Sprintf("%s:%d", host, port)),
		zap.Int("db", db),
	)

	return nil
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Redis 序列化失败", zap.Error(err), zap.String("key", key))
		return err
	}

	if err := Client.Set(ctx, key, data, expiration).Err(); err != nil {
		logger.Error("Redis 设置失败", zap.Error(err), zap.String("key", key))
		return err
	}

	return nil
}

// Get 获取缓存
func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		logger.Error("Redis 获取失败", zap.Error(err), zap.String("key", key))
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		logger.Error("Redis 反序列化失败", zap.Error(err), zap.String("key", key))
		return err
	}

	return nil
}

// Delete 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := Client.Del(ctx, keys...).Err(); err != nil {
		logger.Error("Redis 删除失败", zap.Error(err), zap.Strings("keys", keys))
		return err
	}

	return nil
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	count, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Increment 自增
func Increment(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// ErrCacheMiss 缓存未命中错误
var ErrCacheMiss = fmt.Errorf("cache miss")
