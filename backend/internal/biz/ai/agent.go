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
	// AgentBind 操作
	GetSubAgentIDs(ctx context.Context, agentID int64) ([]int64, error)
	BatchBindSubAgents(ctx context.Context, agentID int64, subAgentIDs []int64) error
	DeleteSubAgentBinds(ctx context.Context, agentID int64) error
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

// GetSubAgentIDs 获取子 Agent IDs
func (uc *AIAgentUseCase) GetSubAgentIDs(ctx context.Context, agentID int64) ([]int64, error) {
	return uc.repo.GetSubAgentIDs(ctx, agentID)
}

// Create 创建 Agent
func (uc *AIAgentUseCase) Create(ctx context.Context, agent *model.AIAgent) error {
	return uc.repo.Create(ctx, agent)
}

// CreateWithSubAgents 创建 Agent 并绑定子 Agent（事务）
func (uc *AIAgentUseCase) CreateWithSubAgents(ctx context.Context, agent *model.AIAgent, subAgentIDs []int64) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		if err := uc.repo.Create(ctx, agent); err != nil {
			return err
		}
		return nil
	})
}

// Update 更新 Agent
func (uc *AIAgentUseCase) Update(ctx context.Context, agent *model.AIAgent) error {
	return uc.repo.Update(ctx, agent)
}

// UpdateWithSubAgents 更新 Agent 并重新绑定子 Agent（事务）
func (uc *AIAgentUseCase) UpdateWithSubAgents(ctx context.Context, agent *model.AIAgent, subAgentIDs []int64) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		if err := uc.repo.Update(ctx, agent); err != nil {
			return err
		}
		// 先删除旧绑定
		if err := uc.repo.DeleteSubAgentBinds(ctx, agent.ID); err != nil {
			return err
		}
		// 再创建新绑定
		if len(subAgentIDs) > 0 {
			if err := uc.repo.BatchBindSubAgents(ctx, agent.ID, subAgentIDs); err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除 Agent（事务：先删绑定关系，再删 Agent）
func (uc *AIAgentUseCase) Delete(ctx context.Context, id int64) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 1. 删除子 Agent 绑定
		if err := uc.repo.DeleteSubAgentBinds(ctx, id); err != nil {
			return err
		}
		// 2. 删除 Agent
		return uc.repo.Delete(ctx, id)
	})
}

// BatchBindSubAgents 批量绑定子 Agent
func (uc *AIAgentUseCase) BatchBindSubAgents(ctx context.Context, agentID int64, subAgentIDs []int64) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 先删除旧绑定
		if err := uc.repo.DeleteSubAgentBinds(ctx, agentID); err != nil {
			return err
		}
		// 再创建新绑定
		if len(subAgentIDs) > 0 {
			return uc.repo.BatchBindSubAgents(ctx, agentID, subAgentIDs)
		}
		return nil
	})
}
