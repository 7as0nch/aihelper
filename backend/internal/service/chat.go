package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	pb "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/pkg/chat"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatService struct {
	pb.UnimplementedChatServer
	log *zap.Logger
	// 可以添加消息存储或其他依赖
}

func NewChatService(logger *zap.Logger) *ChatService {
	return &ChatService{
		log: logger,
	}
}

func (s *ChatService) History(ctx context.Context, req *pb.HistoryRequest) (*pb.HistoryReply, error) {
	return &pb.HistoryReply{}, nil
}
func (s *ChatService) HistoryById(ctx context.Context, req *pb.HistoryRequest) (*pb.MessagesReply, error) {
	return &pb.MessagesReply{}, nil
}
func (s *ChatService) HistoryRenameById(ctx context.Context, req *pb.HistoryRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *ChatService) HistoryDeleteById(ctx context.Context, req *pb.HistoryRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *ChatService) SendStream(req *pb.SendStreamRequest, conn pb.Chat_SendStreamServer) error {
	for {
		err := conn.Send(&pb.SendStreamReply{})
		if err != nil {
			return err
		}
	}
}

func (s *ChatService) SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.log.Error("Read body error:", zap.Error(err))
		return
	}
	defer r.Body.Close()

	var req pb.SendStreamRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.log.Error("Unmarshal body error:", zap.Error(err))
		return
	}

	// Use the abstract Agent interface
	var agent chat.Agent[*pb.SendStreamRequest, *pb.Message] = chat.NewAdkAgent()
	// if agent == nil {
	// 	s.log.Error("Failed to create agent")
	// 	return
	// }

	stream, err := agent.Stream(r.Context(), &req)
	if err != nil {
		s.log.Error("Stream failed:", zap.Error(err))
		return
	}

	for msg := range stream {
		// Escape newlines for SSE data
		// msg.Content = strings.ReplaceAll(msg.Content, "\n", "\\n")
		// msg.ReasoningContent = strings.ReplaceAll(msg.ReasoningContent, "\n", "\\n")
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			s.log.Error("Marshal message error:", zap.Error(err))
			return
		}
		event := fmt.Sprintf("data: %v\n\n", string(jsonMsg))
		if _, err := w.Write([]byte(event)); err != nil {
			s.log.Error("Write error:", zap.Error(err))
			return
		}
		flusher.Flush()
	}

	// Send DONE event
	if _, err := w.Write([]byte("data: [DONE]\n\n")); err != nil {
		s.log.Error("Write error:", zap.Error(err))
	}
	flusher.Flush()
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
	var req pb.SendStreamRequest
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
	for _, msg := range req.History {
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
