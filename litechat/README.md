# LiteChat - 可嵌入式 AI 聊天组件

一个轻量级、可嵌入的 AI 聊天组件，支持拖拽，使用 Shadow DOM + iframe 实现完全的样式隔离。

## 写在前面
> 目前支持demo模式和前端模式，后端模式正在开发中。需要集成后台管理端（配置管理、工作流、用户管理、会话管理、知识库管理、支持插件回调callback等）。
> 后续考虑接入客服功能。\
> // # backend : 纯前端模式：无需登录，只需要接入模型key即可，记录存到本地。\
> // # frontend : 后台模式： 需要登录，接入配套的后台管理。（这里可考虑与宿主完成单点登录配置）\
> // # demo : 样例模式：里面包含各种demo例子。全部使用mock数据。

## 展示

### 首页
![alt text](/images/main_page.png)

### 聊天
![alt text](/images/chat_page.png)

## 特性

- ** 零样式冲突** - Shadow DOM + iframe 双层隔离，互不影响
- ** 完美响应式** - 自适应 PC / 移动端布局
- ** 拖拽浮窗** - 可自由拖动的浮动按钮和聊天窗口
- ** 运行时配置** - 无需重新构建即可动态修改配置
- ** TypeScript** - 提供完整的类型定义支持
- ** 开箱即用** - 无需引入额外 CSS，简单集成

## 安装

```bash
npm install @7as0nch/litechat
# 或
yarn add @7as0nch/litechat
# 或
pnpm add @7as0nch/litechat
```

## 快速开始

### 在 Vue 3 项目中使用

推荐使用我们提供的 `useAiChat` 组合式函数：

```typescript
import { initAiChat } from '@7as0nch/litechat';

// 在组件中
const AiChat = initAiChat({ // InitOptions
  config: { // RuntimeConfig
    VITE_APP_TITLE: '我的 AI 助手',
  },
  defaultOpen: false,
  defaultShow: true,  // 控制浮球是否展示
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

## 配置选项

### 完整配置对象 (LiteChat SDK)

```typescript
interface InitOptions {
  // 1. 基础配置 (Legacy)
  config?: Partial<RuntimeConfig>; // 见下文 RuntimeConfig
  
  // 2. 聊天配置
  chat?: {
    appId?: string;      // Bot ID (对应 VITE_OPENAI_MODEL)
    repoId?: string;
  };
  
  // 3. 系统设置
  setting?: {
    apiBaseUrl?: string; // API 地址 (对应 VITE_API_BASE_URL)
  };
  
  // 4. 认证信息
  auth?: {
    type?: 'token' | 'external';
    token?: string;      // API Key / Token
    onRefreshToken?: (oldToken?: string) => string | Promise<string>;
  };

  // 5. 用户信息 (可选)
  user?: {
    id: string;
    name: string;
    avatar: string;
  };

  // 6. UI 配置
  ui?: {
    layout?: 'pc' | 'mobile'; // 默认 pc
    header?: {
      isNeed?: boolean;       // 是否显示头部
      title?: string;         // 标题
      icon?: string;          // 图标 URL
    };
    footer?: {
      isNeed?: boolean;       // 是否显示底部
      expressionText?: string;// 底部免责声明文本
    };
    chatSlot?: {
      input?: {
        placeholder?: string; // 输入框占位符
        isNeedAudio?: boolean;// 是否显示语音按钮
      };
    };
    uploadBtn?: {
      isNeed?: boolean;       // 是否显示上传按钮
    };
  };

  // 通用选项
  defaultOpen?: boolean;  // 默认打开?
  defaultShow?: boolean;  // 默认显示?
  containerId?: string;   // 挂载容器 ID
}
```

### 环境变量 / RuntimeConfig

```typescript
interface RuntimeConfig {
    // 基础
    VITE_APP_TITLE?: string;           // 应用标题
    VITE_APP_LOGO?: string;            // Logo URL
    VITE_BASE_URL?: string;            // 基础路径
    VITE_AI_TYPE?: 'backend' | 'frontend' | 'demo'; // 运行模式
    
    // API
    VITE_API_BASE_URL?: string;        // 接口地址
    VITE_OPENAI_API_KEY?: string;      // API Key
    VITE_OPENAI_MODEL?: string;        // 模型ID

    // UI 配置 (LiteChat 特性)
    VITE_SHOW_HEADER?: string;         // 'true' | 'false'
    VITE_SHOW_FOOTER?: string;         // 'true' | 'false'
    VITE_FOOTER_TEXT?: string;         // 底部文本
    VITE_INPUT_PLACEHOLDER?: string;   // 输入框提示
    VITE_SHOW_UPLOAD_BTN?: string;     // 'true' | 'false'
    VITE_SHOW_AUDIO_BTN?: string;      // 'true' | 'false'
    
    // 其他
    VITE_ENABLE_QR_LOGIN?: string;
    VITE_FLOAT_BALL_IMAGE?: string;
}
```

## 程序化控制

`initAiChat` 返回 (或通过 Hook 暴露) 的实例包含以下方法：

- `open()`: 打开聊天窗口
- `close()`: 关闭聊天窗口
- `toggle()`: 切换打开/关闭状态
- `show()`: 显示组件 (浮球或窗口)
- `hide()`: 隐藏组件 (完全不可见)
- `unmount()`: 销毁组件并从 DOM 移除

## 样式隔离

组件使用 **Shadow DOM + iframe** 双层隔离技术：

- **Shadow DOM**: 隔离浮动按钮和窗口框架样式
- **iframe**: 隔离聊天内容样式，确保响应式布局正确

完全不会影响宿主页面样式，也不会被宿主页面样式影响。

## 📱 响应式设计

组件窗口宽度为 360px 时自动触发移动端布局：

- 隐藏侧边栏
- 启用移动端工具栏
- 优化触摸交互体验

## 开发

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建库
npm run build:lib

# 预览 demo
open test-widget.html
```

## License

MIT

## 浏览器兼容性

本项目使用现代 JavaScript (ESM) 打包，支持以下浏览器：

| 浏览器 | 最低版本 |
|--------|---------|
| Chrome | 87+ |
| Edge | 88+ |
| Firefox | 78+ |
| Safari | 14+ |
| Opera | 73+ |

**注意事项：**
- 需要支持 ES2020 特性（如 `import()`, `Promise.allSettled`, `??` 等）
- 需要支持 ES Modules (ESM)
- 不支持 IE11 及更早版本
- 移动端浏览器：iOS Safari 14+, Chrome for Android 87+

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

- **邮箱**：7as0nch@gmail.com
- **WeChat**：JasonC12o9
- **QQ**：2538684421