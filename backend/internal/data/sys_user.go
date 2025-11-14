/* *
 * @Author: chengjiang
 * @Date: 2025-11-14 18:37:46
 * @Description:
**/

package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
)

type sysUserRepo struct {
	db Data
}


func NewSysUserRepo(db Data) base.SysUserRepo {
	return &sysUserRepo{}
}

// GetByAccount implements base.SysUserRepo.
func (s *sysUserRepo) GetByAccount(ctx context.Context, account string) (*model.SysUser, error) {
	u := query.SysUser
	user, err := u.WithContext(ctx).Where(u.Name.Eq(account)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}