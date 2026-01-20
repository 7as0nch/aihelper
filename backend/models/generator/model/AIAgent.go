/* *
 * @Author: chengjiang
 * @Date: 2025-12-10 17:01:22
 * @Description: AI Agent 核心
**/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/example/aichat/backend/models"
)

type AIAgent struct {
	models.Model
	Name               string            `json:"name" gorm:"uniqueIndex;size:100"`                     // CnName
	Code               string            `json:"code" gorm:"uniqueIndex;size:100"`                     // EnName
	Description        string            `json:"description" gorm:"size:500"`                          // agent prompt
	AdapterType        AdapterType       `json:"adapter_type" gorm:"size:50;default:2"`                // adk, deepadk 等
	AIModelID          int64             `json:"ai_model_id" gorm:"type:bigint;not null"`              // 关联的 AI 模型 ID: 原始model ID：追溯引用的模型配置。
	AIModel            *AIAgentModel     `json:"ai_model" gorm:"type:text;comment:'模型配置'"`             // model副本，后续通过该获取model配置。
	MaxIteration       int               `json:"max_iteration" gorm:"default:10;comment:'最大迭代次数'"`     // 最大迭代次数
	SystemPrompt       string            `json:"system_prompt" gorm:"type:text;comment:'系统提示词'"`       // 系统提示词
	UserInputPrompt    string            `json:"user_input_prompt" gorm:"type:text;comment:'用户提示词'"`   // 用户提示词
	Status             models.Status     `json:"status" gorm:"default:1;comment:'状态'"`                 // 0: 禁用, 1: 启用
	Type               AgentType         `json:"type" gorm:"default:1;comment:'1 根 Agent, 2 子 Agent'"` // 1 根 Agent, 2 子 Agent: 注意避免相互绑定。
	Order              int               `json:"order" gorm:"default:0"`                               // 排序
	WithWriteTODOs     bool              `json:"with_write_todos" gorm:"default:false"`                // 是否启用写入 TODO 功能
	WithWebSearchAgent bool              `json:"with_web_search_agent" gorm:"default:false"`           // 是否启用 Web 搜索功能
	SystemType         models.SystemType `json:"system_type" gorm:"default:2"`                         // 系统类型: 1-系统内置(不可删除), 2-用户自定义
	ParentID           int64             `json:"parent_id" gorm:"-"`                                   // 父 Agent ID
	SubAIAgents        []*AIAgent        `json:"sub_agents" gorm:"-"`                                  // 子 Agent 列表
}

// TableName 指定表名
func (AIAgent) TableName() string {
	return "ai_agent"
}

type AdapterType uint8

const (
	_                    AdapterType = iota
	AdapterType_ADK                  // 普通 ADK 适配器
	AdapterType_DeepADK              // 深度 ADK 适配器
	AdapterType_Workflow             // 工作流适配器
)

type AgentType uint8

const (
	_              AgentType = iota
	AgentType_Root           // 根 Agent
	AgentType_Sub            // 子 Agent
)

type AIAgentModel struct {
	Category    AIModel_Category `json:"category" gorm:"type:smallint;default:0"` // 模型类别
	ModelType   ModelType         `json:"model_type" gorm:"type:varchar(50);not null"`
	ModelName   string           `json:"model_name" gorm:"type:varchar(100);not null"`
	APIKey      string           `json:"api_key" gorm:"type:varchar(255)"` // 建议加密存储
	BaseURL     string           `json:"base_url" gorm:"type:varchar(255)"`
	Temperature float32          `json:"temperature" gorm:"type:float;default:0.7"`
	TopP        float32          `json:"top_p" gorm:"type:float;default:0.9"`
}

func (m *AIAgentModel) Scan(value interface{}) error {
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
	if len(bytes) == 0 || string(bytes) == "null" {
		if m != nil {
			*m = AIAgentModel{}
		}
		return nil
	}
	if m == nil {
		// 如果接收者是 nil，无法赋值，这不应该发生
		// 但为了安全，我们返回错误
		return errors.New("cannot scan into nil *AIModel")
	}
	return json.Unmarshal(bytes, m)
}

// Value 实现 driver.Valuer 接口，用于将数据写入数据库
func (m AIAgentModel) Value() (driver.Value, error) {
	return json.Marshal(m) // 直接序列化整个结构体
}

type AIApplication struct { // AI Application (原 AIAgentProgram 升级)
	models.Model
	Name        string        `json:"name" gorm:"type:varchar(100);comment:'名称'"`             // 名称
	Code        string        `json:"code" gorm:"uniqueIndex;type:varchar(100);comment:'编码'"` // 编码
	Description string        `json:"description" gorm:"size:500;comment:'描述'"`
	Version     string        `json:"version" gorm:"size:20;default:'1.0.0';comment:'版本号'"`
	Mode        ProgramMode   `json:"mode" gorm:"default:2;comment:'模式'"`              // 模式：1.单agent模式（只有一个agent，采用Ark原理，支持adk或者普通enio）2. 多agent模式（就是deepadk或者adk）
	Status      models.Status `json:"status" gorm:"default:1;comment:'状态'"`            // 2: 禁用, 1: 启用
	Type        ProgramType   `json:"type" gorm:"default:1;comment:'程序类型'"`            // 1. 预定义，2. 自定义
	Scope       Scope         `json:"scope" gorm:"default:1;comment:'作用粒度'"`           // 作用粒度：1.所有人，2.指定角色，3.指定用户
	Schema      string        `json:"schema" gorm:"type:text;comment:'schema'"`          // 拖拽页面布局存档。
	SelfAgent   *AIAgent      `json:"self_agent" gorm:"type:json;comment:'自定义 Agent'"` // 自定义 Agent(通过agent列表获取，重新自定义)
}

func (AIApplication) TableName() string {
	return "ai_application"
}

// 程序类型：1. 预定义，2. 自定义
type ProgramType uint8

const (
	ProgramType_Predefined ProgramType = iota + 1 // 预定义程序：agent配置不可更改。
	ProgramType_Custom                            // 自定义程序: 可以使用SelfAgent，存入。后续无需通过agentbind获取对应的agent信息了。
)

// 作用粒度：1.所有人，2.指定角色，3.指定用户
type Scope uint8

const (
	Scope_All  Scope = iota + 1 // 所有人
	Scope_Role                  // 指定角色
	Scope_User                  // 指定用户
)

// Program 模式：1.单agent模式（只有一个agent，采用Ark原理，支持adk或者普通enio）2. 多agent模式（就是deepadk或者adk）
type ProgramMode uint8

const (
	ProgramMode_Single ProgramMode = iota + 1 // 单agent模式
	ProgramMode_Multi                         // 多agent模式
)

// scan value
func (m *AIAgent) Scan(value interface{}) error {
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

func (a AIAgent) Value() (driver.Value, error) {
	return json.Marshal(a) // 直接序列化整个结构体
}
