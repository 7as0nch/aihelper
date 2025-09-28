package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"gorm.io/gorm"
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/db"
	"github.com/aichat/backend/pkg/mcp"
)

// ExecutionContext 执行上下文
type ExecutionContext struct {
	Context  context.Context
	UserID   uint
	Variables map[string]interface{}
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	Success   bool
	Output    interface{}
	Variables map[string]interface{}
	Error     string
}

// Engine Agent执行引擎
type Engine struct {
	mcpClient *mcp.MCPClient
}

// NewEngine 创建新的Agent执行引擎
func NewEngine() *Engine {
	return &Engine{
		mcpClient: mcp.NewMCPClient(),
	}
}

// ExecuteFunction 执行单个函数工具
func (e *Engine) ExecuteFunction(ctx *ExecutionContext, functionToolID uint) *ExecutionResult {
	// 查询函数工具配置
	var functionTool models.FunctionTool
	result := db.GetDB().First(&functionTool, functionToolID)
	if result.Error != nil {
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to find function tool: %v", result.Error),
		}
	}

	if !functionTool.Enabled {
		return &ExecutionResult{
			Success: false,
			Error:   "Function tool is disabled",
		}
	}

	// 解析配置
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(functionTool.Config), &config); err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse function tool config: %v", err),
		}
	}

	// 根据函数类型执行
	var executionResult *ExecutionResult
	switch functionTool.Type {
	case models.FunctionToolTypeMCP:
		executionResult = e.executeMCPFunction(ctx, &functionTool, config)
	case models.FunctionToolTypeHTTP:
		executionResult = e.executeHTTPFunction(ctx, &functionTool, config)
	case models.FunctionToolTypeLocal:
		executionResult = e.executeLocalFunction(ctx, &functionTool, config)
	default:
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Unsupported function tool type: %s", functionTool.Type),
		}
	}

	return executionResult
}

// ExecuteWorkflow 执行工作流
func (e *Engine) ExecuteWorkflow(ctx *ExecutionContext, workflowID uint) *ExecutionResult {
	// 查询工作流配置
	var workflow models.Workflow
	result := db.GetDB().Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("order asc")
	}).First(&workflow, workflowID)

	if result.Error != nil {
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to find workflow: %v", result.Error),
		}
	}

	if !workflow.Enabled {
		return &ExecutionResult{
			Success: false,
			Error:   "Workflow is disabled",
		}
	}

	// 执行工作流步骤
	variables := make(map[string]interface{})
	for k, v := range ctx.Variables {
		variables[k] = v
	}

	var finalResult *ExecutionResult

	for _, step := range workflow.Steps {
		// 如果有条件，检查条件是否满足
		if step.Condition != "" {
			if !e.evaluateCondition(step.Condition, variables) {
				log.Printf("Skipping step %s due to condition not met", step.StepName)
				continue
			}
		}

		// 执行步骤
		executionCtx := &ExecutionContext{
			Context:   ctx.Context,
			UserID:    ctx.UserID,
			Variables: variables,
		}

		result := e.ExecuteFunction(executionCtx, step.FunctionToolID)
		if !result.Success {
			return &ExecutionResult{
				Success: false,
				Error:   fmt.Sprintf("Step %s failed: %s", step.StepName, result.Error),
			}
		}

		// 应用变量映射
		if step.VariableMapping != "" {
			e.applyVariableMapping(step.VariableMapping, result.Output, variables)
		}

		// 更新变量
		for k, v := range result.Variables {
			variables[k] = v
		}

		finalResult = result
	}

	if finalResult == nil {
		return &ExecutionResult{
			Success:   true,
			Variables: variables,
			Output:    "Workflow executed successfully",
		}
	}

	// 合并变量
	for k, v := range variables {
		finalResult.Variables[k] = v
	}

	return finalResult
}

// 执行MCP函数
func (e *Engine) executeMCPFunction(ctx *ExecutionContext, functionTool *models.FunctionTool, config map[string]interface{}) *ExecutionResult {
	// 从配置中获取MCP服务器名称和工具名称
	serverName, ok := config["server_name"].(string)
	if !ok || serverName == "" {
		return &ExecutionResult{
			Success: false,
			Error:   "MCP server name is required",
		}
	}

	toolName, ok := config["tool_name"].(string)
	if !ok || toolName == "" {
		return &ExecutionResult{
			Success: false,
			Error:   "MCP tool name is required",
		}
	}

	// 准备参数
	params := make(map[string]interface{})
	for k, v := range config {
		if k != "server_name" && k != "tool_name" {
			params[k] = v
		}
	}

	// 注入上下文变量
	for k, v := range ctx.Variables {
		params[k] = v
	}

	// 调用MCP服务
	resp, err := e.mcpClient.Call(serverName, toolName, params)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   fmt.Sprintf("MCP call failed: %v", err),
		}
	}

	// 解析响应
	var output interface{}
	if err := json.Unmarshal(resp, &output); err != nil {
		// 如果解析失败，返回原始响应
		return &ExecutionResult{
			Success:   true,
			Output:    string(resp),
			Variables: ctx.Variables,
		}
	}

	return &ExecutionResult{
		Success:   true,
		Output:    output,
		Variables: ctx.Variables,
	}
}

// 执行HTTP函数（占位函数）
func (e *Engine) executeHTTPFunction(ctx *ExecutionContext, functionTool *models.FunctionTool, config map[string]interface{}) *ExecutionResult {
	// TODO: 实现HTTP函数调用逻辑
	return &ExecutionResult{
		Success:   true,
		Output:    "HTTP function executed successfully",
		Variables: ctx.Variables,
	}
}

// 执行本地函数（占位函数）
func (e *Engine) executeLocalFunction(ctx *ExecutionContext, functionTool *models.FunctionTool, config map[string]interface{}) *ExecutionResult {
	// TODO: 实现本地函数调用逻辑
	return &ExecutionResult{
		Success:   true,
		Output:    "Local function executed successfully",
		Variables: ctx.Variables,
	}
}

// 评估条件（占位函数）
func (e *Engine) evaluateCondition(condition string, variables map[string]interface{}) bool {
	// TODO: 实现条件评估逻辑
	return true
}

// 应用变量映射（占位函数）
func (e *Engine) applyVariableMapping(mapping string, output interface{}, variables map[string]interface{}) {
	// TODO: 实现变量映射逻辑
}