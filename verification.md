# Verification

## 已验证
- 官网已接入 `vue-i18n`，提供中文与英文两套 Landing 文案。
- `mount.ts` 已集成 i18n 与 locale store 初始化。
- Landing 不再使用本地 `isDark`，而是复用全局 theme store。
- 官网已拆分为多个 landing 组件，并新增语言切换、主题切换、滚动章节状态和 Project Lab 板块。
- 已新增轻量 WebGL 背景组件，并对 `prefers-reduced-motion` 做了降级处理。
- `pnpm build` 已通过，说明类型检查和生产构建链路可用。

## 未本地验证
- 未启动浏览器逐项录屏检查桌面端和移动端的滚动观感。
- 未对真实低性能设备验证 WebGL 帧率表现。

## 残留风险
- 当前构建仍提示 chunk 体积偏大，其中 `three.module` 约 725 kB，需要后续继续关注官网首屏资源策略。
- 构建存在一个 CSS `quote` 属性告警，该问题未定位为本次 Landing 改动引入，但建议后续统一排查。
- 仓库当前存在用户已有的其它改动，例如 `backend/go.mod`，本次未处理。