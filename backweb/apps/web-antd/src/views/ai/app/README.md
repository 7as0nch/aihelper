# AI Application 流程编辑器

## 功能概述

可视化的 AI Application 流程编辑器，支持拖拽节点、连线配置，并提供弹框式编辑体验。

## 核心功能

### 1. 弹框模式
- ✅ 支持弹框打开编辑器
- ✅ 弹框可拖拽移动（非全屏模式下）
- ✅ 支持全屏/窗口模式切换
- ✅ 自动保存并刷新列表

### 2. 流程图编辑
- ✅ 可视化节点拖拽
- ✅ 节点连线
- ✅ 节点选中和配置
- ✅ 批量删除节点
- ✅ 缩放和平移画布

### 3. 节点管理
- ✅ 左侧 Agent 节点库
- ✅ 从现有 Agent 创建节点
- ✅ 创建空白节点
- ✅ 右侧节点配置面板

### 4. 数据持久化
- ✅ Schema JSON 序列化
- ✅ 自动加载流程图状态
- ✅ 保存到后端 API

## 文件结构

```
views/ai/app/
├── index.vue                  # 列表页（主入口）
├── app-editor-modal.vue       # 弹框编辑器 ⭐ 新增
├── app-editor.vue             # 原独立页面编辑器（保留用于路由）
├── data.tsx                   # 数据配置
└── components/
    ├── AgentNode.vue          # 自定义节点组件
    ├── AgentLibrary.vue       # Agent 节点库
    └── NodeConfigPanel.vue    # 节点配置面板
```

## 使用方式

### 列表页调用

```vue
<script setup>
import { ref } from 'vue';
import AppEditorModal from './app-editor-modal.vue';

const editorVisible = ref(false);
const currentApplicationId = ref(undefined);

// 新增
function handleAdd() {
  currentApplicationId.value = undefined;
  editorVisible.value = true;
}

// 编辑
function handleEdit(record) {
  currentApplicationId.value = record.id;
  editorVisible.value = true;
}

// 保存成功回调
function handleEditorSuccess() {
  // 刷新列表
  tableApi.query();
}
</script>

<template>
  <AppEditorModal
    :visible="editorVisible"
    :application-id="currentApplicationId"
    @update:visible="(val) => (editorVisible = val)"
    @success="handleEditorSuccess"
  />
</template>
```

## 弹框功能详解

### 拖拽功能
- 📌 点击顶部标题栏可拖动弹框
- 📌 拖拽图标提示用户可拖动
- 📌 全屏模式下禁用拖拽
- 📌 退出全屏自动恢复居中位置

**实现原理：**
```typescript
const handleDragStart = (e: MouseEvent) => {
  if (isFullscreen.value) return;
  isDragging.value = true;
  // 记录起始位置
  dragState.value = {
    startX: e.clientX,
    startY: e.clientY,
    startLeft: parseInt(modalStyle.value.left) || 0,
    startTop: parseInt(modalStyle.value.top) || 0,
  };
  // 监听移动和释放事件
  document.addEventListener('mousemove', handleDragMove);
  document.addEventListener('mouseup', handleDragEnd);
};
```

### 全屏功能
- 📺 点击全屏按钮切换全屏/窗口模式
- 📺 全屏时宽高为 100%
- 📺 窗口模式宽度 90%，高度 85vh
- 📺 自动适配画布视图

**切换代码：**
```typescript
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value;
  if (!isFullscreen.value) {
    // 退出全屏时重置位置
    modalStyle.value = {
      top: '50px',
      left: '50%',
      transform: 'translateX(-50%)',
    };
  }
};
```

## Props 和 Events

### AppEditorModal Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `visible` | `boolean` | `false` | 弹框显示状态 |
| `applicationId` | `number \| string` | - | 应用 ID（编辑模式） |

### AppEditorModal Events

| 事件 | 参数 | 说明 |
|------|------|------|
| `update:visible` | `(value: boolean)` | 更新显示状态 |
| `success` | - | 保存成功 |

## 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Delete` | 删除选中节点 |
| `Ctrl/Cmd + S` | 保存（待实现） |
| `Esc` | 关闭弹框（待实现） |

## 技术栈

- **Vue 3** - Composition API
- **@vue-flow/core** - 流程图核心库
- **Ant Design Vue** - UI 组件库
- **TypeScript** - 类型安全

## 后续优化

- [ ] 添加快捷键支持
- [ ] 节点复制/粘贴功能
- [ ] 流程图导出图片
- [ ] 撤销/重做功能
- [ ] 节点模板保存
- [ ] 流程图验证

## 注意事项

1. **路由保留**：原 `app-editor.vue` 文件保留用于直接路由访问
2. **性能优化**：大型流程图建议启用虚拟滚动
3. **数据同步**：弹框关闭时不自动保存，需手动点击保存按钮
4. **浏览器兼容**：拖拽功能在现代浏览器中测试通过
