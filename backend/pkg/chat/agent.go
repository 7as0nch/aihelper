package chat

import (
	"context"
	"errors"
	"io"

	// "errors"
	"fmt"
	// "io"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	_ "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/pkg/agenttools"
	// tmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 21:40:12
 * @Description:
**/

type AiAgent struct {
	// agent *chain.Chain
	chatModel model.ToolCallingChatModel
	ctx       context.Context
}

func NewAiAgent() *AiAgent {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")

	if openAIAPIKey == "" {
		openAIAPIKey = "9292fee0-90fb-4739-86ff-eb1e886e2823"
	}
	if openAIModelName == "" {
		// openAIModelName = "deepseek-chat"
		openAIModelName = "doubao-seed-1.6-250615"
	}
	if openAIBaseURL == "" {
		openAIBaseURL = "https://api.deepseek.com"
	}
	var ptr float32 = 0.7
	ctx := context.Background()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		// BaseURL: openAIBaseURL,
		Model:  openAIModelName,
		APIKey: openAIAPIKey,
		// Thinking: &tmodel.Thinking{
		// 	Type: tmodel.ThinkingTypeEnabled,
		// },
		Temperature: &ptr,
	})
	if err != nil {
		return nil
	}
	return &AiAgent{
		chatModel: chatModel,
		ctx:       ctx,
	}
}

func NewDeepseekAgent() *AiAgent {
	ctx := context.Background()
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-reasoner",
		APIKey:      "sk-a76a955533c649d6922a5042af6c0492",
		Temperature: 0,
	})
	if err != nil {
		return nil
	}
	var toolinfos []*schema.ToolInfo
	for _, v := range agenttools.GetTools() {
		var toolinfo *schema.ToolInfo
		toolinfo, err = v.Info(ctx)
		if err != nil {
			fmt.Println(err)
			// logs.Errorf("GetToolInfo failed, err=%v", err)
			return nil
		}
		toolinfos = append(toolinfos, toolinfo)
	}
	toolmodel, err := chatmodel.WithTools(toolinfos)
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("WithTools failed, err=%v", err)
		return nil
	}
	return &AiAgent{
		chatModel: toolmodel,
		ctx:       ctx,
	}
}

func (a *AiAgent) GetChatModel() model.ToolCallingChatModel {
	return a.chatModel
}

// func (a *AiAgent) Stream(ctx context.Context, in []*schema.Message, opts ...model.Option) (outStream *schema.StreamReader[*schema.Message], err error) {
// 	// 加入PG助手提示
// 	log.Info("[stream] stream with pg helper prompt")
// 	return a.chatModel.Stream(ctx,
// 		append([]*schema.Message{PgHelperPrompt()}, in...), opts...)
// }

func (a *AiAgent) GetChainAgent() compose.Runnable[[]*schema.Message, *schema.Message] {
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	var toolinfos []*schema.ToolInfo
	for _, v := range agenttools.GetTools() {
		toolinfo, err := v.Info(a.ctx)
		if err != nil {
			fmt.Println(err)
			// logs.Errorf("GetToolInfo failed, err=%v", err)
			return nil
		}
		toolinfos = append(toolinfos, toolinfo)
	}
	toolmodel, err := a.chatModel.WithTools(toolinfos)
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("WithTools failed, err=%v", err)
		return nil
	}
	chain.
		AppendChatTemplate(
			prompt.FromMessages(schema.FString,
				schema.SystemMessage(PgHelperPrompt().Content))).
		AppendChatModel(toolmodel, compose.WithNodeName("chat_model"))
		// AppendToolsNode(agenttools.GetTools(), compose.WithNodeName("tools"))
	toolsNode, err := compose.NewToolNode(a.ctx, &compose.ToolsNodeConfig{
		Tools: agenttools.GetTools(),
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewToolNode failed, err=%v", err)
		return nil
	}
	chain.AppendToolsNode(toolsNode)

	// 编译并运行 chain
	agent, err := chain.Compile(a.ctx)
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("chain.Compile failed, err=%v", err)
		return nil
	}
	return agent
}

func (a *AiAgent) GetReActAgent() *react.Agent {
	// toolCallChecker
	agent, err := react.NewAgent(a.ctx, &react.AgentConfig{
		ToolCallingModel: a.chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: agenttools.GetTools(),
		},
		MessageModifier: func(ctx context.Context, msgs []*schema.Message) []*schema.Message {
			return append([]*schema.Message{PgHelperPrompt()}, msgs...)
		},
		StreamToolCallChecker: toolchecker,
		MaxStep:               10,
		// ToolReturnDirectly: map[string]struct{}{
		// 	// "add_todo": {},
		// 	"list_todo": {},
		// },
	})
	if err != nil {
		fmt.Printf("failed to create agent: %v\n", err)
		return nil
	}
	return agent
}

func (a *AiAgent) GetDeepseekReasoningContent(msg *schema.Message) (string, bool) {
	return deepseek.GetReasoningContent(msg)
}

func toolchecker(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()
	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}

			return false, err
		}

		if len(msg.ToolCalls) > 0 {
			return true, nil
		}
		// msg, err := sr.Recv()
		// if errors.Is(err, io.EOF) {
		// 	return false, nil
		// }
		// if err != nil {
		// 	return false, err
		// }

		// if len(msg.ToolCalls) > 0 {
		// 	return true, nil
		// }

		// if len(msg.Content) == 0 { // skip empty chunks at the front
		// 	continue
		// }

		// return false, nil
	}
	return false, nil
}
