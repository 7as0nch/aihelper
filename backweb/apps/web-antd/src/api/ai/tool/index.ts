import type { ID, PageQuery, PageResult } from '#/api/common';

import { requestClient } from '#/api/request';

import type {
    BatchBindToolsRequest,
    CreateToolRequest,
    ToolInfo,
    UpdateToolRequest,
} from './model';

enum Api {
    root = '/api/ai/tools',
    bindTools = '/api/ai/agents',
}

/**
 * 查询 Tool 分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function toolList(params?: PageQuery) {
    return requestClient.get<PageResult<ToolInfo>>(Api.root, { params });
}

/**
 * 查询 Tool 信息
 * @param toolId Tool ID
 * @returns Tool 信息
 */
export function toolInfo(toolId: ID) {
    return requestClient.get<ToolInfo>(`${Api.root}/${toolId}`);
}

/**
 * Tool 新增
 * @param data 参数
 * @returns ToolInfo
 */
export function toolAdd(data: CreateToolRequest) {
    return requestClient.postWithMsg<ToolInfo>(Api.root, data);
}

/**
 * Tool 更新
 * @param data 参数
 * @returns void
 */
export function toolUpdate(data: UpdateToolRequest) {
    return requestClient.putWithMsg<void>(`${Api.root}/${data.id}`, data);
}

/**
 * Tool 删除
 * @param toolId Tool ID
 * @returns void
 */
export function toolRemove(toolId: ID) {
    return requestClient.deleteWithMsg<void>(`${Api.root}/${toolId}`);
}

/**
 * 批量绑定工具到 Agent
 * @param data 参数
 * @returns void
 */
export function batchBindTools(data: BatchBindToolsRequest) {
    return requestClient.postWithMsg<void>(
        `${Api.bindTools}/${data.agentId}/tools`,
        data,
    );
}
