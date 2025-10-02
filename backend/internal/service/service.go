package service

/* *
 * @Author: chengjiang
 * @Date: 2025-10-01 22:42:39
 * @Description:
**/

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)
