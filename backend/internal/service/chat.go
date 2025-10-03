package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/aichat/backend/api/chat/v1"
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
	s.log.WithContext(ctx).Infof("Sending message from user %s in session %s", req.UserId, req.SessionId)
	
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
	s.log.Infof("Starting message stream for session %s", req.SessionId)
	
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