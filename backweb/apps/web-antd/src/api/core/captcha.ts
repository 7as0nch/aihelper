import { useAppConfig } from '@vben/hooks';

import { requestClient } from '#/api/request';

/**
 * 发送短信验证码
 * @param phonenumber 手机号
 * @returns void
 */
export function sendSmsCode(phonenumber: string) {
  return requestClient.get<void>('/resource/sms/code', {
    params: { phonenumber },
  });
}

/**
 * 发送邮件验证码
 * @param email 邮箱
 * @returns void
 */
export function sendEmailCode(email: string) {
  return requestClient.get<void>('/resource/email/code', {
    params: { email },
  });
}

/**
 * @param img 图片验证码 需要和base64拼接
 * @param captchaEnabled 是否开启
 * @param uuid 验证码ID
 */
export interface CaptchaResponse {
  captchaEnabled: boolean;
  img: string;
  uuid: string;
}

const { sseEnable } = useAppConfig(import.meta.env, import.meta.env.PROD);
/**
 * 图片验证码
 * @returns resp
 */
export function captchaImage() {
  if (!sseEnable) {
    return;
  }
  return requestClient.get<CaptchaResponse>('/auth/code');
}
