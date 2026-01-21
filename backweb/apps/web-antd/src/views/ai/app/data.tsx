import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { DictEnum } from '@vben/constants';

import { Tag } from 'ant-design-vue';

import { renderDict } from '#/utils/render';
import { formatMessageTime } from '#/utils/time';

// Application 状态选项
export const appStatusOptions = [
  { color: 'success', label: '启用', value: 1 },
  { color: 'error', label: '禁用', value: 2 },
];

// Application 模式选项
export const appModeOptions = [
  { label: '单 Agent 模式', value: 1 },
  { label: '多 Agent 模式', value: 2 },
];

// Application 类型选项
export const appTypeOptions = [
  { label: '预定义', value: 1 },
  { label: '自定义', value: 2 },
];

// 作用粒度选项
export const appScopeOptions = [
  { label: '所有人', value: 1 },
  { label: '指定角色', value: 2 },
  { label: '指定用户', value: 3 },
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
      options: appStatusOptions,
    },
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'Select',
    componentProps: {
      options: appModeOptions,
    },
    fieldName: 'mode',
    label: '模式',
  },
  {
    component: 'Select',
    componentProps: {
      options: appTypeOptions,
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
    minWidth: 150,
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
    title: '版本',
    field: 'version',
    width: 100,
    slots: {
      default: ({ row }) => {
        return row.version ? (
          <Tag color="cyan">{row.version}</Tag>
        ) : (
          <span class="text-gray-400">-</span>
        );
      },
    },
  },
  {
    title: '模式',
    field: 'mode',
    width: 140,
    slots: {
      default: ({ row }) => {
        const found = appModeOptions.find((item) => item.value === row.mode);
        return <Tag color="blue">{found?.label ?? row.mode}</Tag>;
      },
    },
  },
  {
    title: '类型',
    field: 'type',
    width: 100,
    slots: {
      default: ({ row }) => {
        const found = appTypeOptions.find((item) => item.value === row.type);
        return <Tag color="purple">{found?.label ?? row.type}</Tag>;
      },
    },
  },
  {
    title: '作用粒度',
    field: 'scope',
    width: 120,
    slots: {
      default: ({ row }) => {
        const found = appScopeOptions.find((item) => item.value === row.scope);
        return <Tag color="orange">{found?.label ?? row.scope}</Tag>;
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
    width: 180,
  },
];

// 基本信息表单
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
      maxlength: 500,
      placeholder: '请输入应用描述信息',
      showCount: true,
    },
  },
  {
    component: 'Input',
    fieldName: 'version',
    label: '版本号',
    componentProps: {
      placeholder: '如：1.0.0',
    },
  },
];

// 应用配置表单
export const appConfigSchema: FormSchemaGetter = () => [
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: appModeOptions,
      placeholder: '请选择模式',
    },
    defaultValue: 1,
    fieldName: 'mode',
    label: '模式',
    rules: 'required',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: appTypeOptions,
      placeholder: '请选择类型',
    },
    defaultValue: 2,
    fieldName: 'type',
    label: '类型',
    rules: 'required',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: appScopeOptions,
      placeholder: '请选择作用粒度',
    },
    defaultValue: 1,
    fieldName: 'scope',
    label: '作用粒度',
    rules: 'required',
  },
  {
    component: 'Select',
    componentProps: {
      allowClear: false,
      options: appStatusOptions,
      placeholder: '请选择状态',
    },
    defaultValue: 1,
    fieldName: 'status',
    label: '状态',
    rules: 'required',
  },
];
