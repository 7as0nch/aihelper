import type {
  ApplicationInfo,
  CreateApplicationRequest,
  ListApplicationsRequest,
  UpdateApplicationRequest,
} from './model';

import type { ID, PageQuery, PageResult } from '#/api/common';

import { requestClient } from '#/api/request';

enum Api {
  root = '/api/ai/applications',
}

/**
 * 查询 Application 分页列表
 * @param params 搜索条件
 * @returns 分页列表
 */
export function applicationList(params?: ListApplicationsRequest & PageQuery) {
  return requestClient.get<PageResult<ApplicationInfo>>(Api.root, { params });
}

/**
 * 查询 Application 信息
 * @param applicationId Application ID
 * @returns Application 信息
 */
export function applicationInfo(applicationId: ID) {
  return requestClient.get<ApplicationInfo>(`${Api.root}/${applicationId}`);
}

/**
 * Application 新增
 * @param data 参数
 * @returns ApplicationInfo
 */
export function applicationAdd(data: CreateApplicationRequest) {
  return requestClient.postWithMsg<ApplicationInfo>(Api.root, data);
}

/**
 * Application 更新
 * @param data 参数
 * @returns void
 */
export function applicationUpdate(data: UpdateApplicationRequest) {
  return requestClient.putWithMsg<void>(`${Api.root}/${data.id}`, data);
}

/**
 * Application 删除
 * @param applicationId Application ID
 * @returns void
 */
export function applicationRemove(applicationId: ID) {
  return requestClient.deleteWithMsg<void>(`${Api.root}/${applicationId}`);
}
