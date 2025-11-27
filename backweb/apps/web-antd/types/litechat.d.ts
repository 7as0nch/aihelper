/*
 * @Author: chengjiang
 * @Date: 2025-11-26 17:00:50
 * @Description:
 */
declare module '@7as0nch/litechat' {
  export function initAiChat(options: {
    config: Record<string, string>;
    defaultOpen?: boolean;
  }): void;
}
