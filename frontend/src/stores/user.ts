import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from 'axios';
import router from '../router';

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

// 登录参数接口
interface LoginParams {
  username: string;
  password: string;
}

// 注册参数接口
interface RegisterParams {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

// 更新用户信息参数接口
interface UpdateUserInfoParams {
  nickname?: string;
  avatar?: string;
  email?: string;
}

export const useUserStore = defineStore('user', () => {
  // 用户信息状态
  const userInfo = ref<UserInfo | null>(null);
  const token = ref<string>('');
  const isLoggedIn = ref<boolean>(false);
  const loading = ref<boolean>(false);

  // 计算属性 - 是否为管理员
  const isAdmin = computed(() => {
    return userInfo.value?.role === 'admin';
  });

  // 计算属性 - 用户基本信息摘要
  const userSummary = computed(() => {
    if (!userInfo.value) {
      return null;
    }
    return {
      nickname: userInfo.value.nickname || userInfo.value.username,
      avatar: userInfo.value.avatar || 'https://via.placeholder.com/40',
      role: userInfo.value.role === 'admin' ? '管理员' : '普通用户'
    };
  });

  // 用户登录
  const login = async (params: LoginParams) => {
    loading.value = true;
    try {
      // const response = await axios.post('/api/users/login', params);
      const response = {data: {token: '123456', message: '登录成功'}}
      if (response.data.token) {
        token.value = response.data.token;
        localStorage.setItem('token', response.data.token);
        // await getUserInfo();
        userInfo.value = {
          id: 1,
          username: params.username,
          email: '<EMAIL></EMAIL>',
          nickname: params.username,
          avatar: 'https://via.placeholder.com/40',
          role: 'admin',
          createTime: new Date().toISOString(),
          updateTime: new Date().toISOString(),
          status: 1
        };
        isLoggedIn.value = true;
        // 存储用户角色信息
        if (userInfo.value) {
          localStorage.setItem('userRole', userInfo.value.role);
        }
        return { success: true, data: response.data };
      }
      return { success: false, message: response.data.message || '登录失败' };
    } catch (error) {
      console.error('登录失败:', error);
      return { success: false, message: '登录失败，请重试' };
    } finally {
      loading.value = false;
    }
  };

  // 用户注册
  const register = async (params: RegisterParams) => {
    loading.value = true;
    try {
      const response = await axios.post('/api/users/register', params);
      if (response.data.success) {
        return { success: true };
      }
      return { success: false, message: response.data.message || '注册失败' };
    } catch (error) {
      console.error('注册失败:', error);
      return { success: false, message: '注册失败，请重试' };
    } finally {
      loading.value = false;
    }
  };

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const response = await axios.get('/api/users/info');
      if (response.data) {
        userInfo.value = response.data;
        isLoggedIn.value = true;
      }
    } catch (error) {
      console.error('获取用户信息失败:', error);
      logout();
    }
  };

  // 更新用户信息
  const updateUserInfo = async (params: UpdateUserInfoParams) => {
    loading.value = true;
    try {
      const response = await axios.put('/api/users/info', params);
      if (response.data.success && response.data.data) {
        // 更新本地用户信息
        userInfo.value = {
          ...userInfo.value,
          ...response.data.data
        } as UserInfo;
        return { success: true, data: response.data.data };
      }
      return { success: false, message: response.data.message || '更新用户信息失败' };
    } catch (error) {
      console.error('更新用户信息失败:', error);
      return { success: false, message: '更新用户信息失败，请重试' };
    } finally {
      loading.value = false;
    }
  };

  // 退出登录
  const logout = () => {
    userInfo.value = null;
    token.value = '';
    isLoggedIn.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('userRole');
    router.push('/login');
  };

  // 管理员登出
  const adminLogout = () => {
    userInfo.value = null;
    token.value = '';
    isLoggedIn.value = false;
    localStorage.removeItem('token');
    localStorage.removeItem('userRole');
    router.push('/admin/login');
  };

  // 初始化
  const initialize = () => {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
      token.value = savedToken;
      getUserInfo();
    }
  };

  return {
    // 状态
    userInfo,
    token,
    isLoggedIn,
    loading,
    // 计算属性
    isAdmin,
    userSummary,
    // 方法
    login,
    register,
    logout,
    adminLogout,
    getUserInfo,
    updateUserInfo,
    initialize
  };
});