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
	AIModel          *UseAIModel      `gorm:"column:ai_model;type:text" json:"ai_model"`
	QuoteId          int64            `gorm:"column:quote_id;type:bigint" json:"quote_id"`
	QuoteContent     string           `gorm:"column:quote_content;type:text" json:"quote_content"`
	QuoteSearchLinks QuoteSearchLinks `gorm:"column:quote_search_links;type:text" json:"quote_search_links"`
	TokenUsage       *TokenUsage      `gorm:"column:token_usage;type:text" json:"token_usage"`
	CallingTools     CallingTools     `gorm:"column:calling_tools;type:text" json:"calling_tools"`
	Attachments      Attachments      `gorm:"column:attachments;type:text" json:"attachments"`
	GenerateTime     string           `gorm:"column:generate_time;type:varchar(50);comment:'生成时间'" json:"generate_time"`
	LikedStatus      LikedStatus      `gorm:"column:liked_status;type:smallint;default:0;comment:'点赞状态'" json:"liked_status"`
}

// TableName AIChatMessage's table name
func (*AIChatMessage) TableName() string {
	return TableNameAIChatMessage
}

type LikedStatus uint8

const (
	_ LikedStatus = iota // 0
	LikedStatus_Liked
	LikedStatus_Unliked
)

// 'user' | 'assistant' | 'human'
type RoleType string

const (
	RoleUser      RoleType = "user"
	RoleAssistant RoleType = "assistant"
	RoleHuman     RoleType = "human"
)

// Custom Types

type AIModel_ThinkingMode string

const (
	AIModel_ThinkingModeSmart AIModel_ThinkingMode = "smart"
	AIModel_ThinkingModeDeep  AIModel_ThinkingMode = "deep"
	AIModel_ThinkingModeQuick AIModel_ThinkingMode = "quick"
)

type AIModel_SearchByWeb uint8

const (
	_                       AIModel_SearchByWeb = iota //smart
	AIModel_SearchByWeb_Yes                            // 搜索网络
	AIModel_SearchByWeb_No                             // 不搜索网络
)

func (m AIModel_SearchByWeb) String() string {
	switch m {
	case AIModel_SearchByWeb_Yes:
		return "yes"
	case AIModel_SearchByWeb_No:
		return "no"
	default:
		return "smart"
	}
}
func AIModel_SearchByWeb_Bool(val bool) AIModel_SearchByWeb {
	if val {
		return AIModel_SearchByWeb_Yes
	}
	return AIModel_SearchByWeb_No
}

type UseAIModel struct {
	ID           string               `json:"id"`
	ModelName    string               `json:"modelName"`
	ThinkingMode AIModel_ThinkingMode `json:"thinkingMode"` // 'smart' | 'deep' | 'quick'
	SearchByWeb  AIModel_SearchByWeb  `json:"searchByWeb"`
}

func (m *UseAIModel) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m UseAIModel) Value() (driver.Value, error) {
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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
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
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
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
