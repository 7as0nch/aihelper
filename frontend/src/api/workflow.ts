import request from './user'; // 导入已配置的axios实例

// 工作流步骤接口
interface WorkflowStep {
  id?: string;
  order: number;
  function_tool_id: string;
  function_name?: string;
  condition?: string;
  variable_mapping?: string;
  config?: Record<string, any>;
}

// 工作流接口
interface Workflow {
  id: string;
  name: string;
  description?: string;
  steps: WorkflowStep[];
  create_time: string;
  update_time: string;
  creator_id?: string;
  status: 'active' | 'inactive' | 'draft';
  run_count?: number;
  success_rate?: number;
}

// 获取工作流列表参数接口
interface GetWorkflowListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: 'active' | 'inactive' | 'draft';
}

// 获取工作流列表响应接口
interface GetWorkflowListResponse {
  code: number;
  message: string;
  data: {
    list: Workflow[];
    total: number;
    page: number;
    pageSize: number;
  };
}

// 获取工作流列表API
export const getWorkflowList = (params?: GetWorkflowListParams): Promise<GetWorkflowListResponse> => {
  return request.get('/admin/workflows', { params });
};

// 创建工作流参数接口
interface CreateWorkflowParams {
  name: string;
  description?: string;
  steps: Array<{
    order: number;
    function_tool_id: string;
    condition?: string;
    variable_mapping?: string;
    config?: Record<string, any>;
  }>;
  status?: 'active' | 'inactive' | 'draft';
}

// 创建工作流响应接口
interface CreateWorkflowResponse {
  code: number;
  message: string;
  data: Workflow;
}

// 创建工作流API
export const createWorkflow = (params: CreateWorkflowParams): Promise<CreateWorkflowResponse> => {
  return request.post('/admin/workflows', params);
};

// 更新工作流参数接口
interface UpdateWorkflowParams {
  id: string;
  name: string;
  description?: string;
  steps: Array<{
    order: number;
    function_tool_id: string;
    condition?: string;
    variable_mapping?: string;
    config?: Record<string, any>;
  }>;
  status?: 'active' | 'inactive' | 'draft';
}

// 更新工作流响应接口
interface UpdateWorkflowResponse {
  code: number;
  message: string;
  data: Workflow;
}

// 更新工作流API
export const updateWorkflow = (params: UpdateWorkflowParams): Promise<UpdateWorkflowResponse> => {
  return request.put(`/admin/workflows/${params.id}`, params);
};

// 删除工作流响应接口
interface DeleteWorkflowResponse {
  code: number;
  message: string;
  data: null;
}

// 删除工作流API
export const deleteWorkflow = (id: string): Promise<DeleteWorkflowResponse> => {
  return request.delete(`/admin/workflows/${id}`);
};

// 获取工作流详情响应接口
interface GetWorkflowDetailResponse {
  code: number;
  message: string;
  data: Workflow;
}

// 获取工作流详情API
export const getWorkflowDetail = (id: string): Promise<GetWorkflowDetailResponse> => {
  return request.get(`/admin/workflows/${id}`);
};

// 激活工作流API
export const activateWorkflow = (id: string): Promise<any> => {
  return request.patch(`/admin/workflows/${id}/activate`);
};

// 停用工作流API
export const deactivateWorkflow = (id: string): Promise<any> => {
  return request.patch(`/admin/workflows/${id}/deactivate`);
};

// 复制工作流API
export const duplicateWorkflow = (id: string): Promise<CreateWorkflowResponse> => {
  return request.post(`/admin/workflows/${id}/duplicate`);
};

// 获取工作流执行记录API
export const getWorkflowExecutions = (workflowId: string, params?: {
  page?: number;
  pageSize?: number;
  startTime?: string;
  endTime?: string;
  status?: 'success' | 'failed' | 'running';
}): Promise<any> => {
  return request.get(`/admin/workflows/${workflowId}/executions`, { params });
};

// 获取工作流执行详情API
export const getWorkflowExecutionDetail = (executionId: string): Promise<any> => {
  return request.get(`/admin/workflows/executions/${executionId}`);
};