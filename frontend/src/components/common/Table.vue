<template>
  <div class="table-container">
    <table class="table">
      <thead>
        <tr>
          <th v-for="column in columns" :key="column.key" class="table-header">
            {{ column.title }}
          </th>
          <th v-if="showActions" class="table-header actions">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in data" :key="item.id" class="table-row">
          <td v-for="column in columns" :key="column.key" class="table-cell">
            <slot :name="column.key" :item="item">
              {{ item[column.key] }}
            </slot>
          </td>
          <td v-if="showActions" class="table-cell actions">
            <button v-if="showEdit" @click="handleEdit(item)" class="btn-edit">编辑</button>
            <button v-if="showDelete" @click="handleDelete(item)" class="btn-delete">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="data.length === 0" class="no-data">暂无数据</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

interface Column {
  key: string;
  title: string;
}

interface DataItem {
  id: string | number;
  [key: string]: any;
}

defineProps<{
  columns: Column[];
  data: DataItem[];
  showActions?: boolean;
  showEdit?: boolean;
  showDelete?: boolean;
}>();

const emit = defineEmits<{
  (e: 'edit', item: DataItem): void;
  (e: 'delete', item: DataItem): void;
}>();

const handleEdit = (item: DataItem) => {
  emit('edit', item);
};

const handleDelete = (item: DataItem) => {
  emit('delete', item);
};
</script>

<style scoped>
.table-container {
  overflow-x: auto;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table-header {
  padding: 1rem 1.5rem;
  text-align: left;
  font-weight: 600;
  color: #495057;
  border-bottom: 2px solid #e9ecef;
  background-color: #f8f9fa;
}

.table-header.actions {
  text-align: center;
}

.table-row:hover {
  background-color: #f8f9fa;
}

.table-cell {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.table-cell.actions {
  text-align: center;
}

.btn-edit,
.btn-delete {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  margin: 0 0.25rem;
  transition: background-color 0.3s;
}

.btn-edit {
  background-color: #ffc107;
  color: #212529;
}

.btn-edit:hover {
  background-color: #e0a800;
}

.btn-delete {
  background-color: #dc3545;
  color: white;
}

.btn-delete:hover {
  background-color: #c82333;
}

.no-data {
  padding: 2rem;
  text-align: center;
  color: #6c757d;
  font-style: italic;
}
</style>