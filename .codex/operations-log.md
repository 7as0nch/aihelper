## 2026-03-24 11:05 CST
- Tool: PowerShell
- Purpose: 检查仓库根目录、litechat 包配置和技能工作流说明
- Params: `Get-ChildItem -Force`、`Get-Content litechat\package.json`、`Get-Content ...workflow-details.md`
- Result: 确认 litechat 使用 pnpm，已有 node_modules；仓库要求保留 `.codex` 过程文件
- Decision: 按 pnpm 生态实现，并补齐 `.codex` 任务记录

## 2026-03-24 11:08 CST
- Tool: PowerShell
- Purpose: 查找现有 i18n、主题和动画参考
- Params: `Get-Content backweb\packages\locales\src\i18n.ts`、`rg -n "IntersectionObserver|requestAnimationFrame|prefers-reduced-motion|Transition" litechat\src`
- Result: litechat 无 i18n；可参考 backweb 的 vue-i18n 结构；动画相关参考较少
- Decision: 在 litechat 内新建轻量 i18n 模块，并自行实现 Landing 的滚动/背景动效

## 2026-03-24 11:12 CST
- Tool: apply_patch fallback
- Purpose: 写入 `.codex` 任务记录文件
- Params: `functions.apply_patch`
- Result: 沙箱刷新失败，无法正常执行补丁工具
- Decision: 本次改为使用 PowerShell 以 UTF-8 无 BOM 写入文件，并继续记录该回退## 2026-03-24 11:22 CST
- Tool: pnpm
- Purpose: 安装官网重构依赖
- Params: `pnpm add vue-i18n three`
- Result: 成功安装 `vue-i18n` 与 `three`
- Decision: 采用轻量动态导入 `three`，避免把 WebGL 逻辑耦合进非官网页面

## 2026-03-24 11:35 CST
- Tool: PowerShell
- Purpose: 写入 litechat i18n、store 和 landing 组件文件
- Params: 多次 `WriteAllText(..., UTF8 No BOM)`
- Result: 已创建 `src/i18n`、`src/locales`、`src/stores/locale.ts` 和 `src/components/landing/*`
- Decision: 官网重构采用组件化拆分，同时保留现有 App 壳和路由结构

## 2026-03-24 11:43 CST
- Tool: pnpm build
- Purpose: 验证官网重构后的类型和构建
- Params: `pnpm build`
- Result: 初次失败，缺少 `three` 类型声明
- Decision: 安装 `@types/three` 后重新构建

## 2026-03-24 11:46 CST
- Tool: pnpm
- Purpose: 补充 three 类型声明
- Params: `pnpm add -D @types/three`
- Result: 成功安装 `@types/three`
- Decision: 保持 AmbientCanvas 的 TypeScript 类型校验开启

## 2026-03-24 11:49 CST
- Tool: pnpm build
- Purpose: 最终验证 litechat 构建
- Params: `pnpm build`
- Result: 构建通过；存在既有 chunk 体积告警，以及一个 CSS `quote` 属性告警
- Decision: 记录为当前残留风险，不阻塞本次官网交付
