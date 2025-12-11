/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIPromptTemplate 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

// AIPromptRepo AIPromptTemplate 数据仓库接口
type AIPromptRepo interface {
	List(ctx context.Context, page, pageSize int32, name string, promptType int) ([]*model.AIPromptTemplate, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AIPromptTemplate, error)
	Create(ctx context.Context, p *model.AIPromptTemplate) error
	Update(ctx context.Context, p *model.AIPromptTemplate) error
	Delete(ctx context.Context, id int64) error
}

// AIPromptUseCase AIPromptTemplate 业务逻辑
type AIPromptUseCase struct {
	repo AIPromptRepo
}

// NewAIPromptUseCase 创建 AIPrompt UseCase
func NewAIPromptUseCase(repo AIPromptRepo) *AIPromptUseCase {
	return &AIPromptUseCase{
		repo: repo,
	}
}

// List 获取提示词模板列表
func (uc *AIPromptUseCase) List(ctx context.Context, page, pageSize int32, name string, promptType int) ([]*model.AIPromptTemplate, int64, error) {
	return uc.repo.List(ctx, page, pageSize, name, promptType)
}

// GetByID 根据 ID 获取提示词模板
func (uc *AIPromptUseCase) GetByID(ctx context.Context, id int64) (*model.AIPromptTemplate, error) {
	return uc.repo.GetByID(ctx, id)
}

// Create 创建提示词模板
func (uc *AIPromptUseCase) Create(ctx context.Context, p *model.AIPromptTemplate) error {
	return uc.repo.Create(ctx, p)
}

// Update 更新提示词模板
func (uc *AIPromptUseCase) Update(ctx context.Context, p *model.AIPromptTemplate) error {
	return uc.repo.Update(ctx, p)
}

// Delete 删除提示词模板
func (uc *AIPromptUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
