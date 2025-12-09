/*
 * @Author: chengjiang
 * @Date: 2025-12-09
 * @Description: DeepADK 适配器 (eino/adk/prebuilt/deep)
 */
package adapter

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/ai/agent"
	"github.com/example/aichat/backend/pkg/ai/chatmodel"
	"github.com/example/aichat/backend/pkg/ai/prints"
	"github.com/example/aichat/backend/pkg/ai/tool"
)

// DeepAdkAdapter DeepADK 适配器
type DeepAdkAdapter struct {
	config    *ai.AgentConfig
	agent     adk.Agent
	subAgents []ai.Agent
}

// NewDeepAdkAdapter 创建 DeepADK 适配器
func NewDeepAdkAdapter(ctx context.Context, config *ai.AgentConfig, subAgents []ai.Agent) (ai.Agent, error) {
	// 创建 ChatModel
	cm, err := chatmodel.NewModel(ctx, chatmodel.ModelConfig{
		ModelType: chatmodel.ModelType(config.ModelConfig.ModelType),
		ModelName: config.ModelConfig.ModelName,
		ApiKey:    config.ModelConfig.APIKey,
		BaseURL:   config.ModelConfig.BaseURL,
		Thinking:  config.ModelConfig.Thinking,
	},
		chatmodel.WithMaxTokens(config.ModelConfig.MaxTokens),
		chatmodel.WithTemperature(config.ModelConfig.Temperature),
		chatmodel.WithTopP(config.ModelConfig.TopP),
		chatmodel.WithDisableThinking(!config.ModelConfig.Thinking),
	)
	if err != nil {
		return nil, err
	}

	// 转换子 Agent 为 adk.Agent
	var adkSubAgents []adk.Agent
	for _, sub := range subAgents {
		if adapter, ok := sub.(*EinoAdapter); ok {
			adkSubAgents = append(adkSubAgents, adapter.GetInternalAgent())
		}
		if adapter, ok := sub.(*DeepAdkAdapter); ok {
			adkSubAgents = append(adkSubAgents, adapter.GetInternalAgent())
		}
	}
	if config.WithWebSearchAgent {
		webSearchAgent, err := agent.NewWebSearchAgent(ctx, config)
		if err != nil {
			return nil, err
		}
		adkSubAgents = append(adkSubAgents, webSearchAgent)
	}

	// 创建 deep Agent
	deepConfig := &deep.Config{
		Name:        config.Name,
		Description: config.Description,
		ChatModel:   cm,
		SubAgents:   adkSubAgents,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tool.GetGlobalTools(), // 工具可以后续扩展
			},
		},
		MaxIteration: config.MaxIteration,
	}

	deepAgent, err := deep.New(ctx, deepConfig)
	if err != nil {
		return nil, err
	}

	return &DeepAdkAdapter{
		config:    config,
		agent:     deepAgent,
		subAgents: subAgents,
	}, nil
}

// Stream 流式输出
func (a *DeepAdkAdapter) Stream(ctx context.Context, req ai.Request) (<-chan ai.Response, error) {
	// 转换消息格式
	messages := a.convertToAdkMessages(append(req.History, req.Message))

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})

	iter := runner.Run(ctx, messages)
	out := make(chan ai.Response)

	go func() {
		defer close(out)
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			prints.EventHandler(event, func(msg *ai.Message, err error) {
				if err != nil {
					out <- ai.Response{Error: err}
					return
				}
				out <- ai.Response{Message: msg}
			})
		}
	}()

	return out, nil
}

// Name 返回 Agent 名称
func (a *DeepAdkAdapter) Name() string {
	return a.config.Name
}

// Close 释放资源
func (a *DeepAdkAdapter) Close() error {
	// 关闭子 Agent
	for _, sub := range a.subAgents {
		_ = sub.Close()
	}
	return nil
}

// GetInternalAgent 获取内部 adk.Agent（供其他适配器使用）
func (a *DeepAdkAdapter) GetInternalAgent() adk.Agent {
	return a.agent
}

// convertToAdkMessages 将自定义消息转换为 adk 消息
func (a *DeepAdkAdapter) convertToAdkMessages(msgs []*ai.Message) []adk.Message {
	var result []adk.Message
	for _, msg := range msgs {
		result = append(result, &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		})
	}
	return result
}
