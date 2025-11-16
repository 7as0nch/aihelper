/* *
 * @Author: chengjiang
 * @Date: 2025-11-16 17:45:51
 * @Description:
**/
package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
)

type sysMenuRepo struct {
	db *Data
}


func NewSysMenuRepo(db *Data) base.SysMenuRepo {
	return &sysMenuRepo{db: db}
}

// GetAll implements base.SysMenuRepo.
func (s *sysMenuRepo) GetAll(ctx context.Context) ([]*model.SysMenu, error) {
	u := query.Use(s.db.db).SysMenu
	menus, err := u.Find()
	if err != nil {
		return nil, err
	}
	return menus, nil
}