# AI助手项目后端

## 项目介绍

这是一个基于Go语言开发的AI助手后端系统，支持集成MCP服务，兼容第三方提供的接口，通过Agent调用执行，并提供后台配置功能，可配置函数工具和工作流。

## 技术栈

- Go 1.24
- Gin Web框架
- GORM ORM框架
- MySQL数据库
- Redis缓存
- JWT认证

## 目录结构

```
backend/
├── main.go              # 主程序入口
├── go.mod               # Go模块定义
├── go.sum               # 依赖版本锁定
├── config.yaml          # 配置文件
├── controllers/         # 控制器
│   ├── user_controller.go    # 用户控制器
│   ├── chat_controller.go    # 聊天控制器
│   └── admin/                # 管理后台控制器
│       ├── function_tool_controller.go    # 函数工具控制器
│       └── workflow_controller.go         # 工作流控制器
├── models/              # 数据模型
│   ├── user.go               # 用户模型
│   ├── chat.go               # 聊天模型
│   ├── function_tool.go      # 函数工具模型
│   └── workflow.go           # 工作流模型
├── pkg/                 # 公共包
│   ├── db/                   # 数据库操作
│   ├── redis/                # Redis操作
│   ├── auth/                 # 认证功能
│   ├── initdata/             # 初始化数据
│   ├── mcp/                  # MCP客户端
│   └── agent/                # Agent执行引擎
```

## 功能特性

### 用户管理
- 用户注册
- 用户登录（支持用户名/邮箱登录）
- JWT认证
- 获取用户信息

### 聊天功能
- 发送消息
- 获取聊天历史
- Agent调用函数工具和工作流

### 管理后台
- MCP服务配置
- 函数工具管理（增删改查）
- 工作流管理（增删改查）

## 配置说明

在`config.yaml`文件中配置以下内容：

```yaml
server:
  port: 8080  # 服务器端口

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: aichat
  charset: utf8mb4
  parseTime: true
  loc: Local

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key  # JWT密钥
  expire_time: 72          # 过期时间（小时）

log:
  level: info
  path: logs/
  max_size: 100
  max_age: 7
  max_backups: 5

mcp:
  enabled: true
  servers:
    - name: mcp1
      url: http://localhost:8001
      token: ""
    - name: mcp2
      url: http://localhost:8002
      token: ""
```

## 安装依赖

```bash
cd backend
# 初始化模块（如果尚未初始化）
go mod init github.com/aichat/backend
# 安装依赖
go mod tidy
```

## 运行项目

确保MySQL和Redis服务已启动，然后执行：

```bash
cd backend
# 直接运行
go run main.go
# 或编译后运行
go build -o aichat
sudo ./aichat
```

## 默认管理员账户

系统启动时会自动创建默认管理员账户：
- 用户名：admin
- 密码：admin123456
- 请在首次登录后修改默认密码

## API接口文档

### 用户接口
- POST `/api/v1/user/login` - 用户登录
- POST `/api/v1/user/register` - 用户注册
- GET `/api/v1/user/info` - 获取用户信息（需要认证）

### 聊天接口
- POST `/api/v1/chat/send` - 发送消息（需要认证）
- GET `/api/v1/chat/history` - 获取聊天历史（需要认证）

### 管理接口（需要管理员权限）
- GET `/api/v1/admin/function/list` - 获取函数工具列表
- POST `/api/v1/admin/function/create` - 创建函数工具
- PUT `/api/v1/admin/function/:id` - 更新函数工具
- DELETE `/api/v1/admin/function/:id` - 删除函数工具

- GET `/api/v1/admin/workflow/list` - 获取工作流列表
- POST `/api/v1/admin/workflow/create` - 创建工作流
- PUT `/api/v1/admin/workflow/:id` - 更新工作流
- DELETE `/api/v1/admin/workflow/:id` - 删除工作流

## 注意事项

1. 确保在生产环境中修改默认的JWT密钥和管理员密码
2. 根据实际需求调整数据库和Redis配置
3. 定期备份数据库以防止数据丢失
4. 配置适当的日志级别和轮转策略