/* *
 * @Author: chengjiang
 * @Date: 2025-12-11 15:44:40
 * @Description:
**/
package ai

import (
	"context"

	pb "github.com/example/aichat/backend/api/ai"
	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/models/generator/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ==================== Tool ====================

// ListTools 获取工具列表
func (s *AIService) ListTools(ctx context.Context, req *pb.ListToolsRequest) (*pb.ListToolsReply, error) {
	tools, total, err := s.toolUC.List(ctx, req.Page, req.PageSize, req.Name, int(req.Type), int(req.Status))
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ToolInfo, 0, len(tools))
	for _, t := range tools {
		list = append(list, toToolInfo(t))
	}

	return &pb.ListToolsReply{
		List:  list,
		Total: total,
	}, nil
}

// GetTool 获取单个工具
func (s *AIService) GetTool(ctx context.Context, req *pb.IdRequest) (*pb.ToolInfo, error) {
	t, err := s.toolUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toToolInfo(t), nil
}

// CreateTool 创建工具
func (s *AIService) CreateTool(ctx context.Context, req *pb.CreateToolRequest) (*pb.ToolInfo, error) {
	t := &model.AITool{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Type:        model.AI_Tool_Type(req.Type),
		Status:      models.Status(req.Status),
		Params:      toToolParams(req.Params),
		MCPUrl:      req.McpUrl,
		MCPToken:    req.McpToken,
	}
	t.New()

	if err := s.toolUC.Create(ctx, t); err != nil {
		return nil, err
	}
	return toToolInfo(t), nil
}

// UpdateTool 更新工具
func (s *AIService) UpdateTool(ctx context.Context, req *pb.UpdateToolRequest) (*emptypb.Empty, error) {
	t := &model.AITool{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Type:        model.AI_Tool_Type(req.Type),
		Status:      models.Status(req.Status),
		Params:      toToolParams(req.Params),
		MCPUrl:      req.McpUrl,
		MCPToken:    req.McpToken,
	}
	t.ID = req.Id

	if err := s.toolUC.Update(ctx, t); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteTool 删除工具
func (s *AIService) DeleteTool(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	if err := s.toolUC.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// BatchBindTools 批量绑定工具到 Agent
func (s *AIService) BatchBindTools(ctx context.Context, req *pb.BatchBindToolsRequest) (*emptypb.Empty, error) {
	if err := s.toolUC.BatchBindTools(ctx, req.AgentId, req.ToolCodes); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ==================== 转换函数 ====================

func toAgentInfo(a *model.AIAgent, subIDs []int64) *pb.AgentInfo {
	return &pb.AgentInfo{
		Id:                 a.ID,
		Name:               a.Name,
		Code:               a.Code,
		Description:        a.Description,
		AdapterType:        int32(a.AdapterType),
		AiModelId:          a.AIModelID,
		MaxIteration:       int32(a.MaxIteration),
		SystemPrompt:       a.SystemPrompt,
		UserInputPrompt:    a.UserInputPrompt,
		Status:             int32(a.Status),
		Type:               int32(a.Type),
		Order:              int32(a.Order),
		WithWriteTodos:     a.WithWriteTODOs,
		WithWebSearchAgent: a.WithWebSearchAgent,
		SystemType:         int32(a.SystemType),
		CreatedAt:          a.CreatedAt.Unix(),
		SubAgentIds:        subIDs,
	}
}

func toModelInfo(m *model.AIModel) *pb.ModelInfo {
	return &pb.ModelInfo{
		Id:          m.ID,
		Category:    int32(m.Category),
		ModelType:   string(m.ModelType),
		ModelName:   m.ModelName,
		ApiKey:      m.APIKey,
		BaseUrl:     m.BaseURL,
		MaxTokens:   int32(m.MaxTokens),
		Temperature: m.Temperature,
		TopP:        m.TopP,
		PriceType:   int32(m.PriceType),
		Price:       m.Price,
		Supplier:    m.Supplier,
		Description: m.Description,
		Status:      int32(m.Status),
		IsDefault:   int32(m.IsDefault),
		CreatedAt:   m.CreatedAt.Unix(),
	}
}

func toPromptInfo(p *model.AIPromptTemplate) *pb.PromptInfo {
	return &pb.PromptInfo{
		Id:          p.ID,
		Type:        int32(p.Type),
		Name:        p.Name,
		Description: p.Description,
		Text:        p.Text,
		CreatedAt:   p.CreatedAt.Unix(),
	}
}

func toToolInfo(t *model.AITool) *pb.ToolInfo {
	params := make([]*pb.ToolParam, 0, len(t.Params))
	for _, p := range t.Params {
		params = append(params, &pb.ToolParam{
			ParamName:    p.ParamName,
			ParamType:    p.ParamType,
			DefaultValue: p.Default,
		})
	}
	return &pb.ToolInfo{
		Id:          t.ID,
		Name:        t.Name,
		Code:        t.Code,
		Description: t.Description,
		SysType:     int32(t.SysType),
		Type:        int32(t.Type),
		Status:      int32(t.Status),
		Params:      params,
		McpUrl:      t.MCPUrl,
		McpToken:    t.MCPToken,
		CreatedAt:   t.CreatedAt.Unix(),
	}
}

func toToolParams(params []*pb.ToolParam) model.AI_Tool_Params {
	result := make(model.AI_Tool_Params, 0, len(params))
	for _, p := range params {
		result = append(result, &model.AI_Tool_Param{
			ParamName: p.ParamName,
			ParamType: p.ParamType,
			Default:   p.DefaultValue,
		})
	}
	return result
}
