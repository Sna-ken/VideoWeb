package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/cloudwego/eino/schema"
)

func BuildContext(ctx context.Context) ([]*schema.Message, error) {
	msgs, err := dao.FindMessageByRedis(ctx, config.WS.ChatRedisKey, 0, config.WS.MaxMessage)
	if err != nil || len(*msgs) == 0 {
		return nil, nil
	}

	messages := make([]*model.Message, 0, len(*msgs))
	for _, data := range *msgs {
		var msg model.Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			continue
		}
		messages = append(messages, &msg)
	}

	slices.Reverse(messages)

	agentname := config.Agent.Name
	result := make([]*schema.Message, 0, len(messages))

	for _, m := range messages {
		if m.Username == agentname {
			result = append(result, schema.AssistantMessage(m.Content, nil))
		} else {
			result = append(result, schema.UserMessage(fmt.Sprintf("%s:%s", m.Username, m.Content)))
		}
	}
	return result, nil
}
