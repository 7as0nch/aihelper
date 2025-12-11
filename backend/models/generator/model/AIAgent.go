/* *
 * @Author: chengjiang
 * @Date: 2025-12-10 17:01:22
 * @Description: AI Agent 核心
**/
package model

import "github.com/example/aichat/backend/models"

type AIAgent struct {
	models.Model
	Name               string            `json:"name" gorm:"uniqueIndex;size:100"`                     // CnName
	Code               string            `json:"code" gorm:"uniqueIndex;size:100"`                     // EnName
	Description        string            `json:"description" gorm:"size:500"`                          // agent prompt
	AdapterType        AdapterType       `json:"adapter_type" gorm:"size:50;default:eino"`             // adk, deepadk 等
	AIModelID          int64             `json:"ai_model_id" gorm:"type:bigint;not null"`              // 关联的 AI 模型 ID
	MaxIteration       int               `json:"max_iteration" gorm:"default:10;comment:'最大迭代次数'"`     // 最大迭代次数
	SystemPrompt       string            `json:"system_prompt" gorm:"type:text;comment:'系统提示词'"`       // 系统提示词
	UserInputPrompt    string            `json:"user_input_prompt" gorm:"type:text;comment:'用户提示词'"`   // 用户提示词
	Status             models.Status     `json:"status" gorm:"default:1;comment:'状态'"`                 // 0: 禁用, 1: 启用
	Type               AgentType         `json:"type" gorm:"default:1;comment:'1 根 Agent, 2 子 Agent'"` // 1 根 Agent, 2 子 Agent: 注意避免相互绑定。
	Order              int               `json:"order" gorm:"default:0"`                               // 排序
	WithWriteTODOs     bool              `json:"with_write_todos" gorm:"default:false"`                // 是否启用写入 TODO 功能
	WithWebSearchAgent bool              `json:"with_web_search_agent" gorm:"default:false"`           // 是否启用 Web 搜索功能
	SystemType         models.SystemType `json:"system_type" gorm:"default:2"`                         // 系统类型: 1-系统内置(不可删除), 2-用户自定义
}

// TableName 指定表名
func (AIAgent) TableName() string {
	return "ai_agent"
}

type AdapterType uint8

const (
	_                   AdapterType = iota
	AdapterType_ADK                 // 普通 ADK 适配器
	AdapterType_DeepADK             // 深度 ADK 适配器
)

type AgentType uint8

const (
	_              AgentType = iota
	AgentType_Root           // 根 Agent
	AgentType_Sub            // 子 Agent
)

type AgentBind struct {
	models.Model
	AgentID int64 `json:"agent_id" gorm:"type:bigint;not null"` // 关联的 AI Agent ID
	SubAgentID int64 `json:"sub_agent_id" gorm:"type:bigint;not null"` // 关联的子 Agent ID
}