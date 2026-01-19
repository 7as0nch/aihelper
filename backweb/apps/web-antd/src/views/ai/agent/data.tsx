import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';

import { Tag } from 'ant-design-vue';

import { modelList } from '#/api/ai/model';
import { renderDict } from '#/utils/render';
import { formatMessageTime } from '#/utils/time';

// Agent 状态选项
export const agentStatusOptions = [
  { color: 'success', label: '启用', value: 1 },
  { color: 'error', label: '禁用', value: 2 },
];

// Agent 类型选项
export const agentTypeOptions = [
  { label: '根 Agent', value: 1 },
  { label: '子 Agent', value: 2 },
];

// 适配器类型选项
export const adapterTypeOptions = [
  { label: 'Eino ADK', value: 1 },
  { label: 'Eino DeepADK', value: 2 },
];

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'name',
    label: '名称',
  },
  {
    component: 'Select',
    componentProps: {
      options: agentStatusOptions,
    },
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'Select',
    componentProps: {
      options: agentTypeOptions,
    },
    fieldName: 'type',
    label: '类型',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60, fixed: 'left' },
  {
    title: '名称',
    field: 'name',
    minWidth: 120,
    fixed: 'left',
  },
  {
    title: '编码',
    field: 'code',
    minWidth: 120,
    slots: {
      default: ({ row }) => {
        return <Tag color="processing">{row.code}</Tag>;
      },
    },
  },
  {
    title: '描述',
    field: 'description',
    minWidth: 200,
  },
  {
    title: '适配器类型',
    field: 'adapterType',
    width: 130,
    slots: {
      default: ({ row }) => {
        const found = adapterTypeOptions.find(
          (item) => item.value === row.adapterType,
        );
        return <Tag color="blue">{found?.label ?? row.adapterType}</Tag>;
      },
    },
  },
  {
    title: '模型 ID',
    field: 'aiModelId',
    width: 100,
    slots: {
      default: ({ row }) => {
        return row.aiModelId ? (
          <Tag color="purple">{row.aiModelId}</Tag>
        ) : (
          <span class="text-gray-400">-</span>
        );
      },
    },
  },
  {
    title: '类型',
    field: 'type',
    width: 100,
    slots: {
      default: ({ row }) => {
        const found = agentTypeOptions.find((item) => item.value === row.type);
        return <Tag color="blue">{found?.label ?? row.type}</Tag>;
      },
    },
  },
  {
    title: '状态',
    field: 'status',
    width: 80,
    slots: {
      default: ({ row }) => {
        return renderDict(row.status, DictEnum.SYS_NORMAL_DISABLE);
      },
    },
  },
  {
    title: '创建时间',
    field: 'createdAt',
    width: 180,
    // formatMessageTime
    slots: {
      default: ({ row }) => {
        return formatMessageTime(row.createdAt);
      },
    },
  },
  {
    field: 'action',
    fixed: 'right',
    slots: { default: 'action' },
    title: '操作',
    width: 150,
  },
];

export const drawerSchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    dependencies: {
      show: () => false,
      triggerFields: [''],
    },
    fieldName: 'id',
    label: 'ID',
  },
  // ========== 基本信息 ==========
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'basicInfo',
    formItemClass: 'col-span-2',
    label: '基本信息',
  },
  {
    component: 'Input',
    fieldName: 'name',
    label: '名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'code',
    label: '编码',
    rules: 'required',
    componentProps: {
      placeholder: '请输入唯一编码',
    },
  },
  {
    component: 'Textarea',
    fieldName: 'description',
    formItemClass: 'col-span-2',
    label: '描述',
    componentProps: {
      maxlength: 200,
      placeholder: '请输入 Agent 描述信息',
      showCount: true,
    },
  },
  // ========== 配置信息 ==========
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'configInfo',
    formItemClass: 'col-span-2',
    label: '配置信息',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: adapterTypeOptions,
      placeholder: '请选择适配器类型',
    },
    defaultValue: 1,
    fieldName: 'adapterType',
    label: '适配器类型',
  },
  {
    component: 'ApiSelect',
    componentProps: {
      api: modelList,
      immediate: true,
      allowClear: true,
      placeholder: '请选择 AI 模型',
      showSearch: true,
      afterFetch: (data: any) => {
        const models = data?.list || data || [];
        return models.map((item: any) => ({
          label: `${item.modelName} (${item.modelType})`,
          value: item.id,
        }));
      },
      onError: (error: any) => {
        console.error('获取模型列表失败:', error);
      },
    },
    fieldName: 'aiModelId',
    label: 'AI 模型',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 1,
      max: 100,
      placeholder: '请输入最大迭代次数',
    },
    defaultValue: 10,
    fieldName: 'maxIteration',
    label: '最大迭代次数',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: agentTypeOptions,
      placeholder: '请选择 Agent 类型',
    },
    defaultValue: 1,
    fieldName: 'type',
    label: 'Agent 类型',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: agentStatusOptions,
      placeholder: '请选择状态',
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 0,
      placeholder: '请输入排序值',
    },
    defaultValue: 0,
    fieldName: 'order',
    label: '排序',
  },
  // ========== 功能开关 ==========
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'featureSwitch',
    formItemClass: 'col-span-2',
    label: '功能开关',
  },
  {
    component: 'Switch',
    componentProps: {
      checkedChildren: '开启',
      unCheckedChildren: '关闭',
    },
    defaultValue: false,
    fieldName: 'withWriteTodos',
    label: '启用 Todos 功能',
  },
  {
    component: 'Switch',
    componentProps: {
      checkedChildren: '开启',
      unCheckedChildren: '关闭',
    },
    defaultValue: false,
    fieldName: 'withWebSearchAgent',
    label: '启用网络搜索',
  },
  // ========== 提示词配置 ==========
  {
    component: 'Divider',
    componentProps: {
      orientation: 'left',
    },
    fieldName: 'promptConfig',
    formItemClass: 'col-span-2',
    label: '提示词配置',
  },
  {
    component: 'Textarea',
    fieldName: 'systemPrompt',
    formItemClass: 'col-span-2',
    label: '系统提示词',
    componentProps: {
      rows: 4,
      placeholder: '请输入系统提示词，用于指导 Agent 的行为',
      showCount: true,
      maxlength: 2000,
    },
  },
  {
    component: 'Textarea',
    fieldName: 'userInputPrompt',
    formItemClass: 'col-span-2',
    label: '用户输入提示词',
    componentProps: {
      rows: 4,
      placeholder: '请输入用户输入提示词，用于格式化用户输入',
      showCount: true,
      maxlength: 2000,
    },
  },
];
