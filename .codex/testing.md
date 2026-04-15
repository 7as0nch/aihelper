# Testing

## 2026-04-15
- Command: `pnpm build`
- Workdir: `litechat`
- Result: 成功
- Scope: 覆盖官网首页重构、sticky 内容区、申请页三种入口、locale 重写与新的 `betaApi` 提交链路

## 2026-04-15
- Command: `go build ./cmd/backend`
- Workdir: `backend`
- Result: 成功
- Scope: 覆盖 protobuf 化的 `BetaApplication` 服务、YAML 邮件配置、数据落库、SMTP 通知、Wire 与 HTTP/GRPC 注册

## 2026-04-15
- Command: `./scripts/gen-api.ps1 -Target api/base/beta_application.proto`
- Workdir: `backend`
- Result: 成功
- Scope: 生成 `beta_application.pb.go`、`beta_application_http.pb.go`、`beta_application_grpc.pb.go`

## 2026-04-15
- Command: `protoc --proto_path=./internal --proto_path=./third_party --go_out=paths=source_relative:./internal ./internal/conf/conf.proto`
- Workdir: `backend`
- Result: 成功
- Scope: 生成新的 `conf.pb.go`，使 `beta.mail` YAML 配置可在代码中读取

## Build Notes
- 前端 `vue-tsc -b && vite build` 已通过，说明官网与申请页的新结构、样式和 TypeScript 改动可以完成生产构建。
- 后端 `go build ./cmd/backend` 已通过，说明新的 protobuf 申请接口、配置读取与依赖注入链路编译正常。
- 前端构建仍提示一个既有 CSS `quote` 属性告警，以及原有的大 bundle / 动态导入提示；这些未阻塞本次交付。

## 2026-04-15
- Command: `pnpm build`
- Workdir: `litechat`
- Result: 成功
- Scope: 覆盖基于 `example_landing_page/Home` 风格的官网与申请页重构、媒体占位模块、页脚样式调整与中英文文案同步

## 2026-04-15
- Command: `pnpm build`
- Workdir: `litechat`
- Result: 成功
- Scope: 覆盖新的公开站头部、首页分段落地页、申请页重构、locale 改写与公共页容器背景调整
## 2026-04-15
- Command: `pnpm build`
- Workdir: `litechat`
- Result: 成功
- Scope: 覆盖 `ApplyBetaView.vue` 表单提交修复，确认 `@submit.prevent="handleSubmit"` 的模板改动未引入新的 TypeScript 或 Vite 构建错误
## 2026-04-15
- Command: `pnpm build`
- Workdir: `litechat`
- Result: 成功
- Scope: 覆盖 `ApplyBetaView.vue` 二次修复，确认 `Form` 改为 `:model="form" + @submit="handleSubmit"` 且新增防重复提交后，前端构建仍然通过