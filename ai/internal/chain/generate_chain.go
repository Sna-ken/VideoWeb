package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sna-ken/videoweb/ai/internal/provider"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/internal/model"
	"github.com/Sna-ken/videoweb/pkg/utils"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

type GenerateChain struct {
	model *openai.ChatModel
}

func NewGenerateChain(ctx context.Context) *GenerateChain {
	return &GenerateChain{
		model: provider.GenerateModel(ctx),
	}
}

// mention触发逻辑:当收到@时触发。判断是否存在不理解词汇，是则搜索，否则正常回答。如果不理解词汇过多，就当干脆复读机复读ovo
func (g *GenerateChain) HandleMention(ctx context.Context, msg *model.Message, history []*schema.Message, dicision *DecideChain) (string, error) {
	agentname := "@" + config.Agent.Name
	content := strings.TrimSpace(
		strings.ReplaceAll(msg.Content, agentname, ""),
	)
	words, err := dicision.ExtractResearch(ctx, content)
	if err != nil {
		return "", fmt.Errorf("generatechain error:%v", err)
	}

	if len(words) == 0 {
		reply, err := g.generate(ctx, content, history, "")
		return reply, err
	}

	wordMeaning := make(map[string]string)
	unKnowncnt := 0
	for _, word := range words {
		meaning, err := utils.Search(ctx, word)
		if err != nil || meaning == "" {
			unKnowncnt++
			continue
		}
		wordMeaning[word] = meaning
	}

	threshold := config.Trigger.UnknownWordThreshold
	if unKnowncnt > int(threshold) {
		return content, nil
	}

	var s strings.Builder
	s.WriteString("以下是消息中部分词汇的含义供参考：\n")
	for word, meaning := range wordMeaning {
		s.WriteString(fmt.Sprintf("- %s: %s\n", word, meaning))
	}
	reply, err := g.generate(ctx, msg.Content, history, s.String())
	return reply, err
}

// silence触发逻辑:当用户提问且冷场一定时间后回答。判断有trigger进行，确认冷场后直接回复即可
func (g *GenerateChain) HandleQuestion(ctx context.Context, msg *model.Message, history []*schema.Message, decision *DecideChain) (string, error) {
	reply, err := g.generate(ctx, msg.Content, history, "")
	return reply, err
}

// emotion触发逻辑:当用户情绪激烈时触发。分析情绪分数，超过阈值则触发。触发后进行安抚和疏导
func (g *GenerateChain) HandleEmotion(ctx context.Context, msg *model.Message, history []*schema.Message) (string, error) {
	guide := "当前有用户情绪波动较大,请对用户情绪进行安抚和疏导。"
	reply, err := g.generate(ctx, msg.Content, history, guide)
	return reply, err
}

func (g *GenerateChain) generate(ctx context.Context, msg string, history []*schema.Message, wordctx string) (string, error) {
	systemPrompt := fmt.Sprintf("你的名字叫:%s。", config.Agent.Name+config.Agent.Persona)

	if wordctx != "" {
		systemPrompt += "\n\n" + wordctx
	}

	messages := make([]*schema.Message, 0, len(history)+2)
	messages = append(messages, schema.SystemMessage(systemPrompt))
	messages = append(messages, history...)

	if msg == "" {
		messages = append(messages, schema.UserMessage("根据上面的聊天记录，以自然的方式参与对话。"))
	} else {
		messages = append(messages, schema.UserMessage(msg))
	}

	resp, err := g.model.Generate(ctx, messages)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}
