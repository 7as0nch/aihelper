package prints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/pkg/ai"
	// "github.com/cloudwego/eino-examples/internal/logs"
)

func EventHandler(event *adk.AgentEvent, handlerFn func(msg *ai.Message, err error)) {
	log.Printf("name: %s\npath: %s", event.AgentName, event.RunPath)
	if event.Output != nil && event.Output.MessageOutput != nil {
		if m := event.Output.MessageOutput.Message; m != nil {
			msg := &ai.Message{
				Role:    ai.RoleType(schema.Assistant),
			}
			if len(m.Content) > 0 {
				if m.Role == schema.Tool {
					log.Printf("\ntool response: %s", m.Content)
				} else {
					log.Printf("\nanswer: %s", m.Content)
					msg.Content = m.Content
				}
			}
			if len(m.ToolCalls) > 0 {
				for _, tc := range m.ToolCalls {
					log.Printf("\ntool name: %s", tc.Function.Name)
					log.Printf("\narguments: %s", tc.Function.Arguments)
					msg.CallingTools = append(msg.CallingTools, &ai.CallingTool{
						Name:         tc.Function.Name,
						FunctionName: tc.Function.Name,
					})
				}
			}
			handlerFn(msg, nil)
		} else if s := event.Output.MessageOutput.MessageStream; s != nil {
			toolMap := map[int][]*schema.Message{}
			var contentStart, thinkingStart bool
			for {
				chunk, err := s.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("error: %v", err)
					handlerFn(nil, err)
					return
				}
				msg := &ai.Message{
					Role:             ai.RoleType(chunk.Role),
				}
				if chunk.Role != schema.Tool {
					msg.ReasoningContent = chunk.ReasoningContent
					msg.Content = chunk.Content
				}
				if chunk.ReasoningContent != "" {
					if !thinkingStart {
						thinkingStart = true
						if chunk.Role == schema.Tool {
							log.Printf("\ntool response: ")
						} else {
							log.Printf("\nThinking: ")
						}
					}
				}
				if chunk.Content != "" {
					if !contentStart {
						contentStart = true
						if chunk.Role == schema.Tool {
							log.Printf("\ntool response: ")
						} else {
							log.Printf("\nanswer: ")
						}
					}
				}

				if len(chunk.ToolCalls) > 0 {
					for _, tc := range chunk.ToolCalls {
						index := tc.Index
						if index == nil {
							log.Fatalf("index is nil")
						}
						toolMap[*index] = append(toolMap[*index], &schema.Message{
							Role: chunk.Role,
							ToolCalls: []schema.ToolCall{
								{
									ID:    tc.ID,
									Type:  tc.Type,
									Index: tc.Index,
									Function: schema.FunctionCall{
										Name:      tc.Function.Name,
										Arguments: tc.Function.Arguments,
									},
								},
							},
						})
						if tc.Function.Name != "" {
							tl := &ai.CallingTool{
								Name:         tc.Function.Name,
								FunctionName: tc.Function.Name,
							}
							msg.CallingTools = append(msg.CallingTools, tl)
							if tc.Function.Arguments != "" {
								var args map[string]interface{}
								if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
									log.Fatalf("Unmarshal arguments failed: %v", err)
									return
								}
								tl.Args = args
							}
						}
					}
				}
				if chunk.ResponseMeta.FinishReason == "stop" {
					if chunk.ResponseMeta.Usage != nil {
						log.Printf("\nusage: %v", chunk.ResponseMeta.Usage)
						msg.TokenUsage = &ai.TokenUsage{
							CurrentTokens: int64(chunk.ResponseMeta.Usage.PromptTokens),
							TotalTokens:   int64(chunk.ResponseMeta.Usage.TotalTokens),
						}
					}
				}
				handlerFn(msg, nil)
			}

			for _, msgs := range toolMap {
				m, err := schema.ConcatMessages(msgs)
				if err != nil {
					log.Fatalf("ConcatMessage failed: %v", err)
					return
				}
				log.Printf("\ntool name: %s", m.ToolCalls[0].Function.Name)
				log.Printf("\narguments: %s", m.ToolCalls[0].Function.Arguments)
			}
		}
	}
	if event.Action != nil {
		if event.Action.TransferToAgent != nil {
			fmt.Printf("\naction: transfer to %v", event.Action.TransferToAgent.DestAgentName)
		}
		if event.Action.Interrupted != nil {
			ii, _ := json.MarshalIndent(event.Action.Interrupted.Data, "  ", "  ")
			log.Println("action: interrupted")
			log.Printf("interrupt snapshot: %v\n", string(ii))
		}
		// if event.Action.Exit {
		// 	fmt.Printf("\naction: exit")
		// }
	}
	if event.Err != nil {
		handlerFn(nil, event.Err)
	}
}
