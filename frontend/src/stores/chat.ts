import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

// 消息接口
interface Message {
  id: string;
  content: string;
  role: 'user' | 'assistant' | 'system';
  timestamp: string;
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
  messages: Message[];
  isActive?: boolean;
}

export const useChatStore = defineStore('chat', () => {
  // 聊天状态
  const chatSessions = ref<ChatSession[]>([]);
  const currentSessionId = ref<string>('');
  const loading = ref<boolean>(false);
  const sending = ref<boolean>(false);
  const error = ref<string>('');
  const inputValue = ref<string>('');
  const showQuickReply = ref<boolean>(false);
  const quickReplyOptions = ref<string[]>([
    '如何优化我的代码性能？',
    '请帮我解释这段代码',
    '推荐一些学习资源',
    '如何解决这个bug？'
  ]);

  // 计算属性 - 当前会话
  const currentSession = computed(() => {
    return chatSessions.value.find(session => session.id === currentSessionId.value) || null;
  });

  // 计算属性 - 当前会话的消息列表
  const currentMessages = computed(() => {
    return currentSession.value?.messages || [];
  });

  // 计算属性 - 会话列表（按更新时间排序）
  const sortedSessions = computed(() => {
    return [...chatSessions.value].sort((a, b) => {
      return new Date(b.updateTime).getTime() - new Date(a.updateTime).getTime();
    });
  });

  // 初始化聊天数据
  const initializeChat = async () => {
    try {
      loading.value = true;
      // 从本地存储加载会话数据
      const savedSessions = localStorage.getItem('chatSessions');
      if (savedSessions) {
        chatSessions.value = JSON.parse(savedSessions);
      }
      
      // 如果没有会话，创建一个新的默认会话
      if (chatSessions.value.length === 0) {
        await createNewSession();
      } else {
        // 设置第一个会话为当前会话
        if (chatSessions.value[0]) {
          currentSessionId.value = chatSessions.value[0].id;
        }
      }
    } catch (err) {
      console.error('初始化聊天数据失败:', err);
      error.value = '初始化聊天数据失败';
    } finally {
      loading.value = false;
    }
  };

  // 保存会话数据到本地存储
  const saveSessionsToStorage = () => {
    localStorage.setItem('chatSessions', JSON.stringify(chatSessions.value));
  };

  // 创建新会话
  const createNewSession = async (title?: string) => {
    try {
      const newSession: ChatSession = {
        id: `chat_${Date.now()}`,
        title: title || '新的聊天',
        createTime: new Date().toISOString(),
        updateTime: new Date().toISOString(),
        messages: [],
        isActive: true
      };

      // 将新会话添加到列表开头
      chatSessions.value.unshift(newSession);
      // 设置为当前会话
      currentSessionId.value = newSession.id;
      // 保存到本地存储
      saveSessionsToStorage();
      
      return newSession;
    } catch (err) {
      console.error('创建新会话失败:', err);
      error.value = '创建新会话失败';
      throw err;
    }
  };

  // 切换会话
  const switchSession = (sessionId: string) => {
    // 清除所有会话的活跃状态
    chatSessions.value.forEach(session => {
      session.isActive = false;
    });
    // 设置当前会话为活跃状态
    const targetSession = chatSessions.value.find(session => session.id === sessionId);
    if (targetSession) {
      targetSession.isActive = true;
      currentSessionId.value = sessionId;
      // 保存到本地存储
      saveSessionsToStorage();
    }
  };

  // 删除会话
  const deleteSession = (sessionId: string) => {
    const index = chatSessions.value.findIndex(session => session.id === sessionId);
    if (index !== -1) {
      chatSessions.value.splice(index, 1);
      // 保存到本地存储
      saveSessionsToStorage();
      
      // 如果删除的是当前会话，切换到第一个会话
      if (currentSessionId.value === sessionId && chatSessions.value.length > 0 && chatSessions.value[0]) {
        currentSessionId.value = chatSessions.value[0].id;
        chatSessions.value[0].isActive = true;
      }
    }
  };

  // 重命名会话
  const renameSession = (sessionId: string, newTitle: string) => {
    const session = chatSessions.value.find(session => session.id === sessionId);
    if (session) {
      session.title = newTitle;
      // 保存到本地存储
      saveSessionsToStorage();
    }
  };

  // 发送消息
  const sendMessage = async (content: string) => {
    if (!content.trim() || !currentSession.value || sending.value) {
      return;
    }

    sending.value = true;
    error.value = '';

    // 创建用户消息
    const userMessage: Message = {
      id: `msg_${Date.now()}_user`,
      content: content.trim(),
      role: 'user',
      timestamp: new Date().toISOString(),
      status: 'sending'
    };

    // 创建助手回复占位消息
    const assistantMessage: Message = {
      id: `msg_${Date.now()}_assistant`,
      content: '',
      role: 'assistant',
      timestamp: new Date().toISOString(),
      status: 'sending'
    };

    try {
      // 添加消息到当前会话
      currentSession.value.messages.push(userMessage, assistantMessage);
      // 更新会话更新时间
      currentSession.value.updateTime = new Date().toISOString();
      // 保存到本地存储
      saveSessionsToStorage();

      // 模拟API请求延迟
      await new Promise(resolve => setTimeout(resolve, 500));

      // 更新用户消息状态
      userMessage.status = 'sent';

      // 模拟助手回复
      let replyContent = "感谢您的提问！这是一个示例回复。在实际应用中，这里会是AI模型生成的回答。";
      
      // 根据用户输入内容生成不同的示例回复
      if (content.includes('代码')) {
        replyContent = "关于代码优化，建议您：\n1. 使用更高效的算法\n2. 避免不必要的重复计算\n3. 合理使用缓存机制\n4. 优化数据库查询";
      } else if (content.includes('学习资源')) {
        replyContent = "推荐的学习资源：\n1. 官方文档\n2. 在线教程平台\n3. 技术书籍\n4. 开源项目";
      } else if (content.includes('bug')) {
        replyContent = "解决bug的步骤：\n1. 重现问题\n2. 分析错误日志\n3. 定位问题代码\n4. 编写修复方案\n5. 测试验证";
      }

      // 更新助手消息内容
      assistantMessage.content = replyContent;
      assistantMessage.status = 'sent';

      // 保存到本地存储
      saveSessionsToStorage();

      // 清空输入框
      inputValue.value = '';
    } catch (err) {
      console.error('发送消息失败:', err);
      error.value = '发送消息失败，请重试';
      // 更新消息状态为失败
      userMessage.status = 'failed';
      assistantMessage.status = 'failed';
    } finally {
      sending.value = false;
    }
  };

  // 清空当前会话消息
  const clearSessionMessages = () => {
    if (currentSession.value) {
      currentSession.value.messages = [];
      // 保存到本地存储
      saveSessionsToStorage();
    }
  };

  // 清空所有会话
  const clearAllSessions = () => {
    chatSessions.value = [];
    currentSessionId.value = '';
    // 清除本地存储
    localStorage.removeItem('chatSessions');
    // 创建一个新的默认会话
    createNewSession();
  };

  // 使用快速回复
  const useQuickReply = (reply: string) => {
    inputValue.value = reply;
    sendMessage(reply);
    showQuickReply.value = false;
  };

  // 导入聊天历史
  const importChatHistory = (file: File) => {
    return new Promise<boolean>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const content = e.target?.result as string;
          const importedSessions: ChatSession[] = JSON.parse(content);
          
          // 验证导入的数据格式
          if (Array.isArray(importedSessions)) {
            // 合并导入的会话到现有会话
            chatSessions.value = [...importedSessions, ...chatSessions.value];
            // 保存到本地存储
            saveSessionsToStorage();
            // 如果没有当前会话，设置第一个会话为当前会话
            if (!currentSessionId.value && chatSessions.value.length > 0 && chatSessions.value[0]) {
              currentSessionId.value = chatSessions.value[0].id;
            }
            resolve(true);
          } else {
            throw new Error('无效的聊天历史数据格式');
          }
        } catch (err) {
          console.error('导入聊天历史失败:', err);
          error.value = '导入聊天历史失败，请检查文件格式';
          reject(err);
        }
      };
      reader.onerror = (err) => {
        console.error('读取文件失败:', err);
        error.value = '读取文件失败';
        reject(err);
      };
      reader.readAsText(file);
    });
  };

  // 导出聊天历史
  const exportChatHistory = () => {
    try {
      const dataStr = JSON.stringify(chatSessions.value, null, 2);
      const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
      
      const exportFileDefaultName = `chat_history_${new Date().toISOString().slice(0,10)}.json`;
      
      const linkElement = document.createElement('a');
      linkElement.setAttribute('href', dataUri);
      linkElement.setAttribute('download', exportFileDefaultName);
      linkElement.click();
      
      return true;
    } catch (err) {
      console.error('导出聊天历史失败:', err);
      error.value = '导出聊天历史失败';
      return false;
    }
  };

  return {
    // 状态
    chatSessions,
    currentSessionId,
    loading,
    sending,
    error,
    inputValue,
    showQuickReply,
    quickReplyOptions,
    // 计算属性
    currentSession,
    currentMessages,
    sortedSessions,
    // 方法
    initializeChat,
    createNewSession,
    switchSession,
    deleteSession,
    renameSession,
    sendMessage,
    clearSessionMessages,
    clearAllSessions,
    useQuickReply,
    importChatHistory,
    exportChatHistory
  };
});