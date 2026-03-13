import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi, type AuthConfig, type UserInfo as User } from '../api/auth';
import { getConfig } from '@/config';
import { getToken, setToken, removeToken } from '@/utils/cookie';

export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null);
    const isAuthenticated = ref(false);
    const showAuthModal = ref(false);
    const config = ref<AuthConfig>({ enableQrLogin: true });

    const consumeQQTokenFromUrl = () => {
        if (typeof window === 'undefined') return;

        const current = new URL(window.location.href);
        const qqToken = current.searchParams.get('qq_token');
        const qqError = current.searchParams.get('qq_error');

        if (qqToken) {
            setToken(qqToken);
        }

        if (qqToken || qqError) {
            current.searchParams.delete('qq_token');
            current.searchParams.delete('qq_error');
            const query = current.searchParams.toString();
            const cleanURL = `${current.pathname}${query ? `?${query}` : ''}${current.hash}`;
            window.history.replaceState({}, '', cleanURL);
        }
    };

    // Initialize from localStorage and load config
    const init = async () => {
        consumeQQTokenFromUrl();

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

        // If token exists but user profile is not loaded yet, recover session from backend.
        if (!isAuthenticated.value && getToken()) {
            try {
                const res = await authApi.getUserInfo();
                user.value = res.user;
                isAuthenticated.value = true;
                localStorage.setItem('user', JSON.stringify(res.user));
            } catch (e) {
                console.error('Failed to recover auth session', e);
                logout();
            }
        }
    };

    const openModal = () => {
        showAuthModal.value = true;
    };

    const closeModal = () => {
        showAuthModal.value = false;
    };

    const checkAuth = (): boolean => {
        // Frontend mode does not require login.
        if (isAuthenticated.value || getConfig('VITE_AI_TYPE') === 'frontend') {
            return true;
        }
        openModal();
        return false;
    };

    const handleLoginSuccess = async (userData: User, token: string) => {
        setToken(token);
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
        closeModal();
    };

    const loginWithPhone = async (phone: string, code: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPhone(phone, code);
            await handleLoginSuccess(response.user, response.accessToken);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const loginWithPassword = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.loginWithPassword(username, password);
            await handleLoginSuccess(response.user, response.accessToken);
            return true;
        } catch (e) {
            console.error('Login failed', e);
            return false;
        }
    };

    const register = async (username: string, password: string): Promise<boolean> => {
        try {
            const response = await authApi.register(username, password);
            await handleLoginSuccess(response.user, response.accessToken);
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

