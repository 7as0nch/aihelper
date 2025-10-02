<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { FormInstance } from 'ant-design-vue';
import { Tabs, TabPane, Form, FormItem, Input, InputNumber, Select, Switch, Button, Alert, Upload } from 'ant-design-vue';
import { PlusOutlined } from '@ant-design/icons';
import type { UploadProps } from 'ant-design-vue';
import { message } from 'ant-design-vue';

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
    { type: 'number', min: 1, max: 100, message: '分页大小在 1 到 100 之间', trigger: 'blur' }
  ],
}) as any;

const apiSettingsRules = ref({
  openaiBaseUrl: [
    { required: true, message: '请输入OpenAI API地址', trigger: 'blur' }
  ],
  openaiModel: [
    { required: true, message: '请选择OpenAI模型', trigger: 'change' }
  ],
  temperature: [
    { required: true, message: '请输入温度值', trigger: 'blur' },
    { type: 'number', min: 0, max: 2, message: '温度值在 0 到 2 之间', trigger: 'blur' }
  ],
  maxTokens: [
    { required: true, message: '请输入最大令牌数', trigger: 'blur' },
    { type: 'number', min: 1, max: 10000, message: '最大令牌数在 1 到 10000 之间', trigger: 'blur' }
  ],
}) as any;

const notificationSettingsRules = ref({
  adminEmail: [
    { required: true, message: '请输入管理员邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  supportEmail: [
    { required: true, message: '请输入支持邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
}) as any;

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

// 表单引用
const basicSettingsRef = ref<FormInstance>();
const apiSettingsRef = ref<FormInstance>();
const notificationSettingsRef = ref<FormInstance>();
const advancedSettingsRef = ref<FormInstance>();

// 模拟保存设置的API函数
const saveSettings = async (_type?: string, _data?: any): Promise<void> => {
  // 模拟API请求延迟
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, 800);
  });
};

// 显示保存成功消息
const showSettingsSavedMessage = () => {
  settingsSaved.value = true;
  
  // 5秒后隐藏保存成功提示
  setTimeout(() => {
    settingsSaved.value = false;
  }, 5000);
};

// 保存基本设置
const saveBasicSettings = async () => {
  try {
    if (basicSettingsRef.value) {
      const valid = await basicSettingsRef.value.validateFields();
      if (valid) {
        savingSettings.value = true;
        // 调用API保存设置
        await saveSettings('basic', basicSettingsForm.value);
        savingSettings.value = false;
        showSettingsSavedMessage();
        message.success('系统基本设置保存成功');
      }
    }
  } catch (error) {
    savingSettings.value = false;
    message.error('保存失败，请检查输入');
    console.error('保存基本设置失败:', error);
  }
};

// 保存API设置
const saveApiSettings = async () => {
  try {
    if (apiSettingsRef.value) {
      const valid = await apiSettingsRef.value.validateFields();
      if (valid) {
        savingSettings.value = true;
        // 调用API保存设置
        await saveSettings('api', apiSettingsForm.value);
        savingSettings.value = false;
        showSettingsSavedMessage();
        message.success('API设置保存成功');
      }
    }
  } catch (error) {
    savingSettings.value = false;
    message.error('保存失败，请检查输入');
    console.error('保存API设置失败:', error);
  }
};

// 保存通知设置
const saveNotificationSettings = async () => {
  try {
    if (notificationSettingsRef.value) {
      const valid = await notificationSettingsRef.value.validateFields();
      if (valid) {
        savingSettings.value = true;
        // 调用API保存设置
        await saveSettings('notification', notificationSettingsForm.value);
        savingSettings.value = false;
        showSettingsSavedMessage();
        message.success('通知设置保存成功');
      }
    }
  } catch (error) {
    savingSettings.value = false;
    message.error('保存失败，请检查输入');
    console.error('保存通知设置失败:', error);
  }
};

// 保存高级设置
const saveAdvancedSettings = async () => {
  try {
    if (advancedSettingsRef.value) {
      const valid = await advancedSettingsRef.value.validateFields();
      if (valid) {
        savingSettings.value = true;
        // 调用API保存设置
        await saveSettings('advanced', advancedSettingsForm.value);
        savingSettings.value = false;
        showSettingsSavedMessage();
        message.success('高级设置保存成功');
      }
    }
  } catch (error) {
    savingSettings.value = false;
    message.error('保存失败，请检查输入');
    console.error('保存高级设置失败:', error);
  }
};

// 文件上传配置
    const uploadOptions: Partial<UploadProps> = {
      name: 'file',
      headers: {
        authorization: 'authorization-text'
      },
      showUploadList: false,
      onChange(info: any) {
        if (info.file.status !== 'uploading') {
          console.log(info.file, info.fileList);
        }
        if (info.file.status === 'done') {
          // 这里假设API返回的数据结构包含fileUrl字段
          if (info.file.response && info.file.response.fileUrl) {
            // 根据上传的action来判断是Logo还是Favicon
            if (info.fileList[0].originFileObj && info.fileList[0].originFileObj.lastModified) {
              // 这里需要根据实际情况判断更新哪个表单字段
              // 为简化示例，我们假设根据文件名或其他方式判断
              const fileUrl = info.file.response.fileUrl;
              if (info.action && info.action.includes('logo')) {
                basicSettingsForm.value.logoUrl = fileUrl;
              } else if (info.action && info.action.includes('favicon')) {
                basicSettingsForm.value.faviconUrl = fileUrl;
              }
            }
            message.success(`${info.file.name} 文件上传成功`);
          }
        } else if (info.file.status === 'error') {
          message.error(`${info.file.name} 文件上传失败`);
        }
      },
      beforeUpload(file: File) {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
        if (!isJpgOrPng) {
          message.error('只能上传 JPG/PNG 格式的图片!');
        }
        const isLt2M = file.size / 1024 / 1024 < 2;
        if (!isLt2M) {
          message.error('图片大小不能超过 2MB!');
        }
        return isJpgOrPng && isLt2M;
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
    <Alert
      v-if="settingsSaved"
      message="保存成功"
      type="success"
      showIcon
      :closable="false"
      style="margin-bottom: 20px;"
    />

    <!-- 设置标签页 -->
    <Tabs v-model="activeTab" type="card" class="settings-tabs">
      <!-- 基本设置 -->
      <TabPane tab="基本设置" key="basic">
        <Form
              ref="basicSettingsRef"
              :model="basicSettingsForm"
              :rules="basicSettingsRules"
              :labelCol="{ span: 6 }"
              :wrapperCol="{ span: 16 }"
              class="settings-form"
            >
          <div class="form-section">
            <h3 class="section-title">系统信息</h3>
            
            <FormItem label="系统名称" name="systemName">
              <Input v-model="basicSettingsForm.systemName" placeholder="请输入系统名称" />
            </FormItem>
            
            <FormItem label="系统描述" name="systemDescription">
              <Input.TextArea v-model="basicSettingsForm.systemDescription" placeholder="请输入系统描述" :rows="3" />
            </FormItem>
            
            <FormItem label="系统Logo">
              <Upload
                v-bind="uploadOptions"
                action="/api/upload/logo"
                class="avatar-uploader"
                :showUploadList="false"
              >
                <img v-if="basicSettingsForm.logoUrl" :src="basicSettingsForm.logoUrl" class="avatar" />
                <div v-else class="avatar-placeholder">
                  <PlusOutlined class="avatar-uploader-icon" />
                  <span class="upload-text">上传Logo</span>
                </div>
              </Upload>
              <div class="upload-hint">建议尺寸: 200x60px, 最大2MB</div>
            </FormItem>
            
            <FormItem label="网站图标">
              <Upload
                v-bind="uploadOptions"
                action="/api/upload/favicon"
                class="avatar-uploader"
                :showUploadList="false"
              >
                <img v-if="basicSettingsForm.faviconUrl" :src="basicSettingsForm.faviconUrl" class="favicon" />
                <div v-else class="avatar-placeholder">
                  <PlusOutlined class="avatar-uploader-icon" />
                  <span class="upload-text">上传图标</span>
                </div>
              </Upload>
              <div class="upload-hint">建议尺寸: 32x32px, 最大1MB</div>
            </FormItem>
            
            <FormItem label="主题主色调">
              <ColorPicker v-model="basicSettingsForm.primaryColor" class="color-picker" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">用户设置</h3>
            
            <FormItem label="允许用户注册">
              <Switch v-model="basicSettingsForm.enableRegistration" />
            </FormItem>
            
            <FormItem label="默认语言" name="defaultLanguage">
              <Select v-model="basicSettingsForm.defaultLanguage" placeholder="请选择默认语言">
                <Select.Option value="zh-CN">简体中文</Select.Option>
                <Select.Option value="en-US">English</Select.Option>
              </Select>
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">显示设置</h3>
            
            <FormItem label="默认分页大小" name="pageSize">
              <InputNumber v-model="basicSettingsForm.pageSize" :min="1" :max="100" />
            </FormItem>
            
            <FormItem label="缓存时长(分钟)">
              <InputNumber v-model="basicSettingsForm.cacheDuration" :min="0" :max="120" />
            </FormItem>
            
            <FormItem label="API超时(秒)">
              <InputNumber v-model="basicSettingsForm.apiTimeout" :min="5" :max="120" />
            </FormItem>
          </div>

          <div class="form-actions">
            <Button type="primary" @click="saveBasicSettings" :loading="savingSettings">保存设置</Button>
          </div>
        </Form>
      </TabPane>

      <!-- API设置 -->
      <TabPane tab="API设置" key="api">
        <Form
          ref="apiSettingsRef"
          :model="apiSettingsForm"
          :rules="apiSettingsRules"
          :labelCol="{ span: 6 }"
          :wrapperCol="{ span: 16 }"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">OpenAI API</h3>
            
            <FormItem label="API Key">
              <Input.Password v-model="apiSettingsForm.openaiApiKey" placeholder="请输入OpenAI API Key" />
            </FormItem>
            
            <FormItem label="API地址" name="openaiBaseUrl">
              <Input v-model="apiSettingsForm.openaiBaseUrl" placeholder="请输入OpenAI API地址" />
            </FormItem>
            
            <FormItem label="默认模型" name="openaiModel">
              <Select v-model="apiSettingsForm.openaiModel" placeholder="请选择默认模型">
                <Select.Option value="gpt-3.5-turbo">gpt-3.5-turbo</Select.Option>
                <Select.Option value="gpt-3.5-turbo-16k">gpt-3.5-turbo-16k</Select.Option>
                <Select.Option value="gpt-4">gpt-4</Select.Option>
                <Select.Option value="gpt-4-32k">gpt-4-32k</Select.Option>
              </Select>
            </FormItem>
            
            <FormItem label="温度值" name="temperature">
              <InputNumber v-model="apiSettingsForm.temperature" :min="0" :max="2" :step="0.1" />
              <div class="setting-hint">控制输出的随机性，值越高越随机</div>
            </FormItem>
            
            <FormItem label="最大令牌数" name="maxTokens">
              <InputNumber v-model="apiSettingsForm.maxTokens" :min="1" :max="10000" />
              <div class="setting-hint">控制模型响应的最大长度</div>
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">速率限制</h3>
            
            <FormItem label="启用速率限制">
              <Switch v-model="apiSettingsForm.enableRateLimit" />
            </FormItem>
            
            <FormItem label="最大请求数" v-if="apiSettingsForm.enableRateLimit">
              <InputNumber v-model="apiSettingsForm.rateLimitRequests" :min="1" :max="1000" />
            </FormItem>
            
            <FormItem label="时间窗口(秒)" v-if="apiSettingsForm.enableRateLimit">
              <InputNumber v-model="apiSettingsForm.rateLimitDuration" :min="1" :max="3600" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">日志记录</h3>
            
            <FormItem label="启用API日志">
              <Switch v-model="apiSettingsForm.enableLogging" />
            </FormItem>
          </div>

          <div class="form-actions">
            <Button type="primary" @click="saveApiSettings" :loading="savingSettings">保存设置</Button>
          </div>
        </Form>
      </TabPane>

      <!-- 通知设置 -->
      <TabPane tab="通知设置" key="notification">
        <Form
          ref="notificationSettingsRef"
          :model="notificationSettingsForm"
          :rules="notificationSettingsRules"
          :labelCol="{ span: 6 }"
          :wrapperCol="{ span: 16 }"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">通知方式</h3>
            
            <FormItem label="启用邮件通知">
              <Switch v-model="notificationSettingsForm.enableEmailNotifications" />
            </FormItem>
            
            <FormItem label="启用短信通知">
              <Switch v-model="notificationSettingsForm.enableSmsNotifications" />
            </FormItem>
            
            <FormItem label="启用推送通知">
              <Switch v-model="notificationSettingsForm.enablePushNotifications" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">联系信息</h3>
            
            <FormItem label="管理员邮箱" name="adminEmail">
              <Input v-model="notificationSettingsForm.adminEmail" placeholder="请输入管理员邮箱" />
            </FormItem>
            
            <FormItem label="支持邮箱" name="supportEmail">
              <Input v-model="notificationSettingsForm.supportEmail" placeholder="请输入支持邮箱" />
            </FormItem>
            
            <FormItem label="邮件模板签名">
              <Input v-model="notificationSettingsForm.emailTemplateSignoff" placeholder="请输入邮件模板签名" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">告警设置</h3>
            
            <FormItem label="启用错误告警">
              <Switch v-model="notificationSettingsForm.enableErrorAlerts" />
            </FormItem>
            
            <FormItem label="启用性能告警">
              <Switch v-model="notificationSettingsForm.enablePerformanceAlerts" />
            </FormItem>
          </div>

          <div class="form-actions">
            <Button type="primary" @click="saveNotificationSettings" :loading="savingSettings">保存设置</Button>
          </div>
        </Form>
      </TabPane>

      <!-- 高级设置 -->
      <TabPane tab="高级设置" key="advanced">
        <Form
          ref="advancedSettingsRef"
          :model="advancedSettingsForm"
          :labelCol="{ span: 6 }"
          :wrapperCol="{ span: 16 }"
          class="settings-form"
        >
          <div class="form-section">
            <h3 class="section-title">系统调试</h3>
            
            <FormItem label="启用调试模式">
              <Switch v-model="advancedSettingsForm.enableDebugMode" />
              <div class="setting-hint">启用后将显示更多调试信息</div>
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">安全设置</h3>
            
            <FormItem label="启用CORS">
              <Switch v-model="advancedSettingsForm.enableCors" />
            </FormItem>
            
            <FormItem label="启用SSL">
              <Switch v-model="advancedSettingsForm.enableSsl" />
            </FormItem>
            
            <FormItem label="SSL证书路径" v-if="advancedSettingsForm.enableSsl">
              <Input v-model="advancedSettingsForm.sslCertPath" placeholder="请输入SSL证书路径" />
            </FormItem>
            
            <FormItem label="SSL私钥路径" v-if="advancedSettingsForm.enableSsl">
              <Input v-model="advancedSettingsForm.sslKeyPath" placeholder="请输入SSL私钥路径" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">会话设置</h3>
            
            <FormItem label="会话超时(分钟)">
              <InputNumber v-model="advancedSettingsForm.sessionTimeout" :min="1" :max="120" />
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">文件上传</h3>
            
            <FormItem label="最大文件大小(MB)">
              <InputNumber v-model="advancedSettingsForm.maxFileSize" :min="1" :max="100" />
            </FormItem>
            
            <FormItem label="允许的文件类型">
              <Input v-model="advancedSettingsForm.allowedFileTypes" placeholder="例如: jpg,jpeg,png,pdf,doc,docx" />
              <div class="setting-hint">多个文件类型用逗号分隔</div>
            </FormItem>
          </div>

          <div class="form-section">
            <h3 class="section-title">数据分析</h3>
            
            <FormItem label="启用数据分析">
              <Switch v-model="advancedSettingsForm.enableAnalytics" />
            </FormItem>
          </div>

          <div class="form-actions">
            <Button type="primary" @click="saveAdvancedSettings" :loading="savingSettings">保存设置</Button>
          </div>
        </Form>
      </TabPane>
    </Tabs>
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