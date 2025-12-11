/* *
 * @Author: chengjiang
 * @Date: 2025-12-10 17:02:59
 * @Description:
**/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/example/aichat/backend/models"
)

type AITool struct {
	models.Model
	Name        string            `json:"name" gorm:"size:100"`
	Code        string            `json:"code" gorm:"size:100;comment:'code'"`
	Description string            `json:"description" gorm:"size:200"`
	SysType     models.SystemType `json:"sys_type" gorm:"type:smallint;default:2"`
	Type        AI_Tool_Type      `json:"type" gorm:"type:smallint;default:1"` // functionCall, MCP
	Status      models.Status     `json:"status" gorm:"type:smallint;default:1;comment:'状态'"`
	Params      AI_Tool_Params    `json:"params" gorm:"type:text;comment:'参数'"`
	MCPUrl      string            `json:"mcp_url" gorm:"size:200;comment:'MCP地址'"`
	MCPToken    string            `json:"mcp_token" gorm:"size:200;comment:'MCP Token'"`
}

func (AITool) TableName() string {
	return "ai_tool"
}

// functionCall / MCP
type AI_Tool_Type uint8

const (
	AI_Tool_Type_FunctionCall AI_Tool_Type = iota + 1 // functionCall
	AI_Tool_Type_MCP                                  // MCP
)

type AI_Tool_Params []*AI_Tool_Param

type AI_Tool_Param struct {
	ParamName string `json:"param_name"`
	ParamType string `json:"param_type"`
	Default   string `json:"default"`
}

// scan value
func (p *AI_Tool_Params) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return fmt.Errorf("cannot scan %T into AI_Tool_Params", value)
	}
}

func (p AI_Tool_Params) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// AIToolAgentBind 工具绑定 Agent
type AIToolAgentBind struct {
	models.Model
	AgentID  uint64 `json:"agent_id" gorm:"type:bigint;not null;comment:'AgentID'"`
	ToolCode string `json:"tool_code" gorm:"type:varchar(100);not null;comment:'工具Code'"`
}

func (AIToolAgentBind) TableName() string {
	return "ai_tool_agent_bind"
}
