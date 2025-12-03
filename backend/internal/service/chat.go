package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	pb "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/internal/biz"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/pkg/auth"
	"github.com/example/aichat/backend/pkg/chat"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatService struct {
	pb.UnimplementedChatServer
	uc  *biz.ChatUsecase
	log *zap.Logger
}

func NewChatService(uc *biz.ChatUsecase, logger *zap.Logger) *ChatService {
	return &ChatService{
		uc:  uc,
		log: logger,
	}
}

func (s *ChatService) History(ctx context.Context, req *pb.HistoryRequest) (*pb.HistoryReply, error) {
	userId := auth.GetUserId(ctx)
	sessions, total, err := s.uc.ListSessions(ctx, userId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var pbSessions []*pb.HistorySession
	for _, session := range sessions {
		pbSessions = append(pbSessions, &pb.HistorySession{
			Id:         session.ID,
			Title:      session.Title,
			UpdateTime: session.UpdatedAt.Unix(),
		})
	}
	return &pb.HistoryReply{
		Sessions: pbSessions,
		Total:    int32(total),
	}, nil
}

func (s *ChatService) HistoryById(ctx context.Context, req *pb.HistoryRequest) (*pb.MessagesReply, error) {
	messages, err := s.uc.GetSessionMessages(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	var pbMessages []*pb.Message
	for _, msg := range messages {
		pbMsg := &pb.Message{
			Id:               msg.ID,
			Role:             string(msg.Role),
			Content:          msg.Content,
			ReasoningContent: msg.ReasoningContent,
			Timestamp:        msg.CreatedAt.Unix(),
			QuoteId:          msg.QuoteId,
			QuoteContent:     msg.QuoteContent,
			IsStreaming:      msg.IsStreaming,
		}

		if msg.AIModel != nil {
			pbMsg.AiModel = &pb.AIModel{
				Id:           msg.AIModel.ID,
				ModelName:    msg.AIModel.ModelName,
				ThinkingMode: msg.AIModel.ThinkingMode,
			}
		}
		if len(msg.QuoteSearchLinks) > 0 {
			var links []*pb.QuoteSearchLink
			for _, link := range msg.QuoteSearchLinks {
				links = append(links, &pb.QuoteSearchLink{
					Url:       link.Url,
					Title:     link.Title,
					Content:   link.Content,
					Highlight: link.Highlight,
				})
			}
			pbMsg.QuoteSearchLinks = links
		}
		if msg.TokenUsage != nil {
			pbMsg.TokenUsage = &pb.TokenUsage{
				CurrentTokens: msg.TokenUsage.CurrentTokens,
				TotalTokens:   msg.TokenUsage.TotalTokens,
			}
		}
		if len(msg.CallingTools) > 0 {
			var tools []*pb.CallingTool
			for _, tool := range msg.CallingTools {
				tools = append(tools, &pb.CallingTool{
					Name:         tool.Name,
					Description:  tool.Description,
					FunctionName: tool.FunctionName,
				})
			}
			pbMsg.CallingTools = tools
		}
		if len(msg.Attachments) > 0 {
			var attachments []*pb.Attachment
			for _, att := range msg.Attachments {
				attachments = append(attachments, &pb.Attachment{
					Id:   att.ID,
					Type: att.Type,
					Name: att.Name,
					Url:  att.Url,
				})
			}
			pbMsg.Attachments = attachments
		}
		pbMessages = append(pbMessages, pbMsg)
	}
	return &pb.MessagesReply{
		Messages: pbMessages,
		Total:    int32(len(messages)),
	}, nil
}

func (s *ChatService) HistoryRenameById(ctx context.Context, req *pb.HistoryRequest) (*emptypb.Empty, error) {
	err := s.uc.RenameSession(ctx, req.SessionId, req.Name)
	return &emptypb.Empty{}, err
}

func (s *ChatService) HistoryDeleteById(ctx context.Context, req *pb.HistoryRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteSession(ctx, req.SessionId)
	return &emptypb.Empty{}, err
}

func (s *ChatService) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.HistorySession, error) {
	userId := auth.GetUserId(ctx)
	session, err := s.uc.CreateSession(ctx, userId, req.Title)
	if err != nil {
		return nil, err
	}
	for _, m := range req.Messages {
		msg := s.uc.NewMessage(session.ID, model.RoleType(m.Role), m.Content)
		msg.New()
		_, err := s.uc.SaveMessage(ctx, msg)
		if err != nil {
			return nil, err
		}
	}
	return &pb.HistorySession{
		Id:         session.ID,
		Title:      session.Title,
		UpdateTime: session.UpdatedAt.Unix(),
	}, nil
}

func (s *ChatService) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.Message, error) {
	msg := s.uc.NewMessage(req.SessionId, model.RoleType(req.Role), req.Content)
	savedMsg, err := s.uc.SaveMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &pb.Message{
		Id:               savedMsg.ID,
		Role:             string(savedMsg.Role),
		Content:          savedMsg.Content,
		ReasoningContent: savedMsg.ReasoningContent,
		Timestamp:        savedMsg.CreatedAt.Unix(),
		QuoteId:          savedMsg.QuoteId,
		QuoteContent:     savedMsg.QuoteContent,
		IsStreaming:      savedMsg.IsStreaming,
	}, nil
}

func (s *ChatService) SendStream(req *pb.SendStreamRequest, conn pb.Chat_SendStreamServer) error {
	// Not implemented for gRPC stream, using SSEHandler
	return nil
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
	if err = protojson.Unmarshal(body, &req); err != nil {
		s.log.Error("Unmarshal body error:", zap.Error(err))
		return
	}

	ctx := r.Context()
	userId := auth.GetUserId(ctx)

	var sessionID int64
	// I'll read `Session-Id` from header.
	sessionIDStr := r.Header.Get("X-Session-Id")
	if sessionIDStr != "" {
		fmt.Sscanf(sessionIDStr, "%d", &sessionID)
	}

	// If sessionID is 0, create new session
	if sessionID == 0 {
		// Create new session
		title := "New Chat"
		if len(req.CurMessage.Content) > 0 {
			runes := []rune(req.CurMessage.Content)
			if len(runes) > 20 {
				title = string(runes[:20])
			} else {
				title = string(runes)
			}
		}
		session, err := s.uc.CreateSession(ctx, userId, title)
		if err != nil {
			s.log.Error("Create session error:", zap.Error(err))
			return
		}
		sessionID = session.ID
	}

	// Save User Message
	userMsg := s.uc.NewMessage(sessionID, model.RoleUser, req.CurMessage.Content)
	if _, err := s.uc.SaveMessage(ctx, userMsg); err != nil {
		s.log.Error("Save user message error:", zap.Error(err))
		return
	}

	// Use the abstract Agent interface
	var agent chat.Agent[*pb.SendStreamRequest, *pb.Message] = chat.NewAdkAgent()

	stream, err := agent.Stream(r.Context(), &req)
	if err != nil {
		s.log.Error("Stream failed:", zap.Error(err))
		return
	}

	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	// Accumulate AI response for saving
	var aiContentBuilder strings.Builder
	var aiReasoningBuilder strings.Builder
	var aiMsg *model.AIChatMessage

	for msg := range stream {
		// Accumulate content
		aiContentBuilder.WriteString(msg.Content)
		aiReasoningBuilder.WriteString(msg.ReasoningContent)

		jsonMsg, err := marshaler.Marshal(msg)
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

	// Save AI Message
	aiMsg = s.uc.NewMessage(sessionID, model.RoleAssistant, aiContentBuilder.String())
	aiMsg.ReasoningContent = aiReasoningBuilder.String()
	// TODO: Save other fields like TokenUsage, etc. if available in the last message or accumulated
	if _, err := s.uc.SaveMessage(ctx, aiMsg); err != nil {
		s.log.Error("Save AI message error:", zap.Error(err))
	}

	// Send DONE event
	if _, err := w.Write([]byte("data: [DONE]\n\n")); err != nil {
		s.log.Error("Write error:", zap.Error(err))
	}
	flusher.Flush()
}
