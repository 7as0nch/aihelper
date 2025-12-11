/**
 * Model 类型定义
 */
export interface ModelInfo {
    id: number;
    category: number;
    modelType: string;
    modelName: string;
    apiKey: string;
    baseUrl: string;
    maxTokens: number;
    temperature: number;
    topP: number;
    priceType: number;
    price: number;
    supplier: string;
    description: string;
    status: number;
    isDefault: number;
    createdAt: number;
}

export interface ListModelsRequest {
    page?: number;
    pageSize?: number;
    modelName?: string;
    status?: number;
}

export interface CreateModelRequest {
    category?: number;
    modelType?: string;
    modelName: string;
    apiKey?: string;
    baseUrl?: string;
    maxTokens?: number;
    temperature?: number;
    topP?: number;
    priceType?: number;
    price?: number;
    supplier?: string;
    description?: string;
    status?: number;
    isDefault?: number;
}

export interface UpdateModelRequest extends CreateModelRequest {
    id: number;
}
