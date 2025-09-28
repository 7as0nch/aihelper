<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUserStore } from '../../stores/user';
import { ElCard, ElRow, ElCol, ElStatistic, ElProgress } from 'element-plus';
import 'element-plus/es/components/card/style/css';
import 'element-plus/es/components/row/style/css';
import 'element-plus/es/components/col/style/css';
import 'element-plus/es/components/statistic/style/css';
import 'element-plus/es/components/progress/style/css';
import 'element-plus/es/components/icon/style/css';
import 'element-plus/es/components/tooltip/style/css';

const userStore = useUserStore();

// 统计数据
const totalUsers = ref(0);
const activeChats = ref(0);
const totalFunctions = ref(0);
const totalWorkflows = ref(0);

// 使用率数据
const aiUsageRate = ref(0);
const storageUsageRate = ref(0);
const systemLoadRate = ref(0);

// 最近活动
const recentActivities = ref([
  { user: '张三', action: '登录系统', time: '5分钟前' },
  { user: '李四', action: '创建新聊天', time: '12分钟前' },
  { user: '王五', action: '更新函数工具', time: '30分钟前' },
  { user: '赵六', action: '创建工作流', time: '1小时前' },
  { user: '管理员', action: '系统维护', time: '2小时前' }
]);

// 初始化数据
onMounted(() => {
  // 模拟获取统计数据
  simulateFetchData();
});

// 模拟获取数据
const simulateFetchData = () => {
  // 使用setTimeout模拟异步请求
  setTimeout(() => {
    totalUsers.value = Math.floor(Math.random() * 500) + 100;
    activeChats.value = Math.floor(Math.random() * 100) + 20;
    totalFunctions.value = Math.floor(Math.random() * 50) + 5;
    totalWorkflows.value = Math.floor(Math.random() * 30) + 2;
    
    aiUsageRate.value = Math.floor(Math.random() * 80) + 10;
    storageUsageRate.value = Math.floor(Math.random() * 70) + 20;
    systemLoadRate.value = Math.floor(Math.random() * 60) + 15;
  }, 500);
};
</script>

<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="dashboard-title">仪表盘</h1>
      <p class="dashboard-subtitle">欢迎回来，{{ userStore.userInfo?.username || '管理员' }}！以下是系统的最新状态概览。</p>
    </div>

    <!-- 统计卡片 -->
    <ElRow :gutter="20" class="stats-row">
      <ElCol :xs="24" :sm="12" :md="6" class="stat-col">
        <ElCard class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" 
                height="24" 
                fill="#409EFF" 
                viewBox="0 0 16 16"
              >
                <path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z"/>
              </svg>
            </div>
            <div class="stat-details">
              <ElStatistic :value="totalUsers" :precision="0" class="stat-value" />
              <p class="stat-label">总用户数</p>
            </div>
          </div>
        </ElCard>
      </ElCol>

      <ElCol :xs="24" :sm="12" :md="6" class="stat-col">
        <ElCard class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" 
                height="24" 
                fill="#67C23A" 
                viewBox="0 0 16 16"
              >
                <path d="M2 5a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5zm3.293 1.293a1 1 0 0 1 1.414 0L8 8.586l1.293-1.293a1 1 0 1 1 1.414 1.414L9.414 10l1.293 1.293a1 1 0 0 1-1.414 1.414L8 11.414l-1.293 1.293a1 1 0 0 1-1.414-1.414L6.586 10 5.293 8.707a1 1 0 0 1 0-1.414z"/>
              </svg>
            </div>
            <div class="stat-details">
              <ElStatistic :value="activeChats" :precision="0" class="stat-value" />
              <p class="stat-label">活跃聊天</p>
            </div>
          </div>
        </ElCard>
      </ElCol>

      <ElCol :xs="24" :sm="12" :md="6" class="stat-col">
        <ElCard class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" 
                height="24" 
                fill="#E6A23C" 
                viewBox="0 0 16 16"
              >
                <path d="M5.854 4.646a.5.5 0 0 1 0 .708L2.707 8l3.147 3.146a.5.5 0 0 1-.708.708l-3.5-3.5a.5.5 0 0 1 0-.708l3.5-3.5a.5.5 0 0 1 .708 0zm4.292 0a.5.5 0 0 0 0 .708L13.293 8l-3.147 3.146a.5.5 0 0 0 .708.708l3.5-3.5a.5.5 0 0 0 0-.708l-3.5-3.5a.5.5 0 0 0-.708 0z"/>
              </svg>
            </div>
            <div class="stat-details">
              <ElStatistic :value="totalFunctions" :precision="0" class="stat-value" />
              <p class="stat-label">函数工具</p>
            </div>
          </div>
        </ElCard>
      </ElCol>

      <ElCol :xs="24" :sm="12" :md="6" class="stat-col">
        <ElCard class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" 
                height="24" 
                fill="#F56C6C" 
                viewBox="0 0 16 16"
              >
                <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
                <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
              </svg>
            </div>
            <div class="stat-details">
              <ElStatistic :value="totalWorkflows" :precision="0" class="stat-value" />
              <p class="stat-label">工作流</p>
            </div>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>

    <!-- 使用率和活动记录 -->
    <ElRow :gutter="20" class="dashboard-content">
      <!-- 使用率图表 -->
      <ElCol :xs="24" :md="12" class="usage-col">
        <ElCard class="usage-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">系统使用率</span>
            </div>
          </template>
          
          <div class="usage-content">
            <div class="usage-item">
              <div class="usage-label">
                <span>AI模型使用率</span>
                <span>{{ aiUsageRate }}%</span>
              </div>
              <ElProgress :percentage="aiUsageRate" :color="aiUsageRate > 80 ? '#F56C6C' : '#409EFF'" />
            </div>
            
            <div class="usage-item">
              <div class="usage-label">
                <span>存储空间使用率</span>
                <span>{{ storageUsageRate }}%</span>
              </div>
              <ElProgress :percentage="storageUsageRate" :color="storageUsageRate > 80 ? '#F56C6C' : '#67C23A'" />
            </div>
            
            <div class="usage-item">
              <div class="usage-label">
                <span>系统负载</span>
                <span>{{ systemLoadRate }}%</span>
              </div>
              <ElProgress :percentage="systemLoadRate" :color="systemLoadRate > 80 ? '#F56C6C' : '#E6A23C'" />
            </div>
          </div>
        </ElCard>
      </ElCol>

      <!-- 最近活动 -->
      <ElCol :xs="24" :md="12" class="activities-col">
        <ElCard class="activities-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近活动</span>
            </div>
          </template>
          
          <div class="activities-content">
            <div v-for="(activity, index) in recentActivities" :key="index" class="activity-item">
              <div class="activity-icon">
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  width="16" 
                  height="16" 
                  fill="#409EFF" 
                  viewBox="0 0 16 16"
                >
                  <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
                  <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
                </svg>
              </div>
              <div class="activity-details">
                <p class="activity-text"><span class="activity-user">{{ activity.user }}</span> {{ activity.action }}</p>
                <p class="activity-time">{{ activity.time }}</p>
              </div>
            </div>
          </div>
        </ElCard>
      </ElCol>
    </ElRow>

    <!-- 系统健康状态 -->
    <div class="health-section">
      <ElCard class="health-card">
        <template #header>
          <div class="card-header">
            <span class="card-title">系统健康状态</span>
          </div>
        </template>
        
        <div class="health-content">
          <div class="health-indicator">
            <div class="health-status">
              <div class="status-dot healthy"></div>
              <span>服务正常运行</span>
            </div>
            <div class="health-message">
              所有服务组件运行正常，系统性能稳定。
            </div>
          </div>
          
          <div class="health-components">
            <div class="component-item">
              <span class="component-name">API服务</span>
              <span class="component-status">正常</span>
            </div>
            <div class="component-item">
              <span class="component-name">数据库</span>
              <span class="component-status">正常</span>
            </div>
            <div class="component-item">
              <span class="component-name">AI模型</span>
              <span class="component-status">正常</span>
            </div>
            <div class="component-item">
              <span class="component-name">文件存储</span>
              <span class="component-status">正常</span>
            </div>
          </div>
        </div>
      </ElCard>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 20px 0;
}

.dashboard-header {
  margin-bottom: 30px;
}

.dashboard-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 10px 0;
}

.dashboard-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  transition: all 0.3s ease;
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.15);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  background-color: var(--bg-secondary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  flex-shrink: 0;
}

.stat-details {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 4px 0 0 0;
}

/* 内容区域 */
.dashboard-content {
  margin-bottom: 30px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 使用率卡片 */
.usage-card,
.activities-card,
.health-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border: none;
  height: 100%;
}

.usage-content {
  padding-top: 20px;
}

.usage-item {
  margin-bottom: 20px;
}

.usage-item:last-child {
  margin-bottom: 0;
}

.usage-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-secondary);
}

/* 活动记录卡片 */
.activities-content {
  max-height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  margin-right: 12px;
  margin-top: 2px;
}

.activity-details {
  flex: 1;
}

.activity-text {
  font-size: 14px;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.activity-user {
  font-weight: 500;
  color: #409EFF;
}

.activity-time {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0;
}

/* 健康状态卡片 */
.health-section {
  margin-top: 20px;
}

.health-content {
  padding-top: 20px;
}

.health-indicator {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background-color: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.health-status {
  display: flex;
  align-items: center;
  margin-right: 20px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 8px;
}

.status-dot.healthy {
  background-color: #67C23A;
}

.health-message {
  font-size: 14px;
  color: var(--text-secondary);
}

.health-components {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.component-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  background-color: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.component-name {
  font-size: 14px;
  color: var(--text-secondary);
}

.component-status {
  font-size: 14px;
  font-weight: 500;
  color: #67C23A;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .dashboard {
    padding: 15px 0;
  }
  
  .dashboard-header {
    margin-bottom: 20px;
  }
  
  .stats-row {
    margin-bottom: 20px;
  }
  
  .dashboard-content {
    margin-bottom: 20px;
  }
  
  .health-components {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-title {
    font-size: 20px;
  }
  
  .stat-content {
    padding: 15px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .usage-item,
  .activity-item {
    margin-bottom: 15px;
  }
}
</style>