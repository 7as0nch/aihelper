/* *
 * @Author: chengjiang
 * @Date: 2026-01-19
 * @Description: AIWorkflow 数据访问层
**/
package ai

import (
	"context"

	bizai "github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
)

type aiWorkflowRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIWorkflowRepo 创建 AIWorkflow 仓库
func NewAIWorkflowRepo(db db.DataRepo, log *zap.Logger) bizai.AIWorkflowRepo {
	return &aiWorkflowRepo{
		db:  db,
		log: log,
	}
}

func (r *aiWorkflowRepo) List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIWorkflow, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AIWorkflow{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.AIWorkflow
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *aiWorkflowRepo) GetByID(ctx context.Context, id int64) (*model.AIWorkflow, error) {
	var m model.AIWorkflow
	if err := r.db.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *aiWorkflowRepo) GetByCode(ctx context.Context, code string) (*model.AIWorkflow, error) {
	var m model.AIWorkflow
	if err := r.db.DB(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *aiWorkflowRepo) Create(ctx context.Context, m *model.AIWorkflow) error {
	return r.db.DB(ctx).Create(m).Error
}

func (r *aiWorkflowRepo) Update(ctx context.Context, m *model.AIWorkflow) error {
	return r.db.DB(ctx).Model(m).Where("id = ?", m.ID).Updates(m).Error
}

func (r *aiWorkflowRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AIWorkflow{}, id).Error
}
