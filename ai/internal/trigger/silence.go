package trigger

import (
	"context"
	"strings"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
)

var QuestionsWord = []string{
	"吗", "怎么", "怎样", "如何", "是否", "能否",
	"会不会", "有没有", "什么", "哪个", "哪里", "谁",
	"为什么", "为啥", "啥时候", "多长时间", "多久",
}

type SilenceTrigger struct{}

func NewSilenceTrigger() *SilenceTrigger {
	return &SilenceTrigger{}
}

func (t *SilenceTrigger) Check(ctx context.Context, msg *model.Message) (*TriggerEvent, error) {
	if ContainsQuestionWord(msg.Content) {
		duration := time.Duration(config.Trigger.SilenceTimeout) * time.Second
		if err := dao.SetTimer(ctx, config.Trigger.SilenceTimerKey, duration); err != nil {
			return nil, err
		}
		return nil, nil //定时结束后由Watcher触发事件
	}
	if err := dao.DeleteTimer(ctx, config.Trigger.SilenceTimerKey); err != nil {
		return nil, err
	}

	return nil, nil
}

func ContainsQuestionWord(content string) bool {
	for _, word := range QuestionsWord {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

func (t *SilenceTrigger) Watcher(ctx context.Context, eventCh chan<- *TriggerEvent) {
	pubsub := config.REDISDB.PSubscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-pubsub.Channel():
			if msg.Payload != config.Trigger.SilenceTimerKey {
				continue
			}
			eventCh <- &TriggerEvent{
				Type:    TriggerQuestion,
				Message: &model.Message{},
			}
		}
	}
}
