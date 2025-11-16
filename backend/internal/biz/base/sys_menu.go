/* *
 * @Author: chengjiang
 * @Date: 2025-11-16 20:13:08
 * @Description:
**/
package base

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

type SysMenuRepo interface {
	GetAll(ctx context.Context) ([]*model.SysMenu, error)
}

type SysMenuUseCase struct {
	menu SysMenuRepo
}

func NewSysMenuUseCase(menu SysMenuRepo) *SysMenuUseCase {
	return &SysMenuUseCase{
		menu: menu,
	}
}

// GetAll
func (s *SysMenuUseCase) GetAll(ctx context.Context) ([]*model.SysMenu, error) {
	return s.menu.GetAll(ctx)
}
