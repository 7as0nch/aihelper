/* *
 * @Author: chengjiang
 * @Date: 2025-11-13 18:22:42
 * @Description: 用户表
**/
package model

import "github.com/example/aichat/backend/models"

const TableNameUser = "user"

type User struct {
	models.Model
	Name     string `gorm:"column:name" json:"name"`         // 名称
	Password string `gorm:"column:password" json:"password"` // 密码
}

// TableName SysMessage's table name
func (*User) TableName() string {
	return TableNameUser
}
