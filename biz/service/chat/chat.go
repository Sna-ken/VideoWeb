package chat

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/google/uuid"
)

func SaveMessage(ctx context.Context, userID string, username string, content string) (*model.Message, error) {
	msg := &model.Message{
		ID:        uuid.New().String(),
		UserID:    userID,
		Username:  username,
		Content:   content,
		Create_at: time.Now(),
	}
	// mysql
	go func() {
		if err := dao.CreateMessage(ctx, msg); err != nil {
			log.Printf("save message to mysql error:%v", err)
		}
	}()

	// redis
	data, err := json.Marshal(msg)
	if err != nil {
		return msg, nil
	}

	if err := dao.SaveMessageToRedis(ctx, config.WS.ChatRedisKey, 0, config.WS.MaxMessage-1, data); err != nil {
		log.Printf("save message to redis error:%v", err)
	}

	return msg, nil
}

func GetRecentMessages(ctx context.Context) ([]*model.Message, error) {
	// redis
	rsl, err := dao.FindMessageByRedis(ctx, config.WS.ChatRedisKey, 0, config.WS.MaxMessage)
	if err == nil && len(*rsl) > 0 {
		messages := make([]*model.Message, 0, len(*rsl))
		for _, data := range *rsl {
			var msg model.Message
			if err := json.Unmarshal([]byte(data), &msg); err != nil {
				continue
			}
			messages = append(messages, &msg)
		}
		// 反转为时间正序
		slices.Reverse(messages)
		return messages, nil
	}

	// mysql
	var messages []*model.Message
	if err := dao.FindMessageByMysql(ctx, &messages, int(config.WS.MaxMessage)); err != nil {
		return nil, err
	}
	slices.Reverse(messages)

	return messages, nil
}
