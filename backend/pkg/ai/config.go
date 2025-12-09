/*
 * @Author: chengjiang
 * @Date: 2025-12-09
 * @Description: AI Agent 配置结构定义（支持数据库存储）
 */
package ai

// AgentConfig 定义 Agent 的配置
type AgentConfig struct {
	Name               string      `json:"name" gorm:"uniqueIndex;size:100"`
	Description        string      `json:"description" gorm:"size:500"`              // agent prompt
	AdapterType        AdapterType `json:"adapter_type" gorm:"size:50;default:eino"` // adk, deepadk 等
	ModelConfig        ModelConfig `json:"model_config" gorm:"foreignKey:AgentID"`
	MaxIteration       int         `json:"max_iteration" gorm:"default:10"`
	IsMaster           bool        `json:"is_master" gorm:"default:false"`
	Status             int         `json:"status" gorm:"default:1"` // 0: 禁用, 1: 启用
	ParentID           int64
	Order              int
	WithWriteTODOs     bool
	WithWebSearchAgent bool
}

// TableName 指定表名
func (AgentConfig) TableName() string {
	return "ai_agent_config"
}

type AdapterType string

const (
	AdapterTypeEino    AdapterType = "adk"
	AdapterTypeDeepAdk AdapterType = "deepadk"
)

// ModelConfig 模型配置
type ModelConfig struct {
	AgentID     int64   `json:"agent_id" gorm:"index"`
	ModelType   string  `json:"model_type" gorm:"size:50"` // ark, openai, deepseek
	ModelName   string  `json:"model_name" gorm:"size:100"`
	APIKey      string  `json:"api_key" gorm:"size:255"` // 建议加密存储
	BaseURL     string  `json:"base_url" gorm:"size:255"`
	MaxTokens   int     `json:"max_tokens" gorm:"default:4096"`
	Temperature float32 `json:"temperature" gorm:"default:0.7"`
	TopP        float32 `json:"top_p" gorm:"default:0.9"`
	Thinking    bool    `json:"thinking" gorm:"default:false"`
	PriceType   string  `json:"price_type" gorm:"size:50"` // 计费方式
}

// TableName 指定表名
func (ModelConfig) TableName() string {
	return "ai_model_config"
}

// PromptConfig 提示词配置
type PromptConfig struct {
	SystemPrompt string            `json:"system_prompt"`
	UserPrompt   string            `json:"user_prompt"`
	Variables    map[string]string `json:"variables"`
}

// ToolConfig 工具配置
type ToolConfig struct {
	ToolType    string                 `json:"tool_type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Params      map[string]interface{} `json:"params"`
}

// AIAgentConfig 完整 AI Agent 配置（用于内存中使用）
type AIAgentConfig struct {
	AgentConfig  AgentConfig     `json:"agent_config"`
	PromptConfig PromptConfig    `json:"prompt_config"`
	Tools        []ToolConfig    `json:"tools"`
	SubAgents    []AIAgentConfig `json:"sub_agents"`
}
