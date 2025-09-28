<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { Button, Loading, Badge } from 'vant';
import Avatar from 'vant';
import Input from 'vant';
import { sendMessage, getMessages, createChat, getChatList } from '../../api/chat';

const router = useRouter();
const userStore = useUserStore();

// 聊天状态
const chatId = ref<string>('');
const messages = ref<Array<{
  id: string;
  content: string;
  type: 'user' | 'ai';
  timestamp: number;
  isLoading?: boolean;
}>>([]);
const inputMessage = ref('');
const sending = ref(false);
const isLoadingHistory = ref(false);
const chatHistory = ref<Array<{
  id: string;
  title: string;
  lastMessage: string;
  timestamp: number;
}>>([]);
const showChatList = ref(false);

// 计算属性
const sortedMessages = computed(() => {
  return [...messages.value].sort((a, b) => a.timestamp - b.timestamp);
});

// 滚动到底部
const scrollToBottom = async () => {
  await nextTick();
  const chatContainer = document.getElementById('chat-container');
  if (chatContainer) {
    chatContainer.scrollTop = chatContainer.scrollHeight;
  }
};

// 获取聊天历史
const fetchChatHistory = async () => {
  try {
    const response = await getChatList();
    if (response.code === 200 && response.data) {
      // 转换ChatSession类型为chatHistory所需类型
      chatHistory.value = response.data.list.map((chat: any) => ({
        id: chat.id,
        title: chat.title,
        lastMessage: '', // 由于API中没有提供，暂时设为空
        timestamp: new Date(chat.updateTime || chat.createTime).getTime()
      })) || [];
    }
  } catch (error) {
    console.error('获取聊天历史失败:', error);
  }
}

// 加载聊天消息
const loadChatMessages = async (id: string) => {
  isLoadingHistory.value = true;
  try {
    const response = await getMessages(id);
    if (response.code === 200 && response.data && response.data.list) {
      chatId.value = id;
      messages.value = response.data.list.map((msg: any) => ({
        id: msg.id || Date.now().toString(),
        content: msg.content,
        type: msg.role === 'user' ? 'user' : 'ai',
        timestamp: new Date(msg.timestamp || Date.now()).getTime()
      }));
      await scrollToBottom();
    }
  } catch (error) {
    console.error('加载聊天消息失败:', error);
  } finally {
    isLoadingHistory.value = false;
  }
};

// 创建新聊天
const startNewChat = async () => {
  try {
    const response = await createChat({ title: '新聊天' });
    if (response.code === 200 && response.data) {
      chatId.value = response.data.id;
      messages.value = [];
      await fetchChatHistory();
      showChatList.value = false;
    }
  } catch (error) {
    console.error('创建新聊天失败:', error);
  }
};

// 发送消息
const handleSendMessage = async () => {
  const message = inputMessage.value.trim();
  if (!message || sending.value) return;

  // 添加用户消息
  const userMessage = {
    id: Date.now().toString(),
    content: message,
    type: 'user' as const,
    timestamp: Date.now()
  };
  messages.value.push(userMessage);
  inputMessage.value = '';
  await scrollToBottom();

  // 如果没有聊天ID，创建新聊天
  if (!chatId.value) {
    const response = await createChat({ title: '新聊天' });
    if (response.code === 200 && response.data) {
      chatId.value = response.data.id;
      await fetchChatHistory();
    }
  }

  // 添加AI正在输入的消息
  const aiMessage = {
    id: (Date.now() + 1).toString(),
    content: '',
    type: 'ai' as const,
    timestamp: Date.now() + 1,
    isLoading: true
  };
  messages.value.push(aiMessage);
  await scrollToBottom();

  // 发送消息到后端
  sending.value = true;
  try {
    const response = await sendMessage({ chatId: chatId.value, content: message });
    if (response.code === 200 && response.data) {
      // 更新AI回复消息
      const index = messages.value.findIndex(msg => msg.id === aiMessage.id);
      if (index !== -1) {
        messages.value[index] = {
          ...aiMessage,
          content: response.data.content,
          isLoading: false
        };
        await scrollToBottom();
        await fetchChatHistory();
      }
    }
  } catch (error) {
    console.error('发送消息失败:', error);
    // 更新为错误消息
    const index = messages.value.findIndex(msg => msg.id === aiMessage.id);
    if (index !== -1) {
      messages.value[index] = {
        ...aiMessage,
        content: '抱歉，暂时无法回复，请稍后再试',
        isLoading: false
      };
      await scrollToBottom();
    }
  } finally {
    sending.value = false;
  }
};

// 退出登录
const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

// 初始化
onMounted(async () => {
  // 检查用户是否已登录
  if (!userStore.isLoggedIn) {
    router.push('/login');
    return;
  }

  // 获取聊天历史
  await fetchChatHistory();

  // 如果有聊天历史，加载最近的聊天
  if (chatHistory.value.length > 0 && chatHistory.value[0]) {
    await loadChatMessages(chatHistory.value[0].id);
  } else {
    // 否则创建新聊天
    await startNewChat();
  }
});

// 监听消息变化，自动滚动到底部
watch(sortedMessages, () => {
  scrollToBottom();
}, { deep: true });
</script>

<template>
  <div class="chat-container">
    <!-- 头部 -->
    <header class="chat-header">
      <div class="header-left">
        <Button 
          plain 
          icon="menu-o" 
          class="menu-button"
          @click="showChatList = !showChatList"
        />
        <h1 class="chat-title">AI Chat</h1>
      </div>
      <div class="header-right">
        <Button 
          plain 
          icon="logout"
          @click="handleLogout"
        />
      </div>
    </header>

    <!-- 聊天列表侧边栏 -->
    <div v-if="showChatList" class="chat-list-sidebar">
      <div class="chat-list-header">
        <h2>聊天记录</h2>
        <Button 
          type="primary" 
          size="small"
          class="new-chat-button"
          @click="startNewChat"
        >
          新聊天
        </Button>
      </div>
      <div class="chat-list">
        <div 
          v-for="chat in chatHistory" 
          :key="chat.id"
          class="chat-item" 
          :class="{ active: chat.id === chatId }"
          @click="loadChatMessages(chat.id)"
        >
          <Badge :content="chat.id === chatId ? '当前' : ''" corner>
            <div class="chat-item-content">
              <div class="chat-item-title">{{ chat.title || '新聊天' }}</div>
              <div class="chat-item-message">{{ chat.lastMessage }}</div>
            </div>
          </Badge>
        </div>
        <div v-if="chatHistory.length === 0" class="empty-chat">
          暂无聊天记录
        </div>
      </div>
    </div>

    <!-- 遮罩层 -->
    <div 
      v-if="showChatList"
      class="overlay"
      @click="showChatList = false"
    />

    <!-- 聊天内容区域 -->
    <main id="chat-container" class="chat-content">
      <div v-if="isLoadingHistory" class="loading-container">
        <Loading type="spinner" color="#667eea" size="24px" />
        <p class="loading-text">加载中...</p>
      </div>
      
      <div v-else class="messages-container">
        <!-- 欢迎消息 -->
        <div v-if="sortedMessages.length === 0" class="welcome-message">
          <div class="welcome-content">
            <Avatar class="welcome-avatar" icon="chat-o" size="60px" />
            <h2 class="welcome-title">欢迎使用AI Chat</h2>
            <p class="welcome-desc">请在下方输入您的问题，我会尽力为您解答</p>
          </div>
        </div>
        
        <!-- 消息列表 -->
        <div v-for="message in sortedMessages" :key="message.id" class="message-wrapper">
          <!-- 用户消息 -->
          <div v-if="message.type === 'user'" class="user-message">
            <Avatar class="message-avatar" icon="user-o" />
            <div class="message-content user-content">
              <p>{{ message.content }}</p>
            </div>
          </div>
          
          <!-- AI消息 -->
          <div v-else class="ai-message">
            <Avatar class="message-avatar" icon="chat-o" />
            <div class="message-content ai-content">
              <p v-if="!message.isLoading">{{ message.content }}</p>
              <div v-else class="loading-message">
                <Loading type="spinner" color="#667eea" size="16px" />
                <span>正在输入...</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 输入区域 -->
    <footer class="chat-input-area">
      <div class="input-container">
        <Input
          v-model="inputMessage"
          type="textarea"
          placeholder="请输入您的问题..."
          :disabled="sending"
          rows="1"
          class="chat-input"
          @keyup.enter.exact="handleSendMessage"
          @keyup.enter.shift.prevent
        />
        <Button
          type="primary"
          icon="send"
          :loading="sending"
          class="send-button"
          @click="handleSendMessage"
          :disabled="!inputMessage.trim() || sending"
        />
      </div>
      <p class="input-tip">Shift + Enter 换行，Enter 发送</p>
    </footer>
  </div>
</template>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-primary);
}

/* 头部 */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: white;
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
}

.menu-button {
  margin-right: 15px;
}

.chat-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

/* 聊天列表侧边栏 */
.chat-list-sidebar {
  position: fixed;
  top: 60px;
  left: 0;
  bottom: 60px;
  width: 280px;
  background-color: white;
  border-right: 1px solid var(--border-color);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  z-index: 50;
  transition: transform 0.3s ease;
}

.chat-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.chat-list-header h2 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.new-chat-button {
  height: 32px;
  font-size: 12px;
  padding: 0 12px;
}

.chat-list {
  height: calc(100% - 80px);
  overflow-y: auto;
}

.chat-item {
  padding: 15px 20px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: background-color 0.3s ease;
}

.chat-item:hover {
  background-color: var(--bg-secondary);
}

.chat-item.active {
  background-color: rgba(102, 126, 234, 0.1);
}

.chat-item-title {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-item-message {
  font-size: 0.8rem;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-chat {
  text-align: center;
  color: var(--text-tertiary);
  padding: 40px 20px;
  font-size: 0.9rem;
}

/* 遮罩层 */
.overlay {
  position: fixed;
  top: 60px;
  left: 0;
  right: 0;
  bottom: 60px;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 40;
}

/* 聊天内容区域 */
.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: var(--bg-primary);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.loading-text {
  margin-top: 10px;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.messages-container {
  max-width: 700px;
  margin: 0 auto;
}

/* 欢迎消息 */
.welcome-message {
  text-align: center;
  padding: 60px 20px;
}

.welcome-content {
  display: inline-block;
  background-color: white;
  padding: 40px 60px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.welcome-avatar {
  margin-bottom: 20px;
}

.welcome-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.welcome-desc {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

/* 消息样式 */
.message-wrapper {
  margin-bottom: 20px;
}

.user-message, .ai-message {
  display: flex;
  align-items: flex-start;
}

.user-message {
  justify-content: flex-end;
}

.user-message .message-avatar {
  order: 2;
  margin-left: 10px;
}

.user-content {
  background-color: #667eea;
  color: white;
  order: 1;
}

.ai-message .message-avatar {
  order: 1;
  margin-right: 10px;
}

.ai-content {
  background-color: white;
  color: var(--text-primary);
  order: 2;
}

.message-avatar {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  line-height: 1.5;
}

.message-content p {
  margin: 0;
  word-wrap: break-word;
}

.loading-message {
  display: flex;
  align-items: center;
  color: var(--text-secondary);
}

.loading-message span {
  margin-left: 8px;
  font-size: 0.9rem;
}

/* 输入区域 */
.chat-input-area {
  background-color: white;
  border-top: 1px solid var(--border-color);
  padding: 15px 20px;
  z-index: 10;
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.chat-input {
  flex: 1;
  min-height: 44px;
  max-height: 120px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  padding: 10px 15px;
  font-size: 14px;
  resize: none;
}

.send-button {
  height: 44px;
  width: 44px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}

.input-tip {
  text-align: center;
  font-size: 0.8rem;
  color: var(--text-tertiary);
  margin-top: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-list-sidebar {
    width: 240px;
  }
  
  .message-content {
    max-width: 80%;
  }
}

@media (max-width: 480px) {
  .chat-header {
    padding: 12px 15px;
  }
  
  .chat-title {
    font-size: 1.1rem;
  }
  
  .chat-content {
    padding: 15px;
  }
  
  .message-content {
    max-width: 85%;
  }
  
  .chat-input-area {
    padding: 12px 15px;
  }
}
</style>