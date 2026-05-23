package trigger

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sna-ken/videoweb/ai/internal/provider"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/cloudwego/eino/schema"
)

type EmotionTrigger struct{}

func NewEmotionTrigger() *EmotionTrigger {
	return &EmotionTrigger{}
}

func (t *EmotionTrigger) Check(ctx context.Context, msg *model.Message) (*TriggerEvent, error) {
	score, err := AnalyzeEomotion(ctx, msg.Content)
	if err != nil {
		return nil, err
	}

	threshold := config.Trigger.EmotionThreshold
	if score < threshold {
		return nil, nil
	}

	return &TriggerEvent{
		Type:    TriggerEmotion,
		Message: msg,
	}, nil
}

func AnalyzeEomotion(ctx context.Context, content string) (float64, error) {
	prompt := fmt.Sprintf("下面这句话的情绪波动程度。0表示完全平静,1表示情绪极度激烈。只返回一个0-1之间的小数,不要有任何其它内容。\n\n句子:%s", content)

	model := provider.DecideModel(ctx)
	resp, err := model.Generate(ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return 0, fmt.Errorf("emotion analyze failed:%v", err)
	}

	score, err := strconv.ParseFloat(strings.TrimSpace(resp.Content), 64)
	if err != nil {
		return 0, fmt.Errorf("parse emotion score failed:%v", err)
	}
	return score, nil
}
