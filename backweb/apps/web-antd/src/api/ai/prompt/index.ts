import type { ID, PageQuery, PageResult } from '#/api/common';

import { requestClient } from '#/api/request';

import type {
    CreatePromptRequest,
    PromptInfo,
    UpdatePromptRequest,
} from './model';

enum Api {
    root = '/api/ai/prompts',
}

/**
 * 查询 Prompt 分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function promptList(params?: PageQuery) {
    return requestClient.get<PageResult<PromptInfo>>(Api.root, { params });
}

/**
 * 查询 Prompt 信息
 * @param promptId Prompt ID
 * @returns Prompt 信息
 */
export function promptInfo(promptId: ID) {
    return requestClient.get<PromptInfo>(`${Api.root}/${promptId}`);
}

/**
 * Prompt 新增
 * @param data 参数
 * @returns PromptInfo
 */
export function promptAdd(data: CreatePromptRequest) {
    return requestClient.postWithMsg<PromptInfo>(Api.root, data);
}

/**
 * Prompt 更新
 * @param data 参数
 * @returns void
 */
export function promptUpdate(data: UpdatePromptRequest) {
    return requestClient.putWithMsg<void>(`${Api.root}/${data.id}`, data);
}

/**
 * Prompt 删除
 * @param promptId Prompt ID
 * @returns void
 */
export function promptRemove(promptId: ID) {
    return requestClient.deleteWithMsg<void>(`${Api.root}/${promptId}`);
}
