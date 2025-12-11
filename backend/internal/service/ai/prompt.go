/* *
 * @Author: chengjiang
 * @Date: 2025-12-11 15:44:53
 * @Description:
**/
package ai

import (
	"context"

	pb "github.com/example/aichat/backend/api/ai"
	"github.com/example/aichat/backend/models/generator/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ==================== Prompt ====================

// ListPrompts 获取提示词模板列表
func (s *AIService) ListPrompts(ctx context.Context, req *pb.ListPromptsRequest) (*pb.ListPromptsReply, error) {
	prompts, total, err := s.promptUC.List(ctx, req.Page, req.PageSize, req.Name, int(req.Type))
	if err != nil {
		return nil, err
	}

	list := make([]*pb.PromptInfo, 0, len(prompts))
	for _, p := range prompts {
		list = append(list, toPromptInfo(p))
	}

	return &pb.ListPromptsReply{
		List:  list,
		Total: total,
	}, nil
}

// GetPrompt 获取单个提示词模板
func (s *AIService) GetPrompt(ctx context.Context, req *pb.IdRequest) (*pb.PromptInfo, error) {
	p, err := s.promptUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toPromptInfo(p), nil
}

// CreatePrompt 创建提示词模板
func (s *AIService) CreatePrompt(ctx context.Context, req *pb.CreatePromptRequest) (*pb.PromptInfo, error) {
	p := &model.AIPromptTemplate{
		Type:        model.AI_Prompt_Type(req.Type),
		Name:        req.Name,
		Description: req.Description,
		Text:        req.Text,
	}
	p.New()

	if err := s.promptUC.Create(ctx, p); err != nil {
		return nil, err
	}
	return toPromptInfo(p), nil
}

// UpdatePrompt 更新提示词模板
func (s *AIService) UpdatePrompt(ctx context.Context, req *pb.UpdatePromptRequest) (*emptypb.Empty, error) {
	p := &model.AIPromptTemplate{
		Type:        model.AI_Prompt_Type(req.Type),
		Name:        req.Name,
		Description: req.Description,
		Text:        req.Text,
	}
	p.ID = req.Id

	if err := s.promptUC.Update(ctx, p); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeletePrompt 删除提示词模板
func (s *AIService) DeletePrompt(ctx context.Context, req *pb.IdRequest) (*emptypb.Empty, error) {
	if err := s.promptUC.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}