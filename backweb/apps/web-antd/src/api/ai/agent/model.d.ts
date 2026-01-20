/**
 * Agent 类型定义
 */
export interface AgentModel {
  category: number;
  modelType: string;
  modelName: string;
  apiKey: string;
  baseUrl: string;
  temperature: number;
  topP: number;
}

export interface AgentInfo {
  id: number;
  name: string;
  code: string;
  description: string;
  adapterType: number; // 适配器类型 1. adk, 2. deepadk
  originalModelId: number; // 来源ai模型id: 原始model ID：追溯引用的模型配置。
  aiModel?: AgentModel; // 通过original_model_id获取的模型配置，支持修改。
  maxIteration: number;
  systemPrompt: string;
  userInputPrompt: string;
  status: number;
  type: number;
  order: number;
  withWriteTodos: boolean;
  withWebSearchAgent: boolean;
  systemType: number;
  createdAt: number;
  subAgentIds: number[];
}

export interface ListAgentsRequest {
  page?: number;
  pageSize?: number;
  name?: string;
  status?: number;
  type?: number;
}

export interface CreateAgentRequest {
  name: string;
  code: string;
  description?: string;
  adapterType?: number;
  originalModelId?: number; // 来源ai模型id: 原始model ID：追溯引用的模型配置。
  aiModel?: AgentModel; // 通过original_model_id获取的模型配置，支持修改。
  maxIteration?: number;
  systemPrompt?: string;
  userInputPrompt?: string;
  status?: number;
  type?: number;
  order?: number;
  withWriteTodos?: boolean;
  withWebSearchAgent?: boolean;
  subAgentIds?: number[];
}

export interface UpdateAgentRequest extends CreateAgentRequest {
  id: number;
}

export interface BatchBindSubAgentsRequest {
  agentId: number;
  subAgentIds: number[];
}
