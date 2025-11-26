import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi, type AuthConfig, type UserInfo as User } from '../api/auth';
import { getConfig } from '@/config';
import { setToken, removeToken } from '@/utils/cookie';

export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null);
    const isAuthenticated = ref(false);
    const showAuthModal = ref(false);
    const config = ref<AuthConfig>({ enableQrLogin: true });

    // Initialize from localStorage and load config
    const init = async () => {
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            user.value = JSON.parse(storedUser);
            isAuthenticated.value = true;
        }

        try {
            config.value = await authApi.getAuthConfig();
        } catch (e) {
            console.error('Failed to load auth config', e);
        }
    };

    const openModal = () => {
        showAuthModal.value = true;
    };

    const closeModal = () => {
        showAuthModal.value = false;
    };

    const checkAuth = (): boolean => {
        // 前端模式不需要登录
        if (isAuthenticated.value || getConfig('VITE_AI_TYPE') === 'frontend') {
            return true;
        }
        openModal();
        return false;
    };

    const handleLoginSuccess = async (userData: User, token: string) => {
        if (!userData) {
            try {
                const res = await authApi.getUserInfo();
                userData = res.user;
            } catch (e) {
                console.error('Failed to load user info', e);
            }
        }
        user.value = userData;
        isAuthenticated.value = true;
        localStorage.setItem('user', JSON.stringify(userData));
        setToken(token);
        closeModal();
    };

    const loginWithPhone = async (phone: string, code: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPhone(phone, code);
            handleLoginSuccess(response.user, response.accessToken);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const loginWithPassword = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPassword(username, password);
            console.log(response);
            handleLoginSuccess(response.user, response.accessToken);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const register = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.register(username, password);
            handleLoginSuccess(response.user, response.accessToken);
            return true;
        } catch (e) {
            console.error('Registration failed', e);
            return false;
        }
    };

    const logout = () => {
        user.value = null;
        isAuthenticated.value = false;
        localStorage.removeItem('user');
        removeToken();
    };

    init();

    return {
        user,
        isAuthenticated,
        showAuthModal,
        config,
        openModal,
        closeModal,
        checkAuth,
        loginWithPhone,
        loginWithPassword,
        register,
        logout
    };
});
