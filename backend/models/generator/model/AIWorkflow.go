package model

import (
	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/pkg/ai"
)

type AIWorkflow struct {
	models.Model
	Name        string             `json:"name" gorm:"size:100;not null;comment:'名称'"`
	Code        string             `json:"code" gorm:"uniqueIndex;size:100;not null;comment:'编码'"`
	Description string             `json:"description" gorm:"size:500;comment:'描述'"`
	Definition  *ai.WorkflowConfig `json:"definition" gorm:"type:text;comment:'工作流定义'"`
	Status      models.Status      `json:"status" gorm:"default:1;comment:'状态'"`
}

func (AIWorkflow) TableName() string {
	return "ai_workflow"
}
