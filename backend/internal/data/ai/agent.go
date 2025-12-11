/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIAgent 数据访问层
**/
package ai

import (
	"context"

	bizai "github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type aiAgentRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIAgentRepo 创建 AIAgent 仓库
func NewAIAgentRepo(db db.DataRepo, log *zap.Logger) bizai.AIAgentRepo {
	return &aiAgentRepo{
		db:  db,
		log: log,
	}
}

// List 获取 Agent 列表
func (r *aiAgentRepo) List(ctx context.Context, page, pageSize int32, name string, status, agentType int) ([]*model.AIAgent, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AIAgent{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}
	if agentType > 0 {
		db = db.Where("type = ?", agentType)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var agents []*model.AIAgent
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// GetByID 根据 ID 获取 Agent
func (r *aiAgentRepo) GetByID(ctx context.Context, id int64) (*model.AIAgent, error) {
	var agent model.AIAgent
	if err := r.db.DB(ctx).First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByCode 根据 Code 获取 Agent
func (r *aiAgentRepo) GetByCode(ctx context.Context, code string) (*model.AIAgent, error) {
	var agent model.AIAgent
	if err := r.db.DB(ctx).Where("code = ?", code).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// Create 创建 Agent
func (r *aiAgentRepo) Create(ctx context.Context, agent *model.AIAgent) error {
	return r.db.DB(ctx).Create(agent).Error
}

// Update 更新 Agent
func (r *aiAgentRepo) Update(ctx context.Context, agent *model.AIAgent) error {
	return r.db.DB(ctx).Model(agent).Where("id = ?", agent.ID).Updates(agent).Error
}

// Delete 删除 Agent
func (r *aiAgentRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AIAgent{}, id).Error
}

// GetSubAgentIDs 获取子 Agent IDs
func (r *aiAgentRepo) GetSubAgentIDs(ctx context.Context, agentID int64) ([]int64, error) {
	var binds []model.AgentBind
	if err := r.db.DB(ctx).Where("agent_id = ?", agentID).Find(&binds).Error; err != nil {
		return nil, err
	}
	ids := make([]int64, len(binds))
	for i, b := range binds {
		ids[i] = b.SubAgentID
	}
	return ids, nil
}

// BatchBindSubAgents 批量绑定子 Agent
func (r *aiAgentRepo) BatchBindSubAgents(ctx context.Context, agentID int64, subAgentIDs []int64) error {
	if len(subAgentIDs) == 0 {
		return nil
	}
	binds := make([]*model.AgentBind, len(subAgentIDs))
	for i, subID := range subAgentIDs {
		binds[i] = &model.AgentBind{
			Model:      models.Model{},
			AgentID:    agentID,
			SubAgentID: subID,
		}
		binds[i].New()
	}
	return r.db.DB(ctx).Create(&binds).Error
}

// DeleteSubAgentBinds 删除子 Agent 绑定
func (r *aiAgentRepo) DeleteSubAgentBinds(ctx context.Context, agentID int64) error {
	return r.db.DB(ctx).Where("agent_id = ?", agentID).Delete(&model.AgentBind{}).Error
}

// AutoMigrate 自动迁移表结构
func (r *aiAgentRepo) AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AIAgent{}, &model.AgentBind{})
}
