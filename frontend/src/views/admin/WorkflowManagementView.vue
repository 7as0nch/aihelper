<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ElTable, ElTableColumn, ElPagination, ElInput, ElSelect, ElOption, ElButton, ElDialog, ElForm, ElFormItem, ElSwitch, ElMessageBox, ElMessage, ElTag } from 'element-plus';
import 'element-plus/es/components/table/style/css';
import 'element-plus/es/components/pagination/style/css';
import 'element-plus/es/components/input/style/css';
import 'element-plus/es/components/select/style/css';
import 'element-plus/es/components/button/style/css';
import 'element-plus/es/components/dialog/style/css';
import 'element-plus/es/components/form/style/css';
import 'element-plus/es/components/form-item/style/css';

import 'element-plus/es/components/switch/style/css';
import 'element-plus/es/components/message-box/style/css';
import 'element-plus/es/components/message/style/css';
import 'element-plus/es/components/tag/style/css';

// 模拟工作流数据
interface Workflow {
  id: number;
  name: string;
  description: string;
  status: 'draft' | 'active' | 'disabled';
  createTime: string;
  updateTime: string;
  runCount: number;
  lastRunTime: string;
  isPublic: boolean;
}

const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searchKeyword = ref('');
const statusFilter = ref('all');
const workflowList = ref<Workflow[]>([]);
const loading = ref(false);

// 对话框状态
const dialogVisible = ref(false);
const dialogTitle = ref('添加工作流');
const isEditMode = ref(false);
const currentWorkflowId = ref<number | null>(null);

// 表单数据
const formData = ref({
  name: '',
  description: '',
  status: 'draft' as 'draft' | 'active' | 'disabled',
  isPublic: false,
  steps: []
});

// 表单验证
const formRules = ref({
  name: [
    { required: true, message: '请输入工作流名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入工作流描述', trigger: 'blur' },
    { min: 5, max: 200, message: '描述长度在 5 到 200 个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择工作流状态', trigger: 'change' }
  ]
});

// 筛选后的工作流列表
const filteredWorkflows = computed(() => {
  return workflowList.value.filter(workflow => {
    // 关键词筛选
    const keywordMatch = !searchKeyword.value || 
      workflow.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      workflow.description.toLowerCase().includes(searchKeyword.value.toLowerCase());
    
    // 状态筛选
    const statusMatch = statusFilter.value === 'all' || workflow.status === statusFilter.value;
    
    return keywordMatch && statusMatch;
  });
});

// 分页后的工作流列表
const paginatedWorkflows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredWorkflows.value.slice(start, end);
});

// 初始化数据
onMounted(() => {
  fetchWorkflows();
});

// 模拟获取工作流数据
const fetchWorkflows = () => {
  loading.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    const statuses: Array<'draft' | 'active' | 'disabled'> = ['draft', 'active', 'disabled'];
    const mockWorkflows: Workflow[] = Array.from({ length: 20 }, (_, i) => ({
        id: i + 1 as number,
        name: `Workflow${i + 1}` as string,
        description: `这是工作流 ${i + 1} 的详细描述，用于自动化执行一系列任务和操作。` as string,
        status: statuses[i % 3] as 'draft' | 'active' | 'disabled',
        createTime: new Date(Date.now() - Math.floor(Math.random() * 90) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN') as string,
        updateTime: new Date(Date.now() - Math.floor(Math.random() * 30) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN') as string,
        runCount: Math.floor(Math.random() * 500) + 5 as number,
        lastRunTime: i % 5 === 0 ? '--' : new Date(Date.now() - Math.floor(Math.random() * 7) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN') as string,
        isPublic: Math.random() > 0.6 as boolean
      } as Workflow));
    
    workflowList.value = mockWorkflows;
    total.value = mockWorkflows.length;
    loading.value = false;
  }, 800);
};

// 搜索工作流
const handleSearch = () => {
  currentPage.value = 1;
};

// 重置筛选条件
const handleReset = () => {
  searchKeyword.value = '';
  statusFilter.value = 'all';
  currentPage.value = 1;
};

// 打开添加工作流对话框
const openAddDialog = () => {
  dialogTitle.value = '添加工作流';
  isEditMode.value = false;
  currentWorkflowId.value = null;
  formData.value = {
    name: '',
    description: '',
    status: 'draft',
    isPublic: false,
    steps: []
  };
  dialogVisible.value = true;
};

// 打开编辑工作流对话框
const openEditDialog = (workflow: Workflow) => {
  dialogTitle.value = '编辑工作流';
  isEditMode.value = true;
  currentWorkflowId.value = workflow.id;
  formData.value = {
    name: workflow.name,
    description: workflow.description,
    status: workflow.status,
    isPublic: workflow.isPublic,
    steps: [] // 简化处理，实际应包含步骤配置
  };
  dialogVisible.value = true;
};

// 保存工作流（添加或编辑）
const saveWorkflow = () => {
  // 模拟API请求延迟
  loading.value = true;
  
  setTimeout(() => {
    if (isEditMode.value && currentWorkflowId.value !== null) {
      // 编辑工作流
      const index = workflowList.value.findIndex(workflow => workflow.id === currentWorkflowId.value);
      if (index !== -1 && workflowList.value[index]) {
        const updatedWorkflow = workflowList.value[index];
        workflowList.value[index] = {
          id: updatedWorkflow.id,
          name: formData.value.name,
          description: formData.value.description,
          status: formData.value.status,
          createTime: updatedWorkflow.createTime,
          updateTime: new Date().toLocaleDateString('zh-CN'),
          runCount: updatedWorkflow.runCount,
          lastRunTime: updatedWorkflow.lastRunTime,
          isPublic: formData.value.isPublic
        };
      }
      ElMessage.success('工作流更新成功');
    } else {
      // 添加工作流
      const newWorkflow: Workflow = {
        id: workflowList.value.length + 1,
        name: formData.value.name,
        description: formData.value.description,
        status: formData.value.status,
        createTime: new Date().toLocaleDateString('zh-CN'),
        updateTime: new Date().toLocaleDateString('zh-CN'),
        runCount: 0,
        lastRunTime: '--',
        isPublic: formData.value.isPublic
      };
      workflowList.value.unshift(newWorkflow);
      total.value++;
      ElMessage.success('工作流添加成功');
    }
    
    dialogVisible.value = false;
    loading.value = false;
  }, 800);
};

// 删除工作流
const deleteWorkflow = (workflowId: number, workflowName: string) => {
  ElMessageBox.confirm(
    `确定要删除工作流「${workflowName}」吗？此操作不可撤销，可能会影响依赖此工作流的服务。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
  .then(() => {
    loading.value = true;
    
    // 模拟API请求延迟
    setTimeout(() => {
      const index = workflowList.value.findIndex(workflow => workflow.id === workflowId);
      if (index !== -1) {
        workflowList.value.splice(index, 1);
        total.value--;
        // 如果删除后当前页没有数据，则返回上一页
        if (paginatedWorkflows.value.length === 0 && currentPage.value > 1) {
          currentPage.value--;
        }
      }
      ElMessage.success('工作流删除成功');
      loading.value = false;
    }, 500);
  })
  .catch(() => {
    // 用户取消删除
  });
};

// 运行工作流
const runWorkflow = (workflowId: number, workflowName: string) => {
  ElMessageBox.confirm(
    `确定要立即运行工作流「${workflowName}」吗？`,
    '运行确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  )
  .then(() => {
    loading.value = true;
    
    // 模拟API请求延迟
    setTimeout(() => {
      const index = workflowList.value.findIndex(workflow => workflow.id === workflowId);
      if (index !== -1 && workflowList.value[index]) {
        const updatedWorkflow = workflowList.value[index];
        workflowList.value[index] = {
          ...updatedWorkflow,
          runCount: updatedWorkflow.runCount + 1,
          lastRunTime: new Date().toLocaleDateString('zh-CN')
        } as Workflow;
      }
      ElMessage.success('工作流开始运行');
      loading.value = false;
    }, 500);
  })
  .catch(() => {
    // 用户取消运行
  });
};
</script>

<template>
  <div class="workflow-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">工作流管理</h1>
      <p class="page-description">管理系统中的所有工作流，包括添加、编辑、删除和运行工作流。</p>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="search-filter-section">
      <div class="search-filter-left">
        <ElInput
          v-model="searchKeyword"
          placeholder="搜索名称或描述"
          prefix-icon="Search"
          clearable
          class="search-input"
          @keyup.enter="handleSearch"
        />
        
        <ElSelect
          v-model="statusFilter"
          placeholder="选择状态"
          class="filter-select"
          @change="handleSearch"
        >
          <ElOption label="全部" value="all" />
          <ElOption label="草稿" value="draft" />
          <ElOption label="活跃" value="active" />
          <ElOption label="禁用" value="disabled" />
        </ElSelect>
        
        <ElButton type="primary" @click="handleSearch" class="search-btn">搜索</ElButton>
        <ElButton @click="handleReset" class="reset-btn">重置</ElButton>
      </div>
      
      <div class="search-filter-right">
        <ElButton type="primary" @click="openAddDialog" class="add-btn">添加工作流</ElButton>
      </div>
    </div>

    <!-- 工作流列表 -->
    <div class="workflow-table-section">
      <ElTable
        v-loading="loading"
        :data="paginatedWorkflows"
        style="width: 100%"
        border
        class="workflow-table"
      >
        <ElTableColumn prop="id" label="ID" width="80" align="center" />
        
        <ElTableColumn prop="name" label="工作流名称" min-width="150" />
        
        <ElTableColumn prop="description" label="描述" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.description" placement="top">
              <div class="description-text">{{ row.description }}</div>
            </el-tooltip>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <ElTag :type="getStatusTagType(row.status)">
              {{ getStatusLabel(row.status) }}
            </ElTag>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="runCount" label="运行次数" width="100" align="center" />
        
        <ElTableColumn prop="lastRunTime" label="最后运行" width="120" />
        
        <ElTableColumn prop="isPublic" label="是否公开" width="100" align="center">
          <template #default="{ row }">
            <ElTag :type="row.isPublic ? 'success' : 'info'">
              {{ row.isPublic ? '公开' : '私有' }}
            </ElTag>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="createTime" label="创建时间" width="120" />
        
        <ElTableColumn prop="updateTime" label="更新时间" width="120" />
        
        <ElTableColumn label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <ElButton 
              type="primary" 
              size="small" 
              @click="runWorkflow(row.id, row.name)" 
              class="run-btn"
              :disabled="row.status !== 'active'"
            >
              运行
            </ElButton>
            <ElButton type="primary" size="small" @click="openEditDialog(row)" class="edit-btn">编辑</ElButton>
            <ElButton type="danger" size="small" @click="deleteWorkflow(row.id, row.name)" class="delete-btn">删除</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
    </div>

    <!-- 分页 -->
    <div class="pagination-section">
      <ElPagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="filteredWorkflows.length"
      />
    </div>

    <!-- 添加/编辑工作流对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      :before-close="() => dialogVisible = false"
    >
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <ElFormItem label="工作流名称" prop="name">
          <ElInput v-model="formData.name" placeholder="请输入工作流名称" />
        </ElFormItem>
        
        <ElFormItem label="描述" prop="description">
          <ElInput v-model="formData.description" placeholder="请输入工作流描述" type="textarea" :rows="3" />
        </ElFormItem>
        
        <ElFormItem label="状态" prop="status">
          <ElSelect v-model="formData.status" placeholder="请选择工作流状态">
            <ElOption label="草稿" value="draft" />
            <ElOption label="活跃" value="active" />
            <ElOption label="禁用" value="disabled" />
          </ElSelect>
        </ElFormItem>
        
        <ElFormItem label="是否公开">
          <ElSwitch v-model="formData.isPublic" />
        </ElFormItem>
        
        <ElFormItem label="工作流步骤">
          <div class="steps-placeholder">
            <p>工作流步骤配置区域</p>
            <p class="placeholder-tip">在此可以配置工作流的具体执行步骤和条件</p>
          </div>
        </ElFormItem>
      </ElForm>
      
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveWorkflow" :loading="loading">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts">
// 辅助函数
const getStatusLabel = (status: string): string => {
  switch (status) {
    case 'draft':
      return '草稿';
    case 'active':
      return '活跃';
    case 'disabled':
      return '禁用';
    default:
      return '未知';
  }
};

const getStatusTagType = (status: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  switch (status) {
    case 'draft':
      return 'info';
    case 'active':
      return 'success';
    case 'disabled':
      return 'danger';
    default:
      return 'info';
  }
};
</script>

<style scoped>
.workflow-management {
  padding: 20px 0;
}

/* 页面标题 */
.page-header {
  margin-bottom: 30px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 10px 0;
}

.page-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

/* 搜索和筛选区域 */
.search-filter-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.search-filter-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.search-input {
  width: 300px;
}

.filter-select {
  width: 150px;
}

.search-btn,
.reset-btn,
.add-btn {
  white-space: nowrap;
}

/* 工作流表格 */
.workflow-table-section {
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 20px;
  margin-bottom: 20px;
}

.workflow-table {
  border: 1px solid var(--border-color);
}

.description-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
}

.run-btn,
.edit-btn,
.delete-btn {
  margin: 0 4px;
}

/* 分页 */
.pagination-section {
  display: flex;
  justify-content: flex-end;
  padding: 20px;
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

/* 步骤配置占位符 */
.steps-placeholder {
  padding: 30px;
  text-align: center;
  background-color: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.steps-placeholder p {
  margin: 0 0 5px 0;
  color: var(--text-secondary);
}

.placeholder-tip {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .search-filter-section {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .search-filter-left {
    justify-content: center;
  }
  
  .search-filter-right {
    display: flex;
    justify-content: center;
  }
  
  .search-input {
    width: 250px;
  }
}

@media (max-width: 768px) {
  .workflow-management {
    padding: 15px 0;
  }
  
  .page-header {
    margin-bottom: 20px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .search-filter-left {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input,
  .filter-select {
    width: 100%;
  }
  
  .workflow-table-section {
    padding: 15px;
    overflow-x: auto;
  }
  
  .pagination-section {
    padding: 15px;
  }
}
</style>