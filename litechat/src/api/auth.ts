

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
export const authApi = {
    loginWithPhone: async (phone: string, _code: string): Promise<LoginResponse> => {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        return {
            token: 'mock-token-' + Date.now(),
            user: {
                id: '1',
                username: phone,
                avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${phone}`
            }
        };
    },

    loginWithPassword: async (username: string, password: string): Promise<LoginResponse> => {
        // Simulate API call
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
    },

    register: async (username: string, _password: string): Promise<LoginResponse> => {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000));

        return {
            token: 'mock-token-' + Date.now(),
            user: {
                id: '1',
                username: username,
                avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`
            }
        };
    },

    getAuthConfig: async (): Promise<AuthConfig> => {
        // In a real app, this might come from the backend
        // For now, we use the env var, but wrapped in a promise to simulate async config loading
        return {
            enableQrLogin: import.meta.env.VITE_ENABLE_QR_LOGIN === 'true'
        };
    }
};
