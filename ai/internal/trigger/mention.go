package trigger

import (
	"context"
	"strings"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
)

type MentionTrigger struct{}

func NewMentionTrigger() *MentionTrigger {
	return &MentionTrigger{}
}

func (t *MentionTrigger) Check(ctx context.Context, msg *model.Message) (*TriggerEvent, error) {
	agentname := config.Agent.Name

	mention := "@" + agentname
	if !strings.Contains(msg.Content, mention) {
		return nil, nil
	}

	return &TriggerEvent{
		Type:    TriggerMention,
		Message: msg,
	}, nil
}
