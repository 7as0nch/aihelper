/*
 * @Author: chengjiang
 * @Date: 2025-11-28 11:47:37
 * @Description: 初始化 AI 助手
 */
import { defineStore } from 'pinia';

import { initAiChat } from '@7as0nch/litechat';

  /**
   * 初始化 AI 助手
   * Initialize the AI assistant
   * @param options 初始化选项
   */
export const useAiStore = defineStore('ai', () => {
  const aiChat = initAiChat({
    config: {
      VITE_APP_TITLE: '我的 AI 助手',
      VITE_API_BASE_URL: 'https://api.example.com',
      VITE_AI_TYPE: "demo",
    },
    defaultOpen: false,
    defaultShow: false,
  });

  /**
   * 显示 AI 助手
   * Show the AI assistant
   */
  function show() {
    aiChat.show();
  }

  /**
   * 隐藏 AI 助手
   * Hide the AI assistant
   */
  function hide() {
    aiChat.hide();
  }

  return {
    show,
    hide,
  };
});
