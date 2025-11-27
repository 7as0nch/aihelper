import axios from 'axios';
// import { getConfig } from '@/config';

import { getToken } from '@/utils/cookie';
import { message } from 'ant-design-vue';
import { useAuthStore } from '@/stores/auth';

export const request = axios.create({
    // Use /api to trigger Vite proxy in dev mode
    baseURL: '/api',
    timeout: 60000,
});

interface ResultBody {
    code: number;
    msg: string;
    data: any;
    newToken?: string;
}

request.interceptors.request.use(
    (config) => {
        const token = getToken();
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

request.interceptors.response.use(
    (response) => {
        const res = response.data as ResultBody;

        // // Handle token refresh
        // if (res.newToken) {
        //     setToken(res.newToken);
        // }

        // Check business code
        if (res.code === 200) {
            return res.data;
        }

        // Handle Auth Error
        if (res.code === 401) {
            const authStore = useAuthStore();
            authStore.logout();
            authStore.openModal();
            message.error('登录已过期，请重新登录');
            return Promise.reject(new Error(res.msg || 'Unauthorized'));
        }

        // Handle other errors
        message.error(res.msg || '请求失败');
        return Promise.reject(new Error(res.msg || 'Error'));
    },
    (error) => {
        // Handle HTTP errors
        if (error.response && error.response.status === 401) {
            const authStore = useAuthStore();
            authStore.logout();
            authStore.openModal();
            message.error('登录已过期，请重新登录');
        } else {
            message.error(error.message || '网络错误');
        }
        return Promise.reject(error);
    }
);

// Wrapper to provide generic type support
const apiRequest = <T = any>(config: any): Promise<T> => {
    return request(config) as Promise<T>;
};

export default apiRequest;
