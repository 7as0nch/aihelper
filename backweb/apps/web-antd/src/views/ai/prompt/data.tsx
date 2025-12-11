import type { FormSchemaGetter } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Tag } from 'ant-design-vue';

// Prompt 类型选项
export const promptTypeOptions = [
    { label: '系统提示', value: 1 },
    { label: '用户提示', value: 2 },
    { label: '模板', value: 3 },
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
            options: promptTypeOptions,
        },
        fieldName: 'type',
        label: '类型',
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
        title: '类型',
        field: 'type',
        width: 100,
        slots: {
            default: ({ row }) => {
                const found = promptTypeOptions.find((item) => item.value === row.type);
                return <Tag color="blue">{found?.label ?? row.type}</Tag>;
            },
        },
    },
    {
        title: '描述',
        field: 'description',
        minWidth: 200,
    },
    {
        title: '内容预览',
        field: 'text',
        minWidth: 300,
        slots: {
            default: ({ row }) => {
                const text = row.text || '';
                const preview = text.length > 100 ? `${text.slice(0, 100)}...` : text;
                return <span class="text-gray-500">{preview}</span>;
            },
        },
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
        component: 'Select',
        componentProps: {
            allowClear: false,
            options: promptTypeOptions,
        },
        defaultValue: 1,
        fieldName: 'type',
        label: '类型',
    },
    {
        component: 'Textarea',
        fieldName: 'description',
        formItemClass: 'col-span-2',
        label: '描述',
        componentProps: {
            rows: 2,
        },
    },
    {
        component: 'Textarea',
        fieldName: 'text',
        formItemClass: 'col-span-2',
        label: '提示词内容',
        rules: 'required',
        componentProps: {
            rows: 10,
            placeholder: '请输入提示词内容...',
        },
    },
];
