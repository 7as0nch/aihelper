/* *
 * @Author: chengjiang
 * @Date: 2026-01-19
 * @Description: AIApplication 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

// AIApplicationRepo AIApplication 数据仓库接口
type AIApplicationRepo interface {
	List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIApplication, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AIApplication, error)
	GetByCode(ctx context.Context, code string) (*model.AIApplication, error)
	Create(ctx context.Context, m *model.AIApplication) error
	Update(ctx context.Context, m *model.AIApplication) error
	Delete(ctx context.Context, id int64) error
}

// AIApplicationUseCase AIApplication 业务逻辑
type AIApplicationUseCase struct {
	repo AIApplicationRepo
}

// NewAIApplicationUseCase 创建 AIApplication UseCase
func NewAIApplicationUseCase(repo AIApplicationRepo) *AIApplicationUseCase {
	return &AIApplicationUseCase{
		repo: repo,
	}
}

func (uc *AIApplicationUseCase) List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIApplication, int64, error) {
	return uc.repo.List(ctx, page, pageSize, name, status)
}

func (uc *AIApplicationUseCase) GetByID(ctx context.Context, id int64) (*model.AIApplication, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *AIApplicationUseCase) GetByCode(ctx context.Context, code string) (*model.AIApplication, error) {
	return uc.repo.GetByCode(ctx, code)
}

func (uc *AIApplicationUseCase) Create(ctx context.Context, m *model.AIApplication) error {
	return uc.repo.Create(ctx, m)
}

func (uc *AIApplicationUseCase) Update(ctx context.Context, m *model.AIApplication) error {
	return uc.repo.Update(ctx, m)
}

func (uc *AIApplicationUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
