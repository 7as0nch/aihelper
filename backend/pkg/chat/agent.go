package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 21:40:12
 * @Description:
**/

type AiAgent struct {
	// agent *chain.Chain
	chatModel *openai.ChatModel
	ctx       context.Context
}

func NewAiAgent() *AiAgent {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")

	if openAIAPIKey == "" {
		openAIAPIKey = "sk-a76a955533c649d6922a5042af6c0492"
	}
	if openAIModelName == "" {
		openAIModelName = "deepseek-chat"
		// openAIModelName = "deepseek-reasoner"
	}
	if openAIBaseURL == "" {
		openAIBaseURL = "https://api.deepseek.com"
	}
	// var ptr float32 = 0.7
	ctx := context.Background()
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     openAIBaseURL,
		Model:       openAIModelName,
		APIKey:      openAIAPIKey,
		// Temperature: &ptr,
	})
	if err != nil {
		return nil
	}
	return &AiAgent{
		chatModel: chatModel,
		ctx:       ctx,
	}
}

func (a *AiAgent) GetChatModel() *openai.ChatModel {
	return a.chatModel
}

func (a *AiAgent) Stream(ctx context.Context, in []*schema.Message, opts ...model.Option) (outStream *schema.StreamReader[*schema.Message], err error) {
	return a.chatModel.Stream(ctx, in, opts...)
}

func (a *AiAgent) GetAgent() compose.Runnable[[]*schema.Message, []*schema.Message] {
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(a.chatModel, compose.WithNodeName("chat_model"))
		// AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(a.ctx)
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("chain.Compile failed, err=%v", err)
		return nil
	}
	return agent
}
