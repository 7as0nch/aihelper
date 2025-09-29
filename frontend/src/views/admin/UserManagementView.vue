<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ElTable, ElTableColumn, ElPagination, ElInput, ElSelect, ElOption, ElButton, ElDialog, ElForm, ElFormItem, ElSwitch, ElMessageBox, ElMessage } from 'element-plus';
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

// 模拟用户数据
interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'user';
  status: boolean;
  createTime: string;
  lastLoginTime: string;
  avatar?: string;
}

const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searchKeyword = ref('');
const roleFilter = ref('all');
const statusFilter = ref('all');
const userList = ref<User[]>([]);
const loading = ref(false);

// 对话框状态
const dialogVisible = ref(false);
const dialogTitle = ref('添加用户');
const isEditMode = ref(false);
const currentUserId = ref<number | null>(null);

// 表单数据
const formData = ref({
  username: '',
  email: '',
  password: '',
  role: 'user' as 'admin' | 'user',
  status: true
});

// 表单验证
const formRules = ref({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: !isEditMode.value, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 30, message: '密码长度在 6 到 30 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
});

// 筛选后的用户列表
const filteredUsers = computed(() => {
  return userList.value.filter(user => {
    // 关键词筛选
    const keywordMatch = !searchKeyword.value || 
      user.username.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      user.email.toLowerCase().includes(searchKeyword.value.toLowerCase());
    
    // 角色筛选
    const roleMatch = roleFilter.value === 'all' || user.role === roleFilter.value;
    
    // 状态筛选
    const statusMatch = statusFilter.value === 'all' || user.status === (statusFilter.value === 'active');
    
    return keywordMatch && roleMatch && statusMatch;
  });
});

// 分页后的用户列表
const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredUsers.value.slice(start, end);
});

// 初始化数据
onMounted(() => {
  fetchUsers();
});

// 模拟获取用户数据
const fetchUsers = () => {
  loading.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    const mockUsers: User[] = Array.from({ length: 50 }, (_, i) => ({
      id: i + 1,
      username: `user${i + 1}`,
      email: `user${i + 1}@example.com`,
      role: i < 5 ? 'admin' : 'user',
      status: Math.random() > 0.2,
      createTime: new Date(Date.now() - Math.floor(Math.random() * 30) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN'),
      lastLoginTime: new Date(Date.now() - Math.floor(Math.random() * 7) * 24 * 60 * 60 * 1000).toLocaleDateString('zh-CN')
    }));
    
    userList.value = mockUsers;
    total.value = mockUsers.length;
    loading.value = false;
  }, 800);
};

// 搜索用户
const handleSearch = () => {
  currentPage.value = 1;
};

// 重置筛选条件
const handleReset = () => {
  searchKeyword.value = '';
  roleFilter.value = 'all';
  statusFilter.value = 'all';
  currentPage.value = 1;
};

// 打开添加用户对话框
const openAddDialog = () => {
  dialogTitle.value = '添加用户';
  isEditMode.value = false;
  currentUserId.value = null;
  formData.value = {
    username: '',
    email: '',
    password: '',
    role: 'user',
    status: true
  };
  dialogVisible.value = true;
};

// 打开编辑用户对话框
const openEditDialog = (user: User) => {
  dialogTitle.value = '编辑用户';
  isEditMode.value = true;
  currentUserId.value = user.id;
  formData.value = {
    username: user.username,
    email: user.email,
    password: '', // 编辑时不显示密码
    role: user.role,
    status: user.status
  };
  dialogVisible.value = true;
};

// 保存用户（添加或编辑）
const saveUser = () => {
  // 模拟API请求延迟
  loading.value = true;
  
  setTimeout(() => {
    if (isEditMode.value && currentUserId.value !== null) {
      // 编辑用户
      const index = userList.value.findIndex(user => user.id === currentUserId.value);
      if (index !== -1 && userList.value[index]) {
        const user = userList.value[index];
        userList.value[index] = {
          id: user.id,
          username: formData.value.username,
          email: formData.value.email,
          role: formData.value.role,
          status: formData.value.status,
          createTime: user.createTime,
          lastLoginTime: user.lastLoginTime,
          avatar: user.avatar,
          ...(formData.value.password && { password: formData.value.password }) // 只有在修改了密码时才更新
        };
      }
      ElMessage.success('用户更新成功');
    } else {
      // 添加用户
      const newUser: User = {
        id: userList.value.length + 1,
        username: formData.value.username,
        email: formData.value.email,
        role: formData.value.role,
        status: formData.value.status,
        createTime: new Date().toLocaleDateString('zh-CN'),
        lastLoginTime: '--',
        avatar: ''
      };
      userList.value.unshift(newUser);
      total.value++;
      ElMessage.success('用户添加成功');
    }
    
    dialogVisible.value = false;
    loading.value = false;
  }, 800);
};

// 删除用户
const deleteUser = (userId: number, username: string) => {
  ElMessageBox.confirm(
    `确定要删除用户「${username}」吗？此操作不可撤销。`,
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
      const index = userList.value.findIndex(user => user.id === userId);
      if (index !== -1) {
        userList.value.splice(index, 1);
        total.value--;
        // 如果删除后当前页没有数据，则返回上一页
        if (paginatedUsers.value.length === 0 && currentPage.value > 1) {
          currentPage.value--;
        }
      }
      ElMessage.success('用户删除成功');
      loading.value = false;
    }, 500);
  })
  .catch(() => {
    // 用户取消删除
  });
};

// 切换用户状态
const toggleUserStatus = (userId: number, currentStatus: boolean) => {
  loading.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    const index = userList.value.findIndex(user => user.id === userId);
    if (index !== -1 && userList.value[index]) {
      userList.value[index].status = !currentStatus;
    }
    ElMessage.success(`用户状态已${!currentStatus ? '启用' : '禁用'}`);
    loading.value = false;
  }, 300);
};
</script>

<template>
  <div class="user-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <p class="page-description">管理系统中的所有用户，包括添加、编辑、删除和查看用户信息。</p>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="search-filter-section">
      <div class="search-filter-left">
        <ElInput
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          prefix-icon="Search"
          clearable
          class="search-input"
          @keyup.enter="handleSearch"
        />
        
        <ElSelect
          v-model="roleFilter"
          placeholder="选择角色"
          class="filter-select"
          @change="handleSearch"
        >
          <ElOption label="全部" value="all" />
          <ElOption label="管理员" value="admin" />
          <ElOption label="普通用户" value="user" />
        </ElSelect>
        
        <ElSelect
          v-model="statusFilter"
          placeholder="选择状态"
          class="filter-select"
          @change="handleSearch"
        >
          <ElOption label="全部" value="all" />
          <ElOption label="活跃" value="active" />
          <ElOption label="禁用" value="inactive" />
        </ElSelect>
        
        <ElButton type="primary" @click="handleSearch" class="search-btn">搜索</ElButton>
        <ElButton @click="handleReset" class="reset-btn">重置</ElButton>
      </div>
      
      <div class="search-filter-right">
        <ElButton type="primary" @click="openAddDialog" class="add-btn">添加用户</ElButton>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="user-table-section">
      <ElTable
        v-loading="loading"
        :data="paginatedUsers"
        style="width: 100%"
        border
        class="user-table"
      >
        <ElTableColumn prop="id" label="用户ID" width="80" align="center" />
        
        <ElTableColumn prop="username" label="用户名" min-width="120">
          <template #default="{ row }">
            <div class="user-info">
              <div class="user-avatar">
                <span>{{ row.username.charAt(0).toUpperCase() }}</span>
              </div>
              <span class="username-text">{{ row.username }}</span>
            </div>
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="email" label="邮箱" min-width="180" />
        
        <ElTableColumn prop="role" label="角色" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'primary' : 'success'">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
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
              @change="toggleUserStatus(row.id, !row.status)"
            />
          </template>
        </ElTableColumn>
        
        <ElTableColumn prop="createTime" label="创建时间" width="120" />
        
        <ElTableColumn prop="lastLoginTime" label="最后登录" width="120" />
        
        <ElTableColumn label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <ElButton type="primary" size="small" @click="openEditDialog(row)" class="edit-btn">编辑</ElButton>
            <ElButton type="danger" size="small" @click="deleteUser(row.id, row.username)" class="delete-btn">删除</ElButton>
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
        :total="filteredUsers.length"
      />
    </div>

    <!-- 添加/编辑用户对话框 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :before-close="() => dialogVisible = false"
    >
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <ElFormItem label="用户名" prop="username">
          <ElInput v-model="formData.username" placeholder="请输入用户名" />
        </ElFormItem>
        
        <ElFormItem label="邮箱" prop="email">
          <ElInput v-model="formData.email" placeholder="请输入邮箱" type="email" />
        </ElFormItem>
        
        <ElFormItem label="密码" prop="password">
          <ElInput v-model="formData.password" placeholder="请输入密码" type="password" />
          <div v-if="isEditMode" class="password-tip">不修改密码请保持为空</div>
        </ElFormItem>
        
        <ElFormItem label="角色" prop="role">
          <ElSelect v-model="formData.role" placeholder="请选择角色">
            <ElOption label="管理员" value="admin" />
            <ElOption label="普通用户" value="user" />
          </ElSelect>
        </ElFormItem>
        
        <ElFormItem label="状态" prop="status">
          <ElSwitch v-model="formData.status" />
        </ElFormItem>
      </ElForm>
      
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="saveUser" :loading="loading">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<style scoped>
.user-management {
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

/* 用户表格 */
.user-table-section {
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 20px;
  margin-bottom: 20px;
}

.user-table {
  border: 1px solid var(--border-color);
}

.user-info {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background-color: #409EFF;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 500;
  margin-right: 8px;
}

.username-text {
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

/* 密码提示 */
.password-tip {
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
  .user-management {
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
  
  .user-table-section {
    padding: 15px;
    overflow-x: auto;
  }
  
  .pagination-section {
    padding: 15px;
  }
}
</style>