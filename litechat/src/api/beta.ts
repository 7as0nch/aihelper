import { getConfig } from '@/config';
import request from '@/utils/request';

export interface BetaApplicationPayload {
  productInterest: 'litechat' | 'aicook' | 'tech-sandbox';
  contactType: 'email' | 'qq';
  contactValue: string;
  useCase: string;
  note?: string;
  sourcePage?: string;
}

export interface BetaApplicationResponse {
  id: number;
  status: string;
  mailStatus: string;
}

export const betaApi = {
  async submitApplication(payload: BetaApplicationPayload): Promise<BetaApplicationResponse> {
    const aiType = getConfig('VITE_AI_TYPE');

    if (aiType === 'backend') {
      return request<BetaApplicationResponse>({
        url: '/beta/applications',
        method: 'post',
        data: payload,
      });
    }

    await new Promise((resolve) => window.setTimeout(resolve, 800));
    return {
      id: Date.now(),
      status: 'submitted',
      mailStatus: 'skipped',
    };
  },
};
