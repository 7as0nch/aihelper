/**
 * Agent 类型定义
 */
export interface AgentInfo {
  id: number;
  name: string;
  code: string;
  description: string;
  adapterType: number; // 适配器类型 1. adk, 2. deepadk
  aiModelId: number;
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
  aiModelId?: number;
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
