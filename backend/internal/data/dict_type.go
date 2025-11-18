/* *
 * @Author: chengjiang
 * @Date: 2025-11-17 14:13:32
 * @Description: 字典类型数据访问层
**/
package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
	"go.uber.org/zap"
)

type dictTypeRepo struct {
	db    DataRepo
	log   *zap.Logger
	query *query.Query  // 存储预编译的查询实例，避免重复获取DB连接
}

func NewDictTypeRepo(db DataRepo, log *zap.Logger) base.DictTypeRepo {
	// 单例模式：复用查询实例，避免频繁获取数据库连接
	return &dictTypeRepo{
		db:    db,
		log:   log,
		query: query.Use(db.GetDB()),
	}
}

// GetAll implements base.DictTypeRepo
func (r *dictTypeRepo) GetAll(ctx context.Context) ([]*model.SysDictType, error) {
	dictTypes, err := r.query.SysDictType.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return dictTypes, nil
}

// GetByTypeCode implements base.DictTypeRepo
func (r *dictTypeRepo) GetByTypeCode(ctx context.Context, typeCode string) (*model.SysDictType, error) {
	dictType, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.DictType.Eq(typeCode)).First()
	if err != nil {
		return nil, err
	}
	return dictType, nil
}

// Add implements base.DictTypeRepo
func (r *dictTypeRepo) Add(ctx context.Context, dictType *model.SysDictType) error {
	err := r.query.SysDictType.WithContext(ctx).Create(dictType)
	if err != nil {
		r.log.Error("Add dictType failed", zap.Error(err))
		return err
	}
	return nil
}

// Update implements base.DictTypeRepo
func (r *dictTypeRepo) Update(ctx context.Context, dictType *model.SysDictType) error {
	rowsAffected, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.ID.Eq(dictType.ID)).Updates(dictType)
	if err != nil {
		r.log.Error("Update dictType failed", zap.Error(err))
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		return nil
	}
	return nil
}

// Remove implements base.DictTypeRepo
func (r *dictTypeRepo) Remove(ctx context.Context, id int64) error {
	rowsAffected, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.ID.Eq(id)).Delete()
	if err != nil {
		r.log.Error("Remove dictType failed", zap.Error(err))
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		return nil
	}
	return nil
}

// GetById implements base.DictTypeRepo
func (r *dictTypeRepo) GetById(ctx context.Context, id int64) (*model.SysDictType, error) {
	dictType, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return dictType, nil
}

// DictTypeList implements base.DictTypeRepo
func (r *dictTypeRepo) DictTypeList(ctx context.Context, pageNum, pageSize int32, dictType, dictName string) ([]*model.SysDictType, int64, error) {
	offset := (pageNum - 1) * pageSize
	dictTypes, count, err := r.query.SysDictType.WithContext(ctx).Where(
		r.query.SysDictType.DictType.Like("%" + dictType + "%"),
		r.query.SysDictType.DictName.Like("%" + dictName + "%"),
	).FindByPage(int(offset), int(pageSize))
	if err != nil {
		return nil, 0, err
	}
	return dictTypes, count, nil
}

// DictTypeById implements base.DictTypeRepo
func (r *dictTypeRepo) DictTypeById(ctx context.Context, id int64) (*model.SysDictType, error) {
	dictType, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return dictType, nil
}

// AddDictType implements base.DictTypeRepo
func (r *dictTypeRepo) AddDictType(ctx context.Context, dictType *model.SysDictType) error {
	err := r.query.SysDictType.WithContext(ctx).Create(dictType)
	if err != nil {
		r.log.Error("Add dict type failed", zap.Error(err))
		return err
	}
	return nil
}

// UpdateDictType implements base.DictTypeRepo
func (r *dictTypeRepo) UpdateDictType(ctx context.Context, dictType *model.SysDictType) error {
	rowsAffected, err := r.query.SysDictType.WithContext(ctx).Where(r.query.SysDictType.ID.Eq(dictType.ID)).Updates(dictType)
	if err != nil {
		r.log.Error("Update dict type failed", zap.Error(err))
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		return nil
	}
	return nil
}

// DeleteDictType implements base.DictTypeRepo.
func (r *dictTypeRepo) DeleteDictType(ctx context.Context, id int64) error {
	// 开始事务 - 这里需要获取新的连接，因为预编译查询可能不适用于事务
	tx := r.query.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	
	// 确保事务在函数结束时被关闭
	defer func() {
		if rc := recover(); rc != nil {
			tx.Rollback()
			panic(rc)
		}
	}()
	
	u := tx.SysDictType
	dictType, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// 先删除该字典类型下的所有字典数据
	dictU := tx.SysDict
	_, err = dictU.Where(dictU.DictType.Eq(dictType.DictType)).Delete()
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除字典类型
	rowsAffected, err := u.Where(u.ID.Eq(id)).Delete()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		tx.Rollback()
		return nil
	}
	
	// 提交事务
	return tx.Commit()
}