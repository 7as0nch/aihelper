package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	pb "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/internal/biz"
	"github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/models/generator/model"
	pkgai "github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/auth"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// sessionStream 封装了单个 AI 生成任务的状态
type sessionStream struct {
	mu        sync.Mutex
	history   []*pb.Message                 // 内存缓冲区，用于支持刷新后的“断点续传”
	clients   map[chan *pb.Message]struct{} // 订阅该任务的所有客户端连接
	done      chan struct{}                 // 信号：AI 生成完成
	aiMsgID   int64                         // 数据库中预留的 AI 消息 ID
	startTime time.Time                     // 任务启动时间，用于耗时统计
}

// broadcast 向所有已连接的客户端广播消息片段
func (ss *sessionStream) broadcast(msg *pb.Message) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.history = append(ss.history, msg)
	for client := range ss.clients {
		select {
		case client <- msg:
		default: // 避免个别慢连接阻塞全局生成速度
		}
	}
}

type ChatService struct {
	pb.UnimplementedChatServer
	uc            *biz.ChatUsecase
	agent         *ai.AIUsecase
	log           *zap.Logger
	activeStreams sync.Map // Key: SessionID (int64) -> Value: *sessionStream
}

func NewChatService(uc *biz.ChatUsecase, agent *ai.AIUsecase, logger *zap.Logger) *ChatService {
	return &ChatService{
		uc:    uc,
		agent: agent,
		log:   logger,
	}
}

// HistoryById 获取会话消息详情，并智能标记流式状态
func (s *ChatService) HistoryById(ctx context.Context, req *pb.HistoryRequest) (*pb.MessagesReply, error) {
	messages, err := s.uc.GetSessionMessages(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	pbMessages := s.convertToPbMessages(messages)

	// 【优化】通过 SessionID 检查内存中是否有正在进行的任务，并标记最后一条 AI 消息
	if len(pbMessages) > 0 {
		lastMsg := pbMessages[len(pbMessages)-1]
		if lastMsg.Role == "assistant" {
			if _, running := s.activeStreams.Load(req.SessionId); running {
				lastMsg.IsStreaming = true
			}
		}
	}

	return &pb.MessagesReply{
		Messages: pbMessages,
		Total:    int32(len(messages)),
	}, nil
}

// SSEHandler 流式对话处理核心入口
func (s *ChatService) SSEHandler(w http.ResponseWriter, r *http.Request) {
	s.setupSSEHeaders(w)
	flusher := w.(http.Flusher)

	req, err := s.parseStreamRequest(r)
	if err != nil {
		s.log.Error("Parse request failed", zap.Error(err))
		return
	}

	sessionID := req.CurSessionID
	// 使用 SessionID 作为唯一 Key，彻底解决刷新导致的重复消息问题
	streamObj, loaded := s.activeStreams.LoadOrStore(sessionID, &sessionStream{
		clients:   make(map[chan *pb.Message]struct{}),
		done:      make(chan struct{}),
		startTime: time.Now(),
	})
	ss := streamObj.(*sessionStream)

	clientChan := make(chan *pb.Message, 100)

	// 如果是新任务（非刷新重连），执行初始化逻辑
	if !loaded {
		s.initAndStartStream(r.Context(), sessionID, req, ss)
	}

	// 1. 自动补发已生成的历史数据（解决“断点续传”）
	s.replayStreamHistory(w, flusher, ss)

	// 2. 注册当前连接到广播列表
	ss.mu.Lock()
	ss.clients[clientChan] = struct{}{}
	ss.mu.Unlock()

	defer func() {
		ss.mu.Lock()
		delete(ss.clients, clientChan)
		close(clientChan)
		ss.mu.Unlock()
	}()

	// 3. 进入监听循环，保持连接直到任务结束
	s.listenAndServe(w, flusher, r.Context(), ss, clientChan)
}

// initAndStartStream 初始化数据库占位符并开启 AI 协程
func (s *ChatService) initAndStartStream(ctx context.Context, sessionID int64, req *pb.SendStreamRequest, ss *sessionStream) {
	// 【即时入库】立即保存用户消息和 AI 占位消息，确保刷新后可见
	userMsg := s.uc.NewMessage(sessionID, model.RoleUser, req.CurMessage.Content)
	userMsg.QuoteId = req.CurMessage.QuoteId
	userMsg.QuoteContent = req.CurMessage.QuoteContent

	aiMsg := &model.AIChatMessage{
		Role:      model.RoleAssistant,
		SessionID: sessionID,
		Content:   "",
	}
	if req.CurMessage.AiModel != nil {
		aiMsg.AIModel = &model.UseAIModel{
			ID:        req.CurMessage.AiModel.Id,
			ModelName: req.CurMessage.AiModel.ModelName,
		}
	}

	userMsg.New()
	aiMsg.New()

	// 使用传入的 ctx (包含用户信息)，确保 beforeCreate 钩子能获取到 userID
	if err := s.uc.BatchSaveMessages(ctx, []*model.AIChatMessage{userMsg, aiMsg}); err != nil {
		s.log.Error("Initial DB save failed", zap.Error(err))
	}

	ss.aiMsgID = aiMsg.ID

	// 透传 Auth 状态启动后台协程
	// 注意：不能直接使用 ctx，因为它是请求相关的，请求结束会被 cancel
	// 我们需要创建一个携带用户信息的背景 Context
	userID := auth.GetUserId(ctx)
	bgCtx := context.WithValue(context.Background(), auth.UserId, userID)
	go s.runAIStream(bgCtx, sessionID, req, ss)
}

// runAIStream AI 生成流水线协程
func (s *ChatService) runAIStream(ctx context.Context, sessionID int64, req *pb.SendStreamRequest, ss *sessionStream) {
	defer func() {
		s.activeStreams.Delete(sessionID)
		close(ss.done)
	}()

	aiStream, err := s.agent.Stream(ctx, s.buildAIRequest(req))
	if err != nil {
		s.broadcastError(ss, err)
		return
	}

	var contentBuilder, reasoningBuilder strings.Builder
	for msg := range aiStream {
		if msg.Error != nil {
			s.log.Error("AI Stream Error", zap.Error(msg.Error))
			s.broadcastError(ss, msg.Error)
			break
		}
		if msg.Message == nil {
			continue
		}

		pbMsg := s.pkgaiToPb(msg.Message)
		ss.broadcast(pbMsg)

		contentBuilder.WriteString(msg.Message.Content)
		reasoningBuilder.WriteString(msg.Message.ReasoningContent)
	}

	// 任务结束，更新数据库中的占位符内容
	// 使用 GetMessage 获取完整记录，避免覆盖其他字段，并确保上下文透传以触发 beforeUpdate 钩子
	aiMsg, err := s.uc.GetMessage(ctx, ss.aiMsgID)
	if err != nil {
		s.log.Error("Get message for update failed", zap.Error(err), zap.Int64("id", ss.aiMsgID))
		return
	}

	aiMsg.Content = contentBuilder.String()
	aiMsg.ReasoningContent = reasoningBuilder.String()
	aiMsg.GenerateTime = time.Since(ss.startTime).String()

	if err := s.uc.UpdateMessage(ctx, aiMsg); err != nil {
		s.log.Error("Final DB update failed", zap.Error(err))
	}
}

// --- 辅助工具方法 ---

func (s *ChatService) setupSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

func (s *ChatService) parseStreamRequest(r *http.Request) (*pb.SendStreamRequest, error) {
	var req pb.SendStreamRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &req, err
	}
	err = protojson.Unmarshal(body, &req)
	return &req, err
}

func (s *ChatService) replayStreamHistory(w http.ResponseWriter, flusher http.Flusher, ss *sessionStream) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	for _, msg := range ss.history {
		jsonMsg, _ := protojson.Marshal(msg)
		fmt.Fprintf(w, "event: delta\ndata: %s\n\n", jsonMsg)
	}
	flusher.Flush()
}

func (s *ChatService) listenAndServe(w http.ResponseWriter, flusher http.Flusher, ctx context.Context, ss *sessionStream, clientChan chan *pb.Message) {
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done(): // 客户端断开连接
			return
		case <-heartbeat.C: // 发送心跳保持 Ingress 长连接
			w.Write([]byte(": keepalive\n\n"))
			flusher.Flush()
		case msg, ok := <-clientChan: // 收到新 AI 内容
			if !ok {
				return
			}
			jsonMsg, _ := protojson.Marshal(msg)
			fmt.Fprintf(w, "event: delta\ndata: %s\n\n", jsonMsg)
			flusher.Flush()
		case <-ss.done: // 后台协程生成完毕
			// 排空最后的消息
			for len(clientChan) > 0 {
				msg := <-clientChan
				jsonMsg, _ := protojson.Marshal(msg)
				fmt.Fprintf(w, "event: delta\ndata: %s\n\n", jsonMsg)
			}
			fmt.Fprintf(w, "event: done\ndata: {}\n\n")
			flusher.Flush()
			return
		}
	}
}

func (s *ChatService) broadcastError(ss *sessionStream, err error) {
	s.log.Error("AI Stream Error", zap.Error(err))
	ss.broadcast(&pb.Message{
		Role:    "assistant",
		Content: fmt.Sprintf("\n[生成出错]: %v", err),
	})
}

// --- 数据模型转换 ---

func (s *ChatService) convertToPbMessages(msgs []*model.AIChatMessage) []*pb.Message {
	var res []*pb.Message
	for _, m := range msgs {
		pbMsg := &pb.Message{
			Id:               m.ID,
			Role:             string(m.Role),
			Content:          m.Content,
			ReasoningContent: m.ReasoningContent,
			Timestamp:        m.CreatedAt.Unix(),
			QuoteId:          m.QuoteId,
			QuoteContent:     m.QuoteContent,
		}
		if m.AIModel != nil {
			pbMsg.AiModel = &pb.AIModel{
				Id:        m.AIModel.ID,
				ModelName: m.AIModel.ModelName,
			}
		}
		res = append(res, pbMsg)
	}
	return res
}

func (s *ChatService) pkgaiToPb(m *pkgai.Message) *pb.Message {
	pbMsg := &pb.Message{
		Role:             string(m.Role),
		Content:          m.Content,
		ReasoningContent: m.ReasoningContent,
	}
	if m.TokenUsage != nil {
		pbMsg.TokenUsage = &pb.TokenUsage{
			CurrentTokens: m.TokenUsage.CurrentTokens,
			TotalTokens:   m.TokenUsage.TotalTokens,
		}
	}
	return pbMsg
}

func (s *ChatService) buildAIRequest(req *pb.SendStreamRequest) pkgai.Request {
	var history []*pkgai.Message
	for _, m := range req.History {
		history = append(history, &pkgai.Message{Role: pkgai.RoleType(m.Role), Content: m.Content})
	}
	return pkgai.Request{
		History: history,
		Message: &pkgai.Message{
			Role:         pkgai.RoleType(req.CurMessage.Role),
			Content:      req.CurMessage.Content,
			QuoteId:      req.CurMessage.QuoteId,
			QuoteContent: req.CurMessage.QuoteContent,
		},
	}
}

// --- gRPC 接口实现占位 (保持兼容性) ---
func (s *ChatService) History(ctx context.Context, req *pb.HistoryRequest) (*pb.HistoryReply, error) {
	userId := auth.GetUserId(ctx)
	sessions, total, err := s.uc.ListSessions(ctx, userId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var pbSessions []*pb.HistorySession
	for _, session := range sessions {
		pbSessions = append(pbSessions, &pb.HistorySession{
			Id: session.ID, Title: session.Title, UpdateTime: session.UpdatedAt.Unix(),
		})
	}
	return &pb.HistoryReply{Sessions: pbSessions, Total: int32(total)}, nil
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
	return &pb.HistorySession{Id: session.ID, Title: session.Title, UpdateTime: session.UpdatedAt.Unix()}, nil
}
func (s *ChatService) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.Message, error) {
	msg := s.uc.NewMessage(req.SessionId, model.RoleType(req.Role), req.Content)
	saved, err := s.uc.SaveMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &pb.Message{Id: saved.ID, Role: string(saved.Role), Content: saved.Content, Timestamp: saved.CreatedAt.Unix()}, nil
}
func (s *ChatService) SendStream(req *pb.SendStreamRequest, conn pb.Chat_SendStreamServer) error {
	return nil
}
