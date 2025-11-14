/* *
 * @Author: chengjiang
 * @Date: 2025-11-13 18:22:42
 * @Description: 用户表
**/
package model

import "github.com/example/aichat/backend/models"

const TableNameSysUser = "sys_user"

type SysUser struct {
	models.Model
	Name     string `gorm:"column:name;type:varchar(50)" json:"name"`         // 名称
	Password string `gorm:"column:password;type:varchar(50)" json:"password"` // 密码
}

// TableName SysMessage's table name
func (*SysUser) TableName() string {
	return TableNameSysUser
}
