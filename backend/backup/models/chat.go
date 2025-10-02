package models

import (
	"time"
	"gorm.io/gorm"
)

// 消息类型
const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeAudio = "audio"
	MessageTypeVideo = "video"
	MessageTypeFile  = "file"
)

// 消息来源
const (
	MessageFromUser      = "user"
	MessageFromAssistant = "assistant"
	MessageFromSystem    = "system"
)

// Chat 聊天会话模型
type Chat struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Model     string         `gorm:"size:100;default:'gpt-3.5-turbo'" json:"model"`
	Status    int            `gorm:"default:1" json:"status"` // 1: 活跃, 0: 已归档
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联字段
	Messages []Message `json:"messages"`
}

// Message 聊天消息模型
type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ChatID    uint           `gorm:"index;not null" json:"chat_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      string         `gorm:"size:20;default:'text'" json:"type"`
	From      string         `gorm:"size:20;not null" json:"from"` // user, assistant, system
	Status    int            `gorm:"default:1" json:"status"` // 1: 已发送, 0: 发送失败
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// 扩展字段
	Metadata map[string]interface{} `gorm:"type:json" json:"metadata"`
}

// TableName 设置表名
func (Chat) TableName() string {
	return "chats"
}

// TableName 设置表名
func (Message) TableName() string {
	return "messages"
}