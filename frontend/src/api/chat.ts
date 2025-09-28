import request from './user'; // 导入已配置的axios实例

// 消息接口
interface Message {
  id: string;
  content: string;
  role: 'user' | 'assistant' | 'system';
  timestamp: string;
  type?: 'text' | 'image' | 'voice' | 'tool_call';
  status?: 'sending' | 'sent' | 'failed';
  toolCalls?: ToolCall[];
}

// 工具调用接口
interface ToolCall {
  id: string;
  type: string;
  function: {
    name: string;
    parameters: Record<string, any>;
  };
  result?: {
    success: boolean;
    data?: any;
    error?: string;
  };
}

// 聊天会话接口
interface ChatSession {
  id: string;
  title: string;
  createTime: string;
  updateTime: string;
  messageCount: number;
  isActive?: boolean;
}

// 获取聊天列表响应接口
interface GetChatListResponse {
  code: number;
  message: string;
  data: {
    list: ChatSession[];
    total: number;
    page: number;
    pageSize: number;
  };
}

// 获取聊天列表API
export const getChatList = (page = 1, pageSize = 20): Promise<GetChatListResponse> => {
  return request.get('/chats', { params: { page, pageSize } });
};

// 创建新聊天参数接口
interface CreateChatParams {
  title: string;
  description?: string;
  initialMessage?: string;
}

// 创建新聊天响应接口
interface CreateChatResponse {
  code: number;
  message: string;
  data: ChatSession;
}

// 创建新聊天API
export const createChat = (params: CreateChatParams): Promise<CreateChatResponse> => {
  return request.post('/chats', params);
};

// 获取聊天详情响应接口
interface GetChatDetailResponse {
  code: number;
  message: string;
  data: {
    chat: ChatSession;
    messages: Message[];
  };
}

// 获取聊天详情API
export const getChatDetail = (chatId: string): Promise<GetChatDetailResponse> => {
  return request.get(`/chats/${chatId}`);
};

// 发送消息参数接口
interface SendMessageParams {
  chatId: string;
  content: string;
  type?: 'text' | 'image' | 'voice';
  useTools?: boolean;
}

// 发送消息响应接口
interface SendMessageResponse {
  code: number;
  message: string;
  data: Message;
}

// 发送消息API
export const sendMessage = (params: SendMessageParams): Promise<SendMessageResponse> => {
  return request.post('/messages', params);
};

// 获取聊天消息响应接口
interface GetMessagesResponse {
  code: number;
  message: string;
  data: {
    list: Message[];
    total: number;
    page: number;
    pageSize: number;
  };
}

// 获取聊天消息API
export const getMessages = (chatId: string, page = 1, pageSize = 50): Promise<GetMessagesResponse> => {
  return request.get(`/chats/${chatId}/messages`, { params: { page, pageSize } });
};

// 更新聊天标题API
export const updateChatTitle = (chatId: string, title: string): Promise<any> => {
  return request.put(`/chats/${chatId}`, { title });
};

// 删除聊天API
export const deleteChat = (chatId: string): Promise<any> => {
  return request.delete(`/chats/${chatId}`);
};

// 清空聊天消息API
export const clearChatMessages = (chatId: string): Promise<any> => {
  return request.delete(`/chats/${chatId}/messages`);
};

// 使用工具调用API
export const callTool = (toolCallId: string, toolName: string, parameters: Record<string, any>): Promise<any> => {
  return request.post('/tools/call', { toolCallId, toolName, parameters });
};

// 获取历史会话统计API
export const getChatStatistics = (): Promise<any> => {
  return request.get('/chats/statistics');
};