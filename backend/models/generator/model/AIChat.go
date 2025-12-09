package model

import "github.com/example/aichat/backend/models"

const TableNameAIChat = "ai_chat"

// AIChat mapped from table <ai_chat>
type AIChat struct {
	models.Model
	Title  string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	UserID int64  `gorm:"column:user_id;type:bigint;not null;index" json:"user_id"`
	// todo 是否涉及转人工处理。
}

// TableName AIChat's table name
func (*AIChat) TableName() string {
	return TableNameAIChat
}
