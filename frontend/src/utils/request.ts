import axios from 'axios';
import { useUserStore } from '../stores/user';

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
request.interceptors.request.use(
  config => {
    const userStore = useUserStore();
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          const userStore = useUserStore();
          userStore.logout();
          break;
        case 403:
          alert('没有权限访问');
          break;
        case 500:
          alert('服务器错误');
          break;
        default:
          alert(error.response.data.message || '请求失败');
      }
    }
    return Promise.reject(error);
  }
);

export default request;