package dao

import (
	"context"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/redis/go-redis/v9"
)

func CreateMessage(ctx context.Context, msg *model.Message) error {
	return config.MYSQLDB.WithContext(ctx).Create(msg).Error
}

func SaveMessageToRedis(ctx context.Context, rediskey string, start int64, end int64, data []byte) error {
	pipe := config.REDISDB.Pipeline()
	pipe.LPush(ctx, rediskey, data)
	pipe.LTrim(ctx, rediskey, start, end)
	_, err := pipe.Exec(ctx)
	return err
}
func FindMessageByRedis(ctx context.Context, rediskey string, start int64, end int64) (*[]string, error) {
	rsl, err := config.REDISDB.LRange(ctx, rediskey, start, end).Result()
	return &rsl, err
}

func FindMessageByMysql(ctx context.Context, messages *[]*model.Message, limit int) error {
	return config.MYSQLDB.WithContext(ctx).Order("create_at DESC").Limit(limit).Find(&messages).Error
}

func SetTimer(ctx context.Context, key string, duration time.Duration) error {
	return config.REDISDB.Set(ctx, key, 1, duration).Err()
}

func DeleteTimer(ctx context.Context, key string) error {
	return config.REDISDB.Del(ctx, key).Err()
}

func IsTimeout(ctx context.Context, key string) (bool, error) {
	err := config.REDISDB.Get(ctx, key).Err()
	if err == redis.Nil {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func IsRateLimited(ctx context.Context, key string) (bool, error) {
	err := config.REDISDB.Get(ctx, key).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
