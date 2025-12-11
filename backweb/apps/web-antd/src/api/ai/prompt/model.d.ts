/**
 * Prompt 类型定义
 */
export interface PromptInfo {
    id: number;
    type: number;
    name: string;
    description: string;
    text: string;
    createdAt: number;
}

export interface ListPromptsRequest {
    page?: number;
    pageSize?: number;
    name?: string;
    type?: number;
}

export interface CreatePromptRequest {
    type?: number;
    name: string;
    description?: string;
    text: string;
}

export interface UpdatePromptRequest extends CreatePromptRequest {
    id: number;
}
