/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIPromptTemplate 数据访问层
**/
package ai

import (
	"context"

	bizai "github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type aiPromptRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIPromptRepo 创建 AIPrompt 仓库
func NewAIPromptRepo(db db.DataRepo, log *zap.Logger) bizai.AIPromptRepo {
	return &aiPromptRepo{
		db:  db,
		log: log,
	}
}

// List 获取提示词模板列表
func (r *aiPromptRepo) List(ctx context.Context, page, pageSize int32, name string, promptType int) ([]*model.AIPromptTemplate, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AIPromptTemplate{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if promptType > 0 {
		db = db.Where("type = ?", promptType)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var prompts []*model.AIPromptTemplate
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&prompts).Error; err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

// GetByID 根据 ID 获取提示词模板
func (r *aiPromptRepo) GetByID(ctx context.Context, id int64) (*model.AIPromptTemplate, error) {
	var p model.AIPromptTemplate
	if err := r.db.DB(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Create 创建提示词模板
func (r *aiPromptRepo) Create(ctx context.Context, p *model.AIPromptTemplate) error {
	return r.db.DB(ctx).Create(p).Error
}

// Update 更新提示词模板
func (r *aiPromptRepo) Update(ctx context.Context, p *model.AIPromptTemplate) error {
	return r.db.DB(ctx).Model(p).Where("id = ?", p.ID).Updates(p).Error
}

// Delete 删除提示词模板
func (r *aiPromptRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AIPromptTemplate{}, id).Error
}

// AutoMigrate 自动迁移表结构
func (r *aiPromptRepo) AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AIPromptTemplate{})
}
