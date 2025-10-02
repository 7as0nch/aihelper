# 前端重构设计文档

## 概述
本文档描述了前端重构项目的设计方案，包括架构设计、组件设计、路由设计等方面的内容。

## 技术架构设计

### 整体架构
1. 采用Vue 3 + TypeScript作为主要技术栈
2. 使用Vite作为构建工具
3. 管理后台使用Element Plus组件库
4. H5端使用Vant组件库
5. 采用Pinia进行状态管理
6. 使用Vue Router进行路由管理

### 目录结构设计
```
src/
├── api/           # 接口请求
├── assets/        # 静态资源
├── components/    # 公共组件
├── composables/   # 组合式函数
├── layouts/       # 布局组件
├── pages/         # 页面组件
├── plugins/       # 插件
├── router/        # 路由配置
├── stores/        # 状态管理
├── styles/        # 样式文件
├── utils/         # 工具函数
└── views/         # 视图组件
```

## 组件设计

### 公共组件
1. Header组件 - 顶部导航栏
2. Footer组件 - 底部信息栏
3. Sidebar组件 - 侧边栏菜单（管理后台）
4. Loading组件 - 加载状态指示器
5. Modal组件 - 弹窗组件
6. Toast组件 - 消息提示组件

### 页面组件设计

#### H5端页面
1. 登录页面 (LoginView)
   - 用户名/密码输入
   - 登录按钮
   - 忘记密码链接
   - 注册跳转链接

2. 注册页面 (RegisterView)
   - 用户名输入
   - 邮箱输入
   - 密码输入
   - 确认密码输入
   - 注册按钮

3. 首页 (HomeView)
   - 欢迎信息
   - 功能导航入口

4. 聊天页面 (ChatView)
   - 消息展示区域
   - 输入框
   - 发送按钮
   - 工具按钮（图片、语音等）

5. 个人中心页面 (ProfileView)
   - 用户信息展示
   - 头像修改
   - 退出登录按钮

#### 管理后台页面
1. 管理员登录页面 (AdminLoginView)
   - 管理员账号登录
   - 验证码输入

2. 管理后台布局 (AdminLayout)
   - 侧边栏菜单
   - 顶部导航
   - 内容区域

3. 仪表盘页面 (DashboardView)
   - 数据统计展示
   - 图表展示

4. 用户管理页面 (UserManagementView)
   - 用户列表展示
   - 用户信息编辑
   - 用户状态管理

5. 系统设置页面 (SystemSettingsView)
   - 系统参数配置
   - 权限管理

6. 工具管理页面 (FunctionToolManagementView)
   - 工具列表
   - 工具配置

7. 工作流管理页面 (WorkflowManagementView)
   - 工作流列表
   - 工作流配置

## 路由设计

### H5端路由
- `/` - 首页
- `/login` - 登录页面
- `/register` - 注册页面
- `/chat` - 聊天页面
- `/profile` - 个人中心

### 管理后台路由
- `/admin/login` - 管理员登录
- `/admin` - 管理后台布局
- `/admin/dashboard` - 仪表盘
- `/admin/users` - 用户管理
- `/admin/settings` - 系统设置
- `/admin/tools` - 工具管理
- `/admin/workflows` - 工作流管理

## 样式设计

### 设计规范
1. 颜色规范
   - 主色调：#409EFF（Element Plus主色）
   - 辅助色：#67C23A（成功）、#E6A23C（警告）、#F56C6C（危险）
   - 中性色：#303133（主要文字）、#606266（常规文字）、#909399（次要文字）

2. 字体规范
   - 主要字体：Helvetica Neue, Helvetica, PingFang SC, Hiragino Sans GB, Microsoft YaHei
   - 字号：12px、14px、16px、18px、20px

3. 间距规范
   - 页面边距：20px
   - 元素间距：10px、15px、20px

### 响应式设计
1. 移动端适配
   - 使用vw/vh单位
   - 媒体查询适配不同屏幕尺寸
   - Flexbox布局

2. 管理后台适配
   - 固定宽度布局
   - 响应式表格

## 性能优化

### 加载优化
1. 代码分割
2. 懒加载组件
3. 图片压缩
4. 静态资源CDN

### 交互优化
1. 防抖节流处理
2. 虚拟滚动列表
3. 骨架屏加载
4. 预加载关键资源

## 安全设计

### 认证授权
1. JWT Token认证
2. 路由守卫权限控制
3. 按钮级别权限控制

### 数据安全
1. XSS防护
2. CSRF防护
3. 敏感信息加密存储