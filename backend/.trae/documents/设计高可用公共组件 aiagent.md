# 设计高可用公共组件 aiagent

## 1. 设计目标
- 封装 eino 框架的 adk，避免业务直接调用
- 提供统一的配置结构体，包括 agent、prompt 和 tool 配置
- 支持多种模型类型
- 支持工具调用
- 支持流式输出
- 便于后续数据库存储

## 2. 配置结构体设计

### 2.1 AgentConfig
```go
type AgentConfig struct {
    Name        string            // Agent 名称
    Description string            // Agent 描述
    ModelConfig ModelConfig       // 模型配置
    MaxIteration int              // 最大迭代次数
    ExtraConfig map[string]interface{} // 额外配置
}
```

### 2.2 PromptConfig
```go
type PromptConfig struct {
    SystemPrompt string            // 系统提示词
    UserPrompt   string            // 用户提示词模板
    Variables    map[string]string // 提示词变量
}
```

### 2.3 ToolConfig
```go
type ToolConfig struct {
    ToolType string                 // 工具类型
    Name     string                 // 工具名称
    Description string              // 工具描述
    Params   map[string]interface{} // 工具参数
}
```

### 2.4 AIAgentConfig
```go
type AIAgentConfig struct {
    AgentConfig  AgentConfig        // Agent 配置
    PromptConfig PromptConfig       // Prompt 配置
    Tools        []ToolConfig       // 工具配置
    SubAgents    []AIAgentConfig    // 子 Agent 配置
}
```

## 3. 公共组件实现

### 3.1 核心结构体
```go
type AIAgent struct {
    config    AIAgentConfig
    agent     adk.Agent
    chatModel model.ToolCallingChatModel
}
```

### 3.2 工厂方法
```go
func NewAIAgent(ctx context.Context, config AIAgentConfig) (*AIAgent, error) {
    // 实现逻辑
}
```

### 3.3 流式输出方法
```go
func (a *AIAgent) Stream(ctx context.Context, req interface{}) (<-chan interface{}, error) {
    // 实现逻辑
}
```

### 3.4 非流式输出方法
```go
func (a *AIAgent) Run(ctx context.Context, req interface{}) (interface{}, error) {
    // 实现逻辑
}
```

## 4. 文件结构

- `/pkg/ai/agent.go`：核心 Agent 接口和实现
- `/pkg/ai/config.go`：配置结构体定义
- `/pkg/ai/factory.go`：工厂方法
- `/pkg/ai/tool.go`：工具相关定义
- `/pkg/ai/prompt.go`：prompt 相关定义

## 5. 实现步骤

1. 创建配置结构体文件 `config.go`
2. 创建工具和 prompt 相关文件
3. 创建核心 Agent 实现文件
4. 修改现有 `deep.go` 依赖新的公共组件
5. 提供简化的 API 供业务调用

## 6. 优势

- 业务无需直接调用 eino 框架的 adk，减少第三方框架污染
- 统一的配置管理，便于后续数据库存储
- 支持多种模型和工具，扩展性强
- 高可用性设计，支持错误处理和重试
- 清晰的接口设计，便于业务使用

## 7. 后续扩展

- 支持配置热更新
- 支持 Agent 池化管理
- 支持监控和日志
- 支持分布式部署

这个设计将实现高可用的公共组件 aiagent，避免直接调用 eino 框架的 adk，减少第三方框架污染，同时便于后续数据库存储和业务调用。