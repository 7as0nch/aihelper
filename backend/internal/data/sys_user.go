/* *
 * @Author: chengjiang
 * @Date: 2025-11-17 14:13:32
 * @Description: 系统用户数据访问层
**/
package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
)

type sysUserRepo struct {
	db    DataRepo
	query *query.Query  // 存储预编译的查询实例
}

func NewSysUserRepo(db DataRepo) base.SysUserRepo {
	// 预编译查询实例，避免每次都重新创建
	return &sysUserRepo{
		db:    db,
		query: query.Use(db.GetDB()),
	}
}

// GetByAccount implements base.SysUserRepo
func (r *sysUserRepo) GetByAccount(ctx context.Context, account string) (*model.SysUser, error) {
	user, err := r.query.SysUser.WithContext(ctx).Where(r.query.SysUser.Name.Eq(account)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetById implements base.SysUserRepo
func (r *sysUserRepo) GetById(ctx context.Context, id int64) (*model.SysUser, error) {
	user, err := r.query.SysUser.WithContext(ctx).Where(r.query.SysUser.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}