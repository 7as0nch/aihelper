/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AITool 业务逻辑层
**/
package ai

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
)

// AIToolRepo AITool 数据仓库接口
type AIToolRepo interface {
	List(ctx context.Context, page, pageSize int32, name string, toolType, status int) ([]*model.AITool, int64, error)
	GetByID(ctx context.Context, id int64) (*model.AITool, error)
	GetByCode(ctx context.Context, code string) (*model.AITool, error)
	Create(ctx context.Context, t *model.AITool) error
	Update(ctx context.Context, t *model.AITool) error
	Delete(ctx context.Context, id int64) error
	// ToolAgentBind 操作
	GetToolCodesByAgentID(ctx context.Context, agentID int64) ([]string, error)
	BatchBindTools(ctx context.Context, agentID int64, toolCodes []string) error
	DeleteToolBinds(ctx context.Context, agentID int64) error
	DeleteToolBindsByCode(ctx context.Context, toolCode string) error
}

// AIToolUseCase AITool 业务逻辑
type AIToolUseCase struct {
	repo AIToolRepo
	tm   base.Transaction
}

// NewAIToolUseCase 创建 AITool UseCase
func NewAIToolUseCase(repo AIToolRepo, tm base.Transaction) *AIToolUseCase {
	return &AIToolUseCase{
		repo: repo,
		tm:   tm,
	}
}

// List 获取工具列表
func (uc *AIToolUseCase) List(ctx context.Context, page, pageSize int32, name string, toolType, status int) ([]*model.AITool, int64, error) {
	return uc.repo.List(ctx, page, pageSize, name, toolType, status)
}

// GetByID 根据 ID 获取工具
func (uc *AIToolUseCase) GetByID(ctx context.Context, id int64) (*model.AITool, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByCode 根据 Code 获取工具
func (uc *AIToolUseCase) GetByCode(ctx context.Context, code string) (*model.AITool, error) {
	return uc.repo.GetByCode(ctx, code)
}

// Create 创建工具
func (uc *AIToolUseCase) Create(ctx context.Context, t *model.AITool) error {
	return uc.repo.Create(ctx, t)
}

// Update 更新工具
func (uc *AIToolUseCase) Update(ctx context.Context, t *model.AITool) error {
	return uc.repo.Update(ctx, t)
}

// Delete 删除工具（事务：先删绑定关系，再删工具）
func (uc *AIToolUseCase) Delete(ctx context.Context, id int64) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 1. 获取工具 code
		tool, err := uc.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		// 2. 删除工具绑定
		if err := uc.repo.DeleteToolBindsByCode(ctx, tool.Code); err != nil {
			return err
		}
		// 3. 删除工具
		return uc.repo.Delete(ctx, id)
	})
}

// GetToolCodesByAgentID 获取 Agent 绑定的工具 Codes
func (uc *AIToolUseCase) GetToolCodesByAgentID(ctx context.Context, agentID int64) ([]string, error) {
	return uc.repo.GetToolCodesByAgentID(ctx, agentID)
}

// BatchBindTools 批量绑定工具到 Agent
func (uc *AIToolUseCase) BatchBindTools(ctx context.Context, agentID int64, toolCodes []string) error {
	return uc.tm.InTx(ctx, func(ctx context.Context) error {
		// 先删除旧绑定
		if err := uc.repo.DeleteToolBinds(ctx, agentID); err != nil {
			return err
		}
		// 再创建新绑定
		if len(toolCodes) > 0 {
			return uc.repo.BatchBindTools(ctx, agentID, toolCodes)
		}
		return nil
	})
}
