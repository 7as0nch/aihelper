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
	"go.uber.org/zap"
)

type sysMenuRepo struct {
	db DataRepo
	log *zap.Logger
}



func NewSysMenuRepo(db DataRepo, log *zap.Logger) base.SysMenuRepo {
	return &sysMenuRepo{db: db, log: log}
}

// GetAll implements base.SysMenuRepo.
func (s *sysMenuRepo) GetAll(ctx context.Context) ([]*model.SysMenu, error) {
	u := query.Use(s.db.GetDB()).SysMenu
	u.Order(u.Sort.Asc())
	menus, err := u.Find()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// GetRouter implements base.SysMenuRepo.
func (s *sysMenuRepo) GetRouter(ctx context.Context) ([]*model.SysMenu, error) {
	u := query.Use(s.db.GetDB()).SysMenu
	menus, err := u.WithContext(ctx).Where(u.Type.In(uint8(model.MenuTypeDir), uint8(model.MenuTypeMenu))).Find()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// Add implements base.SysMenuRepo.
func (s *sysMenuRepo) Add(ctx context.Context, menu *model.SysMenu) error {
	u := query.Use(s.db.GetDB()).SysMenu
	err := u.WithContext(ctx).Create(menu)
	if err != nil {
		s.log.Error("Add sys menu failed", zap.Error(err))
		return err
	}
	return nil
}

// Update implements base.SysMenuRepo.
func (s *sysMenuRepo) Update(ctx context.Context, menu *model.SysMenu) error {
	u := query.Use(s.db.GetDB()).SysMenu
	rowsAffected, err := u.WithContext(ctx).Where(u.ID.Eq(menu.ID)).Updates(menu)
	if err != nil {
		s.log.Error("Update sys menu failed", zap.Error(err))
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		return nil
	}
	return nil
}

// Delete implements base.SysMenuRepo.
func (s *sysMenuRepo) Delete(ctx context.Context, id int64) error {
	u := query.Use(s.db.GetDB()).SysMenu
	rowsAffected, err := u.WithContext(ctx).
	Where(u.ID.Eq(id)).
	Or(u.ParentID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	if rowsAffected.RowsAffected == 0 {
		return nil
	}
	return nil
}

// Get implements base.SysMenuRepo.
func (s *sysMenuRepo) Get(ctx context.Context, id int64) (*model.SysMenu, error) {
	u := query.Use(s.db.GetDB()).SysMenu
	menu, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return menu, nil
}