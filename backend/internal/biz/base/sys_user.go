package base

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-11-14 17:31:58
 * @Description:
**/

type SysUserRepo interface {
	GetByAccount(ctx context.Context, account string) (*model.SysUser, error)
}

type SysUserUseCase struct {
	user SysUserRepo
}

func NewSysUserUseCase(user SysUserRepo) *SysUserUseCase {
	return &SysUserUseCase{
		user: user,
	}
}

// Login
func (s *SysUserUseCase) Login(ctx context.Context, account string) (*model.SysUser, error) {
	user, err := s.user.GetByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return user, nil
}
