import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi, type AuthConfig } from '../api/auth';

export interface User {
    id: string;
    username: string;
    role?: string; // user, admin: default: user.
    avatar?: string;
}

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
        if (isAuthenticated.value) {
            return true;
        }
        openModal();
        return false;
    };

    const handleLoginSuccess = (userData: User, token: string) => {
        user.value = userData;
        isAuthenticated.value = true;
        localStorage.setItem('user', JSON.stringify(userData));
        localStorage.setItem('token', token);
        closeModal();
    };

    const loginWithPhone = async (phone: string, code: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPhone(phone, code);
            handleLoginSuccess(response.user, response.token);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const loginWithPassword = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPassword(username, password);
            handleLoginSuccess(response.user, response.token);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const register = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.register(username, password);
            handleLoginSuccess(response.user, response.token);
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
        localStorage.removeItem('token');
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
