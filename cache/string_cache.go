package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var defaultCacheExpirationDuration = time.Hour * 24 // 缓存过期时间

func SetStringCache(rdb redis.Cmdable, ctx context.Context, key string, value interface{}) error {
	_, err := rdb.Set(ctx, key, value, -1).Result()
	return err
}

func SetStringCacheEx(rdb redis.Cmdable, ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := rdb.SetEX(ctx, key, value, expiration).Result()
	return err
}

func SetStringCacheWithDefaultExpiration(rdb redis.Cmdable, ctx context.Context, key string, value interface{}) error {
	_, err := rdb.SetEX(ctx, key, value, defaultCacheExpirationDuration).Result()
	return err
}

// 如果没有找到数据，返回 redis.Nil 错误
func GetStringCache(rdb redis.Cmdable, ctx context.Context, key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func GetStringCacheValue(rdb redis.Cmdable, ctx context.Context, key string, val interface{}) error {
	err := rdb.Get(ctx, key).Scan(&val)
	if err != nil {
		return err
	}

	return nil
}

func DeleteStringCache(rdb redis.Cmdable, ctx context.Context, key string) error {
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
