package agent

import (
	"context"
	"log"
	"time"

	"github.com/Sna-ken/videoweb/ai/internal/chain"
	"github.com/Sna-ken/videoweb/ai/internal/memory"
	"github.com/Sna-ken/videoweb/ai/internal/trigger"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/internal/model"
)

var (
	DecisionChain   *chain.DecideChain
	GenerationChain *chain.GenerateChain
	MentionTrigger  *trigger.MentionTrigger
	SilenceTrigger  *trigger.SilenceTrigger
	EmotionTrigger  *trigger.EmotionTrigger
)

func InitAgent(ctx context.Context) {
	DecisionChain = chain.NewDecideChain(ctx)
	GenerationChain = chain.NewGenerateChain(ctx)
	MentionTrigger = trigger.NewMentionTrigger()
	SilenceTrigger = trigger.NewSilenceTrigger()
	EmotionTrigger = trigger.NewEmotionTrigger()

	eventCh := make(chan *trigger.TriggerEvent, 10)
	go func() {
		SilenceTrigger.Watcher(ctx, eventCh)
	}()
	go func() {
		for envent := range eventCh {
			HandleEvent(ctx, envent, DecisionChain, GenerationChain)
		}
	}()
}

func OneNewMessage(
	ctx context.Context,
	msg *model.Message,
	mentionTrigger *trigger.MentionTrigger,
	emotionTrigger *trigger.EmotionTrigger,
	silenceTrigger *trigger.SilenceTrigger,
	decideChain *chain.DecideChain,
	generateChain *chain.GenerateChain,
) {
	if msg.UserID == config.Agent.UserID {
		return
	}

	limited, err := dao.IsRateLimited(ctx, config.Trigger.RateLimitKey)
	if err != nil {
		log.Printf("check rate limit failed:%v", err)
		return
	}
	if limited {
		return
	}

	triggers := []trigger.Trigger{mentionTrigger, silenceTrigger, emotionTrigger}
	for _, t := range triggers {
		event, err := t.Check(ctx, msg)
		if err != nil {
			log.Printf("trigger check failed:%v", err)
			continue
		}
		if event != nil {
			HandleEvent(ctx, event, decideChain, generateChain)
			return
		}
	}
}

func HandleEvent(ctx context.Context, event *trigger.TriggerEvent,
	decisionChain *chain.DecideChain, generationChain *chain.GenerateChain) {

	history, err := memory.BuildContext(ctx)
	if err != nil {
		log.Printf("get memory failed:%v", err)
	}

	var reply string

	switch event.Type {
	case trigger.TriggerMention:
		reply, err = generationChain.HandleMention(ctx, event.Message, history, decisionChain)

	case trigger.TriggerQuestion:
		reply, err = generationChain.HandleQuestion(ctx, event.Message, history, decisionChain)

	case trigger.TriggerEmotion:
		reply, err = generationChain.HandleEmotion(ctx, event.Message, history)
	}

	if err != nil {
		log.Printf("generate reply failed:%v", err)
		return
	}
	if reply == "{}" {
		return
	}

	SendReply(ctx, reply)

	limit := time.Duration(config.Trigger.RateLimit) * time.Second
	if err := dao.SetTimer(ctx, config.Trigger.RateLimitKey, limit); err != nil {
		log.Printf("set rate limit failed:%v", err)
	}
}
