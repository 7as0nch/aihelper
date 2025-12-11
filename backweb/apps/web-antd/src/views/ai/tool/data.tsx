import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Tag } from 'ant-design-vue';

// Tool 状态选项
export const toolStatusOptions = [
    { color: 'success', label: '启用', value: 1 },
    { color: 'error', label: '禁用', value: 0 },
];

// Tool 类型选项
export const toolTypeOptions = [
    { label: '内置工具', value: 1 },
    { label: 'MCP 工具', value: 2 },
    { label: '自定义工具', value: 3 },
];

// 系统类型选项
export const sysTypeOptions = [
    { label: '系统', value: 1 },
    { label: '用户', value: 0 },
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
            options: toolTypeOptions,
        },
        fieldName: 'type',
        label: '类型',
    },
    {
        component: 'Select',
        componentProps: {
            options: toolStatusOptions,
        },
        fieldName: 'status',
        label: '状态',
    },
];

export const columns: VxeGridProps['columns'] = [
    { type: 'checkbox', width: 60 },
    {
        title: 'ID',
        field: 'id',
        width: 80,
    },
    {
        title: '名称',
        field: 'name',
        minWidth: 150,
        slots: {
            default: ({ row }) => {
                return <Tag color="processing">{row.name}</Tag>;
            },
        },
    },
    {
        title: '代码',
        field: 'code',
        minWidth: 120,
        slots: {
            default: ({ row }) => {
                return <Tag color="cyan">{row.code}</Tag>;
            },
        },
    },
    {
        title: '描述',
        field: 'description',
        minWidth: 200,
    },
    {
        title: '类型',
        field: 'type',
        width: 100,
        slots: {
            default: ({ row }) => {
                const found = toolTypeOptions.find((item) => item.value === row.type);
                return <Tag color="blue">{found?.label ?? row.type}</Tag>;
            },
        },
    },
    {
        title: '系统类型',
        field: 'sysType',
        width: 100,
        slots: {
            default: ({ row }) => {
                const found = sysTypeOptions.find((item) => item.value === row.sysType);
                return <Tag>{found?.label ?? row.sysType}</Tag>;
            },
        },
    },
    {
        title: '状态',
        field: 'status',
        width: 80,
        slots: { default: 'status' },
    },
    {
        title: '创建时间',
        field: 'createdAt',
        width: 180,
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
        component: 'Select',
        componentProps: {
            allowClear: false,
            options: toolTypeOptions,
        },
        defaultValue: 1,
        fieldName: 'type',
        label: '类型',
    },
    {
        component: 'Select',
        componentProps: {
            allowClear: false,
            options: toolStatusOptions,
        },
        defaultValue: 1,
        fieldName: 'status',
        label: '状态',
    },
    {
        component: 'Textarea',
        fieldName: 'description',
        formItemClass: 'col-span-2',
        label: '描述',
        componentProps: {
            rows: 3,
        },
    },
    {
        component: 'Input',
        fieldName: 'mcpUrl',
        formItemClass: 'col-span-2',
        label: 'MCP URL',
        dependencies: {
            show: (values) => values.type === 2,
            triggerFields: ['type'],
        },
    },
    {
        component: 'Input',
        fieldName: 'mcpToken',
        formItemClass: 'col-span-2',
        label: 'MCP Token',
        dependencies: {
            show: (values) => values.type === 2,
            triggerFields: ['type'],
        },
        componentProps: {
            type: 'password',
        },
    },
];
