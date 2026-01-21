/**
 * Application 类型定义
 * 基于 backend/api/ai/application.proto
 */

import type { AgentModel } from '../agent/model';

/**
 * AppAgent - 应用中的 Agent 配置
 * 映射 proto 的 AppAgent message
 */
export interface AppAgent {
  id?: number;
  name: string;
  code: string;
  description?: string;
  adapterType: number; // 适配器类型 1. adk, 2. deepadk
  originalModelId?: number; // 来源 AI 模型 ID
  aiModel?: AgentModel; // AI 模型配置
  maxIteration?: number; // 最大迭代次数
  systemPrompt?: string; // 系统提示词
  userInputPrompt?: string; // 用户输入提示词
  status: number; // 状态 1: 启用, 2: 禁用
  type: number; // 类型 1: 根 Agent, 2: 子 Agent
  withWriteTodos?: boolean; // 是否启用待办事项
  withWebSearchAgent?: boolean; // 是否启用网络搜索
  systemType?: number; // 1-系统内置, 2-用户自定义
  createdAt?: number;
  subAgents?: AppAgent[]; // 子 Agent 列表
}

/**
 * ApplicationInfo - 应用信息
 * 映射 proto 的 ApplicationInfo message
 */
export interface ApplicationInfo {
  id: number;
  name: string; // 名称
  code: string; // 编码
  description?: string; // 描述
  version?: string; // 版本号
  mode: number; // 模式 1.单agent模式 2.多agent模式
  status: number; // 状态 1: 启用, 2: 禁用
  type: number; // 程序类型 1.预定义 2.自定义
  scope: number; // 作用粒度 1.所有人 2.指定角色 3.指定用户
  selfAgent?: AppAgent; // 自定义 Agent
  createdAt?: number; // 创建时间
  schema?: string; // 拖拽页面布局存档（JSON 字符串）
}

/**
 * 查询应用列表请求
 */
export interface ListApplicationsRequest {
  page?: number;
  pageSize?: number;
  name?: string;
  status?: number;
  type?: number;
  scope?: number;
}

/**
 * 创建应用请求
 */
export interface CreateApplicationRequest {
  name: string;
  code: string;
  description?: string;
  version?: string;
  mode: number;
  status: number;
  type: number;
  scope: number;
  selfAgent?: AppAgent;
  schema?: string;
}

/**
 * 更新应用请求
 */
export interface UpdateApplicationRequest extends CreateApplicationRequest {
  id: number;
}

/**
 * 删除应用请求
 */
export interface DeleteApplicationRequest {
  id: number;
}
