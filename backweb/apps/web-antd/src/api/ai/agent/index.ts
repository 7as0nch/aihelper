import type {
  AgentInfo,
  BatchBindSubAgentsRequest,
  CreateAgentRequest,
  UpdateAgentRequest,
} from './model';

import type { ID, PageQuery, PageResult } from '#/api/common';

import { requestClient } from '#/api/request';

enum Api {
  root = '/api/ai/agents',
}

/**
 * 查询 Agent 分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function agentList(params?: PageQuery) {
  return requestClient.get<PageResult<AgentInfo>>(Api.root, { params });
}

/**
 * 查询 Agent 信息
 * @param agentId Agent ID
 * @returns Agent 信息
 */
export function agentInfo(agentId: ID) {
  return requestClient.get<AgentInfo>(`${Api.root}/${agentId}`);
}

/**
 * Agent 新增
 * @param data 参数
 * @returns AgentInfo
 */
export function agentAdd(data: CreateAgentRequest) {
  return requestClient.postWithMsg<AgentInfo>(Api.root, data);
}

/**
 * Agent 更新
 * @param data 参数
 * @returns void
 */
export function agentUpdate(data: UpdateAgentRequest) {
  return requestClient.putWithMsg<void>(`${Api.root}/${data.id}`, data);
}

/**
 * Agent 删除
 * @param agentId Agent ID
 * @returns void
 */
export function agentRemove(agentId: ID) {
  return requestClient.deleteWithMsg<void>(`${Api.root}/${agentId}`);
}

/**
 * 批量绑定子 Agent
 * @param data 参数
 * @returns void
 */
export function batchBindSubAgents(data: BatchBindSubAgentsRequest) {
  return requestClient.postWithMsg<void>(
    `${Api.root}/${data.agentId}/sub-agents`,
    data,
  );
}
