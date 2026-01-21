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
    field: 'originalModelId',
    width: 100,
    slots: {
      default: ({ row }) => {
        return row.originalModelId ? (
          <Tag color="purple">{row.originalModelId}</Tag>
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

// Agent 配置表单 Schema（用于中间面板，按分组返回）
export const modelConfigSchema: (
  onModelChange?: (value?: number | string) => void,
) => ReturnType<FormSchemaGetter> = (onModelChange) => [
  {
    component: 'Input',
    dependencies: {
      show: () => false,
      triggerFields: [''],
    },
    fieldName: 'id',
    label: 'ID',
  },
  {
    component: 'ApiSelect',
    componentProps: {
      api: modelList,
      immediate: true,
      allowClear: true,
      placeholder: '请选择 AI 模型',
      showSearch: true,
      onChange: (value: number | string | undefined) => {
        onModelChange?.(value);
      },
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
    fieldName: 'originalModelId',
    label: 'AI 模型',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'aiModelBaseUrl',
    label: 'Base URL',
    componentProps: {
      placeholder: '覆盖模型的 Base URL',
    },
  },
  {
    component: 'InputPassword',
    fieldName: 'aiModelApiKey',
    label: 'API Key',
    componentProps: {
      placeholder: '覆盖模型的 API Key',
    },
  },
  {
    component: 'InputNumber',
    fieldName: 'aiModelTemperature',
    label: 'Temperature',
    componentProps: {
      min: 0,
      max: 2,
      step: 0.1,
      placeholder: '覆盖 Temperature',
    },
  },
  {
    component: 'InputNumber',
    fieldName: 'aiModelTopP',
    label: 'Top P',
    componentProps: {
      min: 0,
      max: 1,
      step: 0.1,
      placeholder: '覆盖 Top P',
    },
  },
];

export const basicInfoSchema: FormSchemaGetter = () => [
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
    label: '描述',
    componentProps: {
      maxlength: 200,
      placeholder: '请输入 Agent 描述信息',
      showCount: true,
    },
  },
];

export const agentConfigSchema: FormSchemaGetter = () => [
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
];

export const featuresSchema: FormSchemaGetter = () => [
  {
    component: 'Checkbox',
    componentProps: {
      class: 'ml-auto',
    },
    defaultValue: false,
    fieldName: 'withWriteTodos',
    formItemClass:
      'col-span-1 flex items-center justify-between px-4 py-3 rounded-lg border border-slate-200 bg-white shadow-sm',
    label: '启用待办事项',
  },
  {
    component: 'Checkbox',
    componentProps: {
      class: 'ml-auto',
    },
    defaultValue: false,
    fieldName: 'withWebSearchAgent',
    formItemClass:
      'col-span-1 flex items-center justify-between px-4 py-3 rounded-lg border border-slate-200 bg-white shadow-sm',
    label: '启用网络搜索',
  },
];

// 完整配置 Schema（用于兼容）
export const configSchema: FormSchemaGetter = () => [
  ...modelConfigSchema(),
  ...basicInfoSchema(),
  ...agentConfigSchema(),
  ...featuresSchema(),
];

// 保留旧的 drawerSchema 以兼容（如果需要）
export const drawerSchema: FormSchemaGetter = () => configSchema();
