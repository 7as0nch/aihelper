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
	Type        SysUserType `gorm:"column:type;type:varchar(5)" json:"type"`                // 用户类型（00系统用户 01普通用户）
	Name        string      `gorm:"column:name;type:varchar(50)" json:"name"`               // 名称
	Account     string      `gorm:"column:account;type:varchar(50)" json:"account"`         // 昵称
	Password    string      `gorm:"column:password;type:varchar(50)" json:"password"`       // 密码
	Avatar      string      `gorm:"column:avatar;type:varchar(255)" json:"avatar"`          // 头像
	Email       string      `gorm:"column:email;type:varchar(50)" json:"email"`             // 邮箱
	Phonenumber string      `gorm:"column:phonenumber;type:varchar(20)" json:"phonenumber"` // 手机号
	Remark      string      `gorm:"column:remark;type:varchar(255)" json:"remark"`          // 备注
	Sex         string      `gorm:"column:sex;type:varchar(1)" json:"sex"`
	Status      string      `gorm:"column:status;type:varchar(1)" json:"status"`
}

// TableName SysMessage's table name
func (*SysUser) TableName() string {
	return TableNameSysUser
}

type SysUserType string

const (
	SysUserType_System SysUserType = "Sys" // 系统用户
	SysUserType_Admin  SysUserType = "Adm" // 管理员用户
	SysUserType_Guest  SysUserType = "Gst" // 普通用户
)
