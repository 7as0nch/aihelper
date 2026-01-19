/* *
 * @Author: chengjiang
 * @Date: 2026-01-19
 * @Description: AIApplication 数据访问层
**/
package ai

import (
	"context"

	bizai "github.com/example/aichat/backend/internal/biz/ai"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
)

type aiApplicationRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

// NewAIApplicationRepo 创建 AIApplication 仓库
func NewAIApplicationRepo(db db.DataRepo, log *zap.Logger) bizai.AIApplicationRepo {
	return &aiApplicationRepo{
		db:  db,
		log: log,
	}
}

func (r *aiApplicationRepo) List(ctx context.Context, page, pageSize int32, name string, status int) ([]*model.AIApplication, int64, error) {
	offset := (page - 1) * pageSize
	db := r.db.DB(ctx).Model(&model.AIApplication{})

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

	var list []*model.AIApplication
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *aiApplicationRepo) GetByID(ctx context.Context, id int64) (*model.AIApplication, error) {
	var m model.AIApplication
	if err := r.db.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *aiApplicationRepo) GetByCode(ctx context.Context, code string) (*model.AIApplication, error) {
	var m model.AIApplication
	if err := r.db.DB(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *aiApplicationRepo) Create(ctx context.Context, m *model.AIApplication) error {
	return r.db.DB(ctx).Create(m).Error
}

func (r *aiApplicationRepo) Update(ctx context.Context, m *model.AIApplication) error {
	return r.db.DB(ctx).Model(m).Where("id = ?", m.ID).Updates(m).Error
}

func (r *aiApplicationRepo) Delete(ctx context.Context, id int64) error {
	return r.db.DB(ctx).Delete(&model.AIApplication{}, id).Error
}
