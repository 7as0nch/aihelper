import axios from 'axios';

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
  (config) => {
    // 从本地存储获取token并添加到请求头
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    // 统一处理响应数据
    return response.data;
  },
  (error) => {
    // 统一处理错误
    let errorMessage = '请求失败，请稍后重试';
    if (error.response) {
      // 服务器返回错误
      switch (error.response.status) {
        case 401:
          errorMessage = '请先登录';
          // 清除token并重定向到登录页
          localStorage.removeItem('token');
          if (window.location.pathname.startsWith('/admin')) {
            window.location.href = '/admin/login';
          } else {
            window.location.href = '/login';
          }
          break;
        case 403:
          errorMessage = '没有权限访问';
          break;
        case 404:
          errorMessage = '请求的资源不存在';
          break;
        case 500:
          errorMessage = '服务器错误，请稍后重试';
          break;
        default:
          errorMessage = error.response.data.message || errorMessage;
      }
    } else if (error.request) {
      // 请求发出但没有收到响应
      errorMessage = '网络错误，请检查网络连接';
    }
    
    // 显示错误提示
    console.error('API请求错误:', errorMessage);
    return Promise.reject(new Error(errorMessage));
  }
);

// 登录参数接口
interface LoginParams {
  username: string;
  password: string;
}

// 登录响应接口
interface LoginResponse {
  code: number;
  message: string;
  data: {
    token: string;
    userInfo: {
      id: number;
      username: string;
      email: string;
      role: string;
      status: number;
    }
  };
}

// 登录API
export const login = (params: LoginParams): Promise<LoginResponse> => {
  return request.post('/users/login', params);
};

// 注册参数接口
interface RegisterParams {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

// 注册响应接口
interface RegisterResponse {
  code: number;
  message: string;
  data?: any;
}

// 注册API
export const register = (params: RegisterParams): Promise<RegisterResponse> => {
  return request.post('/users/register', params);
};

// 用户信息接口
interface UserInfo {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  role: string;
  createTime: string;
  updateTime: string;
  status: number;
}

// 获取用户信息响应接口
interface GetUserInfoResponse {
  code: number;
  message: string;
  data: UserInfo;
}

// 获取用户信息API
export const getUserInfo = (): Promise<GetUserInfoResponse> => {
  return request.get('/users/info');
};

// 更新用户信息参数接口
interface UpdateUserInfoParams {
  nickname?: string;
  email?: string;
  avatar?: string;
}

// 更新用户信息响应接口
interface UpdateUserInfoResponse {
  code: number;
  message: string;
  data: UserInfo;
}

// 更新用户信息API
export const updateUserInfo = (params: UpdateUserInfoParams): Promise<UpdateUserInfoResponse> => {
  return request.put('/users/info', params);
};

// 管理员登录API
export const adminLogin = (params: LoginParams): Promise<LoginResponse> => {
  return request.post('/admin/login', params);
};

// 退出登录API
export const logout = (): Promise<any> => {
  return request.post('/users/logout');
};

// 管理员退出登录API
export const adminLogout = (): Promise<any> => {
  return request.post('/admin/logout');
};

export default request;