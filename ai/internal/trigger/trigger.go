package trigger

import (
	"context"

	"github.com/Sna-ken/videoweb/internal/model"
)

type TriggerType string

const (
	TriggerMention  TriggerType = "mention"
	TriggerQuestion TriggerType = "question"
	TriggerEmotion  TriggerType = "emotion"
)

type TriggerEvent struct {
	Type    TriggerType
	Message *model.Message
}

type Trigger interface {
	Check(ctx context.Context, msg *model.Message) (*TriggerEvent, error)
}
