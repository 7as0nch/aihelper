/**
 * Tool 类型定义
 */
export interface ToolParam {
    paramName: string;
    paramType: string;
    defaultValue: string;
}

export interface ToolInfo {
    id: number;
    name: string;
    code: string;
    description: string;
    sysType: number;
    type: number;
    status: number;
    params: ToolParam[];
    mcpUrl: string;
    mcpToken: string;
    createdAt: number;
}

export interface ListToolsRequest {
    page?: number;
    pageSize?: number;
    name?: string;
    type?: number;
    status?: number;
}

export interface CreateToolRequest {
    name: string;
    code: string;
    description?: string;
    type?: number;
    status?: number;
    params?: ToolParam[];
    mcpUrl?: string;
    mcpToken?: string;
}

export interface UpdateToolRequest extends CreateToolRequest {
    id: number;
}

export interface BatchBindToolsRequest {
    agentId: number;
    toolCodes: string[];
}
