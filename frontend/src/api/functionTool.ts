import request from './user'; // 导入已配置的axios实例

// 函数工具接口
interface FunctionTool {
  id: string;
  name: string;
  type: string;
  config: Record<string, any>;
  description?: string;
  enabled: boolean;
  create_time: string;
  update_time: string;
  creator_id?: string;
  usage_count?: number;
  last_used_time?: string;
}

// 获取函数工具列表参数接口
interface GetFunctionToolListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  type?: string;
  enabled?: boolean;
}

// 获取函数工具列表响应接口
interface GetFunctionToolListResponse {
  code: number;
  message: string;
  data: {
    list: FunctionTool[];
    total: number;
    page: number;
    pageSize: number;
  };
}

// 获取函数工具列表API
export const getFunctionToolList = (params?: GetFunctionToolListParams): Promise<GetFunctionToolListResponse> => {
  return request.get('/admin/function-tools', { params });
};

// 创建函数工具参数接口
interface CreateFunctionToolParams {
  name: string;
  type: string;
  config: Record<string, any>;
  description?: string;
  enabled?: boolean;
}

// 创建函数工具响应接口
interface CreateFunctionToolResponse {
  code: number;
  message: string;
  data: FunctionTool;
}

// 创建函数工具API
export const createFunctionTool = (params: CreateFunctionToolParams): Promise<CreateFunctionToolResponse> => {
  return request.post('/admin/function-tools', params);
};

// 更新函数工具参数接口
interface UpdateFunctionToolParams {
  id: string;
  name: string;
  type: string;
  config: Record<string, any>;
  description?: string;
  enabled?: boolean;
}

// 更新函数工具响应接口
interface UpdateFunctionToolResponse {
  code: number;
  message: string;
  data: FunctionTool;
}

// 更新函数工具API
export const updateFunctionTool = (params: UpdateFunctionToolParams): Promise<UpdateFunctionToolResponse> => {
  return request.put(`/admin/function-tools/${params.id}`, params);
};

// 删除函数工具响应接口
interface DeleteFunctionToolResponse {
  code: number;
  message: string;
  data: null;
}

// 删除函数工具API
export const deleteFunctionTool = (id: string): Promise<DeleteFunctionToolResponse> => {
  return request.delete(`/admin/function-tools/${id}`);
};

// 获取函数工具详情响应接口
interface GetFunctionToolDetailResponse {
  code: number;
  message: string;
  data: FunctionTool;
}

// 获取函数工具详情API
export const getFunctionToolDetail = (id: string): Promise<GetFunctionToolDetailResponse> => {
  return request.get(`/admin/function-tools/${id}`);
};

// 启用函数工具API
export const enableFunctionTool = (id: string): Promise<any> => {
  return request.patch(`/admin/function-tools/${id}/enable`);
};

// 禁用函数工具API
export const disableFunctionTool = (id: string): Promise<any> => {
  return request.patch(`/admin/function-tools/${id}/disable`);
};

// 测试函数工具API
export const testFunctionTool = (id: string, params: Record<string, any>): Promise<any> => {
  return request.post(`/admin/function-tools/${id}/test`, params);
};

// 获取函数工具类型列表API
export const getFunctionToolTypes = (): Promise<any> => {
  return request.get('/admin/function-tools/types');
};