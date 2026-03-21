package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	pb "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/internal/biz"
	"github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/internal/consts"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models/generator/model"
	pkgai "github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/auth"
	redislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// sessionStream 封装了单个 AI 生成任务的状态
type sessionStream struct {
	mu        sync.Mutex
	clients   map[chan *pb.Message]struct{} // 订阅该任务的所有客户端连接
	done      chan struct{}                 // 信号：AI 生成完成
	aiMsgID   int64                         // 数据库中预留的 AI 消息 ID
	startTime time.Time                     // 任务启动时间，用于耗时统计
	sessionID int64
	redis     db.RedisRepo
}

// broadcast 向所有已连接的客户端广播消息片段
func (ss *sessionStream) broadcast(msg *pb.Message) {
	if ss.redis != nil {
		if payload, err := protojson.Marshal(msg); err == nil {
			_ = ss.redis.RPush(context.Background(), fmt.Sprintf(consts.RedisKeyChatStreamDelta, ss.sessionID), 30*time.Minute, string(payload))
		}
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
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
	redis         db.RedisRepo
	activeStreams sync.Map // Key: SessionID (int64) -> Value: *sessionStream
}

func NewChatService(uc *biz.ChatUsecase, agent *ai.AIUsecase, logger *zap.Logger, redisRepo db.RedisRepo) *ChatService {
	return &ChatService{
		uc:    uc,
		agent: agent,
		log:   logger,
		redis: redisRepo,
	}
}

type streamStatus string

const (
	streamStatusRunning streamStatus = "running"
	streamStatusDone    streamStatus = "done"
	streamStatusError   streamStatus = "error"
	streamMetaTTL                    = 30 * time.Minute
	streamRunningTTL                 = 2 * time.Minute
)

type streamMeta struct {
	SessionID int64        `json:"session_id"`
	AIMessageID int64      `json:"ai_message_id"`
	StartUnix int64        `json:"start_unix"`
	UpdatedUnix int64      `json:"updated_unix"`
	Status    streamStatus `json:"status"`
	Error     string       `json:"error,omitempty"`
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
			if _, running := s.activeStreams.Load(req.SessionId); running || s.isStreamRunning(ctx, req.SessionId) {
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
		sessionID: sessionID,
		redis:     s.redis,
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
	s.resetStreamCache(ctx, sessionID)

	// 【防御性优化】如果是续传请求（ID > 0），尝试关联已有消息，防止重复创建
	if req.CurMessage.Id > 0 {
		messages, err := s.uc.GetSessionMessages(ctx, sessionID)
		if err == nil {
			for i, m := range messages {
				// 找到用户发出的那条消息，其下一条通常就是对应的 AI 回复占位符
				if m.ID == req.CurMessage.Id && i+1 < len(messages) {
					nextMsg := messages[i+1]
					if nextMsg.Role == model.RoleAssistant {
						ss.aiMsgID = nextMsg.ID
						s.saveStreamMeta(ctx, ss, streamStatusRunning, "")
						s.startAIStream(ctx, sessionID, req, ss)
						return
					}
				}
			}
		}
	}

	// 以下是正常新消息入库流程
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

	if err := s.uc.BatchSaveMessages(ctx, []*model.AIChatMessage{userMsg, aiMsg}); err != nil {
		s.log.Error("Initial DB save failed", zap.Error(err))
	}

	ss.aiMsgID = aiMsg.ID
	s.saveStreamMeta(ctx, ss, streamStatusRunning, "")
	s.startAIStream(ctx, sessionID, req, ss)
}

// 提取启动协程的逻辑
func (s *ChatService) startAIStream(ctx context.Context, sessionID int64, req *pb.SendStreamRequest, ss *sessionStream) {
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
	var latestCallingTools []*pkgai.CallingTool
	var latestQuoteSearchLinks []*pkgai.QuoteSearchLink
	toolIndex := make(map[string]*pkgai.CallingTool)
	toolOrder := make([]string, 0, 4)
	linkIndex := make(map[string]*pkgai.QuoteSearchLink)
	linkOrder := make([]string, 0, 6)
	finalStatus := streamStatusDone
	finalErrText := ""
	for msg := range aiStream {
		if msg.Error != nil {
			s.log.Error("AI Stream Error", zap.Error(msg.Error))
			s.broadcastError(ss, msg.Error)
			finalStatus = streamStatusError
			finalErrText = msg.Error.Error()
			break
		}
		if msg.Message == nil {
			continue
		}

		pbMsg := s.pkgaiToPb(msg.Message)
		ss.broadcast(pbMsg)
		s.saveStreamMeta(ctx, ss, streamStatusRunning, "")

		contentBuilder.WriteString(msg.Message.Content)
		reasoningBuilder.WriteString(msg.Message.ReasoningContent)
		if len(msg.Message.CallingTools) > 0 {
			for _, tool := range msg.Message.CallingTools {
				key := callingToolKey(tool)
				if key == "" {
					continue
				}
				if _, exists := toolIndex[key]; exists {
					continue
				}
				toolIndex[key] = tool
				toolOrder = append(toolOrder, key)
			}
			latestCallingTools = snapshotCallingTools(toolIndex, toolOrder)
		}
		if len(msg.Message.QuoteSearchLinks) > 0 {
			for _, link := range msg.Message.QuoteSearchLinks {
				key := quoteSearchLinkKey(link)
				if key == "" {
					continue
				}
				if _, exists := linkIndex[key]; exists {
					continue
				}
				linkIndex[key] = link
				linkOrder = append(linkOrder, key)
			}
			latestQuoteSearchLinks = snapshotQuoteSearchLinks(linkIndex, linkOrder)
		}
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
	if len(latestCallingTools) > 0 {
		aiMsg.CallingTools = convertCallingToolsPkgToModel(latestCallingTools)
	}
	if len(latestQuoteSearchLinks) > 0 {
		aiMsg.QuoteSearchLinks = convertQuoteSearchLinksPkgToModel(latestQuoteSearchLinks)
	}

	if err := s.uc.UpdateMessage(ctx, aiMsg); err != nil {
		s.log.Error("Final DB update failed", zap.Error(err))
	}
	s.saveStreamMeta(ctx, ss, finalStatus, finalErrText)
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

func (s *ChatService) resetStreamCache(ctx context.Context, sessionID int64) {
	if s.redis == nil {
		return
	}
	_ = s.redis.Del(ctx, fmt.Sprintf(consts.RedisKeyChatStreamMeta, sessionID))
	_ = s.redis.Del(ctx, fmt.Sprintf(consts.RedisKeyChatStreamDelta, sessionID))
}

func (s *ChatService) saveStreamMeta(ctx context.Context, ss *sessionStream, status streamStatus, errText string) {
	if s.redis == nil || ss == nil {
		return
	}

	payload, err := json.Marshal(streamMeta{
		SessionID:  ss.sessionID,
		AIMessageID: ss.aiMsgID,
		StartUnix:  ss.startTime.Unix(),
		UpdatedUnix: time.Now().Unix(),
		Status:     status,
		Error:      errText,
	})
	if err != nil {
		s.log.Warn("marshal stream meta failed", zap.Error(err), zap.Int64("session_id", ss.sessionID))
		return
	}
	if err = s.redis.Set(ctx, fmt.Sprintf(consts.RedisKeyChatStreamMeta, ss.sessionID), string(payload), streamMetaTTL); err != nil {
		s.log.Warn("save stream meta failed", zap.Error(err), zap.Int64("session_id", ss.sessionID))
	}
}

func (s *ChatService) loadStreamMeta(ctx context.Context, sessionID int64) (*streamMeta, error) {
	if s.redis == nil {
		return nil, nil
	}

	value, err := s.redis.Get(ctx, fmt.Sprintf(consts.RedisKeyChatStreamMeta, sessionID))
	if err != nil {
		if errors.Is(err, redislib.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var meta streamMeta
	if err = json.Unmarshal([]byte(value), &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func (s *ChatService) isStreamRunning(ctx context.Context, sessionID int64) bool {
	meta, err := s.loadStreamMeta(ctx, sessionID)
	if err != nil || meta == nil {
		return false
	}
	if meta.Status != streamStatusRunning {
		return false
	}
	return time.Since(time.Unix(meta.UpdatedUnix, 0)) <= streamRunningTTL
}

func (s *ChatService) replayStreamHistory(w http.ResponseWriter, flusher http.Flusher, ss *sessionStream) {
	if s.redis == nil {
		flusher.Flush()
		return
	}

	values, err := s.redis.LRange(context.Background(), fmt.Sprintf(consts.RedisKeyChatStreamDelta, ss.sessionID), 0, -1)
	if err != nil && !errors.Is(err, redislib.Nil) {
		s.log.Warn("replay stream history failed", zap.Error(err), zap.Int64("session_id", ss.sessionID))
		flusher.Flush()
		return
	}
	for _, value := range values {
		fmt.Fprintf(w, "event: delta\ndata: %s\n\n", value)
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
	s.saveStreamMeta(context.Background(), ss, streamStatusError, err.Error())
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
			QuoteSearchLinks: convertQuoteSearchLinksModelToPb(m.QuoteSearchLinks),
			CallingTools:     convertCallingToolsModelToPb(m.CallingTools),
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
		QuoteSearchLinks: convertQuoteSearchLinksPkgToPb(m.QuoteSearchLinks),
		CallingTools:     convertCallingToolsPkgToPb(m.CallingTools),
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
	// 是否需要生成待办计划。 'smart' | 'need' | 'no'
	needTODOPlan := pkgai.NeedTODOPlanSmart
	switch req.NeedTODOPlan {
	case "need":
		needTODOPlan = pkgai.NeedTODOPlanNeed
	case "no":
		needTODOPlan = pkgai.NeedTODOPlanUnNeed
	}
	return pkgai.Request{
		History: history,
		Message: &pkgai.Message{
			Role:         pkgai.RoleType(req.CurMessage.Role),
			Content:      req.CurMessage.Content,
			QuoteId:      req.CurMessage.QuoteId,
			QuoteContent: req.CurMessage.QuoteContent,
			AIModel: &pkgai.AIModel{
				ID:           req.CurMessage.AiModel.Id,
				ModelName:    req.CurMessage.AiModel.ModelName,
				ThinkingMode: req.CurMessage.AiModel.ThinkingMode,
				SearchByWeb:  req.CurMessage.AiModel.SearchByWeb,
			},
		},
		NeedTODOPlan: needTODOPlan,
	}
}

func convertQuoteSearchLinksPkgToModel(links []*pkgai.QuoteSearchLink) model.QuoteSearchLinks {
	if len(links) == 0 {
		return nil
	}
	res := make(model.QuoteSearchLinks, 0, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}
		res = append(res, &model.QuoteSearchLink{
			Url:       link.Url,
			Title:     link.Title,
			Content:   link.Content,
			Highlight: link.Highlight,
		})
	}
	return res
}

func convertCallingToolsPkgToModel(tools []*pkgai.CallingTool) model.CallingTools {
	if len(tools) == 0 {
		return nil
	}
	res := make(model.CallingTools, 0, len(tools))
	for _, tool := range tools {
		if tool == nil {
			continue
		}
		res = append(res, &model.CallingTool{
			Name:         tool.Name,
			Description:  tool.Description,
			FunctionName: tool.FunctionName,
		})
	}
	return res
}

func convertQuoteSearchLinksPkgToPb(links []*pkgai.QuoteSearchLink) []*pb.QuoteSearchLink {
	if len(links) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(links))
	res := make([]*pb.QuoteSearchLink, 0, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}
		key := quoteSearchLinkKey(link)
		if key != "" && seen[key] {
			continue
		}
		if key != "" {
			seen[key] = true
		}
		res = append(res, &pb.QuoteSearchLink{
			Url:       link.Url,
			Title:     link.Title,
			Content:   link.Content,
			Highlight: link.Highlight,
		})
	}
	return res
}

func convertCallingToolsPkgToPb(tools []*pkgai.CallingTool) []*pb.CallingTool {
	if len(tools) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(tools))
	res := make([]*pb.CallingTool, 0, len(tools))
	for _, tool := range tools {
		if tool == nil {
			continue
		}
		key := callingToolKey(tool)
		if key != "" && seen[key] {
			continue
		}
		if key != "" {
			seen[key] = true
		}
		res = append(res, &pb.CallingTool{
			Name:         tool.Name,
			Description:  tool.Description,
			FunctionName: tool.FunctionName,
		})
	}
	return res
}

func convertQuoteSearchLinksModelToPb(links model.QuoteSearchLinks) []*pb.QuoteSearchLink {
	if len(links) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(links))
	res := make([]*pb.QuoteSearchLink, 0, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}
		key := quoteSearchLinkKey(&pkgai.QuoteSearchLink{
			Url:       link.Url,
			Title:     link.Title,
			Content:   link.Content,
			Highlight: link.Highlight,
		})
		if key != "" && seen[key] {
			continue
		}
		if key != "" {
			seen[key] = true
		}
		res = append(res, &pb.QuoteSearchLink{
			Url:       link.Url,
			Title:     link.Title,
			Content:   link.Content,
			Highlight: link.Highlight,
		})
	}
	return res
}

func convertCallingToolsModelToPb(tools model.CallingTools) []*pb.CallingTool {
	if len(tools) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(tools))
	res := make([]*pb.CallingTool, 0, len(tools))
	for _, tool := range tools {
		if tool == nil {
			continue
		}
		key := callingToolKey(&pkgai.CallingTool{
			Name:         tool.Name,
			Description:  tool.Description,
			FunctionName: tool.FunctionName,
		})
		if key != "" && seen[key] {
			continue
		}
		if key != "" {
			seen[key] = true
		}
		res = append(res, &pb.CallingTool{
			Name:         tool.Name,
			Description:  tool.Description,
			FunctionName: tool.FunctionName,
		})
	}
	return res
}

func callingToolKey(tool *pkgai.CallingTool) string {
	if tool == nil {
		return ""
	}
	if tool.FunctionName != "" {
		return "fn:" + tool.FunctionName
	}
	if tool.Name != "" {
		return "name:" + tool.Name
	}
	if tool.Description != "" {
		return "desc:" + tool.Description
	}
	return ""
}

func quoteSearchLinkKey(link *pkgai.QuoteSearchLink) string {
	if link == nil {
		return ""
	}
	if link.Url != "" {
		return "url:" + link.Url
	}
	if link.Title != "" || link.Content != "" {
		return "title:" + link.Title + "|content:" + link.Content
	}
	return ""
}

func snapshotCallingTools(index map[string]*pkgai.CallingTool, order []string) []*pkgai.CallingTool {
	if len(order) == 0 {
		return nil
	}
	res := make([]*pkgai.CallingTool, 0, len(order))
	for _, key := range order {
		if tool, ok := index[key]; ok {
			res = append(res, tool)
		}
	}
	return res
}

func snapshotQuoteSearchLinks(index map[string]*pkgai.QuoteSearchLink, order []string) []*pkgai.QuoteSearchLink {
	if len(order) == 0 {
		return nil
	}
	res := make([]*pkgai.QuoteSearchLink, 0, len(order))
	for _, key := range order {
		if link, ok := index[key]; ok {
			res = append(res, link)
		}
	}
	return res
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
