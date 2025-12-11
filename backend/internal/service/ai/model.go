/* *
 * @Author: chengjiang
 * @Date: 2025-12-11 15:44:29
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

// ==================== Model ====================

// ListModels 获取模型列表
func (s *AIService) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsReply, error) {
	models, total, err := s.modelUC.List(ctx, req.Page, req.PageSize, req.ModelName, int(req.Status))
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ModelInfo, 0, len(models))
	for _, m := range models {
		list = append(list, toModelInfo(m))
	}

	return &pb.ListModelsReply{
		List:  list,
		Total: total,
	}, nil
}

// GetModel 获取单个模型
func (s *AIService) GetModel(ctx context.Context, req *pb.IdRequest) (*pb.ModelInfo, error) {
	m, err := s.modelUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toModelInfo(m), nil
}

// CreateModel 创建模型
func (s *AIService) CreateModel(ctx context.Context, req *pb.CreateModelRequest) (*pb.ModelInfo, error) {
	m := &model.AIModel{
		Category:    model.AIModel_Category(req.Category),
		ModelType:   model.ModelType(req.ModelType),
		ModelName:   req.ModelName,
		APIKey:      req.ApiKey,
		BaseURL:     req.BaseUrl,
		MaxTokens:   int(req.MaxTokens),
		Temperature: req.Temperature,
		TopP:        req.TopP,
		PriceType:   model.AIModel_PriceType(req.PriceType),
		Price:       req.Price,
		Supplier:    req.Supplier,
		Description: req.Description,
		Status:      models.Status(req.Status),
		IsDefault:   models.Status(req.IsDefault),
	}
	m.New()

	if err := s.modelUC.Create(ctx, m); err != nil {
		return nil, err
	}
	return toModelInfo(m), nil
}

// UpdateModel 更新模型
func (s *AIService) UpdateModel(ctx context.Context, req *pb.UpdateModelRequest) (*emptypb.Empty, error) {
	m := &model.AIModel{
		Category:    model.AIModel_Category(req.Category),
		ModelType:   model.ModelType(req.ModelType),
		ModelName:   req.ModelName,
		APIKey:      req.ApiKey,
		BaseURL:     req.BaseUrl,
		MaxTokens:   int(req.MaxTokens),
		Temperature: req.Temperature,
		TopP:        req.TopP,
		PriceType:   model.AIModel_PriceType(req.PriceType),
		Price:       req.Price,
		Supplier:    req.Supplier,
		Description: req.Description,
		Status:      models.Status(req.Status),
		IsDefault:   models.Status(req.IsDefault),
	}
	m.ID = req.Id

	if err := s.modelUC.Update(ctx, m); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteModel 删除模型
func (s *AIService) DeleteModel(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	if err := s.modelUC.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
