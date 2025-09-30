<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';
import 'ant-design-vue/dist/reset.css';
import { Button as AButton } from 'ant-design-vue';
import { Spin as ASpin } from 'ant-design-vue';
import { Avatar as AAvatar } from 'ant-design-vue';
import { Input as AInput } from 'ant-design-vue';
import { Textarea as ATextarea } from 'ant-design-vue';
import { message } from 'ant-design-vue';
import { sendMessage, getMessages, createChat, getChatList } from '../../api/chat';
import { format } from 'date-fns';

const router = useRouter();
const userStore = useUserStore();

// 聊天状态
const chatId = ref<string>('');

// 处理键盘事件
const handleKeyEvent = (event: KeyboardEvent) => {
  if (event.key.toLowerCase() === 'enter') {
    console.log('Enter pressed', 'shiftKey:', event.shiftKey);
    if (event.shiftKey) {
      // 阻止默认行为（Enter 键默认提交表单等行为），允许换行
      event.preventDefault();
      // 手动添加换行符
      inputMessage.value += '\n';
    } else {
      // Enter键发送消息
      event.preventDefault();
      handleSendMessage();
    }
  }
};
const messages = ref<Array<{
  id: string;
  content: string;
  type: 'user' | 'ai';
  timestamp: number;
  isLoading?: boolean;
  status?: 'sending' | 'sent' | 'failed';
}>>([]);
const inputMessage = ref<string>('');
const sending = ref(false);
const isLoadingHistory = ref(false);
const chatHistory = ref<Array<{
  id: string;
  title: string;
  lastMessage: string;
  timestamp: number;
  messageCount: number;
}>>([]);
const showChatList = ref(false);
const showToolMenu = ref(false);
const isMobile = ref(false);
const scrollTimeout = ref<number | null>(null);

// 快速提问建议
const quickSuggestions = [
  { id: '1', content: '如何优化我的代码性能？' },
  { id: '2', content: '请帮我解释这段代码' },
  { id: '3', content: '推荐一些学习资源' }
];

// 计算属性
const sortedMessages = computed(() => {
  return [...messages.value].sort((a, b) => a.timestamp - b.timestamp);
});

const formattedTime = (timestamp: number) => {
  return format(new Date(timestamp), 'HH:mm');
};

const formatDateForHistory = (timestamp: number) => {
  const now = new Date();
  const messageDate = new Date(timestamp);
  const diffDays = Math.floor((now.getTime() - messageDate.getTime()) / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) {
    return format(messageDate, 'HH:mm');
  } else if (diffDays === 1) {
    return '昨天 ' + format(messageDate, 'HH:mm');
  } else if (diffDays < 7) {
    return format(messageDate, 'EEEE HH:mm');
  } else {
    return format(messageDate, 'yyyy-MM-dd HH:mm');
  }
};

// 监听窗口大小变化
const handleResize = () => {
  isMobile.value = window.innerWidth < 768;
  if (isMobile.value && showChatList.value) {
    showChatList.value = false;
  }
};

// 滚动到底部（带动画效果）
const scrollToBottom = async () => {
  await nextTick();
  const chatContainer = document.getElementById('chat-container');
  if (chatContainer) {
    // 清除之前的定时器以避免冲突
    if (scrollTimeout.value) {
      window.clearTimeout(scrollTimeout.value);
    }
    
    scrollTimeout.value = window.setTimeout(() => {
      const endPosition = chatContainer.scrollHeight;
      const startPosition = chatContainer.scrollTop;
      const distance = endPosition - startPosition;
      const duration = 300; // 300ms动画
      let startTime: number | null = null;
      
      const easeInOutQuad = (t: number) => {
        return t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
      };
      
      const scrollStep = (timestamp: number) => {
        if (!startTime) startTime = timestamp;
        const progress = Math.min((timestamp - startTime) / duration, 1);
        const easeProgress = easeInOutQuad(progress);
        
        chatContainer.scrollTop = startPosition + (distance * easeProgress);
        
        if (progress < 1) {
          window.requestAnimationFrame(scrollStep);
        }
      };
      
      window.requestAnimationFrame(scrollStep);
    }, 50);
  }
};

// 获取聊天历史
const fetchChatHistory = async () => {
  try {
    const response = await getChatList();
    if (response && response.code === 200 && response.data && response.data.list) {
      chatHistory.value = response.data.list.map((chat: any) => ({
        id: chat.id,
        title: chat.title,
        lastMessage: chat.lastMessage || '',
        timestamp: new Date(chat.updateTime || chat.createTime).getTime(),
        messageCount: chat.messageCount || 0
      }));
    }
  } catch (error) {
    console.error('获取聊天历史失败:', error);
    message.error('获取聊天历史失败');
  }
};

// 加载聊天消息
const loadChatMessages = async (id: string) => {
  if (id === chatId.value && messages.value.length > 0) {
    // 如果已经是当前聊天并且有消息，不重新加载
    return;
  }
  
  isLoadingHistory.value = true;
  messages.value = []; // 清空当前消息
  
  try {
    const response = await getMessages(id);
    if (response.code === 200 && response.data && response.data.list) {
      chatId.value = id;
      messages.value = response.data.list.map((msg: any) => ({
        id: msg.id || `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        content: msg.content || '',
        type: msg.role === 'user' ? 'user' : 'ai',
        timestamp: new Date(msg.timestamp || Date.now()).getTime(),
        status: msg.status === 'failed' ? 'failed' : 'sent'
      }));
      await scrollToBottom();
    }
  } catch (error) {
    console.error('加载聊天消息失败:', error);
    message.error('加载聊天消息失败');
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
      await scrollToBottom();
      // 聚焦到输入框
      nextTick(() => {
        const input = document.querySelector('.chat-input') as HTMLTextAreaElement;
        if (input) {
          input.focus();
        }
      });
    }
  } catch (error) {
    console.error('创建新聊天失败:', error);
    message.error('创建新聊天失败');
  }
};

// 发送消息
const handleSendMessage = async () => {
  const msg = inputMessage.value.trim();
  if (!msg || sending.value) return;

  // 添加用户消息
  const userMessage = {
    id: `msg_${Date.now()}_user`,
    content: msg,
    type: 'user' as const,
    timestamp: Date.now(),
    status: 'sending' as const
  };
  messages.value.push(userMessage);

  inputMessage.value = '';
  await scrollToBottom();

  // 如果没有聊天ID，创建新聊天
  if (!chatId.value) {
    try {
      const response = await createChat({ title: msg.length > 30 ? msg.substring(0, 30) + '...' : msg });
      if (response.code === 200 && response.data) {
        chatId.value = response.data.id;
        await fetchChatHistory();
      }
    } catch (error) {
      console.error('创建新聊天失败:', error);
      (userMessage as any).status = 'failed';
      message.error('获取聊天历史失败');
      return;
    }
  }

  // 添加AI正在输入的消息
  const aiMessage = {
    id: `msg_${Date.now()}_ai`,
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
    const response = await sendMessage({ chatId: chatId.value, content: msg });
    if (response.code === 200 && response.data) {
      // 更新用户消息状态
      (userMessage as any).status = 'sent';
      
      // 更新AI回复消息
      const index = messages.value.findIndex(msg => msg.id === aiMessage.id);
      if (index !== -1) {
        // 模拟打字机效果
        const content = (response.data && response.data.content) || '抱歉，我暂时无法为您提供这方面的解答。';
        if (!content) {
          aiMessage.isLoading = false;
          aiMessage.content = '抱歉，我暂时无法为您提供这方面的解答。';
          await scrollToBottom();
          await fetchChatHistory();
          return;
        }
        const chars = content.split('');
        let currentIndex = 0;
        
        aiMessage.isLoading = false;
        aiMessage.content = '';
        
        const typeWriter = async () => {
          if (currentIndex < chars.length) {
            aiMessage.content += chars[currentIndex];
            currentIndex++;
            window.setTimeout(typeWriter, Math.random() * 30 + 10); // 10-40ms间隔
          } else {
            await scrollToBottom();
            await fetchChatHistory();
          }
        };
        
        typeWriter();
      }
    }
  } catch (error) {
    console.error('发送消息失败:', error);
    // 更新为错误消息
    (userMessage as any).status = 'failed';
    const index = messages.value.findIndex(msg => msg.id === aiMessage.id);
    if (index !== -1) {
      messages.value[index] = {
        ...aiMessage,
        content: '抱歉，暂时无法回复，请稍后再试',
        isLoading: false,
        status: 'failed' as const
      };
      await scrollToBottom();
    }
    message.error('发送消息失败');
  } finally {
    sending.value = false;
  }
};

// 上传图片
const handleUploadImage = () => {
  showToolMenu.value = false;
  // 在实际应用中，这里会打开文件选择对话框
  message.info('上传图片功能待实现');
};

// 添加链接
const handleAddLink = () => {
  showToolMenu.value = false;
  // 在实际应用中，这里会打开链接输入对话框
  message.info('添加链接功能待实现');
};

// 上传文档
const handleUploadDocument = () => {
  showToolMenu.value = false;
  // 在实际应用中，这里会打开文件选择对话框
  message.info('上传文档功能待实现');
};

// 退出登录
const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

// 发送按钮是否展示：
const notShowSendButton = computed(() => {
  return inputMessage.value.trim().length === 0 || sending.value;
});

// 初始化
onMounted(async () => {
  // 检查用户是否已登录
  if (!userStore.isLoggedIn) {
    router.push('/login');
    return;
  }

  // 设置响应式状态
  handleResize();
  window.addEventListener('resize', handleResize);

  // 获取聊天历史
  await fetchChatHistory();

  // 如果有聊天历史，加载最近的聊天
  if (chatHistory && chatHistory.value.length > 0) {
    if (chatHistory.value[0]) {
      await loadChatMessages(chatHistory.value[0].id);
    }
  } else {
    // 否则创建新聊天
    await startNewChat();
  }
});

// 清理
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
  if (scrollTimeout.value) {
    window.clearTimeout(scrollTimeout.value);
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
        <AButton 
          type="text" 
          icon="menu-fold"
          class="menu-button"
          @click="showChatList = !showChatList"
        />
      </div>
      <h1 class="chat-title">AI Chat</h1>
      <div class="header-right">
        <AButton 
          type="text" 
          icon="logout"
          @click="handleLogout"
        />
      </div>
    </header>

    <!-- 聊天列表侧边栏 -->
    <transition name="slide">
      <div v-if="showChatList" class="chat-list-sidebar">
        <div class="chat-list-header">
          <h2>聊天记录</h2>
          <AButton 
            type="primary" 
            size="small"
            class="new-chat-button"
            @click="startNewChat"
          >
            新聊天
          </AButton>
        </div>
        <div class="chat-list">
          <div 
            v-for="chat in chatHistory" 
            :key="chat.id"
            class="chat-item" 
            :class="{ active: chat.id === chatId }"
            @click="loadChatMessages(chat.id)"
          >
            <div class="chat-item-content">
              <div class="chat-item-title">{{ chat.title || '新聊天' }}</div>
              <div class="chat-item-message">{{ chat.lastMessage || '暂无消息' }}</div>
              <div class="chat-item-time">{{ formatDateForHistory(chat.timestamp) }}</div>
            </div>
          </div>
          <div v-if="chatHistory.length === 0" class="empty-chat">
            <div class="empty-icon">💬</div>
            <div class="empty-text">暂无聊天记录</div>
            <AButton 
              type="primary" 
              size="small"
              class="create-first-chat"
              @click="startNewChat"
            >
              创建第一个聊天
            </AButton>
          </div>
        </div>
      </div>
    </transition>

    <!-- 遮罩层 -->
    <transition name="fade">
      <div 
        v-if="showChatList"
        class="overlay"
        @click="showChatList = false"
      />
    </transition>

    <!-- 聊天内容区域 -->
    <main id="chat-container" class="chat-content">
      <div v-if="isLoadingHistory" class="loading-container">
        <ASpin size="large" tip="加载聊天记录中..." />
      </div>
      
      <div v-else class="messages-container">
        <!-- 欢迎消息 -->
        <transition name="fade-up">
          <div v-if="sortedMessages.length === 0" class="welcome-message">
            <div class="welcome-content">
              <div class="welcome-avatar">🤖</div>
              <h2 class="welcome-title">欢迎使用AI Chat</h2>
              <p class="welcome-desc">我是您的AI助手，有什么可以帮助您的吗？</p>
              <p class="welcome-hint">您可以在下方输入框上方的快速提问区域选择常见问题，或直接输入您的问题。</p>
            </div>
          </div>
        </transition>
        
        <!-- 消息列表 -->
        <div class="message-list">
          <transition-group name="message-fade" tag="div">
            <div v-for="message in sortedMessages" :key="message.id" class="message-wrapper">
              <!-- 用户消息 -->
              <div v-if="message.type === 'user'" class="user-message">
                <div class="message-content user-content" :class="{ 'sending': message.status === 'sending', 'failed': message.status === 'failed' }">
                  <p>{{ message.content }}</p>
                  <span class="message-time">{{ formattedTime(message.timestamp) }}</span>
                  <span v-if="message.status === 'failed'" class="message-status">
                    <span class="anticon anticon-warning"><svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor"><path d="M464 480a48 48 0 1 0 96 0 48 48 0 1 0-96 0zM494.752 816a54.464 54.464 0 0 1-38.528-16.064L186.432 520.192a56.32 56.32 0 0 1-16.064-38.592c0-15.04 5.888-29.696 16.064-40.96l269.76-269.568a56.32 56.32 0 0 1 77.184 0L837.568 440.64a56.32 56.32 0 0 1 16.064 40.96 54.464 54.464 0 0 1-16.064 38.592L533.248 800.064A54.272 54.272 0 0 1 494.752 816z"></path></svg></span>
                  </span>
                </div>
                <AAvatar class="message-avatar" icon="user" :size="40" />
              </div>
              
              <!-- AI消息 -->
              <div v-else class="ai-message">
                <AAvatar class="message-avatar" icon="message" :size="40" />
                <div class="message-content ai-content" :class="{ 'loading': message.isLoading, 'failed': message.status === 'failed' }">
                  <p v-if="!message.isLoading">{{ message.content }}</p>
                  <div v-else class="loading-message">
                    <div class="loading-dots">
                      <span class="dot"></span>
                      <span class="dot"></span>
                      <span class="dot"></span>
                    </div>
                    <span>正在输入...</span>
                  </div>
                  <span class="message-time">{{ formattedTime(message.timestamp) }}</span>
                  <span v-if="message.status === 'failed'" class="message-status">
                    <i class="van-icon van-icon-warning-o"></i>
                  </span>
                </div>
              </div>
            </div>
          </transition-group>
        </div>
      </div>
    </main>
    
    <!-- 快速提问横向滚动区域 -->
    <div class="quick-suggestions-container" style="margin-left: 5px;">
      <div class="quick-suggestions-scroll">
        <AButton
          v-for="suggestion in quickSuggestions"
          :key="suggestion.id"
          type="default"
          size="small"
          class="quick-suggestion-btn"
          @click="() => { inputMessage = suggestion.content; handleSendMessage(); }"
        >
          {{ suggestion.content }}
        </AButton>
      </div>
    </div>
    <!-- 输入区域 -->
    <footer class="chat-input-area">
      
      <div class="input-container">
        <!-- 工具按钮 -->
        <div class="tool-button-container">
          <div class="tool-button-wrapper">
            <AButton 
              type="text" 
              icon="plus"
              class="tool-button"
              @click="showToolMenu = !showToolMenu"
            />
            <!-- 工具菜单 -->
            <transition name="fade">
              <div v-if="showToolMenu" class="tool-menu">
                <div class="tool-menu-item" @click="handleUploadImage">
                  <span class="anticon anticon-picture"><svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor"><path d="M870.7 386.9a8 8 0 0 0-6.5-3.9H736V256c0-17.7-14.3-32-32-32H320c-17.7 0-32 14.3-32 32v97H160c-17.7 0-32 14.3-32 32v416c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V419.8c0-17.7-14.3-32.9-32-32.9zm-425.2 74.7c11.3 0 20.4-9.1 20.4-20.4s-9.1-20.4-20.4-20.4-20.4 9.1-20.4 20.4 9.1 20.4 20.4 20.4z m304.6 304.6c-2.9 0-5.6-1.8-6.7-4.6L823.7 528a8.03 8.03 0 0 0 .9-11.3l-39.8-50.8c-3.2-4.1-9.6-4.1-12.9 0L652.3 592c-2.9 3.7-7.6 5.8-12.2 5.8H406.4c-4.6 0-9.3-2.1-12.2-5.8L312 518.8c-3.2-4.1-9.6-4.1-12.9 0l-39.8 50.8c-3.9 5-2.2 12.3 2.8 15.3l167.7 131.7c1.1 2.7 3.8 4.6 6.7 4.6h246.6z m-82.7-439.7H398.4V288h279.2v41.8z"></path></svg></span>
                  <span>上传图片</span>
                </div>
                <div class="tool-menu-item" @click="handleAddLink">
                  <span class="anticon anticon-link"><svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor"><path d="M928 544c0 23.7-10.2 44.8-26.5 59.8L405.7 905.9c-12.9 11.3-33.5 9.9-45.1-3.4-11.1-12.8-9.6-33.2 3.6-45.2L839.8 533.3a8 8 0 0 0 2.6-5.9v-42.7c0-6.5-7.4-10.3-12.7-6.5l-43.2 34.5c-22.4 17.9-52.1 17.9-74.5 0L194 312.3c-22.4-17.9-52.1-17.9-74.5 0L103 346.9c-5.3 4.3-12.7 0.4-12.7-6.5v-42.7c0-6.5 7.4-10.3 12.7-6.5l43.2 34.5c17.1 13.7 38 21.6 59.8 21.6s42.7-7.8 59.8-21.6L530.7 114.3c11.1-12 31.5-13.5 44.1-3.4 12.7 10.2 13.5 30.6 2.3 42.7L239.1 462.6a8 8 0 0 0-2.6 5.9v42.7c0 6.5 7.4 10.3 12.7 6.5l46.3-37c17.1-13.7 38-21.6 59.8-21.6s42.7 7.8 59.8 21.6l325.6 258.2c12.4 9.9 32.8 8.7 44.1-3.4 11.1-12.8 9.6-33.2-3.6-45.2L593.3 533.3a8 8 0 0 0-2.6-5.9v-42.7c0-6.5 7.4-10.3 12.7-6.5l46.3 37c23.2 18.5 58.2 18.5 81.4 0l46.3-37c5.3-4.3 12.7-0.4 12.7 6.5v42.7c0 .2.1.5.1.7z"></path></svg></span>
                  <span>添加链接</span>
                </div>
                <div class="tool-menu-item" @click="handleUploadDocument">
                  <span class="anticon anticon-file-text"><svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor"><path d="M854.6 262.6L637.4 46.4c-6-6-14.1-9.4-22.6-9.4H192c-17.7 0-32 14.3-32 32v832c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V285.3c0-8.5-3.4-16.6-9.4-22.7zM790.2 326H602V137.8L790.2 326zm1.8 562H232V136h302v216a42 42 0 0 0 42 42h216v494z"></path><path d="M344 408H232c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h112c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8zm576-84v48c0 4.4-3.6 8-8 8H632c-4.4 0-8-3.6-8-8v-48c0-4.4 3.6-8 8-8h280c4.4 0 8 3.6 8 8zm-240 96v48c0 4.4-3.6 8-8 8H472c-4.4 0-8-3.6-8-8v-48c0-4.4 3.6-8 8-8h104c4.4 0 8 3.6 8 8zm240 96v48c0 4.4-3.6 8-8 8H632c-4.4 0-8-3.6-8-8v-48c0-4.4 3.6-8 8-8h280c4.4 0 8 3.6 8 8zm-240 96v48c0 4.4-3.6 8-8 8H472c-4.4 0-8-3.6-8-8v-48c0-4.4 3.6-8 8-8h104c4.4 0 8 3.6 8 8z"></path></svg></span>
                  <span>上传文档</span>
                </div>
              </div>
            </transition>
          </div>
        </div>
        
        <!-- 输入框 -->
        <ATextarea
          :value="inputMessage"
          @input="(value) => { inputMessage = value.target.value ? value.target.value : '' }"
          placeholder="请输入您的问题..."
          :disabled="sending"
          :rows="1"
          class="chat-input"
          @keydown="handleKeyEvent"
          :auto-size="{ minRows: 1, maxRows: 10 }"
        />
        
        <!-- 发送按钮 -->
        <AButton
          type="primary"
          icon="↩️"
          :loading="sending"
          class="send-button"
          @click="handleSendMessage"
          :disabled="notShowSendButton"
        />
      </div>
      <p class="input-tip">Shift + Enter 换行，Enter 发送</p>
    </footer>
  </div>
</template>

<style scoped>
/* CSS 变量定义 */
:root {
  --primary-color: #667eea;
  --primary-light: #0324f5;
  --text-primary: #333333;
  --text-secondary: #666;
  --text-tertiary: #999;
  --bg-primary: #f7f8fa;
  --bg-secondary: #f0f2f5;
  --border-color: #e0e0e0;
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-full: 24px;
}

/* 基础样式 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-primary);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

/* 头部 */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: white;
  border-bottom: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  z-index: 10;
  transition: background-color 0.3s ease;
}

.header-left {
  display: flex;
  align-items: center;
}

.menu-button {
  margin-right: 15px;
  font-size: 20px;
  color: var(--text-secondary);
  transition: color 0.3s ease;
}

.menu-button:hover {
  color: var(--primary-color);
}

.chat-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.header-right .van-button {
  color: var(--text-secondary);
  font-size: 18px;
  transition: color 0.3s ease;
}

.header-right .van-button:hover {
  color: var(--primary-color);
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
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.08);
  z-index: 50;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.chat-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
  background-color: white;
}

.chat-list-header h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.new-chat-button {
  height: 32px;
  font-size: 12px;
  padding: 0 12px;
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  transition: all 0.3s ease;
}

.new-chat-button:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-light);
}

.chat-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}

.chat-item {
  padding: 12px 20px;
  cursor: pointer;
  border-bottom: 1px solid transparent;
  transition: all 0.3s ease;
  position: relative;
}

.chat-item:hover {
  background-color: var(--bg-secondary);
}

.chat-item.active {
  background-color: rgba(102, 126, 234, 0.1);
  border-left: 3px solid var(--primary-color);
  padding-left: 17px;
}

.chat-item-content {
  display: flex;
  flex-direction: column;
}

.chat-item-title {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-item-message {
  font-size: 0.85rem;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 2px;
}

.chat-item-time {
  font-size: 0.75rem;
  color: var(--text-tertiary);
}

.empty-chat {
  text-align: center;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 0.9rem;
  color: var(--text-tertiary);
  margin-bottom: 20px;
}

.create-first-chat {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  transition: all 0.3s ease;
}

.create-first-chat:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-light);
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
  backdrop-filter: blur(2px);
}

/* 聊天内容区域 */
.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: var(--bg-primary);
  background-image: radial-gradient(#e5e7eb 1px, transparent 0);
  background-size: 20px 20px;
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
  max-width: 760px;
  margin: 0 auto;
}

/* 欢迎消息 */
.welcome-message {
  text-align: center;
  padding: 60px 20px;
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-content {
  background-color: white;
  padding: 40px 60px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.welcome-content:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.welcome-avatar {
  font-size: 60px;
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
  margin: 0 0 30px 0;
  line-height: 1.5;
}

.welcome-suggestions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.suggestion-btn {
  text-align: left;
  background-color: var(--bg-secondary);
  border-color: var(--border-color);
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.suggestion-btn:hover {
  background-color: var(--bg-primary);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

/* 消息列表 */
.message-list {
  padding: 20px 0;
}

.message-wrapper {
  margin-bottom: 24px;
}

.user-message, .ai-message {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.user-message {
  justify-content: flex-end;
}

.ai-message {
  justify-content: flex-start;
}

.message-avatar {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  background-color: var(--bg-secondary);
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s ease;
}

.message-avatar:hover {
  transform: scale(1.05);
}

.message-content {
  max-width: 70%;
  padding: 16px 20px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  line-height: 1.5;
  position: relative;
  word-wrap: break-word;
  transition: all 0.3s ease;
}

.user-content {
  background-color: var(--primary-color);
  color: white;
  border-bottom-right-radius: 4px;
}

.ai-content {
  background-color: white;
  color: var(--text-primary);
  border-bottom-left-radius: 4px;
}

.message-content:hover {
  box-shadow: var(--shadow-md);
}

.message-content.sending {
  opacity: 0.8;
}

.message-content.failed {
  opacity: 0.8;
  border: 1px solid #ff4d4f;
}

.message-content.loading {
  opacity: 0.8;
}

.message-content p {
    margin: 0;
    font-size: 0.95rem;
    line-height: 1.6;
    white-space: pre-wrap; /* 保留文本中的空白和换行符 */
    word-break: break-word; /* 确保长单词不会破坏布局 */
  }

.message-time {
  font-size: 0.75rem;
  opacity: 0.6;
  margin-left: 8px;
}

.message-status {
  position: absolute;
  top: 12px;
  right: 12px;
  color: #ff4d4f;
}

.loading-message {
  display: flex;
  align-items: center;
  color: var(--text-secondary);
  gap: 8px;
}

.loading-dots {
  display: flex;
  gap: 4px;
  align-items: center;
}

.loading-dots .dot {
  width: 6px;
  height: 6px;
  background-color: var(--text-secondary);
  border-radius: 50%;
  animation: loading 1.4s ease-in-out infinite;
}

.loading-dots .dot:nth-child(1) {
  animation-delay: 0s;
}

.loading-dots .dot:nth-child(2) {
  animation-delay: 0.2s;
}

.loading-dots .dot:nth-child(3) {
  animation-delay: 0.4s;
}

/* 输入区域 */
.chat-input-area {
  background-color: white;
  border-top: 1px solid var(--border-color);
  padding: 15px 20px;
  z-index: 10;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.05);
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  position: relative;
}

/* 工具按钮和菜单 */
.tool-button-container {
  position: relative;
}

.tool-button-wrapper {
  position: relative;
}

.tool-button {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  margin-right: 8px;
}

.tool-button:hover {
  background-color: var(--bg-secondary);
  color: var(--primary-color);
}

/* 快速提问区域 */
.quick-suggestions-container {
  margin: 10px 0;
  padding: 5px 0;
}

.quick-suggestions-scroll {
  display: flex;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 5px 0;
  scrollbar-width: thin;
  scrollbar-color: var(--border-color) transparent;
  gap: 8px;
}

.quick-suggestions-scroll::-webkit-scrollbar {
  height: 4px;
}

.quick-suggestions-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.quick-suggestions-scroll::-webkit-scrollbar-thumb {
  background-color: var(--border-color);
  border-radius: 20px;
}

/* border-color: var(--border-color); text-align: center; padding: 6px 12px; */
.quick-suggestion-btn {
  white-space: nowrap;
  background-color: white;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  flex-shrink: 0;
  justify-content: center;
}

.quick-suggestion-btn:hover {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
  transform: translateY(-1px);
}

.quick-suggestion-btn:active {
  transform: translateY(0);
}

/* 工具菜单样式 */

.tool-menu {
  position: absolute;
  bottom: 50px;
  left: 0;
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  min-width: 160px;
  z-index: 100;
}

.tool-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  font-size: 0.9rem;
  color: var(--text-primary);
}

.tool-menu-item:hover {
  background-color: var(--bg-secondary);
  color: var(--primary-color);
}

.tool-menu-item i {
  font-size: 16px;
}

/* 输入框 */
.chat-input {
  flex: 1;
  min-height: 44px;
  max-height: 120px;
  border-radius: var(--radius-full);
  border: 1px solid var(--border-color);
  padding: 10px 20px;
  font-size: 14px;
  resize: none;
  outline: none;
  transition: all 0.3s ease;
  background-color: var(--bg-secondary);
}

.chat-input:focus {
  border-color: var(--primary-color);
  background-color: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.chat-input:disabled {
  background-color: #f5f5f5;
  color: var(--text-tertiary);
}

/* 发送按钮 */
.send-button {
  height: 44px;
  width: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  transition: all 0.3s ease;
  font-size: 16px;
}

.send-button:hover:not(:disabled) {
  background-color: var(--primary-light);
  border-color: var(--primary-light);
  transform: scale(1.05);
}

.send-button:active:not(:disabled) {
  transform: scale(0.95);
}

.send-button:disabled {
  background-color: #d9d9d9;
  border-color: #d9d9d9;
  cursor: not-allowed;
}

.input-tip {
  text-align: center;
  font-size: 0.8rem;
  color: var(--text-tertiary);
  margin-top: 8px;
}

/* 动画效果 */
@keyframes loading {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.3;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 过渡效果 */
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from {
  transform: translateX(-100%);
}

.slide-leave-to {
  transform: translateX(-100%);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-up-enter-active,
.fade-up-leave-active {
  transition: all 0.4s ease;
}

.fade-up-enter-from,
.fade-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.message-fade-enter-active,
.message-fade-leave-active {
  transition: all 0.3s ease;
}

.message-fade-enter-from,
.message-fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-list-sidebar {
    width: 240px;
  }
  
  .message-content {
    max-width: 80%;
  }
  
  .welcome-content {
    padding: 30px 40px;
  }
  
  .welcome-suggestions {
    gap: 8px;
  }
  
  .suggestion-btn {
    padding: 10px 16px;
    font-size: 0.85rem;
  }
  
  /* 快速提问区域在平板视图的调整 */
  .quick-suggestions-container {
    margin: 8px 0;
  }
  
  .quick-suggestions-scroll {
    gap: 6px;
  }
  
  .quick-suggestion-btn {
    font-size: 0.85rem;
    padding: 6px 12px;
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
    padding: 15px 10px;
  }
  
  .message-content {
    max-width: 85%;
    padding: 12px 16px;
  }
  
  .chat-input-area {
    padding: 12px 15px;
  }
  
  .tool-button {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }
  
  .send-button {
    height: 36px;
    width: 36px;
    font-size: 14px;
  }
  
  .chat-input {
    min-height: 36px;
  }
  
  .welcome-content {
    padding: 20px 30px;
  }
  
  .welcome-avatar {
    font-size: 48px;
    margin-bottom: 16px;
  }
  
  .welcome-title {
    font-size: 1.3rem;
  }
  
  .welcome-desc {
    font-size: 0.9rem;
    margin-bottom: 20px;
  }
  
  /* 快速提问区域在手机视图的调整 */
  .quick-suggestions-container {
    margin: 6px 0;
    padding: 4px 0;
  }
  
  .quick-suggestions-scroll {
    gap: 6px;
    padding: 4px 0;
  }
  
  .quick-suggestion-btn {
    font-size: 0.8rem;
    padding: 5px 10px;
    border-radius: var(--radius-md);
  }
}
</style>