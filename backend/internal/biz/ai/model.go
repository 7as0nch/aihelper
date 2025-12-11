/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIModel 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

// AIModelRepo AIModel 数据仓库接口
type AIModelRepo interface {
	List(ctx context.Context, page, pageSize int32, modelName string, status int) ([]*model.AIModel, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AIModel, error)
	GetDefault(ctx context.Context) (*model.AIModel, error)
	Create(ctx context.Context, m *model.AIModel) error
	Update(ctx context.Context, m *model.AIModel) error
	Delete(ctx context.Context, id int64) error
	SetDefault(ctx context.Context, id int64) error
	ClearDefault(ctx context.Context) error
}

// AIModelUseCase AIModel 业务逻辑
type AIModelUseCase struct {
	repo AIModelRepo
}

// NewAIModelUseCase 创建 AIModel UseCase
func NewAIModelUseCase(repo AIModelRepo) *AIModelUseCase {
	return &AIModelUseCase{
		repo: repo,
	}
}

// List 获取模型列表
func (uc *AIModelUseCase) List(ctx context.Context, page, pageSize int32, modelName string, status int) ([]*model.AIModel, int64, error) {
	return uc.repo.List(ctx, page, pageSize, modelName, status)
}

// GetByID 根据 ID 获取模型
func (uc *AIModelUseCase) GetByID(ctx context.Context, id int64) (*model.AIModel, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetDefault 获取默认模型
func (uc *AIModelUseCase) GetDefault(ctx context.Context) (*model.AIModel, error) {
	return uc.repo.GetDefault(ctx)
}

// Create 创建模型
func (uc *AIModelUseCase) Create(ctx context.Context, m *model.AIModel) error {
	return uc.repo.Create(ctx, m)
}

// Update 更新模型
func (uc *AIModelUseCase) Update(ctx context.Context, m *model.AIModel) error {
	return uc.repo.Update(ctx, m)
}

// Delete 删除模型
func (uc *AIModelUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
