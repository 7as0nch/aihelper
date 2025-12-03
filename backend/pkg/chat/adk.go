/* *
 * @Author: chengjiang
 * @Date: 2025-11-03 14:08:14
 * @Description:
**/
package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"

	// "github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	v1 "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/pkg/agenttools"
	"github.com/example/aichat/backend/pkg/prints"
	"github.com/go-kratos/kratos/v2/log"
	arkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

type AdkAgent struct {
	agent adk.Agent
}

func newAdkModel(ctx context.Context) model.ToolCallingChatModel {
	chatmodel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-reasoner",
		APIKey:      "sk-dc49fec5d27a416f8758ece703aed2ff",
		Temperature: 0.7,
	})
	if err != nil {
		return nil
	}
	return chatmodel
}

func newAdkReActModel(ctx context.Context) model.ToolCallingChatModel {
	chatmodel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-reasoner",
		APIKey:  "sk-dc49fec5d27a416f8758ece703aed2ff",
		Thinking: &arkmodel.Thinking{
			Type: arkmodel.ThinkingTypeEnabled,
		},
	})
	if err != nil {
		return nil
	}
	return chatmodel
}

func NewAdkAgent() *AdkAgent {
	ctx := context.Background()

	globalAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "globalAgent",
		Description: "你是一个全局工具调用助手:可获取系统配置，如当前系统时间。", // English description
		Instruction: `
			1. 调用时间工具不需要传任何参数，直接获取即可。
		`,
		Model: newAdkModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agenttools.GetGlobalTools(),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewChatModelAgent failed, err=%v", err)
		return nil
	}
	bussinessAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "bussinessAgent",
		Description: "你是一个业务工具调用助手:可查询用户报表，渠道ROI等。", // English description
		Instruction: `可支持查询ROI，用户报表等功能。`,
		Model:       newAdkModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agenttools.GetTools(),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewChatModelAgent failed, err=%v", err)
		return nil
	}
	// helloAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
	// 	Name:        "helloAgent",
	// 	Description: "你是一个招呼智能体，打招呼的时候使用", // English description
	// 	Instruction: `客客气气打招呼，可以调用打招呼的工具`,
	// 	Model:       newAdkReActModel(ctx),
	// 	ToolsConfig: adk.ToolsConfig{
	// 		ToolsNodeConfig: compose.ToolsNodeConfig{
	// 			Tools: agenttools.McpTools(),
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	// logs.Errorf("NewChatModelAgent failed, err=%v", err)
	// 	return nil
	// }
	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "HostAgent",
		Description: "你是一个主管智能助手，可以调用其他智能体处理用户的需求", // English description
		Instruction: `如果需要获取当前时间用globalAgent，业务（用户报表，ROI等）相关调用bussinessAgent进行相互协作，最后总结由host主agent返回给用户。`,
		Model:       newAdkModel(ctx),
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewChatModelAgent failed, err=%v", err)
		return nil
	}
	res, err := adk.SetSubAgents(ctx, a, []adk.Agent{globalAgent, bussinessAgent})
	// res, err := supervisor.New(ctx, &supervisor.Config{
	// 	Supervisor: a,
	// 	SubAgents: []adk.Agent{globalAgent, bussinessAgent},
	// })
	if err != nil {
		return nil
	}
	return &AdkAgent{agent: res}
}

func (a *AdkAgent) Run(ctx context.Context, messages []adk.Message, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})
	return runner.Run(ctx, messages, opts...)
}

func (a *AdkAgent) Stream(ctx context.Context, req *v1.SendStreamRequest) (<-chan *v1.Message, error) {
	var messages []*schema.Message
	messages = append(messages, PgHelperPrompt())
	for _, msg := range req.History {
		if msg.Role == "system" {
			continue
		}
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		})
	}

	if req.CurMessage != nil {
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(req.CurMessage.Role),
			Content: req.CurMessage.Content,
		})
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})

	adkMessages := make([]adk.Message, len(messages))
	copy(adkMessages, messages)

	iter := runner.Run(ctx, adkMessages)

	out := make(chan *v1.Message)

	go func() {
		defer close(out)
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				if errors.Is(event.Err, context.Canceled) || strings.Contains(event.Err.Error(), "context canceled") {
					log.Infof("Stream canceled: %v", event.Err)
				} else {
					log.Errorf("Stream error: %v", event.Err)
				}
				break
			}

			// Convert event to pb.Message
			if event.Output != nil && event.Output.MessageOutput != nil {
				// msgOut := event.Output.MessageOutput
				// role := string(msgOut.Message.Role)

				prints.EventHandler(event, func(thinking, content string, err error) {
					if err == io.EOF {
						return
					}

					pbMsg := &v1.Message{
						// Role: role, // Role might not be available in every chunk or might be redundant
						Content:          content,
						Role:             "assisant",
						ReasoningContent: thinking,
					}
					out <- pbMsg
				})
			}
		}
	}()

	return out, nil
}
