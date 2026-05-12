package dao

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
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
