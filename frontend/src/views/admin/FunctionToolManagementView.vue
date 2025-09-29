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
import 'element-plus/es/components/input-number/style/css';
import 'element-plus/es/components/switch/style/css';
import 'element-plus/es/components/message-box/style/css';
import 'element-plus/es/components/message/style/css';
import 'element-plus/es/components/tag/style/css';

// 模拟函数工具数据
interface FunctionTool {
  id: number;
  name: string;
  description: string;
  type: 'api' | 'script' | 'custom';
  status: boolean;
  createTime: string;
  updateTime: string;
  usageCount: number;
  isPublic: boolean;
}

const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searchKeyword = ref('');
const typeFilter = ref('all');
const statusFilter = ref('all');
const functionList = ref<FunctionTool[]>([]);
const loading = ref(false);

// 对话框状态
const dialogVisible = ref(false);
const dialogTitle = ref('添加函数工具');
const isEditMode = ref(false);
const currentFunctionId = ref<number | null>(null);

// 表单数据
const formData = ref({
  name: '',
  description: '',
  type: 'api' as 'api' | 'script' | 'custom',
  status: true,
  isPublic: false,
  code: ''
});

// 表单验证
const formRules = ref({
  name: [
    { required: true, message: '请输入函数工具名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入函数工具描述', trigger: 'blur' },
    { min: 5, max: 200, message: '描述长度在 5 到 200 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择函数工具类型', trigger: 'change' }
  ],
  code: [
    { required: true, message: '请输入函数代码', trigger: 'blur' },
    { min: 10, message: '代码不能少于 10 个字符', trigger: 'blur' }
  ]
});

// 筛选后的函数工具列表
const filteredFunctions = computed(() => {
  return functionList.value.filter(func => {
    // 关键词筛选
    const keywordMatch = !searchKeyword.value || 
      func.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      func.description.toLowerCase().includes(searchKeyword.value.toLowerCase());
    
    // 类型筛选
    const typeMatch = typeFilter.value === 'all' || func.type === typeFilter.value;
    
    // 状态筛选
    const statusMatch = statusFilter.value === 'all' || func.status === (statusFilter.value === 'active');
    
    return keywordMatch && typeMatch && statusMatch;
  });
});

// 分页后的函数工具列表
const paginatedFunctions = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredFunctions.value.slice(start, end);
});

// 初始化数据
onMounted(() => {
  fetchFunctions();
});

// 模拟获取函数工具数据
const fetchFunctions = () => {
  loading.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    const mockFunctions: FunctionTool[] = Array.from({ length: 30 }, (_, i) => ({
      id: i + 1,
      name: `FunctionTool${i + 1}`,
      description: `这是函数工具 ${i + 1} 的详细描述，用于处理特定的业务逻辑和数据转换。`,
      type: ['api', 'script', 'custom'][i % 3] as 'api' | 'script' | 'custom',
      status: Math.random() > 0.1,
      createTime: new Date(Date.now() - Math.floor(Math.random() * 60) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN'),
      updateTime: new Date(Date.now() - Math.floor(Math.random() * 30) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN'),
      usageCount: Math.floor(Math.random() * 1000) + 10,
      isPublic: Math.random() > 0.5
    }));
    
    functionList.value = mockFunctions;
    total.value = mockFunctions.length;
    loading.value = false;
  }, 800);
};

// 搜索函数工具
const handleSearch = () => {
  currentPage.value = 1;
};

// 重置筛选条件
const handleReset = () => {
  searchKeyword.value = '';
  typeFilter.value = 'all';
  statusFilter.value = 'all';
  currentPage.value = 1;
};

// 打开添加函数工具对话框
const openAddDialog = () => {
  dialogTitle.value = '添加函数工具';
  isEditMode.value = false;
  currentFunctionId.value = null;
  formData.value = {
    name: '',
    description: '',
    type: 'api',
    status: true,
    isPublic: false,
    code: ''
  };
  dialogVisible.value = true;
};

// 打开编辑函数工具对话框
const openEditDialog = (func: FunctionTool) => {
  dialogTitle.value = '编辑函数工具';
  isEditMode.value = true;
  currentFunctionId.value = func.id;
  formData.value = {
    name: func.name,
    description: func.description,
    type: func.type,
    status: func.status,
    isPublic: func.isPublic,
    code: `// 函数代码示例\nfunction ${func.name.toLowerCase()}() {\n  // 实现逻辑\n  return true;\n}`
  };
  dialogVisible.value = true;
};

// 保存函数工具（添加或编辑）
const saveFunction = () => {
  // 模拟API请求延迟
  loading.value = true;
  
  setTimeout(() => {
    if (isEditMode.value && currentFunctionId.value !== null) {
      // 编辑函数工具
      const index = functionList.value.findIndex(func => func.id === currentFunctionId.value);
      if (index !== -1) {
        if (functionList.value[index]) {
          functionList.value[index] = {
            id: functionList.value[index].id,
            name: formData.value.name,
            description: formData.value.description,
            type: formData.value.type,
            status: formData.value.status,
            isPublic: formData.value.isPublic,
            updateTime: new Date().toLocaleDateString('zh-CN'),
            createTime: functionList.value[index].createTime || new Date().toLocaleDateString('zh-CN'),
            usageCount: functionList.value[index].usageCount || 0
          };
        }
      }
      ElMessage.success('函数工具更新成功');
    } else {
      // 添加函数工具
      const newFunction: FunctionTool = {
        id: functionList.value.length + 1,
        name: formData.value.name,
        description: formData.value.description,
        type: formData.value.type,
        status: formData.value.status,
        createTime: new Date().toLocaleDateString('zh-CN'),
        updateTime: new Date().toLocaleDateString('zh-CN'),
        usageCount: 0,
        isPublic: formData.value.isPublic
      };
      functionList.value.unshift(newFunction);
      total.value++;
      ElMessage.success('函数工具添加成功');
    }
    
    dialogVisible.value = false;
    loading.value = false;
  }, 800);
};

// 删除函数工具
const deleteFunction = (functionId: number, functionName: string) => {
  ElMessageBox.confirm(
    `确定要删除函数工具「${functionName}」吗？此操作不可撤销，可能会影响依赖此函数的服务。`,
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
      const index = functionList.value.findIndex(func => func.id === functionId);
      if (index !== -1) {
        functionList.value.splice(index, 1);
        total.value--;
        // 如果删除后当前页没有数据，则返回上一页
        if (paginatedFunctions.value.length === 0 && currentPage.value > 1) {
          currentPage.value--;
        }
      }
      ElMessage.success('函数工具删除成功');
      loading.value = false;
    }, 500);
  })
  .catch(() => {
    // 用户取消删除
  });
};

// 切换函数工具状态
const toggleFunctionStatus = (functionId: number, currentStatus: boolean) => {
  loading.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    const index = functionList.value.findIndex(func => func.id === functionId);
    if (index !== -1) {
      if (functionList.value[index]) {
        functionList.value[index].status = !currentStatus;
      }
    }
    ElMessage.success(`函数工具状态已${!currentStatus ? '启用' : '禁用'}`);
    loading.value = false;
  }, 300);
};
</script>

<template>
  <div class="function-tool-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">函数工具管理</h1>
      <p class="page-description">管理系统中的所有函数工具，包括添加、编辑、删除和查看函数工具信息。</p>
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
          v-model="typeFilter"
          placeholder="选择类型"
          class="filter-select"
          @change="handleSearch"
        >
          <ElOption label="全部" value="all" />
          <ElOption label="API调用" value="api" />
          <ElOption label="脚本执行" value="script" />
          <ElOption label="自定义函数" value="custom" />
        </ElSelect>
        
        <ElSelect
          v-model="statusFilter"
          placeholder="选择状态"
          class="filter-select"
          @change="handleSearch"
        >
          <ElOption label="全部" value="all" />
          <ElOption label="启用" value="active" />
          <ElOption label="禁用" value="inactive" />
        </ElSelect>
        
        <ElButton type="primary" @click="handleSearch" class="search-btn">搜索</ElButton>
        <ElButton @click="handleReset" class="reset-btn">重置</ElButton>
      </div>
      
      <div class="search-filter-right">
        <ElButton type="primary" @click="openAddDialog" class="add-btn">添加函数工具</ElButton>
      </div>
    </div>

    <!-- 函数工具列表 -->
    <div class="function-table-section">
      <ElTable
        v-loading="loading"
        :data="paginatedFunctions"
        style="width: 100%"
        border
        class="function-table"
      >
        <ElTableColumn prop="id" label="ID" width="80" align="center" />
        
        <ElTableColumn prop="name" label="函数名称" min-width="150" />
        
        <ElTableColumn prop="description" label="描述" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.description" placement="top">
              <div class="description-text">{{ row.description }}</div>
            </el-tooltip>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="type" label="类型" width="120" align="center">
          <template #default="{ row }">
            <ElTag :type="getTagType(row.type)">
              {{ getTypeLabel(row.type) }}
            </ElTag>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="usageCount" label="使用次数" width="100" align="center" />
        
        <ElTableColumn prop="isPublic" label="是否公开" width="100" align="center">
          <template #default="{ row }">
            <ElTag :type="row.isPublic ? 'success' : 'info'">
              {{ row.isPublic ? '公开' : '私有' }}
            </ElTag>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <ElSwitch
              v-model="row.status"
              :active-value="true"
              :inactive-value="false"
              active-text="启用"
              inactive-text="禁用"
              @change="toggleFunctionStatus(row.id, !row.status)"
            />
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="createTime" label="创建时间" width="120" />
        
        <ElTableColumn prop="updateTime" label="更新时间" width="120" />
        
        <ElTableColumn label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <ElButton type="primary" size="small" @click="openEditDialog(row)" class="edit-btn">编辑</ElButton>
            <ElButton type="danger" size="small" @click="deleteFunction(row.id, row.name)" class="delete-btn">删除</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
    </div>

    <!-- 分页 -->
    <div class="pagination-section">
      <ElPagination
        :current-page="currentPage"
        @current-change="currentPage = $event"
        :page-size="pageSize"
        @size-change="pageSize = $event"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="filteredFunctions.length"
      />
    </div>

    <!-- 添加/编辑函数工具对话框 -->
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
        <ElFormItem label="函数名称" prop="name">
          <ElInput v-model="formData.name" placeholder="请输入函数工具名称" />
        </ElFormItem>
        
        <ElFormItem label="描述" prop="description">
          <ElInput v-model="formData.description" placeholder="请输入函数工具描述" type="textarea" :rows="3" />
        </ElFormItem>
        
        <ElFormItem label="类型" prop="type">
          <ElSelect v-model="formData.type" placeholder="请选择函数工具类型">
            <ElOption label="API调用" value="api" />
            <ElOption label="脚本执行" value="script" />
            <ElOption label="自定义函数" value="custom" />
          </ElSelect>
        </ElFormItem>
        
        <ElFormItem label="函数代码" prop="code">
          <ElInput v-model="formData.code" placeholder="请输入函数代码" type="textarea" :rows="8" />
          <div class="code-tip">请确保代码格式正确，符合系统要求。</div>
        </ElFormItem>
        
        <ElFormItem label="是否公开">
          <ElSwitch v-model="formData.isPublic" />
        </ElFormItem>
        
        <ElFormItem label="状态">
          <ElSwitch v-model="formData.status" />
        </ElFormItem>
      </ElForm>
      
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveFunction" :loading="loading">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts">
// 辅助函数
const getTypeLabel = (type: string): string => {
  switch (type) {
    case 'api':
      return 'API调用';
    case 'script':
      return '脚本执行';
    case 'custom':
      return '自定义函数';
    default:
      return '未知';
  }
};

const getTagType = (type: string): 'primary' | 'success' | 'warning' | 'info' => {
  switch (type) {
    case 'api':
      return 'primary';
    case 'script':
      return 'success';
    case 'custom':
      return 'warning';
    default:
      return 'info';
  }
};
</script>

<style scoped>
.function-tool-management {
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

/* 函数工具表格 */
.function-table-section {
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 20px;
  margin-bottom: 20px;
}

.function-table {
  border: 1px solid var(--border-color);
}

.description-text {
  display: -webkit-box;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
}

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

/* 代码提示 */
.code-tip {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 4px;
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
  .function-tool-management {
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
  
  .function-table-section {
    padding: 15px;
    overflow-x: auto;
  }
  
  .pagination-section {
    padding: 15px;
  }
}
</style>