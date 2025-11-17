/* *
 * @Author: chengjiang
 * @Date: 2025-11-14 18:37:46
 * @Description:
**/

package data

import (
	"context"
	"errors"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type sysUserRepo struct {
	db DataRepo
}


func NewSysUserRepo(db DataRepo) base.SysUserRepo {
	return &sysUserRepo{
		db: db,
	}
}

// GetByAccount implements base.SysUserRepo.
func (s *sysUserRepo) GetByAccount(ctx context.Context, account string) (*model.SysUser, error) {
	u := query.Use(s.db.GetDB()).SysUser
	user, err := u.Where(u.Name.Eq(account)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *sysUserRepo) GetById(ctx context.Context, id int64) (*model.SysUser, error) {
	u := query.Use(s.db.GetDB()).SysUser
	user, err := u.Where(u.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Errorf("get user by id error: %v", err)
		return nil, err
	}
	return user, nil
}