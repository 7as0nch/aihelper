package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/example/aichat/backend/models"
)

const TableNameAIChatMessage = "ai_chat_message"

// AIChatMessage mapped from table <ai_chat_message>
type AIChatMessage struct {
	models.Model
	SessionID        int64            `gorm:"column:session_id;type:bigint;not null;index" json:"session_id"`
	Role             RoleType         `gorm:"column:role;type:varchar(50);not null" json:"role"`
	Content          string           `gorm:"column:content;type:text" json:"content"`
	ReasoningContent string           `gorm:"column:reasoning_content;type:text" json:"reasoning_content"`
	AIModel          *AIModel         `gorm:"column:ai_model;type:text" json:"ai_model"`
	QuoteId          string           `gorm:"column:quote_id;type:varchar(255)" json:"quote_id"`
	QuoteContent     string           `gorm:"column:quote_content;type:text" json:"quote_content"`
	QuoteSearchLinks QuoteSearchLinks `gorm:"column:quote_search_links;type:text" json:"quote_search_links"`
	TokenUsage       *TokenUsage      `gorm:"column:token_usage;type:text" json:"token_usage"`
	CallingTools     CallingTools     `gorm:"column:calling_tools;type:text" json:"calling_tools"`
	Attachments      Attachments      `gorm:"column:attachments;type:text" json:"attachments"`
	IsStreaming      bool             `gorm:"column:is_streaming;type:boolean;default:false" json:"is_streaming"`
}

// TableName AIChatMessage's table name
func (*AIChatMessage) TableName() string {
	return TableNameAIChatMessage
}

// 'user' | 'assistant' | 'human'
type RoleType string

const (
	RoleUser      RoleType = "user"
	RoleAssistant RoleType = "assistant"
	RoleHuman     RoleType = "human"
)

// Custom Types

type AIModel struct {
	ID           string `json:"id"`
	ModelName    string `json:"modelName"`
	ThinkingMode string `json:"thinkingMode"` // 'smart' | 'deep' | 'quick'
}

func (m *AIModel) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m AIModel) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type QuoteSearchLink struct {
	Url       string   `json:"url"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Highlight []string `json:"highlight"`
}

type QuoteSearchLinks []*QuoteSearchLink

func (m *QuoteSearchLinks) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m QuoteSearchLinks) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

type TokenUsage struct {
	CurrentTokens int64 `json:"currentTokens"`
	TotalTokens   int64 `json:"totalTokens"`
}

func (m *TokenUsage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m TokenUsage) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type CallingTool struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	FunctionName string `json:"functionName"`
}

type CallingTools []*CallingTool

func (m *CallingTools) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m CallingTools) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

type Attachment struct {
	ID   string `json:"id"`
	Type string `json:"type"` // 'image' | 'file'
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Attachments []*Attachment

func (m *Attachments) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m Attachments) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
