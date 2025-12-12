import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';

import { Tag } from 'ant-design-vue';

import { renderDict } from '#/utils/render';
import { formatMessageTime } from '#/utils/time';

// Model 状态选项
export const modelStatusOptions = [
  { color: 'success', label: '启用', value: 1 },
  { color: 'error', label: '禁用', value: 2 },
];

// Model 分类选项
export const modelCategoryOptions = [
  { label: '文本模型', value: 1 },
  { label: '图像模型', value: 2 },
  { label: '多模态模型', value: 3 },
];

// Model 类型选项
export const modelTypeOptions = [
  { label: '豆包ark', value: 'ark' },
  { label: 'Open AI', value: 'openai' },
  { label: 'DeepSeek', value: 'deepseek' },
];

// 价格类型选项
export const priceTypeOptions = [
  { label: '按 Token', value: 1 },
  { label: '按次', value: 2 },
  { label: '免费', value: 0 },
];

// 是否默认选项
export const isDefaultOptions = [
  { label: '是', value: 1 },
  { label: '否', value: 2 },
];

export const querySchema: FormSchemaGetter = () => [
  {
    component: 'Input',
    fieldName: 'modelName',
    label: '模型名称',
  },
  {
    component: 'Select',
    componentProps: {
      options: modelStatusOptions,
    },
    fieldName: 'status',
    label: '状态',
  },
];

export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 60, fixed: 'left' },
  {
    title: '名称',
    field: 'modelName',
    minWidth: 150,
    fixed: 'left',
  },
  {
    title: '类型',
    field: 'modelType',
    width: 120,
    slots: {
      default: ({ row }) => {
        const found = modelTypeOptions.find(
          (item) => item.value === row.modelType,
        );
        return <Tag color="blue">{found?.label ?? row.modelType}</Tag>;
      },
    },
  },
  {
    title: '分类',
    field: 'category',
    width: 100,
    slots: {
      default: ({ row }) => {
        const found = modelCategoryOptions.find(
          (item) => item.value === row.category,
        );
        return <Tag color="blue">{found?.label ?? row.category}</Tag>;
      },
    },
  },
  {
    title: '供应商',
    field: 'supplier',
    width: 120,
  },
  {
    title: 'Base URL',
    field: 'baseUrl',
    minWidth: 200,
  },
  {
    title: '最大 Token',
    field: 'maxTokens',
    width: 100,
  },
  {
    title: '默认',
    field: 'isDefault',
    width: 80,
    slots: {
      default: ({ row }) => {
        return row.isDefault === 1 ? (
          <Tag color="success">是</Tag>
        ) : (
          <Tag color="error">否</Tag>
        );
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
    title: '价格类型',
    field: 'priceType',
    width: 100,
    slots: {
      default: ({ row }) => {
        const found = priceTypeOptions.find(
          (item) => item.value === row.priceType,
        );
        return <Tag color="blue">{found?.label ?? row.priceType}</Tag>;
      },
    },
  },
  {
    title: '价格',
    field: 'price',
    width: 100,
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
  {
    component: 'Input',
    fieldName: 'modelName',
    label: '模型名称',
    rules: 'required',
  },
  {
    component: 'Select',
    fieldName: 'modelType',
    label: '模型类型',
    rules: 'required',
    componentProps: {
      allowClear: false,
      options: modelTypeOptions,
    },
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: modelCategoryOptions,
    },
    defaultValue: 1,
    fieldName: 'category',
    label: '分类',
  },
  {
    component: 'Input',
    fieldName: 'supplier',
    label: '供应商',
  },
  {
    component: 'Input',
    fieldName: 'baseUrl',
    label: 'Base URL',
    formItemClass: 'col-span-2',
  },
  {
    component: 'Input',
    fieldName: 'apiKey',
    label: 'API Key',
    formItemClass: 'col-span-2',
    componentProps: {
      type: 'password',
    },
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 0,
    },
    defaultValue: 4096,
    fieldName: 'maxTokens',
    label: '最大 Token',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 0,
      max: 2,
      step: 0.1,
    },
    defaultValue: 0.7,
    fieldName: 'temperature',
    label: 'Temperature',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 0,
      max: 1,
      step: 0.1,
    },
    defaultValue: 1,
    fieldName: 'topP',
    label: 'Top P',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: priceTypeOptions,
    },
    defaultValue: 1,
    fieldName: 'priceType',
    label: '价格类型',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 0,
      step: 0.01,
    },
    defaultValue: 0,
    fieldName: 'price',
    label: '价格',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: modelStatusOptions,
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: isDefaultOptions,
    },
    defaultValue: 0,
    fieldName: 'isDefault',
    label: '是否默认',
  },
  {
    component: 'Textarea',
    fieldName: 'description',
    formItemClass: 'col-span-2',
    componentProps: {
      maxLength: 255,
    },
    label: '描述',
  },
];
