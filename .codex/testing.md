# Testing

## 2026-03-24
- Command: `pnpm build`
- Workdir: `litechat`
- First result: 失败，缺少 `three` 的类型声明
- Fix: 执行 `pnpm add -D @types/three`
- Second result: 成功

## Build Notes
- Vite 生产构建通过。
- 构建输出包含独立的 `three.module` chunk，说明官网 WebGL 已被拆分为单独资源。
- 构建阶段仍有 CSS `quote` 属性告警和 chunk 体积告警，但未阻塞产物生成。