package graph

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel1' in graph 'pghelper'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &deepseek.ChatModelConfig{
		APIKey:      "sk-a76a955533c649d6922a5042af6c0492",
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-reasoner",
		Temperature: 0.5}
	cm, err = deepseek.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

// newChatModel1 component initialization function of node 'ChatModel2' in graph 'pghelper'
func newChatModel1(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &deepseek.ChatModelConfig{
		APIKey:      "sk-a76a955533c649d6922a5042af6c0492",
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-reasoner",
		Temperature: 0.5}
	cm, err = deepseek.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
