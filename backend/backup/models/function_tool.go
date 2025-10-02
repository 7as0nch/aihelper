package models

import (
	"time"
	"gorm.io/gorm"
)

// FunctionToolType 函数工具类型
type FunctionToolType string

const (
	FunctionToolTypeHTTP  FunctionToolType = "http"
	FunctionToolTypeMCP   FunctionToolType = "mcp"
	FunctionToolTypeLocal FunctionToolType = "local"
)

// FunctionTool 函数工具模型
type FunctionTool struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"` // 0表示系统级工具
	Name        string         `gorm:"size:100;not null;unique" json:"name"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Type        FunctionToolType `gorm:"size:20;not null" json:"type"`
	Config      string         `gorm:"type:text;not null" json:"config"` // JSON格式的配置
	Parameters  string         `gorm:"type:text;not null" json:"parameters"` // JSON Schema格式的参数定义
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联字段
	WorkflowSteps []WorkflowStep `gorm:"foreignKey:FunctionToolID" json:"workflow_steps"`
}

// TableName 设置表名
func (FunctionTool) TableName() string {
	return "function_tools"
}