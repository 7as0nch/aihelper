/* *
 * @Author: chengjiang
 * @Date: 2026-01-19
 * @Description: AIWorkflow 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

// AIWorkflowRepo AIWorkflow 数据仓库接口
type AIWorkflowRepo interface {
	List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIWorkflow, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AIWorkflow, error)
	GetByCode(ctx context.Context, code string) (*model.AIWorkflow, error)
	Create(ctx context.Context, m *model.AIWorkflow) error
	Update(ctx context.Context, m *model.AIWorkflow) error
	Delete(ctx context.Context, id int64) error
}

// AIWorkflowUseCase AIWorkflow 业务逻辑
type AIWorkflowUseCase struct {
	repo AIWorkflowRepo
}

// NewAIWorkflowUseCase 创建 AIWorkflow UseCase
func NewAIWorkflowUseCase(repo AIWorkflowRepo) *AIWorkflowUseCase {
	return &AIWorkflowUseCase{
		repo: repo,
	}
}

func (uc *AIWorkflowUseCase) List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIWorkflow, int64, error) {
	return uc.repo.List(ctx, page, pageSize, name, status)
}

func (uc *AIWorkflowUseCase) GetByID(ctx context.Context, id int64) (*model.AIWorkflow, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *AIWorkflowUseCase) GetByCode(ctx context.Context, code string) (*model.AIWorkflow, error) {
	return uc.repo.GetByCode(ctx, code)
}

func (uc *AIWorkflowUseCase) Create(ctx context.Context, m *model.AIWorkflow) error {
	return uc.repo.Create(ctx, m)
}

func (uc *AIWorkflowUseCase) Update(ctx context.Context, m *model.AIWorkflow) error {
	return uc.repo.Update(ctx, m)
}

func (uc *AIWorkflowUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
