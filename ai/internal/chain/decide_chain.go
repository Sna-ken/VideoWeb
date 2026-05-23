package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sna-ken/videoweb/ai/internal/provider"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

type DecideChain struct {
	model *openai.ChatModel
}

func NewDecideChain(ctx context.Context) *DecideChain {
	return &DecideChain{
		model: provider.DecideModel(ctx),
	}
}

func (d *DecideChain) ExtractResearch(ctx context.Context, content string) ([]string, error) {
	prompt := fmt.Sprintf("判断下面句子中是否含有你不确定含义的词语,如果有,以JSON数组返回这些词,例如[\"yyds\", \"何意味\"]。。如果没有,就返回空数组[]。只返回JSON,不要有其它内容\n\n句子:%s", content)

	resp, err := d.model.Generate(ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("decidechain error:%v", err)
	}

	var words []string
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Content)), &words); err != nil {
		return []string{}, nil
	}

	return words, nil
}
