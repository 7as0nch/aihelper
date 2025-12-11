/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AITool 数据访问层
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

type aiToolRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIToolRepo 创建 AITool 仓库
func NewAIToolRepo(db db.DataRepo, log *zap.Logger) bizai.AIToolRepo {
	return &aiToolRepo{
		db:  db,
		log: log,
	}
}

// List 获取工具列表
func (r *aiToolRepo) List(ctx context.Context, page, pageSize int32, name string, toolType, status int) ([]*model.AITool, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AITool{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if toolType > 0 {
		db = db.Where("type = ?", toolType)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tools []*model.AITool
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&tools).Error; err != nil {
		return nil, 0, err
	}

	return tools, total, nil
}

// GetByID 根据 ID 获取工具
func (r *aiToolRepo) GetByID(ctx context.Context, id int64) (*model.AITool, error) {
	var t model.AITool
	if err := r.db.DB(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByCode 根据 Code 获取工具
func (r *aiToolRepo) GetByCode(ctx context.Context, code string) (*model.AITool, error) {
	var t model.AITool
	if err := r.db.DB(ctx).Where("code = ?", code).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// Create 创建工具
func (r *aiToolRepo) Create(ctx context.Context, t *model.AITool) error {
	return r.db.DB(ctx).Create(t).Error
}

// Update 更新工具
func (r *aiToolRepo) Update(ctx context.Context, t *model.AITool) error {
	return r.db.DB(ctx).Model(t).Where("id = ?", t.ID).Updates(t).Error
}

// Delete 删除工具
func (r *aiToolRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AITool{}, id).Error
}

// GetToolCodesByAgentID 获取 Agent 绑定的工具 Codes
func (r *aiToolRepo) GetToolCodesByAgentID(ctx context.Context, agentID int64) ([]string, error) {
	var binds []model.AIToolAgentBind
	if err := r.db.DB(ctx).Where("agent_id = ?", agentID).Find(&binds).Error; err != nil {
		return nil, err
	}
	codes := make([]string, len(binds))
	for i, b := range binds {
		codes[i] = b.ToolCode
	}
	return codes, nil
}

// BatchBindTools 批量绑定工具到 Agent
func (r *aiToolRepo) BatchBindTools(ctx context.Context, agentID int64, toolCodes []string) error {
	if len(toolCodes) == 0 {
		return nil
	}
	binds := make([]*model.AIToolAgentBind, len(toolCodes))
	for i, code := range toolCodes {
		binds[i] = &model.AIToolAgentBind{
			Model:    models.Model{},
			AgentID:  uint64(agentID),
			ToolCode: code,
		}
		binds[i].New()
	}
	return r.db.DB(ctx).Create(&binds).Error
}

// DeleteToolBinds 删除 Agent 的所有工具绑定
func (r *aiToolRepo) DeleteToolBinds(ctx context.Context, agentID int64) error {
	return r.db.DB(ctx).Where("agent_id = ?", agentID).Delete(&model.AIToolAgentBind{}).Error
}

// DeleteToolBindsByCode 根据工具 Code 删除绑定
func (r *aiToolRepo) DeleteToolBindsByCode(ctx context.Context, toolCode string) error {
	return r.db.DB(ctx).Where("tool_code = ?", toolCode).Delete(&model.AIToolAgentBind{}).Error
}

// AutoMigrate 自动迁移表结构
func (r *aiToolRepo) AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AITool{}, &model.AIToolAgentBind{})
}
