/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AIModel 数据访问层
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

type aiModelRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIModelRepo 创建 AIModel 仓库
func NewAIModelRepo(db db.DataRepo, log *zap.Logger) bizai.AIModelRepo {
	return &aiModelRepo{
		db:  db,
		log: log,
	}
}

// List 获取模型列表
func (r *aiModelRepo) List(ctx context.Context, page, pageSize int32, modelName string, status int) ([]*model.AIModel, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AIModel{})

	if modelName != "" {
		db = db.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.AIModel
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

// GetByID 根据 ID 获取模型
func (r *aiModelRepo) GetByID(ctx context.Context, id int64) (*model.AIModel, error) {
	var m model.AIModel
	if err := r.db.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetDefault 获取默认模型
func (r *aiModelRepo) GetDefault(ctx context.Context) (*model.AIModel, error) {
	var m model.AIModel
	if err := r.db.DB(ctx).Where("is_default = ?", models.Status_Enabled).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Create 创建模型
func (r *aiModelRepo) Create(ctx context.Context, m *model.AIModel) error {
	return r.db.DB(ctx).Create(m).Error
}

// Update 更新模型
func (r *aiModelRepo) Update(ctx context.Context, m *model.AIModel) error {
	return r.db.DB(ctx).Model(m).Where("id = ?", m.ID).Updates(m).Error
}

// Delete 删除模型
func (r *aiModelRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AIModel{}, id).Error
}

// SetDefault 设置默认模型
func (r *aiModelRepo) SetDefault(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Model(&model.AIModel{}).Where("id = ?", id).Update("is_default", models.Status_Enabled).Error
}

// ClearDefault 清除默认模型
func (r *aiModelRepo) ClearDefault(ctx context.Context) error {
	return r.db.DB(ctx).Model(&model.AIModel{}).Where("is_default = ?", models.Status_Enabled).Update("is_default", models.Status_Disabled).Error
}

// AutoMigrate 自动迁移表结构
func (r *aiModelRepo) AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AIModel{})
}
