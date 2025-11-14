package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/pkg/chat"
	"go.uber.org/zap"
)

type ChatService struct {
	v1.UnimplementedChatServer
	log *zap.Logger
	// 可以添加消息存储或其他依赖
}

func NewChatService(logger *zap.Logger) *ChatService {
	return &ChatService{
		log: logger,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, req *v1.SendMessageRequest) (*v1.SendMessageReply, error) {
	// s.log.WithContext(ctx).Infof("Sending message from user %s in session %s", req.UserId, req.SessionId)

	// 这里应该实现实际的消息发送逻辑
	// 例如：保存到数据库，推送到消息队列等
	messageId := fmt.Sprintf("msg_%d", time.Now().Unix())

	// 模拟处理时间
	time.Sleep(10 * time.Millisecond)

	s.log.Info("Message sent successfully with ID", zap.String("messageId", messageId))

	return &v1.SendMessageReply{
		MessageId: messageId,
		Success:   true,
	}, nil
}

func (s *ChatService) GetMessages(ctx context.Context, req *v1.GetMessagesRequest) (*v1.GetMessagesReply, error) {
	s.log.Info("Getting messages for session", zap.String("sessionId", req.SessionId), zap.Int32("page", req.Page))

	// 这里应该实现从数据库获取消息的逻辑
	// 模拟返回一些消息
	messages := make([]*v1.ChatMessage, 0)

	// 模拟分页数据
	for i := 0; i < 10; i++ {
		messages = append(messages, &v1.ChatMessage{
			Id:        fmt.Sprintf("msg_%d_%d", req.Page, i),
			Content:   fmt.Sprintf("This is message %d on page %d", i, req.Page),
			UserId:    "user_123",
			SessionId: req.SessionId,
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
			Type:      "user",
		})
	}

	return &v1.GetMessagesReply{
		Messages: messages,
		Total:    100, // 模拟总数量
	}, nil
}

func (s *ChatService) StreamMessages(req *v1.StreamMessagesRequest, stream v1.Chat_StreamMessagesServer) error {
	s.log.Info("Starting message stream for session", zap.String("sessionId", req.SessionId))

	// 这里应该实现实际的流式消息推送逻辑
	// 例如：从消息队列订阅消息并推送给客户端
	// 模拟发送几条消息
	for i := 0; i < 5; i++ {
		message := &v1.StreamMessagesReply{
			Message: &v1.ChatMessage{
				Id:        fmt.Sprintf("stream_msg_%d", i),
				Content:   fmt.Sprintf("Streamed message %d", i),
				UserId:    "system",
				SessionId: req.SessionId,
				Timestamp: time.Now().Format(time.RFC3339),
				Type:      "assistant",
			},
		}

		if err := stream.Send(message); err != nil {
			s.log.Error("Failed to send streamed message", zap.Error(err))
			return err
		}

		// 模拟延迟
		time.Sleep(1 * time.Second)
	}

	s.log.Info("Message stream completed for session", zap.String("sessionId", req.SessionId))
	return nil
}

// StreamMessagesToWebSocket streams messages to a WebSocket connection
func (s *ChatService) StreamMessagesToWebSocket(sessionId string, sendMessage func(message map[string]interface{}), sendComplete func()) {
	s.log.Info("Starting message stream for session %s via WebSocket", zap.String("sessionId", sessionId))

	// 模拟发送几条消息
	for i := 0; i < 5; i++ {
		message := map[string]interface{}{
			"id":         fmt.Sprintf("stream_msg_%d", i),
			"content":    fmt.Sprintf("Streamed message %d", i),
			"user_id":    "system",
			"session_id": sessionId,
			"timestamp":  time.Now().Format(time.RFC3339),
			"type":       "assistant",
		}

		sendMessage(message)

		// 模拟延迟
		time.Sleep(1 * time.Second)
	}

	sendComplete()
}

func (s *ChatService) SSEHandler(w http.ResponseWriter, r *http.Request) {
	chat.NewAdkAgent().SSEHandler(w, r)
}

func (s *ChatService) SSEHandlerOld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	// 拿到post body
	body := r.Body
	defer body.Close()
	// 读取post body
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		s.log.Error("Read body error:", zap.Error(err))
		return
	}
	var req v1.SendMessageRequest
	if err = json.Unmarshal(bodyBytes, &req); err != nil {
		s.log.Error("Unmarshal body error:", zap.Error(err))
		return
	}
	chatAgent := chat.NewDeepseekAgent()
	if chatAgent == nil {
		s.log.Info("NewAiAgent failed")
		return
	}
	var messages []*schema.Message
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			continue
		}
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		})
	}
	agentStream, err := chatAgent.GetReActAgent().Stream(context.Background(),
		messages, agent.WithComposeOptions(compose.WithCallbacks(chat.GetCallback())))
	if err != nil {
		s.log.Info("Stream failed:", zap.Error(err))
		return
	}
	defer agentStream.Close()
	for {
		msg, err := agentStream.Recv()
		if err != nil {
			if err == io.EOF {
				s.log.Info("Stream completed")
				event := "data: [DONE]\n\n"
				if _, err = w.Write([]byte(event)); err != nil {
					s.log.Error("Write error:", zap.Error(err))
					return
				}
				flusher.Flush()
				break
			}
			s.log.Info("Recv failed:", zap.Error(err))
			return
		}
		var content string
		if msg.ReasoningContent != "" {
			content = fmt.Sprintf("</think>%v</think>", msg.ReasoningContent)
		} else {
			content = fmt.Sprintf("%v", msg.Content)
		}
		// 处理换行符，将其转换为\\n以便在SSE中正确传输
		content = strings.ReplaceAll(content, "\n", "\\n")
		event := "data: " + content + "\n\n"

		if _, err := w.Write([]byte(event)); err != nil {
			s.log.Error("Write error:", zap.Error(err))
			return
		}
		flusher.Flush()
	}
}