# Verification

## 已验证
- 官网首页已重构为更客观、简约的落地页，公开 CTA 统一收口到 `/apply`。
- 官网内容区已改为 sticky/pin 叙事布局，产品方向通过固定容器和滚动切换呈现。
- 首页、申请页和页脚中不再提供 `LiteChat` / `AICook` 的公开直达入口。
- 申请页已提供三种方式：QQ 联系、邮箱联系、网页申请。
- 网页申请 API 已切换为 protobuf 定义，接口来自 `backend/api/base/beta_application.proto`。
- 后端已通过生成代码注册 `BetaApplication` 的 HTTP/GRPC 服务，不再使用手写 `HandleFunc` 作为主入口。
- 邮件配置已迁移到 `backend/internal/conf/conf.proto` 与 `backend/configs/config.yaml.example` 的 `beta.mail` 节点。
- `pnpm build`、`go build ./cmd/backend`、`gen-api.ps1` 和 `conf.proto` 生成命令均已通过。

## 未本地验证
- 未启动浏览器逐项检查桌面端和移动端的真实视觉效果、sticky 动画节奏与滚动手感。
- 未对新申请接口做真实 HTTP 提交联调，也未验证数据库中的自动建表结果。
- 未填入真实 YAML 邮件配置，因此 SMTP 提醒只完成了代码接入，未进行实发验证。

## 残留风险
- `beta.mail` 配置为空时，申请仍会入库，但邮件提醒会降级为跳过发送。
- 前端构建仍提示一个既有 CSS `quote` 属性告警，以及原有的大 bundle / 动态导入提示，本次未处理。
- 工作区中存在用户原有未提交修改（如 `.gitignore`、`README.md` 等），本次实现未回退这些改动。

## 本轮前端补充验证
- 官网首页已改为参考 `example_landing_page/Home` 的分段式落地页结构，包含 Banner、产品摘要、sticky/pin 展示区、理念区与申请收口区。
- 首页与申请页均已预留明确的图片 / 动图占位模块，后续可直接替换为产品截图、流程图或 GIF。
- 申请页已与首页统一视觉语言，仍保留 QQ 联系、邮箱联系、网页申请三种方式。
- 本轮前端重构后的 `pnpm build` 已通过，未引入新的构建错误。

## 2026-04-15 官网重做补充
- 首页与申请页已统一到新的母品牌站视觉，不再保留公开页主题切换。
- 首页结构已改为参考 `example_landing_page/Home` 的分段式节奏：顶部导航、主 Banner、品牌定位、产品区、工作方式、最近在做、申请收口。
- 申请页已保留 QQ、邮箱、网页申请三种入口，并保持 `betaApi.submitApplication` 的原有提交契约。
- 新增的公开站共用头部组件为 `litechat/src/components/landing/PublicSiteHeader.vue`，页脚与语言切换也同步适配新视觉变量。
- `pnpm build` 已通过，本轮未引入新的前端构建错误。
## 2026-04-15 申请页提交修复
- `ApplyBetaView.vue` 已将 `handleSubmit` 绑定到 `Form` 的 `@submit.prevent`，提交按钮仅保留 `html-type="submit"`。
- 申请页表单现在点击提交按钮和在输入框中按回车会共用同一条提交逻辑，不再依赖按钮点击事件单独触发。
- 本轮已再次执行 `litechat` 下的 `pnpm build`，构建通过，未引入新的前端编译错误。
## 2026-04-15 申请页提交二次修复
- `ApplyBetaView.vue` 已将 `Form` 补充 `:model="form"`，并改为监听 `@submit="handleSubmit"`，避免继续依赖 `@submit.prevent` 的组件事件修饰写法。
- `handleSubmit` 已新增事件兜底与 `isSubmitting` 防重入判断，降低按钮连点或重复回车造成并发提交的风险。
- 本轮再次执行 `litechat` 下的 `pnpm build`，构建通过，未引入新的前端编译错误。