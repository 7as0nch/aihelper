package models

import (
	"time"
	"gorm.io/gorm"
)

// 用户角色
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// User 用户模型
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"size:50;unique;not null" json:"username"`
	Password      string         `gorm:"size:255;not null" json:"-"` // 不返回密码
	Email         string         `gorm:"size:100;unique;not null" json:"email"`
	Role          string         `gorm:"size:20;default:'user'" json:"role"`
	Status        int            `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	Avatar        string         `gorm:"size:255" json:"avatar"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段，不返回

	// 关联字段
	Chats        []Chat          `json:"chats"`
	FunctionTools []FunctionTool `json:"function_tools"`
	Workflows    []Workflow      `json:"workflows"`
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加密码加密等逻辑
	return nil
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}