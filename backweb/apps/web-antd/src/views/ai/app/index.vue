<script setup lang="ts">
import type { VbenFormProps } from '@vben/common-ui';

import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ApplicationInfo } from '#/api/ai/application/model';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import { getVxePopupContainer } from '@vben/utils';

import { Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid, vxeCheckboxChecked } from '#/adapter/vxe-table';
import { applicationList, applicationRemove } from '#/api/ai/application';

import AppEditorModal from './app-editor-modal.vue';
import { columns, querySchema } from './data';

const formOptions: VbenFormProps = {
  commonConfig: {
    labelWidth: 80,
    componentProps: {
      allowClear: true,
    },
  },
  schema: querySchema(),
  wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
};

const gridOptions: VxeGridProps = {
  checkboxConfig: {
    highlight: true,
    reserve: true,
  },
  columns,
  height: 'auto',
  keepSource: true,
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues = {}) => {
        return await applicationList({
          pageNum: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
      },
    },
  },
  rowConfig: {
    keyField: 'id',
  },
  id: 'ai-application-index',
};

const [BasicTable, tableApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

// 弹框状态
const editorVisible = ref(false);
const currentApplicationId = ref<number | string | undefined>(undefined);

function handleAdd() {
  currentApplicationId.value = undefined;
  editorVisible.value = true;
}

async function handleEdit(record: ApplicationInfo) {
  currentApplicationId.value = record.id;
  editorVisible.value = true;
}

function handleEditorSuccess() {
  tableApi.query();
}

async function handleDelete(row: ApplicationInfo) {
  await applicationRemove(row.id);
  await tableApi.query();
}

function handleMultiDelete() {
  const rows = tableApi.grid.getCheckboxRecords();
  const ids = rows.map((row: ApplicationInfo) => row.id);
  Modal.confirm({
    title: '提示',
    okType: 'danger',
    content: `确认删除选中的${ids.length}条记录吗？`,
    onOk: async () => {
      for (const id of ids) {
        await applicationRemove(id);
      }
      await tableApi.query();
    },
  });
}
</script>

<template>
  <Page :auto-content-height="true">
    <BasicTable table-title="AI 应用列表">
      <template #toolbar-tools>
        <Space>
          <a-button
            :disabled="!vxeCheckboxChecked(tableApi)"
            danger
            type="primary"
            @click="handleMultiDelete"
          >
            批量删除
          </a-button>
          <a-button type="primary" @click="handleAdd"> 新增应用 </a-button>
        </Space>
      </template>
      <template #status="{ row }">
        <a-tag :color="row.status === 1 ? 'success' : 'error'">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </a-tag>
      </template>
      <template #action="{ row }">
        <Space>
          <ghost-button @click.stop="handleEdit(row)"> 编辑流程 </ghost-button>
          <Popconfirm
            :get-popup-container="getVxePopupContainer"
            placement="left"
            title="确认删除？"
            @confirm="handleDelete(row)"
          >
            <ghost-button danger @click.stop=""> 删除 </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </BasicTable>

    <!-- 编辑器弹框 -->
    <AppEditorModal
      :visible="editorVisible"
      :application-id="currentApplicationId"
      @update:visible="(val) => (editorVisible = val)"
      @success="handleEditorSuccess"
    />
  </Page>
</template>
