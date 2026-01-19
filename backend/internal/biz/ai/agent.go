/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIAgent 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
)

// AIAgentRepo AIAgent 数据仓库接口
type AIAgentRepo interface {
	List(ctx context.Context, page, pageSize int32, name string, status, agentType int) ([]*model.AIAgent, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AIAgent, error)
	GetByCode(ctx context.Context, code string) (*model.AIAgent, error)
	Create(ctx context.Context, agent *model.AIAgent) error
	Update(ctx context.Context, agent *model.AIAgent) error
	Delete(ctx context.Context, id int64) error
}

// AIAgentUseCase AIAgent 业务逻辑
type AIAgentUseCase struct {
	repo AIAgentRepo
	tm   base.Transaction
}

// NewAIAgentUseCase 创建 AIAgent UseCase
func NewAIAgentUseCase(repo AIAgentRepo, tm base.Transaction) *AIAgentUseCase {
	return &AIAgentUseCase{
		repo: repo,
		tm:   tm,
	}
}

// List 获取 Agent 列表
func (uc *AIAgentUseCase) List(ctx context.Context, page, pageSize int32, name string, status, agentType int) ([]*model.AIAgent, int64, error) {
	return uc.repo.List(ctx, page, pageSize, name, status, agentType)
}

// GetByID 根据 ID 获取 Agent
func (uc *AIAgentUseCase) GetByID(ctx context.Context, id int64) (*model.AIAgent, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByCode 根据 Code 获取 Agent
func (uc *AIAgentUseCase) GetByCode(ctx context.Context, code string) (*model.AIAgent, error) {
	return uc.repo.GetByCode(ctx, code)
}

// Delete 删除 Agent
func (uc *AIAgentUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

// Create 创建 Agent
func (uc *AIAgentUseCase) Create(ctx context.Context, agent *model.AIAgent) error {
	return uc.repo.Create(ctx, agent)
}

// Update 更新 Agent
func (uc *AIAgentUseCase) Update(ctx context.Context, agent *model.AIAgent) error {
	return uc.repo.Update(ctx, agent)
}
