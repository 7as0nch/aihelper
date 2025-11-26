# LiteChat - 可嵌入式 AI 聊天组件

一个轻量级、可嵌入的 AI 聊天组件，支持拖拽、全屏、响应式布局，使用 Shadow DOM + iframe 实现完全的样式隔离。

## 写在前面
> 目前支持demo模式和前端模式，后端模式正在开发中。需要集成后台管理端（配置管理、工作流、用户管理、会话管理、知识库管理等）。
> 后续考虑接入客服功能。

## 展示

### 首页
![alt text](/images/main_page.png)

### 聊天
![alt text](/images/chat_page.png)

## ✨ 特性

- **🎨 零样式冲突** - Shadow DOM + iframe 双层隔离，互不影响
- **📱 完美响应式** - 自适应 PC / 移动端布局
- **🎯 拖拽浮窗** - 可自由拖动的浮动按钮和聊天窗口
- **⚙️ 运行时配置** - 无需重新构建即可动态修改配置
- **🔒 TypeScript** - 提供完整的类型定义支持
- **📦 开箱即用** - 无需引入额外 CSS，简单集成

## 📦 安装

```bash
npm install @7as0nch/litechat
# 或
yarn add @7as0nch/litechat
# 或
pnpm add @7as0nch/litechat
```

## 🚀 快速开始

### 在 Vue 3 项目中使用

推荐使用我们提供的 `useAiChat` 组合式函数：

```typescript
import { useAiChat } from '@7as0nch/litechat/vue';

// 在组件中
const { open, close, show, hide, toggle } = useAiChat({
  config: {
    VITE_APP_TITLE: '我的 AI 助手',
  },
  defaultOpen: false,
  containerId: 'my-chat-container' // 可选
});
```



### 在普通 HTML 中使用 (ESM)

```html
<script type="module">
  import { initAiChat } from 'https://unpkg.com/@7as0nch/litechat';
  
  const widget = initAiChat({
    config: { VITE_APP_TITLE: '客服助手' },
    containerId: 'chat-root',
    defaultOpen: true
  });
  
  // 控制组件
  // widget.open();
  // widget.hide();
</script>
```

### 使用 CDN (ESM)

```html
<script type="module">
  import { initAiChat } from 'https://unpkg.com/@7as0nch/litechat/dist-widget/litechat-widget.es.js';
  
  initAiChat({
    config: { VITE_APP_TITLE: 'My Chat' }
  });
</script>
```

## ⚙️ 配置选项

```typescript
interface InitOptions {
  // 运行时配置
  config?: {
    VITE_APP_TITLE?: string;         // 应用标题
    VITE_APP_LOGO?: string;           // Logo URL
    VITE_API_BASE_URL?: string;       // API 基础 URL
    // ... 其他配置
  };
  
  // 是否默认打开
  defaultOpen?: boolean;  // 默认 false

  // 是否默认显示 (浮球或窗口)
  defaultShow?: boolean;  // 默认 true
  
  // 挂载容器 ID
  containerId?: string;   // 默认 'ai-chat-widget-root'
}
```

## 🎮 程序化控制

`initAiChat` 返回 (或通过 Hook 暴露) 的实例包含以下方法：

- `open()`: 打开聊天窗口
- `close()`: 关闭聊天窗口
- `toggle()`: 切换打开/关闭状态
- `show()`: 显示组件 (浮球或窗口)
- `hide()`: 隐藏组件 (完全不可见)
- `unmount()`: 销毁组件并从 DOM 移除

## 🎨 样式隔离

组件使用 **Shadow DOM + iframe** 双层隔离技术：

- **Shadow DOM**: 隔离浮动按钮和窗口框架样式
- **iframe**: 隔离聊天内容样式，确保响应式布局正确

完全不会影响宿主页面样式，也不会被宿主页面样式影响。

## 📱 响应式设计

组件窗口宽度为 360px 时自动触发移动端布局：

- 隐藏侧边栏
- 启用移动端工具栏
- 优化触摸交互体验

## 🔧 开发

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建库
npm run build:lib

# 预览 demo
open demo-widget.html
open demo-npm.html
```

## 📄 License

MIT

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

- **邮箱**：7as0nch@gmail.com
- **WeChat**：JasonC12o9
- **QQ**：2538684421