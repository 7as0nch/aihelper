# Requirements Document

## Introduction

本项目旨在开发一套高度可定制的智能助手系统，核心功能定位为类客服形态的AI聊天机器人。系统通过集成MCP服务实现与第三方接口的兼容，并提供专业的管理后台支持MCP和工作流配置。系统采用前后端分离架构，前端分为H5端和管理后台，后端采用单体架构。

## Alignment with Product Vision

[由于steering文档不存在，此处暂时留空]

## Requirements

### Requirement 1: MCP服务集成

**User Story:** 作为系统管理员，我希望能集成自主研发的基于PostgreSQL的MCP服务，以便生成用户反馈内容的报表。

#### Acceptance Criteria

1. WHEN 用户请求报表生成 THEN 系统SHALL调用pgmcp服务查询数据库
2. IF 数据库查询成功 THEN 系统SHALL返回查询参数给Excel处理Agent
3. WHEN Excel处理Agent接收到参数 THEN 系统SHALL为用户导出Excel报表

### Requirement 2: 工作流执行引擎

**User Story:** 作为系统管理员，我希望配置AI助手的工作流，以便AI助手能够依据预设工作流执行多轮接口调用。

#### Acceptance Criteria

1. WHEN 配置工作流 THEN 系统SHALL支持定义多步骤的工作流
2. IF 工作流包含数据库查询步骤 THEN 系统SHALL调用pgmcp服务
3. WHEN 工作流包含Excel导出步骤 THEN 系统SHALL调用Excel智能体
4. IF 工作流执行完成 THEN 系统SHALL向用户反馈完整的操作执行结果

### Requirement 3: 管理后台

**User Story:** 作为系统管理员，我希望有一个专业的管理后台，以便为AI助手配置可调用的MCP及工作流。

#### Acceptance Criteria

1. WHEN 访问管理后台 THEN 系统SHALL提供MCP服务配置界面
2. IF 配置MCP服务 THEN 系统SHALL支持PostgreSQL MCP服务的参数设置
3. WHEN 创建工作流 THEN 系统SHALL提供可视化工作流编辑器
4. IF 保存工作流配置 THEN 系统SHALL持久化存储工作流定义

### Requirement 4: H5前端界面

**User Story:** 作为普通用户，我希望通过H5界面与AI助手交互，以便获得智能客服服务。

#### Acceptance Criteria

1. WHEN 用户访问H5界面 THEN 系统SHALL显示聊天界面
2. IF 用户发送消息 THEN 系统SHALL调用AI助手处理并返回响应
3. WHEN 用户请求报表 THEN 系统SHALL执行预设工作流并返回结果
4. IF 工作流执行成功 THEN 系统SHALL提供Excel报表下载链接

### Requirement 5: 前后端分离架构

**User Story:** 作为开发人员，我希望系统采用前后端分离架构，以便实现更好的可维护性和扩展性。

#### Acceptance Criteria

1. WHEN 部署前端 THEN 系统SHALL使用vue3+TypeScript+vite+antv+tailwindcss技术栈
2. IF 部署后端 THEN 系统SHALL使用go1.24+kratos+postgre+redis+rocketmq+socketio技术栈
3. WHEN 前后端通信 THEN 系统SHALL通过API接口进行数据交换
4. IF 需要实时通信 THEN 系统SHALL使用socketio实现

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: 每个文件应具有单一、明确定义的职责
- **Modular Design**: 组件、工具和服务应隔离且可重用
- **Dependency Management**: 最小化模块间的相互依赖
- **Clear Interfaces**: 定义组件和层之间的清晰契约

### Performance
- 系统响应时间应在2秒内
- 支持并发用户数不少于1000
- 数据库查询响应时间应在500毫秒内

### Security
- 实现用户身份认证和授权
- 数据传输使用HTTPS加密
- 防止SQL注入和XSS攻击
- API接口需要身份验证

### Reliability
- 系统可用性达到99.9%
- 实现错误处理和异常恢复机制
- 支持数据备份和恢复

### Usability
- 管理后台界面参照vben框架风格
- H5界面提供友好的用户体验
- 支持响应式设计，适配不同设备
- 提供清晰的操作指引和错误提示