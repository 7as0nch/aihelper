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
 * @Date: 2026-01-19
 * @Description: AIApplication 服务层
**/

// ==================== Application ====================

// ListApplications 获取 Application 列表
func (s *AIService) ListApplications(ctx context.Context, req *pb.ListApplicationsRequest) (*pb.ListApplicationsReply, error) {
	// 注意：applicationUC.List 目前只支持 name 和 status，type 和 scope 参数暂时忽略
	applications, total, err := s.applicationUC.List(ctx, req.Page, req.PageSize, req.Name, int(req.Status))
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ApplicationInfo, 0, len(applications))
	for _, app := range applications {
		list = append(list, toApplicationInfo(app))
	}

	return &pb.ListApplicationsReply{
		List:  list,
		Total: total,
	}, nil
}

// GetApplication 获取单个 Application
func (s *AIService) GetApplication(ctx context.Context, req *pb.IdRequest) (*pb.ApplicationInfo, error) {
	application, err := s.applicationUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toApplicationInfo(application), nil
}

// CreateApplication 创建 Application
func (s *AIService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.ApplicationInfo, error) {
	application := &model.AIApplication{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Version:     req.Version,
		Mode:        model.ProgramMode(req.Mode),
		Status:      models.Status(req.Status),
		Type:        model.ProgramType(req.Type),
		Scope:       model.Scope(req.Scope),
		Schema:      req.Schema, // Schema 字段暂时为空，后续可通过其他接口更新
	}

	// 转换 SelfAgent
	if req.SelfAgent != nil {
		application.SelfAgent = toModelAIAgent(req.SelfAgent)
	}

	application.New()

	if err := s.applicationUC.Create(ctx, application); err != nil {
		return nil, err
	}

	return toApplicationInfo(application), nil
}

// UpdateApplication 更新 Application
func (s *AIService) UpdateApplication(ctx context.Context, req *pb.UpdateApplicationRequest) (*emptypb.Empty, error) {
	application := &model.AIApplication{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Version:     req.Version,
		Mode:        model.ProgramMode(req.Mode),
		Status:      models.Status(req.Status),
		Type:        model.ProgramType(req.Type),
		Scope:       model.Scope(req.Scope),
		Schema:      req.Schema,
	}
	application.ID = req.Id

	// 转换 SelfAgent
	if req.SelfAgent != nil {
		application.SelfAgent = toModelAIAgent(req.SelfAgent)
	}

	if err := s.applicationUC.Update(ctx, application); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteApplication 删除 Application
func (s *AIService) DeleteApplication(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	if err := s.applicationUC.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ==================== 转换函数 ====================

// toApplicationInfo 将 model.AIApplication 转换为 pb.ApplicationInfo
func toApplicationInfo(app *model.AIApplication) *pb.ApplicationInfo {
	info := &pb.ApplicationInfo{
		Id:          app.ID,
		Name:        app.Name,
		Code:        app.Code,
		Description: app.Description,
		Version:     app.Version,
		Mode:        int32(app.Mode),
		Status:      int32(app.Status),
		Type:        int32(app.Type),
		Scope:       int32(app.Scope),
		Schema:      app.Schema,
		CreatedAt:   app.CreatedAt.Unix(),
	}

	// 转换 SelfAgent
	if app.SelfAgent != nil {
		info.SelfAgent = toAppAgent(app.SelfAgent)
	}

	return info
}

// toAppAgent 将 model.AIAgent 转换为 pb.AppAgent
func toAppAgent(agent *model.AIAgent) *pb.AppAgent {
	appAgent := &pb.AppAgent{
		Id:                 agent.ID,
		Name:               agent.Name,
		Code:               agent.Code,
		Description:        agent.Description,
		AdapterType:        int32(agent.AdapterType),
		OriginalModelId:    agent.AIModelID,
		MaxIteration:       int32(agent.MaxIteration),
		SystemPrompt:       agent.SystemPrompt,
		UserInputPrompt:    agent.UserInputPrompt,
		Status:             int32(agent.Status),
		Type:               int32(agent.Type),
		WithWriteTodos:     agent.WithWriteTODOs,
		WithWebSearchAgent: agent.WithWebSearchAgent,
		SystemType:         int32(agent.SystemType),
		CreatedAt:          agent.CreatedAt.Unix(),
	}

	// 转换 AIModel
	if agent.AIModel != nil {
		appAgent.AiModel = &pb.AgentModel{
			Category:    int32(agent.AIModel.Category),
			ModelType:   string(agent.AIModel.ModelType),
			ModelName:   agent.AIModel.ModelName,
			ApiKey:      agent.AIModel.APIKey,
			BaseUrl:     agent.AIModel.BaseURL,
			Temperature: agent.AIModel.Temperature,
			TopP:        agent.AIModel.TopP,
		}
	}

	// 转换子 Agent 列表
	if len(agent.SubAIAgents) > 0 {
		subAgents := make([]*pb.AppAgent, 0, len(agent.SubAIAgents))
		for _, subAgent := range agent.SubAIAgents {
			subAgents = append(subAgents, toAppAgent(subAgent))
		}
		appAgent.SubAgents = subAgents
	}

	return appAgent
}

// toModelAIAgent 将 pb.AppAgent 转换为 model.AIAgent
func toModelAIAgent(appAgent *pb.AppAgent) *model.AIAgent {
	agent := &model.AIAgent{
		Name:               appAgent.Name,
		Code:               appAgent.Code,
		Description:        appAgent.Description,
		AdapterType:        model.AdapterType(appAgent.AdapterType),
		AIModelID:          appAgent.OriginalModelId,
		MaxIteration:       int(appAgent.MaxIteration),
		SystemPrompt:       appAgent.SystemPrompt,
		UserInputPrompt:    appAgent.UserInputPrompt,
		Status:             models.Status(appAgent.Status),
		Type:               model.AgentType(appAgent.Type),
		WithWriteTODOs:     appAgent.WithWriteTodos,
		WithWebSearchAgent: appAgent.WithWebSearchAgent,
		SystemType:         models.SystemType(appAgent.SystemType),
	}

	// 如果有 ID，设置 ID（用于更新场景）
	if appAgent.Id > 0 {
		agent.ID = appAgent.Id
	}

	// 转换 AIModel
	if appAgent.AiModel != nil {
		agent.AIModel = &model.AIAgentModel{
			Category:    model.AIModel_Category(appAgent.AiModel.Category),
			ModelType:   model.ModelType(appAgent.AiModel.ModelType),
			ModelName:   appAgent.AiModel.ModelName,
			APIKey:      appAgent.AiModel.ApiKey,
			BaseURL:     appAgent.AiModel.BaseUrl,
			Temperature: appAgent.AiModel.Temperature,
			TopP:        appAgent.AiModel.TopP,
		}
	}

	// 转换子 Agent 列表
	if len(appAgent.SubAgents) > 0 {
		subAgents := make([]*model.AIAgent, 0, len(appAgent.SubAgents))
		for _, subAgent := range appAgent.SubAgents {
			subAgents = append(subAgents, toModelAIAgent(subAgent))
		}
		agent.SubAIAgents = subAgents
	}

	return agent
}
