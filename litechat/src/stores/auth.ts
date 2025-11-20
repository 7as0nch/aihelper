import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface User {
    id: string;
    username: string;
    avatar?: string;
}

export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null);
    const isAuthenticated = ref(false);

    const showAuthModal = ref(false);

    // Initialize from localStorage
    const init = () => {
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            user.value = JSON.parse(storedUser);
            isAuthenticated.value = true;
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

    const login = async (username: string, password: string): Promise<boolean> => {
        // Mock API call
        await new Promise(resolve => setTimeout(resolve, 1000));

        if (username && password) {
            const mockUser: User = {
                id: '1',
                username: username,
                avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
            };

            user.value = mockUser;
            isAuthenticated.value = true;
            localStorage.setItem('user', JSON.stringify(mockUser));
            closeModal();
            return true;
        }
        return false;
    };

    const register = async (username: string, password: string): Promise<boolean> => {
        // Mock API call
        await new Promise(resolve => setTimeout(resolve, 1000));

        if (username && password) {
            // Auto login after register
            return await login(username, password);
        }
        return false;
    };

    const logout = () => {
        user.value = null;
        isAuthenticated.value = false;
        localStorage.removeItem('user');
        // We don't redirect here, let the component handle it or use a global guard
    };

    init();

    return {
        user,
        isAuthenticated,
        showAuthModal,
        openModal,
        closeModal,
        checkAuth,
        login,
        register,
        logout
    };
});
