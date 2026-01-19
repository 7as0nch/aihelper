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
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		AdapterType: model.AdapterType(req.AdapterType),
		AIModelID:   req.AiModelId,
		// AIModel: &model.AIModel{
		// 	Category:    model.AIModel_Category(req.Category),
		// 	ModelType:   model.ModelType(req.ModelType),
		// 	ModelName:   req.ModelName,
		// 	APIKey:      req.ApiKey,
		// 	BaseURL:     req.BaseUrl,
		// 	MaxTokens:   int(req.MaxTokens),
		// 	Temperature: req.Temperature,
		// 	TopP:        req.TopP,
		// 	PriceType:   model.AIModel_PriceType(req.PriceType),
		// 	Price:       req.Price,
		// 	Supplier:    req.Supplier,
		// 	Description: req.Description,
		// },
		MaxIteration:       int(req.MaxIteration),
		SystemPrompt:       req.SystemPrompt,
		UserInputPrompt:    req.UserInputPrompt,
		Status:             models.Status(req.Status),
		Type:               model.AgentType(req.Type),
		Order:              int(req.Order),
		WithWriteTODOs:     req.WithWriteTodos,
		WithWebSearchAgent: req.WithWebSearchAgent,
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
		AIModelID:          req.AiModelId,
		MaxIteration:       int(req.MaxIteration),
		SystemPrompt:       req.SystemPrompt,
		UserInputPrompt:    req.UserInputPrompt,
		Status:             models.Status(req.Status),
		Type:               model.AgentType(req.Type),
		Order:              int(req.Order),
		WithWriteTODOs:     req.WithWriteTodos,
		WithWebSearchAgent: req.WithWebSearchAgent,
	}
	agent.ID = req.Id

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

// BatchBindSubAgents 批量绑定子 Agent
func (s *AIService) BatchBindSubAgents(ctx context.Context, req *pb.BatchBindSubAgentsRequest) (*emptypb.Empty, error) {
	
	return &emptypb.Empty{}, nil
}
