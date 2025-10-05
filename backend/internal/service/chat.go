package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/pkg/chat"
	"github.com/go-kratos/kratos/v2/log"
)

type ChatService struct {
	v1.UnimplementedChatServer
	log *log.Helper
	// 可以添加消息存储或其他依赖
}

func NewChatService(logger log.Logger) *ChatService {
	return &ChatService{
		log: log.NewHelper(logger),
	}
}

func (s *ChatService) SendMessage(ctx context.Context, req *v1.SendMessageRequest) (*v1.SendMessageReply, error) {
	// s.log.WithContext(ctx).Infof("Sending message from user %s in session %s", req.UserId, req.SessionId)

	// 这里应该实现实际的消息发送逻辑
	// 例如：保存到数据库，推送到消息队列等
	messageId := fmt.Sprintf("msg_%d", time.Now().Unix())

	// 模拟处理时间
	time.Sleep(10 * time.Millisecond)

	s.log.WithContext(ctx).Infof("Message sent successfully with ID: %s", messageId)

	return &v1.SendMessageReply{
		MessageId: messageId,
		Success:   true,
	}, nil
}

func (s *ChatService) GetMessages(ctx context.Context, req *v1.GetMessagesRequest) (*v1.GetMessagesReply, error) {
	s.log.WithContext(ctx).Infof("Getting messages for session %s, page %d", req.SessionId, req.Page)

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
	ctx := stream.Context()
	s.log.WithContext(ctx).Infof("Starting message stream for session %s", req.SessionId)

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
			s.log.Errorf("Failed to send streamed message: %v", err)
			return err
		}

		// 模拟延迟
		time.Sleep(1 * time.Second)
	}

	s.log.Infof("Message stream completed for session %s", req.SessionId)
	return nil
}

// StreamMessagesToWebSocket streams messages to a WebSocket connection
func (s *ChatService) StreamMessagesToWebSocket(sessionId string, sendMessage func(message map[string]interface{}), sendComplete func()) {
	s.log.Infof("Starting message stream for session %s via WebSocket", sessionId)

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

// ensureMarkdownLineBreaks 确保Markdown格式元素有正确的换行
func ensureMarkdownLineBreaks(content string) string {
	// 处理标题，确保标题前后有空行
	// 处理分割线，确保前后有空行
	// 处理代码块，确保前后有空行
	
	// 在标题前添加空行（如果不在开头）
	content = addLineBreakBeforePattern(content, `(?m)^#{1,6}\s`)
	
	// 在分割线前后添加空行
	content = addLineBreakAroundPattern(content, `(?m)^[-*]{3,}\s*$`)
	
	// 在代码块前后添加空行
	content = addLineBreakAroundPattern(content, "(?m)^ *```[a-z]* *\\s*$")
	
	return content
}

// addLineBreakBeforePattern 在匹配模式前添加空行
func addLineBreakBeforePattern(content, pattern string) string {
	// 编译正则表达式
	re := regexp.MustCompile(pattern)
	
	// 在匹配模式前添加空行（如果不在开头）
	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 如果匹配在开头，不需要添加空行
		if strings.HasPrefix(content, match) {
			return match
		}
		// 否则在前面添加空行
		return "\n\n" + match
	})
}

// addLineBreakAroundPattern 在匹配模式前后添加空行
func addLineBreakAroundPattern(content, pattern string) string {
	// 编译正则表达式
	re := regexp.MustCompile(pattern)
	
	// 在匹配模式前后添加空行
	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 在前后添加空行，但避免重复的空行
		return "\n\n" + strings.TrimSpace(match) + "\n\n"
	})
}

func (s *ChatService) SSEHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Info("Read body error:", err)
		return
	}
	var req v1.SendMessageRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		log.Info("Unmarshal body error:", err)
		return
	}
	chatAgent := chat.NewAiAgent()
	if chatAgent == nil {
		log.Info("NewAiAgent failed")
		return
	}
	var messages []*schema.Message
	for _, msg := range req.Messages {
		
		messages = append(messages, &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		})
	}
	agentStream, err := chatAgent.GetChatModel().Stream(context.Background(), messages)
	if err != nil {
		log.Info("Stream failed:", err)
		return
	}
	defer agentStream.Close()
	for {
		msg, err := agentStream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Info("Stream completed")
				event := "data: [DONE]\n\n"
				if _, err = w.Write([]byte(event)); err != nil {
					log.Info("Write error:", err)
					return
				}
				flusher.Flush()
				break
			}
			log.Info("Recv failed:", err)
			return
		}
		// 处理Markdown格式，确保标题和分割线等元素有正确的换行
		content := fmt.Sprintf("%v", msg.Content)
		// 确保标题前后有空行
		content = ensureMarkdownLineBreaks(content)
		// 处理换行符，将其转换为\\n以便在SSE中正确传输
		content = strings.ReplaceAll(content, "\n", "\\n")
		event := "data: " + content + "\n\n"

		
		if _, err := w.Write([]byte(event)); err != nil {
			log.Info("Write error:", err)
			return
		}
		flusher.Flush()
	}
}
