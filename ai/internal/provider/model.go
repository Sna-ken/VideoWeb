package provider

import (
	"context"
	"log"
	"sync"

	"github.com/Sna-ken/videoweb/config"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

var (
	decideOnce    sync.Once
	generateOnce  sync.Once
	decideModel   *openai.ChatModel
	generateModel *openai.ChatModel
)

func DecideModel(ctx context.Context) *openai.ChatModel {
	decideOnce.Do(func() {
		m, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			APIKey:  config.Model.APIKey,
			BaseURL: config.Model.BaseURL,
			Model:   config.Model.DecideModel,
		})
		if err != nil {
			log.Fatalf("failed to create decide model: %v", err)
		}
		decideModel = m
	})
	return decideModel
}

func GenerateModel(ctx context.Context) *openai.ChatModel {
	generateOnce.Do(func() {
		m, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			APIKey:  config.Model.APIKey,
			BaseURL: config.Model.BaseURL,
			Model:   config.Model.GenerateModel,
		})
		if err != nil {
			log.Fatalf("failed to create generate model: %v", err)
		}
		generateModel = m
	})
	return generateModel
}
