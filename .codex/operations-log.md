## 2026-04-14 18:40 CST
- Tool: `Get-ChildItem`, `rg`, `Get-Content`
- Purpose: 定位官网、申请页、文案文件与后端入口
- Result: 确认前端主入口在 `litechat/src/views/LandingView.vue` 与 `ApplyBetaView.vue`，文案在 `src/locales/*/landing.ts`
- Decision: 采用前后端同时改造，移除公开产品直达入口

## 2026-04-14 18:48 CST
- Tool: `Get-Content`
- Purpose: 检查后端 HTTP 注册、响应封装、数据层模式
- Result: 确认后端为 Kratos + GORM + Wire，自定义响应封装在 `backend/pkg/auth/encoder.go`
- Decision: 新增匿名 `POST /beta/applications`，使用独立业务实体并在构造仓库时自动迁移

## 2026-04-14 18:55 CST
- Tool: `request_user_input`
- Purpose: 确认表单字段、入口策略和邮件提醒范围
- Result: 采用基础四项、全部移除公开直达入口、提交后入库并发邮件提醒
- Decision: 网页申请字段固定为产品方向、联系方式、使用场景、补充说明
## 2026-04-15 00:15 CST
- Tool: PowerShell `Set-Content`
- Purpose: 写入官网、申请页、locale、footer 与后端申请链路实现
- Result: 完成前后端核心文件重写与新文件新增
- Decision: 因 `apply_patch` 工具在当前环境不可用，改用 UTF-8 的 PowerShell 直接写文件

## 2026-04-15 00:32 CST
- Tool: `pnpm build`
- Purpose: 验证 litechat 前端重构
- Result: 成功；仅保留既有 CSS `quote` 告警和 bundle 提示
- Decision: 记录为非阻塞遗留项

## 2026-04-15 00:39 CST
- Tool: `go build ./cmd/backend`
- Purpose: 验证 backend 申请接口编译
- Result: 首次因 `kerrors.MethodNotAllowed` 不存在失败，替换为 `kerrors.New(405, ...)` 后成功
- Decision: 保持匿名申请接口的自定义 HTTP 处理器方案
## 2026-04-15 10:05 CST
- Tool: `gen-api.ps1`, `protoc`
- Purpose: 将内测申请接口切换为 protobuf，并把邮件配置迁入 YAML 配置树
- Result: 新增 `api/base/beta_application.proto`，生成对应 pb/http/grpc 文件；`conf.proto` 与 `conf.pb.go` 已更新
- Decision: 废弃手写 `HandleFunc` 主路径，改用生成注册函数 `RegisterBetaApplicationHTTPServer`

## 2026-04-15 10:42 CST
- Tool: `Get-Content`, PowerShell `Set-Content`, `pnpm build`
- Purpose: 按 `example_landing_page/Home` 的结构重做官网与申请页前端风格，并预留图片/动图位置
- Result: `LandingView.vue`、`ApplyBetaView.vue`、`LandingFooter.vue` 与中英文 landing 文案已重构为统一的 Ant Landing 风格骨架；首页保留 sticky/pin 内容区并新增多个媒体占位模块
- Decision: 保留现有内测收口逻辑与三种申请方式，只替换视觉语言、版式节奏与素材承载方式

## 2026-04-15 11:36 CST
- Tool: `Get-Content`, `rg`, `git status --short`, `git diff`
- Purpose: 重新核对首页、申请页、locale、页脚与工作区脏文件边界，避免覆盖用户已有改动
- Result: 确认本轮只触达 litechat 公开官网相关文件；后端和其他脏文件保持不动
- Decision: 维持“首页 + 申请页 + 共用公开站组件 + locale”这一最小改动面

## 2026-04-15 11:49 CST
- Tool: `apply_patch`, PowerShell `WriteAllText`
- Purpose: 写入新的公开站头部、首页、申请页、页脚与 locale
- Result: `apply_patch` 在当前环境持续触发 sandbox refresh 失败，因此改用 UTF-8 无 BOM 的 PowerShell 写入目标文件
- Decision: 继续遵循不碰无关文件的原则，只覆盖本轮重做涉及的公开站文件

## 2026-04-15 11:54 CST
- Tool: `pnpm build`
- Purpose: 验证按 `example_landing_page/Home` 节奏重做后的官网公开页
- Result: 构建成功；仅保留既有 CSS `quote` 告警、动态导入提示和 chunk 体积提示
- Decision: 将这些告警记录为现有遗留项，不阻塞本轮交付
## 2026-04-15 12:08 CST
- Tool: `Get-Content`, `rg`, `pnpm build`
- Purpose: 排查并修复 `ApplyBetaView.vue` 表单区域的提交交互问题
- Result: 确认 `handleSubmit` 只绑定在提交按钮点击事件上，`Form` 的 `@submit.prevent` 未接处理器；已改为由表单统一接管提交，并移除按钮上的重复绑定
- Decision: 与仓库内 `AuthModal.vue` 的表单模式保持一致，确保点击按钮和按回车都走同一条提交链路
## 2026-04-15 12:22 CST
- Tool: `Get-Content`, `rg`, `pnpm build`
- Purpose: 继续加固 `ApplyBetaView.vue` 的表单提交流程
- Result: 确认 ant-design-vue `Form` 本地实现支持 `submit` / `finish` 事件；已将页面改为 `:model="form" + @submit="handleSubmit"`，并在处理器中补充重复提交保护
- Decision: 保留现有手写校验与错误展示，只修正事件绑定和并发提交风险，避免扩大改动面