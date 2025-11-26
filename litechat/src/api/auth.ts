
import { getConfig } from '@/config';

export interface UserInfo {
    id: string;
    userName: string;
    avatar?: string;
    roles?: string[];
}

export interface GetUserResponse {
    user: UserInfo;
}

export interface LoginResponse {
    accessToken: string;
    user: UserInfo;
}

export interface AuthConfig {
    enableQrLogin: boolean;
}

// Define the type for authApi
export interface AuthApi {
    loginWithPhone: (phone: string, code: string) => Promise<LoginResponse>;
    loginWithPassword: (username: string, password: string) => Promise<LoginResponse>;
    register: (username: string, password: string) => Promise<LoginResponse>;
    getAuthConfig: () => Promise<AuthConfig>;
    getUserInfo: () => Promise<GetUserResponse>;
    getLoginMethods: () => Promise<{ enablePasswordLogin: boolean; enableQrLogin: boolean }>;
}

// Mock API implementation for now
import request from '@/utils/request';

// Mock API implementation for now
export const authApi: AuthApi = {
    loginWithPhone: async (phone: string, code: string): Promise<LoginResponse> => {
        const aiType = getConfig('VITE_AI_TYPE');

        // Demo mode: Mock login
        if (aiType === 'demo') {
            await new Promise(resolve => setTimeout(resolve, 1000));
            return {
                accessToken: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    userName: phone,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${phone}`
                }
            };
        }

        // Backend mode (Backend API): Real login
        if (aiType === 'backend') {
            return request<LoginResponse>({
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
                accessToken: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    userName: username,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
                }
            };
        }

        // Backend mode (Backend API): Real login
        if (aiType === 'backend') {
            return request<LoginResponse>({
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
                accessToken: 'mock-token-' + Date.now(),
                user: {
                    id: '1',
                    userName: username,
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
                }
            };
        }

        // Backend mode (Backend API): Real register
        if (aiType === 'backend') {
            return request<LoginResponse>({
                url: '/auth/register',
                method: 'post',
                data: { username, password }
            });
        }

        throw new Error('Register not supported in pure frontend mode');
    },

    getAuthConfig: async (): Promise<AuthConfig> => {
        // const aiType = getConfig('VITE_AI_TYPE');

        // if (aiType === 'backend') {
        //     // Try to get config from backend, fallback to env
        //     try {
        //         const res = await request<AuthConfig>({ url: '/auth/config', method: 'get' });
        //         return res;
        //     } catch (e) {
        //         console.warn('Failed to fetch auth config from backend, using env');
        //     }
        // }

        return {
            enableQrLogin: getConfig('VITE_ENABLE_QR_LOGIN') === 'true'
        };
    },

    getUserInfo: async (): Promise<GetUserResponse> => {
        const aiType = getConfig('VITE_AI_TYPE');

        // Demo mode: Mock user info
        if (aiType === 'demo') {
            await new Promise(resolve => setTimeout(resolve, 1000));
            return {
                user: {
                    id: '1',
                    userName: 'demo',
                    avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${'demo'}`
                }
            };
        }

        // Backend mode (Backend API): Real user info
        if (aiType === 'backend') {
            return request<GetUserResponse>({
                url: '/auth/user/info',
                method: 'get'
            });
        }

        throw new Error('User info not supported in pure frontend mode');
    },

    async getLoginMethods() {
        return {
            enablePasswordLogin: true,
            enableQrLogin: getConfig('VITE_ENABLE_QR_LOGIN') === 'true'
        };
    }
};
