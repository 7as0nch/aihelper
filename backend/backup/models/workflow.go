package models

import (
	"time"
	"gorm.io/gorm"
)

// Workflow 工作流模型
type Workflow struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"` // 0表示系统级工作流
	Name        string         `gorm:"size:100;not null;unique" json:"name"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联字段
	Steps []WorkflowStep `json:"steps"`
}

// WorkflowStep 工作流步骤模型
type WorkflowStep struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	WorkflowID      uint           `gorm:"index;not null" json:"workflow_id"`
	FunctionToolID  uint           `gorm:"index;not null" json:"function_tool_id"`
	StepName        string         `gorm:"size:100;not null" json:"step_name"`
	Order           int            `gorm:"not null" json:"order"`
	Condition       string         `gorm:"type:text" json:"condition"` // 条件表达式
	VariableMapping string         `gorm:"type:text" json:"variable_mapping"` // 变量映射配置
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       *gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联字段
	FunctionTool FunctionTool `json:"function_tool"`
}

// TableName 设置表名
func (Workflow) TableName() string {
	return "workflows"
}

// TableName 设置表名
func (WorkflowStep) TableName() string {
	return "workflow_steps"
}