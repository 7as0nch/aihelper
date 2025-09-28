<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElTabs, ElTabPane, ElForm, ElFormItem, ElInput, ElInputNumber, ElSelect, ElOption, ElSwitch, ElColorPicker, ElButton, ElMessage, ElAlert, ElUpload } from 'element-plus';
import type { UploadProps } from 'element-plus';
import 'element-plus/es/components/tabs/style/css';
import 'element-plus/es/components/tab-pane/style/css';
import 'element-plus/es/components/form/style/css';
import 'element-plus/es/components/form-item/style/css';
import 'element-plus/es/components/input/style/css';
import 'element-plus/es/components/input-number/style/css';
import 'element-plus/es/components/select/style/css';
import 'element-plus/es/components/switch/style/css';
import 'element-plus/es/components/color-picker/style/css';
import 'element-plus/es/components/button/style/css';
import 'element-plus/es/components/message/style/css';
import 'element-plus/es/components/alert/style/css';
import 'element-plus/es/components/upload/style/css';

// 系统基本设置表单数据
const basicSettingsForm = ref({
  systemName: 'AI Chat Assistant',
  systemDescription: '基于人工智能的智能对话助手平台',
  logoUrl: '',
  faviconUrl: '',
  primaryColor: '#409eff',
  enableRegistration: true,
  defaultLanguage: 'zh-CN',
  pageSize: 10,
  cacheDuration: 30,
  apiTimeout: 30,
});

// API设置表单数据
const apiSettingsForm = ref({
  openaiApiKey: '',
  openaiBaseUrl: 'https://api.openai.com/v1',
  openaiModel: 'gpt-3.5-turbo',
  temperature: 0.7,
  maxTokens: 2000,
  enableRateLimit: true,
  rateLimitRequests: 100,
  rateLimitDuration: 60,
  enableLogging: true,
});

// 通知设置表单数据
const notificationSettingsForm = ref({
  enableEmailNotifications: true,
  enableSmsNotifications: false,
  enablePushNotifications: true,
  adminEmail: '',
  supportEmail: '',
  emailTemplateSignoff: 'AI Chat Assistant 团队',
  enableErrorAlerts: true,
  enablePerformanceAlerts: true,
});

// 高级设置表单数据
const advancedSettingsForm = ref({
  enableDebugMode: false,
  enableCors: true,
  enableSsl: false,
  sslCertPath: '',
  sslKeyPath: '',
  sessionTimeout: 30,
  maxFileSize: 10,
  allowedFileTypes: 'jpg,jpeg,png,pdf,doc,docx',
  enableAnalytics: true,
});

// 表单验证规则
const basicSettingsRules = ref({
  systemName: [
    { required: true, message: '请输入系统名称', trigger: 'blur' },
    { min: 2, max: 50, message: '系统名称长度在 2 到 50 个字符之间', trigger: 'blur' }
  ],
  systemDescription: [
    { required: true, message: '请输入系统描述', trigger: 'blur' },
    { min: 5, max: 200, message: '系统描述长度在 5 到 200 个字符之间', trigger: 'blur' }
  ],
  defaultLanguage: [
    { required: true, message: '请选择默认语言', trigger: 'change' }
  ],
  pageSize: [
    { required: true, message: '请输入默认分页大小', trigger: 'blur' },
    { type: 'number' as const, min: 1, max: 100, message: '分页大小在 1 到 100 之间', trigger: 'blur' }
  ],
});

const apiSettingsRules = ref({
  openaiBaseUrl: [
    { required: true, message: '请输入OpenAI API地址', trigger: 'blur' }
  ],
  openaiModel: [
    { required: true, message: '请选择OpenAI模型', trigger: 'change' }
  ],
  temperature: [
    { required: true, message: '请输入温度值', trigger: 'blur' },
    { type: 'number' as const, min: 0, max: 2, message: '温度值在 0 到 2 之间', trigger: 'blur' }
  ],
  maxTokens: [
    { required: true, message: '请输入最大令牌数', trigger: 'blur' },
    { type: 'number' as const, min: 1, max: 10000, message: '最大令牌数在 1 到 10000 之间', trigger: 'blur' }
  ],
});

const notificationSettingsRules = ref({
  adminEmail: [
    { required: true, message: '请输入管理员邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  supportEmail: [
    { required: true, message: '请输入支持邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
});

// 表单提交状态
const savingSettings = ref(false);
const settingsSaved = ref(false);

// 加载设置数据
onMounted(() => {
  loadSettings();
});

// 模拟加载设置数据
const loadSettings = () => {
  // 模拟API请求延迟
  setTimeout(() => {
    // 这里应该是从API获取设置数据的逻辑
    // 目前使用预设的模拟数据
  }, 500);
};

// 保存系统基本设置
const saveBasicSettings = () => {
  savingSettings.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    savingSettings.value = false;
    settingsSaved.value = true;
    ElMessage.success('系统基本设置保存成功');
    
    // 5秒后隐藏保存成功提示
    setTimeout(() => {
      settingsSaved.value = false;
    }, 5000);
  }, 800);
};

// 保存API设置
const saveApiSettings = () => {
  savingSettings.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    savingSettings.value = false;
    settingsSaved.value = true;
    ElMessage.success('API设置保存成功');
    
    // 5秒后隐藏保存成功提示
    setTimeout(() => {
      settingsSaved.value = false;
    }, 5000);
  }, 800);
};

// 保存通知设置
const saveNotificationSettings = () => {
  savingSettings.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    savingSettings.value = false;
    settingsSaved.value = true;
    ElMessage.success('通知设置保存成功');
    
    // 5秒后隐藏保存成功提示
    setTimeout(() => {
      settingsSaved.value = false;
    }, 5000);
  }, 800);
};

// 保存高级设置
const saveAdvancedSettings = () => {
  savingSettings.value = true;
  
  // 模拟API请求延迟
  setTimeout(() => {
    savingSettings.value = false;
    settingsSaved.value = true;
    ElMessage.success('高级设置保存成功');
    
    // 5秒后隐藏保存成功提示
    setTimeout(() => {
      settingsSaved.value = false;
    }, 5000);
  }, 800);
};

// 文件上传配置
const uploadOptions: Partial<UploadProps> = {
  headers: {
    'Authorization': 'Bearer token'
  },
  showFileList: false,
  beforeUpload: (file) => {
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      ElMessage.error('文件大小不能超过 2MB!');
      return false;
    }
    return true;
  },
  onSuccess: (_response, _uploadFile) => {
    // 这里应该处理上传成功后的逻辑
    ElMessage.success('文件上传成功');
  },
  onError: (_error, _uploadFile) => {
    ElMessage.error('文件上传失败');
  }
};
</script>

<template>
  <div class="system-settings">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <p class="page-description">配置系统的各项参数和选项，包括基本信息、API设置、通知设置和高级设置。</p>
    </div>

    <!-- 保存成功提示 -->
    <ElAlert
      v-if="settingsSaved"
      title="保存成功"
      type="success"
      :closable="false"
      style="margin-bottom: 20px;"
    />

    <!-- 设置标签页 -->
    <ElTabs v-model="activeTab" type="card" class="settings-tabs">
      <!-- 基本设置 -->
      <ElTabPane label="基本设置" name="basic">
        <ElForm
          ref="basicSettingsRef"
          :model="basicSettingsForm"
          :rules="basicSettingsRules"
          label-width="150px"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">系统信息</h3>
            
            <ElFormItem label="系统名称" prop="systemName">
              <ElInput v-model="basicSettingsForm.systemName" placeholder="请输入系统名称" />
            </ElFormItem>
            
            <ElFormItem label="系统描述" prop="systemDescription">
              <ElInput v-model="basicSettingsForm.systemDescription" placeholder="请输入系统描述" type="textarea" :rows="3" />
            </ElFormItem>
            
            <ElFormItem label="系统Logo">
              <ElUpload
                v-bind="uploadOptions"
                action="/api/upload/logo"
                class="avatar-uploader"
                :show-file-list="false"
              >
                <img v-if="basicSettingsForm.logoUrl" :src="basicSettingsForm.logoUrl" class="avatar" />
                <div v-else class="avatar-placeholder">
                  <i class="el-icon-plus avatar-uploader-icon"></i>
                  <span class="upload-text">上传Logo</span>
                </div>
              </ElUpload>
              <div class="upload-hint">建议尺寸: 200x60px, 最大2MB</div>
            </ElFormItem>
            
            <ElFormItem label="网站图标">
              <ElUpload
                v-bind="uploadOptions"
                action="/api/upload/favicon"
                class="avatar-uploader"
                :show-file-list="false"
              >
                <img v-if="basicSettingsForm.faviconUrl" :src="basicSettingsForm.faviconUrl" class="favicon" />
                <div v-else class="avatar-placeholder">
                  <i class="el-icon-plus avatar-uploader-icon"></i>
                  <span class="upload-text">上传图标</span>
                </div>
              </ElUpload>
              <div class="upload-hint">建议尺寸: 32x32px, 最大1MB</div>
            </ElFormItem>
            
            <ElFormItem label="主题主色调">
              <ElColorPicker v-model="basicSettingsForm.primaryColor" class="color-picker" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">用户设置</h3>
            
            <ElFormItem label="允许用户注册">
              <ElSwitch v-model="basicSettingsForm.enableRegistration" />
            </ElFormItem>
            
            <ElFormItem label="默认语言" prop="defaultLanguage">
              <ElSelect v-model="basicSettingsForm.defaultLanguage" placeholder="请选择默认语言">
                <ElOption label="简体中文" value="zh-CN" />
                <ElOption label="English" value="en-US" />
              </ElSelect>
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">显示设置</h3>
            
            <ElFormItem label="默认分页大小" prop="pageSize">
              <ElInputNumber v-model="basicSettingsForm.pageSize" :min="1" :max="100" />
            </ElFormItem>
            
            <ElFormItem label="缓存时长(分钟)">
              <ElInputNumber v-model="basicSettingsForm.cacheDuration" :min="0" :max="120" />
            </ElFormItem>
            
            <ElFormItem label="API超时(秒)">
              <ElInputNumber v-model="basicSettingsForm.apiTimeout" :min="5" :max="120" />
            </ElFormItem>
          </div>

          <div class="form-actions">
            <ElButton type="primary" @click="saveBasicSettings" :loading="savingSettings">保存设置</ElButton>
          </div>
        </ElForm>
      </ElTabPane>

      <!-- API设置 -->
      <ElTabPane label="API设置" name="api">
        <ElForm
          ref="apiSettingsRef"
          :model="apiSettingsForm"
          :rules="apiSettingsRules"
          label-width="150px"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">OpenAI API</h3>
            
            <ElFormItem label="API Key">
              <ElInput v-model="apiSettingsForm.openaiApiKey" placeholder="请输入OpenAI API Key" show-password />
            </ElFormItem>
            
            <ElFormItem label="API地址" prop="openaiBaseUrl">
              <ElInput v-model="apiSettingsForm.openaiBaseUrl" placeholder="请输入OpenAI API地址" />
            </ElFormItem>
            
            <ElFormItem label="默认模型" prop="openaiModel">
              <ElSelect v-model="apiSettingsForm.openaiModel" placeholder="请选择默认模型">
                <ElOption label="gpt-3.5-turbo" value="gpt-3.5-turbo" />
                <ElOption label="gpt-3.5-turbo-16k" value="gpt-3.5-turbo-16k" />
                <ElOption label="gpt-4" value="gpt-4" />
                <ElOption label="gpt-4-32k" value="gpt-4-32k" />
              </ElSelect>
            </ElFormItem>
            
            <ElFormItem label="温度值" prop="temperature">
              <ElInputNumber v-model="apiSettingsForm.temperature" :min="0" :max="2" :step="0.1" />
              <div class="setting-hint">控制输出的随机性，值越高越随机</div>
            </ElFormItem>
            
            <ElFormItem label="最大令牌数" prop="maxTokens">
              <ElInputNumber v-model="apiSettingsForm.maxTokens" :min="1" :max="10000" />
              <div class="setting-hint">控制模型响应的最大长度</div>
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">速率限制</h3>
            
            <ElFormItem label="启用速率限制">
              <ElSwitch v-model="apiSettingsForm.enableRateLimit" />
            </ElFormItem>
            
            <ElFormItem label="最大请求数" v-if="apiSettingsForm.enableRateLimit">
              <ElInputNumber v-model="apiSettingsForm.rateLimitRequests" :min="1" :max="1000" />
            </ElFormItem>
            
            <ElFormItem label="时间窗口(秒)" v-if="apiSettingsForm.enableRateLimit">
              <ElInputNumber v-model="apiSettingsForm.rateLimitDuration" :min="1" :max="3600" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">日志记录</h3>
            
            <ElFormItem label="启用API日志">
              <ElSwitch v-model="apiSettingsForm.enableLogging" />
            </ElFormItem>
          </div>

          <div class="form-actions">
            <ElButton type="primary" @click="saveApiSettings" :loading="savingSettings">保存设置</ElButton>
          </div>
        </ElForm>
      </ElTabPane>

      <!-- 通知设置 -->
      <ElTabPane label="通知设置" name="notification">
        <ElForm
          ref="notificationSettingsRef"
          :model="notificationSettingsForm"
          :rules="notificationSettingsRules"
          label-width="150px"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">通知方式</h3>
            
            <ElFormItem label="启用邮件通知">
              <ElSwitch v-model="notificationSettingsForm.enableEmailNotifications" />
            </ElFormItem>
            
            <ElFormItem label="启用短信通知">
              <ElSwitch v-model="notificationSettingsForm.enableSmsNotifications" />
            </ElFormItem>
            
            <ElFormItem label="启用推送通知">
              <ElSwitch v-model="notificationSettingsForm.enablePushNotifications" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">联系信息</h3>
            
            <ElFormItem label="管理员邮箱" prop="adminEmail">
              <ElInput v-model="notificationSettingsForm.adminEmail" placeholder="请输入管理员邮箱" />
            </ElFormItem>
            
            <ElFormItem label="支持邮箱" prop="supportEmail">
              <ElInput v-model="notificationSettingsForm.supportEmail" placeholder="请输入支持邮箱" />
            </ElFormItem>
            
            <ElFormItem label="邮件模板签名">
              <ElInput v-model="notificationSettingsForm.emailTemplateSignoff" placeholder="请输入邮件模板签名" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">告警设置</h3>
            
            <ElFormItem label="启用错误告警">
              <ElSwitch v-model="notificationSettingsForm.enableErrorAlerts" />
            </ElFormItem>
            
            <ElFormItem label="启用性能告警">
              <ElSwitch v-model="notificationSettingsForm.enablePerformanceAlerts" />
            </ElFormItem>
          </div>

          <div class="form-actions">
            <ElButton type="primary" @click="saveNotificationSettings" :loading="savingSettings">保存设置</ElButton>
          </div>
        </ElForm>
      </ElTabPane>

      <!-- 高级设置 -->
      <ElTabPane label="高级设置" name="advanced">
        <ElForm
          ref="advancedSettingsRef"
          :model="advancedSettingsForm"
          label-width="150px"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">系统调试</h3>
            
            <ElFormItem label="启用调试模式">
              <ElSwitch v-model="advancedSettingsForm.enableDebugMode" />
              <div class="setting-hint">启用后将显示更多调试信息</div>
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">安全设置</h3>
            
            <ElFormItem label="启用CORS">
              <ElSwitch v-model="advancedSettingsForm.enableCors" />
            </ElFormItem>
            
            <ElFormItem label="启用SSL">
              <ElSwitch v-model="advancedSettingsForm.enableSsl" />
            </ElFormItem>
            
            <ElFormItem label="SSL证书路径" v-if="advancedSettingsForm.enableSsl">
              <ElInput v-model="advancedSettingsForm.sslCertPath" placeholder="请输入SSL证书路径" />
            </ElFormItem>
            
            <ElFormItem label="SSL私钥路径" v-if="advancedSettingsForm.enableSsl">
              <ElInput v-model="advancedSettingsForm.sslKeyPath" placeholder="请输入SSL私钥路径" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">会话设置</h3>
            
            <ElFormItem label="会话超时(分钟)">
              <ElInputNumber v-model="advancedSettingsForm.sessionTimeout" :min="1" :max="120" />
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">文件上传</h3>
            
            <ElFormItem label="最大文件大小(MB)">
              <ElInputNumber v-model="advancedSettingsForm.maxFileSize" :min="1" :max="100" />
            </ElFormItem>
            
            <ElFormItem label="允许的文件类型">
              <ElInput v-model="advancedSettingsForm.allowedFileTypes" placeholder="例如: jpg,jpeg,png,pdf,doc,docx" />
              <div class="setting-hint">多个文件类型用逗号分隔</div>
            </ElFormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">数据分析</h3>
            
            <ElFormItem label="启用数据分析">
              <ElSwitch v-model="advancedSettingsForm.enableAnalytics" />
            </ElFormItem>
          </div>

          <div class="form-actions">
            <ElButton type="primary" @click="saveAdvancedSettings" :loading="savingSettings">保存设置</ElButton>
          </div>
        </ElForm>
      </ElTabPane>
    </ElTabs>
  </div>
</template>

<script lang="ts">

// 激活的标签页
const activeTab = ref('basic');
</script>

<style scoped>
.system-settings {
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

/* 设置标签页 */
.settings-tabs {
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* 设置表单 */
.settings-form {
  padding: 20px;
}

.form-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 20px;
  padding-bottom: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 15px 0;
  padding-left: 10px;
  border-left: 4px solid var(--primary-color);
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  justify-content: flex-start;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

/* 文件上传样式 */
.avatar-uploader {
  display: flex;
  align-items: center;
}

.avatar {
  width: 200px;
  height: 60px;
  border-radius: var(--radius-sm);
  object-fit: contain;
  border: 1px solid var(--border-color);
  background-color: var(--bg-secondary);
}

.favicon {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-sm);
  object-fit: contain;
  border: 1px solid var(--border-color);
  background-color: var(--bg-secondary);
}

.avatar-placeholder {
  width: 200px;
  height: 60px;
  border: 1px dashed var(--border-color);
  border-radius: var(--radius-sm);
  background-color: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s;
}

.avatar-placeholder:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.avatar-uploader-icon {
  font-size: 24px;
  color: var(--text-tertiary);
}

.upload-text {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 5px;
}

.upload-hint {
  margin-top: 5px;
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 设置提示 */
.setting-hint {
  margin-top: 5px;
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 颜色选择器 */
.color-picker {
  width: 100px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .settings-form {
    padding: 15px;
  }
  
  .form-section {
    margin-bottom: 25px;
  }
}

@media (max-width: 768px) {
  .system-settings {
    padding: 15px 0;
  }
  
  .page-header {
    margin-bottom: 20px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .settings-form {
    padding: 10px;
  }
  
  .form-section {
    margin-bottom: 20px;
    padding-bottom: 15px;
  }
  
  .section-title {
    font-size: 14px;
  }
  
  .avatar {
    width: 150px;
    height: 50px;
  }
  
  .avatar-placeholder {
    width: 150px;
    height: 50px;
  }
}
</style>