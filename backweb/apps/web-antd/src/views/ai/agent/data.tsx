import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';

import { Tag } from 'ant-design-vue';

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
    title: '注册类型',
    field: 'adapterType',
    width: 120,
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
  {
    component: 'Input',
    fieldName: 'name',
    label: '名称',
    rules: 'required',
  },
  {
    component: 'Input',
    fieldName: 'code',
    label: '代码',
    rules: 'required',
  },
  {
    component: 'Textarea',
    fieldName: 'description',
    formItemClass: 'col-span-2',
    label: '描述',
    componentProps: {
      maxlength: 200,
      showCount: true,
    },
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: adapterTypeOptions,
    },
    defaultValue: 1,
    fieldName: 'adapterType',
    label: '适配器类型',
  },
  {
    component: 'InputNumber',
    fieldName: 'aiModelId',
    label: '模型 ID',
  },
  {
    component: 'InputNumber',
    componentProps: {
      min: 1,
      max: 100,
    },
    defaultValue: 10,
    fieldName: 'maxIteration',
    label: '最大迭代',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: agentTypeOptions,
    },
    defaultValue: 1,
    fieldName: 'type',
    label: '类型',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: agentStatusOptions,
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'InputNumber',
    defaultValue: 0,
    fieldName: 'order',
    label: '排序',
  },
  {
    component: 'Switch',
    defaultValue: false,
    fieldName: 'withWriteTodos',
    label: '开启 Todos',
  },
  {
    component: 'Switch',
    defaultValue: false,
    fieldName: 'withWebSearchAgent',
    label: '网络搜索',
  },
  {
    component: 'Textarea',
    fieldName: 'systemPrompt',
    formItemClass: 'col-span-2',
    label: '系统提示词',
    componentProps: {
      rows: 4,
    },
  },
  {
    component: 'Textarea',
    fieldName: 'userInputPrompt',
    formItemClass: 'col-span-2',
    label: '用户输入提示',
    componentProps: {
      rows: 4,
    },
  },
];
