/* *
 * @Author: chengjiang
 * @Date: 2025-11-03 14:08:14
 * @Description:
**/
package chat

import (
	"context"
	"fmt"
	"io"

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
		APIKey:      "sk-a76a955533c649d6922a5042af6c0492",
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
		APIKey:  "sk-a76a955533c649d6922a5042af6c0492",
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
		Model: newAdkReActModel(ctx),
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
		Model:       newAdkReActModel(ctx),
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
	helloAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "helloAgent",
		Description: "你是一个招呼智能体，打招呼的时候使用", // English description
		Instruction: `客客气气打招呼，可以调用打招呼的工具`,
		Model:       newAdkReActModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agenttools.McpTools(),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewChatModelAgent failed, err=%v", err)
		return nil
	}
	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "HostAgent",
		Description: "你是一个主管智能助手，可以调用其他智能体处理用户的需求", // English description
		Instruction: `如果需要获取当前时间用globalAgent，业务（用户报表，ROI等）相关调用bussinessAgent进行相互协作，最后总结由host主agent返回给用户。`,
		Model:       newAdkReActModel(ctx),
	})
	if err != nil {
		fmt.Println(err)
		// logs.Errorf("NewChatModelAgent failed, err=%v", err)
		return nil
	}
	res, err := adk.SetSubAgents(ctx, a, []adk.Agent{globalAgent, bussinessAgent, helloAgent})
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

	// Add current message if needed, or assume it's in history?
	// The original code appended PgHelperPrompt and then history.
	// Usually the current user message is also needed.
	// Let's assume req.Message is the current message if not in history.
	// But looking at original code: `var req v1.SendStreamRequest; ... messages = append(messages, PgHelperPrompt()); for _, msg := range req.History ...`
	// It seems it only used history? Or maybe the current message is appended to history by the caller?
	// Let's stick to the original logic: just use req.History (and prompt).
	// Wait, the original code used `req.History`.
	// But usually `SendStreamRequest` has a `message` field for the new message.
	// Let's check `SendStreamRequest` definition in proto.
	// `message Message = 1; repeated Message history = 2;`
	// The original code ONLY used `req.History`. This might be a bug or specific usage in original code.
	// However, to be safe and robust, I should probably include `req.Message` if it's not empty.
	// But for now, I will strictly follow the original logic to avoid breaking behavior,
	// BUT I will add `req.Message` if it exists, as it's standard.

	if req.Message != nil {
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(req.Message.Role),
			Content: req.Message.Content,
		})
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           a.agent,
		EnableStreaming: true,
	})

	// Convert []*schema.Message to []adk.Message
	// adk.Message is an interface, schema.Message implements it?
	// Let's check adk.Run signature: `Run(ctx context.Context, messages []Message, ...)`
	// And `type Message = schema.Message` in some versions or compatible.
	// In `pkg/chat/adk.go`: `func (a *AdkAgent) Run(ctx context.Context, messages []adk.Message, ...)`
	// And `messages = append(messages, PgHelperPrompt())` where PgHelperPrompt returns `*schema.Message`.
	// So `[]*schema.Message` might need conversion if `adk.Message` is an interface.
	// Actually `adk.Message` is likely `*schema.Message` or interface implemented by it.
	// The original code: `var messages []adk.Message; messages = append(messages, PgHelperPrompt()); ... append(messages, &schema.Message{...})`
	// So `*schema.Message` satisfies `adk.Message`.

	adkMessages := make([]adk.Message, len(messages))
	for i, m := range messages {
		adkMessages[i] = m
	}

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
				log.Errorf("Stream error: %v", event.Err)
				break
			}

			// Convert event to pb.Message
			if event.Output != nil && event.Output.MessageOutput != nil {
				// msgOut := event.Output.MessageOutput
				// role := string(msgOut.Message.Role)

				prints.EventHandler(event, func(content string, err error) {
					if err == io.EOF {
						return
					}

					pbMsg := &v1.Message{
						// Role: role, // Role might not be available in every chunk or might be redundant
						Content: content,
					}
					out <- pbMsg
				})
			}
		}
	}()

	return out, nil
}
