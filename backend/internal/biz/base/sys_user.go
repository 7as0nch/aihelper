package base

import (
	"context"
	"errors"

	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/pkg/auth"
	"github.com/go-kratos/kratos/v2/log"
)

/* *
 * @Author: chengjiang
 * @Date: 2025-11-14 17:31:58
 * @Description:
**/

type SysUserRepo interface {
	GetByAccount(ctx context.Context, account string) (*model.SysUser, error)
	GetById(ctx context.Context, id int64) (*model.SysUser, error)
}

type SysUserUseCase struct {
	user SysUserRepo
	auth auth.AuthRepo
}

func NewSysUserUseCase(user SysUserRepo, auth auth.AuthRepo) *SysUserUseCase {
	return &SysUserUseCase{
		user: user,
		auth: auth,
	}
}

// Login
func (s *SysUserUseCase) Login(ctx context.Context, account string) (string, error) {
	user, err := s.user.GetByAccount(ctx, account)
	if err != nil {
		return "", err
	}
	token, err := s.auth.NewToken(ctx, user.ID, user.Account, user.Phonenumber)
	if err != nil {
		log.Errorf("获取Token失败: %v", err)
		return "", errors.New("获取Token失败")
	}
	return token, nil
}

// GetInfo
func (s *SysUserUseCase) GetInfo(ctx context.Context) (*model.SysUser, error) {
	uid, _ := ctx.Value(auth.UserId).(int64)
	user, err := s.user.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Logout
func (s *SysUserUseCase) Logout(ctx context.Context) error {
	// TODO 实现登出逻辑
	return nil
}
