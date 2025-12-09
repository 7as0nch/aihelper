package prints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/pkg/ai"
	// "github.com/cloudwego/eino-examples/internal/logs"
)

func EventHandler(event *adk.AgentEvent, handlerFn func(msg *ai.Message, err error)) {
	log.Printf("name: %s\npath: %s", event.AgentName, event.RunPath)
	if event.Output != nil && event.Output.MessageOutput != nil {
		if m := event.Output.MessageOutput.Message; m != nil {
			if len(m.Content) > 0 {
				if m.Role == schema.Tool {
					log.Printf("\ntool response: %s", m.Content)
				} else {
					log.Printf("\nanswer: %s", m.Content)
				}
			}
			if len(m.ToolCalls) > 0 {
				for _, tc := range m.ToolCalls {
					log.Printf("\ntool name: %s", tc.Function.Name)
					log.Printf("\narguments: %s", tc.Function.Arguments)
				}
			}
		} else if s := event.Output.MessageOutput.MessageStream; s != nil {
			toolMap := map[int][]*schema.Message{}
			var contentStart, thinkingStart bool
			charNumOfOneRow := 0
			maxCharNumOfOneRow := 120
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
				if chunk.ReasoningContent != "" {
					if !thinkingStart {
						thinkingStart = true
						if chunk.Role == schema.Tool {
							log.Printf("\ntool response: ")
						} else {
							log.Printf("\nThinking: ")
						}
					}

					charNumOfOneRow += len(chunk.ReasoningContent)
					if strings.Contains(chunk.ReasoningContent, "\n") {
						charNumOfOneRow = 0
					} else if charNumOfOneRow >= maxCharNumOfOneRow {
						charNumOfOneRow = 0
					}
					var tk *ai.TokenUsage
					if chunk.ResponseMeta != nil && chunk.ResponseMeta.Usage != nil {
						tk = &ai.TokenUsage{
							CurrentTokens: int64(chunk.ResponseMeta.Usage.CompletionTokens),
							TotalTokens:   int64(chunk.ResponseMeta.Usage.TotalTokens),
						}
					}
					handlerFn(&ai.Message{
						Role:             ai.RoleType(chunk.Role),
						ReasoningContent: chunk.ReasoningContent,
						TokenUsage:       tk,
					}, nil)
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
					charNumOfOneRow += len(chunk.Content)
					if strings.Contains(chunk.Content, "\n") {
						charNumOfOneRow = 0
					} else if charNumOfOneRow >= maxCharNumOfOneRow {
						charNumOfOneRow = 0
					}
					var tk *ai.TokenUsage
					if chunk.ResponseMeta != nil && chunk.ResponseMeta.Usage != nil {
						tk = &ai.TokenUsage{
							CurrentTokens: int64(chunk.ResponseMeta.Usage.CompletionTokens),
							TotalTokens:   int64(chunk.ResponseMeta.Usage.TotalTokens),
						}
					}
					handlerFn(&ai.Message{
						Role:    ai.RoleType(chunk.Role),
						Content: chunk.Content,
						TokenUsage: tk,
					}, nil)
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
					}
				}
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
