package ai

import (
	"context"

	pb "github.com/example/aichat/backend/api/ai"
	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/models/generator/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-12-11 15:44:19
 * @Description:
**/

// ==================== Agent ====================

// ListAgents 获取 Agent 列表
func (s *AIService) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsReply, error) {
	agents, total, err := s.agentUC.List(ctx, req.Page, req.PageSize, req.Name, int(req.Status), int(req.Type))
	if err != nil {
		return nil, err
	}

	list := make([]*pb.AgentInfo, 0, len(agents))
	for _, a := range agents {
		// 获取子 Agent IDs
		list = append(list, toAgentInfo(a, nil))
	}

	return &pb.ListAgentsReply{
		List:  list,
		Total: total,
	}, nil
}

// GetAgent 获取单个 Agent
func (s *AIService) GetAgent(ctx context.Context, req *pb.IdRequest) (*pb.AgentInfo, error) {
	agent, err := s.agentUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toAgentInfo(agent, nil), nil
}

// CreateAgent 创建 Agent
func (s *AIService) CreateAgent(ctx context.Context, req *pb.CreateAgentRequest) (*pb.AgentInfo, error) {
	agent := &model.AIAgent{
		Name:               req.Name,
		Code:               req.Code,
		Description:        req.Description,
		AdapterType:        model.AdapterType(req.AdapterType),
		AIModelID:          req.OriginalModelId,
		MaxIteration:       int(req.MaxIteration),
		SystemPrompt:       req.SystemPrompt,
		UserInputPrompt:    req.UserInputPrompt,
		Status:             models.Status(req.Status),
		Type:               model.AgentType(req.Type),
		WithWriteTODOs:     req.WithWriteTodos,
		WithWebSearchAgent: req.WithWebSearchAgent,
	}

	// 处理 ai_model 配置
	if req.AiModel != nil {
		agent.AIModel = toModelAIAgentModel(req.AiModel)
	}

	agent.New()

	if err := s.agentUC.Create(ctx, agent); err != nil {
		return nil, err
	}

	return toAgentInfo(agent, nil), nil
}

// UpdateAgent 更新 Agent
func (s *AIService) UpdateAgent(ctx context.Context, req *pb.UpdateAgentRequest) (*emptypb.Empty, error) {
	agent := &model.AIAgent{
		Name:               req.Name,
		Code:               req.Code,
		Description:        req.Description,
		AdapterType:        model.AdapterType(req.AdapterType),
		AIModelID:          req.OriginalModelId,
		MaxIteration:       int(req.MaxIteration),
		SystemPrompt:       req.SystemPrompt,
		UserInputPrompt:    req.UserInputPrompt,
		Status:             models.Status(req.Status),
		Type:               model.AgentType(req.Type),
		WithWriteTODOs:     req.WithWriteTodos,
		WithWebSearchAgent: req.WithWebSearchAgent,
	}
	agent.ID = req.Id

	// 处理 ai_model 配置
	if req.AiModel != nil {
		agent.AIModel = toModelAIAgentModel(req.AiModel)
	}

	if err := s.agentUC.Update(ctx, agent); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteAgent 删除 Agent
func (s *AIService) DeleteAgent(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	if err := s.agentUC.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// toPbAgentModel 将 model.AIAgentModel 转换为 pb.AgentModel
func toPbAgentModel(m *model.AIAgentModel) *pb.AgentModel {
	if m == nil {
		return nil
	}
	return &pb.AgentModel{
		Category:    int32(m.Category),
		ModelType:   string(m.ModelType),
		ModelName:   m.ModelName,
		ApiKey:      m.APIKey,
		BaseUrl:     m.BaseURL,
		Temperature: m.Temperature,
		TopP:        m.TopP,
	}
}

// ==================== 转换函数 ====================

func toAgentInfo(a *model.AIAgent, subIDs []int64) *pb.AgentInfo {
	info := &pb.AgentInfo{
		Id:                 a.ID,
		Name:               a.Name,
		Code:               a.Code,
		Description:        a.Description,
		AdapterType:        int32(a.AdapterType),
		OriginalModelId:    a.AIModelID,
		MaxIteration:       int32(a.MaxIteration),
		SystemPrompt:       a.SystemPrompt,
		UserInputPrompt:    a.UserInputPrompt,
		Status:             int32(a.Status),
		Type:               int32(a.Type),
		WithWriteTodos:     a.WithWriteTODOs,
		WithWebSearchAgent: a.WithWebSearchAgent,
		SystemType:         int32(a.SystemType),
		CreatedAt:          a.CreatedAt.Unix(),
		// ParentAgentId:      0,
	}

	// 转换 ai_model 配置
	if a.AIModel != nil {
		info.AiModel = toPbAgentModel(a.AIModel)
	}

	return info
}

// toModelAIAgentModel 将 pb.AgentModel 转换为 model.AIAgentModel
func toModelAIAgentModel(pbModel *pb.AgentModel) *model.AIAgentModel {
	if pbModel == nil {
		return nil
	}
	return &model.AIAgentModel{
		Category:    model.AIModel_Category(pbModel.Category),
		ModelType:   model.ModelType(pbModel.ModelType),
		ModelName:   pbModel.ModelName,
		APIKey:      pbModel.ApiKey,
		BaseURL:     pbModel.BaseUrl,
		Temperature: pbModel.Temperature,
		TopP:        pbModel.TopP,
	}
}