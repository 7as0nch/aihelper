import type { ID, PageQuery, PageResult } from '#/api/common';

import { requestClient } from '#/api/request';

import type {
    CreateModelRequest,
    ModelInfo,
    UpdateModelRequest,
} from './model';

enum Api {
    root = '/api/ai/models',
}

/**
 * 查询 Model 分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function modelList(params?: PageQuery) {
    return requestClient.get<PageResult<ModelInfo>>(Api.root, { params });
}

/**
 * 查询 Model 信息
 * @param modelId Model ID
 * @returns Model 信息
 */
export function modelInfo(modelId: ID) {
    return requestClient.get<ModelInfo>(`${Api.root}/${modelId}`);
}

/**
 * Model 新增
 * @param data 参数
 * @returns ModelInfo
 */
export function modelAdd(data: CreateModelRequest) {
    return requestClient.postWithMsg<ModelInfo>(Api.root, data);
}

/**
 * Model 更新
 * @param data 参数
 * @returns void
 */
export function modelUpdate(data: UpdateModelRequest) {
    return requestClient.putWithMsg<void>(`${Api.root}/${data.id}`, data);
}

/**
 * Model 删除
 * @param modelId Model ID
 * @returns void
 */
export function modelRemove(modelId: ID) {
    return requestClient.deleteWithMsg<void>(`${Api.root}/${modelId}`);
}
