
import { getConfig } from '@/config';

export interface LoginResponse {
    token: string;
    user: {
        id: string;
        username: string;
        avatar?: string;
    };
}

export interface AuthConfig {
    enableQrLogin: boolean;
}

// Mock API implementation for now
import request from '@/utils/request';

// Mock API implementation for now
export const authApi = {
    loginWithPhone: async (phone: string, code: string): Promise<LoginResponse> => {
        const aiType = getConfig('VITE_AI_TYPE');

        // Demo mode: Mock login
        if (aiType === 'demo') {
            await new Promise(resolve => setTimeout(resolve, 1000));
            return {
                token: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    username: phone,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${phone}`
                }
            };
        }

        // Frontend mode (Backend API): Real login
        if (aiType === 'frontend') {
            return request({
                url: '/auth/login/phone',
                method: 'post',
                data: { phone, code }
            });
        }

        // Backend mode (Pure Frontend): No login usually, but if called, maybe mock or error
        throw new Error('Login not supported in pure frontend mode');
    },

    loginWithPassword: async (username: string, password: string): Promise<LoginResponse> => {
        const aiType = getConfig('VITE_AI_TYPE');

        // Demo mode: Mock login
        if (aiType === 'demo') {
            await new Promise(resolve => setTimeout(resolve, 1000));
            if (password === 'error') throw new Error('Invalid credentials');
            return {
                token: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    username: username,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
                }
            };
        }

        // Frontend mode (Backend API): Real login
        if (aiType === 'frontend') {
            return request({
                url: '/auth/login',
                method: 'post',
                data: { username, password }
            });
        }

        throw new Error('Login not supported in pure frontend mode');
    },

    register: async (username: string, password: string): Promise<LoginResponse> => {
        const aiType = getConfig('VITE_AI_TYPE');

        // Demo mode: Mock register
        if (aiType === 'demo') {
            await new Promise(resolve => setTimeout(resolve, 1000));
            return {
                token: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    username: username,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
                }
            };
        }

        // Frontend mode (Backend API): Real register
        if (aiType === 'frontend') {
            return request({
                url: '/auth/register',
                method: 'post',
                data: { username, password }
            });
        }

        throw new Error('Register not supported in pure frontend mode');
    },

    getAuthConfig: async (): Promise<AuthConfig> => {
        const aiType = getConfig('VITE_AI_TYPE');

        if (aiType === 'frontend') {
            // Try to get config from backend, fallback to env
            try {
                const res = await request({ url: '/auth/config', method: 'get' });
                return res as unknown as AuthConfig;
            } catch (e) {
                console.warn('Failed to fetch auth config from backend, using env');
            }
        }

        return {
            enableQrLogin: getConfig('VITE_ENABLE_QR_LOGIN') === 'true'
        };
    },

    async getLoginMethods() {
        return {
            enablePasswordLogin: true,
            enableQrLogin: getConfig('VITE_ENABLE_QR_LOGIN') === 'true'
        };
    }
};
